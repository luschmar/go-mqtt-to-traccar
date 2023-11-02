package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Message struct {
	MeasurementId    string          `json:"measurementId"`
	MeasurementValue json.RawMessage `json:"measurementValue"`
	Type             string          `json:"type"`
}

type DecodedPayload struct {
	Latitude  float64     `json:"Latitude"`
	Longitude float64     `json:"Longitude"`
	HDOP      float32     `json:"HDOP"`
	BatV      float64     `json:"BatV"`
	Altitude  float64     `json:"Altitude"`
	Alarm     string      `json:"ALARM_status"`
	MD        string      `json:"MD"`
	Messages  [][]Message `json:"messages"`
}

type UplinkMessage struct {
	DecodedPayload DecodedPayload `json:"decoded_payload"`
}

type EndDeviceIds struct {
	DeviceId string `json:"device_id"`
}

type TTNData struct {
	EndDeviceIds  EndDeviceIds  `json:"end_device_ids"`
	UplinkMessage UplinkMessage `json:"uplink_message"`
	ReceivedAt    string        `json:"received_at"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var ttnData TTNData
	json.Unmarshal([]byte(msg.Payload()), &ttnData)
	getUrl := ttnDataToUrl(ttnData)

	fmt.Println(getUrl)

	resp, err := http.Get(getUrl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
	panic(err)
}

func main() {
	broker := getenv("MQTT_HOST", "10.0.0.10:3983")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", broker))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)

	// Wait for a signal to exit the program gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Unsubscribe(getTopic())
	client.Disconnect(250)
}

func sub(client mqtt.Client) {
	topic := getTopic()
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}

func ttnDataToUrl(ttnData TTNData) string {
	if len(ttnData.UplinkMessage.DecodedPayload.Messages) > 0 {
		fmt.Printf("Try extract Messages value: %s\n", ttnData.UplinkMessage.DecodedPayload.Messages)
		var deviceId = ttnData.EndDeviceIds.DeviceId
		latitude := ""
		longitude := ""
		batt := ""
		timestamp := ""
		alarm := ""
		var message = ttnData.UplinkMessage.DecodedPayload.Messages[0]
		for i := range message {
			// Lat
			if message[i].MeasurementId == "4198" {
				latitude = strings.Replace(string(message[i].MeasurementValue), "\"", "", -1)
			}
			// Long
			if message[i].MeasurementId == "4197" {
				longitude = strings.Replace(string(message[i].MeasurementValue), "\"", "", -1)
			}
			// Bat
			if message[i].MeasurementId == "3000" {
				batt = strings.Replace(string(message[i].MeasurementValue), "\"", "", -1)
			}
			// Bat
			if message[i].Type == "Timestamp" {
				timestamp = strings.Replace(string(message[i].MeasurementValue), "\"", "", -1)
			}
			// Bat
			if message[i].Type == "4200" {
				var evalAlarm = strings.Replace(string(message[i].MeasurementValue), "\"", "", -1)
				if evalAlarm == "1" {
					alarm = "&alarm=general"
				}
			}
		}
		if longitude != "" && latitude != "" {
			return fmt.Sprintf("http://%s/?id=%s&lat=%s&lon=%s&batt=%s&timestamp=%s%s",
				getenv("TC_HOST", "10.0.0.10:3055"),
				deviceId,
				latitude,
				longitude,
				batt,
				timestamp,
				alarm)
		}
	}

	if strings.EqualFold(ttnData.UplinkMessage.DecodedPayload.Alarm, "true") {
		return fmt.Sprintf("http://%s/?id=%s&lat=%f&lon=%f&hdop=%f&batt=%d&alarm=general&md=%s",
			getenv("TC_HOST", "10.0.0.10:3055"),
			ttnData.EndDeviceIds.DeviceId,
			ttnData.UplinkMessage.DecodedPayload.Latitude,
			ttnData.UplinkMessage.DecodedPayload.Longitude,
			ttnData.UplinkMessage.DecodedPayload.HDOP,
			battVToLevel(ttnData.UplinkMessage.DecodedPayload.BatV),
			ttnData.UplinkMessage.DecodedPayload.MD)
	}
	return fmt.Sprintf("http://%s/?id=%s&lat=%f&lon=%f&hdop=%f&batt=%d&md=%s",
		getenv("TC_HOST", "10.0.0.10:3055"),
		ttnData.EndDeviceIds.DeviceId,
		ttnData.UplinkMessage.DecodedPayload.Latitude,
		ttnData.UplinkMessage.DecodedPayload.Longitude,
		ttnData.UplinkMessage.DecodedPayload.HDOP,
		battVToLevel(ttnData.UplinkMessage.DecodedPayload.BatV),
		ttnData.UplinkMessage.DecodedPayload.MD)
}

func battVToLevel(battV float64) int16 {
	// > 4.0v : 80% ~ 100%
	// 3.85v ~3.99v: 60% ~ 80%
	// 3.70v ~ 3.84v: 40% ~ 60%
	// 3.40v ~ 3.69v: 20% ~ 40%
	// < 3.39v: 0~20%
	if battV > 4.0 {
		return 100
	}

	if battV > 3.85 {
		return 70
	}

	if battV > 3.7 {
		return 50
	}

	if battV > 3.4 {
		return 30
	}
	return 10
}
func getTopic() string {
	return getenv("MQTT_TOPIC", "ttn/devices/+/up")
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
