package ad

func (r *Repository) AddUserGroup(domain string, UserDN string, GroupDN string) error {
	return r.ads[domain].AddUserToGroup(UserDN, GroupDN)

}
func (r *Repository) DelUserGroup(domain string, UserDN string, GroupDN string) error {
	return r.ads[domain].DelUserFromGroup(UserDN, GroupDN)

}
