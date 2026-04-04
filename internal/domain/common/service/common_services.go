package service

import (
	"time"

	"im.turms/server/internal/domain/common/repository"
)

// ExpirableEntityService maps to ExpirableEntityService.java
// @MappedFrom ExpirableEntityService
type ExpirableEntityService struct {
	ExpirableEntityRepository *repository.ExpirableEntityRepository
}

// @MappedFrom getEntityExpirationDate()
func (s *ExpirableEntityService) GetEntityExpirationDate() *time.Time {
	if s.ExpirableEntityRepository == nil {
		return nil
	}
	return s.ExpirableEntityRepository.GetEntityExpirationDate()
}

// UserDefinedAttributesService maps to UserDefinedAttributesService.java
// @MappedFrom UserDefinedAttributesService
type UserDefinedAttributesService struct {
	knownAttributes                 map[string]interface{}
	sourceNameToAttributeProperties map[string]interface{}
	immutableAttributes             []string
	ignoreUnknownAttributesOnUpsert bool
}

// @MappedFrom updateGlobalProperties(UserDefinedAttributesProperties properties)
func (s *UserDefinedAttributesService) UpdateGlobalProperties(properties interface{}) {
	// Logic to update global properties (to be fully implemented with specific types)
}

// @MappedFrom parseAttributesForUpsert(Map<String, Value> userDefinedAttributes)
func (s *UserDefinedAttributesService) ParseAttributesForUpsert(userDefinedAttributes map[string]interface{}) map[string]interface{} {
	// Logic to parse attributes for upsert (to be fully implemented with specific types)
	return nil
}

func (s *UserDefinedAttributesService) parseAttributes(ignoreUnknownAttributes bool, inputAttributes map[string]interface{}) map[string]interface{} {
	return nil
}

func (s *UserDefinedAttributesService) findUserDefinedAttributes(immutableAttributesForUpsert []string) ([]string, error) {
	return nil, nil
}
