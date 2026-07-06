package store

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/service"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func OpenPostgres(ctx context.Context, dsn string) (*PostgresStore, error) {
	cfg, err := pgxpool.ParseConfig(strings.TrimSpace(dsn))
	if err != nil {
		return nil, err
	}
	if cfg.MaxConns == 0 {
		cfg.MaxConns = 20
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}

func NewPostgres(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}

func (s *PostgresStore) Health(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *PostgresStore) Overview(ctx context.Context) (domain.OverviewData, error) {
	channels, err := s.ListChannels(ctx)
	if err != nil {
		return domain.OverviewData{}, err
	}
	stats, err := s.overviewStats(ctx)
	if err != nil {
		return domain.OverviewData{}, err
	}
	incidents, err := s.recentIncidents(ctx)
	if err != nil {
		return domain.OverviewData{}, err
	}
	traffic, err := s.trafficSeries(ctx)
	if err != nil {
		return domain.OverviewData{}, err
	}
	return domain.OverviewData{Stats: stats, Channels: channels, Incidents: incidents, Traffic: traffic}, nil
}

func (s *PostgresStore) overviewStats(ctx context.Context) (domain.OverviewStats, error) {
	var stats domain.OverviewStats
	if err := s.pool.QueryRow(ctx, `
SELECT
  COUNT(*) FILTER (WHERE status <> 'disabled')::int,
  COUNT(*)::int,
  COALESCE(SUM(messages_today), 0)::int
FROM integration_channels`).Scan(&stats.ActiveChannels, &stats.TotalChannels, &stats.MessagesIngestedToday); err != nil {
		return domain.OverviewStats{}, err
	}
	if err := s.pool.QueryRow(ctx, `SELECT COALESCE(SUM(invocations_24h), 0)::int FROM integration_ai_policies`).Scan(&stats.AIActionsToday); err != nil {
		return domain.OverviewStats{}, err
	}
	if err := s.pool.QueryRow(ctx, `
SELECT COUNT(*)::int
FROM integration_outbound_commands
WHERE status IN ('pending', 'requires_approval')`).Scan(&stats.OutboxPending); err != nil {
		return domain.OverviewStats{}, err
	}
	var totalEvents, failedEvents int
	if err := s.pool.QueryRow(ctx, `
SELECT COUNT(*)::int, COUNT(*) FILTER (WHERE status = 'failed')::int
FROM integration_message_events
WHERE time >= now() - interval '24 hours'`).Scan(&totalEvents, &failedEvents); err != nil {
		return domain.OverviewStats{}, err
	}
	if totalEvents > 0 {
		stats.ErrorRate = float64(failedEvents) / float64(totalEvents)
	}
	if err := s.pool.QueryRow(ctx, `
SELECT COALESCE(percentile_disc(0.95) WITHIN GROUP (ORDER BY latency_ms), 0)::int
FROM integration_message_events
WHERE time >= now() - interval '24 hours'`).Scan(&stats.P95LatencyMs); err != nil {
		return domain.OverviewStats{}, err
	}
	return stats, nil
}

func (s *PostgresStore) ListChannels(ctx context.Context) ([]domain.Channel, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, connector_id, kind, name, status, receive_capabilities, send_capabilities,
       last_sync_at, error_count_24h, messages_today, active_conversations, created_at, updated_at
FROM integration_channels
ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanChannels(rows)
}

