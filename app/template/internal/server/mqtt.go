package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"go_kratos_template/app/template/internal/conf"
	"go_kratos_template/pkg/client"
)

func NewMqttClient(c *conf.Server, logger log.Logger) *client.Mqtt {
	mqttConf := client.MQTTConfig{
		Network:  c.Mqtt.Network,
		Addr:     c.Mqtt.Addr,
		ClientID: c.Mqtt.ClientID,
		Username: c.Mqtt.Username,
		Password: c.Mqtt.Password,

		AutoReconnect:        c.Mqtt.AutoReconnect,
		MaxReconnectInterval: c.Mqtt.MaxReconnectInterval.AsDuration(),
	}
	return client.NewMqttClient(&mqttConf, logger)
}
