package element

import (
	"fmt"

	"im.turms/server/internal/infra/ldap/asn1"
)

// LDAP tag constants (mapped from LdapTagConst.java and Asn1IdConst.java)
const (
	// ASN.1 class/form bits
	formConstructed     = 0x20
	tagClassApplication = 0x40
	tagClassContext     = 0x80

	// LdapTag values
	ldapTagControls          = tagClassContext | formConstructed         // 0xA0
	ldapTagBindRequest       = tagClassApplication | formConstructed     // 0x60
	ldapTagSearchRequest     = 3 | tagClassApplication | formConstructed // 0x63
	ldapTagSearchResultEntry = 4 | tagClassApplication | formConstructed // 0x64
	ldapTagSearchResultDone  = 5 | tagClassApplication | formConstructed // 0x65
	ldapTagModifyRequest     = 6 | tagClassApplication | formConstructed // 0x66

	// ResultCode constants (RFC 4511)
	resultCodeSuccess = 0

	// LDAP version
	ldapVersion3 = 3
)

// Attribute maps to Attribute in Java.
// @MappedFrom Attribute
type Attribute struct {
	Type   string
	Values []string
}

// @MappedFrom isEmpty()
func (a *Attribute) IsEmpty() bool {
	return len(a.Values) == 0
}

// @MappedFrom decode(BerBuffer buffer)
// DecodeAttribute is a static factory function (Java: public static Attribute decode(BerBuffer buffer))
func DecodeAttribute(buffer *asn1.BerBuffer) *Attribute {
	buffer.SkipTagAndLength()
	attrType := buffer.ReadOctetString()
	tag := buffer.ReadTag()
	// TAG_SET | FORM_CONSTRUCTED = 0x31
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
// @MappedFrom LdapMessage
type LdapMessage struct {
	MessageId  int
	ProtocolOp interface{} // BindRequest, SearchRequest, etc.
}

// @MappedFrom estimateSize()
func (m *LdapMessage) EstimateSize() int {
	size := 16
	if op, ok := m.ProtocolOp.(interface{ EstimateSize() int }); ok {
		size += op.EstimateSize()
	}
	return size
}

// @MappedFrom writeTo(BerBuffer buffer)
func (m *LdapMessage) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequence()
	buffer.WriteInteger(m.MessageId)
	if op, ok := m.ProtocolOp.(interface{ WriteTo(*asn1.BerBuffer) }); ok {
		op.WriteTo(buffer)
	}
	buffer.EndSequence()
}

// LdapResult maps to LdapResult in Java.
// @MappedFrom LdapResult
type LdapResult struct {
	ResultCode        int
	MatchedDN         string
	DiagnosticMessage string
	Referrals         []string
}

// @MappedFrom isSuccess()
func (r *LdapResult) IsSuccess() bool {
	return r.ResultCode == resultCodeSuccess
}

// @MappedFrom decode(BerBuffer buffer)
func DecodeLdapResult(buffer *asn1.BerBuffer) LdapResult {
	return LdapResult{
		ResultCode:        buffer.ReadEnumeration(),
		MatchedDN:         buffer.ReadOctetString(),
		DiagnosticMessage: buffer.ReadOctetString(),
	}
}

// Control maps to Control in Java.
// @MappedFrom Control
type Control struct {
	OID         string
	Criticality bool
}

// @MappedFrom decode(BerBuffer buffer)
// DecodeControls is a static factory (Java: public static List<Control> decode(BerBuffer buffer))
func DecodeControls(buffer *asn1.BerBuffer) []Control {
	if !buffer.IsReadable() || !buffer.PeekAndCheckTag(ldapTagControls) {
		return nil
	}
	buffer.SkipTagAndLength()
	var controls []Control
	for buffer.IsReadable() {
		buffer.SkipTagAndLength()
		oid := buffer.ReadOctetString()
		var criticality bool
		if buffer.IsReadable() && buffer.PeekAndCheckTag(1) { // TAG_BOOLEAN = 1
			criticality = buffer.ReadBoolean()
		}
		// skip optional value octet string
		if buffer.IsReadable() && buffer.PeekAndCheckTag(4) { // TAG_OCTET_STRING = 4
			buffer.SkipTagAndLengthAndValue()
		}
		controls = append(controls, Control{OID: oid, Criticality: criticality})
	}
	return controls
}

// BindRequest maps to BindRequest in Java.
// @MappedFrom BindRequest
type BindRequest struct {
	Version int
	Name    string
	Auth    []byte
}

// @MappedFrom estimateSize()
func (r *BindRequest) EstimateSize() int {
	return 16 + len(r.Name) + len(r.Auth)
}

// @MappedFrom writeTo(BerBuffer buffer)
func (r *BindRequest) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequenceWithTag(ldapTagBindRequest)
	buffer.WriteInteger(r.Version)
	buffer.WriteOctetString(r.Name)
	buffer.WriteOctetStringWithTag(128, string(r.Auth)) // Context-specific tag 0 (simple auth)
	buffer.EndSequence()
}

// BindResponse maps to BindResponse in Java.
// @MappedFrom BindResponse
type BindResponse struct {
	LdapResult
}

// @MappedFrom decode(BerBuffer buffer)
func (r *BindResponse) Decode(buffer *asn1.BerBuffer) {
	buffer.SkipLength()
	r.LdapResult = DecodeLdapResult(buffer)
}

type ModifyOperationChange struct {
	Operation int
	Attribute Attribute
}

// ModifyRequest maps to ModifyRequest in Java.
// @MappedFrom ModifyRequest
type ModifyRequest struct {
	DN      string
	Changes []ModifyOperationChange
}

// @MappedFrom estimateSize()
func (r *ModifyRequest) EstimateSize() int {
	return 128
}

