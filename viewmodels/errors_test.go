package viewmodels_test

import (
	"testing"

	"github.com/eduardohoraciosanto/simple-blockchain/viewmodels"
)

func TestErrorInterface(t *testing.T) {
	vErr := viewmodels.Error{
		Code:        "some-code",
		Description: "Some Description",
	}
	if vErr.Error() != "some-code:Some Description" {
		t.Fatalf("Unexpected Stringed Error")
	}
}
