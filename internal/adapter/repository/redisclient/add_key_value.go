package redisclient

import "time"

func (r *Repository) AddKeyValue(key string, field interface{}, td time.Duration) error {
	return r.cl.Set(ctx, key, field, td).Err()
}

// func (r *Repository) GetKeyValue(key string) (interface{}, error) {
// 	return r.cl.Get(ctx, key).Result()
// }
