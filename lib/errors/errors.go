package errors

import "fmt"

func PrintIfError(err error, msg string) {
	if err != nil {
		fmt.Println(msg + ": " + err.Error())
	}
}
