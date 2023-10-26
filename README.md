# go-mqtt-to-traccar
Transfer data from a LGT-92 Device to TTN/MQTT into Traccar.

## Setup

0) LGT-92 Setup
1) TTN Setup
2) Local MQTT (Mosquitto) Setup
3) Run Docker local
4) Enjoy

# LGT-92 Setup

## Update Firmware to v1.6.7 

Download Firmware from: https://www.dragino.com/downloads/downloads/LGT_92/firmware/v1.6.7/EU868.hex

## Connect with ST-LINK/V2

```
LGT-92 black -> ST-LINK/V2 STM32 Pin 20 (GND)
LGT-92 white -> ST-LINK/V2 STM32 Pin 7 (TMS_SWDIO)
LGT-92 green -> ST-LINK/V2 STM32 Pin 9 (TCK_SWCLK)
LGT-92 red -> ST-LINK/V2 STM32 Pin 19 (VDD)
```

Run following command:

```openocd -f /usr/share/openocd/scripts/interface/stlink-v2.cfg -c "transport select hla_swd" -f /usr/share/openocd/scripts/target/stm32l0.cfg -c init -c 'program EU868.hex verify reset exit'```

https://tobru.ch/firmware-update-on-dragino-lgt-92-with-st-link-v2/

## Serial Console

Wire following:
```
LGT-92 black -> UartSBee GND
LGT-92 green -> UartSBee RX
LGT-92 white -> UartSBee TX
```

Execute follwing commands:
```
AT+VER
AT+FDR
AT+SGM=0
```

AT+VER: Image Version and Frequency Band
AT+FDR: Factory Data Reset
AT+SGM=0: enable the motion sensor to roll/pitch info from accelerometer / and HDOP/ Altitude -> Because payload is larger -> 18 bytes this won't work on DR0 on US915/AU915

https://www.dragino.com/downloads/downloads/LGT_92/LGT-92_LoRa_GPS_Tracker_UserManual_v1.6.8.pdf

# TTN Setup

https://www.dragino.com/downloads/downloads/LGT_92/Decoder/LGT92-v1.6.4_decoder_TTN.txt

```
Javascript Decoder
```

## Traccar Setup

Filter out zero values; insert following configs in `traccar.xml`
```
    <entry key='filter.enable'>true</entry>
    <entry key='filter.zero'>true</entry>
```


## Run Docker

MQTT_HOST
MQTT_TOPIC
TC_HOST


## Example MQTT Message

```
{
    "end_device_ids":
    {
        "device_id": "???",
        "application_ids":
        {
            "application_id": "???"
        },
        "dev_eui": "???",
        "join_eui": "???",
        "dev_addr": "???"
    },
    "correlation_ids":
    [
        "???"
    ],
    "received_at": "???Z",
    "uplink_message":
    {
        "session_key_id": "???",
        "f_port": ???,
        "f_cnt": ???,
        "frm_payload": "???",
        "decoded_payload":
        {
            "ALARM_status": "???",
            "Altitude": ???,
            "BatV": ???,
            "FW": ???,
            "HDOP": ???,
            "LON": "???",
            "Latitude": ???,
            "Longitude": ???,
            "MD": "???",
            "Pitch": ???,
            "Roll": ???
        },
        "rx_metadata":
        [
            {
                "gateway_ids":
                {
                    "gateway_id": "???",
                    "eui": "???"
                },
                "time": "???Z",
                "timestamp": ???,
                "rssi": ???,
                "channel_rssi": ???,
                "snr": ???,
                "location":
                {
                    "latitude": ???,
                    "longitude": ???,
                    "altitude": ???,
                    "source": "SOURCE_REGISTRY"
                },
                "uplink_token": "???",
                "received_at": "???Z"
            },
            {
                "gateway_ids":
                {
                    "gateway_id": "???",
                    "eui": "???"
                },
                "time": "???Z",
                "timestamp": ???,
                "rssi": ???,
                "channel_rssi": ???,
                "snr": ???,
                "location":
                {
                    "latitude": ???,
                    "longitude": ???,
                    "altitude": ???,
                    "source": "SOURCE_REGISTRY"
                },
                "uplink_token": "???",
                "received_at": "???Z"
            }
        ],
        "settings":
        {
            "data_rate":
            {
                "lora":
                {
                    "bandwidth": ???,
                    "spreading_factor": ???,
                    "coding_rate": "???"
                }
            },
            "frequency": "???",
            "timestamp": ???,
            "time": "???Z"
        },
        "received_at": "???Z",
        "consumed_airtime": "?.?????s",
        "locations":
        {
            "frm-payload":
            {
                "latitude": ?.?????,
                "longitude": ?.?????,
                "source": "SOURCE_GPS"
            }
        },
        "network_ids":
        {
            "net_id": "???",
            "tenant_id": "???",
            "cluster_id": "???",
            "cluster_address": "???.cloud.thethings.network"
        }
    }
}
```