package kstrings

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

const (
	DELIM_START = '{'
	DELIM_STOP  = '}'
	DELIM_STR   = "{}"
	ESCAPE_CHAR = '\\'
)

type FormattingTuple struct {
	message   string
	throwable error
	args      []any
}

func NewFormattingTuple(message string, args []any, throwable error) *FormattingTuple {
	return &FormattingTuple{message: message, args: args, throwable: throwable}
}

func (its *FormattingTuple) Message() string  { return its.message }
func (its *FormattingTuple) Args() []any      { return its.args }
func (its *FormattingTuple) Throwable() error { return its.throwable }

////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////

func Println(messagePattern string, args ...any) {
	log.Println(FormatString(messagePattern, args...))
}

func FormatString(messagePattern string, args ...any) string {
	return SliceFormat(messagePattern, args...).Message()
}

// @bref Performs single argument substitution for the 'messagePattern' passed as
// parameter.
//
// For example,
//
// MessageFormatter.format(&quot;Hi {}.&quot;, &quot;there&quot;);
//
// will return the string "Hi there.".
//
// @param messagePattern
//
//	The message pattern which will be parsed and formatted
//
// @param arg
//
//	The argument to be substituted in place of the formatting anchor
//
// @return The formatted message
func Format(messagePattern string, arg any) *FormattingTuple {
	return ArrayFormat(messagePattern, []any{arg})
}

// Performs a two argument substitution for the 'messagePattern' passed as
// parameter.
//
// For example,
//
// MessageFormatter.format(&quot;Hi {}. My name is {}.&quot;, &quot;Alice&quot;, &quot;Bob&quot;);
//
// will return the string "Hi Alice. My name is Bob.".
//
// @param messagePattern
//
//	The message pattern which will be parsed and formatted
//
// @param arg1
//
//	The argument to be substituted in place of the first formatting
//	anchor
//
// @param arg2
//
//	The argument to be substituted in place of the second formatting
//	anchor
//
// @return The formatted message
func FormatArgs(messagePattern string, arg1, arg2 any) *FormattingTuple {
	return ArrayFormat(messagePattern, []any{arg1, arg2})
}

func SliceFormat(messagePattern string, args ...any) *FormattingTuple {
	return ArrayFormat(messagePattern, args)
}

func ArrayFormat(messagePattern string, argArray []any) *FormattingTuple {
	throwableCandidate := ThrowableCandidate(argArray)
	args := argArray
	// var err error = nil
	if throwableCandidate != nil {
		args, _ = TrimmedCopy(argArray)
	}
	return ArrayFormatWithError(messagePattern, args, throwableCandidate)
}

func ArrayFormatWithError(messagePattern string, argArray []any, throwable error) *FormattingTuple {
	if TrimSpace(messagePattern) == "" {
		return NewFormattingTuple("", argArray, throwable)
	}

	if argArray == nil {
		return NewFormattingTuple(messagePattern, nil, nil)
	}

	i, j, y := int(0), int(0), int(0)
	var sbuf strings.Builder
	sbuf.Grow(len(messagePattern) + 50)

	for y = 0; y < len(argArray); y = y + 1 {

		j = IndexOf(messagePattern, DELIM_STR, i)

		if j == -1 {
			// no more variables
			if i == 0 { // this is a simple string
				return NewFormattingTuple(messagePattern, argArray, throwable)
			} else { // add the tail string which contains no variables and return
				// the result.
				sbuf.WriteString(messagePattern[i:])
				//sbuf.append(messagePattern, i, messagePattern.length())
				return NewFormattingTuple(sbuf.String(), argArray, throwable)
			}
		} else {
			if isEscapedDelimeter(messagePattern, j) {
				if !isDoubleEscaped(messagePattern, j) {
					y-- // DELIM_START was escaped, thus should not be incremented
					sbuf.WriteString(messagePattern[i : j-1])
					//sbuf.append(messagePattern, i, j-1)
					sbuf.WriteByte(DELIM_START)
					// sbuf.append(DELIM_START)
					i = j + 1
				} else {
					// The escape character preceding the delimiter start is
					// itself escaped: "abc x:\\{}"
					// we have to consume one backward slash
					sbuf.WriteString(messagePattern[i : j-1])
					// sbuf.append(messagePattern, i, j-1)
					paraMap := make(map[any]any)
					deeplyAppendParameter(&sbuf, argArray[y], paraMap)
					i = j + 2
				}
			} else {
				// normal case
				sbuf.WriteString(messagePattern[i:j])
				// sbuf.append(messagePattern, i, j)
				paraMap := make(map[any]any)
				deeplyAppendParameter(&sbuf, argArray[y], paraMap)
				i = j + 2
			}
		}
	}
	// append the characters following the last {} pair.
	sbuf.WriteString(messagePattern[i:])
	// sbuf.append(messagePattern, i, messagePattern.length())
	return NewFormattingTuple(sbuf.String(), argArray, throwable)
}

