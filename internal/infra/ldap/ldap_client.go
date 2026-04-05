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

type LdapClient struct {
	conn            net.Conn
	pendingRequests sync.Map // map[int32]*pendingRequest
	nextMessageID   int32
	
	readBuffer      *asn1.BerBuffer
	
	onClosed        func(error)
	closed          atomic.Bool
}

type pendingRequest struct {
	id       int32
	response chan interface{}
	decoder  func(buffer *asn1.BerBuffer) (interface{}, bool, error) // bool: isComplete
	result   interface{}
}

func NewLdapClient(addr string, useTLS bool, tlsConfig *tls.Config, timeout time.Duration) (*LdapClient, error) {
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

func (c *LdapClient) Addr() string {
	if c.conn == nil {
		return ""
	}
	return c.conn.RemoteAddr().String()
}

func (c *LdapClient) nextID() int32 {
	return atomic.AddInt32(&c.nextMessageID, 1)
}

func (c *LdapClient) SendRequest(op element.ProtocolOperation, decoder func(*asn1.BerBuffer) (interface{}, bool, error)) (chan interface{}, int32) {
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
	}
	
	buf := asn1.NewBerBuffer(op.EstimateSize() + 16)
	msg.WriteTo(buf)
	
	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		c.pendingRequests.Delete(id)
		close(respChan)
		return nil, id
	}
	
	return respChan, id
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
			// Protocol error
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
		
		// Skip optional controls (we don't handle them for now or skip them)
		if c.readBuffer.ReaderIndex() < msgEnd {
			_ = element.DecodeControls(c.readBuffer)
		}
		
		// Skip any trailing bytes in this message frame
		if c.readBuffer.ReaderIndex() < msgEnd {
			c.readBuffer.SkipBytes(msgEnd - c.readBuffer.ReaderIndex())
		}
		
		if complete {
			req.response <- res
			c.pendingRequests.Store(req.id, req)
			c.pendingRequests.Delete(req.id)
		} else {
			// For SearchResult, it might take multiple entries.
			// Current logic stores the intermediate state in req.result (implicit in the decoder's closure if needed)
			// But for simplicity, the decoder should update its state or return the full state.
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

// Helper methods for common operations

func (c *LdapClient) Bind(name, password string) (*element.BindResponse, error) {
	req := &element.BindRequest{
		Version:  element.LdapVersion3,
		Name:     name,
		Password: password,
	}
	
	ch, _ := c.SendRequest(req, func(buf *asn1.BerBuffer) (interface{}, bool, error) {
		return element.DecodeBindResponse(buf), true, nil
	})
	
	if ch == nil {
		return nil, errors.New("failed to send bind request")
	}
	
	res := <-ch
	if err, ok := res.(error); ok {
		return nil, err
	}
	return res.(*element.BindResponse), nil
}

func (c *LdapClient) Search(baseDN string, scope element.LdapScope, filter string, attrs []string) (*element.SearchResult, error) {
	req := &element.SearchRequest{
		BaseDN:     baseDN,
		Scope:      scope,
		DerefAlias: element.DerefNever,
		Filter:     filter,
		Attributes: attrs,
	}
	
	var accumulatedEntries []element.SearchResultEntry
	
	ch, _ := c.SendRequest(req, func(buf *asn1.BerBuffer) (interface{}, bool, error) {
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
	
	if ch == nil {
		return nil, errors.New("failed to send search request")
	}
	
	res := <-ch
	if err, ok := res.(error); ok {
		return nil, err
	}
	return res.(*element.SearchResult), nil
}
