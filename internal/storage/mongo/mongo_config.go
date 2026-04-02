package mongo

type MongoConfig struct {
}

// @MappedFrom adminMongoClient(TurmsPropertiesManager propertiesManager)
func (c *MongoConfig) AdminMongoClient() interface{} {
	return nil
}

// @MappedFrom userMongoClient(TurmsPropertiesManager propertiesManager)
func (c *MongoConfig) UserMongoClient() interface{} {
	return nil
}

// @MappedFrom groupMongoClient(TurmsPropertiesManager propertiesManager)
func (c *MongoConfig) GroupMongoClient() interface{} {
	return nil
}

// @MappedFrom conversationMongoClient(TurmsPropertiesManager propertiesManager)
func (c *MongoConfig) ConversationMongoClient() interface{} {
	return nil
}

// @MappedFrom messageMongoClient(TurmsPropertiesManager propertiesManager)
func (c *MongoConfig) MessageMongoClient() interface{} {
	return nil
}

// @MappedFrom conferenceMongoClient(TurmsPropertiesManager propertiesManager)
func (c *MongoConfig) ConferenceMongoClient() interface{} {
	return nil
}
