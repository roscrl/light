package browser

import (
	"testing"

	"github.com/go-rod/rod"
)

const (
	localhost = "http://localhost:"
)

func newBrowserWithCleanup(t *testing.T) *rod.Browser {
	t.Helper()

	browser := rod.New().MustConnect()

	t.Cleanup(func() {
		browser.MustClose()
	})

	return browser
}
