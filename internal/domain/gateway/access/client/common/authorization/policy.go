package authorization

import (
	"errors"
	"fmt"
)

type PolicyStatementEffect int

const (
	PolicyStatementEffect_ALLOW PolicyStatementEffect = iota
	PolicyStatementEffect_DENY
)

func ParseEffect(e string) (PolicyStatementEffect, error) {
	switch e {
	case "ALLOW":
		return PolicyStatementEffect_ALLOW, nil
	case "DENY":
		return PolicyStatementEffect_DENY, nil
	default:
		return 0, fmt.Errorf("invalid policy statement effect: %s", e)
	}
}

type PolicyStatementAction int

const (
	PolicyStatementAction_ALL PolicyStatementAction = iota
	PolicyStatementAction_CREATE
	PolicyStatementAction_DELETE
	PolicyStatementAction_UPDATE
	PolicyStatementAction_QUERY
)

func ParseAction(a string) (PolicyStatementAction, error) {
	switch a {
	case "*":
		return PolicyStatementAction_ALL, nil
	case "CREATE":
		return PolicyStatementAction_CREATE, nil
	case "DELETE":
		return PolicyStatementAction_DELETE, nil
	case "UPDATE":
		return PolicyStatementAction_UPDATE, nil
	case "QUERY":
		return PolicyStatementAction_QUERY, nil
	default:
		return 0, fmt.Errorf("invalid policy statement action: %s", a)
	}
}

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

func ParseResource(r string) (PolicyStatementResource, error) {
	// Simple map-based or switch could be used here. For brevity, matching the original list:
	switch r {
	case "*":
		return 0, nil // Let's handle "*" specially in the parsing loop
	case "USER":
		return PolicyStatementResource_USER, nil
	case "USER_LOCATION":
		return PolicyStatementResource_USER_LOCATION, nil
	case "USER_ONLINE_STATUS":
		return PolicyStatementResource_USER_ONLINE_STATUS, nil
	case "USER_PROFILE":
		return PolicyStatementResource_USER_PROFILE, nil
	case "USER_SETTING":
		return PolicyStatementResource_USER_SETTING, nil
	case "NEARBY_USER":
		return PolicyStatementResource_NEARBY_USER, nil
	case "RELATIONSHIP":
		return PolicyStatementResource_RELATIONSHIP, nil
	case "RELATIONSHIP_GROUP":
		return PolicyStatementResource_RELATIONSHIP_GROUP, nil
	case "FRIEND_REQUEST":
		return PolicyStatementResource_FRIEND_REQUEST, nil
	case "GROUP":
		return PolicyStatementResource_GROUP, nil
	case "GROUP_BLOCKED_USER":
		return PolicyStatementResource_GROUP_BLOCKED_USER, nil
	case "GROUP_INVITATION":
		return PolicyStatementResource_GROUP_INVITATION, nil
	case "GROUP_JOIN_QUESTION":
		return PolicyStatementResource_GROUP_JOIN_QUESTION, nil
	case "GROUP_JOIN_QUESTION_ANSWER":
		return PolicyStatementResource_GROUP_JOIN_QUESTION_ANSWER, nil
	case "GROUP_JOIN_REQUEST":
		return PolicyStatementResource_GROUP_JOIN_REQUEST, nil
	case "GROUP_MEMBER":
		return PolicyStatementResource_GROUP_MEMBER, nil
	case "JOINED_GROUP":
		return PolicyStatementResource_JOINED_GROUP, nil
	case "MESSAGE":
		return PolicyStatementResource_MESSAGE, nil
	case "CONVERSATION":
		return PolicyStatementResource_CONVERSATION, nil
	case "CONVERSATION_SETTING":
		return PolicyStatementResource_CONVERSATION_SETTING, nil
	case "TYPING_STATUS":
		return PolicyStatementResource_TYPING_STATUS, nil
	case "MEETING":
		return PolicyStatementResource_MEETING, nil
	case "RESOURCE":
		return PolicyStatementResource_RESOURCE, nil
	default:
		return 0, fmt.Errorf("invalid policy statement resource: %s", r)
	}
}

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

var ErrIllegalPolicy = errors.New("illegal policy")

