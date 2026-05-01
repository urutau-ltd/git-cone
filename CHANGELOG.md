# Changelog

## v0.13.3 - 2026-05-01

- Fix `ssh audit` so negotiated `hostkey`, `cipher`, `kex`, and `kex-pq` are reported correctly from wrapped `*ssh.ServerConn` session metadata.
- Reject unknown SSH public keys during authentication; only registered user keys or configured bootstrap admin keys may complete public-key auth.
- Add defense-in-depth validation in the SSH authentication middleware so unknown public keys are denied even if the initial auth gate is bypassed.
- Fix test infrastructure panic in `pkg/test.RandomPort()` when local listeners are unavailable.
- Make SSH session tests independent from env var overrides and skip listener-dependent cases cleanly in restricted environments.
