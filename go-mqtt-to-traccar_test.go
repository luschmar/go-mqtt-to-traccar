package main

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestToParameterList(t *testing.T) {
	data := TTNData{EndDeviceIds{}, UplinkMessage{}, ""}

	ttnDataToUrl(data)
}
