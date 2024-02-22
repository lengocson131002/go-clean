package server

import (
	"context"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/pipeline"
	pb "github.com/lengocson131002/go-clean/presentation/grpc/pb"
)

type T24AccountServer struct {
	pb.UnimplementedT24AccountServiceServer
	logger logger.Logger
}

func NewT24AccountServer(logger logger.Logger) pb.T24AccountServiceServer {
	return &T24AccountServer{
		logger: logger,
	}
}

// OpenT24Account implements grpc.T24AccountServiceServer.
func (t *T24AccountServer) OpenT24Account(ctx context.Context, req *pb.OpenT24AccountRequest) (*pb.OpenT24AccountResponse, error) {
	pipelineReq := &domain.OpenAccountRequest{
		CIF:             int(req.CIF),
		AccountTitle:    req.AccountTitle,
		ShortName:       req.ShortName,
		Category:        req.Category,
		RmCode:          req.RmCode,
		BranchCode:      req.BranchCode,
		PostingRestrict: req.PostingRestrict,
		Program:         req.Program,
		Currency:        req.Currency,
	}

	res, err := pipeline.Send[*domain.OpenAccountRequest, *domain.OpenAccountResponse](ctx, pipelineReq)
	if err != nil {
		return nil, err
	}

	return &pb.OpenT24AccountResponse{
		CIF:    int32(res.CIF),
		Status: res.Status,
	}, nil
}
