package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/usecase/outbound"
)

type openAccountHandler struct {
	t24mq outbound.T24MQClient
}

func NewOpenAccountHandler(
	t24Cfg *bootstrap.T24Config,
	t24mq outbound.T24MQClient,

) domain.OpenAccountHandler {

	return &openAccountHandler{
		t24mq: t24mq,
	}
}

func (h *openAccountHandler) Handle(ctx context.Context, req *domain.OpenAccountRequest) (*domain.OpenAccountResponse, error) {
	t24Req := &outbound.T24MQOpenAccountRequest{
		CIF:             req.CIF,
		AccountTitle:    req.AccountTitle,
		ShortName:       req.ShortName,
		Category:        req.Category,
		RmCode:          req.RmCode,
		BranchCode:      req.BranchCode,
		PostingRestrict: req.PostingRestrict,
		Program:         req.Program,
		Currency:        req.Currency,
	}
	t24Res, err := h.t24mq.ExceuteOpenAccount(ctx, t24Req)
	if err != nil {
		return nil, err
	}

	return &domain.OpenAccountResponse{
		CIF:    t24Res.CIF,
		Status: t24Res.Status,
	}, nil
}
