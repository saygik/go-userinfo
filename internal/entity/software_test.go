package entity

import "testing"

func TestNormalizeSoftUserAccountStatus(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "active", input: SoftUserAccountStatusActive, want: SoftUserAccountStatusActive},
		{name: "disabled with spaces", input: "  disabled  ", want: SoftUserAccountStatusDisabled},
		{name: "blocked case insensitive", input: "BLOCKED", want: SoftUserAccountStatusBlocked},
		{name: "deleted", input: SoftUserAccountStatusDeleted, want: SoftUserAccountStatusDeleted},
		{name: "empty defaults to active", input: "", want: SoftUserAccountStatusActive},
		{name: "unknown defaults to active", input: "removed", want: SoftUserAccountStatusActive},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := NormalizeSoftUserAccountStatus(test.input); got != test.want {
				t.Fatalf("NormalizeSoftUserAccountStatus(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}
