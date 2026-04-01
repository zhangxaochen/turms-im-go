package discovery

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"im.turms/server/internal/testingutil"
)
type mockListener struct {
	called int
}

func (m *mockListener) OnMembersChange() {
	m.called++
}



func TestDiscoveryService(t *testing.T) {
	ctx := context.Background()
	rc, cleanup := testingutil.SetupRedis(t)
	defer cleanup()

	// 2. Clear out any previous data for our test cluster
	keys, _ := rc.RDB.Keys(ctx, "turms:cluster:test_cluster:*").Result()
	if len(keys) > 0 {
		rc.RDB.Del(ctx, keys...)
	}

	localMember := &Member{
		ClusterID:        "test_cluster",
		NodeID:           "node_mock_1",
		NodeType:         NodeTypeService,
		IsLeaderEligible: true,
		IsActive:         true,
	}

	ds1 := NewDiscoveryService(rc, localMember, 1)

	lis := &mockListener{}
	ds1.AddListener(lis)

	// 3. Start Discovery Service
	err := ds1.Start()
	assert.NoError(t, err)

	// Verify local member was registered and is returned
	members := ds1.GetMembers()
	assert.Len(t, members, 1)
	assert.Equal(t, "node_mock_1", members[0].NodeID)

	// 4. Test multiple nodes
	localMember2 := &Member{
		ClusterID:        "test_cluster",
		NodeID:           "node_mock_2",
		NodeType:         NodeTypeGateway,
		IsLeaderEligible: true, // But starts later, shouldn't usurp leader initially
		IsActive:         true,
	}
	ds2 := NewDiscoveryService(rc, localMember2, 1)
	err = ds2.Start()
	assert.NoError(t, err)

	// Wait for background sync loop to catch the new node
	time.Sleep(1200 * time.Millisecond)

	// ds1 should now see ds2
	members1 := ds1.GetMembers()
	assert.Len(t, members1, 2)
	assert.Condition(t, func() bool {
		return members1[0].NodeID == "node_mock_2" || members1[1].NodeID == "node_mock_2"
	})

	// 5. Test Leader Election
	time.Sleep(1200 * time.Millisecond)
	// Both should agree on leader
	leader1 := ds1.GetLeaderID()
	leader2 := ds2.GetLeaderID()
	assert.NotEmpty(t, leader1)
	assert.Equal(t, leader1, leader2)
	assert.Equal(t, "node_mock_1", leader1, "First node to start should be leader")

	// 6. Test failover
	ds1.Stop()                  // Node 1 dies
	time.Sleep(4 * time.Second) // Wait for TTL to expire (TTL is 3s for 1s interval)

	// ds2 should realize it's alone and become leader
	members2_after := ds2.GetMembers()
	assert.Len(t, members2_after, 1, "Node 1 should be gone")
	leader2_after := ds2.GetLeaderID()
	assert.Equal(t, "node_mock_2", leader2_after, "Node 2 should take over leadership")

	ds2.Stop()
}
