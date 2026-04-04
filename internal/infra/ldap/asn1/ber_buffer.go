package asn1

import (
	"fmt"
)

// LDAP tag constants
const (
	TagBoolean      = 0x01
	TagInteger      = 0x02
	TagOctetString  = 0x04
	TagEnumerated   = 0x0A
	TagSequence     = 0x30
	FormConstructed = 0x20
)

// BerBuffer maps to BerBuffer in Java.
// @MappedFrom BerBuffer
type BerBuffer struct {
	buf                        []byte
	readerIdx                  int
	sequenceLengthWriterIndexes []int
	currentSequenceLengthIndex  int
}

func NewBerBuffer(initialCapacity int) *BerBuffer {
	return &BerBuffer{
		buf:                        make([]byte, 0, initialCapacity),
		readerIdx:                  0,
		sequenceLengthWriterIndexes: make([]int, 0, 8),
		currentSequenceLengthIndex:  0,
	}
}

func NewBerBufferFromBytes(data []byte) *BerBuffer {
	return &BerBuffer{
		buf:       data,
		readerIdx: 0,
	}
}

// @MappedFrom skipTag()
func (b *BerBuffer) SkipTag() {
	b.SkipBytes(1)
}

// @MappedFrom skipTagAndLength()
func (b *BerBuffer) SkipTagAndLength() {
	b.SkipBytes(1)
	b.ReadLength()
}

// @MappedFrom skipTagAndLengthAndValue()
func (b *BerBuffer) SkipTagAndLengthAndValue() {
	b.SkipBytes(1)
	length := b.ReadLength()
	b.SkipBytes(length)
}

// @MappedFrom readTag()
func (b *BerBuffer) ReadTag() int {
	if b.readerIdx >= len(b.buf) {
		return 0
	}
	tag := int(b.buf[b.readerIdx])
	b.readerIdx++
	return tag
}

// @MappedFrom peekAndCheckTag(int tag)
func (b *BerBuffer) PeekAndCheckTag(tag int) bool {
	if b.readerIdx < len(b.buf) {
		return int(b.buf[b.readerIdx]) == tag
	}
	return false
}

// @MappedFrom skipLength()
func (b *BerBuffer) SkipLength() {
	b.ReadLength()
}

// @MappedFrom skipLengthAndValue()
func (b *BerBuffer) SkipLengthAndValue() {
	length := b.ReadLength()
	b.SkipBytes(length)
}

