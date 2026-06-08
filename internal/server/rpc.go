package server

import (
	"context"

	pb "github.com/raiashpanda007/kairo/internal/pb"
	"github.com/raiashpanda007/kairo/shared/types"
	"github.com/raiashpanda007/kairo/shared/utils/fs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const QUEUE_META_DATA string = "./queue_config.json"

func (s *KairoServerStruct) CreateQueue(ctx context.Context, req *pb.CreateQueueRequest) (*pb.CreateQueueResponse, error) {

	s.Logger.Info("Create New Queue Requested")

	if req.GetQueueName() == "" {
		return nil, status.Error(codes.InvalidArgument, "queue name is required")
	}

	var QueueData types.QueueMetaData
	QueueData.Name = req.GetQueueName()

	if req.VisibilityTimeout == nil {
		QueueData.VisibilityTimeOut = 10 * 1000
	} else {
		if req.GetVisibilityTimeout() <= 0 || req.GetVisibilityTimeout() >= 10*60*1000 {
			return nil, status.Error(codes.InvalidArgument, "visibility timeout must be between 0 and 10 minutes (in ms)")
		}
		QueueData.VisibilityTimeOut = req.GetVisibilityTimeout()
	}

	if req.BufferSizeLimit == nil {
		QueueData.BufferSizeLimit = 100
	} else {
		if req.GetBufferSizeLimit() <= 0 || req.GetBufferSizeLimit() >= 10000 {
			return nil, status.Error(codes.InvalidArgument, "buffer size limit must be between 0 and 10000")
		}
		QueueData.BufferSizeLimit = req.GetBufferSizeLimit()
	}

	// Read existing queue config.
	queueMetaData, err := fs.ReadFile[[]types.QueueMetaData](QUEUE_META_DATA)
	if err != nil {
		s.Logger.Error("Failed to read queue config" + err.Error())
		return nil, status.Error(codes.Internal, "unable to read queue config: "+err.Error())
	}

	// Reject duplicate queue names.
	for _, v := range queueMetaData {
		if v.Name == req.GetQueueName() {
			return nil, status.Errorf(codes.AlreadyExists, "queue %q already exists", req.GetQueueName())
		}
	}

	queueMetaData = append(queueMetaData, QueueData)

	if err := fs.WriteJSONFile("./", QUEUE_META_DATA, queueMetaData, 0644); err != nil {
		s.Logger.Error("Failed to persist queue config" + err.Error())
		return nil, status.Error(codes.Internal, "unable to save queue config: "+err.Error())
	}

	s.QueueManager.StartQueue(ctx, QueueData)

	return &pb.CreateQueueResponse{
		Message: "queue created successfully",
		Status:  pb.ResponseStatus_OK_STATUS,
		MetaData: &pb.QueueMetaData{
			QueueName:         QueueData.Name,
			VisibilityTimeout: QueueData.VisibilityTimeOut,
			BufferSizeLimit:   QueueData.BufferSizeLimit,
		},
	}, nil
}

func (s *KairoServerStruct) Enqueue(ctx context.Context, req *pb.EnqueueRequest) (*pb.EnqueueResponse, error) {
	queueName := req.GetQueueName()
	message := req.GetMessage()
	queueChanMap := s.QueueManager.GetQueueMsgChanMap()

	queueChan, ok := queueChanMap[queueName]

	if !ok {
		return nil, status.Errorf(codes.NotFound, "queue %q not found", queueName)
	}

	queueChan <- message

	return &pb.EnqueueResponse{
		Message: "SUCCESSFULLY MESSAGE ENQUEUED",
		Status:  pb.ResponseStatus_OK_STATUS,
	}, nil
}

// func (s *KairoServerStruct) GetAllQueues(ctx context.Context, req *pb.GetAllQueueRequest) (*pb.GetAllQueueResponse, error) {}

// func (s *KairoServerStruct) DeleteQueue(ctx context.Context, req *pb.DeleteQueueRequest) (*pb.DeleteQueueResponse, error) {}

// func (s *KairoServerStruct) BulkEnqueue(ctx context.Context, req *pb.BulkEnqueueRequest) (*pb.EnqueueResponse, error) {}
