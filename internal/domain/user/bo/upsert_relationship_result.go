package bo

type UpsertRelationshipResult struct {
	IsCreated  bool
	GroupIndex *int32
}

var (
	UpsertRelationshipResultCreated    = UpsertRelationshipResult{IsCreated: true}
	UpsertRelationshipResultNotCreated = UpsertRelationshipResult{IsCreated: false}
)