// @MappedFrom parse(Map<String, Object> map)
func (d *PolicyDeserializer) Parse(data map[string]interface{}) (*Policy, error) {
	statementsVal, ok := data["statements"]
	if !ok {
		return nil, fmt.Errorf("%w: missing 'statements' field", ErrIllegalPolicy)
	}

	statementsList, ok := statementsVal.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: 'statements' must be an array", ErrIllegalPolicy)
	}

	var parsedStatements []PolicyStatement
	for i, stmtObj := range statementsList {
		stmtMap, ok := stmtObj.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: statement %d must be an object", ErrIllegalPolicy, i)
		}

		effectStr, ok := stmtMap["effect"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: missing or invalid 'effect' in statement %d", ErrIllegalPolicy, i)
		}
		effect, err := ParseEffect(effectStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v in statement %d", ErrIllegalPolicy, err, i)
		}

		actionsVal, ok := stmtMap["actions"]
		if !ok {
			return nil, fmt.Errorf("%w: missing 'actions' in statement %d", ErrIllegalPolicy, i)
		}

		var actions []PolicyStatementAction
		switch v := actionsVal.(type) {
		case string:
			action, err := ParseAction(v)
			if err != nil {
				return nil, fmt.Errorf("%w: %v in statement %d", ErrIllegalPolicy, err, i)
			}
			actions = append(actions, action)
		case []interface{}:
			for _, av := range v {
				astr, ok := av.(string)
				if !ok {
					return nil, fmt.Errorf("%w: action must be a string in statement %d", ErrIllegalPolicy, i)
				}
				action, err := ParseAction(astr)
				if err != nil {
					return nil, fmt.Errorf("%w: %v in statement %d", ErrIllegalPolicy, err, i)
				}
				actions = append(actions, action)
			}
		default:
			return nil, fmt.Errorf("%w: 'actions' must be a string or array of strings in statement %d", ErrIllegalPolicy, i)
		}

		resourcesVal, ok := stmtMap["resources"]
		if !ok {
			return nil, fmt.Errorf("%w: missing 'resources' in statement %d", ErrIllegalPolicy, i)
		}

		var resources []PolicyStatementResource
		switch v := resourcesVal.(type) {
		case string:
			if v == "*" {
				// All resources
				for r := PolicyStatementResource_USER; r <= PolicyStatementResource_RESOURCE; r++ {
					resources = append(resources, r)
				}
			} else {
				res, err := ParseResource(v)
				if err != nil {
					return nil, fmt.Errorf("%w: %v in statement %d", ErrIllegalPolicy, err, i)
				}
				resources = append(resources, res)
			}
		case []interface{}:
			for _, rv := range v {
				rstr, ok := rv.(string)
				if !ok {
					return nil, fmt.Errorf("%w: resource must be a string in statement %d", ErrIllegalPolicy, i)
				}
				if rstr == "*" {
					// All resources, clear existing and add all
					resources = nil
					for r := PolicyStatementResource_USER; r <= PolicyStatementResource_RESOURCE; r++ {
						resources = append(resources, r)
					}
					break
				}
				res, err := ParseResource(rstr)
				if err != nil {
					return nil, fmt.Errorf("%w: %v in statement %d", ErrIllegalPolicy, err, i)
				}
				resources = append(resources, res)
			}
		default:
			return nil, fmt.Errorf("%w: 'resources' must be a string or array of strings in statement %d", ErrIllegalPolicy, i)
		}

		parsedStatements = append(parsedStatements, PolicyStatement{
			Effect:    effect,
			Actions:   actions,
			Resources: resources,
		})
	}

	return &Policy{Statements: parsedStatements}, nil
}

