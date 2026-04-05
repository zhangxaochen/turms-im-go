package session

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	adminconstant "im.turms/server/internal/domain/admin/constant"
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common/authorization"
	"im.turms/server/internal/domain/gateway/config"
	"im.turms/server/internal/domain/gateway/session/bo"
	userservice "im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/ldap"
	"im.turms/server/internal/infra/ldap/element"
	"im.turms/server/internal/infra/plugin"
)

// SessionIdentityAccessManagementSupport maps to the Java support interface.
type SessionIdentityAccessManagementSupport interface {
	VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error)
	UpdateGlobalProperties(properties interface{}) bool
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

	manager := &SessionIdentityAccessManager{
		support:       support,
		userService:   userService,
		policyManager: policyManager,
		pluginManager: pluginManager,
	}
	manager.enableIdentityAccessManagement = manager.support.UpdateGlobalProperties(&iamProps)
	return manager
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
	if m.support != nil {
		m.enableIdentityAccessManagement = m.support.UpdateGlobalProperties(&iamProps)
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
		result, err := m.pluginManager.InvokeExtensionPoints(ctx, "UserAuthenticator", "Authenticate", loginInfo)
		if err != nil {
			return nil, err
		}
		if result != nil {
			if *result {
				return bo.GrantedWithAllPermissions, nil
			}
			return bo.LoginAuthenticationFailed, nil
		}
		// If result == nil, fall through to default handler!
	}

	if m.support != nil {
		return m.support.VerifyAndGrant(ctx, loginInfo)
	}

	return bo.LoginAuthenticationFailed, nil
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

func (m *PasswordSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) bool {
	props, ok := properties.(*config.IdentityAccessManagementProperties)
	if !ok {
		return false
	}
	return props.Enabled
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
	permissions := make(map[int32]bool, len(allowedRequestTypes))
	for _, rt := range allowedRequestTypes {
		permissions[rt] = true
	}
	return bo.NewUserPermissionInfo(constant.ResponseStatusCode_OK, permissions), nil
}

func (m *HttpSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) bool {
	props, ok := properties.(*config.IdentityAccessManagementProperties)
	if !ok {
		return false
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

	return props.Enabled
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
			permissions := make(map[int32]bool, len(allowedRequestTypes))
			for _, rt := range allowedRequestTypes {
				permissions[rt] = true
			}
			return bo.NewUserPermissionInfo(constant.ResponseStatusCode_OK, permissions), nil
		}
	}

	// Granted with all permissions by default if valid
	return bo.GrantedWithAllPermissions, nil
}

func (m *JwtSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) bool {
	props, ok := properties.(*config.IdentityAccessManagementProperties)
	if !ok {
		return false
	}
	m.algorithm = props.Jwt.Algorithm
	m.secretKey = []byte(props.Jwt.SecretKey)

	return props.Enabled
}

// LdapSessionIdentityAccessManager maps to LdapSessionIdentityAccessManager in Java.
type LdapSessionIdentityAccessManager struct {
	mu          sync.RWMutex
	adminClient *ldap.LdapClient
	userClient  *ldap.LdapClient
	baseDN      string
	userFilter  string

	userBindMu    sync.Mutex
	policyManager *authorization.PolicyManager
}

func NewLdapSessionIdentityAccessManager(
	props *config.LdapIdentityAccessManagementProperties,
	policyManager *authorization.PolicyManager,
) *LdapSessionIdentityAccessManager {
	if !strings.Contains(props.User.SearchFilter, "{0}") {
		panic(fmt.Errorf("illegal argument: The User Search Filter must contain {0} to substitute the user ID: %s", props.User.SearchFilter))
	}

	adminClient, err := ldap.NewLdapClient(props.Admin.Host, props.Admin.Port, props.Admin.UseTLS, nil, 5*time.Second)
	if err != nil {
		panic(fmt.Errorf("illegal state: failed to connect to admin LDAP server: %w", err))
	}
	
	userClient, err := ldap.NewLdapClient(props.User.Host, props.User.Port, props.User.UseTLS, nil, 5*time.Second)
	if err != nil {
		panic(fmt.Errorf("illegal state: failed to connect to user LDAP server: %w", err))
	}

	// Java parity: bind the admin client once at startup
	if props.Admin.Username != "" {
		ok, err := adminClient.Bind(false, props.Admin.Username, props.Admin.Password)
		if err != nil || !ok {
			panic(fmt.Errorf("illegal state: admin bind failed (ok=%v): %v", ok, err))
		}
		
		// Perform startup health check (search)
		_, err = adminClient.Search(
			props.BaseDN,
			element.ScopeWholeSubtree,
			element.DerefAlways,
			1,
			0, // timeout
			false,
			[]string{"dn"},
			strings.ReplaceAll(props.User.SearchFilter, "{0}", "health_check"),
		)
		if err != nil {
			fmt.Printf("WARN: startup admin search health check failed (this might be expected depending on LDAP config): %v\n", err)
		}
	}

	return &LdapSessionIdentityAccessManager{
		adminClient:   adminClient,
		userClient:    userClient,
		baseDN:        props.BaseDN,
		userFilter:    props.User.SearchFilter,
		policyManager: policyManager,
	}
}

