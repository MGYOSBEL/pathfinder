package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Options struct {
	Server string
	Topic  string
	QoS    byte
}

type MqttClient struct {
	options Options
	client  MQTT.Client
}

type Handler func(payload []byte)

func NewClient(o Options) *MqttClient {
	opts := MQTT.NewClientOptions().AddBroker(o.Server)
	c := MQTT.NewClient(opts)
	return &MqttClient{
		client:  c,
		options: o,
	}
}

func (c *MqttClient) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MqttClient) Subscribe(h Handler) error {
	token := c.client.Subscribe(c.options.Topic, c.options.QoS, func(cl MQTT.Client, m MQTT.Message) {
		h(m.Payload())
	})
	token.Wait()
	return token.Error()
}

func (c *MqttClient) Publish(topic string, m any) error {
	token := c.client.Publish(topic, c.options.QoS, false, m)
	token.Wait()
	return token.Error()
}

func (c *MqttClient) Disconnect() {
	c.client.Disconnect(100)
}
