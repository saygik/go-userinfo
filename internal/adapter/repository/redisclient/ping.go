package redisclient

// Ping проверяет доступность Redis.
func (r *Repository) Ping() error {
	return r.cl.Ping(ctx).Err()
}
