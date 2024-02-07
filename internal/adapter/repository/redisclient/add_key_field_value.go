package redisclient

func (r *Repository) AddKeyFieldValue(key string, field string, value []byte) error {
	return r.cl.HSet(ctx, key, field, value).Err()
}
