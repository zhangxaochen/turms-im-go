package mongo

type MongoCollectionMigrator struct {
}

// @MappedFrom migrate(Set<String> existingCollectionNames)
func (m *MongoCollectionMigrator) Migrate() error {
	return nil
}
