package element

import "im.turms/server/internal/infra/ldap/asn1"

// Attribute maps to Attribute in Java.
// @MappedFrom Attribute
type Attribute struct {
}

// @MappedFrom decode(BerBuffer buffer)
func (a *Attribute) Decode(buffer *asn1.BerBuffer) {
}

// LdapMessage maps to LdapMessage in Java.
// @MappedFrom LdapMessage
type LdapMessage struct {
}

// @MappedFrom estimateSize()
func (m *LdapMessage) EstimateSize() int {
	return 0
}

// @MappedFrom writeTo(BerBuffer buffer)
func (m *LdapMessage) WriteTo(buffer *asn1.BerBuffer) {
}

// LdapResult maps to LdapResult in Java.
// @MappedFrom LdapResult
type LdapResult struct {
}

// @MappedFrom isSuccess()
func (r *LdapResult) IsSuccess() bool {
	return false
}

// Control maps to Control in Java.
// @MappedFrom Control
type Control struct {
}

// @MappedFrom decode(BerBuffer buffer)
func (c *Control) Decode(buffer *asn1.BerBuffer) {
}

// BindRequest maps to BindRequest in Java.
// @MappedFrom BindRequest
type BindRequest struct {
}

// @MappedFrom estimateSize()
func (r *BindRequest) EstimateSize() int {
	return 0
}

// @MappedFrom writeTo(BerBuffer buffer)
func (r *BindRequest) WriteTo(buffer *asn1.BerBuffer) {
}

// BindResponse maps to BindResponse in Java.
// @MappedFrom BindResponse
type BindResponse struct {
}

// @MappedFrom decode(BerBuffer buffer)
func (r *BindResponse) Decode(buffer *asn1.BerBuffer) {
}

// ModifyRequest maps to ModifyRequest in Java.
// @MappedFrom ModifyRequest
type ModifyRequest struct {
}

// @MappedFrom estimateSize()
func (r *ModifyRequest) EstimateSize() int {
	return 0
}

// @MappedFrom writeTo(BerBuffer buffer)
func (r *ModifyRequest) WriteTo(buffer *asn1.BerBuffer) {
}

// ModifyResponse maps to ModifyResponse in Java.
// @MappedFrom ModifyResponse
type ModifyResponse struct {
}

// @MappedFrom decode(BerBuffer buffer)
func (r *ModifyResponse) Decode(buffer *asn1.BerBuffer) {
}

// Filter maps to Filter in Java.
// @MappedFrom Filter
type Filter struct {
}

// @MappedFrom write(BerBuffer buffer, String filter)
func (f *Filter) Write(buffer *asn1.BerBuffer, filter string) {
}

// SearchRequest maps to SearchRequest in Java.
// @MappedFrom SearchRequest
type SearchRequest struct {
}

// @MappedFrom estimateSize()
func (r *SearchRequest) EstimateSize() int {
	return 0
}

// @MappedFrom writeTo(BerBuffer buffer)
func (r *SearchRequest) WriteTo(buffer *asn1.BerBuffer) {
}

// SearchResult maps to SearchResult in Java.
// @MappedFrom SearchResult
type SearchResult struct {
}

// @MappedFrom decode(BerBuffer buffer)
func (r *SearchResult) Decode(buffer *asn1.BerBuffer) {
}

// @MappedFrom isComplete()
func (r *SearchResult) IsComplete() bool {
	return false
}
