package redisclient

func (r *Repository) GetKeyFieldValue(key string, field string) (string, error) {
	return r.cl.HGet(ctx, key, field).Result()
}

func (r *Repository) GetKeyFieldsValue(key string, fields []string) ([]interface{}, error) {
	return r.cl.HMGet(ctx, key, fields...).Result()
}
