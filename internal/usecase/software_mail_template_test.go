package usecase

import (
	"strings"
	"testing"
)

func TestRenderSoftwareMailTemplate(t *testing.T) {
	values := map[string]string{
		"permission": "Доступ к отчётам",
		"login":      "ivanov",
		"fio":        "Иванов Иван",
		"system":     "1С",
		"enddate":    "01.12.2026",
	}

	got := renderSoftwareMailTemplate(
		`{{fio}} ({{login}}): «{{permission}}» в {{system}} истекает {{enddate}}.`,
		values,
	)
	want := `Иванов Иван (ivanov): «Доступ к отчётам» в 1С истекает 01.12.2026.`
	if got != want {
		t.Fatalf("renderSoftwareMailTemplate() = %q, want %q", got, want)
	}
}

func TestRenderSoftwareMailTemplateKeepsUnknownPlaceholders(t *testing.T) {
	got := renderSoftwareMailTemplate("{{known}} {{unknown}}", map[string]string{"known": "ok"})
	want := "ok {{unknown}}"
	if got != want {
		t.Fatalf("renderSoftwareMailTemplate() = %q, want %q", got, want)
	}
}

func TestSoftwareAdminMatchesUser(t *testing.T) {
	cases := []struct {
		name        string
		admin       map[string]interface{}
		currentUser string
		want        bool
	}{
		{
			name:        "matches login case-insensitively",
			admin:       map[string]interface{}{"name": "Ivanov"},
			currentUser: "ivanov",
			want:        true,
		},
		{
			name:        "matches upn",
			admin:       map[string]interface{}{"name": "ivanov", "userPrincipalName": "ivanov@brnv.rw"},
			currentUser: "ivanov@brnv.rw",
			want:        true,
		},
		{
			name:        "does not match other user",
			admin:       map[string]interface{}{"name": "petrov"},
			currentUser: "ivanov",
			want:        false,
		},
		{
			name:        "ignores empty current user",
			admin:       map[string]interface{}{"name": ""},
			currentUser: "  ",
			want:        false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := softwareAdminMatchesUser(tc.admin, tc.currentUser); got != tc.want {
				t.Fatalf("softwareAdminMatchesUser() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDefaultSoftwareMailTemplate(t *testing.T) {
	u := &UseCase{}
	template := u.defaultSoftwareMailTemplate()

	if template.Subject != defaultSoftwareMailSubject {
		t.Fatalf("subject = %q, want %q", template.Subject, defaultSoftwareMailSubject)
	}
	if template.Body != defaultSoftwareMailBody {
		t.Fatalf("body = %q, want %q", template.Body, defaultSoftwareMailBody)
	}
	if template.IsCustom {
		t.Fatal("default template must not be marked as custom")
	}
	if !strings.Contains(template.Body, "{{system}}") {
		t.Fatalf("default body should contain system placeholder: %q", template.Body)
	}
}
