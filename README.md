# 🍦 git-cone

`git-cone` is a security-hardened hard fork of
[`soft-serve`](https://github.com/charmbracelet/soft-serve). It is intended to
stay as drop-in compatible as practical while focusing on security fixes and
operational hardening, not on growing the core product surface.

`cone` is the primary CLI. `soft` remains available as a compatibility layer.

This fork includes security work imported and adapted from
[dvrd](https://github.com/dvrd)'s `soft-serve` branch, including the patch set
described in local commit `8eb4b04` ("apply dvrd rounds 46-59 fixes"). That
import covered fixes such as SSRF/JWT hardening and several backend/store
corrections.

Additional hardening was then implemented directly in `git-cone`:

- SSH was hardened using stricter KEX, cipher, and MAC defaults in
  [pkg/ssh/ssh.go](pkg/ssh/ssh.go).
- SSH stdin was hardened using input rate limiting in
  [pkg/ssh/middleware.go](pkg/ssh/middleware.go).
- File serving was hardened by removing the `sendFile` TOCTOU window in
  [pkg/web/git.go](pkg/web/git.go).
- User deletion was hardened by fixing repo and row deletion ordering in
  [pkg/backend/user.go](pkg/backend/user.go).

As of `v0.13.X` the internals between `soft-serve` and `git-cone` have diverged
a bit more, please don't take this list as exhaustive.

## Highlights

- Pure Go build, including SQLite via `modernc.org/sqlite`
- SSH TUI and SSH command interface
- HTTP Git smart protocol and LFS
- Dual env-prefix support: `GIT_CONE_*` overrides `SOFT_SERVE_*`
- Optional Gotify notifications for security-relevant events
- Strict mode for hardened deployments
- `git://` disabled by default
- `cone audit`, `repo verify`, and `/health`

## Screenshots

![screenshot1](.repo-assets/git-cone-screenshot.png)
![screenshot2](.repo-assets/git-cone-screenshot2.png)

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

1. Start the server with `GIT_CONE_INITIAL_ADMIN_KEYS` pointing to your public
   key.
2. Connect with `ssh -p 23231 git@host` for the TUI.
3. Run admin commands over SSH, for example:
   `ssh -p 23231 host user create alice` `ssh -p 23231 host repo create demo`
   `ssh -p 23231 host audit`

## Docker

> [!IMPORTANT]
> The `latest` tag is literally the latest image built wheter it was tagged
> not, this includes development builds. Use a pinned version.

```bash
docker pull ghcr.io/urutau-ltd/git-cone:<tag>
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

This fork aims to remain a practical drop-in replacement for recent `soft-serve`
deployments.

| Before         | After                                                  |
| -------------- | ------------------------------------------------------ |
| `soft serve`   | `cone serve` or `soft serve`                           |
| `soft browse`  | `cone browse` or `soft browse`                         |
| `SOFT_SERVE_*` | `GIT_CONE_*` preferred, `SOFT_SERVE_*` still supported |

What changed on purpose:

- the preferred binary name is now `cone`
- `soft` still works and maps to the same implementation
- the default server name is `Git Cone`
- the image stores data in `/git-cone/data`
- `git://` is off by default
- the server is SQLite-only

### Differences From `soft-serve`

This fork intentionally diverges from upstream in a few places. These are the
ones operators usually need to know before a migration:

| Area                   | `soft-serve` expectation                     | `git-cone` behavior                                                                  |
| ---------------------- | -------------------------------------------- | ------------------------------------------------------------------------------------ |
| Primary binary         | `soft`                                       | `cone` is preferred; `soft` remains as compatibility wrapper                         |
| Environment prefix     | `SOFT_SERVE_*`                               | `GIT_CONE_*` is preferred; `SOFT_SERVE_*` still works                                |
| Data path in container | `/soft-serve`                                | `/git-cone/data`                                                                     |
| Default server name    | Upstream default                             | `Git Cone`                                                                           |
| Database backends      | Upstream had more room for alternate drivers | SQLite-only                                                                          |
| `git://` daemon        | Historically available by default            | Disabled by default                                                                  |
| Strict mode            | Not present                                  | Available via `security.strict`                                                      |
| SSH crypto defaults    | Upstream defaults                            | Hardened KEX/cipher/MAC policy, including post-quantum KEX for newer OpenSSH clients |
| Health endpoint        | Not present                                  | `GET /health` returns JSON                                                           |
| Audit command          | Not present                                  | `ssh host audit`                                                                     |
| Repo integrity check   | Not present                                  | `ssh host repo verify <repo>`                                                        |
| Notifications          | No upstream support                          | Optional Gotify notifications added by this fork                                     |

Behavior that stays intentionally compatible:

- SSH TUI remains the main interface
- Git over SSH and HTTP still work the same way
- `soft serve`, `soft browse`, and `SOFT_SERVE_*` still work
- the SSH command surface stays close to upstream, with additive hardening
  features

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

## CLI Reference

This section documents the actual command surface exposed by the current fork.
It covers:

- the local binary CLI: `cone` and `soft`
- the SSH command interface users run against the server
- the repo subcommands under `repo`
- the internal transport commands Git uses under SSH

### Local Binary CLI

`cone` is the primary binary. `soft` is a compatibility wrapper over the same
implementation.

#### Root commands

| Command                  | Purpose                                         | Notes                  |
| ------------------------ | ----------------------------------------------- | ---------------------- |
| `cone serve`             | Start the server                                | `soft serve` works too |
| `cone browse [PATH]`     | Open the local TUI against a repository on disk | Defaults to `.`        |
| `cone admin migrate`     | Migrate the SQLite schema to the latest version | Local admin command    |
| `cone admin rollback`    | Roll back the previous database migration       | Local admin command    |
| `cone admin sync-hooks`  | Rewrite server-managed hooks for all repos      | Local admin command    |
| `cone hook pre-receive`  | Internal hook entrypoint                        | Hidden; called by Git  |
| `cone hook update`       | Internal hook entrypoint                        | Hidden; called by Git  |
| `cone hook post-receive` | Internal hook entrypoint                        | Hidden; called by Git  |
| `cone hook post-update`  | Internal hook entrypoint                        | Hidden; called by Git  |
| `cone man`               | Generate a manpage                              | Hidden                 |

#### Local command flags

`cone serve`

- `--sync-hooks` Rewrite hooks for all repositories before the server starts.

`cone browse [PATH]`

- no extra flags

`cone admin migrate`

- no extra flags

`cone admin rollback`

- no extra flags

`cone admin sync-hooks`

- no extra flags

`cone hook ...`

- `--config` Deprecated and ignored. Hook execution reads from the loaded config
  and environment.

### SSH Usage Model

The SSH interface has three modes:

1. `ssh -p 23231 host` Opens the interactive TUI.
2. `ssh -p 23231 host <repo>` Opens the TUI directly on a readable repository.
3. `ssh -p 23231 host <command> ...` Runs a non-interactive SSH command.

The SSH command help is dynamic and uses the configured public SSH URL to print
the correct host and port in examples.

### SSH Top-Level Commands

These are the public non-interactive commands available over SSH.

| Command                            | Purpose                                              | Access                                                              |
| ---------------------------------- | ---------------------------------------------------- | ------------------------------------------------------------------- |
| `audit`                            | Show server/session audit information                | Any authenticated user; limited output for unauthenticated sessions |
| `doctor`                           | Show effective server security and runtime settings  | Admin                                                               |
| `info`                             | Show information about the current user              | Authenticated user                                                  |
| `jwt [repository1 repository2...]` | Mint a JWT scoped to the given audience/repositories | Authenticated user                                                  |
| `pubkey ...`                       | Manage your own SSH public keys                      | Authenticated user                                                  |
| `repo ...`                         | Manage repositories and inspect repo content         | Varies by subcommand                                                |
| `set-username USERNAME`            | Change your own username                             | Authenticated user                                                  |
| `settings ...`                     | Read or change server-wide access settings           | Admin                                                               |
| `token ...`                        | Manage your own access tokens                        | Authenticated user                                                  |
| `user ...`                         | Manage users                                         | Admin                                                               |

### SSH Command Details

#### `audit`

```text
ssh -p 23231 host audit
```

Prints:

- server version
- current username and admin state when available
- remote client address and SSH client version
- active public key fingerprint
- active public key algorithm
- public key age when the DB has `created_at`
- auth mode and whether the session is keyless
- negotiated hostkey, cipher, KEX, and whether the KEX is post-quantum when the
  session exposes them
- owned and collaborator repo counts

Unauthenticated or keyless sessions only get the server version line.

#### `doctor`

```text
ssh -p 23231 host doctor
```

Prints the effective server-side settings that matter for operations and
hardening, including:

- strict mode state
- effective SSH/HTTP/stats/git listen addresses and public URLs
- LFS and LFS-over-SSH state
- hook timeout and SSH timeouts
- host/client key paths and whether those files exist
- hardened SSH KEX, cipher, and MAC policy

#### `info`

```text
ssh -p 23231 host info
```

Prints the current user's account details.

#### `jwt [repository1 repository2...]`

```text
ssh -p 23231 host jwt
ssh -p 23231 host jwt repo1 repo2
```

Creates a signed JWT using the server JWK pair. The listed repositories become
the JWT audience.

#### `set-username USERNAME`

```text
ssh -p 23231 host set-username new-name
```

Changes the username of the authenticated user.

### `pubkey` Commands

Manage the public keys attached to the current account.

| Command                        | Purpose                               |
| ------------------------------ | ------------------------------------- |
| `pubkey add AUTHORIZED_KEY`    | Add a public key to your account      |
| `pubkey remove AUTHORIZED_KEY` | Remove a public key from your account |
| `pubkey list`                  | List your current public keys         |

Aliases:

- `pubkeys`
- `publickey`
- `publickeys`

Examples:

```bash
ssh -p 23231 host pubkey list
ssh -p 23231 host pubkey add "ssh-ed25519 AAAA..."
ssh -p 23231 host pubkey remove "ssh-ed25519 AAAA..."
```

### `token` Commands

Manage HTTP/API access tokens for the current user.

| Command             | Purpose                                                     |
| ------------------- | ----------------------------------------------------------- |
| `token create NAME` | Create a token and print it once                            |
| `token list`        | List token IDs, names, creation dates, and expiration state |
| `token delete ID`   | Delete a token by numeric ID                                |

Aliases:

- `token`: alias group `access-token`
- `token list`: alias `ls`
- `token delete`: aliases `rm`, `remove`

Flags:

`token create NAME`

- `--expires-in` Expiration duration such as `1y`, `3mo`, `2w`, `5d4h`, or
  `1h30m`.

Examples:

```bash
ssh -p 23231 host token create "ci bot"
ssh -p 23231 host token create "pipe clone token" --expires-in 90d
ssh -p 23231 host token list
ssh -p 23231 host token delete 3
```

### `settings` Commands

These are server-wide settings stored in the database. They are admin-only.

| Command                               | Purpose                           |
| ------------------------------------- | --------------------------------- |
| `settings allow-keyless [true         | false]`                           |
| `settings anon-access [ACCESS_LEVEL]` | Get or set anonymous access level |

Valid `ACCESS_LEVEL` values:

- `no-access`
- `read-only`
- `read-write`
- `admin-access`

Examples:

```bash
ssh -p 23231 host settings allow-keyless false
ssh -p 23231 host settings anon-access no-access
```

Note:

- when `security.strict` is enabled, `allow-keyless` is forced off
- when `security.strict` is enabled, `anon-access` is forced to `no-access`

### `user` Commands

All `user` commands are admin-only.

| Command                                      | Purpose                  |
| -------------------------------------------- | ------------------------ |
| `user create USERNAME`                       | Create a new user        |
| `user delete USERNAME`                       | Delete a user            |
| `user list`                                  | List all users           |
| `user add-pubkey USERNAME AUTHORIZED_KEY`    | Add a key to a user      |
| `user remove-pubkey USERNAME AUTHORIZED_KEY` | Remove a key from a user |
| `user set-admin USERNAME [true               | false]`                  |
| `user info USERNAME`                         | Show user details        |
| `user set-username USERNAME NEW_USERNAME`    | Rename a user            |

Aliases:

- `user`: alias group `users`
- `user list`: alias `ls`

Flags:

`user create USERNAME`

- `-a`, `--admin` Create the user as admin.
- `-k`, `--key` Attach an initial public key.

Examples:

```bash
ssh -p 23231 host user create alice
ssh -p 23231 host user create pipe-bot --key "ssh-ed25519 AAAA..."
ssh -p 23231 host user create release-bot --admin
ssh -p 23231 host user set-admin alice true
ssh -p 23231 host user info alice
```

### `repo` Commands

`repo` is the largest command group.

Aliases for the group:

- `repos`
- `repository`
- `repositories`

The subcommands below are available under:

```text
ssh -p 23231 host repo ...
```

#### Repo listing and metadata

| Command                                     | Purpose                                 | Access                  |
| ------------------------------------------- | --------------------------------------- | ----------------------- |
| `repo list`                                 | List readable repositories              | Readable                |
| `repo info REPOSITORY`                      | Print repo metadata, branches, and tags | Readable                |
| `repo description REPOSITORY [DESCRIPTION]` | Get or set description                  | Write/admin             |
| `repo project-name REPOSITORY [NAME]`       | Get or set project name                 | Write/admin             |
| `repo private REPOSITORY [true              | false]`                                 | Get or set private flag |
| `repo hidden REPOSITORY [true               | false]`                                 | Get or set hidden flag  |
| `repo is-mirror REPOSITORY`                 | Report whether the repo is a mirror     | Readable                |

Flags:

`repo list`

- `-a`, `--all` Include hidden repositories that are otherwise readable.

Examples:

```bash
ssh -p 23231 host repo list
ssh -p 23231 host repo list --all
ssh -p 23231 host repo info demo
ssh -p 23231 host repo description demo "Internal deployment repo"
ssh -p 23231 host repo private demo true
```

#### Repo creation and lifecycle

| Command                           | Purpose                               | Access      |
| --------------------------------- | ------------------------------------- | ----------- |
| `repo create REPOSITORY`          | Create a repository                   | Write/admin |
| `repo import REPOSITORY REMOTE`   | Import a repository from a remote URL | Write/admin |
| `repo delete REPOSITORY`          | Delete a repository                   | Admin       |
| `repo rename REPOSITORY NEW_NAME` | Rename a repository                   | Admin       |
| `repo verify REPOSITORY`          | Run `git fsck --full`                 | Write/admin |

Flags:

`repo create REPOSITORY`

- `-p`, `--private`
- `-d`, `--description`
- `-n`, `--name`
- `-H`, `--hidden`

`repo import REPOSITORY REMOTE`

- `--lfs`
- `--lfs-endpoint`
- `-m`, `--mirror`
- `-p`, `--private`
- `-d`, `--description`
- `-n`, `--name`
- `-H`, `--hidden`

Examples:

```bash
ssh -p 23231 host repo create demo
ssh -p 23231 host repo create secret-repo --private --description "confidential"
ssh -p 23231 host repo import upstream https://example.com/repo.git --mirror
ssh -p 23231 host repo verify demo
ssh -p 23231 host repo rename demo demo-archive
```

#### Repo content inspection

| Command                                   | Purpose                         | Access   |
| ----------------------------------------- | ------------------------------- | -------- |
| `repo blob REPOSITORY [REFERENCE] [PATH]` | Print a file from a commit/tree | Readable |
| `repo tree REPOSITORY [REFERENCE] [PATH]` | Print the repository tree       | Readable |
| `repo commit REPOSITORY SHA`              | Print a commit diff             | Readable |

Flags:

`repo blob REPOSITORY [REFERENCE] [PATH]`

- `-r`, `--raw`
- `-l`, `--linenumber`
- `-c`, `--color`

`repo commit REPOSITORY SHA`

- `-c`, `--color`
- `-p`, `--patch`

Examples:

```bash
ssh -p 23231 host repo blob demo HEAD README.md
ssh -p 23231 host repo blob demo HEAD README.md --raw
ssh -p 23231 host repo tree demo HEAD
ssh -p 23231 host repo commit demo 0123456789abcdef --patch
```

#### Branch commands

| Command                                   | Purpose                       | Access      |
| ----------------------------------------- | ----------------------------- | ----------- |
| `repo branch list REPOSITORY`             | List branches                 | Readable    |
| `repo branch default REPOSITORY [BRANCH]` | Get or set the default branch | Write/admin |
| `repo branch delete REPOSITORY BRANCH`    | Delete a branch               | Write/admin |

Examples:

```bash
ssh -p 23231 host repo branch list demo
ssh -p 23231 host repo branch default demo main
ssh -p 23231 host repo branch delete demo old-branch
```

#### Tag commands

| Command                          | Purpose      | Access      |
| -------------------------------- | ------------ | ----------- |
| `repo tag list REPOSITORY`       | List tags    | Readable    |
| `repo tag delete REPOSITORY TAG` | Delete a tag | Write/admin |

Examples:

```bash
ssh -p 23231 host repo tag list demo
ssh -p 23231 host repo tag delete demo v0.1.0
```

#### Collaborator commands

| Command                                       | Purpose               | Access |
| --------------------------------------------- | --------------------- | ------ |
| `repo collab add REPOSITORY USERNAME [LEVEL]` | Add a collaborator    | Admin  |
| `repo collab remove REPOSITORY USERNAME`      | Remove a collaborator | Admin  |
| `repo collab list REPOSITORY`                 | List collaborators    | Admin  |

Valid collaborator levels:

- `no-access`
- `read-only`
- `read-write`
- `admin-access`

If omitted, `repo collab add` defaults to `read-write`.

Examples:

```bash
ssh -p 23231 host repo collab add demo alice read-only
ssh -p 23231 host repo collab add demo bob admin-access
ssh -p 23231 host repo collab list demo
ssh -p 23231 host repo collab remove demo alice
```

#### Webhook commands

All webhook commands are admin-only for the target repository.

| Command                                                               | Purpose                |
| --------------------------------------------------------------------- | ---------------------- |
| `repo webhook list REPOSITORY`                                        | List webhooks          |
| `repo webhook create REPOSITORY URL`                                  | Create a webhook       |
| `repo webhook delete REPOSITORY WEBHOOK_ID`                           | Delete a webhook       |
| `repo webhook update REPOSITORY WEBHOOK_ID`                           | Update a webhook       |
| `repo webhook deliveries list REPOSITORY WEBHOOK_ID`                  | List delivery attempts |
| `repo webhook deliveries redeliver REPOSITORY WEBHOOK_ID DELIVERY_ID` | Redeliver one attempt  |
| `repo webhook deliveries get REPOSITORY WEBHOOK_ID DELIVERY_ID`       | Inspect one delivery   |

Aliases:

- `repo webhook`: alias `webhooks`
- `repo webhook deliveries`: aliases `delivery`, `deliver`

Flags:

`repo webhook create REPOSITORY URL`

- `-e`, `--events`
- `-s`, `--secret`
- `-a`, `--active`
- `-c`, `--content-type`

`repo webhook update REPOSITORY WEBHOOK_ID`

- `-e`, `--events`
- `-s`, `--secret`
- `-a`, `--active`
- `-c`, `--content-type`
- `-u`, `--url`

`--content-type` accepts:

- `json`
- `form`

Examples:

```bash
ssh -p 23231 host repo webhook list demo
ssh -p 23231 host repo webhook create demo https://example.com/hook --events push --secret supersecret
ssh -p 23231 host repo webhook update demo 4 --active false
ssh -p 23231 host repo webhook deliveries list demo 4
ssh -p 23231 host repo webhook deliveries get demo 4 550e8400-e29b-41d4-a716-446655440000
```

### Internal SSH Transport Commands

These are not normally typed by humans. Git and Git LFS invoke them over SSH.

| Command                               | Purpose                   |
| ------------------------------------- | ------------------------- |
| `git-upload-pack REPO`                | Clone/fetch over SSH      |
| `git-upload-archive REPO`             | Archive access over SSH   |
| `git-receive-pack REPO`               | Push over SSH             |
| `git-lfs-authenticate REPO OPERATION` | Git LFS auth handshake    |
| `git-lfs-transfer REPO OPERATION`     | Git LFS transfer protocol |

If LFS is disabled, the LFS commands are not registered. If `lfs.ssh_enabled` is
false, `git-lfs-transfer` is not registered.

The Git transport URLs and repository workflow stay the same:

- SSH clone/push: `ssh://host:23231/<repo>.git`
- HTTP clone/push: `https://host/<repo>.git`

## Authentication and Strict Mode

`strict=true` does not disable token-based HTTP access. It does this:

- forces anonymous access to `no-access`
- disables keyless access
- clamps SSH timeouts
- clamps HTTP CORS to `http.public_url`
- keeps the stats listener on loopback when enabled

That means:

- SSH always requires an authorized key
- HTTP Git/LFS requires valid credentials
- access tokens still work for automation such as `pipe`
- `pipe` can keep using internal HTTP with a token behind Caddy; no server-local
  TLS changes are required for that flow

With `strict=true`, a non-private repo is not anonymously readable. Today there
is no per-repo “public override” when global anonymous access is forced off.

This is intentional hardening. If you need anonymous read access for non-private
repos, do not enable strict mode.

## Gotify Notifications

`soft-serve` did not ship with Gotify support. `git-cone` adds it.

Use it if you want the server to send notifications to Gotify for selected
events. This can be useful on its own, or alongside tools such as `pipe`.

When disabled, the server does not make notification network calls. When
enabled, notification delivery failures do not block SSH, Git, HTTP, hooks, or
pushes.

### Events

Current events emitted by the fork:

| Event                            | Trigger                                   | Priority |
| -------------------------------- | ----------------------------------------- | -------- |
| `git-cone: new user`             | A new user is created                     | `5`      |
| `git-cone: webhook failure`      | A webhook fails 3 times in a row          | `5`      |
| `git-cone: push to private repo` | A user pushes to a private repository     | `3`      |
| `git-cone: auth failure burst`   | 5 failed auth attempts from one IP in 60s | `7`      |

What these mean in practice:

- `new user`: a new account was created
- `webhook failure`: one webhook failed 3 times in a row
- `push to private repo`: someone pushed to a private repository
- `auth failure burst`: one IP accumulated 5 failed auth attempts in 60 seconds

### Configuration

YAML:

```yaml
notify:
  gotify:
    enabled: true
    url: "https://gotify.example.com"
    token: "YOUR_GOTIFY_APP_TOKEN"
    priority: 5
```

Environment variables:

- `GIT_CONE_NOTIFY_GOTIFY_ENABLED=true`
- `GIT_CONE_NOTIFY_GOTIFY_URL=https://gotify.example.com`
- `GIT_CONE_NOTIFY_GOTIFY_TOKEN=YOUR_GOTIFY_APP_TOKEN`
- `GIT_CONE_NOTIFY_GOTIFY_PRIORITY=5`

Compatibility variables with the old prefix also work:

- `SOFT_SERVE_NOTIFY_GOTIFY_ENABLED`
- `SOFT_SERVE_NOTIFY_GOTIFY_URL`
- `SOFT_SERVE_NOTIFY_GOTIFY_TOKEN`
- `SOFT_SERVE_NOTIFY_GOTIFY_PRIORITY`

The old prefix support here is provided by `git-cone`'s dual-prefix config
loader. It exists for migration convenience; it is not inherited Gotify support
from upstream.

### Example Compose

```yaml
services:
  soft-serve:
    image: ghcr.io/urutau-ltd/git-cone:latest
    environment:
      - GIT_CONE_NOTIFY_GOTIFY_ENABLED=true
      - GIT_CONE_NOTIFY_GOTIFY_URL=http://gotify:80
      - GIT_CONE_NOTIFY_GOTIFY_TOKEN=${GOTIFY_TOKEN}
      - GIT_CONE_NOTIFY_GOTIFY_PRIORITY=5
```

Notes:

- internal services such as `pipe` should talk directly to `http://gotify:80`
- avoid sending trusted service-to-service traffic through Anubis
- notification delivery failures do not block Git or SSH operations

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

This keeps working with `security.strict=true`. Strict mode disables anonymous
and keyless access, but it does not disable HTTP token auth for Git/LFS.

If you also run Gotify, let `pipe` talk to `http://gotify:80` directly on the
internal network. Do not put internal service-to-service traffic through Anubis.

## Hooks

Global hooks live under `$GIT_CONE_DATA_PATH/hooks/`. The server writes a sample
`update.sample` hook on first start.

## Configuration

The generated `config.yaml` is the canonical reference for the settings this
fork currently exposes.

Notable defaults:

```yaml
db:
  driver: "sqlite"

git:
  listen_addr: "" # git:// disabled by default

hooks:
  timeout: 30

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
- `GIT_CONE_HOOKS_TIMEOUT`
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

The repo no longer ships systemd-centric packaging inherited from upstream. If
you need host-managed services, optional examples live in
[docs/service-managers.md](docs/service-managers.md) for:

- SysVinit
- OpenRC
- Runit
- GNU Shepherd

## License

[MIT](LICENSE)
