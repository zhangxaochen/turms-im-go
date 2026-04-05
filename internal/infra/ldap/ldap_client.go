package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"im.turms/server/internal/infra/ldap/asn1"
	"im.turms/server/internal/infra/ldap/element"
)

// LdapClient is a native LDAP client supporting MessageID multiplexing over a single TCP connection.
// @MappedFrom im.turms.gateway.infra.ldap.LdapClient
type LdapClient struct {
	host   string
	port   int
	useTLS bool

	conn           net.Conn
	pendingRequests sync.Map // map[int32]*pendingRequest

	// RFC 4511 C.1.5: messageID of requests MUST be non-zero (zero reserved for Notice of Disconnection)
	nextMessageID int32

	readBuffer *asn1.BerBuffer

	// writeMu protects conn.Write from concurrent goroutines interleaving partial writes
	writeMu sync.Mutex

	onClosed func(error)
	closed   atomic.Bool
}

type pendingRequest struct {
	id       int32
	response chan interface{}
	decoder  func(buffer *asn1.BerBuffer) (interface{}, bool, error) // bool: isComplete
}

func NewLdapClient(host string, port int, useTLS bool, tlsConfig *tls.Config, timeout time.Duration) (*LdapClient, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	dialer := &net.Dialer{Timeout: timeout}
	var conn net.Conn
	var err error

	if useTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}

	if err != nil {
		return nil, err
	}

	client := &LdapClient{
		host:          host,
		port:          port,
		useTLS:        useTLS,
		conn:          conn,
		nextMessageID: 1,
		readBuffer:    asn1.NewBerBuffer(4096),
	}

	go client.readLoop()

	return client, nil
}

func (c *LdapClient) Close() error {
	if c.closed.Swap(true) {
		return nil
	}
	return c.conn.Close()
}

func (c *LdapClient) IsConnected() bool {
	return c.conn != nil && !c.closed.Load()
}

func (c *LdapClient) nextID() int32 {
	return atomic.AddInt32(&c.nextMessageID, 1)
}

// sendRequest sends an LDAP request with optional controls and returns a channel for the response.
func (c *LdapClient) sendRequest(op element.ProtocolOperation, controls []element.Control, decoder func(*asn1.BerBuffer) (interface{}, bool, error)) (chan interface{}, int32, error) {
	if c.closed.Load() {
		return nil, 0, errors.New("LDAP client is closed")
	}

	id := c.nextID()
	respChan := make(chan interface{}, 1)

	c.pendingRequests.Store(id, &pendingRequest{
		id:       id,
		response: respChan,
		decoder:  decoder,
	})

	msg := element.LdapMessage{
		MessageId:         int(id),
		ProtocolOperation: op,
		Controls:          controls,
	}

	buf := asn1.NewBerBuffer(op.EstimateSize() + 16)
	msg.WriteTo(buf)

	c.writeMu.Lock()
	_, err := c.conn.Write(buf.Bytes())
	c.writeMu.Unlock()

	if err != nil {
		c.pendingRequests.Delete(id)
		close(respChan)
		return nil, id, fmt.Errorf("failed to write LDAP request: %w", err)
	}

	return respChan, id, nil
}

func (c *LdapClient) readLoop() {
	tmpBuf := make([]byte, 4096)
	for {
		n, err := c.conn.Read(tmpBuf)
		if err != nil {
			c.handleClose(err)
			return
		}

		c.readBuffer.Append(tmpBuf[:n])
		c.processReadBuffer()
	}
}

func (c *LdapClient) processReadBuffer() {
	for c.readBuffer.IsReadable() {
		c.readBuffer.MarkReaderIndex()

		tag := c.readBuffer.ReadTag()
		if tag != 0x30 { // LdapMessage must be a SEQUENCE
			c.handleClose(errors.New("LDAP protocol error: expected SEQUENCE tag"))
			return
		}

		length := c.readBuffer.TryReadLengthIfReadable()
		if length == -1 {
			c.readBuffer.ResetReaderIndex()
			return
		}

		if !c.readBuffer.IsReadableWithCount(length) {
			c.readBuffer.ResetReaderIndex()
			return
		}

		// Framed message complete
		msgEnd := c.readBuffer.ReaderIndex() + length

		messageID := c.readBuffer.ReadInteger()

		val, ok := c.pendingRequests.Load(int32(messageID))
		if !ok {
			// Unknown message ID? Skip it.
			c.readBuffer.SkipBytes(msgEnd - c.readBuffer.ReaderIndex())
			continue
		}

		req := val.(*pendingRequest)

		// Decode protocol op
		res, complete, err := req.decoder(c.readBuffer)
		if err != nil {
			req.response <- err
			c.pendingRequests.Delete(req.id)
			c.readBuffer.SkipBytes(msgEnd - c.readBuffer.ReaderIndex())
			continue
		}

		// Skip optional controls
		if c.readBuffer.ReaderIndex() < msgEnd {
			_ = element.DecodeControls(c.readBuffer)
		}

		// Skip any trailing bytes in this message frame
		if c.readBuffer.ReaderIndex() < msgEnd {
			c.readBuffer.SkipBytes(msgEnd - c.readBuffer.ReaderIndex())
		}

		if complete {
			req.response <- res
			c.pendingRequests.Delete(req.id)
		}
	}
}

