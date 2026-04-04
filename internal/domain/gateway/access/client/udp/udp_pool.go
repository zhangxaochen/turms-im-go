package udp

import (
	"encoding/binary"
	"sync"

	"im.turms/server/internal/domain/common/constant"
)

var (
	codePool         = make(map[constant.ResponseStatusCode][]byte)
	notificationPool = make(map[UdpNotificationType][]byte)
	poolMu           sync.RWMutex
)

func init() {
	// Initialize notification pool
	// Ordinal + 1 per Java implementation
	notificationPool[OpenConnectionNotification] = []byte{byte(OpenConnectionNotification) + 1}
}

func GetBufferFromStatusCode(code constant.ResponseStatusCode) []byte {
	poolMu.RLock()
	buf, ok := codePool[code]
	poolMu.RUnlock()
	if ok {
		return buf
	}

	poolMu.Lock()
	defer poolMu.Unlock()
	// Double check
	if buf, ok = codePool[code]; ok {
		return buf
	}

	if code == constant.ResponseStatusCode_OK {
		buf = []byte{}
	} else {
		buf = make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(code))
	}
	codePool[code] = buf
	return buf
}

func GetBufferFromNotificationType(t UdpNotificationType) []byte {
	poolMu.RLock()
	defer poolMu.RUnlock()
	return notificationPool[t]
}
