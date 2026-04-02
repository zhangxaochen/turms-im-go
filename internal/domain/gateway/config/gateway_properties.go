package config

// GatewayProperties maps to turms properties for the gateway module.
type GatewayProperties struct {
	ClientAPI *ClientAPIProperties
}

// ClientAPIProperties maps client API configurations.
type ClientAPIProperties struct {
	ReturnReasonForServerError bool
}

// NewGatewayProperties creates default properties.
func NewGatewayProperties() *GatewayProperties {
	return &GatewayProperties{
		ClientAPI: &ClientAPIProperties{
			ReturnReasonForServerError: false, // Default matching Java behaviour via security
		},
	}
}
