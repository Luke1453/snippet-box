package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// indicates that its a helper function for testing
	// used to print correct method name in test results
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}
