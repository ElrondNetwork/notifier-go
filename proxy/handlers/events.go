package handlers

import (
	"net/http"
	"time"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/pubsub"
	"github.com/gin-gonic/gin"
)

const (
	baseEventsEndpoint   = "/events"
	pushEventsEndpoint   = "/push"
	revertEventsEndpoint = "/revert"

	setnxRetryMs = 500

	revertKeyPrefix = "revert_"
)

type eventsHandler struct {
	notifierHub dispatcher.Hub
	config      config.ConnectorApiConfig
	redlock     *pubsub.RedlockWrapper
}

// NewEventsHandler registers handlers for the /events group
func NewEventsHandler(
	notifierHub dispatcher.Hub,
	groupHandler *GroupHandler,
	config config.ConnectorApiConfig,
	redlock *pubsub.RedlockWrapper,
) error {
	h := &eventsHandler{
		notifierHub: notifierHub,
		config:      config,
		redlock:     redlock,
	}

	endpoints := []EndpointHandler{
		{Method: http.MethodPost, Path: pushEventsEndpoint, HandlerFunc: h.pushEvents},
		{Method: http.MethodPost, Path: revertEventsEndpoint, HandlerFunc: h.revertEvents},
	}

	endpointGroupHandler := EndpointGroupHandler{
		Root:             baseEventsEndpoint,
		Middlewares:      h.createMiddlewares(),
		EndpointHandlers: endpoints,
	}

	groupHandler.AddEndpointGroupHandler(endpointGroupHandler)

	return nil
}

func (h *eventsHandler) pushEvents(c *gin.Context) {
	var blockEvents data.BlockEvents

	err := c.Bind(&blockEvents)
	if err != nil {
		JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shouldProcessEvents := true
	if h.config.CheckDuplicates {
		shouldProcessEvents = h.tryCheckProcessedOrRetry(blockEvents.Hash)
	}

	if blockEvents.Events != nil && shouldProcessEvents {
		log.Info("received events for block",
			"block hash", blockEvents.Hash,
			"shouldProcess", shouldProcessEvents,
		)
		h.notifierHub.BroadcastChan() <- blockEvents
	} else {
		log.Info("received duplicated events for block",
			"block hash", blockEvents.Hash,
			"ignoring", true,
		)
	}

	JsonResponse(c, http.StatusOK, nil, "")
}

func (h *eventsHandler) revertEvents(c *gin.Context) {
	var revertBlock data.RevertBlock

	err := c.Bind(&revertBlock)
	if err != nil {
		JsonResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shouldProcessRevert := true
	if h.config.CheckDuplicates {
		revertKey := revertKeyPrefix + revertBlock.Hash
		shouldProcessRevert = h.tryCheckProcessedOrRetry(revertKey)
	}

	if shouldProcessRevert {
		log.Info("received revert event for block",
			"block hash", revertBlock.Hash,
			"shouldProcess", shouldProcessRevert,
		)
		h.notifierHub.BroadcastRevertChan() <- revertBlock
	} else {
		log.Info("received duplicated revert event for block",
			"block hash", revertBlock.Hash,
			"ignoring", true,
		)
	}

	JsonResponse(c, http.StatusOK, nil, "")
}

func (h *eventsHandler) createMiddlewares() []gin.HandlerFunc {
	var middleware []gin.HandlerFunc

	if h.config.Username != "" && h.config.Password != "" {
		basicAuth := gin.BasicAuth(gin.Accounts{
			h.config.Username: h.config.Password,
		})
		middleware = append(middleware, basicAuth)
	}

	return middleware
}

func (h *eventsHandler) tryCheckProcessedOrRetry(blockHash string) bool {
	var err error
	var setSuccessful bool

	for {
		setSuccessful, err = h.redlock.IsBlockProcessed(blockHash)

		if err != nil {
			if !h.redlock.HasConnection() {
				time.Sleep(time.Millisecond * setnxRetryMs)
				continue
			}
		}

		break
	}

	return err == nil && setSuccessful
}
