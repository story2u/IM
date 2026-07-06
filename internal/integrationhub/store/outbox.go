package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/service"
)

const outboxSelect = `
SELECT id, created_at, updated_at, channel_id, channel_kind, conversation_id, conversation_label,
       message_id, message_type, sender, delivery_method, status, retry_count, last_error, payload, idempotency_key
FROM integration_outbound_commands`

func (s *PostgresStore) ListOutbox(ctx context.Context, status string) ([]domain.OutboxItem, error) {
	status = strings.TrimSpace(status)
	rows, err := s.pool.Query(ctx, outboxSelect+`
WHERE status <> 'canceled'
  AND ($1 = '' OR $1 = 'all' OR status = $1)
ORDER BY created_at DESC
LIMIT 300`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOutboxItems(rows)
}

func (s *PostgresStore) MoveOutbox(ctx context.Context, id string, status domain.OutboxStatus, incrementRetry bool, now time.Time, audit domain.AuditLogEntry) (domain.OutboxItem, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.OutboxItem{}, err
	}
	defer rollback(ctx, tx)

	item, err := scanOutboxItem(tx.QueryRow(ctx, outboxSelect+` WHERE id = $1 FOR UPDATE`, id))
	if err != nil {
		return domain.OutboxItem{}, mapNotFound(err)
	}
	if err := validateOutboxTransition(item.Status, status, incrementRetry); err != nil {
		return domain.OutboxItem{}, err
	}
	retryCount := item.RetryCount
	if incrementRetry {
		retryCount++
	}
	var lastError *string = item.LastError
	if status == domain.OutboxSending {
		lastError = nil
	}
	approvedAt := (*time.Time)(nil)
	canceledAt := (*time.Time)(nil)
	if status == domain.OutboxSending && item.Status == domain.OutboxRequiresApproval {
		approvedAt = &now
	}
	if status == domain.OutboxCanceled {
		canceledAt = &now
	}
	updated, err := scanOutboxItem(tx.QueryRow(ctx, `
UPDATE integration_outbound_commands
SET status = $2,
    retry_count = $3,
    last_error = $4,
    updated_at = $5,
    approved_at = COALESCE($6::timestamptz, approved_at),
    canceled_at = COALESCE($7::timestamptz, canceled_at)
WHERE id = $1
RETURNING id, created_at, updated_at, channel_id, channel_kind, conversation_id, conversation_label,
          message_id, message_type, sender, delivery_method, status, retry_count, last_error, payload, idempotency_key`,
		id, status, retryCount, lastError, now, approvedAt, canceledAt))
	if err != nil {
		return domain.OutboxItem{}, err
	}
	audit.Channel = &updated.Channel
	if err := insertAudit(ctx, tx, audit); err != nil {
		return domain.OutboxItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.OutboxItem{}, err
	}
	return updated, nil
}

func validateOutboxTransition(current domain.OutboxStatus, next domain.OutboxStatus, incrementRetry bool) error {
	if incrementRetry {
		if current != domain.OutboxFailed {
			return fmt.Errorf("%w: only failed outbox messages can be retried", service.ErrValidation)
		}
		return nil
	}
	switch next {
	case domain.OutboxSending:
		if current != domain.OutboxRequiresApproval {
			return fmt.Errorf("%w: only messages requiring approval can be approved", service.ErrValidation)
		}
	case domain.OutboxCanceled:
		if current != domain.OutboxPending && current != domain.OutboxRequiresApproval {
			return fmt.Errorf("%w: only pending or approval-required messages can be canceled", service.ErrValidation)
		}
	}
	return nil
}
