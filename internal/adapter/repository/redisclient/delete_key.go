package redisclient

import "errors"

func (r *Repository) Delete(key string) error {
	deleted, err := r.cl.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return errors.New("error deleting key from redis")
	}
	return nil
}
