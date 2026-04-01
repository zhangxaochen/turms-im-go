package idgen

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/**
 * The flake ID is designed for turms. ID size: 64 bits.
 *
 * 1 bit for the sign of ID. Always 0.
 * 41 bits for timestamp (69 years).
 * 4 bits for data center ID (16).
 * 8 bits for worker ID (256).
 * 10 bits for sequenceNumber (1,024/ms).
 */

const (
	Epoch int64 = 1602547200000 // 2020-10-13 00:00:00 in UTC

	TimestampBits      = 41
	DataCenterIDBits   = 4
	WorkerIDBits       = 8
	SequenceNumberBits = 10

	TimestampLeftShift = SequenceNumberBits + WorkerIDBits + DataCenterIDBits
	DataCenterIDShift  = SequenceNumberBits + WorkerIDBits
	WorkerIDShift      = SequenceNumberBits

	SequenceNumberMask = (1 << SequenceNumberBits) - 1

	MaxDataCenterID = 1 << DataCenterIDBits
	MaxWorkerID     = 1 << WorkerIDBits
)

type SnowflakeIdGenerator struct {
	mu             sync.Mutex
	lastTimestamp  int64
	sequenceNumber int64

	dataCenterId int64
	workerId     int64
}

// NewSnowflakeIdGenerator initializes a new Thread-Safe Snowflake Generator.
func NewSnowflakeIdGenerator(dataCenterId int64, workerId int64) (*SnowflakeIdGenerator, error) {
	gen := &SnowflakeIdGenerator{}
	if err := gen.UpdateNodeInfo(dataCenterId, workerId); err != nil {
		return nil, err
	}

	// Randomize the sequenceNumber on init to decrease the chance of collision.
	gen.sequenceNumber = int64(rand.Int31())

	return gen, nil
}

// UpdateNodeInfo updates node ID safely or returns an error if over limits.
func (g *SnowflakeIdGenerator) UpdateNodeInfo(dataCenterId int64, workerId int64) error {
	if dataCenterId < 0 || dataCenterId >= MaxDataCenterID {
		return fmt.Errorf("data center ID must be in the range: [0, %d), but got: %d", MaxDataCenterID, dataCenterId)
	}
	if workerId < 0 || workerId >= MaxWorkerID {
		return fmt.Errorf("worker ID must be in the range: [0, %d), but got: %d", MaxWorkerID, workerId)
	}

	g.dataCenterId = dataCenterId
	g.workerId = workerId
	return nil
}

// NextIncreasingId generates monotonically increasing 64-bit ID.
func (g *SnowflakeIdGenerator) NextIncreasingId() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.sequenceNumber++
	sequenceNum := g.sequenceNumber & SequenceNumberMask

	now := time.Now().UnixMilli()
	if now > g.lastTimestamp {
		g.lastTimestamp = now
	}

	if sequenceNum == 0 {
		g.lastTimestamp++
	}

	timestamp := g.lastTimestamp - Epoch

	return (timestamp << TimestampLeftShift) | (g.dataCenterId << DataCenterIDShift) | (g.workerId << WorkerIDShift) | sequenceNum
}

// NextLargeGapId generates IDs with large gaps for better sharding distribution in MongoDB.
func (g *SnowflakeIdGenerator) NextLargeGapId() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.sequenceNumber++
	sequenceNum := g.sequenceNumber & SequenceNumberMask

	now := time.Now().UnixMilli()
	if now > g.lastTimestamp {
		g.lastTimestamp = now
	}

	if sequenceNum == 0 {
		g.lastTimestamp++
	}

	timestamp := g.lastTimestamp - Epoch

	return (sequenceNum << (TimestampBits + DataCenterIDBits + WorkerIDBits)) |
		(timestamp << (DataCenterIDBits + WorkerIDBits)) |
		(g.dataCenterId << WorkerIDBits) | g.workerId
}
