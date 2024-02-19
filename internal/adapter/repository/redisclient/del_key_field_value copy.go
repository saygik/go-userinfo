package redisclient

func (r *Repository) DelKeyField(key string, field string) error {
	return r.cl.HDel(ctx, key, field).Err()
}