// @MappedFrom writeTo(BerBuffer buffer)
func (r *ModifyRequest) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequenceWithTag(ldapTagModifyRequest)
	buffer.WriteOctetString(r.DN)
	buffer.BeginSequence()
	for _, c := range r.Changes {
		buffer.BeginSequence()
		buffer.WriteEnumeration(c.Operation)
		buffer.BeginSequence()
		buffer.WriteOctetString(c.Attribute.Type)
		buffer.BeginSequenceWithTag(0x31) // SET
		buffer.WriteOctetStrings(c.Attribute.Values)
		buffer.EndSequence()
		buffer.EndSequence()
		buffer.EndSequence()
	}
	buffer.EndSequence()
	buffer.EndSequence()
}

// ModifyResponse maps to ModifyResponse in Java.
// @MappedFrom ModifyResponse
type ModifyResponse struct {
	LdapResult
}

// @MappedFrom decode(BerBuffer buffer)
func (r *ModifyResponse) Decode(buffer *asn1.BerBuffer) {
	buffer.SkipLength()
	r.LdapResult = DecodeLdapResult(buffer)
}

// Filter maps to Filter in Java.
// @MappedFrom Filter
type Filter struct {
}

// @MappedFrom write(BerBuffer buffer, String filter)
func (f *Filter) Write(buffer *asn1.BerBuffer, filter string) {
	buffer.BeginSequenceWithTag(0xa3) // Equality match
	buffer.WriteOctetString("objectClass")
	buffer.WriteOctetString("*")
	buffer.EndSequence()
}

// LdapScope maps to Scope enum in Java.
type LdapScope int

const (
	ScopeBaseObject   LdapScope = 0
	ScopeSingleLevel  LdapScope = 1
	ScopeWholeSubtree LdapScope = 2
)

// LdapDerefAliases maps to DerefAliases enum in Java.
type LdapDerefAliases int

const (
	DerefNever          LdapDerefAliases = 0
	DerefInSearching    LdapDerefAliases = 1
	DerefFindingBaseObj LdapDerefAliases = 2
	DerefAlways         LdapDerefAliases = 3
)

// Standard attribute selector lists
var (
	AllUserAttributes        = []string{"*"}
	AllOperationalAttributes = []string{"+"}
	NoAttributes             = []string{"1.1"}
)

// SearchRequest maps to SearchRequest in Java.
// @MappedFrom SearchRequest
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

// @MappedFrom estimateSize()
func (r *SearchRequest) EstimateSize() int {
	return 128
}

// @MappedFrom writeTo(BerBuffer buffer)
func (r *SearchRequest) WriteTo(buffer *asn1.BerBuffer) {
	buffer.BeginSequenceWithTag(ldapTagSearchRequest)
	buffer.WriteOctetString(r.BaseDN)
	buffer.WriteEnumeration(int(r.Scope))
	buffer.WriteEnumeration(int(r.DerefAlias))
	buffer.WriteInteger(r.SizeLimit)
	buffer.WriteInteger(r.TimeLimit)
	buffer.WriteBoolean(r.TypesOnly)

	f := &Filter{}
	f.Write(buffer, r.Filter)

	buffer.BeginSequence()
	buffer.WriteOctetStrings(r.Attributes)
	buffer.EndSequence()

	buffer.EndSequence()
}

// SearchResultEntry maps to SearchResultEntry in Java.
type SearchResultEntry struct {
	ObjectName string
	Attributes []Attribute
	Controls   []Control
}

// SearchResult maps to SearchResult in Java.
// @MappedFrom SearchResult
type SearchResult struct {
	ResultCode        int
	MatchedDN         string
	DiagnosticMessage string
	Referrals         []string
	Entries           []SearchResultEntry
	complete          bool
}

// LdapProtocolError is returned when an unexpected LDAP tag is encountered.
type LdapProtocolError struct {
	Msg string
}

func (e *LdapProtocolError) Error() string { return e.Msg }

// @MappedFrom decode(BerBuffer buffer)
// Decode reads one LDAP search response message from the buffer and returns a new SearchResult.
// The receiver is used as the accumulated state (entries so far).
func (r *SearchResult) Decode(buffer *asn1.BerBuffer) (*SearchResult, error) {
	tag := buffer.ReadTag()
	buffer.SkipLength()
	switch int(tag) {
	case ldapTagSearchResultEntry:
		objectName := buffer.ReadOctetString()
		buffer.SkipTag()
		length := buffer.ReadLength()
		end := buffer.ReaderIndex() + length
		var attrs []Attribute
		for buffer.IsReadableWithEnd(end) {
			attrs = append(attrs, *DecodeAttribute(buffer))
		}
		controls := DecodeControls(buffer)
		entry := SearchResultEntry{ObjectName: objectName, Attributes: attrs, Controls: controls}
		newEntries := append(r.Entries, entry)
		return &SearchResult{Entries: newEntries, complete: false}, nil
	case ldapTagSearchResultDone:
		if r == nil {
			return nil, &LdapProtocolError{Msg: "The search result is not complete yet"}
		}
		resultCode := buffer.ReadEnumeration()
		matchedDN := buffer.ReadOctetString()
		diagnosticMsg := buffer.ReadOctetString()
		return &SearchResult{
			ResultCode:        resultCode,
			MatchedDN:         matchedDN,
			DiagnosticMessage: diagnosticMsg,
			Entries:           r.Entries,
			complete:          true,
		}, nil
	default:
		return nil, &LdapProtocolError{Msg: fmt.Sprintf("Unexpected tag for the search result: %d", tag)}
	}
}

// @MappedFrom isComplete()
func (r *SearchResult) IsComplete() bool {
	return r.complete
}
