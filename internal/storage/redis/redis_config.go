package redis

type RedisConfig struct {
}

// @MappedFrom newSequenceIdRedisClientManager(RedisProperties properties)
func (c *RedisConfig) NewSequenceIdRedisClientManager() interface{} {
	return nil
}

// @MappedFrom sequenceIdRedisClientManager()
func (c *RedisConfig) SequenceIdRedisClientManager() interface{} {
	return nil
}
