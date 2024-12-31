package factory

import (
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/disabled"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher/hub"
	"github.com/multiversx/mx-chain-notifier-go/filters"
)

// CreateHub creates a common hub component
func CreateHub(apiType string) (dispatcher.Hub, error) {
	switch apiType {
	case common.MessageQueuePublisherType:
		return &disabled.Hub{}, nil
	case common.WSPublisherType:
		return createHub()
	case common.WSPublisherTypeV2:
		return &disabled.Hub{}, nil
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createHub() (dispatcher.Hub, error) {
	args := hub.ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
	return hub.NewCommonHub(args)
}
