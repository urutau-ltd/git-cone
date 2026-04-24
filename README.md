# git-cone

`git-cone` is a security-hardened hard fork of `soft-serve`.
It is intended to stay as drop-in compatible as practical while focusing on
security fixes and operational hardening, not on growing the core product
surface.

`cone` is the primary CLI. `soft` remains available as a compatibility layer.

## Highlights

- Pure Go build, including SQLite via `modernc.org/sqlite`
- SSH TUI and SSH command interface
- HTTP Git smart protocol and LFS
- Dual env-prefix support: `GIT_CONE_*` overrides `SOFT_SERVE_*`
- Optional Gotify notifications for security-relevant events
- Strict mode for hardened deployments
- `git://` disabled by default

## Quick Start

Build locally:

```bash
make build
./dist/cone serve
```

Or enter the Guix development shell first:

```bash
make shell
make build
make test
```

## Docker

```bash
docker pull ghcr.io/urutau-ltd/git-cone:latest
```

Minimal Compose example:

```yaml
services:
  git-cone:
    image: ghcr.io/urutau-ltd/git-cone:latest
    ports:
      - "23231:23231"
      - "127.0.0.1:23232:23232"
    volumes:
      - git-cone-data:/git-cone/data
      - ./git-cone/hooks:/git-cone/data/hooks
    environment:
      - GIT_CONE_DATA_PATH=/git-cone/data
      - GIT_CONE_INITIAL_ADMIN_KEYS=ssh-ed25519 AAAA...
      - GIT_CONE_SSH_PUBLIC_URL=ssh://git.example.com
      - GIT_CONE_HTTP_PUBLIC_URL=https://git.example.com
      - GIT_CONE_NAME=Git Cone
      - GIT_CONE_SECURITY_STRICT=true
    restart: unless-stopped

volumes:
  git-cone-data:
```

## Compatibility

This fork aims to remain a practical drop-in replacement for recent
`soft-serve` deployments.

| Before | After |
|---|---|
| `soft serve` | `cone serve` or `soft serve` |
| `soft browse` | `cone browse` or `soft browse` |
| `SOFT_SERVE_*` | `GIT_CONE_*` preferred, `SOFT_SERVE_*` still supported |

For existing Compose stacks, the least disruptive migration is:

- keep the service name as `soft-serve`
- keep the volume name as `soft-serve-data`
- switch the image to `ghcr.io/urutau-ltd/git-cone:latest`
- mount that volume at `/git-cone/data`
- set `GIT_CONE_DATA_PATH=/git-cone/data`

## Authentication and Strict Mode

`strict=true` does not disable token-based HTTP access.
It does this:

- forces anonymous access to `no-access`
- disables keyless access
- clamps SSH timeouts

That means:

- SSH always requires an authorized key
- HTTP Git/LFS requires valid credentials
- access tokens still work for automation such as `pipe`

With `strict=true`, a non-private repo is not anonymously readable.
Today there is no per-repo “public override” when global anonymous access is forced off.

## Pipe

For `pipe`, use an access token and internal HTTP:

```text
http://pipe-bot:${PIPE_GIT_TOKEN}@soft-serve:23232
```

Recommended setup:

1. Create a `pipe-bot` user.
2. Grant it read-only or admin access to the repos it needs.
3. Create a token with `ssh cone token create`.
4. Store that token in your deployment env as `PIPE_GIT_TOKEN`.

If you also run Gotify, let `pipe` talk to `http://gotify:80` directly on the
internal network. Do not put internal service-to-service traffic through Anubis.

## Hooks

Global hooks live under `$GIT_CONE_DATA_PATH/hooks/`.
The server writes a sample `update.sample` hook on first start.

## Configuration

The generated `config.yaml` is the canonical reference for the settings this
fork currently exposes.

Notable defaults:

```yaml
git:
  listen_addr: ""   # git:// disabled by default

security:
  strict: false

notify:
  gotify:
    enabled: false
```

## Development

This repo is Guix-friendly and includes a maintained `manifest.scm`.

```bash
make shell
make build
make test
make test-all
make image
```

Files worth knowing:

- `manifest.scm`: Guix dev environment
- `Makefile`: common dev and CI entrypoints
- `.pipe.yml`: project pipeline

## Service Managers

The repo no longer ships systemd-centric packaging inherited from upstream.
If you need host-managed services, optional examples live in
[docs/service-managers.md](docs/service-managers.md) for:

- SysVinit
- OpenRC
- Runit
- GNU Shepherd

## License

[MIT](LICENSE)
