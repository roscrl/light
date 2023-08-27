package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-rod/rod/lib/input"
	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/app"
	"github.com/roscrl/light/core/db"
	"github.com/roscrl/light/core/support/ulid"
)

func TestHome(t *testing.T) {
	t.Parallel()

	cfg := config.NewTestConfig()
	cfg.SqliteDBPath = "file::memory:?cache=shared"

	is, app := app.NewStartedTestAppWithCleanup(t, cfg)

	ulid1, ulid2 := ulid.NewString(), ulid.NewString()
	{ // migrate and seed in memory database
		db.RunMigrations(app.DB, db.PathMigrations)

		_, err := app.DB.Exec("INSERT INTO todos (id, task, status) VALUES (?, ?, ?)", ulid1, "important todo!", "pending")
		is.NoErr(err)

		_, err = app.DB.Exec("INSERT INTO todos (id, task, status) VALUES (?, ?, ?)", ulid2, "also important!", "done")
		is.NoErr(err)
	}

	browser := newBrowserWithCleanup(t)

	{ // has correct not found page message
		notFoundPage := browser.MustPage("http://localhost:" + app.Cfg.Port + "/does-not-exist").MustWaitStable()
		notFoundMessage := notFoundPage.MustElement("[data-testid='notfound_message']").MustText()

		is.Equal(strings.TrimSpace(notFoundMessage), "Oops, that page does not exist.")
	}

	homePage := browser.MustPage("http://localhost:" + app.Cfg.Port).MustWaitStable()

	{ // has correct home page title
		homePageTitle := homePage.MustElement("[data-testid='page_title']").MustText()
		is.Equal(homePageTitle, "Todos!")
	}

	todoItems := homePage.MustElements("#todo_list > turbo-frame")

	{ // check there is a todo list with two items from seed data
		is.Equal(len(todoItems), 2)
	}

	{ // check todo items have correct text
		todoItem1Text := todoItems[0].MustElement("form > turbo-frame > label > span").MustText()
		is.Equal(strings.TrimSpace(todoItem1Text), "important todo!")

		todoItem2Text := todoItems[1].MustElement("form > turbo-frame > label > span").MustText()
		is.Equal(strings.TrimSpace(todoItem2Text), "also important!")
	}

	{ // check todo items have correct status
		todoItem1InputElement := todoItems[0].MustElement(fmt.Sprintf("#todo-%v > label > input[type=checkbox]", ulid1))
		todoItem1Checked := todoItem1InputElement.MustAttribute("checked")
		is.Equal(todoItem1Checked, nil) // checked attribute is not present

		todoItem2InputElement := todoItems[1].MustElement(fmt.Sprintf("#todo-%v > label > input[type=checkbox]", ulid2))
		todoItem2Checked := todoItem2InputElement.MustAttribute("checked")

		is.Equal(*todoItem2Checked, "") // checked attribute is present but empty value which means it is checked
	}

	{ // check adding empty todo shows form error message
		homePage.MustElement("#todo_form_frame > form > div > label > input").MustInput("").MustType(input.Enter)
		homePage.MustWaitStable()

		todoFormError := homePage.MustElement("[data-testid='todo_form_error']").MustText()
		is.Equal(strings.TrimSpace(todoFormError), "task cannot be empty.")
	}

	{ // check adding todo prepends to todo list
		homePage.MustElement("#todo_form_frame > form > div > label > input").MustInput("my new todo!").MustType(input.Enter)
		homePage.MustWaitStable()

		newTodoItems := homePage.MustElements("#todo_list > turbo-frame")
		totalTodoItems := len(newTodoItems)
		is.Equal(totalTodoItems, 3)

		text := newTodoItems[0].MustElement("form > turbo-frame > label > span").MustText()

		is.Equal(strings.TrimSpace(text), "my new todo!")
	}
}
