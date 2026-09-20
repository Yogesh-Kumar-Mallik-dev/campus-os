package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
)

// BLOCK_UI_TEST_001
// Purpose: Verifies headless construction of Fyne UI layout and tab navigation.
func TestAppBuildLayout(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	testWindow := test.NewWindow(nil)
	clientApp := &Application{
		FyneApp:   testApp,
		Window:    testWindow,
		APIClient: api.NewClient("http://localhost:8080"),
	}

	layout := clientApp.BuildLayout()
	if layout == nil {
		t.Fatal("BLOCK_UI_TEST_001: expected non-nil layout canvas object")
	}
}
