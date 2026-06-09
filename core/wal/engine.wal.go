package core_wal

import (
	"github.com/raiashpanda007/kairo/shared/types"
	"go.uber.org/zap"
)

func InitWalEngine(logger *zap.Logger, walChan chan types.QueueMessage) {

	logger.Info("Message reached to WAL engine.")

}
