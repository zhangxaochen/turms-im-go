package config

import (
	"im.turms/server/internal/infra/property"
)

type IdentityAccessManagementProperties struct {
	Enabled bool                                   `json:"enabled" yaml:"enabled"`
	Type    IdentityAccessManagementType           `json:"type" yaml:"type"`
	Http    HttpIdentityAccessManagementProperties `json:"http" yaml:"http"`
	Jwt     JwtIdentityAccessManagementProperties  `json:"jwt" yaml:"jwt"`
	Ldap    LdapIdentityAccessManagementProperties `json:"ldap" yaml:"ldap"`
}

type IdentityAccessManagementType string

const (
	IdentityAccessManagementType_PASSWORD IdentityAccessManagementType = "PASSWORD"
	IdentityAccessManagementType_HTTP     IdentityAccessManagementType = "HTTP"
	IdentityAccessManagementType_JWT      IdentityAccessManagementType = "JWT"
	IdentityAccessManagementType_LDAP     IdentityAccessManagementType = "LDAP"
	IdentityAccessManagementType_NOOP     IdentityAccessManagementType = "NOOP"
)

type HttpIdentityAccessManagementProperties struct {
	Request        HttpIdentityAccessManagementRequestProperties `json:"request" yaml:"request"`
	Authentication HttpAuthenticationProperties                  `json:"authentication" yaml:"authentication"`
}

type HttpIdentityAccessManagementRequestProperties struct {
	URL           string              `json:"url" yaml:"url"`
	Headers       map[string]string   `json:"headers" yaml:"headers"`
	TimeoutMillis int                 `json:"timeoutMillis" yaml:"timeoutMillis"`
	HttpMethod    property.HttpMethod `json:"httpMethod" yaml:"httpMethod"`
}

type HttpAuthenticationProperties struct {
	ResponseExpectation HttpAuthenticationResponseExpectationProperties `json:"responseExpectation" yaml:"responseExpectation"`
}

type HttpAuthenticationResponseExpectationProperties struct {
	StatusCodes []string               `json:"statusCodes" yaml:"statusCodes"`
	Headers     map[string]string      `json:"headers" yaml:"headers"`
	BodyFields  map[string]interface{} `json:"bodyFields" yaml:"bodyFields"`
}

type JwtIdentityAccessManagementProperties struct {
	Algorithm      string                      `json:"algorithm" yaml:"algorithm"`
	SecretKey      string                      `json:"secretKey" yaml:"secretKey"`
	Authentication JwtAuthenticationProperties `json:"authentication" yaml:"authentication"`
}

type JwtAuthenticationProperties struct {
	Expectation JwtAuthenticationExpectationProperties `json:"expectation" yaml:"expectation"`
}

type JwtAuthenticationExpectationProperties struct {
	CustomPayloadClaims map[string]interface{} `json:"customPayloadClaims" yaml:"customPayloadClaims"`
}

// SearchFilterPlaceholderUserID is the placeholder in the user search filter
// that will be replaced with the actual user ID.
// @MappedFrom LdapIdentityAccessManagementUserProperties.SEARCH_FILTER_PLACEHOLDER_USER_ID
const SearchFilterPlaceholderUserID = "${userId}"

// LdapIdentityAccessManagementProperties mirrors Java's LdapIdentityAccessManagementProperties.
// @MappedFrom LdapIdentityAccessManagementProperties
type LdapIdentityAccessManagementProperties struct {
	BaseDN string                                      `json:"baseDn" yaml:"baseDn"`
	Admin  LdapIdentityAccessManagementAdminProperties `json:"admin" yaml:"admin"`
	User   LdapIdentityAccessManagementUserProperties  `json:"user" yaml:"user"`
}

// LdapIdentityAccessManagementAdminProperties mirrors Java's LdapIdentityAccessManagementAdminProperties.
// @MappedFrom LdapIdentityAccessManagementAdminProperties
type LdapIdentityAccessManagementAdminProperties struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	UseTLS   bool   `json:"useTls" yaml:"useTls"`
}

// LdapIdentityAccessManagementUserProperties mirrors Java's LdapIdentityAccessManagementUserProperties.
// @MappedFrom LdapIdentityAccessManagementUserProperties
type LdapIdentityAccessManagementUserProperties struct {
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	SearchFilter string `json:"searchFilter" yaml:"searchFilter"`
	UseTLS       bool   `json:"useTls" yaml:"useTls"`
}

func NewIdentityAccessManagementProperties() *IdentityAccessManagementProperties {
	return &IdentityAccessManagementProperties{
		Type: IdentityAccessManagementType_PASSWORD,
		Http: HttpIdentityAccessManagementProperties{
			Request: HttpIdentityAccessManagementRequestProperties{
				TimeoutMillis: 5000,
				HttpMethod:    property.HttpMethod_POST,
			},
			Authentication: HttpAuthenticationProperties{
				ResponseExpectation: HttpAuthenticationResponseExpectationProperties{
					StatusCodes: []string{"200"},
				},
			},
		},
	}
}
