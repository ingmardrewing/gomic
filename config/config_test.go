package config

import "testing"

func TestIsDev(t *testing.T) {
	Stage = "dev"
	expected := true
	actual := IsDev()
	if actual != expected {
		t.Errorf("Expected %t, but got %t", expected, actual)
	}
}

func TestIsProd(t *testing.T) {
	Stage = "prod"
	expected := true
	actual := IsProd()
	if actual != expected {
		t.Errorf("Expected %t, but got %t", expected, actual)
	}
}

func TestIsTest(t *testing.T) {
	Stage = "test"
	expected := true
	actual := IsTest()
	if actual != expected {
		t.Errorf("Expected %t, but got %t", expected, actual)
	}
}