func ThrowableCandidate(args []any) error {
	if len(args) == 0 {
		return nil
	}

	lastEntry := args[len(args)-1]

	throwable, ok := lastEntry.(error)
	if ok {
		return throwable
	}

	return nil
}

// @bref Helper method to get all but the last element of an array
//
// @param argArray The arguments from which we want to remove the last element
//
// @return a copy of the array without the last element
func TrimmedCopy(args []any) ([]any, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("%s", "non-sensical empty or null argument array")
	}

	trimmedLen := len(args) - 1
	trimmed := make([]any, 0, trimmedLen)

	if trimmedLen > 0 {
		copy(trimmed, args[:trimmedLen])
	}

	return trimmed, nil
}

////////////////////////////////////////////////////////////////////////

func isEscapedDelimeter(messagePattern string, delimeterStartIndex int) bool {
	if delimeterStartIndex == 0 {
		return false
	}
	potentialEscape := messagePattern[delimeterStartIndex-1]
	return potentialEscape == ESCAPE_CHAR
}

func isDoubleEscaped(messagePattern string, delimeterStartIndex int) bool {
	return delimeterStartIndex >= 2 && messagePattern[delimeterStartIndex-2] == ESCAPE_CHAR
}

func deeplyAppendParameter(sbuf *strings.Builder, o any, seenMap map[any]any) {
	if o == nil {
		sbuf.WriteString("null")
		return
	}

	objType := reflect.TypeOf(o)
	if objType.Kind() == reflect.Slice || objType.Kind() == reflect.Array {
		switch o := o.(type) {
		case []bool:
			booleanArrayAppend(sbuf, o)
		case []byte:
			byteArrayAppend(sbuf, o)
		case []rune:
			charArrayAppend(sbuf, o)
		case []int16:
			shortArrayAppend(sbuf, o)
		case []int:
			intArrayAppend(sbuf, o)
		case []int64:
			longArrayAppend(sbuf, o)
		case []float32:
			floatArrayAppend(sbuf, o)
		case []float64:
			doubleArrayAppend(sbuf, o)
		case []string:
			stringArrayAppend(sbuf, o)
		default:
			objectArrayAppend(sbuf, o.([]any), seenMap)
		}
	} else {
		safeObjectAppend(sbuf, o)
	}

}

func safeObjectAppend(sbuf *strings.Builder, o any) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("SLF4J: Failed toString() invocation on an object of type [%T]\n", o)
			fmt.Println(r)
			sbuf.WriteString("[FAILED toString()]")
		}
	}()

	oAsString := fmt.Sprintf("%v", o)
	sbuf.WriteString(oAsString)
}

func objectArrayAppend(sbuf *strings.Builder, a []any, seenMap map[any]any) {
	sbuf.WriteRune('[')
	key := fmt.Sprintf("%p", a)
	if _, ok := seenMap[key]; !ok {
		seenMap[key] = nil
		len := len(a)
		for i := 0; i < len; i++ {
			deeplyAppendParameter(sbuf, a[i], seenMap)
			if i != len-1 {
				sbuf.WriteString(", ")
			}
		}
		// allow repeats in siblings
		delete(seenMap, key)
	} else {
		sbuf.WriteString("...")
	}
	sbuf.WriteRune(']')
}

func booleanArrayAppend(sbuf *strings.Builder, a []bool) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(fmt.Sprintf("%v", val))
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}

func byteArrayAppend(sbuf *strings.Builder, a []byte) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(fmt.Sprintf("%v", val))
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}

func charArrayAppend(sbuf *strings.Builder, a []rune) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(fmt.Sprintf("%c", val))
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}

func shortArrayAppend(sbuf *strings.Builder, a []int16) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(fmt.Sprintf("%v", val))
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}

func intArrayAppend(sbuf *strings.Builder, a []int) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(fmt.Sprintf("%d", val))
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}

func longArrayAppend(sbuf *strings.Builder, a []int64) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(fmt.Sprintf("%d", val))
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}

func floatArrayAppend(sbuf *strings.Builder, a []float32) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(fmt.Sprintf("%f", val))
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}

func doubleArrayAppend(sbuf *strings.Builder, a []float64) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(fmt.Sprintf("%f", val))
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}

func stringArrayAppend(sbuf *strings.Builder, a []string) {
	sbuf.WriteRune('[')
	for i, val := range a {
		sbuf.WriteString(val)
		if i != len(a)-1 {
			sbuf.WriteString(", ")
		}
	}
	sbuf.WriteRune(']')
}
