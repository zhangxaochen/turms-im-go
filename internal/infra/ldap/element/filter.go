package element

import (
	"bytes"
	"fmt"

	"im.turms/server/internal/infra/ldap/asn1"
)

// Filter tags
const (
	typeAnd             = 0xA0
	typeOr              = 0xA1
	typeNot             = 0xA2
	typeEquality        = 0xA3
	typeSubstring       = 0xA4
	typeGreater         = 0xA5
	typeLess            = 0xA6
	typeApproximate     = 0xA8
	typeExtensibleMatch = 0xA9
	typePresent         = 0x87

	extensibleMatchingRule  = 0x81
	extensibleMatchingType  = 0x82
	extensibleMatchingValue = 0x83
	extensibleMatchingDN    = 0x84

	substringInitial = 0x80
	substringAny     = 0x81
	substringFinal   = 0x82
)

type filterContext struct {
	readIndex int
}

// WriteFilter encodes an LDAP filter string into a BerBuffer.
// @MappedFrom write(BerBuffer buffer, String filter)
func WriteFilter(buffer *asn1.BerBuffer, filterStr string) {
	filter := []byte(filterStr)
	writeFilterInner(buffer, filter, len(filter))
}

func writeFilterInner(buffer *asn1.BerBuffer, filter []byte, filterEndIndex int) {
	currentParensIndex := 0
	ctx := &filterContext{readIndex: 0}

	for ctx.readIndex < filterEndIndex {
		switch filter[ctx.readIndex] {
		case '(':
			ctx.readIndex++
			currentParensIndex++
			if ctx.readIndex >= filterEndIndex {
				panic("Unbalanced parenthesis")
			}
			switch filter[ctx.readIndex] {
			case '&':
				writeFilterSet(buffer, filter, typeAnd, ctx, filterEndIndex)
				currentParensIndex--
			case '|':
				writeFilterSet(buffer, filter, typeOr, ctx, filterEndIndex)
				currentParensIndex--
			case '!':
				writeFilterSet(buffer, filter, typeNot, ctx, filterEndIndex)
				currentParensIndex--
			default:
				balance := 1
				escape := false
				currentFilterEndIndex := ctx.readIndex
				for currentFilterEndIndex < filterEndIndex && balance > 0 {
					if !escape {
						if filter[currentFilterEndIndex] == '(' {
							balance++
						} else if filter[currentFilterEndIndex] == ')' {
							balance--
						}
					}
					escape = filter[currentFilterEndIndex] == '\\' && !escape
					if balance > 0 {
						currentFilterEndIndex++
					}
				}
				if balance != 0 {
					panic("Unbalanced parenthesis")
				}
				writeLeafFilter(buffer, filter, ctx.readIndex, currentFilterEndIndex)
				ctx.readIndex = currentFilterEndIndex + 1
				currentParensIndex--
			}
		case ')':
			buffer.EndSequence()
			ctx.readIndex++
			currentParensIndex--
		case ' ':
			ctx.readIndex++
		default:
			writeLeafFilter(buffer, filter, ctx.readIndex, filterEndIndex)
			ctx.readIndex = filterEndIndex
		}
		if currentParensIndex < 0 {
			panic("Unbalanced parenthesis")
		}
	}
	if currentParensIndex != 0 {
		panic("Unbalanced parenthesis")
	}
}

func writeFilterSet(buffer *asn1.BerBuffer, filter []byte, filterType int, ctx *filterContext, filterEnd int) {
	ctx.readIndex++
	buffer.BeginSequenceWithTag(filterType)
	closingParenIndex := findClosingParenIndex(filter, ctx.readIndex, filterEnd)
	writeFiltersInSet(buffer, filter, filterType, ctx.readIndex, closingParenIndex)
	ctx.readIndex = closingParenIndex + 1
	buffer.EndSequence()
}

func writeFiltersInSet(buffer *asn1.BerBuffer, filter []byte, filterType int, start int, end int) {
	readIdx := start
	currentFilterCount := 0
	for readIdx < end {
		c := filter[readIdx]
		if c == ' ' || c == '(' {
			if c == '(' {
				closingParenIndex := findClosingParenIndex(filter, readIdx+1, end)
				if filterType == typeNot && currentFilterCount > 0 {
					panic("The filter \"!\" cannot be followed by more than one filter")
				}
				length := closingParenIndex - readIdx
				writeFilterInner(buffer, filter[readIdx:closingParenIndex+1], length+1)
				currentFilterCount++
				readIdx = closingParenIndex + 1
			} else {
				readIdx++
			}
			continue
		}
		readIdx++
	}
}

