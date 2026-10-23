package http

import "testing"

func TestIsTrustedOrigin(t *testing.T) {
	cases := []struct {
		name   string
		origin string
		host   string
		want   bool
	}{
		{"same host", "https://userinfoapi.brnv.rw", "userinfoapi.brnv.rw", true},
		{"same domain", "https://userinfo.brnv.rw", "userinfoapi.brnv.rw", true},
		{"same domain with port", "https://w-058.brnv.rw", "w-058-258049004.brnv.rw:9099", true},
		{"local development", "http://localhost:3000", "localhost:9099", true},
		{"local development to remote api", "http://localhost:3000", "userinfoapi.brnv.rw", true},
		{"local development 127.0.0.1 to remote api", "http://127.0.0.1:3000", "userinfoapi.brnv.rw", true},
		{"foreign domain", "https://evil.com", "userinfoapi.brnv.rw", false},
		{"lookalike domain", "https://brnv.rw.evil.com", "userinfoapi.brnv.rw", false},
		{"empty origin", "", "userinfoapi.brnv.rw", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isTrustedOrigin(tc.origin, tc.host); got != tc.want {
				t.Fatalf("isTrustedOrigin(%q, %q) = %v, want %v", tc.origin, tc.host, got, tc.want)
			}
		})
	}
}
