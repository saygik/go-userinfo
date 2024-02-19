package redisclient

func (r *Repository) GetKeyValue(key string) (string, error) {
	return r.cl.Get(ctx, key).Result()
}
