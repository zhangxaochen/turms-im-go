package element

import (
	"im.turms/server/internal/infra/ldap/asn1"
)

// LDAP tag constants (mapped from LdapTagConst.java and Asn1IdConst.java)
const (
	// Tag Class/Form
	FormConstructed     = 0x20
	TagClassApplication = 0x40
	TagClassContext     = 0x80

	// Operation Tags
	LdapTagBindRequest           = TagClassApplication | FormConstructed | 0
	LdapTagBindResponse          = TagClassApplication | FormConstructed | 1
	LdapTagUnbindRequest         = TagClassApplication | 2
	LdapTagSearchRequest         = TagClassApplication | FormConstructed | 3
	LdapTagSearchResultEntry     = TagClassApplication | FormConstructed | 4
	LdapTagSearchResultDone      = TagClassApplication | FormConstructed | 5
	LdapTagModifyRequest         = TagClassApplication | FormConstructed | 6
	LdapTagModifyResponse        = TagClassApplication | FormConstructed | 7
	LdapTagAddRequest            = TagClassApplication | FormConstructed | 8
	LdapTagAddResponse           = TagClassApplication | FormConstructed | 9
	LdapTagDelRequest            = TagClassApplication | 10
	LdapTagDelResponse           = TagClassApplication | FormConstructed | 11
	LdapTagModifyDNRequest       = TagClassApplication | FormConstructed | 12
	LdapTagModifyDNResponse      = TagClassApplication | FormConstructed | 13
	LdapTagCompareRequest        = TagClassApplication | FormConstructed | 14
	LdapTagCompareResponse       = TagClassApplication | FormConstructed | 15
	LdapTagAbandonRequest        = TagClassApplication | 16
	LdapTagSearchResultReference = TagClassApplication | FormConstructed | 19
	LdapTagExtendedRequest       = TagClassApplication | FormConstructed | 23
	LdapTagExtendedResponse      = TagClassApplication | FormConstructed | 24
	LdapTagIntermediateResponse  = TagClassApplication | FormConstructed | 25

	LdapTagControls = TagClassContext | FormConstructed | 0 // 0xA0
)

const (
	ResultCodeSuccess            = 0
	ResultCodeInvalidCredentials  = 49
	LdapVersion3                 = 3
)

// Control OID constants
// @MappedFrom ControlOidConst.java
const (
	ControlOidFastBind = "1.2.840.113556.1.4.1781"
)

// FastBindControl is the pre-built control for LDAP fast bind.
// @MappedFrom LdapClient.REQUEST_CONTROLS_FAST_BIND
var FastBindControl = Control{OID: ControlOidFastBind, Criticality: false}

// SearchRequestNoAttributes is the "1.1" OID meaning no attributes should be returned.
// @MappedFrom SearchRequest.NO_ATTRIBUTES
var SearchRequestNoAttributes = []string{"1.1"}

// DerefAliases is an alias for LdapDerefAliases to match Java naming.
type DerefAliases = LdapDerefAliases

// ProtocolOperation interface for all LDAP operations
type ProtocolOperation interface {
	WriteTo(buffer *asn1.BerBuffer)
	EstimateSize() int
}

// Attribute maps to Attribute in Java.
type Attribute struct {
	Type   string
	Values []string
}

func (a *Attribute) IsEmpty() bool {
	return len(a.Values) == 0
}

func (a *Attribute) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequence()
	buffer.WriteOctetString(a.Type)
	buffer.BeginSequenceWithTag(asn1.TagSequence)
	buffer.WriteOctetStrings(a.Values)
	buffer.EndSequence()
	buffer.EndSequence()
}

func DecodeAttribute(buffer *asn1.BerBuffer) *Attribute {
	buffer.SkipTagAndLength()
	attrType := buffer.ReadOctetString()
	tag := buffer.ReadTag()
	const tagSetConstructed = 0x31
	if int(tag) != tagSetConstructed {
		buffer.SkipLengthAndValue()
		return &Attribute{Type: attrType, Values: []string{}}
	}
	length := buffer.ReadLength()
	end := buffer.ReaderIndex() + length
	var values []string
	for buffer.IsReadableWithEnd(end) {
		values = append(values, buffer.ReadOctetString())
	}
	return &Attribute{Type: attrType, Values: values}
}

