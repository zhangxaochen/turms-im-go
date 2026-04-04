package asn1

// BerBuffer maps to BerBuffer in Java.
// @MappedFrom BerBuffer
type BerBuffer struct {
}

// @MappedFrom skipTag()
func (b *BerBuffer) SkipTag() {
}

// @MappedFrom skipTagAndLength()
func (b *BerBuffer) SkipTagAndLength() {
}

// @MappedFrom skipTagAndLengthAndValue()
func (b *BerBuffer) SkipTagAndLengthAndValue() {
}

// @MappedFrom readTag()
func (b *BerBuffer) ReadTag() int {
	return 0
}

// @MappedFrom peekAndCheckTag(int tag)
func (b *BerBuffer) PeekAndCheckTag(tag int) bool {
	return false
}

// @MappedFrom skipLength()
func (b *BerBuffer) SkipLength() {
}

// @MappedFrom skipLengthAndValue()
func (b *BerBuffer) SkipLengthAndValue() {
}

// @MappedFrom writeLength(int length)
func (b *BerBuffer) WriteLength(length int) {
}

// @MappedFrom readLength()
func (b *BerBuffer) ReadLength() int {
	return 0
}

// @MappedFrom tryReadLengthIfReadable()
func (b *BerBuffer) TryReadLengthIfReadable() int {
	return 0
}

// @MappedFrom beginSequence()
func (b *BerBuffer) BeginSequence() {
}

// @MappedFrom beginSequence(int tag)
func (b *BerBuffer) BeginSequenceWithTag(tag int) {
}

// @MappedFrom endSequence()
func (b *BerBuffer) EndSequence() {
}

// @MappedFrom writeBoolean(boolean value)
func (b *BerBuffer) WriteBoolean(value bool) {
}

// @MappedFrom writeBoolean(int tag, boolean value)
func (b *BerBuffer) WriteBooleanWithTag(tag int, value bool) {
}

// @MappedFrom readBoolean()
func (b *BerBuffer) ReadBoolean() bool {
	return false
}

// @MappedFrom writeInteger(int value)
func (b *BerBuffer) WriteInteger(value int) {
}

// @MappedFrom writeInteger(int tag, int value)
func (b *BerBuffer) WriteIntegerWithTag(tag int, value int) {
}

// @MappedFrom readInteger()
func (b *BerBuffer) ReadInteger() int {
	return 0
}

// @MappedFrom readIntWithTag(int tag)
func (b *BerBuffer) ReadIntWithTag(tag int) int {
	return 0
}

// @MappedFrom writeOctetString(String value)
func (b *BerBuffer) WriteOctetString(value string) {
}

// @MappedFrom writeOctetString(byte[] value)
func (b *BerBuffer) WriteOctetStringBytes(value []byte) {
}

// @MappedFrom writeOctetString(int tag, byte[] value)
func (b *BerBuffer) WriteOctetStringBytesWithTag(tag int, value []byte) {
}

// @MappedFrom writeOctetString(byte[] value, int start, int length)
func (b *BerBuffer) WriteOctetStringBytesRange(value []byte, start int, length int) {
}

// @MappedFrom writeOctetString(int tag, byte[] value, int start, int length)
func (b *BerBuffer) WriteOctetStringBytesRangeWithTag(tag int, value []byte, start int, length int) {
}

// @MappedFrom writeOctetString(int tag, String value)
func (b *BerBuffer) WriteOctetStringWithTag(tag int, value string) {
}

// @MappedFrom writeOctetStrings(List<String> values)
func (b *BerBuffer) WriteOctetStrings(values []string) {
}

// @MappedFrom readOctetString()
func (b *BerBuffer) ReadOctetString() string {
	return ""
}

// @MappedFrom readOctetStringWithTag(int tag)
func (b *BerBuffer) ReadOctetStringWithTag(tag int) string {
	return ""
}

// @MappedFrom readOctetStringWithLength(int length)
func (b *BerBuffer) ReadOctetStringWithLength(length int) string {
	return ""
}

// @MappedFrom writeEnumeration(int value)
func (b *BerBuffer) WriteEnumeration(value int) {
}

// @MappedFrom readEnumeration()
func (b *BerBuffer) ReadEnumeration() int {
	return 0
}

// @MappedFrom getBytes()
func (b *BerBuffer) GetBytes() []byte {
	return nil
}

// @MappedFrom skipBytes(int length)
func (b *BerBuffer) SkipBytes(length int) {
}

// @MappedFrom refCnt()
func (b *BerBuffer) RefCnt() int {
	return 0
}

// @MappedFrom retain()
func (b *BerBuffer) Retain() {
}

// @MappedFrom retain(int increment)
func (b *BerBuffer) RetainIncrement(increment int) {
}

// @MappedFrom touch()
func (b *BerBuffer) Touch() {
}

// @MappedFrom touch(Object hint)
func (b *BerBuffer) TouchWithHint(hint any) {
}

// @MappedFrom release()
func (b *BerBuffer) Release() bool {
	return false
}

// @MappedFrom release(int decrement)
func (b *BerBuffer) ReleaseDecrement(decrement int) bool {
	return false
}

// @MappedFrom isReadable(int length)
func (b *BerBuffer) IsReadableLen(length int) bool {
	return false
}

// @MappedFrom isReadable()
func (b *BerBuffer) IsReadable() bool {
	return false
}

// @MappedFrom isReadableWithEnd(int end)
func (b *BerBuffer) IsReadableWithEnd(end int) bool {
	return false
}

// @MappedFrom readerIndex()
func (b *BerBuffer) ReaderIndex() int {
	return 0
}
