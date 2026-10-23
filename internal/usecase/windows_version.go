package usecase

import (
	"strconv"
	"strings"
)

type WindowsVersionInfo struct {
	Version   string `json:"version"`
	Supported bool   `json:"supported"`
}

// соответствие build-номера Windows маркетинговому имени версии (клиентские ОС)
var windowsBuildToHuman = map[int]string{
	28000: "26H1", // Windows 11
	26200: "25H2", // Windows 11
	26100: "24H2", // Windows 11 24H2
	22631: "23H2", // Windows 11 23H2
	22621: "22H2", // Windows 11 22H2
	22000: "21H2", // Windows 11 21H2

	19045: "22H2", // Windows 10 22H2
	19044: "21H2", // Windows 10 21H2
	19043: "21H1", // Windows 10 21H1
	19042: "20H2", // Windows 10 20H2
	19041: "2004",
	18363: "1909",
	18362: "1903",
	17763: "1809",
	17134: "1803",
	16299: "1709",
	15063: "1703",
	14393: "1607",
	10586: "1511",
	10240: "1507",

	9600: "8.1", // Windows 8.1
	9200: "8",   // Windows 8
	7601: "7 SP1", // Windows 7 SP1 / Server 2008 R2
	7600: "7",

	// Server-only ветка (для случаев, когда имя ОС не содержит "server")
	20348: "Server 2022",
	6003:  "Server 2008",
	3790:  "Server 2003",
	2195:  "Server 2000",
}

// serverBuildToHuman - маркетинговые имена серверных сборок.
// Часть сборок (26100, 17763, 14393) совпадает с клиентскими Windows,
// поэтому имя версии зависит от operatingSystem, а не только от build.
var serverBuildToHuman = map[int]string{
	26100: "Server 2025",
	20348: "Server 2022",
	17763: "Server 2019",
	14393: "Server 2016",
	9600:  "Server 2012 R2",
	9200:  "Server 2012",
	7601:  "Server 2008 R2",
	7600:  "Server 2008",
	6003:  "Server 2008",
	6002:  "Server 2008",
	6001:  "Server 2008",
	3790:  "Server 2003",
	2195:  "Server 2000",
}

// неподдерживаемые сборки (общая политика для клиентских и серверных ОС)
var unsupportedBuilds = map[int]bool{
	19041: true, // 2004
	19042: true, // 2004
	19043: true, // 2004
	18363: true, // 1909
	18362: true, // 1903
	17763: true, // 1809
	17134: true, // 1803
	16299: true, // 1709
	15063: true, // 1703
	14393: true, // 1607
	10586: true, // 1511
	10240: true, // 1507
	7601:  true, // 7 SP1
	7600:  true, // 7 SP1
	9600:  true, // Server 2012 R2
	6003:  true, // Server 2008
	3790:  true, // Server 2003
	2195:  true, // Server 2000
}

// isServerOSName определяет серверную ОС по строке operatingSystem.
func isServerOSName(osName string) bool {
	return strings.Contains(strings.ToLower(osName), "server")
}

// versionNumber выдёргивает номер сборки из строки operatingSystemVersion ("10.0 (26100)" -> "26100").
func versionNumber(osVersion string) string {
	osVersion = strings.TrimSpace(osVersion)
	if osVersion == "" {
		return ""
	}

	// оставляем только числовые последовательности и берём последнюю (build)
	parts := strings.FieldsFunc(osVersion, func(r rune) bool {
		return r < '0' || r > '9'
	})
	if len(parts) == 0 {
		return ""
	}

	buildStr := parts[len(parts)-1]
	return buildStr
}

// windowsVersionToHuman преобразует operatingSystem + operatingSystemVersion AD
// в человекочитаемый вид: 24H2 / 22H2 / Server 2025 / Server 2012 R2 и т.п.
// Примеры:
//   - "Windows 11 Pro", "10.0 (26100)"          -> "24H2"
//   - "Windows Server 2025 Standard", "10.0 (26100)" -> "Server 2025"
//   - "Windows Server 2022 Standard", "10.0 (20348)" -> "Server 2022"
func windowsVersionToHuman(osName, osVersion string) string {
	build, err := strconv.Atoi(versionNumber(osVersion))
	if err != nil {
		return ""
	}

	// Для серверных ОС приоритет у серверных названий: build 26100 -
	// это одновременно Windows 11 24H2 и Windows Server 2025.
	if isServerOSName(osName) {
		if human, ok := serverBuildToHuman[build]; ok {
			return human
		}
	}

	if human, ok := windowsBuildToHuman[build]; ok {
		return human
	}

	return ""
}

func windowsVersionToHumanWithSupport(osName, osVersion string) WindowsVersionInfo {
	build, err := strconv.Atoi(versionNumber(osVersion))
	if err != nil {
		return WindowsVersionInfo{Version: "", Supported: false}
	}

	version := windowsVersionToHuman(osName, osVersion)

	// Дополнительно проверяем, что версия известна
	supported := version != "" && !unsupportedBuilds[build]

	return WindowsVersionInfo{
		Version:   version,
		Supported: supported,
	}
}