func findClosingParenIndex(filter []byte, start int, end int) int {
	depth := 1
	escape := false
	closingParenIndex := start
	for closingParenIndex < end && depth > 0 {
		c := filter[closingParenIndex]
		if !escape {
			if c == '(' {
				depth++
			} else if c == ')' {
				depth--
			}
		}
		escape = c == '\\' && !escape
		if depth > 0 {
			closingParenIndex++
		}
	}
	if depth != 0 {
		panic("Unbalanced parenthesis")
	}
	return closingParenIndex
}

func writeLeafFilter(buffer *asn1.BerBuffer, filter []byte, start int, end int) {
	equalIndex := bytes.IndexByte(filter[start:end], '=')
	if equalIndex == -1 {
		panic("Missing \"=\"")
	}
	equalIndex += start

	filterType := -1
	filterTypeEndIndex := equalIndex

	if equalIndex > start {
		switch filter[equalIndex-1] {
		case '<':
			filterType = typeLess
			filterTypeEndIndex = equalIndex - 1
		case '>':
			filterType = typeGreater
			filterTypeEndIndex = equalIndex - 1
		case '~':
			filterType = typeApproximate
			filterTypeEndIndex = equalIndex - 1
		case ':':
			filterType = typeExtensibleMatch
			filterTypeEndIndex = equalIndex - 1
		default:
			filterType = -1
			filterTypeEndIndex = equalIndex
		}
	}

	valueStartIndex := equalIndex + 1

	if filterTypeEndIndex == equalIndex {
		// Potential equality, present, or substring
		if findUnescaped(filter, valueStartIndex, end) == -1 {
			filterType = typeEquality
		} else if end == valueStartIndex+1 && filter[valueStartIndex] == '*' {
			filterType = typePresent
			buffer.WriteOctetStringBytesWithStartAndTag(filterType, filter, start, filterTypeEndIndex-start)
			return
		} else {
			writeSubstringFilter(buffer, filter, start, filterTypeEndIndex, valueStartIndex, end)
			return
		}
	}

	if filterType == typeExtensibleMatch {
		writeExtensibleMatchFilter(buffer, filter, start, filterTypeEndIndex, valueStartIndex, end)
	} else {
		buffer.BeginSequenceWithTag(filterType)
		buffer.WriteOctetStringBytesWithStartAndTag(asn1.TagOctetString, filter, start, filterTypeEndIndex-start)

		unescaped := unescapeFilterValue(filter, valueStartIndex, end)
		if unescaped == nil {
			buffer.WriteOctetStringBytesWithStartAndTag(asn1.TagOctetString, filter, valueStartIndex, end-valueStartIndex)
		} else {
			buffer.WriteOctetStringBytes(unescaped)
		}
		buffer.EndSequence()
	}
}

func writeSubstringFilter(buffer *asn1.BerBuffer, filter []byte, typeStart, typeEnd, valueStart, valueEnd int) {
	buffer.BeginSequenceWithTag(typeSubstring)
	buffer.WriteOctetStringBytesWithStartAndTag(asn1.TagOctetString, filter, typeStart, typeEnd-typeStart)
	buffer.BeginSequence()

	isFirst := true
	previousIndex := valueStart
	for {
		index := findUnescaped(filter, previousIndex, valueEnd)
		if index == -1 {
			break
		}

		if previousIndex < index {
			// There is content before this star
			unescaped := unescapeFilterValue(filter, previousIndex, index)
			if isFirst {
				// First component: write as substringInitial
				if unescaped == nil {
					buffer.WriteOctetStringBytesWithStartAndTag(substringInitial, filter, previousIndex, index-previousIndex)
				} else {
					buffer.WriteOctetStringBytesWithTag(substringInitial, unescaped)
				}
			} else {
				// Middle component: write as substringAny
				if unescaped == nil {
					buffer.WriteOctetStringBytesWithStartAndTag(substringAny, filter, previousIndex, index-previousIndex)
				} else {
					buffer.WriteOctetStringBytesWithTag(substringAny, unescaped)
				}
			}
		}

		isFirst = false
		previousIndex = index + 1
	}

	// Last component (after the last star): write as substringFinal
	if previousIndex < valueEnd {
		unescaped := unescapeFilterValue(filter, previousIndex, valueEnd)
		if unescaped == nil {
			buffer.WriteOctetStringBytesWithStartAndTag(substringFinal, filter, previousIndex, valueEnd-previousIndex)
		} else {
			buffer.WriteOctetStringBytesWithTag(substringFinal, unescaped)
		}
	}

	buffer.EndSequence()
	buffer.EndSequence()
}

