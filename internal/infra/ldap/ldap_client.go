package ldap

import (
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/go-ldap/ldap/v3"
)

// ModifyOperationChange maps to ModifyOperationChange in Java
type ModifyOperationChange struct {
	Type      uint // ldap.ModifyRequestOpAdd, etc.
	Attribute string
	Values    []string
}

// LdapClient maps to LdapClient in Java but uses github.com/go-ldap/ldap/v3 under the hood.
// @MappedFrom LdapClient
type LdapClient struct {
	Addr       string
	Conn       *ldap.Conn
	UseTLS     bool
	SkipVerify bool
	mu         sync.RWMutex
}

func NewLdapClient(addr string, useTLS bool, skipVerify bool) *LdapClient {
	return &LdapClient{
		Addr:       addr,
		UseTLS:     useTLS,
		SkipVerify: skipVerify,
	}
}

// @MappedFrom isConnected()
func (c *LdapClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Conn != nil && !c.Conn.IsClosing()
}

// @MappedFrom connect()
func (c *LdapClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Connection-sharing semantics (do not reconnect if already established)
	if c.Conn != nil && !c.Conn.IsClosing() {
		return nil
	}

	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}

	var l *ldap.Conn
	var err error

	if c.UseTLS {
		l, err = ldap.DialTLS("tcp", c.Addr, &tls.Config{InsecureSkipVerify: c.SkipVerify})
	} else {
		l, err = ldap.DialURL(fmt.Sprintf("ldap://%s", c.Addr))
	}

	if err != nil {
		return err
	}
	c.Conn = l
	return nil
}

// @MappedFrom bind(boolean useFastBind, String dn, String password)
func (c *LdapClient) Bind(useFastBind bool, dn string, password string) (bool, error) {
	c.mu.RLock()
	conn := c.Conn
	c.mu.RUnlock()

	if conn == nil {
		return false, fmt.Errorf("connection is not established")
	}

	// We use Simple Bind regardless of useFastBind since go-ldap abstracts this well,
	// but we could send specific controls if fastBind optimization was critical.
	err := conn.Bind(dn, password)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// @MappedFrom search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter)
func (c *LdapClient) Search(baseDn string, scope int, derefAliases int, sizeLimit int, timeLimit int, typeOnly bool, attributes []string, filter string) (*ldap.SearchResult, error) {
	c.mu.RLock()
	conn := c.Conn
	c.mu.RUnlock()

	if conn == nil {
		return nil, fmt.Errorf("connection not established")
	}

	req := ldap.NewSearchRequest(
		baseDn,
		scope,
		derefAliases,
		sizeLimit,
		timeLimit,
		typeOnly,
		filter,
		attributes,
		nil,
	)
	return conn.Search(req)
}

// @MappedFrom modify(String dn, List<ModifyOperationChange> changes)
func (c *LdapClient) Modify(dn string, changes []ModifyOperationChange) error {
	if len(changes) == 0 {
		return nil
	}

	req := ldap.NewModifyRequest(dn, nil)
	for _, change := range changes {
		if change.Type == ldap.AddAttribute && len(change.Values) == 0 {
			return fmt.Errorf("INVALID_ATTRIBUTE_SYNTAX: ADD operation requires at least one value for attribute %s", change.Attribute)
		}
		switch change.Type {
		case ldap.AddAttribute:
			req.Add(change.Attribute, change.Values)
		case ldap.DeleteAttribute:
			req.Delete(change.Attribute, change.Values)
		case ldap.ReplaceAttribute:
			req.Replace(change.Attribute, change.Values)
		}
	}

	c.mu.RLock()
	conn := c.Conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("connection not established")
	}
	return conn.Modify(req)
}

func (c *LdapClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}
}
