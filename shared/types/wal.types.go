package types

type WALDataStructure struct {
	CRC       uint32
	RecordLen uint32
	Payload   []byte
}
