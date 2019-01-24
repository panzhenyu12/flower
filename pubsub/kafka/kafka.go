package kafka

import (
	"errors"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/golang/protobuf/proto"
	systemconfig "github.com/panzhenyu12/flower/config"
	"github.com/panzhenyu12/flower/pubsub"
	"golang.org/x/net/context"
)

var (
	RequiredAcks = sarama.WaitForAll
)

type Publisher struct {
	producer sarama.SyncProducer
	topic    string
}

func NewPublisher(cfg *Config) (pubsub.Publisher, error) {
	var err error
	p := &Publisher{}

	if len(cfg.Topic) == 0 {
		return p, errors.New("topic name is required")
	}
	p.topic = cfg.Topic

	sconfig := cfg.Config
	if sconfig == nil {
		sconfig = sarama.NewConfig()
		sconfig.Producer.Retry.Max = cfg.MaxRetry
		sconfig.Producer.RequiredAcks = RequiredAcks
	}
	sconfig.Producer.Return.Successes = true
	p.producer, err = sarama.NewSyncProducer(cfg.BrokerHosts, sconfig)
	return p, err
}

func (p *Publisher) Publish(ctx context.Context, key string, m proto.Message) error {
	mb, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return p.PublishRaw(ctx, key, mb)
}

func (p *Publisher) PublishRaw(_ context.Context, key string, m []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(m),
	}
	_, _, err := p.producer.SendMessage(msg)
	return err
}

func (p *Publisher) Stop() error {
	return p.producer.Close()
}

type (
	subscriber struct {
		cnsmr     *cluster.Consumer
		topic     string
		partition int32

		offset          func() int64
		broadcastOffset func(int64)

		kerr error

		stop chan chan error
	}

	subMessage struct {
		message         *sarama.ConsumerMessage
		broadcastOffset func(int64)
	}
)

func (m *subMessage) Message() []byte {
	return m.message.Value
}

func (m *subMessage) ExtendDoneDeadline(time.Duration) error {
	return nil
}

func (m *subMessage) Done() error {
	m.broadcastOffset(m.message.Offset)
	return nil
}

func NewSubscriber(cfg *Config, offsetProvider func() int64, offsetBroadcast func(int64)) (pubsub.Subscriber, error) {
	var (
		err error
	)
	s := &subscriber{
		offset:          offsetProvider,
		broadcastOffset: offsetBroadcast,
		partition:       cfg.Partition,
		stop:            make(chan chan error, 1),
	}

	if len(cfg.BrokerHosts) == 0 {
		return s, errors.New("at least 1 broker host is required")
	}

	if len(cfg.Topic) == 0 {
		return s, errors.New("topic name is required")
	}
	s.topic = cfg.Topic

	sconfig := cluster.NewConfig()
	sconfig.Consumer.Return.Errors = true
	sconfig.Group.Return.Notifications = true
	s.cnsmr, err = cluster.NewConsumer(cfg.BrokerHosts, systemconfig.GetConfig().KafkaGroupID, []string{}, sconfig)
	return s, err
}

func (s *subscriber) Start() <-chan pubsub.SubscriberMessage {
	output := make(chan pubsub.SubscriberMessage, 1000)
	go func() {
		for {
			select {
			case exit := <-s.stop:
				exit <- s.cnsmr.Close()
				return
			case kerr := <-s.cnsmr.Errors():
				s.kerr = kerr
				return
			case msg, _ := <-s.cnsmr.Messages():
				output <- &subMessage{
					message:         msg,
					broadcastOffset: s.broadcastOffset,
				}
			}
		}
	}()
	return output
}

func (s *subscriber) Stop() error {
	exit := make(chan error)
	s.stop <- exit
	// close result from the partition consumer
	err := <-exit
	if err != nil {
		return err
	}
	return s.cnsmr.Close()
}

func (s *subscriber) Err() error {
	return s.kerr
}

func GetPartitions(brokerHosts []string, topic string) (partitions []int32, err error) {
	if len(brokerHosts) == 0 {
		return partitions, errors.New("at least 1 broker host is required")
	}

	if len(topic) == 0 {
		return partitions, errors.New("topic name is required")
	}

	var cnsmr sarama.Consumer
	cnsmr, err = sarama.NewConsumer(brokerHosts, sarama.NewConfig())
	if err != nil {
		return partitions, err
	}

	defer func() {
		if cerr := cnsmr.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	return cnsmr.Partitions(topic)
}
