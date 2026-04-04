package mongo

import (
	"context"
	"fmt"
	"log/slog"
)

// TurmsMongoClient is a stub for the MongoDB client used in the Turms gateway.
// @MappedFrom im.turms.server.common.storage.mongo.TurmsMongoClient
type TurmsMongoClient struct {
	name         string
	uri          string
	clusterTypes []string // allowed cluster types, empty = any
}

// RegisterEntitiesByClasses registers entity classes for the mongo client.
// @MappedFrom registerEntitiesByClasses(Class<?>... classes)
func (c *TurmsMongoClient) RegisterEntitiesByClasses(classes ...string) {
	slog.Debug("MongoClient registered entities", "client", c.name, "classes", classes)
}

// Close releases the mongo client resources.
func (c *TurmsMongoClient) Close() {
	slog.Info("MongoClient closed", "client", c.name)
}

// Ping verifies the connection.
func (c *TurmsMongoClient) Ping(ctx context.Context) error {
	return nil
}

// MongoProperties approximates im.turms.server.common.infra.property.env.common.mongo.MongoProperties.
type MongoProperties struct {
	URI string
}

// IdentityAccessManagementType approximates the Java enum.
type IdentityAccessManagementType int

const (
	IdentityAccessManagementTypePassword IdentityAccessManagementType = iota
	IdentityAccessManagementTypeLDAP
	IdentityAccessManagementTypeHTTP
	IdentityAccessManagementTypeJWT
	IdentityAccessManagementTypeNoop
)

// GatewayMongoProperties groups mongo configs for the gateway.
type GatewayMongoProperties struct {
	Admin MongoProperties
	User  MongoProperties
}

// SessionIdentityAccessManagementProperties holds IAM settings.
type SessionIdentityAccessManagementProperties struct {
	Enabled bool
	Type    IdentityAccessManagementType
}

// MongoConfig holds factory logic for gateway Mongo clients.
// @MappedFrom im.turms.gateway.storage.mongo.MongoConfig (Spring @Configuration)
type MongoConfig struct {
	adminProps MongoProperties
	userProps  MongoProperties
	iamProps   SessionIdentityAccessManagementProperties
}

// NewMongoConfig creates a MongoConfig from extracted gateway properties.
func NewMongoConfig(
	adminProps MongoProperties,
	userProps MongoProperties,
	iamProps SessionIdentityAccessManagementProperties,
) *MongoConfig {
	return &MongoConfig{
		adminProps: adminProps,
		userProps:  userProps,
		iamProps:   iamProps,
	}
}

// getMongoClient builds a TurmsMongoClient from the given properties.
// @MappedFrom BaseMongoConfig.getMongoClient(MongoProperties, String, Set<ClusterType>)
func getMongoClient(props MongoProperties, name string, allowedClusterTypes []string) (*TurmsMongoClient, error) {
	if props.URI == "" {
		return nil, fmt.Errorf("mongodb URI must not be empty for client %q", name)
	}
	slog.Info("Creating TurmsMongoClient", "name", name, "uri", props.URI)
	return &TurmsMongoClient{
		name:         name,
		uri:          props.URI,
		clusterTypes: allowedClusterTypes,
	}, nil
}

// AdminMongoClient creates the admin MongoDB client and registers Admin and AdminRole entities.
// @MappedFrom adminMongoClient(TurmsPropertiesManager propertiesManager)
func (c *MongoConfig) AdminMongoClient() (*TurmsMongoClient, error) {
	client, err := getMongoClient(c.adminProps, "admin", nil)
	if err != nil {
		return nil, err
	}
	client.RegisterEntitiesByClasses("Admin", "AdminRole")
	return client, nil
}

// UserMongoClient creates the user MongoDB client and registers the User entity,
// but only if IAM is enabled and the type is PASSWORD.
// Returns nil (no error) if IAM is disabled or type is not PASSWORD.
// @MappedFrom userMongoClient(TurmsPropertiesManager propertiesManager)
func (c *MongoConfig) UserMongoClient() (*TurmsMongoClient, error) {
	if !c.iamProps.Enabled || c.iamProps.Type != IdentityAccessManagementTypePassword {
		// Intentionally nil – matches Java @Nullable @Bean behaviour
		return nil, nil
	}
	client, err := getMongoClient(
		c.userProps,
		"user",
		[]string{"SHARDED", "LOAD_BALANCED"}, // matches Java Set.of(ClusterType.SHARDED, ClusterType.LOAD_BALANCED)
	)
	if err != nil {
		return nil, err
	}
	client.RegisterEntitiesByClasses("User")
	return client, nil
}

// GroupMongoClient stub – not part of the gateway module.
func (c *MongoConfig) GroupMongoClient() interface{} { return nil }

// ConversationMongoClient stub – not part of the gateway module.
func (c *MongoConfig) ConversationMongoClient() interface{} { return nil }

// MessageMongoClient stub – not part of the gateway module.
func (c *MongoConfig) MessageMongoClient() interface{} { return nil }

// ConferenceMongoClient stub – not part of the gateway module.
func (c *MongoConfig) ConferenceMongoClient() interface{} { return nil }
