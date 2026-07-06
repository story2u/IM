package connectors

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"im-go/internal/integrationhub/domain"
)

type WeComAdapter struct{}

func (WeComAdapter) Kind() domain.ChannelKind {
	return domain.ChannelWeCom
}

type WeComCallbackPayload struct {
	EventID        string         `json:"eventId"`
	MsgID          string         `json:"msgId"`
	ExternalUserID string         `json:"externalUserId"`
	UserID         string         `json:"userId"`
	SenderName     string         `json:"senderName"`
	ConversationID string         `json:"conversationId"`
	Content        string         `json:"content"`
	MsgType        string         `json:"msgType"`
	CreateTime     int64          `json:"createTime"`
	Raw            map[string]any `json:"-"`
}

func (a WeComAdapter) ParseInbound(_ context.Context, raw json.RawMessage) (domain.InboundEvent, error) {
	var payload WeComCallbackPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return domain.InboundEvent{}, err
	}
	var rawMap map[string]any
	_ = json.Unmarshal(raw, &rawMap)
	payload.Raw = rawMap

	externalID := firstNonBlank(payload.EventID, payload.MsgID)
	eventType := "message.received"
	if payload.MsgType != "" {
		eventType = "message." + strings.ToLower(payload.MsgType) + ".received"
	}
	receivedAt := time.Now().UTC()
	if payload.CreateTime > 0 {
		receivedAt = time.Unix(payload.CreateTime, 0).UTC()
	}
	conversationID := firstNonBlank(payload.ConversationID, payload.ExternalUserID, payload.UserID, "unknown")
	conversationLabel := firstNonBlank(payload.SenderName, payload.ExternalUserID, payload.UserID, "WeCom contact")

	return domain.InboundEvent{
		ConnectorKind:   a.Kind(),
		EventType:       eventType,
		ExternalEventID: externalID,
		ReceivedAt:      receivedAt,
		Status:          domain.EventPending,
		Normalized: map[string]any{
			"conversationId":    conversationID,
			"conversationLabel": conversationLabel,
			"sender":            conversationLabel,
			"content":           payload.Content,
			"messageType":       firstNonBlank(payload.MsgType, "text"),
		},
		AdapterPayload: map[string]any{
			"provider": "wecom",
			"payload":  payload.Raw,
		},
	}, nil
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
