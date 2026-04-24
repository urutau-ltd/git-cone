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
- `cone audit`, `repo verify`, and `/health`

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

First-time SSH admin flow:

1. Start the server with `GIT_CONE_INITIAL_ADMIN_KEYS` pointing to your public key.
2. Connect with `ssh -p 23231 git@host` for the TUI.
3. Run admin commands over SSH, for example:
   `ssh -p 23231 host user create alice`
   `ssh -p 23231 host repo create demo`
   `ssh -p 23231 host audit`

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

Container notes:

- data lives at `/git-cone/data`
- hooks live at `/git-cone/data/hooks`
- the image provides both `cone` and `soft`
- `/health` is intended for local container health checks such as Docker/Dozzle

## Compatibility

This fork aims to remain a practical drop-in replacement for recent
`soft-serve` deployments.

| Before | After |
|---|---|
| `soft serve` | `cone serve` or `soft serve` |
| `soft browse` | `cone browse` or `soft browse` |
| `SOFT_SERVE_*` | `GIT_CONE_*` preferred, `SOFT_SERVE_*` still supported |

What changed on purpose:

- the preferred binary name is now `cone`
- `soft` still works and maps to the same implementation
- the default server name is `Git Cone`
- the image stores data in `/git-cone/data`
- `git://` is off by default
- the server is SQLite-only

For existing Compose stacks, the least disruptive migration is:

- keep the service name as `soft-serve`
- keep the volume name as `soft-serve-data`
- switch the image to `ghcr.io/urutau-ltd/git-cone:latest`
- mount that volume at `/git-cone/data`
- set `GIT_CONE_DATA_PATH=/git-cone/data`

Example drop-in replacement:

```yaml
services:
  soft-serve:
    image: ghcr.io/urutau-ltd/git-cone:latest
    ports:
      - "23231:23231"
      - "23232:23232"
    volumes:
      - soft-serve-data:/git-cone/data
      - ./soft-serve/hooks:/git-cone/data/hooks
    environment:
      - GIT_CONE_DATA_PATH=/git-cone/data
      - GIT_CONE_INITIAL_ADMIN_KEYS=${SOFT_SERVE_ADMIN_KEY}
      - GIT_CONE_SSH_PUBLIC_URL=ssh://git.example.com
      - GIT_CONE_HTTP_PUBLIC_URL=https://git.example.com
      - GIT_CONE_NAME=Git Cone
      - GIT_CONE_SECURITY_STRICT=true
    healthcheck:
      test: ["CMD", "curl", "-fsS", "http://127.0.0.1:23232/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 15s
```

## Commands

The SSH command interface is still the main control plane.

Common commands:

- `ssh -p 23231 host info`
- `ssh -p 23231 host user create <name>`
- `ssh -p 23231 host repo create <repo>`
- `ssh -p 23231 host repo list`
- `ssh -p 23231 host token create "<label>"`
- `ssh -p 23231 host pubkey add`

Fork-specific or newly documented commands:

- `ssh -p 23231 host audit`
  Prints server version, current user, active key fingerprint, optional key age,
  negotiated SSH details when available, and repo counts.
- `ssh -p 23231 host repo verify <repo>`
  Runs `git fsck --full` on a repository. Intended for admins and users with write access.

CLI entrypoints:

- `cone serve`: start the server
- `cone browse`: browse repositories locally
- `soft serve`: compatibility alias
- `soft browse`: compatibility alias

The Git transport URLs and repository workflow stay the same:

- SSH clone/push: `ssh://host:23231/<repo>.git`
- HTTP clone/push: `https://host/<repo>.git`

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

This is intentional hardening. If you need anonymous read access for non-private
repos, do not enable strict mode.

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
db:
  driver: "sqlite"

git:
  listen_addr: ""   # git:// disabled by default

security:
  strict: false

notify:
  gotify:
    enabled: false
```

Environment prefixes:

- `GIT_CONE_*` is preferred
- `SOFT_SERVE_*` remains supported for compatibility
- if both are set, `GIT_CONE_*` wins

Useful variables:

- `GIT_CONE_DATA_PATH`
- `GIT_CONE_INITIAL_ADMIN_KEYS`
- `GIT_CONE_SSH_PUBLIC_URL`
- `GIT_CONE_HTTP_PUBLIC_URL`
- `GIT_CONE_NAME`
- `GIT_CONE_SECURITY_STRICT`
- `GIT_CONE_NOTIFY_GOTIFY_ENABLED`
- `GIT_CONE_NOTIFY_GOTIFY_URL`
- `GIT_CONE_NOTIFY_GOTIFY_TOKEN`

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
- `internal/cli/cli.go`: shared CLI wiring for `cone` and `soft`
- `pkg/config/config.go`: defaults and env loading
- `pkg/ssh/cmd/`: SSH command handlers

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
