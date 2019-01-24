package pubsub

import (
	"time"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type Publisher interface {
	Publish(context.Context, string, proto.Message) error
	PublishRaw(context.Context, string, []byte) error
}

type MultiPublisher interface {
	Publisher
	PublishMulti(context.Context, []string, []proto.Message) error
	PublishMultiRaw(context.Context, []string, [][]byte) error
}

type Subscriber interface {
	Start() <-chan SubscriberMessage
	Err() error
	Stop() error
}

type SubscriberMessage interface {
	Message() []byte
	ExtendDoneDeadline(time.Duration) error
	Done() error
}
