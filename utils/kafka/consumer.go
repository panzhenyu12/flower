package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/golang/glog"
	systemconfig "flower/config"
	"github.com/pkg/errors"
)

type KafkaConsumer struct {
	*cluster.Consumer
	quit chan bool
}

func NewKafkaConsumer(brokers []string, topics []string) (*KafkaConsumer, error) {
	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(brokers, systemconfig.GetConfig().KafkaGroupID, topics, config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &KafkaConsumer{
		Consumer: consumer,
		quit:     make(chan bool),
	}, nil
}

func (consumer *KafkaConsumer) Start(f func(*sarama.ConsumerMessage)) {
	go func(consumer *KafkaConsumer, f func(*sarama.ConsumerMessage)) {
		for {
			select {
			case <-consumer.quit:
				consumer.Close()
				return
			case err := <-consumer.Errors():
				glog.Warningln("Kafka err:", err)
				<-time.After(time.Second)
				break
			case info := <-consumer.Notifications():
				glog.Infoln("Kafka notification:", info)
				break
			case msg := <-consumer.Messages():
				glog.Infof("%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				f(msg)
			}
		}
	}(consumer, f)
}

func (consumer *KafkaConsumer) Stop() {
	consumer.quit <- true
}
