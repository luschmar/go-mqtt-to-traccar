# go-mqtt-to-traccar
Export values from MQTT to Traccar

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