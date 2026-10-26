package oauth2authentik

import (
	"encoding/base64"
	"testing"
)

// Вектор из RFC 7636 (Appendix B).
func TestCodeChallenge(t *testing.T) {
	const verifier = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	const want = "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"

	o := OAuth2{}
	if got := o.CodeChallenge(verifier); got != want {
		t.Fatalf("CodeChallenge(%q) = %q, want %q", verifier, got, want)
	}
}

// Не-JWT и подделанный токен не должны приводить к панике.
func TestIntrospectRejectsInvalidTokens(t *testing.T) {
	o := OAuth2{}

	for _, token := range []string{"", "-", "not-a-jwt", "a.b.c"} {
		if _, err := o.IntrospectOAuth2Token(token, false); err == nil {
			t.Fatalf("IntrospectOAuth2Token(%q) вернул nil-ошибку", token)
		}
	}
}

// Refresh без токена не паникует.
func TestRefreshEmptyToken(t *testing.T) {
	o := OAuth2{}
	if _, err := o.Refresh(""); err == nil {
		t.Fatal("Refresh(\"\") вернул nil-ошибку")
	}
}

// Отзыв без revocation endpoint - no-op без ошибки.
func TestRevokeWithoutEndpoint(t *testing.T) {
	o := OAuth2{}
	if err := o.Revoke("token"); err != nil {
		t.Fatalf("Revoke() = %v, want nil", err)
	}
}

// iss из непроверенного JWT используется только для выбора verifier'а.
func TestTokenIssuer(t *testing.T) {
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"iss":"https://sso.brnv.rw/application/o/grafana/"}`))
	token := "header." + payload + ".signature"

	if got := tokenIssuer(token); got != "https://sso.brnv.rw/application/o/grafana/" {
		t.Fatalf("tokenIssuer() = %q", got)
	}
	for _, bad := range []string{"", "not-a-jwt", "a.b.c", "a.!!!.c"} {
		if got := tokenIssuer(bad); got != "" {
			t.Fatalf("tokenIssuer(%q) = %q, want empty", bad, got)
		}
	}
}

func TestIsTrustedIssuer(t *testing.T) {
	o := OAuth2{mainIssuer: "https://sso.brnv.rw/application/o/userinfotest/"}

	cases := []struct {
		name   string
		issuer string
		want   bool
	}{
		{"main application", "https://sso.brnv.rw/application/o/userinfotest/", true},
		{"grafana application", "https://sso.brnv.rw/application/o/grafana/", true},
		{"another application", "https://sso.brnv.rw/application/o/my-app_2/", true},
		{"other host", "https://evil.com/application/o/grafana/", false},
		{"other scheme", "http://sso.brnv.rw/application/o/grafana/", false},
		{"not an authentik app path", "https://sso.brnv.rw/oauth2/grafana/", false},
		{"path traversal", "https://sso.brnv.rw/application/o/../../evil/", false},
		{"empty slug", "https://sso.brnv.rw/application/o/", false},
		{"userinfo in url", "https://user@sso.brnv.rw/application/o/grafana/", false},
		{"query in url", "https://sso.brnv.rw/application/o/grafana/?x=1", false},
		{"empty issuer", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := o.isTrustedIssuer(tc.issuer); got != tc.want {
				t.Fatalf("isTrustedIssuer(%q) = %v, want %v", tc.issuer, got, tc.want)
			}
		})
	}
}
