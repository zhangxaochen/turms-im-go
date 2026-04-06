package service

import (
	"fmt"
	"sort"
	"strings"
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

// UserDefinedAttributeProperties maps to UserDefinedAttributeProperties.java
type UserDefinedAttributeProperties struct {
	SourceName string
	StoredName string
	Immutable  bool
	Value      interface{}
}

// UserDefinedAttributesProperties maps to UserDefinedAttributesProperties.java
type UserDefinedAttributesProperties struct {
	AllowedAttributes               []*UserDefinedAttributeProperties
	IgnoreUnknownAttributesOnUpsert bool
}

// UserDefinedAttributesService maps to UserDefinedAttributesService.java
// @MappedFrom UserDefinedAttributesService
type UserDefinedAttributesService struct {
	knownAttributes                 map[string]struct{}
	sourceNameToAttributeProperties map[string]*UserDefinedAttributeProperties
	immutableAttributes             map[string]struct{}
	ignoreUnknownAttributesOnUpsert bool
}

// NewUserDefinedAttributesService creates a new UserDefinedAttributesService.
func NewUserDefinedAttributesService() *UserDefinedAttributesService {
	return &UserDefinedAttributesService{
		knownAttributes:                 make(map[string]struct{}),
		sourceNameToAttributeProperties: make(map[string]*UserDefinedAttributeProperties),
		immutableAttributes:             make(map[string]struct{}),
	}
}

// @MappedFrom updateGlobalProperties(UserDefinedAttributesProperties properties)
func (s *UserDefinedAttributesService) UpdateGlobalProperties(properties *UserDefinedAttributesProperties) error {
	if properties == nil {
		return nil
	}

	attributeList := properties.AllowedAttributes
	if len(attributeList) == 0 {
		return nil
	}

	newSourceNameToAttrProps := make(map[string]*UserDefinedAttributeProperties, len(attributeList))
	newImmutableAttributes := make(map[string]struct{})

	for _, attr := range attributeList {
		sourceName := attr.SourceName
		storedName := attr.StoredName

		// Bug fix: Create a copy instead of mutating the input attr directly.
		// Java builds a new object via .toBuilder().storedName(sourceName).build().
		attrCopy := &UserDefinedAttributeProperties{
			SourceName: attr.SourceName,
			StoredName: attr.StoredName,
			Immutable:  attr.Immutable,
			Value:      attr.Value,
		}

		if storedName == "" {
			attrCopy.StoredName = sourceName
			if _, exists := newSourceNameToAttrProps[sourceName]; exists {
				return fmt.Errorf("found a duplicate attribute: %s", sourceName)
			}
			newSourceNameToAttrProps[sourceName] = attrCopy
		} else {
			if _, exists := newSourceNameToAttrProps[storedName]; exists {
				return fmt.Errorf("found a duplicate attribute: %s", storedName)
			}
			newSourceNameToAttrProps[storedName] = attrCopy
		}

		if attr.Immutable {
			newImmutableAttributes[sourceName] = struct{}{}
		}
	}

	s.sourceNameToAttributeProperties = newSourceNameToAttrProps
	s.knownAttributes = make(map[string]struct{}, len(newSourceNameToAttrProps))
	for k := range newSourceNameToAttrProps {
		s.knownAttributes[k] = struct{}{}
	}
	s.immutableAttributes = newImmutableAttributes
	s.ignoreUnknownAttributesOnUpsert = properties.IgnoreUnknownAttributesOnUpsert
	return nil
}

// parseValue is a placeholder for CustomValueService.parseValue logic
func (s *UserDefinedAttributesService) parseValue(valType interface{}, name string, value interface{}) interface{} {
	// Simplified parsing map
	return value
}

// @MappedFrom parseAttributes
func (s *UserDefinedAttributesService) parseAttributes(ignoreUnknownAttributes bool, inputAttributes map[string]interface{}) (map[string]interface{}, error) {
	if len(inputAttributes) == 0 {
		return make(map[string]interface{}), nil
	}

	if len(s.sourceNameToAttributeProperties) == 0 {
		if ignoreUnknownAttributes {
			return make(map[string]interface{}), nil
		}
		var keys []string
		for k := range inputAttributes {
			keys = append(keys, k)
		}
		// Bug fix: Sort keys for deterministic error messages (Java uses TreeSet).
		sort.Strings(keys)
		return nil, fmt.Errorf("unknown attributes: %v", keys) // Should use exception.NewTurmsError usually
	}

	outputAttributes := make(map[string]interface{}, len(inputAttributes))

	if len(inputAttributes) <= len(s.sourceNameToAttributeProperties) {
		for sourceName, value := range inputAttributes {
			attrProps, exists := s.sourceNameToAttributeProperties[sourceName]
			if !exists {
				if ignoreUnknownAttributes {
					continue
				}
				return nil, fmt.Errorf("unknown attribute: %s", sourceName)
			}
			outputAttributes[attrProps.StoredName] = s.parseValue(attrProps.Value, sourceName, value)
		}
	} else {
		if !ignoreUnknownAttributes {
			var unknownKeys []string
			for k := range inputAttributes {
				if _, exists := s.sourceNameToAttributeProperties[k]; !exists {
					unknownKeys = append(unknownKeys, k)
				}
			}
			// Bug fix: Sort unknown keys for deterministic error messages (Java uses TreeSet).
			sort.Strings(unknownKeys)
			return nil, fmt.Errorf("unknown attributes: %v", unknownKeys) // Usually an error
		}

		for sourceName, attrProps := range s.sourceNameToAttributeProperties {
			val, exists := inputAttributes[sourceName]
			if !exists || val == nil {
				continue
			}
			outputAttributes[attrProps.StoredName] = s.parseValue(attrProps.Value, sourceName, val)
		}
	}

	return outputAttributes, nil
}

// findUserDefinedAttributes is an abstract hook mimicking Java's Mono<List<String>> findUserDefinedAttributes
func (s *UserDefinedAttributesService) findUserDefinedAttributes(immutableAttributesForUpsert []string) ([]string, error) {
	// Must be implemented by subclasses.
	return nil, nil // Return no attributes block
}

// @MappedFrom parseAttributesForUpsert(Map<String, Value> userDefinedAttributes)
func (s *UserDefinedAttributesService) ParseAttributesForUpsert(userDefinedAttributes map[string]interface{}) (map[string]interface{}, error) {
	if userDefinedAttributes == nil {
		return nil, fmt.Errorf("userDefinedAttributes must not be null")
	}
	if len(userDefinedAttributes) == 0 {
		return make(map[string]interface{}), nil
	}

	var immutableAttributesForUpsert []string
	if len(s.immutableAttributes) > 0 {
		for name := range userDefinedAttributes {
			if _, isImmutable := s.immutableAttributes[name]; isImmutable {
				immutableAttributesForUpsert = append(immutableAttributesForUpsert, name)
			}
		}
	}

	if len(immutableAttributesForUpsert) == 0 {
		return s.parseAttributes(s.ignoreUnknownAttributesOnUpsert, userDefinedAttributes)
	}

	existingAttributes, err := s.findUserDefinedAttributes(immutableAttributesForUpsert)
	if err != nil {
		return nil, err
	}

	if len(existingAttributes) == 0 {
		return s.parseAttributes(s.ignoreUnknownAttributesOnUpsert, userDefinedAttributes)
	}

	existingMap := make(map[string]struct{}, len(existingAttributes))
	for _, ea := range existingAttributes {
		existingMap[ea] = struct{}{}
	}

	var conflicted []string
	for _, curAttr := range immutableAttributesForUpsert {
		if _, exists := existingMap[curAttr]; exists {
			conflicted = append(conflicted, curAttr)
		}
	}

	if len(conflicted) > 0 {
		// Bug fix: Sort conflicted attributes for deterministic error messages (Java sorts).
		sort.Strings(conflicted)
		return nil, fmt.Errorf("cannot update existing immutable attributes: %s", strings.Join(conflicted, ", "))
	}

	return s.parseAttributes(s.ignoreUnknownAttributesOnUpsert, userDefinedAttributes)
}
