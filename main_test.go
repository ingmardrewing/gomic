package main

import "testing"

/*
func TestIsNewFile(t *testing.T) {
	result := isNewFile("DevAbode_0001.png")
	if result {
		t.Error("Expected result to be false, but it is true")
	}

	result = isNewFile(".DS_Store")
	if result {
		t.Error("Expected result to be false, but it is true")
	}

	result = isNewFile("DevAbode_0011.png")
	if result != true {
		t.Error("Expected result to be true, but it is false")
	}
}
*/

func TestIsRelevant(t *testing.T) {
	result := isRelevant(".DS_Store")
	if result {
		t.Error("Expected result to be false, but it is true")
	}

	result = isRelevant("thumb_DevAbode_0085.png")
	if result {
		t.Error("Expected result to be false, but it is true")
	}

	result = isRelevant("DevAbode_0085.png")
	if !result {
		t.Error("Expected result to be true, but it is false")
	}
}
