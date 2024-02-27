package outbound

import (
	"context"
)

type T24MQOpenAccountRequest struct {
	CIF             int
	AccountTitle    string
	ShortName       string
	Category        string
	RmCode          string
	BranchCode      string
	PostingRestrict string
	Program         string
	Currency        string
}

type T24MQOpenAccountResponse struct {
	CIF    int
	Status string
}
type T24MQClient interface {
	ExceuteOpenAccount(ctx context.Context, request *T24MQOpenAccountRequest) (*T24MQOpenAccountResponse, error)
}
