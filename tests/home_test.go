package tests

import (
	"github.com/roscrl/light/core"
	"testing"

	"github.com/go-rod/rod"
	"github.com/matryer/is"

	"github.com/roscrl/light/config"
)

func TestHome(t *testing.T) {
	_, s := is.New(t), core.NewServer(config.TestConfig())

	s.Start()
	defer s.Stop()

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("http://localhost:" + s.Cfg.Port).MustWaitStable()

	page.MustScreenshot("screenshot.png")
}
