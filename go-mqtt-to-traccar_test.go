package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func Test_ttnDataToUrl_empty(t *testing.T) {
	data := TTNData{EndDeviceIds{}, UplinkMessage{}, ""}

	result := ttnDataToUrl(data)
	assert.Equal(t, "http://10.0.0.10:3055/?id=&lat=0.000000&lon=0.000000&hdop=0.000000&batt=10&md=", result)
}

func Test_ttnDataToUrl_data(t *testing.T) {
	decoded := DecodedPayload{1.0, 2.0, 3.0, 4.0, 5.0, "FALSE", "A", nil}
	uplink := UplinkMessage{decoded}
	data := TTNData{EndDeviceIds{}, uplink, ""}

	result := ttnDataToUrl(data)
	assert.Equal(t, "http://10.0.0.10:3055/?id=&lat=1.000000&lon=2.000000&hdop=3.000000&batt=70&md=A", result)
}

func Test_ttnDataToUrl_t1000(t *testing.T) {
	messages := [][]Message{
		{Message{"4198", []byte("10"), "Lat"}, Message{"4197", []byte("20"), "Long"}},
	}

	decoded := DecodedPayload{0.0, 0.0, 0.0, 0.0, 0.0, "FALSE", "A", messages}
	uplink := UplinkMessage{decoded}
	data := TTNData{EndDeviceIds{"device"}, uplink, ""}

	result := ttnDataToUrl(data)

	assert.True(t, strings.HasPrefix(result, "http://10.0.0.10:3055/?id=device&lat=10&lon=20&batt=&timestamp="))
}

func Test_parse(t *testing.T) {
	data := `{"end_device_ids":{"device_id":"device","application_ids":{"application_id":"void"},"dev_eui":"XXX","join_eui":"XXX","dev_addr":"XXX"},"correlation_ids":["gs:uplink:XXX"],"received_at":"XXX","uplink_message":{"session_key_id":"XXX","f_port":0,"f_cnt":0,"frm_payload":"XXX","decoded_payload":{"err":0,"messages":[[{"measurementId":"4200","measurementValue":"0","type":"SOS Event"},{"measurementId":"4197","measurementValue":4.444444,"type":"Longitude"},{"measurementId":"4198","measurementValue":"3.333333","type":"Latitude"},{"measurementId":"3000","measurementValue":98,"type":"Battery"},{"measurementValue":1698957833000,"type":"Timestamp"}]],"payload":"XXX","valid":true},"rx_metadata":[{"gateway_ids":{"gateway_id":"XXX","eui":"XXX"},"time":"XXXZ","timestamp":3493860637,"rssi":0,"channel_rssi":0,"snr":0,"uplink_token":"XXX","channel_index":0,"gps_time":"XXXZ","received_at":"XXXZ"},{"gateway_ids":{"gateway_id":"XXX","eui":"XXX"},"time":"XXXZ","timestamp":1379597011,"rssi":0,"channel_rssi":0,"snr":0,"location":{"latitude":0,"longitude":0,"source":"SOURCE_REGISTRY"},"uplink_token":"XXX","channel_index":0,"gps_time":"XXX","received_at":"XXX"}],"settings":{"data_rate":{"lora":{"bandwidth":0,"spreading_factor":0,"coding_rate":"4/5"}},"frequency":"0","timestamp":0,"time":"XXXZ"},"received_at":"XXXZ","confirmed":true,"consumed_airtime":"0.0s","version_ids":{"brand_id":"XXX","model_id":"XXX","hardware_version":"1.0","firmware_version":"1.0","band_id":"XXX"},"network_ids":{"net_id":"000000","ns_id":"XXX","tenant_id":"xxx","cluster_id":"xxx","cluster_address":"xxx"}}}`
	var ttnData TTNData
	json.Unmarshal([]byte(data), &ttnData)

	result := ttnDataToUrl(ttnData)

	assert.True(t, strings.HasPrefix(result, "http://10.0.0.10:3055/?id=device&lat=3.333333&lon=4.444444&batt=98&timestamp="))
}

func Test_transformDataToGoogleLocation(t *testing.T) {
	data := `[{"mac":"MAC-1","rssi":"-12"},{"mac":"MAC-2","rssi":"-13"}]`
	result := transformDataToGoogleLocation(data)

	geoLocationWifi := []GeoLocationWifi{GeoLocationWifi{"MAC-1", "-12", "0"}, GeoLocationWifi{"MAC-2", "-13", "0"}}

	assert.Equal(t, GeoLocationRequest{geoLocationWifi}, result)
}
