package core_wal

import (
	"os"

	"github.com/raiashpanda007/kairo/shared/types"
	"github.com/raiashpanda007/kairo/shared/utils/fs"
	"go.uber.org/zap"
)

const WAL_LOG_FILE_DIR string = "./core/wal/walFiles/"

func InitWalEngine(logger *zap.Logger, walChan chan types.QueueMessage, queueName string) {

	walFilePath, err := fs.CreateFile(queueName+".log", WAL_LOG_FILE_DIR)

	if err != nil {
		logger.Error("Unable to find the wal folder for :: "+queueName, zap.Error(err))
	}

	// TODO: Open the wal file in sync mode and append only mode. all the messages will be written in that file.

	logger.Info("Message reached to WAL engine.")

	var openedFile *os.File = nil

	for {
		select {
		case msgs := <-walChan:

			openedFile, err = os.OpenFile(
				walFilePath,
				os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC,
				0644,
			)

			if err != nil {
				// FIX: remove panic and Add a proper return value.
				panic("Unable to open the wal file in SYNC Mode" + err.Error())
			}

			// TODO: Find out the best way of writting OUR LOG. RN just add message.
			_, err = openedFile.Write([]byte(msgs.Message))

			if err != nil {
				// FIX: remove panic and Add a proper return value.
				panic("Unable to write the WAL log in the system" + err.Error())
			}
		}
	}

	// TODO: LATER on make sure that attach ctx.signal with closing the file . That won't change anything much but just to better in system.
	defer openedFile.Close()
}
