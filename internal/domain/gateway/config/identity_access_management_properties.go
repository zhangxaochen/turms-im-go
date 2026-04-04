package config

import (
	"im.turms/server/internal/infra/property"
)

type IdentityAccessManagementProperties struct {
	Type   IdentityAccessManagementType        `json:"type" yaml:"type"`
	Http   HttpIdentityAccessManagementProperties `json:"http" yaml:"http"`
	Jwt    JwtIdentityAccessManagementProperties  `json:"jwt" yaml:"jwt"`
	Ldap   LdapIdentityAccessManagementProperties `json:"ldap" yaml:"ldap"`
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
	URL           string            `json:"url" yaml:"url"`
	Headers       map[string]string `json:"headers" yaml:"headers"`
	TimeoutMillis int               `json:"timeoutMillis" yaml:"timeoutMillis"`
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
	Algorithm string `json:"algorithm" yaml:"algorithm"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type LdapIdentityAccessManagementProperties struct {
	URL        string `json:"url" yaml:"url"`
	BaseDN     string `json:"baseDN" yaml:"baseDN"`
	UserFilter string `json:"userFilter" yaml:"userFilter"`
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
