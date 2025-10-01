package redisclient

func (r *Repository) AddKeyFieldValue(key string, field string, value []byte) error {
	return r.cl.HSet(ctx, key, field, value).Err()
}

func (r *Repository) Rename(oldKey, newKey string) error {
	return r.cl.Rename(ctx, oldKey, newKey).Err()
}
func (r *Repository) RenameNX(oldKey, newKey string) (bool, error) {
	return r.cl.RenameNX(ctx, oldKey, newKey).Result()
}
func (r *Repository) Unlink(keys ...string) (int64, error) {
	return r.cl.Unlink(ctx, keys...).Result()
}
