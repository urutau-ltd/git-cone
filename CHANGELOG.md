# Changelog

## v0.14.0 - 2026-10-26

- Sync with the final upstream Soft Serve `v0.12.4` feature set: anonymous and
  keyless access overrides, default repository creation, `repo mirror`,
  multi-location README discovery, hardened import/mirror validation, and the
  upstream security fixes.
- Keep the fork's hardening and operations layer on top: strict mode, hardened
  SSH algorithms and rate limiting, `GIT_CONE_*` precedence, Gotify alerts,
  `audit`, `doctor`, and `repo verify`.
- Move to Go 1.27.3, apply `go fix` modernizers, and refresh dependencies.
- Use Zig's C compiler (`CC="zig cc"`) for cgo-based tooling in the Guix
  environment.
- Remove the retired Codeberg pipe configuration and the upstream Taskfile in
  favor of the fork's GitHub Actions and Makefile.

## v0.13.3 - 2026-05-01

- Fix `ssh audit` so negotiated `hostkey`, `cipher`, `kex`, and `kex-pq` are
  reported correctly from wrapped `*ssh.ServerConn` session metadata.
- Reject unknown SSH public keys during authentication; only registered user
  keys or configured bootstrap admin keys may complete public-key auth.
- Add defense-in-depth validation in the SSH authentication middleware so
  unknown public keys are denied even if the initial auth gate is bypassed.
- Fix test infrastructure panic in `pkg/test.RandomPort()` when local listeners
  are unavailable.
- Make SSH session tests independent from env var overrides and skip
  listener-dependent cases cleanly in restricted environments.
