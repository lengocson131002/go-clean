package data

import (
	"context"
	"fmt"

	"github.com/lengocson131002/go-clean/pkg/database"
	ot "github.com/lengocson131002/go-clean/pkg/trace/opentelemetry"
	"go.opentelemetry.io/otel/trace"
)

// Interface for metarepository
type MasterDataRepository interface {
	GetTemplateRequest(ctx context.Context, templateName string) (string, error)
}

type templateEntity struct {
	tName     string `db:"template_name"`
	tRequest  string `db:"template_request"`
	tResponse string `db:"template_response"`
}

type MasterDataDatabase struct {
	DB *database.Gdbc
}

type masterDataRepository struct {
	db     *database.Gdbc
	tracer *ot.OpenTelemetryTracer
}

func NewMasterDataRepository(db *MasterDataDatabase, tracer *ot.OpenTelemetryTracer) MasterDataRepository {
	return &masterDataRepository{
		db:     db.DB,
		tracer: tracer,
	}
}

func (repo *masterDataRepository) GetTemplateRequest(ctx context.Context, templateName string) (string, error) {
	ctx, span := repo.tracer.StartSpanFromContext(ctx, "Get T24 template from master data database", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	sql := "SELECT TEMPLATE_NAME, TEMPLATE_REQUEST, TEMPLATE_RESPONSE FROM GW_XSLTEMPLATES WHERE TEMPLATE_NAME = $1"

	row := repo.db.QueryRow(ctx, sql, templateName)
	if row == nil {
		return "", fmt.Errorf("Template not found")
	}

	if row.Err() != nil {
		return "", fmt.Errorf("failed to get template: %v", row.Err())
	}

	tempEntity := new(templateEntity)

	row.Scan(&tempEntity.tName, &tempEntity.tRequest, &tempEntity.tResponse)

	return tempEntity.tRequest, nil
}
