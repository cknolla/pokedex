package test

import "fmt"

func WriteStringDiff(expectedString, receivedString string) string {
	return fmt.Sprintf(
		"\nExpected:\t'%s'\nGot:\t\t'%s'",
		expectedString,
		receivedString,
	)
}
