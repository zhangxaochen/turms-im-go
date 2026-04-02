package mongo

type MongoConfig struct {
}

func (c *MongoConfig) AdminMongoClient() interface{} {
	return nil
}

func (c *MongoConfig) UserMongoClient() interface{} {
	return nil
}

func (c *MongoConfig) GroupMongoClient() interface{} {
	return nil
}

func (c *MongoConfig) ConversationMongoClient() interface{} {
	return nil
}

func (c *MongoConfig) MessageMongoClient() interface{} {
	return nil
}

func (c *MongoConfig) ConferenceMongoClient() interface{} {
	return nil
}
