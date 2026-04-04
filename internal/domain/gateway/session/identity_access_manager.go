package session

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	adminconstant "im.turms/server/internal/domain/admin/constant"
	commonconstant "im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common/authorization"
	"im.turms/server/internal/domain/gateway/config"
	"im.turms/server/internal/domain/gateway/session/bo"
	userservice "im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/infra/ldap"
)

// SessionIdentityAccessManagementSupport maps to the Java support interface.
type SessionIdentityAccessManagementSupport interface {
	VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error)
	UpdateGlobalProperties(properties interface{})
}

// SessionIdentityAccessManager orchestrates the authentication.
// @MappedFrom SessionIdentityAccessManager
type SessionIdentityAccessManager struct {
	enableIdentityAccessManagement bool
	support                        SessionIdentityAccessManagementSupport
	userService                    userservice.UserService
	policyManager                  *authorization.PolicyManager
	pluginManager                  *plugin.PluginManager
}

func NewSessionIdentityAccessManager(
	gatewayProperties *config.GatewayProperties,
	userService userservice.UserService,
	pluginManager *plugin.PluginManager,
) *SessionIdentityAccessManager {
	iamProps := gatewayProperties.IdentityAccessManagement
	policyManager := &authorization.PolicyManager{}
	support := createSupport(iamProps, userService, policyManager)

	return &SessionIdentityAccessManager{
		enableIdentityAccessManagement: iamProps.Enabled,
		support:                        support,
		userService:                    userService,
		policyManager:                  policyManager,
		pluginManager:                  pluginManager,
	}
}

func createSupport(
	iamProps *config.IdentityAccessManagementProperties,
	userService userservice.UserService,
	policyManager *authorization.PolicyManager,
) SessionIdentityAccessManagementSupport {
	switch iamProps.Type {
	case config.IdentityAccessManagementType_PASSWORD:
		return &PasswordSessionIdentityAccessManager{
			userService: userService,
		}
	case config.IdentityAccessManagementType_HTTP:
		return NewHttpSessionIdentityAccessManager(&iamProps.Http, policyManager)
	case config.IdentityAccessManagementType_JWT:
		return NewJwtSessionIdentityAccessManager(&iamProps.Jwt, policyManager)
	case config.IdentityAccessManagementType_LDAP:
		return NewLdapSessionIdentityAccessManager(&iamProps.Ldap, policyManager)
	case config.IdentityAccessManagementType_NOOP:
		return &NoopSessionIdentityAccessManager{}
	default:
		return &NoopSessionIdentityAccessManager{}
	}
}

