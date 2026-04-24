package sshpolicy

import gossh "golang.org/x/crypto/ssh"

func HardenedKeyExchanges() []string {
	return []string{
		gossh.KeyExchangeMLKEM768X25519,
		gossh.KeyExchangeCurve25519,
		"curve25519-sha256@libssh.org",
		gossh.KeyExchangeDH16SHA512,
		"diffie-hellman-group18-sha512",
	}
}

func HardenedCiphers() []string {
	return []string{
		"chacha20-poly1305@openssh.com",
		"aes256-gcm@openssh.com",
		"aes128-gcm@openssh.com",
		"aes256-ctr",
		"aes192-ctr",
		"aes128-ctr",
	}
}

func HardenedMACs() []string {
	return []string{
		"hmac-sha2-256-etm@openssh.com",
		"hmac-sha2-512-etm@openssh.com",
		"hmac-sha2-256",
	}
}

func IsPostQuantumKEX(name string) bool {
	switch name {
	case gossh.KeyExchangeMLKEM768X25519, "sntrup761x25519-sha512", "sntrup761x25519-sha512@openssh.com":
		return true
	default:
		return false
	}
}
