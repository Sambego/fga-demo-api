package errors

import (
	"fmt"
)

var ErrorUnauthorized = fmt.Errorf("unauthorized")
var ErrorForbidden = fmt.Errorf("forbidden")
