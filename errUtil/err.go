package errUtil

import (
	"fmt"
)

func CheckErr(msg string, data any, err error) {
	if err != nil {
		fmt.Sprintf("[ERROR] %v %v Error: %v\n", msg, data, err)
	}
}
