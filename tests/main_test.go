package tests

import (
	"testing"

	"go.uber.org/goleak"

	"github.com/go-rod/rod"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func newBrowserWithCleanup(t *testing.T) *rod.Browser {
	t.Helper()

	browser := rod.New().MustConnect()

	t.Cleanup(func() {
		browser.MustClose()
	})

	return browser
}