func (s *PostgresStore) CreateChannel(ctx context.Context, channel domain.Channel) (domain.Channel, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.Channel{}, err
	}
	defer rollback(ctx, tx)
	if channel.ConnectorID != nil {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_connectors (id, kind, name, status, config, created_at, updated_at)
VALUES ($1, $2, $3, $4, '{}'::jsonb, $5, $6)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, status = EXCLUDED.status, updated_at = EXCLUDED.updated_at`,
			*channel.ConnectorID, channel.Kind, channel.Name, channel.Status, channel.CreatedAt, channel.UpdatedAt); err != nil {
			return domain.Channel{}, err
		}
	}
	_, err = tx.Exec(ctx, `
INSERT INTO integration_channels (
  id, connector_id, kind, name, status, receive_capabilities, send_capabilities,
  last_sync_at, error_count_24h, messages_today, active_conversations, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		channel.ID, channel.ConnectorID, channel.Kind, channel.Name, channel.Status,
		receiveToStrings(channel.ReceiveCapabilities), sendToStrings(channel.SendCapabilities),
		channel.LastSyncAt, channel.ErrorCount24h, channel.MessagesToday, channel.ActiveConversations,
		channel.CreatedAt, channel.UpdatedAt)
	if err != nil {
		return domain.Channel{}, err
	}
	if err := insertAudit(ctx, tx, domain.AuditLogEntry{
		ID:        "aud_" + channel.ID,
		Time:      channel.CreatedAt,
		Actor:     "System",
		ActorType: domain.ActorSystem,
		Action:    "Created channel",
		Target:    channel.ID,
		Channel:   &channel.Kind,
		Result:    domain.AuditSuccess,
	}); err != nil {
		return domain.Channel{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Channel{}, err
	}
	return channel, nil
}

func (s *PostgresStore) TouchChannel(ctx context.Context, id string, now time.Time, audit domain.AuditLogEntry) (domain.Channel, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.Channel{}, err
	}
	defer rollback(ctx, tx)
	channel, err := updateChannel(ctx, tx, id, "", now)
	if err != nil {
		return domain.Channel{}, err
	}
	audit.Channel = &channel.Kind
	if err := insertAudit(ctx, tx, audit); err != nil {
		return domain.Channel{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Channel{}, err
	}
	return channel, nil
}

func (s *PostgresStore) UpdateChannelStatus(ctx context.Context, id string, status domain.ChannelStatus, now time.Time, audit domain.AuditLogEntry) (domain.Channel, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.Channel{}, err
	}
	defer rollback(ctx, tx)
	channel, err := updateChannel(ctx, tx, id, status, now)
	if err != nil {
		return domain.Channel{}, err
	}
	audit.Channel = &channel.Kind
	if err := insertAudit(ctx, tx, audit); err != nil {
		return domain.Channel{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Channel{}, err
	}
	return channel, nil
}

func updateChannel(ctx context.Context, tx pgx.Tx, id string, status domain.ChannelStatus, now time.Time) (domain.Channel, error) {
	query := `
UPDATE integration_channels
SET last_sync_at = $2, updated_at = $2
WHERE id = $1
RETURNING id, connector_id, kind, name, status, receive_capabilities, send_capabilities,
          last_sync_at, error_count_24h, messages_today, active_conversations, created_at, updated_at`
	args := []any{id, now}
	if status != "" {
		query = `
UPDATE integration_channels
SET status = $2, last_sync_at = $3, updated_at = $3
WHERE id = $1
RETURNING id, connector_id, kind, name, status, receive_capabilities, send_capabilities,
          last_sync_at, error_count_24h, messages_today, active_conversations, created_at, updated_at`
		args = []any{id, status, now}
	}
	channel, err := scanChannel(tx.QueryRow(ctx, query, args...))
	if err != nil {
		return domain.Channel{}, mapNotFound(err)
	}
	return channel, nil
}

func (s *PostgresStore) MessageFlow(ctx context.Context, filter service.MessageEventFilter) ([]domain.PipelineStageStats, []domain.MessageEvent, error) {
	statsRows, err := s.pool.Query(ctx, `
SELECT stage, label, throughput_per_min, failures_1h, avg_latency_ms
FROM integration_pipeline_stage_stats
ORDER BY CASE stage
  WHEN 'connector' THEN 1 WHEN 'ingest' THEN 2 WHEN 'normalize' THEN 3 WHEN 'store' THEN 4
  WHEN 'sop_ai' THEN 5 WHEN 'outbox' THEN 6 WHEN 'delivery' THEN 7 ELSE 99 END`)
	if err != nil {
		return nil, nil, err
	}
	defer statsRows.Close()
	stats := make([]domain.PipelineStageStats, 0)
	for statsRows.Next() {
		var item domain.PipelineStageStats
		if err := statsRows.Scan(&item.Stage, &item.Label, &item.ThroughputPerMin, &item.Failures1h, &item.AvgLatencyMs); err != nil {
			return nil, nil, err
		}
		stats = append(stats, item)
	}
	if err := statsRows.Err(); err != nil {
		return nil, nil, err
	}

	events, err := s.messageEvents(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return stats, events, nil
}

func (s *PostgresStore) messageEvents(ctx context.Context, filter service.MessageEventFilter) ([]domain.MessageEvent, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, time, channel_kind, direction, conversation_id, conversation_label, event_type, status, latency_ms, trace_id
FROM integration_message_events
WHERE ($1 = '' OR $1 = 'all' OR channel_kind = $1)
  AND ($2 = '' OR $2 = 'all' OR status = $2)
  AND ($3 = '' OR $3 = 'all' OR event_type = $3)
  AND ($4 = '' OR lower(trace_id) LIKE '%' || lower($4) || '%')
ORDER BY time DESC
LIMIT 200`, strings.TrimSpace(filter.Channel), strings.TrimSpace(filter.Status), strings.TrimSpace(filter.EventType), strings.TrimSpace(filter.TraceID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMessageEvents(rows)
}

func (s *PostgresStore) Observability(ctx context.Context) (domain.ObservabilityData, error) {
	channels, err := s.ListChannels(ctx)
	if err != nil {
		return domain.ObservabilityData{}, err
	}
	events, err := s.messageEvents(ctx, service.MessageEventFilter{})
	if err != nil {
		return domain.ObservabilityData{}, err
	}
	traffic, err := s.trafficSeries(ctx)
	if err != nil {
		return domain.ObservabilityData{}, err
	}
	stats, err := s.overviewStats(ctx)
	if err != nil {
		return domain.ObservabilityData{}, err
	}
	return domain.ObservabilityData{Channels: channels, MessageEvents: events, Traffic: traffic, Stats: stats}, nil
}

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return service.ErrNotFound
	}
	return err
}

func rollback(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}

func encodeMap(value map[string]any) []byte {
	if value == nil {
		return []byte(`{}`)
	}
	data, err := json.Marshal(value)
	if err != nil {
		return []byte(`{}`)
	}
	return data
}

func decodeMap(data []byte) map[string]any {
	if len(data) == 0 {
		return map[string]any{}
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return map[string]any{}
	}
	if value == nil {
		return map[string]any{}
	}
	return value
}
