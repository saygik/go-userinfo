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
