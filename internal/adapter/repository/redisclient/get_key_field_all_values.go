package redisclient

func (r *Repository) GetKeyFieldAll(key string) (map[string]string, error) {
	return r.cl.HGetAll(ctx, key).Result()
}
