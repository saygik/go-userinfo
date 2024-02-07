package redisclient

func (r *Repository) ClearAllDomainsUsers() {
	r.cl.Del(ctx, "ad")
	r.cl.Del(ctx, "adc")
}
