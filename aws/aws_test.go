package aws

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetAwsSession(t *testing.T) {
	sess := getAwsSession()
	sType := fmt.Sprintf("%s", reflect.TypeOf(sess))

	expected := "*session.Session"
	if sType != expected {
		t.Errorf("Expected session to be %s, but got %s", expected, sType)
	}
}
