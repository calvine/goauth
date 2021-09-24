package memory

import (
	"context"
	"fmt"

	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(alr.GetName()).Start(ctx, "LogMessage")
	span.SetAttributes(attribute.String("db", alr.GetType()))
	defer span.End()
	alr.logMessages = append(alr.logMessages, message)
	if alr.printToStdOut {
		fmt.Printf("AUDIT LOG MESSAGE: %v\n\n", message)
	}
	return nil
}
