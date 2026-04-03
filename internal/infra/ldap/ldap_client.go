package ldap

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

// LdapClient maps to LdapClient in Java but uses github.com/go-ldap/ldap/v3 under the hood.
// @MappedFrom LdapClient
type LdapClient struct {
	Addr       string
	Conn       *ldap.Conn
	UseTLS     bool
	SkipVerify bool
}

func NewLdapClient(addr string, useTLS bool, skipVerify bool) *LdapClient {
	return &LdapClient{
		Addr:       addr,
		UseTLS:     useTLS,
		SkipVerify: skipVerify,
	}
}

// @MappedFrom connect()
func (c *LdapClient) Connect() error {
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
func (c *LdapClient) Bind(useFastBind bool, dn string, password string) error {
	if c.Conn == nil {
		return fmt.Errorf("URL string is missing or connection is not established")
	}

	// We use Simple Bind regardless of useFastBind since go-ldap abstracts this well.
	return c.Conn.Bind(dn, password)
}

// @MappedFrom modify(String dn, List<ModifyOperationChange> changes)
// For simplicity, we accept go-ldap's ModifyRequest natively or build it here.
func (c *LdapClient) Modify(req *ldap.ModifyRequest) error {
	if c.Conn == nil {
		return fmt.Errorf("connection not established")
	}
	return c.Conn.Modify(req)
}

func (c *LdapClient) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
