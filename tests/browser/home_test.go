package browser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"

	"github.com/roscrl/light/app"
	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/helpers/ulid"
	"github.com/roscrl/light/db"

	_ "github.com/roscrl/light/core/utils/testutil"
)

func TestHome(t *testing.T) {
	t.Parallel()

	cfg := config.NewTestConfig()
	cfg.SqliteDBPath = fmt.Sprintf("file:%s?mode=memory&cache=shared", ulid.NewString())

	is, app := app.NewStartedTestAppWithCleanup(t, cfg)

	todoID1, todoID2 := ulid.NewString(), ulid.NewString()
	{ // migrate and seed in memory database
		db.RunMigrations(app.DB, db.PathMigrations)

		_, err := app.DB.Exec("INSERT INTO todos (id, task, status) VALUES (?, ?, ?)", todoID1, "important todo!", "pending")
		is.NoErr(err)

		_, err = app.DB.Exec("INSERT INTO todos (id, task, status) VALUES (?, ?, ?)", todoID2, "also important!", "done")
		is.NoErr(err)
	}

	browser := newBrowserWithCleanup(t)

	{ // has correct not found page message
		notFoundPage := browser.MustPage(localhost + app.Cfg.Port + "/does-not-exist")
		notFoundMessage := notFoundPage.MustElement("[data-testid='notfound_message']").MustText()

		is.Equal(strings.TrimSpace(notFoundMessage), "Oops, that page does not exist.")

		notFoundPage.MustClose()
	}

	homePage := browser.MustPage(localhost + app.Cfg.Port)

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
		todoItem1InputElement := todoItems[0].MustElement(fmt.Sprintf("#todo-%v > label > input[type=checkbox]", todoID1))
		todoItem1Checked := todoItem1InputElement.MustAttribute("checked")
		is.Equal(todoItem1Checked, nil) // checked attribute is not present

		todoItem2InputElement := todoItems[1].MustElement(fmt.Sprintf("#todo-%v > label > input[type=checkbox]", todoID2))
		todoItem2Checked := todoItem2InputElement.MustAttribute("checked")

		is.Equal(*todoItem2Checked, "") // checked attribute is present but empty value which means it is checked
	}

	{ // check adding empty todo shows form error message
		homePage.MustElement("#todo_form_frame > form > div > label > input").MustInput("").MustType(input.Enter)

		todoFormError := homePage.MustElement("[data-testid='todo_form_error']").MustText()
		is.Equal(strings.TrimSpace(todoFormError), "task cannot be empty.")
	}

	var addedTodo *rod.Element

	{ // check adding todo prepends to todo list
		homePage.MustElement("#todo_form_frame > form > div > label > input").MustInput("my new todo!").MustType(input.Enter)

		homePage.MustWaitElementsMoreThan("#todo_list > turbo-frame", 2)
		addedTodo = homePage.MustElements("#todo_list > turbo-frame").First()

		text := addedTodo.MustElement("form > turbo-frame > label > span").MustText()

		is.Equal(strings.TrimSpace(text), "my new todo!")
	}

	{ // check updating added todo status and task name works
		addedTodo.MustElement("form > turbo-frame > label > input[type=checkbox]").MustClick() // check the checkbox

		homePage.MustWaitIdle()

		outerID := addedTodo.MustAttribute("id")
		addedTodoID := strings.Split(*outerID, "-")[1]

		todoEditButtonSelector := fmt.Sprintf("#todo-%v > div", addedTodoID)
		homePage.MustElement(todoEditButtonSelector).MustClick()

		homePage.MustWaitIdle()

		todoInputSelector := fmt.Sprintf("#todo-%v > label > label > input[type=text]", addedTodoID)

		homePage.MustElement(todoInputSelector).MustSelectAllText().MustInput("").MustInput("updated todo!").MustType(input.Enter)

		homePage.MustWaitIdle()

		updatedTextSelector := fmt.Sprintf("#todo-%v > label > span", addedTodoID)
		updatedText := homePage.MustElement(updatedTextSelector).MustText()
		is.Equal(strings.TrimSpace(updatedText), "updated todo!")
	}
}

func TestHomeSearch(t *testing.T) {
	t.Parallel()

	cfg := config.NewTestConfig()
	cfg.SqliteDBPath = fmt.Sprintf("file:%s?mode=memory&cache=shared", ulid.NewString())

	is, app := app.NewStartedTestAppWithCleanup(t, cfg)

	todoID1, todoID2 := ulid.NewString(), ulid.NewString()
	{ // migrate and seed in memory database
		db.RunMigrations(app.DB, db.PathMigrations)

		_, err := app.DB.Exec("INSERT INTO todos (id, task, status) VALUES (?, ?, ?)", todoID1, "important todo!", "pending")
		is.NoErr(err)

		_, err = app.DB.Exec("INSERT INTO todos (id, task, status) VALUES (?, ?, ?)", todoID2, "also important!", "done")
		is.NoErr(err)
	}

	browser := newBrowserWithCleanup(t)
	homePage := browser.MustPage(localhost + app.Cfg.Port)

	{ // searching for only one todo item shows only one todo item
		formSearchInput := homePage.MustElement("[data-testid='todo_form_search_input']")

		formSearchInput.MustInput("also").MustType(input.Enter)

		homePage.MustWaitIdle()

		todoItems := homePage.MustElements("#todo_list > turbo-frame")
		is.Equal(len(todoItems), 1)

		todoItemText := todoItems[0].MustElement("form > turbo-frame > label > span").MustText()
		is.Equal(strings.TrimSpace(todoItemText), "also important!")
	}
}
