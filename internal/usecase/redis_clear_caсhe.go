package usecase

func (u *UseCase) ClearRedisCaсhe() {
	//u.r.clearAllDomainsUsers()
	u.redis.ClearAllDomainsUsers()
	adl := u.ad.DomainList()
	for _, one := range adl {
		u.redis.DelKeyField("ad", one.Name)
		u.redis.DelKeyField("adc", one.Name)
	}
}
