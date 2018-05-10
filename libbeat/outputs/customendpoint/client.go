package customendpoint

import (
	"time"

	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/codec"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/beat"
	"strconv"
	"net/http"
	"net/url"
)

type client struct {
	*http.Client
	beat     beat.Info
	url	     string
	observer outputs.Observer
	username string
	password string
	token	 string
	codec    codec.Codec
	timeout  time.Duration
}

func newClient(
	hc 			*http.Client,
	beat     	beat.Info,
	observer 	outputs.Observer,
	timeout 	time.Duration,
	host		string,
	port		int,
	path		string,
	token	 	string,
	user 		string,
	pass 		string,
	codec 		codec.Codec,
) *client {
	return &client{
		Client:   hc,
		beat:     beat,
		observer: observer,
		url:      host+":"+strconv.Itoa(port)+path,
		token: 	  token,
		timeout:  timeout,
		username: user,
		password: pass,
		codec:    codec,
	}
}

func (c *client) Close() error { return nil }

func (c *client) Publish(batch publisher.Batch) error {
	defer batch.ACK()

	st := c.observer
	events := batch.Events()
	st.NewBatch(len(events))

	dropped := 0
	for i := range events {
		event := &events[i]

		serializedEvent, err := c.codec.Encode(c.beat.Beat, &event.Content)
		if err != nil {
			if event.Guaranteed() {
				logp.Critical("Failed to serialize the event: %v", err)
			} else {
				logp.Warn("Failed to serialize the event: %v", err)
			}

			dropped++
			continue
		}

		if _, err = c.PostForm(c.url, url.Values{
				"token": {c.token},
				"doc_type": {"network_packet"},
				"doc": {string(serializedEvent[:])}});
				err != nil {
			st.WriteError(err)

			if event.Guaranteed() {
				logp.Critical("Sending event to api failed with: %v", err)
			} else {
				logp.Warn("Sending event to api failed with: %v", err)
			}

			dropped++
			continue
		}

		st.WriteBytes(len(serializedEvent) + 1)
	}

	st.Dropped(dropped)
	st.Acked(len(events) - dropped)

	return nil
}
