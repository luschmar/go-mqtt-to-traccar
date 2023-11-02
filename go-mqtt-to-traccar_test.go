package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func Test_ttnDataToUrl_empty(t *testing.T) {
	data := TTNData{EndDeviceIds{}, UplinkMessage{}, ""}

	result := ttnDataToUrl(data, "aa/a/aa/aa")
	assert.Equal(t, "http://10.0.0.10:3055/?id=&lat=0.000000&lon=0.000000&hdop=0.000000&batt=10&md=", result)
}

func Test_ttnDataToUrl_data(t *testing.T) {
	decoded := DecodedPayload{1.0, 2.0, 3.0, 4.0, 5.0, "FALSE", "A", nil}
	uplink := UplinkMessage{decoded}
	data := TTNData{EndDeviceIds{}, uplink, ""}

	result := ttnDataToUrl(data, "a/b/c/d")
	assert.Equal(t, "http://10.0.0.10:3055/?id=&lat=1.000000&lon=2.000000&hdop=3.000000&batt=70&md=A", result)
}

func Test_ttnDataToUrl_t1000(t *testing.T) {
	messages := make([]Message, 2)
	messages[0] = Message{"4198", []byte("10"), "Lat"}
	messages[1] = Message{"4197", []byte("20"), "Long"}

	decoded := DecodedPayload{0.0, 0.0, 0.0, 0.0, 0.0, "FALSE", "A", messages}
	uplink := UplinkMessage{decoded}
	data := TTNData{EndDeviceIds{}, uplink, ""}

	result := ttnDataToUrl(data, "a/b/c/d")
	assert.Equal(t, "http://10.0.0.10:3055/?id=c&lat=10&lon=20&batt=&timestamp=", result)
}
