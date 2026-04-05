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
	buf                         []byte
	readerIdx                   int
	sequenceLengthWriterIndexes []int
	currentSequenceLengthIndex  int
	marks                       []int
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
		panic("Insufficient data: cannot read tag")
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
	if length <= 0x7F {
		b.buf = append(b.buf, byte(length))
	} else if length <= 0xFF {
		b.buf = append(b.buf, 0x81, byte(length))
	} else if length <= 0xFFFF {
		b.buf = append(b.buf, 0x82, byte(length>>8), byte(length))
	} else if length <= 0xFF_FFFF {
		b.buf = append(b.buf, 0x83, byte(length>>16), byte(length>>8), byte(length))
	} else {
		b.buf = append(b.buf, 0x84, byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
	}
}

// @MappedFrom readLength()
func (b *BerBuffer) ReadLength() int {
	if b.readerIdx >= len(b.buf) {
		panic("readLength: insufficient data")
	}
	lenByte := int(b.buf[b.readerIdx])
	b.readerIdx++
	if (lenByte & 0x80) == 0 {
		return lenByte
	}
	numBytes := lenByte & 0x7F
	if numBytes == 0 {
		panic("Indefinite length is not supported")
	}
	if numBytes > 4 {
		panic(fmt.Sprintf("The length (%d) is too long", numBytes))
	}
	if b.readerIdx+numBytes > len(b.buf) {
		panic("Insufficient data")
	}
	length := 0
	for i := 0; i < numBytes; i++ {
		length = (length << 8) | int(b.buf[b.readerIdx])
		b.readerIdx++
	}
	if length < 0 {
		panic("Invalid length bytes")
	}
	return length
}

// @MappedFrom tryReadLengthIfReadable()
func (b *BerBuffer) TryReadLengthIfReadable() int {
	if b.readerIdx >= len(b.buf) {
		return -1
	}
	lenByte := int(b.buf[b.readerIdx])
	if (lenByte & 0x80) == 0 {
		b.readerIdx++
		return lenByte
	}
	numBytes := lenByte & 0x7F
	if numBytes == 0 {
		panic("Indefinite length is not supported")
	}
	if numBytes > 4 {
		panic(fmt.Sprintf("The length (%d) is too long", numBytes))
	}
	if b.readerIdx+1+numBytes > len(b.buf) {
		return -1
	}
	b.readerIdx++
	length := 0
	for i := 0; i < numBytes; i++ {
		length = (length << 8) | int(b.buf[b.readerIdx])
		b.readerIdx++
	}
	if length < 0 {
		panic("Invalid length bytes")
	}
	return length
}

// @MappedFrom beginSequence()
func (b *BerBuffer) BeginSequence() {
	b.BeginSequenceWithTag(TagSequence | FormConstructed)
}

// @MappedFrom beginSequence(int tag)
func (b *BerBuffer) BeginSequenceWithTag(tag int) {
	b.buf = append(b.buf, byte(tag))
	// 3 = 1 (for the byte length of value length) + 2 (for the value length up to 64k)
	lengthIdx := len(b.buf)
	b.sequenceLengthWriterIndexes = append(b.sequenceLengthWriterIndexes, lengthIdx)
	b.currentSequenceLengthIndex++
	b.buf = append(b.buf, 0, 0, 0)
}

// @MappedFrom endSequence()
func (b *BerBuffer) EndSequence() {
	if b.currentSequenceLengthIndex <= 0 {
		panic("Unbalanced sequences")
	}
	b.currentSequenceLengthIndex--
	lengthIdx := b.sequenceLengthWriterIndexes[b.currentSequenceLengthIndex]
	b.sequenceLengthWriterIndexes = b.sequenceLengthWriterIndexes[:b.currentSequenceLengthIndex]

	valueWriterIndexStart := lengthIdx + 3
	valueLength := len(b.buf) - valueWriterIndexStart

	if valueLength <= 0xFFFF {
		b.buf[lengthIdx] = 0x82
		b.buf[lengthIdx+1] = byte(valueLength >> 8)
		b.buf[lengthIdx+2] = byte(valueLength)
	} else {
		panic(fmt.Sprintf("Expecting the sequence value length to be less than or equal to 64k, but got %d", valueLength))
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
		b.buf = append(b.buf, 0)
	}
}

// @MappedFrom readBoolean()
func (b *BerBuffer) ReadBoolean() bool {
	actualTag := b.ReadTag()
	if actualTag != TagBoolean {
		panic(fmt.Sprintf("Expecting tag: %d, but got: %d", TagBoolean, actualTag))
	}
	length := b.ReadLength()
	if length > 1 {
		panic("The boolean is too large")
	}
	if b.readerIdx >= len(b.buf) {
		panic("Insufficient data")
	}
	val := b.buf[b.readerIdx]
	b.readerIdx++
	return val != 0
}

// @MappedFrom writeInteger(int value)
func (b *BerBuffer) WriteInteger(value int) {
	b.WriteIntegerWithTag(TagInteger, value)
}

