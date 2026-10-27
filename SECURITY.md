# Security Policy

`git-cone` is a security-hardened fork of Soft Serve. Upstream Soft Serve has
been sunset and no longer receives security updates, so fixes for this fork are
maintained here.

## Supported versions

| Version    | Supported          |
| ---------- | ------------------ |
| `v0.14.x`  | :white_check_mark: |
| `<= v0.13` | :x:                |

Only the latest release receives security fixes. Please reproduce issues against
the newest tagged release before reporting them.

## Reporting a vulnerability

**Please do not open a public issue for a suspected vulnerability.**

Preferred reporting channel: GitHub private vulnerability reporting.

- <https://github.com/urutau-ltd/git-cone/security/advisories/new>

If you cannot use GitHub advisories, contact the maintainers privately at
`softwarelibre@urutau-ltd.org`. The Codeberg/Forgejo mirror is available at
<https://sl.urutau-ltd.org/urutau-ltd/git-cone>.

Please include:

- the affected version or commit;
- the affected component (SSH, HTTP Git, LFS, webhooks, authentication, etc.);
- reproduction steps or a proof of concept;
- the expected and observed impact;
- any suggested remediation or patch.

We aim to acknowledge reports within 3 business days and provide an initial
assessment within 7 business days. Please allow time for a fix and coordinated
disclosure before publishing details.

## Scope

Reports are especially welcome for:

- SSH authentication, session isolation, command authorization, and the TUI;
- HTTP Git smart protocol, LFS, bearer tokens, JWTs, and access tokens;
- repository name/path handling, storage confinement, and symlink escapes;
- repository privacy, anonymous/keyless access, and collaborator permissions;
- webhook delivery, SSRF protections, repository imports, and mirror syncs;
- configuration loading, `GIT_CONE_*` overrides, and secret handling;
- dependency vulnerabilities that affect reachable code.

General hardening observations and defense-in-depth suggestions are also useful,
but please keep them out of public issue trackers until triaged.

## Hardening baseline

The fork already enables, among other things:

- hardened SSH KEX, cipher, and MAC defaults;
- SSH input rate limiting and rejection of unknown public keys;
- SSRF validation for webhooks, LFS endpoints, imports, and mirror syncs;
- SSH host-key verification and optional import/mirror host allowlists;
- strict mode, which disables keyless access and clamps exposed surfaces;
- security headers for HTTP responses and hardened LFS authorization.

See the README and generated `config.yaml` for the available security settings.