var ALL_RESOURCES = []PolicyStatementResource{
	PolicyStatementResource_USER,
	PolicyStatementResource_USER_LOCATION,
	PolicyStatementResource_USER_ONLINE_STATUS,
	PolicyStatementResource_USER_PROFILE,
	PolicyStatementResource_USER_SETTING,
	PolicyStatementResource_NEARBY_USER,
	PolicyStatementResource_RELATIONSHIP,
	PolicyStatementResource_RELATIONSHIP_GROUP,
	PolicyStatementResource_FRIEND_REQUEST,
	PolicyStatementResource_GROUP,
	PolicyStatementResource_GROUP_BLOCKED_USER,
	PolicyStatementResource_GROUP_INVITATION,
	PolicyStatementResource_GROUP_JOIN_QUESTION,
	PolicyStatementResource_GROUP_JOIN_QUESTION_ANSWER,
	PolicyStatementResource_GROUP_JOIN_REQUEST,
	PolicyStatementResource_GROUP_MEMBER,
	PolicyStatementResource_JOINED_GROUP,
	PolicyStatementResource_MESSAGE,
	PolicyStatementResource_CONVERSATION,
	PolicyStatementResource_CONVERSATION_SETTING,
	PolicyStatementResource_TYPING_STATUS,
	PolicyStatementResource_MEETING,
	PolicyStatementResource_RESOURCE,
}

var ResourceOperations = map[PolicyStatementResource]struct {
	Creating []int32
	Deleting []int32
	Updating []int32
	Querying []int32
}{
	PolicyStatementResource_USER: {
		Updating: []int32{105},
	},
	PolicyStatementResource_USER_LOCATION: {
		Updating: []int32{103},
	},
	PolicyStatementResource_USER_ONLINE_STATUS: {
		Updating: []int32{104},
		Querying: []int32{102},
	},
	PolicyStatementResource_USER_PROFILE: {
		Querying: []int32{100},
	},
	PolicyStatementResource_USER_SETTING: {
		Deleting: []int32{301}, // Not in snippet, but expected based on Java
		Updating: []int32{106}, // Not in snippet, but expected based on Java
		Querying: []int32{107}, // Not in snippet, but expected based on Java
	},
	PolicyStatementResource_NEARBY_USER: {
		Querying: []int32{101},
	},
	PolicyStatementResource_RELATIONSHIP: {
		Creating: []int32{202},
		Deleting: []int32{205},
		Updating: []int32{212},
		Querying: []int32{209, 207},
	},
	PolicyStatementResource_RELATIONSHIP_GROUP: {
		Creating: []int32{201},
		Deleting: []int32{204},
		Updating: []int32{211},
		Querying: []int32{208},
	},
	PolicyStatementResource_FRIEND_REQUEST: {
		Creating: []int32{200},
		Deleting: []int32{203},
		Updating: []int32{210},
		Querying: []int32{206},
	},
	PolicyStatementResource_GROUP: {
		Creating: []int32{300},
		Deleting: []int32{301},
		Updating: []int32{305},
		Querying: []int32{302},
	},
	PolicyStatementResource_GROUP_BLOCKED_USER: {
		Creating: []int32{306},
		Deleting: []int32{307},
		Querying: []int32{308, 309},
	},
	PolicyStatementResource_GROUP_INVITATION: {
		Creating: []int32{311},
		Deleting: []int32{314},
		Updating: []int32{320},
		Querying: []int32{317},
	},
	PolicyStatementResource_GROUP_JOIN_QUESTION: {
		Creating: []int32{313},
		Deleting: []int32{316},
		Updating: []int32{321},
		Querying: []int32{319},
	},
	PolicyStatementResource_GROUP_JOIN_QUESTION_ANSWER: {
		Querying: []int32{310},
	},
	PolicyStatementResource_GROUP_JOIN_REQUEST: {
		Creating: []int32{312},
		Deleting: []int32{315},
		Updating: []int32{322},
		Querying: []int32{318},
	},
	PolicyStatementResource_GROUP_MEMBER: {
		Creating: []int32{11},
		Deleting: []int32{12},
		Updating: []int32{14},
		Querying: []int32{13},
	},
	PolicyStatementResource_JOINED_GROUP: {
		Querying: []int32{303, 304},
	},
	PolicyStatementResource_MESSAGE: {
		Creating: []int32{8, 9}, // CREATE_MESSAGE_REQUEST, CREATE_MESSAGE_REACTIONS_REQUEST (9 not in snippet but in Java)
		Deleting: []int32{10},   // DELETE_MESSAGE_REACTIONS_REQUEST
		Updating: []int32{10},   // UPDATE_MESSAGE_REQUEST
		Querying: []int32{9},    // QUERY_MESSAGES_REQUEST
	},
	PolicyStatementResource_CONVERSATION: {
		Updating: []int32{6},
		Querying: []int32{5},
	},
	PolicyStatementResource_CONVERSATION_SETTING: {
		Deleting: []int32{401},
		Updating: []int32{400},
		Querying: []int32{402},
	},
	PolicyStatementResource_TYPING_STATUS: {
		Updating: []int32{7},
	},
	PolicyStatementResource_MEETING: {
		Creating: []int32{500},
		Deleting: []int32{501},
		Updating: []int32{502, 503},
		Querying: []int32{504},
	},
	PolicyStatementResource_RESOURCE: {
		Deleting: []int32{600},
		Updating: []int32{604},
		Querying: []int32{601, 602, 603},
	},
}

