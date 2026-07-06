package store

import (
	"context"

	"im-go/internal/integrationhub/domain"
)

func (s *PostgresStore) RecordInboundEvent(ctx context.Context, event domain.InboundEvent, audit domain.AuditLogEntry) (domain.InboundEvent, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.InboundEvent{}, err
	}
	defer rollback(ctx, tx)
	var channelID *string
	if event.ChannelID != "" {
		channelID = &event.ChannelID
	}
	var eventError *string
	if event.Error != nil {
		eventError = event.Error
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO integration_inbound_events (
  id, channel_id, connector_kind, event_type, external_event_id, received_at,
  normalized_payload, adapter_payload, trace_id, status, error
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (connector_kind, external_event_id) WHERE external_event_id IS NOT NULL AND external_event_id <> ''
DO UPDATE SET
  received_at = EXCLUDED.received_at,
  normalized_payload = EXCLUDED.normalized_payload,
  adapter_payload = EXCLUDED.adapter_payload,
  trace_id = EXCLUDED.trace_id,
  status = EXCLUDED.status,
  error = EXCLUDED.error`,
		event.ID, channelID, event.ConnectorKind, event.EventType, event.ExternalEventID, event.ReceivedAt,
		encodeMap(event.Normalized), encodeMap(event.AdapterPayload), event.TraceID, event.Status, eventError); err != nil {
		return domain.InboundEvent{}, err
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO integration_message_events (
  id, time, channel_kind, direction, conversation_id, conversation_label, event_type, status, latency_ms, trace_id
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT (id) DO NOTHING`,
		"evt_"+event.ID, event.ReceivedAt, event.ConnectorKind, domain.DirectionInbound,
		stringFromMap(event.Normalized, "conversationId", "unknown"),
		stringFromMap(event.Normalized, "conversationLabel", "Inbound event"),
		event.EventType, event.Status, 0, event.TraceID); err != nil {
		return domain.InboundEvent{}, err
	}
	if err := insertAudit(ctx, tx, audit); err != nil {
		return domain.InboundEvent{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.InboundEvent{}, err
	}
	return event, nil
}

func stringFromMap(values map[string]any, key string, fallback string) string {
	if values == nil {
		return fallback
	}
	value, ok := values[key]
	if !ok {
		return fallback
	}
	if text, ok := value.(string); ok && text != "" {
		return text
	}
	return fallback
}
