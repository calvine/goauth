package memory

import (
	"context"
	"fmt"

	"github.com/calvine/goauth/core/apptelemetry"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
)

type auditLogRepo struct {
	logMessages   []models.AuditLog
	printToStdOut bool
}

func NewMemoryAuditLogRepo(printToStdOut bool) repo.AuditLogRepo {
	auditLog := make([]models.AuditLog, 0)
	return &auditLogRepo{auditLog, printToStdOut}

}

func (auditLogRepo) GetName() string {
	return "auditLogRepo"
}

func (auditLogRepo) GetType() string {
	return dataSourceType
}

func (alr *auditLogRepo) LogMessage(ctx context.Context, message models.AuditLog) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, alr.GetName(), "LogMessage", alr.GetType())
	defer span.End()
	alr.logMessages = append(alr.logMessages, message)
	if alr.printToStdOut {
		fmt.Printf("AUDIT LOG MESSAGE: %v\n\n", message)
	}
	span.AddEvent("log message stored")
	return nil
}
