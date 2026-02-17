package ad

import "strings"

func (r *Repository) GetDomainAdminsGLPI(domain string) int {
	if val, ok := r.adconfigs[domain]; ok {
		return val.AdminGLPIGroup
	} else {
		return 0
	}
}

// "mattermostLogChannel": "gqnnu5asjfdhzkysbrqxor1hde",
func (r *Repository) GetDomainMattermostLogChannel(domain string) string {
	if val, ok := r.adconfigs[domain]; ok {
		return val.MattermostLogChannel
	} else {
		return ""
	}
}

func (r *Repository) GetMattermostLogChannelsByPrefix(inputPath string) []string {

	var channels []string
	domains := r.adconfigs
	for _, config := range domains {
		if config == nil {
			continue
		}

		// Проверяем что GlpiRegionPrefix является префиксом inputPath
		prefix := strings.TrimSpace(config.GlpiRegionPrefix)
		if prefix != "" &&
			strings.HasPrefix(inputPath, prefix) &&
			config.MattermostLogChannel != "" {

			channels = append(channels, config.MattermostLogChannel)
		}
	}

	return channels
}
