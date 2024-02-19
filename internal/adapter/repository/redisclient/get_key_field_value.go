package redisclient

func (r *Repository) GetKeyFieldValue(key string, field string) (string, error) {
	return r.cl.HGet(ctx, key, field).Result()
}
