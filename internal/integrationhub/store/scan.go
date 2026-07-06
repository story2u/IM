package store

import (
	"github.com/jackc/pgx/v5"

	"im-go/internal/integrationhub/domain"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func scanChannel(row rowScanner) (domain.Channel, error) {
	var channel domain.Channel
	var kind, status string
	var receive, send []string
	if err := row.Scan(
		&channel.ID, &channel.ConnectorID, &kind, &channel.Name, &status, &receive, &send,
		&channel.LastSyncAt, &channel.ErrorCount24h, &channel.MessagesToday, &channel.ActiveConversations,
		&channel.CreatedAt, &channel.UpdatedAt,
	); err != nil {
		return domain.Channel{}, err
	}
	channel.Kind = domain.ChannelKind(kind)
	channel.Status = domain.ChannelStatus(status)
	channel.ReceiveCapabilities = receiveFromStrings(receive)
	channel.SendCapabilities = sendFromStrings(send)
	return channel, nil
}

func scanChannels(rows pgx.Rows) ([]domain.Channel, error) {
	channels := make([]domain.Channel, 0)
	for rows.Next() {
		channel, err := scanChannel(rows)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	return channels, rows.Err()
}

func scanMessageEvents(rows pgx.Rows) ([]domain.MessageEvent, error) {
	events := make([]domain.MessageEvent, 0)
	for rows.Next() {
		var event domain.MessageEvent
		var channel, direction, status string
		if err := rows.Scan(
			&event.ID, &event.Time, &channel, &direction, &event.ConversationID, &event.ConversationLabel,
			&event.EventType, &status, &event.LatencyMs, &event.TraceID,
		); err != nil {
			return nil, err
		}
		event.Channel = domain.ChannelKind(channel)
		event.Direction = domain.MessageDirection(direction)
		event.Status = domain.MessageEventStatus(status)
		events = append(events, event)
	}
	return events, rows.Err()
}

func scanConversation(row rowScanner) (domain.Conversation, error) {
	var conversation domain.Conversation
	var channel, aiStatus, sopStage string
	var tags []string
	if err := row.Scan(
		&conversation.ID, &conversation.ChannelID, &channel, &conversation.ContactName, &conversation.ContactHandle,
		&conversation.LastMessagePreview, &conversation.LastMessageAt, &conversation.AssignedOperator,
		&aiStatus, &sopStage, &conversation.SOPWorkflowID, &conversation.SOPWorkflowName,
		&conversation.Unread, &tags, &conversation.CreatedAt, &conversation.UpdatedAt,
	); err != nil {
		return domain.Conversation{}, err
	}
	conversation.Channel = domain.ChannelKind(channel)
	conversation.AIStatus = domain.AIStatus(aiStatus)
	conversation.SOPStage = domain.SOPStage(sopStage)
	conversation.Tags = tags
	return conversation, nil
}

func scanConversations(rows pgx.Rows) ([]domain.Conversation, error) {
	conversations := make([]domain.Conversation, 0)
	for rows.Next() {
		conversation, err := scanConversation(rows)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}
	return conversations, rows.Err()
}

func scanMessage(row rowScanner) (domain.ConversationMessage, error) {
	var message domain.ConversationMessage
	var channel, direction string
	if err := row.Scan(
		&message.ID, &message.ConversationID, &message.ChannelID, &channel, &direction,
		&message.Author, &message.Content, &message.MessageType, &message.IsAIGenerated,
		&message.ExternalMessageID, &message.Time,
	); err != nil {
		return domain.ConversationMessage{}, err
	}
	message.Channel = domain.ChannelKind(channel)
	message.Direction = domain.MessageDirection(direction)
	return message, nil
}

func scanMessages(rows pgx.Rows) ([]domain.ConversationMessage, error) {
	messages := make([]domain.ConversationMessage, 0)
	for rows.Next() {
		message, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func scanOutboxItem(row rowScanner) (domain.OutboxItem, error) {
	var item domain.OutboxItem
	var channel, delivery, status string
	var payload []byte
	if err := row.Scan(
		&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.ChannelID, &channel, &item.ConversationID,
		&item.ConversationLabel, &item.MessageID, &item.MessageType, &item.Sender, &delivery,
		&status, &item.RetryCount, &item.LastError, &payload, &item.IdempotencyKey,
	); err != nil {
		return domain.OutboxItem{}, err
	}
	item.Channel = domain.ChannelKind(channel)
	item.DeliveryMethod = domain.SendCapability(delivery)
	item.Status = domain.OutboxStatus(status)
	item.Payload = decodeMap(payload)
	return item, nil
}

func scanOutboxItems(rows pgx.Rows) ([]domain.OutboxItem, error) {
	items := make([]domain.OutboxItem, 0)
	for rows.Next() {
		item, err := scanOutboxItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func receiveToStrings(values []domain.ReceiveCapability) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, string(value))
	}
	return out
}

func receiveFromStrings(values []string) []domain.ReceiveCapability {
	out := make([]domain.ReceiveCapability, 0, len(values))
	for _, value := range values {
		out = append(out, domain.ReceiveCapability(value))
	}
	return out
}

func sendToStrings(values []domain.SendCapability) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, string(value))
	}
	return out
}

func sendFromStrings(values []string) []domain.SendCapability {
	out := make([]domain.SendCapability, 0, len(values))
	for _, value := range values {
		out = append(out, domain.SendCapability(value))
	}
	return out
}

func channelKindsToStrings(values []domain.ChannelKind) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, string(value))
	}
	return out
}

func channelKindsFromStrings(values []string) []domain.ChannelKind {
	out := make([]domain.ChannelKind, 0, len(values))
	for _, value := range values {
		out = append(out, domain.ChannelKind(value))
	}
	return out
}
