package redisclient

import "time"

func (r *Repository) AddKeyValue(key string, field interface{}, td time.Duration) error {
	return r.cl.Set(ctx, key, field, td).Err()
}