// LdapMessage maps to LdapMessage in Java.
type LdapMessage struct {
	MessageId         int
	ProtocolOperation ProtocolOperation
	Controls          []Control
}

func (m *LdapMessage) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequence()
	buffer.WriteInteger(m.MessageId)

	if m.ProtocolOperation == nil {
		panic("ProtocolOperation must not be nil")
	}
	m.ProtocolOperation.WriteTo(buffer)

	if len(m.Controls) > 0 {
		buffer.BeginSequenceWithTag(LdapTagControls)
		for _, c := range m.Controls {
			c.WriteTo(buffer)
		}
		buffer.EndSequence()
	}
	buffer.EndSequence()
}

// Control maps to Control in Java.
type Control struct {
	OID         string
	Criticality bool
	Value       []byte // Optional value octet string
}

func (c *Control) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequence()
	buffer.WriteOctetString(c.OID)
	if c.Criticality {
		buffer.WriteBoolean(true)
	}
	if len(c.Value) > 0 {
		buffer.WriteOctetStringBytes(c.Value)
	}
	buffer.EndSequence()
}

func DecodeControls(buffer *asn1.BerBuffer) []Control {
	if !buffer.IsReadable() || !buffer.PeekAndCheckTag(LdapTagControls) {
		return nil
	}
	buffer.ReadTag()
	length := buffer.ReadLength()
	end := buffer.ReaderIndex() + length
	var controls []Control
	for buffer.IsReadableWithEnd(end) {
		buffer.SkipTagAndLength()
		oid := buffer.ReadOctetString()
		var criticality bool
		if buffer.IsReadableWithEnd(end) && buffer.PeekAndCheckTag(asn1.TagBoolean) {
			criticality = buffer.ReadBoolean()
		}
		var value []byte
		if buffer.IsReadableWithEnd(end) && buffer.PeekAndCheckTag(asn1.TagOctetString) {
			value = []byte(buffer.ReadOctetString())
		}
		if oid != "" {
			controls = append(controls, Control{OID: oid, Criticality: criticality, Value: value})
		}
	}
	return controls
}

// LdapResult maps to the common LDAP result structure in responses.
type LdapResult struct {
	ResultCode        int
	MatchedDN         string
	DiagnosticMessage string
	Referrals         []string
}

func DecodeLdapResult(buffer *asn1.BerBuffer) LdapResult {
	res := LdapResult{
		ResultCode:        buffer.ReadEnumeration(),
		MatchedDN:         buffer.ReadOctetString(),
		DiagnosticMessage: buffer.ReadOctetString(),
	}
	if buffer.IsReadable() && buffer.PeekAndCheckTag(TagClassContext|FormConstructed|3) {
		buffer.ReadTag()
		length := buffer.ReadLength()
		end := buffer.ReaderIndex() + length
		for buffer.IsReadableWithEnd(end) {
			res.Referrals = append(res.Referrals, buffer.ReadOctetString())
		}
	}
	return res
}

func (r *LdapResult) IsSuccess() bool {
	return r.ResultCode == ResultCodeSuccess
}

// BindRequest
type BindRequest struct {
	Version  int
	Name     string
	Password string
}

func (r *BindRequest) EstimateSize() int {
	return 32 + len(r.Name) + len(r.Password)
}

func (r *BindRequest) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequenceWithTag(LdapTagBindRequest)
	buffer.WriteInteger(r.Version)
	buffer.WriteOctetString(r.Name)
	// simple auth (context tag 0)
	buffer.WriteOctetStringWithTag(TagClassContext|0, r.Password)
	buffer.EndSequence()
}

// BindResponse
type BindResponse struct {
	LdapResult
	ServerSaslCreds []byte // Optional
}

func DecodeBindResponse(buffer *asn1.BerBuffer) *BindResponse {
	buffer.SkipTag()
	buffer.SkipLength()
	res := &BindResponse{
		LdapResult: DecodeLdapResult(buffer),
	}
	if buffer.IsReadable() && buffer.PeekAndCheckTag(TagClassContext|7) {
		res.ServerSaslCreds = []byte(buffer.ReadOctetStringWithTag(TagClassContext | 7))
	}
	return res
}

