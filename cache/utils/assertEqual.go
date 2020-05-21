package utils

import (
	"testing"
	"fmt"
	"reflect"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		t.Logf("TEST PASSED")
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func AssertDeepEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if reflect.DeepEqual(a, b) {
		t.Logf("TEST PASSED")
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}