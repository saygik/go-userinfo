package usecase

import (
	"strconv"
	"strings"
)

type WindowsVersionInfo struct {
	Version   string `json:"version"`
	Supported bool   `json:"supported"`
}

// соответствие build-номера Windows маркетинговому имени версии
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

	7601: "7 SP1", // Windows 7 SP1 / Server 2008 R2

	// Server-only ветка
	20348: "Server 2022",
	9600:  "Server 2012 R2",
	6003:  "Server 2008",
	3790:  "Server 2003",
	2195:  "Server 2000",
}

// windowsVersionToHuman преобразует строку operatingSystemVersion AD в вид 24H2 / 22H2 / Server 2012 R2 и т.п.
// Примеры входа:
//   - "10.0 (26100)"  -> "24H2"
//   - "10.0 (19045)"  -> "22H2"
//   - "6.3 (9600)"    -> "Server 2012 R2"
//   - "10.0 (20348)"  -> "Server 2022"
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
func windowsVersionToHuman(osVersion string) string {

	buildStr := versionNumber(osVersion)
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		return ""
	}

	if human, ok := windowsBuildToHuman[build]; ok {
		return human
	}

	return ""
}

func windowsVersionToHumanWithSupport(osVersion string) WindowsVersionInfo {
	buildStr := versionNumber(osVersion)
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		return WindowsVersionInfo{Version: "", Supported: false}
	}

	version := ""
	if human, ok := windowsBuildToHuman[build]; ok {
		version = human
	}

	// Список неподдерживаемых билдов
	unsupportedBuilds := map[int]bool{
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

	supported := !unsupportedBuilds[build]

	// Дополнительно проверяем, что версия известна
	if version == "" {
		supported = false
	}

	return WindowsVersionInfo{
		Version:   version,
		Supported: supported,
	}
}
