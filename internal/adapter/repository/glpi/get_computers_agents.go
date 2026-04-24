package glpi

import (
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetComputersAgents(domain string) (agentsMap map[string]entity.GLPI_Computer_Agent, err error) {
	agentsMap = make(map[string]entity.GLPI_Computer_Agent)
	agents := []entity.GLPI_Computer_Agent{}
	sql := fmt.Sprintf(`
	SELECT
     c.id,
     c.name,
     a.remote_addr AS ip,
     a.version,
     a.last_contact,
     d.name         AS domain
     FROM glpi_computers c
	JOIN glpi_agents a ON c.id = a.items_id
	LEFT JOIN glpi_domains_items di ON c.id = di.items_id AND di.itemtype = 'Computer'
	LEFT JOIN glpi_domains d ON d.id = di.domains_id
	WHERE a.itemtype = 'Computer' AND d.name='%s'`, domain)
	_, err = r.db.Select(&agents, sql)
	if err != nil {
		return nil, err
	}
	for _, agent := range agents {
		agentsMap[strings.ToUpper(agent.Name)] = agent
	}
	return agentsMap, err
}

func (r *Repository) GetComputersTags(domain string) (map[string]string, error) {
	tagsMap := make(map[string]string)
	computersTags := []struct {
		Name string `db:"name"`
		Tags string `db:"tags"`
	}{}
	sql := fmt.Sprintf(`
SELECT
    c.name,
    GROUP_CONCAT(DISTINCT t.name SEPARATOR ', ') AS tags
FROM glpi_computers c
LEFT JOIN glpi_domains_items di ON c.id = di.items_id AND di.itemtype = 'Computer'
LEFT JOIN glpi_domains d ON d.id = di.domains_id
JOIN glpi_plugin_tag_tagitems ti ON c.id = ti.items_id AND ti.itemtype = 'Computer'
LEFT JOIN glpi_plugin_tag_tags t ON ti.plugin_tag_tags_id = t.id
WHERE d.name = '%s'
GROUP BY c.name`, domain)
	_, err := r.db.Select(&computersTags, sql)
	if err != nil {
		return nil, err
	}
	for _, computer := range computersTags {
		tagsMap[strings.ToUpper(computer.Name)] = computer.Tags
	}
	return tagsMap, err
}
