package authorization

type PolicyStatementEffect int

const (
	PolicyStatementEffect_ALLOW PolicyStatementEffect = iota
	PolicyStatementEffect_DENY
)

type PolicyStatementAction int

const (
	PolicyStatementAction_ALL PolicyStatementAction = iota
	PolicyStatementAction_CREATE
	PolicyStatementAction_DELETE
	PolicyStatementAction_UPDATE
	PolicyStatementAction_QUERY
)

type PolicyStatementResource int

const (
	PolicyStatementResource_USER PolicyStatementResource = iota
	PolicyStatementResource_USER_LOCATION
	PolicyStatementResource_USER_ONLINE_STATUS
	PolicyStatementResource_USER_PROFILE
	PolicyStatementResource_USER_SETTING
	PolicyStatementResource_NEARBY_USER
	PolicyStatementResource_RELATIONSHIP
	PolicyStatementResource_RELATIONSHIP_GROUP
	PolicyStatementResource_FRIEND_REQUEST
	PolicyStatementResource_GROUP
	PolicyStatementResource_GROUP_BLOCKED_USER
	PolicyStatementResource_GROUP_INVITATION
	PolicyStatementResource_GROUP_JOIN_QUESTION
	PolicyStatementResource_GROUP_JOIN_QUESTION_ANSWER
	PolicyStatementResource_GROUP_JOIN_REQUEST
	PolicyStatementResource_GROUP_MEMBER
	PolicyStatementResource_JOINED_GROUP
	PolicyStatementResource_MESSAGE
	PolicyStatementResource_CONVERSATION
	PolicyStatementResource_CONVERSATION_SETTING
	PolicyStatementResource_TYPING_STATUS
	PolicyStatementResource_MEETING
	PolicyStatementResource_RESOURCE
)

// @MappedFrom PolicyStatement
type PolicyStatement struct {
	Effect    PolicyStatementEffect
	Actions   []PolicyStatementAction
	Resources []PolicyStatementResource
}

// @MappedFrom PolicyStatement(PolicyStatementEffect effect, Set<PolicyStatementAction> actions, Set<PolicyStatementResource> resources)
func NewPolicyStatement(effect PolicyStatementEffect, actions []PolicyStatementAction, resources []PolicyStatementResource) *PolicyStatement {
	return &PolicyStatement{Effect: effect, Actions: actions, Resources: resources}
}

// @MappedFrom Policy
type Policy struct {
	Statements []PolicyStatement
}

// @MappedFrom Policy(List<PolicyStatement> statements)
func NewPolicy(statements []PolicyStatement) *Policy {
	return &Policy{Statements: statements}
}

// @MappedFrom PolicyDeserializer
type PolicyDeserializer struct {
}

// @MappedFrom parse(Map<String, Object> map)
func (d *PolicyDeserializer) Parse(data map[string]interface{}) (*Policy, error) {
	// Dummy parser pending full implementation
	return &Policy{}, nil
}

var ResourceOperations = map[PolicyStatementResource]struct {
	Creating []int32
	Deleting []int32
	Updating []int32
	Querying []int32
}{
	// Simple map for future population
}