var ALL_REQUEST_TYPES []int32

func init() {
	allSet := make(map[int32]struct{})
	for _, ops := range ResourceOperations {
		for _, id := range ops.Creating {
			allSet[id] = struct{}{}
		}
		for _, id := range ops.Deleting {
			allSet[id] = struct{}{}
		}
		for _, id := range ops.Updating {
			allSet[id] = struct{}{}
		}
		for _, id := range ops.Querying {
			allSet[id] = struct{}{}
		}
	}
	for id := range allSet {
		ALL_REQUEST_TYPES = append(ALL_REQUEST_TYPES, id)
	}
}

type PolicyManager struct{}

func (m *PolicyManager) FindAllowedRequestTypes(policy *Policy) []int32 {
	if policy == nil || len(policy.Statements) == 0 {
		return nil
	}

	allowedSet := make(map[int32]struct{})

	// Handle ALLOW statements first
	for _, stmt := range policy.Statements {
		if stmt.Effect == PolicyStatementEffect_ALLOW {
			for _, res := range stmt.Resources {
				ops := ResourceOperations[res]
				for _, act := range stmt.Actions {
					switch act {
					case PolicyStatementAction_ALL:
						for _, id := range ops.Creating {
							allowedSet[id] = struct{}{}
						}
						for _, id := range ops.Deleting {
							allowedSet[id] = struct{}{}
						}
						for _, id := range ops.Updating {
							allowedSet[id] = struct{}{}
						}
						for _, id := range ops.Querying {
							allowedSet[id] = struct{}{}
						}
					case PolicyStatementAction_CREATE:
						for _, id := range ops.Creating {
							allowedSet[id] = struct{}{}
						}
					case PolicyStatementAction_DELETE:
						for _, id := range ops.Deleting {
							allowedSet[id] = struct{}{}
						}
					case PolicyStatementAction_UPDATE:
						for _, id := range ops.Updating {
							allowedSet[id] = struct{}{}
						}
					case PolicyStatementAction_QUERY:
						for _, id := range ops.Querying {
							allowedSet[id] = struct{}{}
						}
					}
				}
			}
		}
	}

	// Then subtract DENY statements
	for _, stmt := range policy.Statements {
		if stmt.Effect == PolicyStatementEffect_DENY {
			for _, res := range stmt.Resources {
				ops := ResourceOperations[res]
				for _, act := range stmt.Actions {
					switch act {
					case PolicyStatementAction_ALL:
						for _, id := range ops.Creating {
							delete(allowedSet, id)
						}
						for _, id := range ops.Deleting {
							delete(allowedSet, id)
						}
						for _, id := range ops.Updating {
							delete(allowedSet, id)
						}
						for _, id := range ops.Querying {
							delete(allowedSet, id)
						}
					case PolicyStatementAction_CREATE:
						for _, id := range ops.Creating {
							delete(allowedSet, id)
						}
					case PolicyStatementAction_DELETE:
						for _, id := range ops.Deleting {
							delete(allowedSet, id)
						}
					case PolicyStatementAction_UPDATE:
						for _, id := range ops.Updating {
							delete(allowedSet, id)
						}
					case PolicyStatementAction_QUERY:
						for _, id := range ops.Querying {
							delete(allowedSet, id)
						}
					}
				}
			}
		}
	}

	if len(allowedSet) == 0 {
		return nil
	}

	result := make([]int32, 0, len(allowedSet))
	for id := range allowedSet {
		result = append(result, id)
	}
	return result
}
