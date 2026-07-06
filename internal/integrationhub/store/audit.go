package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/service"
)

func insertAudit(ctx context.Context, tx pgx.Tx, entry domain.AuditLogEntry) error {
	var channel *string
	if entry.Channel != nil {
		value := string(*entry.Channel)
		channel = &value
	}
	_, err := tx.Exec(ctx, `
INSERT INTO integration_audit_logs (id, time, actor, actor_type, action, target, channel_kind, result, ip, trace_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT (id) DO NOTHING`,
		entry.ID, entry.Time, entry.Actor, entry.ActorType, entry.Action, entry.Target, channel, entry.Result, entry.IP, entry.TraceID)
	return err
}

func (s *PostgresStore) AuditLog(ctx context.Context, filter service.AuditFilter) ([]domain.AuditLogEntry, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, time, actor, actor_type, action, target, channel_kind, result, ip, trace_id
FROM integration_audit_logs
WHERE ($1 = '' OR $1 = 'all' OR actor_type = $1)
  AND ($2 = '' OR lower(action) LIKE '%' || lower($2) || '%' OR lower(target) LIKE '%' || lower($2) || '%')
ORDER BY time DESC
LIMIT 300`, filter.ActorType, filter.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := make([]domain.AuditLogEntry, 0)
	for rows.Next() {
		var entry domain.AuditLogEntry
		var actorType, result string
		var channel *string
		if err := rows.Scan(
			&entry.ID, &entry.Time, &entry.Actor, &actorType, &entry.Action, &entry.Target,
			&channel, &result, &entry.IP, &entry.TraceID,
		); err != nil {
			return nil, err
		}
		entry.ActorType = domain.AuditActorType(actorType)
		entry.Result = domain.AuditResult(result)
		if channel != nil {
			value := domain.ChannelKind(*channel)
			entry.Channel = &value
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
