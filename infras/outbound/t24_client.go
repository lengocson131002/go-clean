package outbound

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/ibm-messaging/mq-golang-jms20/jms20subset"
	"github.com/ibm-messaging/mq-golang-jms20/mqjms"
	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/infras/data"
	"github.com/lengocson131002/go-clean/pkg/t24/response/parser/txn"
	"github.com/lengocson131002/go-clean/pkg/trace"
	"github.com/lengocson131002/go-clean/pkg/xslt"
	"github.com/lengocson131002/go-clean/usecase/outbound"
)

const (
	OPEN_CURRENT_ACCOUNT = "OpenCurrentAccount"
)

type t24MqClient struct {
	t24Cfg *bootstrap.T24Config
	mRepo  data.MasterDataRepository
	xslt   xslt.Xslt
	tracer trace.Tracer
}

func NewT24MqClient(
	t24Config *bootstrap.T24Config,
	xslt xslt.Xslt,
	mRepo data.MasterDataRepository,
	tracer trace.Tracer,
) outbound.T24MQClient {
	return &t24MqClient{
		t24Cfg: t24Config,
		mRepo:  mRepo,
		xslt:   xslt,
		tracer: tracer,
	}
}

type t24MQOpenAccountXmlRequest struct {
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

type templateEntity struct {
	name     string `db:"template_name"`
	request  string `db:"template_request"`
	response string `db:"template_response"`
}

// ExceuteOpenAccount implements outbound.T24MQClient.
func (c *t24MqClient) ExceuteOpenAccount(ctx context.Context, request *outbound.T24MQOpenAccountRequest) (*outbound.T24MQOpenAccountResponse, error) {
	var result = new(outbound.T24MQOpenAccountResponse)

	ctx, finish := c.tracer.StartExternalTrace(ctx, "call t24 to open account")
	defer finish(ctx, trace.WithExternalResponse(result))

	// Step 1: Get template from T24
	temReq, err := c.mRepo.GetTemplateRequest(ctx, OPEN_CURRENT_ACCOUNT)
	if err != nil {
		return nil, err
	}

	// Step 2: Create T24 request
	xml, err := xml.Marshal(&t24MQOpenAccountXmlRequest{
		CIF:             request.CIF,
		AccountTitle:    request.AccountTitle,
		ShortName:       request.ShortName,
		Category:        request.Category,
		Program:         request.Program,
		RmCode:          request.RmCode,
		BranchCode:      request.BranchCode,
		PostingRestrict: request.PostingRestrict,
		Currency:        request.Currency,
		T24User:         c.t24Cfg.Username,
	})

	if err != nil {
		return nil, err
	}

	t24OfReq, err := c.xslt.Transform([]byte(temReq), xml)
	if err != nil {
		return nil, err
	}

	// Step 3: Send an receive message via MQ
	cf := mqjms.ConnectionFactoryImpl{
		Hostname:    c.t24Cfg.MqHost,
		PortNumber:  c.t24Cfg.MqPort,
		QMName:      c.t24Cfg.MqManager,
		ChannelName: c.t24Cfg.MqChannel,
		UserName:    c.t24Cfg.MqUsername,
		Password:    c.t24Cfg.MqPassword,
	}

	context, errCtx := cf.CreateContext()
	if context != nil {
		defer context.Close()
	}

	if errCtx != nil {
		return nil, fmt.Errorf("Error when creating MQ context: %w", errCtx)
	}

	requestQueue := context.CreateQueue(c.t24Cfg.MqNameIn)
	replyQueue := context.CreateQueue(c.t24Cfg.MqNameOut)

	msg := context.CreateTextMessageWithString(strings.TrimSpace(string(t24OfReq)))

	if err := context.CreateProducer().Send(requestQueue, msg); err != nil {
		return nil, fmt.Errorf("Error when sending message to MQ: %w", err)
	}

	msgSelector := "JMSCorrelationID = '" + msg.GetJMSMessageID() + "'"
	consumer, err := context.CreateConsumerWithSelector(replyQueue, msgSelector)
	if consumer != nil {
		defer consumer.Close()
	}

	if err != nil {
		return nil, fmt.Errorf("Error creating MQ reply consumer: %w", err)
	}

	respMsg, err := consumer.Receive(int32(c.t24Cfg.MqTimeout))
	if err != nil {
		return nil, fmt.Errorf("Error receive reply message: %w", err)
	}

	switch msg := respMsg.(type) {
	case jms20subset.TextMessage:
		t24ResParser := txn.NewTransactionResponseDataParser()
		res, err := t24ResParser.ParseResponseData("", *msg.GetText())
		if err != nil {
			return nil, err
		}

		if len(res.ResponseRecord.ResponseFields) == 0 {
			return nil, fmt.Errorf("No response record from T24")
		}

		result.CIF = request.CIF
		result.Status = res.ResponseCommon.Status

		return result, nil

	default:
		return nil, fmt.Errorf("Error receive reply message: Not text message")
	}

}
