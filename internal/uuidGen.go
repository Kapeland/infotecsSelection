package internal

import (
	"fmt"
	guuid "github.com/google/uuid"
)

func GetUUID() string {
	id := guuid.New()
	return fmt.Sprintf("%s", id.String())
}
