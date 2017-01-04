package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/maddin2016/perfmonbeat/config"
)

type Perfmonbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

type Data struct {
	name  string
	value *float64
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Perfmonbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Perfmonbeat) Run(b *beat.Beat) error {
	logp.Info("perfmonbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	query, err := GetHandle(bt.config.Counters)
	if err != nil {
		logp.Err("%v", err)
		bt.Stop()
	}
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		data, err := query.ReadData()
		if err != nil {
			logp.Err("%v", err)
		}
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"data":       data,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
	}
}

func (bt *Perfmonbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
