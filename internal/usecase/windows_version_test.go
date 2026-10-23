package usecase

import "testing"

func TestWindowsVersionToHuman(t *testing.T) {
	cases := []struct {
		name      string
		osName    string
		osVersion string
		want      string
	}{
		{"Windows 11 24H2", "Windows 11 Pro", "10.0 (26100)", "24H2"},
		{"Windows Server 2025", "Windows Server 2025 Standard", "10.0 (26100)", "Server 2025"},
		{"Windows Server 2025 Datacenter", "Windows Server 2025 Datacenter", "10.0 (26100)", "Server 2025"},
		{"Windows Server 2022", "Windows Server 2022 Standard", "10.0 (20348)", "Server 2022"},
		{"Windows Server 2019", "Windows Server 2019 Standard", "10.0 (17763)", "Server 2019"},
		{"Windows Server 2012 R2", "Windows Server 2012 R2 Standard", "6.3 (9600)", "Server 2012 R2"},
		{"Windows 11 23H2", "Windows 11 Pro", "10.0 (22631)", "23H2"},
		{"Windows 10 22H2", "Windows 10 Pro", "10.0 (19045)", "22H2"},
		{"Windows 8.1", "Windows 8.1 Pro", "6.3 (9600)", "8.1"},
		{"Windows 8", "Windows 8 Pro", "6.2 (9200)", "8"},
		{"empty version", "Windows 11 Pro", "   ", ""},
		{"no digits", "Windows Server 2025 Standard", "unknown", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := windowsVersionToHuman(tc.osName, tc.osVersion); got != tc.want {
				t.Fatalf("windowsVersionToHuman(%q, %q) = %q, want %q", tc.osName, tc.osVersion, got, tc.want)
			}
		})
	}
}

func TestWindowsVersionToHumanWithSupport(t *testing.T) {
	if got := windowsVersionToHumanWithSupport("Windows Server 2025 Standard", "10.0 (26100)"); !got.Supported || got.Version != "Server 2025" {
		t.Fatalf("Server 2025: got %+v, want supported \"Server 2025\"", got)
	}
	if got := windowsVersionToHumanWithSupport("Windows 10 Pro", "10.0 (19041)"); got.Supported {
		t.Fatalf("build 19041 must be unsupported, got %+v", got)
	}
	if got := windowsVersionToHumanWithSupport("Windows Server 2025 Standard", ""); got.Supported || got.Version != "" {
		t.Fatalf("empty version: got %+v, want unsupported empty", got)
	}
}

func TestOSFamilyName(t *testing.T) {
	cases := []struct {
		osName string
		want   string
	}{
		{"Windows 11 Pro", "Windows 11"},
		{"Windows 10 Pro", "Windows 10"},
		{"Windows Server 2025 Standard", "Windows Server 2025"},
		// AD иногда отдаёт название с двойными пробелами
		{"Windows Server  2025  Datacenter", "Windows Server 2025"},
		{"Windows Server 2022 Standard", "Windows Server 2022"},
		{"Windows Server 2012 R2 Standard", "Windows Server 2012 R2"},
		// Незнакомая версия не должна исчезать из статистики
		{"Windows Server 2028 Standard", "Windows Server (другая версия)"},
		{"Windows 12 Pro", "Windows (другая версия)"},
		{"Linux", ""},
		{"", ""},
	}

	for _, tc := range cases {
		t.Run(tc.osName, func(t *testing.T) {
			if got := osFamilyName(tc.osName); got != tc.want {
				t.Fatalf("osFamilyName(%q) = %q, want %q", tc.osName, got, tc.want)
			}
		})
	}
}

func TestOSFamilyRank(t *testing.T) {
	if osFamilyRank("Windows Server 2025 Standard") >= osFamilyRank("Windows Server 2022 Standard") {
		t.Fatal("Windows Server 2025 must be ranked above Windows Server 2022")
	}
	if osFamilyRank("Windows Server  2025  Datacenter") != osFamilyRank("Windows Server 2025 Standard") {
		t.Fatal("rank must not depend on extra spaces")
	}
	if osFamilyRank("Windows XP") >= osFamilyRank("Some unknown OS") {
		t.Fatal("known families must be ranked above unknown ones")
	}
}
