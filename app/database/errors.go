package database

import (
	"fmt"
	"strings"
)

type ErrInvalidDataStoreName []string

func (ds ErrInvalidDataStoreName) Error() error {
	return fmt.Errorf("datastore: invalid datastore name. Must be one of: %s", strings.Join(ds, ", "))
}
