package aws

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ingmardrewing/gomic/config"
)

func TestGetAwsSession(t *testing.T) {
	sess := getAwsSession()
	sType := fmt.Sprintf("%s", reflect.TypeOf(sess))

	expected := "*session.Session"
	if sType != expected {
		t.Errorf("Expected session to be %s, but got %s", expected, sType)
	}
}

func TestGetThumbnailPaths(t *testing.T) {
	config.ReadDirect("/Users/drewing/Sites/gomic.yaml")
	ap := getAwsPage()
	expectedLocal := "/Users/drewing/Desktop/devabo_de_uploads/comicstrips/thumb_DevAbode_0085.png"
	expectedRemote := "comicstrips/thumb_DevAbode_0085.png"

	local, remote := getThumbnailPaths(ap)

	if local != expectedLocal {
		t.Errorf("Expected %s, but got %s", expectedLocal, local)
	}

	if remote != expectedRemote {
		t.Errorf("Expected %s, but got %s", expectedRemote, remote)
	}
}

func TestGetFilePaths(t *testing.T) {
	config.ReadDirect("/Users/drewing/Sites/gomic.yaml")
	ap := getAwsPage()
	expectedLocal := "/Users/drewing/Desktop/devabo_de_uploads/comicstrips/DevAbode_0085.png"
	expectedRemote := "comicstrips/DevAbode_0085.png"

	local, remote := getFilePaths(ap)

	if local != expectedLocal {
		t.Errorf("Expected %s, but got %s", expectedLocal, local)
	}

	if remote != expectedRemote {
		t.Errorf("Expected %s, but got %s", expectedRemote, remote)
	}
}

// ******  mocking

type pageMock struct{}

func (p pageMock) ImageFilename() string {
	return "DevAbode_0085.png"
}

func getAwsPage() AwsPage {
	return &pageMock{}
}
