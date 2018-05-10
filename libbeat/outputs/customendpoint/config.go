package customendpoint

import (
	"time"
	"fmt"

	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/transport"
	"github.com/elastic/beats/libbeat/outputs/codec"
)

type config struct {
	Host           string                `config:"host"`
	Port           int                   `config:"port"`
	Path		   string				 `config:"path"`
	Token		   string				 `config:"token"`
	MaxRetries     int                   `config:"max_retries"`
	PacketShipper  string                `config:"default_packet_shipper"`
	TLS            *outputs.TLSConfig    `config:"tls"`
	Timeout        time.Duration         `config:"timeout"`
	Proxy          transport.ProxyConfig `config:",inline"`
	Codec          codec.Config 		 `config:"codec"`
	Username       string                `config:"username"`
	Password       string                `config:"password"`
}

var (
	defaultConfig = config{
		Host:   		"localhost",
		Port:           5000,
		Path:			"send-data",
		MaxRetries:     3,
		PacketShipper:  "packetbeat",
	}
)

func (c *config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("hostname %v is invalid",
			c.Host)
	}

	return nil
}
