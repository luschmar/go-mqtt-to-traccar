package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"syscall"
	// "time"
	"encoding/json"
	"net/http"
)

type DecodedPayload struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	HDOP      string  `json:"HDOP"`
	Altitude  float64 `json:"Altitude"`
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

	// send https://www.traccar.org/osmand/
	fmt.Println(ttnData.EndDeviceIds.DeviceId)
	fmt.Println(ttnData.UplinkMessage.DecodedPayload.Latitude)
	// parseTime, err := time.Parse(time.RFC1123, ttnData.ReceivedAt)
	// http://10.0.0.10:3055/?id=%s&lat=%f&lon=%f&timestamp=%d&hdop=%d&altitude=%d&speed=%d
	getUrl := fmt.Sprintf("http://localhost:3055/?id=%s&lat=%f&lon=%f",
		ttnData.EndDeviceIds.DeviceId,
		ttnData.UplinkMessage.DecodedPayload.Latitude,
		ttnData.UplinkMessage.DecodedPayload.Longitude) //,
	//parseTime.Unix(),
	//0,//ttnData.UplinkMessage.DecodedPayload.HDOP,
	//0,//ttnData.UplinkMessage.DecodedPayload.Altitude,
	//0)

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
}

func main() {
	var broker = "localhost"
	var port = 3983
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
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

	client.Unsubscribe("ttn/devices/+/up")
	client.Disconnect(250)
}

func sub(client mqtt.Client) {
	topic := "ttn/devices/+/up"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
