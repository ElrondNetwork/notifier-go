package factory

import (
	"github.com/multiversx/mx-chain-communication-go/websocket/data"
	factoryHost "github.com/multiversx/mx-chain-communication-go/websocket/factory"
	"github.com/multiversx/mx-chain-core-go/marshal"
	marshalFactory "github.com/multiversx/mx-chain-core-go/marshal/factory"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/disabled"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher/ws"
	"github.com/multiversx/mx-chain-notifier-go/process"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

// CreateWSHandler creates websocket handler component based on api type
func CreateWSHandler(apiType string, wsDispatcher dispatcher.Dispatcher, marshaller marshal.Marshalizer) (dispatcher.WSHandler, error) {
	switch apiType {
	case common.MessageQueuePublisherType:
		return &disabled.WSHandler{}, nil
	case common.WSPublisherType:
		return createWSHandler(wsDispatcher, marshaller)
	case common.WSPublisherTypeV2:
		return &disabled.WSHandler{}, nil
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createWSHandler(wsDispatcher dispatcher.Dispatcher, marshaller marshal.Marshalizer) (dispatcher.WSHandler, error) {
	upgrader, err := ws.NewWSUpgraderWrapper(readBufferSize, writeBufferSize)
	if err != nil {
		return nil, err
	}

	args := ws.ArgsWebSocketProcessor{
		Dispatcher: wsDispatcher,
		Upgrader:   upgrader,
		Marshaller: marshaller,
	}
	return ws.NewWebSocketProcessor(args)
}

// CreateWSObserverConnector will create the web socket connector for observer node communication
func CreateWSObserverConnector(
	config config.WebSocketConfig,
	facade process.EventsFacadeHandler,
) (process.WSClient, error) {
	if config.Enabled {
		return createWsObsConnector(config, facade)
	}

	return &disabled.WSHandler{}, nil
}

func createWsObsConnector(
	config config.WebSocketConfig,
	facade process.EventsFacadeHandler,
) (process.WSClient, error) {
	marshaller, err := marshalFactory.NewMarshalizer(config.DataMarshallerType)
	if err != nil {
		return nil, err
	}

	host, err := createWsHost(config, marshaller)
	if err != nil {
		return nil, err
	}

	payloadHandler, err := CreatePayloadHandler(marshaller, facade)
	if err != nil {
		return nil, err
	}

	err = host.SetPayloadHandler(payloadHandler)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func createWsHost(wsConfig config.WebSocketConfig, wsMarshaller marshal.Marshalizer) (factoryHost.FullDuplexHost, error) {
	return factoryHost.CreateWebSocketHost(factoryHost.ArgsWebSocketHost{
		WebSocketConfig: data.WebSocketConfig{
			URL:                        wsConfig.URL,
			WithAcknowledge:            wsConfig.WithAcknowledge,
			Mode:                       wsConfig.Mode,
			RetryDurationInSec:         int(wsConfig.RetryDurationInSec),
			BlockingAckOnError:         wsConfig.BlockingAckOnError,
			AcknowledgeTimeoutInSec:    wsConfig.AcknowledgeTimeoutInSec,
			DropMessagesIfNoConnection: wsConfig.DropMessagesIfNoConnection,
		},
		Marshaller: wsMarshaller,
		Log:        log,
	})
}