func (m *SessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {
	props, ok := properties.(*config.GatewayProperties)
	if !ok {
		return
	}
	iamProps := props.IdentityAccessManagement
	m.enableIdentityAccessManagement = iamProps.Enabled
	if m.support != nil {
		m.support.UpdateGlobalProperties(&iamProps)
	}
}

// @MappedFrom verifyAndGrant(int version, Long userId, String password, DeviceType deviceType, Map<String, String> deviceDetails, UserStatus userStatus, Location location, String ip)
func (m *SessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	if loginInfo.UserID == adminconstant.AdminRequesterId {
		return bo.LoginAuthenticationFailed, nil
	}

	if !m.enableIdentityAccessManagement {
		return bo.GrantedWithAllPermissions, nil
	}

	// 1. Check PluginManager for UserAuthenticator hooks
	if m.pluginManager != nil && m.pluginManager.HasRunningExtensions("UserAuthenticator") {
		ok, err := m.pluginManager.InvokeExtensionPoints(ctx, "UserAuthenticator", "Authenticate", loginInfo)
		if err != nil {
			return nil, err
		}
		if ok {
			return bo.GrantedWithAllPermissions, nil
		}
		return bo.LoginAuthenticationFailed, nil
	}

	if m.support != nil {
		return m.support.VerifyAndGrant(ctx, loginInfo)
	}

	return bo.GrantedWithAllPermissions, nil
}

// PasswordSessionIdentityAccessManager maps to PasswordSessionIdentityAccessManager in Java.
type PasswordSessionIdentityAccessManager struct {
	userService userservice.UserService
}

// @MappedFrom verifyAndGrant(UserLoginInfo userLoginInfo)
func (m *PasswordSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	if loginInfo.Password == nil || *loginInfo.Password == "" {
		return bo.LoginAuthenticationFailed, nil
	}

	password, err := m.userService.FindPassword(ctx, loginInfo.UserID)
	if err != nil {
		return nil, err
	}
	if password == nil {
		return bo.LoginAuthenticationFailed, nil
	}

	active, err := m.userService.IsActiveAndNotDeleted(ctx, loginInfo.UserID)
	if err != nil {
		return nil, err
	}
	if !active {
		return bo.LoggingInUserNotActive, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(*password), []byte(*loginInfo.Password))
	if err != nil {
		return bo.LoginAuthenticationFailed, nil
	}

	return bo.GrantedWithAllPermissions, nil
}

func (m *PasswordSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {
}

// HttpSessionIdentityAccessManager maps to HttpSessionIdentityAccessManager in Java.
type HttpSessionIdentityAccessManager struct {
	client              *http.Client
	url                 string
	method              string
	expectedStatusCodes map[string]struct{}
	expectedHeaders     map[string]string
	expectedBodyFields  map[string]interface{}
	policyManager       *authorization.PolicyManager
	policyDeserializer  *authorization.PolicyDeserializer
}

func NewHttpSessionIdentityAccessManager(
	props *config.HttpIdentityAccessManagementProperties,
	policyManager *authorization.PolicyManager,
) *HttpSessionIdentityAccessManager {
	reqProps := props.Request
	respProps := props.Authentication.ResponseExpectation

	statusCodes := make(map[string]struct{})
	for _, sc := range respProps.StatusCodes {
		statusCodes[sc] = struct{}{}
	}

	return &HttpSessionIdentityAccessManager{
		client: &http.Client{
			Timeout: time.Duration(reqProps.TimeoutMillis) * time.Millisecond,
		},
		url:                 reqProps.URL,
		method:              string(reqProps.HttpMethod),
		expectedStatusCodes: statusCodes,
		expectedHeaders:     respProps.Headers,
		expectedBodyFields:  respProps.BodyFields,
		policyManager:       policyManager,
		policyDeserializer:  &authorization.PolicyDeserializer{},
	}
}

func (m *HttpSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	body, err := json.Marshal(loginInfo)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, m.method, m.url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return bo.LoginAuthenticationFailed, nil
	}
	defer resp.Body.Close()

	// Check status codes
	statusStr := fmt.Sprintf("%d", resp.StatusCode)
	if _, ok := m.expectedStatusCodes[statusStr]; !ok {
		return bo.LoginAuthenticationFailed, nil
	}

	// Check headers
	for k, v := range m.expectedHeaders {
		if resp.Header.Get(k) != v {
			return bo.LoginAuthenticationFailed, nil
		}
	}

	// Read body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return bo.LoginAuthenticationFailed, nil
	}

	var respMap map[string]interface{}
	if err := json.Unmarshal(respBody, &respMap); err != nil {
		return bo.LoginAuthenticationFailed, nil
	}

	// Check body fields
	for k, v := range m.expectedBodyFields {
		if val, ok := respMap[k]; !ok || fmt.Sprint(val) != fmt.Sprint(v) {
			return bo.LoginAuthenticationFailed, nil
		}
	}

	// Decode policy
	policy, err := m.policyDeserializer.Parse(respMap)
	if err != nil {
		return nil, err
	}

	allowedRequestTypes := m.policyManager.FindAllowedRequestTypes(policy)
	permissions := make(map[any]bool, len(allowedRequestTypes))
	for _, rt := range allowedRequestTypes {
		permissions[rt] = true
	}
	return bo.NewUserPermissionInfo(commonconstant.ResponseStatusCode_OK, permissions), nil
}

