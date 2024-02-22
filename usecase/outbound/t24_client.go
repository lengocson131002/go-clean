package outbound

import (
	"context"
	"encoding/xml"
)

type T24MQOpenAccountRequest struct {
	XMLName         xml.Name `xml:"ROOT"` // root xml
	CIF             int      `xml:"CIF"`
	AccountTitle    string   `xml:"accountTitle"`
	ShortName       string   `xml:"shortName"`
	Category        string   `xml:"category"`
	RmCode          string   `xml:"rmCode"`
	BranchCode      string   `xml:"branchCode"`
	PostingRestrict string   `xml:"postingRestrict"`
	Program         string   `xml:"program"`
	Currency        string   `xml:"currency"`
	T24User         string   `xml:"t24User"`
}

type T24MQOpenAccountResponse struct {
	CIF    int
	Status string
}
type T24MQClient interface {
	ExceuteOpenAccount(ctx context.Context, request *T24MQOpenAccountRequest) (*T24MQOpenAccountResponse, error)
}
