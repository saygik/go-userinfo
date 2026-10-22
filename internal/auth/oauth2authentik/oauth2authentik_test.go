package oauth2authentik

import "testing"

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
