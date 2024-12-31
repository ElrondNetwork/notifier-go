package factory

import (
	"github.com/multiversx/mx-chain-core-go/marshal"
	marshalFactory "github.com/multiversx/mx-chain-core-go/marshal/factory"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/multiversx/mx-chain-notifier-go/rabbitmq"
	"github.com/multiversx/mx-chain-notifier-go/ws"
)

// CreatePublisher creates publisher component
func CreatePublisher(
	apiType string,
	config config.MainConfig,
	marshaller marshal.Marshalizer,
	commonHub dispatcher.Hub,
) (process.Publisher, error) {
	switch apiType {
	case common.MessageQueuePublisherType:
		return createRabbitMqPublisher(config.RabbitMQ, marshaller)
	case common.WSPublisherType:
		return createWSPublisher(commonHub)
	case common.WSPublisherTypeV2:
		return createWSPublisherV2(config.ExternalWebSocketConnector)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createRabbitMqPublisher(
	config config.RabbitMQConfig,
	marshaller marshal.Marshalizer,
) (rabbitmq.PublisherService, error) {
	rabbitClient, err := rabbitmq.NewRabbitMQClient(config.Url)
	if err != nil {
		return nil, err
	}

	rabbitMqPublisherArgs := rabbitmq.ArgsRabbitMqPublisher{
		Client:     rabbitClient,
		Config:     config,
		Marshaller: marshaller,
	}
	rabbitPublisher, err := rabbitmq.NewRabbitMqPublisher(rabbitMqPublisherArgs)
	if err != nil {
		return nil, err
	}

	return process.NewPublisher(rabbitPublisher)
}

func createWSPublisher(commonHub dispatcher.Hub) (process.Publisher, error) {
	return process.NewPublisher(commonHub)
}

func createWSPublisherV2(
	config config.WebSocketConfig,
) (process.Publisher, error) {
	marshaller, err := marshalFactory.NewMarshalizer(config.DataMarshallerType)
	if err != nil {
		return nil, err
	}

	host, err := createWsHost(config, marshaller)
	if err != nil {
		return nil, err
	}

	args := ws.WSPublisherArgs{
		Marshaller: marshaller,
		WSConn:     host,
	}

	wsPublisher, err := ws.NewWSPublisher(args)
	if err != nil {
		return nil, err
	}

	return process.NewPublisher(wsPublisher)
}
