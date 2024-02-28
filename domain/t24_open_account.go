package domain

import (
	"context"
)

// Domain Model

// Events

// OPen t24 account
type OpenAccountRequest struct {
	CIF             int    `json:"CIF"`
	AccountTitle    string `json:"accountTitle"`
	ShortName       string `json:"shortName"`
	Category        string `json:"category"`
	RmCode          string `json:"rmCode"`
	BranchCode      string `json:"branchCode"`
	PostingRestrict string `json:"postingRestrict"`
	Program         string `json:"program"`
	Currency        string `json:"currency"`
}

type OpenAccountResponse struct {
	CIF    int    `json:"CIF"`
	Status string `json:"status"`
}

type OpenAccountHandler interface {
	Handle(ctx context.Context, req *OpenAccountRequest) (*OpenAccountResponse, error)
}
