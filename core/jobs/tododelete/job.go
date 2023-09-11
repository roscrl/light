package tododelete

import (
	"fmt"
	"log"
)

const (
	id = "id"
)

func Args(todoID string) map[string]any {
	return map[string]any{
		id: todoID,
	}
}

func Run(args map[string]any) error {
	log.Println(args)

	return fmt.Errorf("something went wrong")
}