func (c *LdapClient) handleClose(err error) {
	c.closed.Store(true)
	c.pendingRequests.Range(func(key, value interface{}) bool {
		req := value.(*pendingRequest)
		req.response <- fmt.Errorf("connection closed: %w", err)
		return true
	})
	if c.onClosed != nil {
		c.onClosed(err)
	}
}

// Bind sends an LDAP Bind request. If useFastBind is true, it sends with the LDAP_SERVER_FAST_BIND_OID control.
// @MappedFrom LdapClient.bind(boolean useFastBind, String dn, String password)
func (c *LdapClient) Bind(useFastBind bool, dn, password string) (bool, error) {
	req := &element.BindRequest{
		Version:  element.LdapVersion3,
		Name:     dn,
		Password: password,
	}

	var controls []element.Control
	if useFastBind {
		controls = []element.Control{element.FastBindControl}
	}

	ch, _, err := c.sendRequest(req, controls, func(buf *asn1.BerBuffer) (interface{}, bool, error) {
		return element.DecodeBindResponse(buf), true, nil
	})

	if err != nil {
		return false, err
	}

	res := <-ch
	if err, ok := res.(error); ok {
		return false, err
	}
	resp := res.(*element.BindResponse)
	if resp.IsSuccess() {
		return true, nil
	}
	if resp.ResultCode == element.ResultCodeInvalidCredentials {
		return false, nil
	}
	return false, &LdapException{
		ResultCode:        resp.ResultCode,
		DiagnosticMessage: resp.DiagnosticMessage,
	}
}

// Search sends an LDAP Search request.
// @MappedFrom LdapClient.search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter)
func (c *LdapClient) Search(
	baseDN string,
	scope element.LdapScope,
	derefAliases element.DerefAliases,
	sizeLimit int,
	timeLimit int,
	typesOnly bool,
	attributes []string,
	filter string,
) (*element.SearchResult, error) {
	req := &element.SearchRequest{
		BaseDN:       baseDN,
		Scope:        scope,
		DerefAlias:   derefAliases,
		SizeLimit:    sizeLimit,
		TimeLimit:    timeLimit,
		TypesOnly:    typesOnly,
		Filter:       filter,
		Attributes:   attributes,
	}

	var accumulatedEntries []element.SearchResultEntry

	ch, _, err := c.sendRequest(req, nil, func(buf *asn1.BerBuffer) (interface{}, bool, error) {
		tag := buf.PeekTag()
		switch int(tag) {
		case element.LdapTagSearchResultEntry:
			entry := element.DecodeSearchResultEntry(buf)
			accumulatedEntries = append(accumulatedEntries, *entry)
			return nil, false, nil
		case element.LdapTagSearchResultDone:
			buf.SkipTag()
			buf.SkipLength()
			res := element.DecodeLdapResult(buf)
			return &element.SearchResult{
				LdapResult: res,
				Entries:    accumulatedEntries,
			}, true, nil
		default:
			return nil, false, fmt.Errorf("unexpected tag in search response: %d", tag)
		}
	})

	if err != nil {
		return nil, err
	}

	res := <-ch
	if err, ok := res.(error); ok {
		return nil, err
	}
	return res.(*element.SearchResult), nil
}

// LdapException represents an error from LDAP server.
// @MappedFrom im.turms.gateway.infra.ldap.LdapException
type LdapException struct {
	ResultCode        int
	DiagnosticMessage string
}

func (e *LdapException) Error() string {
	return fmt.Sprintf("LDAP error (code=%d): %s", e.ResultCode, e.DiagnosticMessage)
}
