package redisclient

import "errors"

func (r *Repository) Delete(key string) error {
	//	r.cl.Get(ctx, key).Result()

	dd, _ := r.GetKeyValue(key)
	ff := dd
	_ = ff
	deleted, err := r.cl.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return errors.New("error deleting key from redis")
	}
	return nil
}
