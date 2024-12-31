package ws

import (
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

type WSConnection interface {
	Send(payload []byte, topic string) error
	Close() error
	IsInterfaceNil() bool
}

var log = logger.GetOrCreate("ws")

type WSPublisherArgs struct {
	Marshaller marshal.Marshalizer
	WSConn     WSConnection
}

type wsPublisher struct {
	marshaller marshal.Marshalizer
	wsConn     WSConnection
}

// NewWSPublisher will create a new instance of websocket publisher
func NewWSPublisher(args WSPublisherArgs) (*wsPublisher, error) {
	return &wsPublisher{
		marshaller: args.Marshaller,
		wsConn:     args.WSConn,
	}, nil
}

// Publish will publish logs and events to websocket clients
func (w *wsPublisher) Publish(events data.BlockEvents) {
	eventBytes, err := w.marshaller.Marshal(events)
	if err != nil {
		log.Error("failure marshalling events", "err", err.Error())
		return
	}

	w.wsConn.Send(eventBytes, outport.TopicSaveBlock)
}

// PublishRevert will publish revert event to websocket clients
func (w *wsPublisher) PublishRevert(revertBlock data.RevertBlock) {
}

// PublishFinalized will publish finalized event to websocket clients
func (w *wsPublisher) PublishFinalized(finalizedBlock data.FinalizedBlock) {
}

// PublishTxs will publish txs event to websocket clients
func (w *wsPublisher) PublishTxs(blockTxs data.BlockTxs) {
}

// PublishScrs will publish scrs event to websocket clients
func (w *wsPublisher) PublishScrs(blockScrs data.BlockScrs) {
}

// PublishBlockEventsWithOrder will publish block events with order to websocket clients
func (w *wsPublisher) PublishBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
}

// Close will trigger to close ws publisher
func (w *wsPublisher) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (w *wsPublisher) IsInterfaceNil() bool {
	return w == nil
}
