package core_wal

import (
	"bytes"
	"context"
	"encoding/binary"
	"hash/crc32"
	"os"

	"github.com/raiashpanda007/kairo/shared/types"
	"github.com/raiashpanda007/kairo/shared/utils/fs"
	"go.uber.org/zap"
)

const WAL_LOG_FILE_DIR string = "./core/wal/walFiles/"

func EncodeWALData(msgData types.QueueMessage) (types.WALDataStructure, error) {

	var payloadBuffer bytes.Buffer
	msgIDBytes := []byte(msgData.MsgId)
	msgIDLength := uint32(len(msgIDBytes))
	msgBytes := []byte(msgData.Message)
	msgLength := uint32(len(msgBytes))

	err := binary.Write(&payloadBuffer, binary.LittleEndian, msgIDLength)
	if err != nil {

		return types.WALDataStructure{}, err
	}
	payloadBuffer.Write(msgIDBytes)
	err = binary.Write(&payloadBuffer, binary.LittleEndian, msgLength)

	if err != nil {
		return types.WALDataStructure{}, err
	}
	payloadBuffer.Write(msgBytes)

	payload := payloadBuffer.Bytes()

	crc := crc32.ChecksumIEEE(payload)
	recordLen := uint32(len(payload))

	return types.WALDataStructure{
		CRC:       crc,
		RecordLen: recordLen,
		Payload:   payload,
	}, nil

}

func InitWalEngine(ctx context.Context, logger *zap.Logger, walChan chan types.QueueMessage, queueName string) {

	walFilePath, err := fs.CreateFile(queueName+".log", WAL_LOG_FILE_DIR)

	if err != nil {
		logger.Error("Unable to find the wal folder for :: "+queueName, zap.Error(err))
	}

	logger.Info("Message reached to WAL engine.")

	var openedFile *os.File = nil
	openedFile, err = os.OpenFile(
		walFilePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC,
		0644,
	)

	if err != nil {
		// FIX: remove panic and Add a proper return value.
		panic("Unable to open the wal file in SYNC Mode" + err.Error())
	}

	for {
		select {
		case msgs := <-walChan:

			// TODO: Find out the best way of writting OUR LOG. RN just add message.

			encodedData, err := EncodeWALData(msgs)
			if err != nil {
				// FIX: remove panic and Add a proper return value.
				panic("Can't encode data" + err.Error())
			}

			err = binary.Write(openedFile, binary.LittleEndian, encodedData.CRC)
			if err != nil {
				// FIX: remove panic and Add a proper return value.
				panic("Unable to write the WAL log in the system" + err.Error())
			}

			err = binary.Write(openedFile, binary.LittleEndian, encodedData.RecordLen)

			if err != nil {
				panic("Unable to write the WAL log in the system" + err.Error())
			}

			err = binary.Write(openedFile, binary.LittleEndian, encodedData.Payload)

			if err != nil {
				panic("Unable to write the WAL log in the system" + err.Error())
			}

		case <-ctx.Done():
			defer openedFile.Close()
			return
		}
	}
}

/*
	WAL Record Binary Format

		We manually encode the record instead of directly serializing a Go struct
			because strings are variable-sized and Go structs contain runtime-specific
				memory layouts (such as string headers and pointers) that are not suitable
					for durable storage.

						Every record is converted into a deterministic byte layout:

								[ID_LEN][ID_BYTES][MSG_LEN][MSG_BYTES]

									ID_LEN and MSG_LEN are stored first so the decoder knows exactly how many
										bytes belong to each field during recovery.

											After building the payload, we calculate a CRC checksum over the payload
												bytes and store it in the WAL header. During recovery we recalculate the
													CRC and compare it with the stored value to detect corruption or partial
														writes caused by crashes.

							This gives us:
							- Variable-sized messages
							- Deterministic recovery
							- Corruption detection
																							- Full control over the on-disk format
*/
