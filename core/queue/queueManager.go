package core_queue

import (
	"context"
	"sync"

	"github.com/raiashpanda007/kairo/shared/types"
	"github.com/raiashpanda007/kairo/shared/utils/fs"
	"go.uber.org/zap"
)

// TODO: Implement a queue manager that can handle multiple queues, each with its own configuration. This manager should be able to create, delete, and manage queues based on the metadata stored in the queue_config.json file. It should also handle the visibility timeout and buffer size limit for each queue when processing messages.

// First Queue Manager will read the json file -> Create an in memory map of queueName and it's channel in which the messages will be pushed -> then call a function start queue each on a sepearte routine

type QueueManager struct {
	mu              sync.Mutex
	logger          *zap.Logger
	queueMsgChanMap map[string]chan string
}

func NewQueueManager(logger *zap.Logger) QueueManager {

	return QueueManager{
		mu:              sync.Mutex{},
		logger:          logger,
		queueMsgChanMap: make(map[string]chan string),
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
func (qm *QueueManager) runQueue(ctx context.Context, queueData types.QueueMetaData, queueMsgChan chan string) {

	for {
		msg := <-queueMsgChan

		qm.logger.Info("Msg logged :: " + msg)
	}
}

func (qm *QueueManager) StartQueue(ctx context.Context, queue types.QueueMetaData) {
	qm.logger.Info("Starting queue :: " + queue.Name)
	qm.mu.Lock()
	var queueMsgChannel = make(chan string, queue.BufferSizeLimit)
	qm.queueMsgChanMap[queue.Name] = queueMsgChannel
	qm.mu.Unlock()

	go qm.runQueue(ctx, queue, queueMsgChannel)
}

func (qm *QueueManager) GetQueueMsgChanMap() map[string]chan string {
	return qm.queueMsgChanMap
}
