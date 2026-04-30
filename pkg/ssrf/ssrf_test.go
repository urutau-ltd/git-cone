package ssrf

import (
	"net"
	"testing"
)

func TestIsPrivateOrInternal(t *testing.T) {
	tests := []struct {
		name   string
		ip     string
		isPriv bool
	}{
		{"Google DNS", "8.8.8.8", false},
		{"Cloudflare DNS", "1.1.1.1", false},
		{"Public IPv6", "2001:4860:4860::8888", false},
		{"127.0.0.1", "127.0.0.1", true},
		{"127.1.2.3", "127.1.2.3", true},
		{"::1", "::1", true},
		{"10.0.0.1", "10.0.0.1", true},
		{"192.168.1.1", "192.168.1.1", true},
		{"172.16.0.1", "172.16.0.1", true},
		{"172.31.255.255", "172.31.255.255", true},
		{"169.254.169.254", "169.254.169.254", true},
		{"169.254.1.1", "169.254.1.1", true},
		{"fe80::1", "fe80::1", true},
		{"0.0.0.0", "0.0.0.0", true},
		{"255.255.255.255", "255.255.255.255", true},
		{"240.0.0.1", "240.0.0.1", true},
		{"100.64.0.1", "100.64.0.1", true},
		{"100.127.255.255", "100.127.255.255", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.ip)
			}
			if got := isPrivateOrInternal(ip); got != tt.isPriv {
				t.Errorf("isPrivateOrInternal(%s) = %v, want %v", tt.ip, got, tt.isPriv)
			}
		})
	}
}

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		want     bool
	}{
		{"localhost", "localhost", true},
		{"LOCALHOST", "LOCALHOST", true},
		{"localhost.localdomain", "localhost.localdomain", true},
		{"test.localhost", "test.localhost", true},
		{"example.com", "example.com", false},
		{"localhos", "localhos", false},
		{"localhost.com", "localhost.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLocalhost(tt.hostname); got != tt.want {
				t.Errorf("isLocalhost(%s) = %v, want %v", tt.hostname, got, tt.want)
			}
		})
	}
}

func TestValidateIPBeforeDial(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{"public IP", "8.8.8.8", false},
		{"private IP", "192.168.1.1", true},
		{"loopback", "127.0.0.1", true},
		{"link-local", "169.254.169.254", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.ip)
			}
			err := ValidateIPBeforeDial(ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPBeforeDial(%s) error = %v, wantErr %v", tt.ip, err, tt.wantErr)
			}
		})
	}
}
