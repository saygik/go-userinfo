package usecase

import (
	"sort"
	"strconv"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

// GetADComputersVersions возвращает количество компьютеров в домене, сгруппированное по паре
// operatingSystem + operatingSystemVersion.
// Дополнительно вычисляет operatingSystemVersionHuman (например, 24H2, Server 2012 R2).
func (u *UseCase) GetADComputersVersions(domain, user string) ([]entity.ComputerVersionCount, error) {
	if !u.ad.IsDomainExist(domain) {
		return nil, u.Error("домен " + domain + " отсутствует в системе")
	}

	access := u.GetAccessToResource(domain, user)
	if access == -1 {
		return nil, u.Error("у вас нет прав на просмотр информации по домену " + domain)
	}

	comps, err := u.ad.GetDomainComputers(domain)
	if err != nil {
		return nil, err
	}

	type key struct {
		os    string
		raw   string
		human string
	}

	counts := map[key]int{}
	for _, c := range comps {
		raw, _ := c["operatingSystemVersion"].(string)
		if strings.TrimSpace(raw) == "" {
			// Пустые / пробельные версии ОС пропускаем совсем
			continue
		}
		osName, _ := c["operatingSystem"].(string)
		human := windowsVersionToHuman(raw)
		counts[key{os: osName, raw: raw, human: human}]++
	}

	res := make([]entity.ComputerVersionCount, 0, len(counts))
	for k, cnt := range counts {
		family := osFamilyName(k.os)
		versionName := versionNumber(k.raw)
		if versionName != "" {
			// Проверяем LTSC/LTSB и заменяем
			ltscNames := map[string]string{
				"10240": "LTSB 2015",
				"14393": "LTSB 2016",
				"17763": "LTSC 2019",
				"19044": "LTSC 2021",
				"26100": "LTSC 2024",
			}
			if name, ok := ltscNames[versionName]; ok {
				versionName = name + " (" + versionName + ")"
			} else {
				versionName = "(" + versionName + ")"
			}
		}
		res = append(res, entity.ComputerVersionCount{
			OperatingSystem:             k.os,
			OperatingSystemFamily:       family,
			OperatingSystemVersion:      versionName,
			OperatingSystemVersionHuman: k.human,
			Count:                       cnt,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		a := res[i]
		b := res[j]

		// 1. Сначала группируем по семейству ОС (Windows 11, затем 10, затем 8.x, 7 и т.д.)
		fa := osFamilyRank(a.OperatingSystem)
		fb := osFamilyRank(b.OperatingSystem)
		if fa != fb {
			return fa < fb
		}

		// 2. Внутри семейства сортируем по build-номеру версии (operatingSystemVersion) по убыванию
		ba := extractBuildNumber(a.OperatingSystemVersion)
		bb := extractBuildNumber(b.OperatingSystemVersion)
		if ba != bb {
			return ba > bb
		}

		// 3. Для стабильности — по имени ОС и строке версии
		if a.OperatingSystem != b.OperatingSystem {
			return a.OperatingSystem < b.OperatingSystem
		}
		if a.OperatingSystemVersion != b.OperatingSystemVersion {
			return a.OperatingSystemVersion < b.OperatingSystemVersion
		}

		// 4. В крайнем случае — по количеству (по убыванию)
		return a.Count > b.Count
	})

	return res, nil
}

// osFamilyRank задаёт порядок семейств ОС, чтобы однотипные ОС шли рядом.
// Меньшее значение — выше в списке.
func osFamilyRank(osName string) int {
	s := strings.ToLower(osName)
	switch {
	case strings.Contains(s, "windows 11"):
		return 10
	case strings.Contains(s, "windows 10"):
		return 20
	case strings.Contains(s, "windows 8.1"):
		return 30
	case strings.Contains(s, "windows 8"):
		return 40
	case strings.Contains(s, "windows 7"):
		return 50
	case strings.Contains(s, "windows vista"):
		return 60
	case strings.Contains(s, "windows xp"):
		return 70
	case strings.Contains(s, "server 2022"):
		return 110
	case strings.Contains(s, "server 2019"):
		return 120
	case strings.Contains(s, "server 2016"):
		return 130
	case strings.Contains(s, "server 2012"):
		return 140
	case strings.Contains(s, "server 2008"):
		return 150
	case strings.Contains(s, "server 2003"):
		return 160
	case strings.Contains(s, "server 2000"):
		return 170
	default:
		return 1000
	}
}

// osFamilyName возвращает человекочитаемое название семейства ОС для блока (Windows 11, Windows 10, Windows 7, Windows Server 2012 и т.д.).
func osFamilyName(osName string) string {
	s := strings.ToLower(osName)
	switch {
	case strings.Contains(s, "windows 11"):
		return "Windows 11"
	case strings.Contains(s, "windows 10"):
		return "Windows 10"
	case strings.Contains(s, "windows 8.1"):
		return "Windows 8.1"
	case strings.Contains(s, "windows 8"):
		return "Windows 8"
	case strings.Contains(s, "windows 7"):
		return "Windows 7"
	case strings.Contains(s, "windows vista"):
		return "Windows Vista"
	case strings.Contains(s, "windows xp"):
		return "Windows XP"
	case strings.Contains(s, "server 2022"):
		return "Windows Server 2022"
	case strings.Contains(s, "server 2019"):
		return "Windows Server 2019"
	case strings.Contains(s, "server 2016"):
		return "Windows Server 2016"
	case strings.Contains(s, "server 2012 r2"):
		return "Windows Server 2012 R2"
	case strings.Contains(s, "server 2012"):
		return "Windows Server 2012"
	case strings.Contains(s, "server 2008 r2"):
		return "Windows Server 2008 R2"
	case strings.Contains(s, "server 2008"):
		return "Windows Server 2008"
	case strings.Contains(s, "server 2003"):
		return "Windows Server 2003"
	case strings.Contains(s, "server 2000"):
		return "Windows Server 2000"
	default:
		return ""
	}
}

// extractBuildNumber выдёргивает номер сборки из operatingSystemVersion вида "10.0 (19045)".
// Если распарсить не удалось — возвращает 0.
func extractBuildNumber(osVersion string) int {
	osVersion = strings.TrimSpace(osVersion)
	if osVersion == "" {
		return 0
	}
	parts := strings.FieldsFunc(osVersion, func(r rune) bool {
		return r < '0' || r > '9'
	})
	if len(parts) == 0 {
		return 0
	}
	buildStr := parts[len(parts)-1]
	n, err := strconv.Atoi(buildStr)
	if err != nil {
		return 0
	}
	return n
}
