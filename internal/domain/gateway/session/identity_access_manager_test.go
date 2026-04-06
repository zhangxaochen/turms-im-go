package session

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"im.turms/server/internal/domain/gateway/config"
	"im.turms/server/internal/domain/gateway/session/bo"
)

func startMockLdapServer(b *testing.B) string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatalf("failed to listen: %v", err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go handleMockLdapConn(conn)
		}
	}()
	return l.Addr().String()
}

func handleMockLdapConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				//
			}
			return
		}

		if n <= 0 {
			continue
		}

		// Quick hack to parse msgID from universal ASN.1 request:
		// 30 <len> 02 01 <msgId> <op> ...
		// We just assume the msgId is at index 4 (if length is 1 byte) or index 5 (if length is 2 bytes).
		// We'll search for 0x02, 0x01, and the next byte is msgId.
		msgID := byte(1)
		for i := 0; i < n-2; i++ {
			if buf[i] == 0x02 && buf[i+1] == 0x01 {
				msgID = buf[i+2]
				break
			}
		}

		// Very naive check for operation type
		// If it contains 0x60/0x61 etc.
		// Wait, 0x60 is bind request. 0x63 is search request.
		isBind := false
		isSearch := false
		for i := 0; i < n; i++ {
			if buf[i] == 0x60 {
				isBind = true
				break
			}
			if buf[i] == 0x63 {
				isSearch = true
				break
			}
		}

		if isBind {
			_, _ = conn.Write(serializeBindResponse(int(msgID)))
		} else if isSearch {
			_, _ = conn.Write(serializeSearchEntry(int(msgID)))
			_, _ = conn.Write(serializeSearchDone(int(msgID)))
		}
	}
}

func serializeBindResponse(msgID int) []byte {
	msg := []byte{
		0x02, 0x01, byte(msgID), // MessageID
		0x61, 0x07, // BindResponse
		0x0a, 0x01, 0x00, // ResultCode (success = 0)
		0x04, 0x00, // MatchedDN
		0x04, 0x00, // ErrorMessage
	}
	buf := []byte{0x30, byte(len(msg))}
	buf = append(buf, msg...)
	return buf
}

func serializeSearchEntry(msgID int) []byte {
	dn := "uid=1,ou=users,dc=example,dc=com"
	entryMsg := []byte{0x04, byte(len(dn))}
	entryMsg = append(entryMsg, []byte(dn)...)
	entryMsg = append(entryMsg, 0x30, 0x00) // Attributes

	entryWrap := []byte{0x64, byte(len(entryMsg))}
	entryWrap = append(entryWrap, entryMsg...)

	msg := []byte{0x02, 0x01, byte(msgID)}
	msg = append(msg, entryWrap...)

	buf := []byte{0x30, byte(len(msg))}
	buf = append(buf, msg...)
	return buf
}

func serializeSearchDone(msgID int) []byte {
	msg := []byte{
		0x02, 0x01, byte(msgID), // MessageID
		0x65, 0x07, // SearchResultDone
		0x0a, 0x01, 0x00, // ResultCode (success = 0)
		0x04, 0x00, // MatchedDN
		0x04, 0x00, // ErrorMessage
	}
	buf := []byte{0x30, byte(len(msg))}
	buf = append(buf, msg...)
	return buf
}

func BenchmarkLdapSessionIdentityAccessManager_VerifyAndGrant(b *testing.B) {
	addr := startMockLdapServer(b)

	host, portStr, _ := net.SplitHostPort(addr)
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	props := &config.IdentityAccessManagementProperties{
		Type:    config.IdentityAccessManagementType_LDAP,
		Enabled: true,
	}
	props.Ldap.Admin.Host = host
	props.Ldap.Admin.Port = port
	props.Ldap.User.Host = host
	props.Ldap.User.Port = port
	props.Ldap.BaseDN = "dc=example,dc=com"
	props.Ldap.User.SearchFilter = "(uid={0})"

	mgr := &LdapSessionIdentityAccessManager{}
	mgr.UpdateGlobalProperties(props)

	pass := "secret"
	userID := int64(1)
	loginInfo := &bo.UserLoginInfo{
		UserID:   &userID,
		Password: &pass,
	}
	ctx := context.Background()

	b.ResetTimer()
	var successCount int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := mgr.VerifyAndGrant(ctx, loginInfo)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			} else {
				b.Errorf("VerifyAndGrant failed: %v", err)
			}
		}
	})

	b.ReportMetric(float64(successCount), "auths")
}
