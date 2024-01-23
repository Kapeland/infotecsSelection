package uuid

import (
	"fmt"
	guuid "github.com/google/uuid"
)

func CreateUUID() string {
	id := guuid.New()
	return fmt.Sprintf("%s", id.String())
}