func (m *LdapSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	if loginInfo.Password == nil || *loginInfo.Password == "" {
		return bo.LoginAuthenticationFailed, nil
	}

	m.mu.RLock()
	adminClient := m.adminClient
	userClient := m.userClient
	baseDN := m.baseDN
	userFilter := m.userFilter
	m.mu.RUnlock()

	if adminClient == nil || userClient == nil {
		return bo.LoginAuthenticationFailed, fmt.Errorf("LDAP clients not initialized")
	}

	// 1. Search for the user DN using adminClient
	filter := strings.ReplaceAll(userFilter, "{0}", fmt.Sprintf("%d", loginInfo.UserID))
	searchResult, err := adminClient.Search(
		baseDN,
		element.ScopeWholeSubtree,
		element.DerefAlways,
		2,     // sizeLimit: 2 so we can fail if multiple returned
		0,     // timeLimit
		false, // typesOnly
		nil,   // NO_ATTRIBUTES
		filter,
	)
	if err != nil {
		return bo.LoginAuthenticationFailed, nil
	}

	if len(searchResult.Entries) == 0 {
		return bo.LoggingInUserNotActive, nil
	}
	if len(searchResult.Entries) > 1 {
		return nil, exception.NewTurmsError(
			int32(constant.ResponseStatusCode_SERVER_INTERNAL_ERROR),
			fmt.Sprintf("More than 1 entry found for the user (%d), which means that the filter \"%s\" is wrong", loginInfo.UserID, userFilter),
		)
	}

	userDN := searchResult.Entries[0].ObjectName

	// 2. Bind with the user DN and password
	// RFC 4511: 4.2.1. Processing of the Bind Request
	// Serialize bind requests using userBindMu to mimic Java's TaskScheduler schedule behavior
	m.userBindMu.Lock()
	ok, err := userClient.Bind(true, userDN, *loginInfo.Password)
	m.userBindMu.Unlock()

	if err != nil || !ok {
		return bo.LoginAuthenticationFailed, nil
	}

	return bo.GrantedWithAllPermissions, nil
}

func (m *LdapSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) bool {
	props, ok := properties.(*config.IdentityAccessManagementProperties)
	if !ok {
		return false
	}
	
	if !props.Enabled || props.Type != config.IdentityAccessManagementType_LDAP {
		return false
	}
	
	if !strings.Contains(props.Ldap.User.SearchFilter, "{0}") {
		fmt.Printf("WARN: The User Search Filter must contain {0} to substitute the user ID: %s\n", props.Ldap.User.SearchFilter)
		return false
	}

	adminClient, err := ldap.NewLdapClient(props.Ldap.Admin.Host, props.Ldap.Admin.Port, props.Ldap.Admin.UseTLS, nil, 5*time.Second)
	if err != nil {
		fmt.Printf("WARN: failed to connect to admin LDAP server: %v\n", err)
		return false
	}
	
	userClient, err := ldap.NewLdapClient(props.Ldap.User.Host, props.Ldap.User.Port, props.Ldap.User.UseTLS, nil, 5*time.Second)
	if err != nil {
		fmt.Printf("WARN: failed to connect to user LDAP server: %v\n", err)
		adminClient.Close()
		return false
	}

	if props.Ldap.Admin.Username != "" {
		ok, err := adminClient.Bind(false, props.Ldap.Admin.Username, props.Ldap.Admin.Password)
		if err != nil || !ok {
			fmt.Printf("WARN: admin bind failed (ok=%v): %v\n", ok, err)
			adminClient.Close()
			userClient.Close()
			return false
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.adminClient != nil {
		m.adminClient.Close()
	}
	if m.userClient != nil {
		m.userClient.Close()
	}

	m.adminClient = adminClient
	m.userClient = userClient
	m.baseDN = props.Ldap.BaseDN
	m.userFilter = props.Ldap.User.SearchFilter

	return props.Enabled
}

// NoopSessionIdentityAccessManager maps to NoopSessionIdentityAccessManager in Java.
type NoopSessionIdentityAccessManager struct{}

func (m *NoopSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	// Java: returns GRANTED_WITH_ALL_PERMISSIONS (all permissions)
	return bo.GrantedWithAllPermissions, nil
}
func (m *NoopSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) bool {
	return true
}
