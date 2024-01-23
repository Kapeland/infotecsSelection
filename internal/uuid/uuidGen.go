package uuid

import (
	"fmt"
	guuid "github.com/google/uuid"
)

func CreateUUID() string {
	id := guuid.New()
	return fmt.Sprintf("%s", id.String())
}

// Checks the validity of the URL.
func CheckUUID(varUUID string) error {
	return guuid.Validate(varUUID)
}