func writeExtensibleMatchFilter(buffer *asn1.BerBuffer, filter []byte, matchStart, matchEnd, valueStart, valueEnd int) {
	// Parse extensible match: [attr][:dn][:matchingRule]:=value
	// or :matchingRule:=value or :dn:=value
	attrStart := matchStart
	attrEnd := matchEnd
	var matchingRule []byte
	var attrType []byte
	matchDN := false

	// Find the colon that precedes '=' (extensible match separator)
	pos := matchStart
	for pos < matchEnd {
		if filter[pos] == ':' {
			colonPos := pos
			// Check for ":dn" flag
			if colonPos+3 <= matchEnd && bytes.Equal(filter[colonPos:colonPos+3], []byte(":dn")) {
				// Check if it's ":dn:" or ":dn" at end of attr portion
				if colonPos+3 >= matchEnd || filter[colonPos+3] == ':' {
					matchDN = true
					pos = colonPos + 3
					if pos < matchEnd && filter[pos] == ':' {
						pos++
					}
					continue
				}
			}
			// It could be a matching rule: ":matchingRule"
			ruleStart := colonPos + 1
			ruleEnd := matchEnd
			// Look for another colon after this one (for ":dn:matchingRule" or ":matchingRule:dn")
			nextColon := bytes.IndexByte(filter[ruleStart:matchEnd], ':')
			if nextColon != -1 {
				nextColon += ruleStart
				// Check if the part between colons is "dn"
				if bytes.Equal(filter[ruleStart:nextColon], []byte("dn")) {
					matchDN = true
					ruleEnd = ruleStart // no matching rule in this segment
					pos = nextColon + 1
					continue
				}
				ruleEnd = nextColon
			}
			matchingRule = filter[ruleStart:ruleEnd]
			attrEnd = colonPos
			break
		}
		pos++
	}

	// If there was a colon, extract attr type from before the first colon
	if attrEnd < matchEnd {
		attrType = filter[attrStart:attrEnd]
	}

	buffer.BeginSequenceWithTag(typeExtensibleMatch)

	// Write matching rule first (tag 0x81) if present
	if len(matchingRule) > 0 {
		buffer.WriteOctetStringBytesWithTag(extensibleMatchingRule, matchingRule)
	}

	// Write attribute type (tag 0x82) if present
	if len(attrType) > 0 {
		buffer.WriteOctetStringBytesWithTag(extensibleMatchingType, attrType)
	}

	// Write value (tag 0x83) - mandatory
	unescaped := unescapeFilterValue(filter, valueStart, valueEnd)
	if unescaped == nil {
		buffer.WriteOctetStringBytesWithStartAndTag(extensibleMatchingValue, filter, valueStart, valueEnd-valueStart)
	} else {
		buffer.WriteOctetStringBytesWithTag(extensibleMatchingValue, unescaped)
	}

	// Write dn flag (tag 0x84)
	buffer.WriteBooleanWithTag(extensibleMatchingDN, matchDN)

	buffer.EndSequence()
}

func unescapeFilterValue(filter []byte, start, end int) []byte {
	hasEscape := false
	for i := start; i < end; i++ {
		if filter[i] == '\\' {
			hasEscape = true
			break
		}
	}
	if !hasEscape {
		return nil
	}

	res := make([]byte, 0, end-start)
	for i := start; i < end; i++ {
		if filter[i] == '\\' && i+2 < end {
			var b byte
			fmt.Sscanf(string(filter[i+1:i+3]), "%02x", &b)
			res = append(res, b)
			i += 2
		} else {
			res = append(res, filter[i])
		}
	}
	return res
}

func findUnescaped(str []byte, start, end int) int {
	for i := start; i < end; i++ {
		if str[i] == '*' {
			backslashes := 0
			for j := i - 1; j >= start && str[j] == '\\'; j-- {
				backslashes++
			}
			if backslashes%2 == 0 {
				return i
			}
		}
	}
	return -1
}
