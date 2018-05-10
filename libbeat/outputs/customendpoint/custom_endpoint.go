package customendpoint

import (
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/codec"
	"time"
	"net/http"
)

var debugf = logp.MakeDebug("customendpoint")

const (
	defaultWaitRetry    = 1 * time.Second
	defaultMaxWaitRetry = 60 * time.Second
)

func init() {
	debugf("Registering customendpoint....")
	outputs.RegisterType("customendpoint", makeCustomEndpointOut)
}

// makeCustomEndpointOut instantiates a new CustomEndpoint output instance.
func makeCustomEndpointOut(
	beat beat.Info,
	observer outputs.Observer,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return outputs.Fail(err)
	}

	// disable bulk support in publisher pipeline
	cfg.SetInt("bulk_max_size", -1, -1)

	//tls, err := outputs.LoadTLSConfig(config.TLS)
	//if err != nil {
	//	return outputs.Fail(err)
	//}

	//transp := &transport.Config{
	//	Timeout: config.Timeout,
	//	Proxy:   &config.Proxy,
	//	TLS:     tls,
	//	Stats:   observer,
	//}

	enc, err := codec.CreateEncoder(beat, config.Codec)
	if err != nil {
		return outputs.Fail(err)
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    config.Timeout * time.Second,
		//DisableCompression: true,
	}
	hc := &http.Client{Transport: tr}
	if err != nil {
		return outputs.Fail(err)
	}

	client := newClient(hc, beat, observer, config.Timeout, config.Host, config.Port, config.Path,
		config.Token, config.Username, config.Password, enc)

	//return outputs.Success(-1, 0, fo)
	return outputs.Success(-1, config.MaxRetries, client)
}
