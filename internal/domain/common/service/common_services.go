package service

// ExpirableEntityService maps to ExpirableEntityService.java
// @MappedFrom ExpirableEntityService
type ExpirableEntityService struct {
}

// @MappedFrom getEntityExpirationDate()
func (s *ExpirableEntityService) GetEntityExpirationDate() {
}

// UserDefinedAttributesService maps to UserDefinedAttributesService.java
// @MappedFrom UserDefinedAttributesService
type UserDefinedAttributesService struct {
}

// @MappedFrom updateGlobalProperties(UserDefinedAttributesProperties properties)
func (s *UserDefinedAttributesService) UpdateGlobalProperties() {
}

// @MappedFrom parseAttributesForUpsert(Map<String, Value> userDefinedAttributes)
func (s *UserDefinedAttributesService) ParseAttributesForUpsert() {
}