func (m *HttpSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {
	props, ok := properties.(*config.IdentityAccessManagementProperties)
	if !ok {
		return
	}
	m.url = props.Http.Request.URL
	m.method = string(props.Http.Request.HttpMethod)
	m.client.Timeout = time.Duration(props.Http.Request.TimeoutMillis) * time.Millisecond

	m.expectedStatusCodes = make(map[string]struct{})
	for _, sc := range props.Http.Authentication.ResponseExpectation.StatusCodes {
		m.expectedStatusCodes[sc] = struct{}{}
	}
	m.expectedHeaders = props.Http.Authentication.ResponseExpectation.Headers
	m.expectedBodyFields = props.Http.Authentication.ResponseExpectation.BodyFields
}

// JwtSessionIdentityAccessManager maps to JwtSessionIdentityAccessManager in Java.
type JwtSessionIdentityAccessManager struct {
	algorithm      string
	secretKey      []byte
	policyManager  *authorization.PolicyManager
	expectedClaims map[string]interface{}
}

func NewJwtSessionIdentityAccessManager(
	props *config.JwtIdentityAccessManagementProperties,
	policyManager *authorization.PolicyManager,
) *JwtSessionIdentityAccessManager {
	return &JwtSessionIdentityAccessManager{
		algorithm:     props.Algorithm,
		secretKey:     []byte(props.SecretKey),
		policyManager: policyManager,
		// In a real system, these would be populated from properties too
		expectedClaims: make(map[string]interface{}),
	}
}

func (m *JwtSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	if loginInfo.Password == nil || *loginInfo.Password == "" {
		return bo.LoginAuthenticationFailed, nil
	}

	tokenString := *loginInfo.Password
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check algorithm
		if token.Method.Alg() != m.algorithm {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil || !token.Valid {
		return bo.LoginAuthenticationFailed, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return bo.LoginAuthenticationFailed, nil
	}

	// Validate userId (subject)
	sub, ok := claims["sub"].(string)
	if !ok || sub != fmt.Sprintf("%d", loginInfo.UserID) {
		return bo.LoginAuthenticationFailed, nil
	}

	// Validate expected custom claims
	for k, v := range m.expectedClaims {
		if cv, ok := claims[k]; !ok || fmt.Sprint(cv) != fmt.Sprint(v) {
			return bo.LoginAuthenticationFailed, nil
		}
	}

	// Check policy if it exists in claims
	if policyData, ok := claims["policy"].(map[string]interface{}); ok {
		policyDeserializer := &authorization.PolicyDeserializer{}
		policy, err := policyDeserializer.Parse(policyData)
		if err == nil {
			allowedRequestTypes := m.policyManager.FindAllowedRequestTypes(policy)
			permissions := make(map[any]bool, len(allowedRequestTypes))
			for _, rt := range allowedRequestTypes {
				permissions[rt] = true
			}
			return bo.NewUserPermissionInfo(commonconstant.ResponseStatusCode_OK, permissions), nil
		}
	}

	// Granted with all permissions by default if valid
	return bo.GrantedWithAllPermissions, nil
}

func (m *JwtSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {
	props, ok := properties.(*config.IdentityAccessManagementProperties)
	if !ok {
		return
	}
	m.algorithm = props.Jwt.Algorithm
	m.secretKey = []byte(props.Jwt.SecretKey)
}

// LdapSessionIdentityAccessManager maps to LdapSessionIdentityAccessManager in Java.
type LdapSessionIdentityAccessManager struct {
	client        *ldap.LdapClient
	baseDN        string
	userFilter    string
	policyManager *authorization.PolicyManager
}

func NewLdapSessionIdentityAccessManager(
	props *config.LdapIdentityAccessManagementProperties,
	policyManager *authorization.PolicyManager,
) *LdapSessionIdentityAccessManager {
	return &LdapSessionIdentityAccessManager{
		client:        ldap.NewLdapClient(props.URL, false, true), // Defaulting to no TLS for now
		baseDN:        props.BaseDN,
		userFilter:    props.UserFilter,
		policyManager: policyManager,
	}
}

func (m *LdapSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	if loginInfo.Password == nil || *loginInfo.Password == "" {
		return bo.LoginAuthenticationFailed, nil
	}

	if !m.client.IsConnected() {
		if err := m.client.Connect(); err != nil {
			return bo.LoginAuthenticationFailed, fmt.Errorf("failed to connect to LDAP: %w", err)
		}
	}

	// 1. Search for the user DN using userId and userFilter
	filter := strings.ReplaceAll(m.userFilter, "{0}", fmt.Sprintf("%d", loginInfo.UserID))
	searchResult, err := m.client.Search(m.baseDN, 2, 0, 1, 0, false, []string{"dn"}, filter)
	if err != nil || len(searchResult.Entries) != 1 {
		return bo.LoginAuthenticationFailed, nil
	}
	userDN := searchResult.Entries[0].DN

	// 2. Bind with the user DN and password
	ok, err := m.client.Bind(false, userDN, *loginInfo.Password)
	if err != nil || !ok {
		return bo.LoginAuthenticationFailed, nil
	}

	return bo.GrantedWithAllPermissions, nil
}

func (m *LdapSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {
	props, ok := properties.(*config.IdentityAccessManagementProperties)
	if !ok {
		return
	}
	m.baseDN = props.Ldap.BaseDN
	m.userFilter = props.Ldap.UserFilter
	if m.client.Addr != props.Ldap.URL {
		m.client.Close()
		m.client = ldap.NewLdapClient(props.Ldap.URL, false, true)
	}
}

// NoopSessionIdentityAccessManager maps to NoopSessionIdentityAccessManager in Java.
type NoopSessionIdentityAccessManager struct{}

func (m *NoopSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	// Java: returns GRANTED_WITH_ALL_PERMISSIONS (all permissions)
	return bo.GrantedWithAllPermissions, nil
}
func (m *NoopSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {}
