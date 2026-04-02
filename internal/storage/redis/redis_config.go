package redis

type RedisConfig struct {
}

func (c *RedisConfig) NewSequenceIdRedisClientManager() interface{} {
	return nil
}

func (c *RedisConfig) SequenceIdRedisClientManager() interface{} {
	return nil
}
