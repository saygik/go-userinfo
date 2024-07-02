package ad

func (r *Repository) GetUser(domain string, user string) (map[string]interface{}, error) {
	return r.ads[domain].GetUserInfo(user)

}