// @MappedFrom writeLength(int length)
func (b *BerBuffer) WriteLength(length int) {
	if length <= 127 {
		b.buf = append(b.buf, byte(length))
	} else if length <= 255 {
		b.buf = append(b.buf, 0x81, byte(length))
	} else if length <= 65535 {
		b.buf = append(b.buf, 0x82, byte(length>>8), byte(length))
	} else if length <= 16777215 {
		b.buf = append(b.buf, 0x83, byte(length>>16), byte(length>>8), byte(length))
	} else {
		b.buf = append(b.buf, 0x84, byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
	}
}

// @MappedFrom readLength()
func (b *BerBuffer) ReadLength() int {
	if b.readerIdx >= len(b.buf) {
		return 0
	}
	lenByte := int(b.buf[b.readerIdx])
	b.readerIdx++
	if (lenByte & 0x80) == 0 {
		return lenByte
	}
	numBytes := lenByte & 0x7F
	if numBytes == 0 {
		return 0 // indefinite length not fully supported, assume 0 for structural simplicity
	}
	length := 0
	for i := 0; i < numBytes; i++ {
		if b.readerIdx >= len(b.buf) {
			break
		}
		length = (length << 8) | int(b.buf[b.readerIdx])
		b.readerIdx++
	}
	return length
}

// @MappedFrom tryReadLengthIfReadable()
func (b *BerBuffer) TryReadLengthIfReadable() int {
	if b.readerIdx >= len(b.buf) {
		return -1
	}
	return b.ReadLength()
}

// @MappedFrom beginSequence()
func (b *BerBuffer) BeginSequence() {
	b.BeginSequenceWithTag(TagSequence)
}

// @MappedFrom beginSequence(int tag)
func (b *BerBuffer) BeginSequenceWithTag(tag int) {
	b.buf = append(b.buf, byte(tag))
	// Reserve 2 bytes for length (0x82, LengthHigh, LengthLow = 3 bytes total)
	b.sequenceLengthWriterIndexes = append(b.sequenceLengthWriterIndexes, len(b.buf))
	b.currentSequenceLengthIndex++
	b.buf = append(b.buf, 0x82, 0, 0)
}

// @MappedFrom endSequence()
func (b *BerBuffer) EndSequence() {
	if b.currentSequenceLengthIndex > 0 {
		b.currentSequenceLengthIndex--
		lengthIdx := b.sequenceLengthWriterIndexes[b.currentSequenceLengthIndex]
		b.sequenceLengthWriterIndexes = b.sequenceLengthWriterIndexes[:b.currentSequenceLengthIndex]
		
		// The length of the sequence is the total buffer size minus the start of the sequence body
		sequenceLength := len(b.buf) - (lengthIdx + 3)
		
		if sequenceLength <= 65535 {
			b.buf[lengthIdx] = 0x82
			b.buf[lengthIdx+1] = byte(sequenceLength >> 8)
			b.buf[lengthIdx+2] = byte(sequenceLength)
		} else {
			// This simplified version only reserves 2 bytes for length.
			// Production implementation would shift bytes if length > 65535.
			// Doing best effort for now by capping or expecting standard payload sizes.
			panic(fmt.Sprintf("Sequence too long: %d", sequenceLength))
		}
	}
}

// @MappedFrom writeBoolean(boolean value)
func (b *BerBuffer) WriteBoolean(value bool) {
	b.WriteBooleanWithTag(TagBoolean, value)
}

// @MappedFrom writeBoolean(int tag, boolean value)
func (b *BerBuffer) WriteBooleanWithTag(tag int, value bool) {
	b.buf = append(b.buf, byte(tag), 1)
	if value {
		b.buf = append(b.buf, 0xFF)
	} else {
		b.buf = append(b.buf, 0x00)
	}
}

// @MappedFrom readBoolean()
func (b *BerBuffer) ReadBoolean() bool {
	b.SkipTagAndLength()
	if b.readerIdx < len(b.buf) {
		val := b.buf[b.readerIdx]
		b.readerIdx++
		return val != 0
	}
	return false
}

// @MappedFrom writeInteger(int value)
func (b *BerBuffer) WriteInteger(value int) {
	b.WriteIntegerWithTag(TagInteger, value)
}

// @MappedFrom writeInteger(int tag, int value)
func (b *BerBuffer) WriteIntegerWithTag(tag int, value int) {
	b.buf = append(b.buf, byte(tag))
	if value >= -128 && value <= 127 {
		b.buf = append(b.buf, 1, byte(value))
	} else if value >= -32768 && value <= 32767 {
		b.buf = append(b.buf, 2, byte(value>>8), byte(value))
	} else if value >= -8388608 && value <= 8388607 {
		b.buf = append(b.buf, 3, byte(value>>16), byte(value>>8), byte(value))
	} else {
		b.buf = append(b.buf, 4, byte(value>>24), byte(value>>16), byte(value>>8), byte(value))
	}
}

// @MappedFrom readInteger()
func (b *BerBuffer) ReadInteger() int {
	return b.ReadIntWithTag(TagInteger)
}

// @MappedFrom readIntWithTag(int tag)
func (b *BerBuffer) ReadIntWithTag(tag int) int {
	readT := b.ReadTag()
	if readT != tag {
		return 0
	}
	length := b.ReadLength()
	if length == 0 || b.readerIdx+length > len(b.buf) {
		return 0
	}
	
	val := 0
	isNegative := (b.buf[b.readerIdx] & 0x80) != 0
	for i := 0; i < length; i++ {
		val = (val << 8) | int(b.buf[b.readerIdx])
		b.readerIdx++
	}
	if isNegative {
		val |= ^((1 << (length * 8)) - 1)
	}
	return val
}

// @MappedFrom writeOctetString(String value)
func (b *BerBuffer) WriteOctetString(value string) {
	b.WriteOctetStringWithTag(TagOctetString, value)
}

// @MappedFrom writeOctetString(byte[] value)
func (b *BerBuffer) WriteOctetStringBytes(value []byte) {
	b.WriteOctetStringBytesWithTag(TagOctetString, value)
}

// @MappedFrom writeOctetString(int tag, byte[] value)
func (b *BerBuffer) WriteOctetStringBytesWithTag(tag int, value []byte) {
	b.buf = append(b.buf, byte(tag))
	b.WriteLength(len(value))
	b.buf = append(b.buf, value...)
}

// @MappedFrom writeOctetString(byte[] value, int start, int length)
func (b *BerBuffer) WriteOctetStringBytesWithStart(value []byte, start int, length int) {
	b.WriteOctetStringBytesWithStartAndTag(TagOctetString, value, start, length)
}

// @MappedFrom writeOctetString(int tag, byte[] value, int start, int length)
func (b *BerBuffer) WriteOctetStringBytesWithStartAndTag(tag int, value []byte, start int, length int) {
	b.buf = append(b.buf, byte(tag))
	b.WriteLength(length)
	b.buf = append(b.buf, value[start:start+length]...)
}

// @MappedFrom writeOctetString(int tag, String value)
func (b *BerBuffer) WriteOctetStringWithTag(tag int, value string) {
	b.buf = append(b.buf, byte(tag))
	b.WriteLength(len(value))
	b.buf = append(b.buf, value...)
}

// @MappedFrom writeOctetStrings(List<String> values)
func (b *BerBuffer) WriteOctetStrings(values []string) {
	for _, v := range values {
		b.WriteOctetStringWithTag(TagOctetString, v)
	}
}

// @MappedFrom readOctetString()
func (b *BerBuffer) ReadOctetString() string {
	return b.ReadOctetStringWithTag(TagOctetString)
}

// @MappedFrom readOctetStringWithTag(int tag)
func (b *BerBuffer) ReadOctetStringWithTag(tag int) string {
	readT := b.ReadTag()
	if readT != tag {
		return ""
	}
	length := b.ReadLength()
	return b.ReadOctetStringWithLength(length)
}

// @MappedFrom readOctetStringWithLength(int length)
func (b *BerBuffer) ReadOctetStringWithLength(length int) string {
	if length == 0 || b.readerIdx+length > len(b.buf) {
		return ""
	}
	str := string(b.buf[b.readerIdx : b.readerIdx+length])
	b.readerIdx += length
	return str
}

// @MappedFrom writeEnumeration(int value)
func (b *BerBuffer) WriteEnumeration(value int) {
	b.WriteIntegerWithTag(TagEnumerated, value)
}

// @MappedFrom readEnumeration()
func (b *BerBuffer) ReadEnumeration() int {
	return b.ReadIntWithTag(TagEnumerated)
}

// @MappedFrom getBytes()
func (b *BerBuffer) GetBytes() []byte {
	return b.buf
}

// @MappedFrom skipBytes(int length)
func (b *BerBuffer) SkipBytes(length int) {
	if b.readerIdx+length <= len(b.buf) {
		b.readerIdx += length
	} else {
		b.readerIdx = len(b.buf)
	}
}

// @MappedFrom close()
func (b *BerBuffer) Close() {
	// Releasing buffers to pool can be added here
}

// @MappedFrom refCnt()
func (b *BerBuffer) RefCnt() int {
	return 1
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
func (b *BerBuffer) TouchHint(hint interface{}) {
}

// @MappedFrom release()
func (b *BerBuffer) Release() bool {
	return true
}

// @MappedFrom release(int decrement)
func (b *BerBuffer) ReleaseDecrement(decrement int) bool {
	return true
}

// @MappedFrom isReadable(int length)
func (b *BerBuffer) IsReadableLen(length int) bool {
	return b.readerIdx+length <= len(b.buf)
}

// @MappedFrom isReadable()
func (b *BerBuffer) IsReadable() bool {
	return b.readerIdx < len(b.buf)
}

// @MappedFrom isReadableWithEnd(int end)
func (b *BerBuffer) IsReadableWithEnd(end int) bool {
	return b.readerIdx < end && b.readerIdx < len(b.buf)
}

// @MappedFrom readerIndex()
func (b *BerBuffer) ReaderIndex() int {
	return b.readerIdx
}
