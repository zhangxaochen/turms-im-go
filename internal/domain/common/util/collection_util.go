package util

import (
	"fmt"
	"reflect"
)

// ContainsAllLooseComparison checks if actual contains expected with type lenience.
func ContainsAllLooseComparison(actual, expected interface{}) bool {
	if expected == nil {
		return true // nil expected means all conditions met
	}
	if actual == nil {
		return false
	}

	expVal := reflect.ValueOf(expected)
	actVal := reflect.ValueOf(actual)

	switch expVal.Kind() {
	case reflect.Map:
		if actVal.Kind() != reflect.Map {
			return false
		}
		for _, k := range expVal.MapKeys() {
			expV := expVal.MapIndex(k)
			actV := actVal.MapIndex(k)
			if !actV.IsValid() {
				// Try falling back to string coercion for keys
				kStr := fmt.Sprint(k.Interface())
				found := false
				for _, ak := range actVal.MapKeys() {
					if fmt.Sprint(ak.Interface()) == kStr {
						actV = actVal.MapIndex(ak)
						found = true
						break
					}
				}
				if !found || !actV.IsValid() {
					return false
				}
			}
			if !ContainsAllLooseComparison(actV.Interface(), expV.Interface()) {
				return false
			}
		}
		return true

	case reflect.Slice, reflect.Array:
		if actVal.Kind() != reflect.Slice && actVal.Kind() != reflect.Array {
			return false
		}
		// Expected elements must all be present in actual
		for i := 0; i < expVal.Len(); i++ {
			expV := expVal.Index(i)
			found := false
			for j := 0; j < actVal.Len(); j++ {
				actV := actVal.Index(j)
				if ContainsAllLooseComparison(actV.Interface(), expV.Interface()) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true

	default:
		// Scalar comparison using formatted string
		return fmt.Sprint(actual) == fmt.Sprint(expected)
	}
}
