package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"im-go/internal/integrationhub/domain"
)

func (s *PostgresStore) SeedIfEmpty(ctx context.Context, snapshot domain.Snapshot) error {
	var count int
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM integration_channels`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer rollback(ctx, tx)
	if err := seedConnectors(ctx, tx, snapshot.Connectors); err != nil {
		return err
	}
	if err := seedChannels(ctx, tx, snapshot.Channels); err != nil {
		return err
	}
	if err := seedChannelAccounts(ctx, tx, snapshot.ChannelAccounts); err != nil {
		return err
	}
	if err := seedSOPWorkflows(ctx, tx, snapshot.SOPWorkflows); err != nil {
		return err
	}
	if err := seedConversations(ctx, tx, snapshot.Conversations); err != nil {
		return err
	}
	if err := seedMessages(ctx, tx, snapshot.Messages); err != nil {
		return err
	}
	if err := seedAIPolicies(ctx, tx, snapshot.AIPolicies); err != nil {
		return err
	}
	if err := seedOutbox(ctx, tx, snapshot.OutboxItems); err != nil {
		return err
	}
	if err := seedMessageEvents(ctx, tx, snapshot.MessageEvents); err != nil {
		return err
	}
	if err := seedAudit(ctx, tx, snapshot.AuditLog); err != nil {
		return err
	}
	if err := seedIncidents(ctx, tx, snapshot.Incidents); err != nil {
		return err
	}
	if err := seedPipelineStats(ctx, tx, snapshot.PipelineStats); err != nil {
		return err
	}
	if err := seedTraffic(ctx, tx, snapshot.TrafficSeries); err != nil {
		return err
	}
	if err := seedSettings(ctx, tx, snapshot.Settings); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func seedConnectors(ctx context.Context, tx pgx.Tx, connectors []domain.Connector) error {
	for _, connector := range connectors {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_connectors (id, kind, name, status, config, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7)
ON CONFLICT (id) DO NOTHING`,
			connector.ID, connector.Kind, connector.Name, connector.Status, encodeMap(connector.Config), connector.CreatedAt, connector.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

func seedChannels(ctx context.Context, tx pgx.Tx, channels []domain.Channel) error {
	for _, channel := range channels {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_channels (
  id, connector_id, kind, name, status, receive_capabilities, send_capabilities,
  last_sync_at, error_count_24h, messages_today, active_conversations, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
ON CONFLICT (id) DO NOTHING`,
			channel.ID, channel.ConnectorID, channel.Kind, channel.Name, channel.Status,
			receiveToStrings(channel.ReceiveCapabilities), sendToStrings(channel.SendCapabilities),
			channel.LastSyncAt, channel.ErrorCount24h, channel.MessagesToday, channel.ActiveConversations,
			channel.CreatedAt, channel.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

func seedChannelAccounts(ctx context.Context, tx pgx.Tx, accounts []domain.ChannelAccount) error {
	for _, account := range accounts {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_channel_accounts (id, channel_id, display_name, external_account_id, status, config, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (id) DO NOTHING`,
			account.ID, account.ChannelID, account.DisplayName, account.ExternalAccountID, account.Status,
			encodeMap(account.Config), account.CreatedAt, account.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

func seedSOPWorkflows(ctx context.Context, tx pgx.Tx, workflows []domain.SOPWorkflow) error {
	for _, workflow := range workflows {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_sop_workflows (
  id, name, trigger_expr, channels, active_conversations, completion_rate, sla_minutes, status, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT (id) DO NOTHING`,
			workflow.ID, workflow.Name, workflow.Trigger, channelKindsToStrings(workflow.Channels),
			workflow.ActiveConversations, workflow.CompletionRate, workflow.SLAMinutes, workflow.Status,
			workflow.CreatedAt, workflow.UpdatedAt); err != nil {
			return err
		}
		for i, step := range workflow.Steps {
			position := step.Position
			if position == 0 {
				position = i + 1
			}
			if _, err := tx.Exec(ctx, `
INSERT INTO integration_sop_steps (
  id, workflow_id, position, name, condition_expr, ai_action, human_action, timeout_minutes, fallback
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
ON CONFLICT (id) DO NOTHING`,
				step.ID, workflow.ID, position, step.Name, step.Condition, step.AIAction, step.HumanAction,
				step.TimeoutMinutes, step.Fallback); err != nil {
				return err
			}
		}
	}
	return nil
}

func seedConversations(ctx context.Context, tx pgx.Tx, conversations []domain.Conversation) error {
	for _, conversation := range conversations {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_conversations (
  id, channel_id, channel_kind, contact_name, contact_handle, last_message_preview, last_message_at,
  assigned_operator, ai_status, sop_stage, sop_workflow_id, sop_workflow_name, unread, tags, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
ON CONFLICT (id) DO NOTHING`,
			conversation.ID, conversation.ChannelID, conversation.Channel, conversation.ContactName, conversation.ContactHandle,
			conversation.LastMessagePreview, conversation.LastMessageAt, conversation.AssignedOperator,
			conversation.AIStatus, conversation.SOPStage, conversation.SOPWorkflowID, conversation.SOPWorkflowName,
			conversation.Unread, conversation.Tags, conversation.CreatedAt, conversation.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

func seedMessages(ctx context.Context, tx pgx.Tx, messages []domain.ConversationMessage) error {
	for _, message := range messages {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_messages (
  id, conversation_id, channel_id, channel_kind, direction, author, content, message_type, is_ai_generated, external_message_id, created_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (id) DO NOTHING`,
			message.ID, message.ConversationID, message.ChannelID, message.Channel, message.Direction, message.Author,
			message.Content, firstNonBlank(message.MessageType, "text"), message.IsAIGenerated, message.ExternalMessageID,
			message.Time); err != nil {
			return err
		}
	}
	return nil
}

func seedAIPolicies(ctx context.Context, tx pgx.Tx, policies []domain.AIPolicy) error {
	for _, policy := range policies {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_ai_policies (
  id, kind, name, enabled, priority, trigger_condition, fallback_strategy,
  success_rate_7d, invocations_24h, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (id) DO NOTHING`,
			policy.ID, policy.Kind, policy.Name, policy.Enabled, policy.Priority, policy.TriggerCondition,
			policy.FallbackStrategy, policy.SuccessRate7d, policy.Invocations24h, policy.CreatedAt, policy.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

func seedOutbox(ctx context.Context, tx pgx.Tx, items []domain.OutboxItem) error {
	for _, item := range items {
		idempotencyKey := item.IdempotencyKey
		if idempotencyKey == "" {
			idempotencyKey = "seed-" + item.ID
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_outbound_commands (
  id, conversation_id, channel_id, channel_kind, conversation_label, message_id, message_type, sender,
  delivery_method, status, retry_count, last_error, payload, idempotency_key, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
ON CONFLICT (id) DO NOTHING`,
			item.ID, item.ConversationID, item.ChannelID, item.Channel, item.ConversationLabel, item.MessageID,
			item.MessageType, item.Sender, item.DeliveryMethod, item.Status, item.RetryCount, item.LastError,
			encodeMap(item.Payload), idempotencyKey, item.CreatedAt, firstTime(item.UpdatedAt, item.CreatedAt)); err != nil {
			return err
		}
	}
	return nil
}

func seedMessageEvents(ctx context.Context, tx pgx.Tx, events []domain.MessageEvent) error {
	for _, event := range events {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_message_events (
  id, time, channel_kind, direction, conversation_id, conversation_label, event_type, status, latency_ms, trace_id
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT (id) DO NOTHING`,
			event.ID, event.Time, event.Channel, event.Direction, event.ConversationID, event.ConversationLabel,
			event.EventType, event.Status, event.LatencyMs, event.TraceID); err != nil {
			return err
		}
	}
	return nil
}

func seedAudit(ctx context.Context, tx pgx.Tx, entries []domain.AuditLogEntry) error {
	for _, entry := range entries {
		if err := insertAudit(ctx, tx, entry); err != nil {
			return err
		}
	}
	return nil
}

func seedIncidents(ctx context.Context, tx pgx.Tx, incidents []domain.RecentIncident) error {
	for _, incident := range incidents {
		var channel *string
		if incident.Channel != nil {
			value := string(*incident.Channel)
			channel = &value
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_incidents (id, time, severity, summary, channel_kind)
VALUES ($1,$2,$3,$4,$5)
ON CONFLICT (id) DO NOTHING`,
			incident.ID, incident.Time, incident.Severity, incident.Summary, channel); err != nil {
			return err
		}
	}
	return nil
}

func seedPipelineStats(ctx context.Context, tx pgx.Tx, stats []domain.PipelineStageStats) error {
	for _, item := range stats {
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_pipeline_stage_stats (stage, label, throughput_per_min, failures_1h, avg_latency_ms)
VALUES ($1,$2,$3,$4,$5)
ON CONFLICT (stage) DO UPDATE SET
  label = EXCLUDED.label,
  throughput_per_min = EXCLUDED.throughput_per_min,
  failures_1h = EXCLUDED.failures_1h,
  avg_latency_ms = EXCLUDED.avg_latency_ms,
  updated_at = now()`,
			item.Stage, item.Label, item.ThroughputPerMin, item.Failures1h, item.AvgLatencyMs); err != nil {
			return err
		}
	}
	return nil
}

func seedTraffic(ctx context.Context, tx pgx.Tx, points []domain.TrafficPoint) error {
	for i, point := range points {
		hourIndex := i
		if parsed, err := strconv.Atoi(strings.TrimSuffix(point.Hour, ":00")); err == nil {
			hourIndex = parsed
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO integration_traffic_points (hour_index, hour_label, inbound, outbound)
VALUES ($1,$2,$3,$4)
ON CONFLICT (hour_index) DO UPDATE SET
  hour_label = EXCLUDED.hour_label,
  inbound = EXCLUDED.inbound,
  outbound = EXCLUDED.outbound,
  updated_at = now()`,
			hourIndex, point.Hour, point.Inbound, point.Outbound); err != nil {
			return err
		}
	}
	return nil
}

func seedSettings(ctx context.Context, tx pgx.Tx, settings domain.PlatformSettings) error {
	if _, err := tx.Exec(ctx, `
INSERT INTO integration_platform_settings (
  id, workspace_name, timezone, default_language, environment, region, retention_days, webhook_url, enabled_providers
) VALUES ('default',$1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (id) DO NOTHING`,
		firstNonBlank(settings.WorkspaceName, "IM Integration Platform"),
		firstNonBlank(settings.Timezone, "UTC"),
		firstNonBlank(settings.DefaultLanguage, "en"),
		firstNonBlank(settings.Environment, "development"),
		firstNonBlank(settings.Region, "local"),
		settings.RetentionDays,
		settings.WebhookURL,
		settings.EnabledProviders); err != nil {
		return fmt.Errorf("seed settings: %w", err)
	}
	return nil
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func firstTime(values ...time.Time) time.Time {
	for _, value := range values {
		if !value.IsZero() {
			return value
		}
	}
	return values[len(values)-1]
}
