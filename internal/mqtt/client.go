package mqtt

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	client MQTT.Client
}

func NewClient() *MqttClient {
	opts := MQTT.NewClientOptions().AddBroker("192.168.0.16:1883")
	c := MQTT.NewClient(opts)
	return &MqttClient{
		client: c,
	}
}

func (c *MqttClient) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MqttClient) Subscribe() error {
	c.client.Subscribe("#", 0, func(c MQTT.Client, m MQTT.Message) {
		fmt.Printf("Received message in %s: %s\n", m.Topic(), m.Payload())
	})
	return nil
}
