package lib

import "fmt"

func Info(msgStr string) {
	fmt.Printf("\nGot: %v\n", msgStr)
}
func Err(errStr string, err error) {
	fmt.Printf("\nError: %v :\n %v\n", errStr, err)
}
