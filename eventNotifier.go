package notifier

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core/statistics"
	nodeData "github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/indexer"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/proxy/client"
)

var log = logger.GetOrCreate("outport/eventNotifier")

const (
	pushEventEndpoint = "/events/push"
)

type eventNotifier struct {
	isNilNotifier bool
	httpClient    client.HttpClient
	marshalizer   marshal.Marshalizer
}

type EventNotifierArgs struct {
	HttpClient  client.HttpClient
	Marshalizer marshal.Marshalizer
}

// NewEventNotifier creates a new instance of the eventNotifier
// It implements all methods of process.Indexer
func NewEventNotifier(args EventNotifierArgs) (*eventNotifier, error) {
	return &eventNotifier{
		isNilNotifier: false,
		httpClient:    args.HttpClient,
		marshalizer:   args.Marshalizer,
	}, nil
}

func (en *eventNotifier) SaveBlock(args *indexer.ArgsSaveBlockData) {
	log.Debug("SaveBlock called")
	if args.TransactionsPool == nil {
		return
	}

	log.Debug("checking if block has logs", "num logs", len(args.TransactionsPool.Logs))
	var logEvents []nodeData.EventHandler
	for _, handler := range args.TransactionsPool.Logs {
		if !handler.IsInterfaceNil() {
			logEvents = append(logEvents, handler.GetLogEvents()...)
		}
	}

	if len(logEvents) == 0 {
		return
	}

	var events []data.Event
	for _, eventHandler := range logEvents {
		if !eventHandler.IsInterfaceNil() {
			var topics []string
			for _, topic := range eventHandler.GetTopics() {
				topics = append(topics, string(topic))
			}
			events = append(events, data.Event{
				Address:    string(eventHandler.GetAddress()),
				Identifier: string(eventHandler.GetIdentifier()),
				Data:       string(eventHandler.GetData()),
				Topics:     topics,
			})
		}
	}

	log.Debug("extracted events from block logs", "num events", len(events))

	err := en.httpClient.Post(pushEventEndpoint, events, nil)
	if err != nil {
		log.Error("error while posting event data", "err", err.Error())
	}
}

func (en *eventNotifier) Close() error {
	return nil
}

func (en *eventNotifier) RevertIndexedBlock(header nodeData.HeaderHandler, body nodeData.BodyHandler) {
}

func (en *eventNotifier) SaveRoundsInfo(rf []*indexer.RoundInfo) {
}

func (en *eventNotifier) SaveValidatorsRating(indexID string, validatorsRatingInfo []*indexer.ValidatorRatingInfo) {
}

func (en *eventNotifier) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) {
}

func (en *eventNotifier) UpdateTPS(tpsBenchmark statistics.TPSBenchmark) {
}

func (en *eventNotifier) SaveAccounts(timestamp uint64, accounts []state.UserAccountHandler) {
}

func (en *eventNotifier) SetTxLogsProcessor(txLogsProc process.TransactionLogProcessorDatabase) {
}

func (en *eventNotifier) IsNilIndexer() bool {
	return en.isNilNotifier
}

func (en *eventNotifier) IsInterfaceNil() bool {
	return en == nil
}
