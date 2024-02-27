package data

import (
	"context"
	"fmt"

	"github.com/lengocson131002/go-clean/pkg/database"
	"github.com/lengocson131002/go-clean/pkg/trace"
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
	tracer trace.Tracer
}

func NewMasterDataRepository(db *MasterDataDatabase, tracer trace.Tracer) MasterDataRepository {
	return &masterDataRepository{
		db:     db.DB,
		tracer: tracer,
	}
}

func (repo *masterDataRepository) GetTemplateRequest(ctx context.Context, templateName string) (string, error) {
	sql := "SELECT TEMPLATE_NAME, TEMPLATE_REQUEST, TEMPLATE_RESPONSE FROM GW_XSLTEMPLATES WHERE TEMPLATE_NAME = $1"

	ctx, finish := repo.tracer.StartDatabaseTrace(
		ctx,
		"get template request from master database",
		trace.WithDBTableName("GW_XSLTEMPLATES"),
		trace.WithDBSql(sql),
	)

	defer finish(ctx)

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
