package store

import (
	"context"
	"strings"

	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/service"
)

const conversationSelect = `
SELECT id, channel_id, channel_kind, contact_name, contact_handle, last_message_preview, last_message_at,
       assigned_operator, ai_status, sop_stage, sop_workflow_id, sop_workflow_name, unread, tags, created_at, updated_at
FROM integration_conversations`

const messageSelect = `
SELECT id, conversation_id, channel_id, channel_kind, direction, author, content, message_type,
       is_ai_generated, external_message_id, created_at
FROM integration_messages`

func (s *PostgresStore) Conversations(ctx context.Context, filter service.ConversationFilter) ([]domain.Conversation, []domain.ConversationMessage, error) {
	rows, err := s.pool.Query(ctx, conversationSelect+`
WHERE ($1 = '' OR $1 = 'all' OR channel_kind = $1)
  AND ($2 = '' OR lower(contact_name || ' ' || contact_handle) LIKE '%' || lower($2) || '%')
ORDER BY last_message_at DESC
LIMIT 200`, strings.TrimSpace(filter.Channel), strings.TrimSpace(filter.Query))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	conversations, err := scanConversations(rows)
	if err != nil {
		return nil, nil, err
	}
	if len(conversations) == 0 {
		return conversations, []domain.ConversationMessage{}, nil
	}
	ids := make([]string, 0, len(conversations))
	for _, conversation := range conversations {
		ids = append(ids, conversation.ID)
	}
	messageRows, err := s.pool.Query(ctx, messageSelect+`
WHERE conversation_id = ANY($1)
ORDER BY created_at ASC`, ids)
	if err != nil {
		return nil, nil, err
	}
	defer messageRows.Close()
	messages, err := scanMessages(messageRows)
	if err != nil {
		return nil, nil, err
	}
	return conversations, messages, nil
}

func (s *PostgresStore) GetConversation(ctx context.Context, id string) (domain.Conversation, error) {
	conversation, err := scanConversation(s.pool.QueryRow(ctx, conversationSelect+` WHERE id = $1`, id))
	if err != nil {
		return domain.Conversation{}, mapNotFound(err)
	}
	return conversation, nil
}

func (s *PostgresStore) ConversationMessages(ctx context.Context, conversationID string) ([]domain.ConversationMessage, error) {
	if _, err := s.GetConversation(ctx, conversationID); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, messageSelect+`
WHERE conversation_id = $1
ORDER BY created_at ASC`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMessages(rows)
}

func (s *PostgresStore) QueueOutboundMessage(ctx context.Context, message domain.ConversationMessage, item domain.OutboxItem, event domain.MessageEvent, audit domain.AuditLogEntry) (domain.OutboxItem, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.OutboxItem{}, err
	}
	defer rollback(ctx, tx)
	if _, err := tx.Exec(ctx, `
INSERT INTO integration_messages (
  id, conversation_id, channel_id, channel_kind, direction, author, content, message_type, is_ai_generated, external_message_id, created_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		message.ID, message.ConversationID, message.ChannelID, message.Channel, message.Direction, message.Author,
		message.Content, firstNonBlank(message.MessageType, "text"), message.IsAIGenerated, message.ExternalMessageID, message.Time); err != nil {
		return domain.OutboxItem{}, err
	}
	if _, err := tx.Exec(ctx, `
UPDATE integration_conversations
SET last_message_preview = $2, last_message_at = $3, updated_at = $3
WHERE id = $1`, message.ConversationID, message.Content, message.Time); err != nil {
		return domain.OutboxItem{}, err
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO integration_outbound_commands (
  id, conversation_id, channel_id, channel_kind, conversation_label, message_id, message_type, sender,
  delivery_method, status, retry_count, last_error, payload, idempotency_key, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		item.ID, item.ConversationID, item.ChannelID, item.Channel, item.ConversationLabel, item.MessageID,
		item.MessageType, item.Sender, item.DeliveryMethod, item.Status, item.RetryCount, item.LastError,
		encodeMap(item.Payload), item.IdempotencyKey, item.CreatedAt, item.UpdatedAt); err != nil {
		return domain.OutboxItem{}, err
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO integration_message_events (
  id, time, channel_kind, direction, conversation_id, conversation_label, event_type, status, latency_ms, trace_id
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		event.ID, event.Time, event.Channel, event.Direction, event.ConversationID, event.ConversationLabel,
		event.EventType, event.Status, event.LatencyMs, event.TraceID); err != nil {
		return domain.OutboxItem{}, err
	}
	if err := insertAudit(ctx, tx, audit); err != nil {
		return domain.OutboxItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.OutboxItem{}, err
	}
	return item, nil
}
