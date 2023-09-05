package browser

import (
	"strings"
	"testing"

	"github.com/roscrl/light/app"
	"github.com/roscrl/light/config"

	_ "github.com/roscrl/light/core/utils/testutil"
)

func TestNotFound(t *testing.T) {
	t.Parallel()

	is, app := app.NewStartedTestAppWithCleanup(t, config.NewTestConfig())

	browser := newBrowserWithCleanup(t)

	notFoundPage := browser.MustPage(localhost + app.Cfg.Port + "/does-not-exist")
	notFoundMessage := notFoundPage.MustElement("[data-testid='notfound_message']").MustText()

	is.Equal(strings.TrimSpace(notFoundMessage), "Oops, that page does not exist.")

	notFoundPage.MustClose()
}
