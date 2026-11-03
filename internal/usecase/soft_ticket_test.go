package usecase

import (
	"strings"
	"testing"
)

func TestNormalizeSoftTicketID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "empty", input: "", want: ""},
		{name: "spaces only", input: "   ", want: ""},
		{name: "digits with spaces", input: " 123456 ", want: "123456"},
		{name: "leading zeros", input: "00123", want: "00123"},
		{name: "max length", input: strings.Repeat("1", softTicketIDMaxLen), want: strings.Repeat("1", softTicketIDMaxLen)},
		{name: "too long", input: strings.Repeat("1", softTicketIDMaxLen+1), wantErr: true},
		{name: "letters", input: "12ab34", wantErr: true},
		{name: "url", input: "https://support.rw/front/ticket.form.php?id=123", wantErr: true},
		{name: "negative", input: "-123", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := normalizeSoftTicketID(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatalf("normalizeSoftTicketID(%q) = %q, want error", test.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeSoftTicketID(%q) returned error: %v", test.input, err)
			}
			if got != test.want {
				t.Fatalf("normalizeSoftTicketID(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}