// @MappedFrom writeInteger(int tag, int value)
func (b *BerBuffer) WriteIntegerWithTag(tag int, value int) {
	b.buf = append(b.buf, byte(tag))
	v := uint32(value)
	if value < 0 {
		if (v & 0xFFFF_FF80) == 0xFFFF_FF80 {
			b.buf = append(b.buf, 1, byte(v))
		} else if (v & 0xFFFF_8000) == 0xFFFF_8000 {
			b.buf = append(b.buf, 2, byte(v>>8), byte(v))
		} else if (v & 0xFF80_0000) == 0xFF80_0000 {
			b.buf = append(b.buf, 3, byte(v>>16), byte(v>>8), byte(v))
		} else {
			b.buf = append(b.buf, 4, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		}
	} else {
		if (v & 0x0000_007F) == v {
			b.buf = append(b.buf, 1, byte(v))
		} else if (v & 0x0000_7FFF) == v {
			b.buf = append(b.buf, 2, byte(v>>8), byte(v))
		} else if (v & 0x007F_FFFF) == v {
			b.buf = append(b.buf, 3, byte(v>>16), byte(v>>8), byte(v))
		} else {
			b.buf = append(b.buf, 4, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		}
	}
}

// @MappedFrom readInteger()
func (b *BerBuffer) ReadInteger() int {
	return b.ReadIntWithTag(TagInteger)
}

// @MappedFrom readIntWithTag(int tag)
func (b *BerBuffer) ReadIntWithTag(tag int) int {
	actualTag := b.ReadTag()
	if actualTag != tag {
		panic(fmt.Sprintf("Expecting tag: %d, but got: %d", tag, actualTag))
	}
	length := b.ReadLength()
	if length == 0 {
		return 0
	}
	if length > 4 {
		panic("The integer is too long")
	}
	if b.readerIdx+length > len(b.buf) {
		panic("Insufficient data")
	}
	firstByte := b.buf[b.readerIdx]
	b.readerIdx++
	value := int(firstByte & 0x7F)
	for i := 1; i < length; i++ {
		value = (value << 8) | int(b.buf[b.readerIdx])
		b.readerIdx++
	}
	if (firstByte & 0x80) != 0 {
		value = -value
	}
	return value
}

// @MappedFrom writeOctetString(String value)
func (b *BerBuffer) WriteOctetString(value string) {
	b.WriteOctetStringWithTag(TagOctetString, value)
}

// @MappedFrom writeOctetString(byte[] value)
func (b *BerBuffer) WriteOctetStringBytes(value []byte) {
	b.buf = append(b.buf, TagOctetString)
	b.WriteLength(len(value))
	b.buf = append(b.buf, value...)
}

// @MappedFrom writeOctetString(int tag, byte[] value)
func (b *BerBuffer) WriteOctetStringBytesWithTag(tag int, value []byte) {
	b.buf = append(b.buf, byte(tag))
	b.WriteLength(len(value))
	b.buf = append(b.buf, value...)
}

// @MappedFrom writeOctetString(byte[] value, int start, int length)
func (b *BerBuffer) WriteOctetStringBytesWithStart(value []byte, start int, length int) {
	b.buf = append(b.buf, TagOctetString)
	b.WriteLength(length)
	b.buf = append(b.buf, value[start:start+length]...)
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
	lengthIdx := len(b.buf)
	b.buf = append(b.buf, 0, 0, 0)
	
	startIdx := len(b.buf)
	b.buf = append(b.buf, value...)
	valueLength := len(b.buf) - startIdx

	if valueLength <= 0xFFFF {
		b.buf[lengthIdx] = 0x82
		b.buf[lengthIdx+1] = byte(valueLength >> 8)
		b.buf[lengthIdx+2] = byte(valueLength)
	} else {
		panic(fmt.Sprintf("Expecting the sequence value length to be less than or equal to 64k, but got %d", valueLength))
	}
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
	actualTag := b.ReadTag()
	if actualTag != tag {
		panic(fmt.Sprintf("Encountered ASN.1 tag %d (expected tag %d)", actualTag, tag))
	}
	length := b.ReadLength()
	return b.ReadOctetStringWithLength(length)
}

// @MappedFrom readOctetStringWithLength(int length)
func (b *BerBuffer) ReadOctetStringWithLength(length int) string {
	if length == 0 {
		return ""
	}
	if b.readerIdx+length > len(b.buf) {
		panic("Insufficient data")
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
	return b.readerIdx < end
}

func (b *BerBuffer) ReaderIndex() int {
	return b.readerIdx
}

func (b *BerBuffer) WriterIndex() int {
	return len(b.buf)
}

func (b *BerBuffer) Append(data []byte) {
	b.buf = append(b.buf, data...)
}

func (b *BerBuffer) PeekTag() byte {
	if b.readerIdx >= len(b.buf) {
		return 0
	}
	return b.buf[b.readerIdx]
}

func (b *BerBuffer) MarkReaderIndex() {
	b.marks = append(b.marks, b.readerIdx)
}

func (b *BerBuffer) ResetReaderIndex() {
	if len(b.marks) > 0 {
		b.readerIdx = b.marks[len(b.marks)-1]
		b.marks = b.marks[:len(b.marks)-1]
	}
}

func (b *BerBuffer) Bytes() []byte {
	return b.buf
}

func (b *BerBuffer) IsReadableWithCount(n int) bool {
	return b.readerIdx+n <= len(b.buf)
}

func (b *BerBuffer) WriteTag(tag int) {
	b.buf = append(b.buf, byte(tag))
}

const TagSet = 0x31