// UnbindRequest
type UnbindRequest struct{}

func (r *UnbindRequest) EstimateSize() int { return 2 }
func (r *UnbindRequest) WriteTo(buffer *asn1.BerBuffer) {
	buffer.WriteTag(LdapTagUnbindRequest)
	buffer.WriteLength(0)
}

// SearchRequest
type LdapScope int

const (
	ScopeBaseObject   LdapScope = 0
	ScopeSingleLevel  LdapScope = 1
	ScopeWholeSubtree LdapScope = 2
)

type LdapDerefAliases int

const (
	DerefNever          LdapDerefAliases = 0
	DerefInSearching    LdapDerefAliases = 1
	DerefFindingBaseObj LdapDerefAliases = 2
	DerefAlways         LdapDerefAliases = 3
)

type SearchRequest struct {
	BaseDN     string
	Scope      LdapScope
	DerefAlias LdapDerefAliases
	SizeLimit  int
	TimeLimit  int
	TypesOnly  bool
	Attributes []string
	Filter     string
}

func (r *SearchRequest) EstimateSize() int {
	return 128 + len(r.Filter)
}

func (r *SearchRequest) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequenceWithTag(LdapTagSearchRequest)
	buffer.WriteOctetString(r.BaseDN)
	buffer.WriteEnumeration(int(r.Scope))
	buffer.WriteEnumeration(int(r.DerefAlias))
	buffer.WriteInteger(r.SizeLimit)
	buffer.WriteInteger(r.TimeLimit)
	buffer.WriteBoolean(r.TypesOnly)

	WriteFilter(buffer, r.Filter)

	buffer.BeginSequence()
	buffer.WriteOctetStrings(r.Attributes)
	buffer.EndSequence()

	buffer.EndSequence()
}

// SearchResultEntry
type SearchResultEntry struct {
	ObjectName string
	Attributes []Attribute
}

func DecodeSearchResultEntry(buffer *asn1.BerBuffer) *SearchResultEntry {
	buffer.SkipTag()
	buffer.SkipLength()
	objectName := buffer.ReadOctetString()
	buffer.SkipTag() // SEQUENCE OF partial attributes
	length := buffer.ReadLength()
	end := buffer.ReaderIndex() + length
	var attrs []Attribute
	for buffer.IsReadableWithEnd(end) {
		attrs = append(attrs, *DecodeAttribute(buffer))
	}
	return &SearchResultEntry{ObjectName: objectName, Attributes: attrs}
}

// SearchResult
type SearchResult struct {
	LdapResult
	Entries []SearchResultEntry
}

// ModifyOperation
type ModifyOperation int

const (
	ModifyAdd     ModifyOperation = 0
	ModifyDelete  ModifyOperation = 1
	ModifyReplace ModifyOperation = 2
)

type Change struct {
	Operation ModifyOperation
	Attribute Attribute
}

type ModifyRequest struct {
	DN      string
	Changes []Change
}

func (r *ModifyRequest) EstimateSize() int {
	return len(r.DN) + len(r.Changes)*32
}

func (r *ModifyRequest) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequenceWithTag(LdapTagModifyRequest)
	buffer.WriteOctetString(r.DN)
	buffer.BeginSequence()
	for _, c := range r.Changes {
		buffer.BeginSequence()
		buffer.WriteEnumeration(int(c.Operation))
		c.Attribute.WriteTo(buffer)
		buffer.EndSequence()
	}
	buffer.EndSequence()
	buffer.EndSequence()
}

type ModifyResponse struct {
	LdapResult
}

func DecodeModifyResponse(buffer *asn1.BerBuffer) *ModifyResponse {
	buffer.SkipTag()
	buffer.SkipLength()
	return &ModifyResponse{LdapResult: DecodeLdapResult(buffer)}
}

type LdapProtocolError struct {
	Msg string
}

func (e *LdapProtocolError) Error() string { return e.Msg }
