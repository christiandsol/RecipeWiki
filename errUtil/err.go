package errUtil

import (
	"fmt"
)

func PrintErr(msg string, data any, err error) {
	if err != nil {
		fmt.Sprintf("[ERROR] %v %v Error: %v\n", msg, data, err)
	}
}
