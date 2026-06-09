package core_queue

import (
	"context"
	"sync"

	core_wal "github.com/raiashpanda007/kairo/core/wal"
	"github.com/raiashpanda007/kairo/shared/types"
	"github.com/raiashpanda007/kairo/shared/utils/fs"
	"go.uber.org/zap"
)

// First Queue Manager will read the json file -> Create an in memory map of queueName and it's channel in which the messages will be pushed -> then call a function start queue each on a sepearte routine
type QueueManager struct {
	mu              sync.Mutex
	logger          *zap.Logger
	queueMsgChanMap map[string]chan types.QueueMessage
}

func NewQueueManager(logger *zap.Logger) QueueManager {

	return QueueManager{
		mu:              sync.Mutex{},
		logger:          logger,
		queueMsgChanMap: make(map[string]chan types.QueueMessage),
	}
}

func (qm *QueueManager) Start(ctx context.Context) {
	queueMetaData, err := fs.ReadFile[[]types.QueueMetaData]("./queue_config.json")

	if err != nil {
		qm.logger.Error("Unable to read config", zap.Error(err))
	}

	for _, queue := range queueMetaData {
		qm.StartQueue(ctx, queue)
	}

}

// TODO: Complete this function
func (qm *QueueManager) runQueue(ctx context.Context, queueData types.QueueMetaData, queueMsgChan chan types.QueueMessage) {

	for {
		msg := <-queueMsgChan
		/*
			TODO: We need that to manage in-memory queue as but before we will save the msg in wal , and then we will push in channel but if our buffer (routine ) is already filled then we will
			keep a pointer on our wal then whenever channel is empty (this will be another routine that will check constantly the moment channel have some space it will start pushing messages from
			wal.)
		*/

		// TODO: 1. First create a in memory channel that will keep queueMsgs.
		// TODO: 2. then push msgs Wal engine channel .
		// TODO: 3. Then complete WAL engine.
		qm.logger.Info("Msg logged :: " + msg.Message)
	}
}

func (qm *QueueManager) StartQueue(ctx context.Context, queue types.QueueMetaData) {
	qm.logger.Info("Starting queue :: " + queue.Name)
	qm.mu.Lock()
	// Keeping the buffer size of channel as 10000 for now, can be made dynamic in future based on the queue config.
	var queueMsgChannel = make(chan types.QueueMessage, 10000)
	qm.queueMsgChanMap[queue.Name] = queueMsgChannel
	qm.mu.Unlock()

	walMsgChan := make(chan types.QueueMessage, 10000)

	go core_wal.InitWalEngine(qm.logger, walMsgChan)
	go qm.runQueue(ctx, queue, queueMsgChannel)
}

func (qm *QueueManager) GetQueueMsgChanMap() map[string]chan types.QueueMessage {
	return qm.queueMsgChanMap
}
