package connector

import (
	"fmt"
	"strings"
	"time"

	"im-go/internal/tasks"
)

// OutboundTaskOptions supplies deployment-level connector defaults.
type OutboundTaskOptions struct {
	ConnectorID string
	Channel     string
	TenantID    string
	AccountID   string
	EndpointID  string
}

// DeliveryReceiptOptions supplies deployment-level receipt defaults.
type DeliveryReceiptOptions struct {
	ConnectorID string
	Channel     string
	TenantID    string
}

// OutboundMessageFromTask maps durable send tasks to the connector send contract.
func OutboundMessageFromTask(record tasks.Record, options OutboundTaskOptions) (OutboundMessage, bool) {
	taskID := strings.TrimSpace(record.TaskID)
	if taskID == "" {
		return OutboundMessage{}, false
	}
	messageType, ok := messageTypeFromTask(record.TaskType)
	if !ok {
		return OutboundMessage{}, false
	}
	payload := record.Payload
	if payload == nil {
		payload = map[string]any{}
	}
	traceID := stringPointerValue(record.TraceID)
	messageID := firstNonBlank(traceID, taskID)
	receiver := firstNonBlank(payloadText(payload, "receiver"), payloadText(payload, "username"))
	receiverName := firstNonBlank(payloadText(payload, "receiver_name"), receiver)
	channel := firstNonBlank(options.Channel, ChannelInternalWebhook)
	channelUserID := firstNonBlank(stringPointerValue(record.ChannelUserID), payloadText(payload, "channel_user_id"), payloadText(payload, "wework_user_id"), stringPointerValue(record.WeWorkUserID))
	outbound := OutboundMessage{
		MessageID:      messageID,
		TraceID:        traceID,
		IdempotencyKey: idempotencyKey(record, payload),
		ConnectorID:    firstNonBlank(options.ConnectorID, "default-send-connector"),
		Channel:        channel,
		TenantID:       firstNonBlank(options.TenantID, "default"),
		AccountID:      strings.TrimSpace(options.AccountID),
		ChannelUserID:  channelUserID,
		EndpointID:     firstNonBlank(options.EndpointID, record.Target.DeviceID, record.Target.AgentID),
		Target: ContactIdentity{
			ExternalUserID: receiver,
			DisplayName:    receiverName,
			Remark:         payloadText(payload, "aliases"),
		},
		Conversation: ConversationBinding{
			ConversationID: firstNonBlank(payloadText(payload, "conversation_id"), payloadText(payload, "session_id")),
			Type:           "single",
			DisplayName:    receiverName,
		},
		MessageType: messageType,
		Content:     payloadText(payload, "text"),
		CreatedAt:   record.CreatedAt.UTC(),
		Metadata: map[string]any{
			"task_id":   taskID,
			"task_type": strings.TrimSpace(record.TaskType),
			"source":    strings.TrimSpace(record.Source),
			"agent_id":  strings.TrimSpace(record.Target.AgentID),
			"device_id": strings.TrimSpace(record.Target.DeviceID),
		},
	}
	if outbound.CreatedAt.IsZero() {
		outbound.CreatedAt = time.Time{}
	}
	if senderID := payloadText(payload, "sender_id"); senderID != "" {
		outbound.Metadata["sender_id"] = senderID
	}
	if channelUserID != "" {
		outbound.Metadata["channel_user_id"] = channelUserID
		if weworkUserID := firstNonBlank(payloadText(payload, "wework_user_id"), stringPointerValue(record.WeWorkUserID)); weworkUserID != "" {
			outbound.Metadata["wework_user_id"] = weworkUserID
		}
	}
	if messageType != MessageTypeText {
		media := mediaFromTask(taskID, messageType, payload)
		if len(media) == 0 {
			return OutboundMessage{}, false
		}
		outbound.Media = media
	}
	return outbound, true
}

// DeliveryReceiptFromTask maps terminal task state to the connector receipt contract.
func DeliveryReceiptFromTask(record tasks.Record, options DeliveryReceiptOptions) (DeliveryReceipt, bool) {
	traceID := stringPointerValue(record.TraceID)
	taskID := strings.TrimSpace(record.TaskID)
	if traceID == "" && taskID == "" {
		return DeliveryReceipt{}, false
	}
	status := strings.ToLower(strings.TrimSpace(string(record.Status)))
	receiptStatus := ""
	switch status {
	case string(tasks.StatusSuccess):
		receiptStatus = ReceiptDelivered
	case string(tasks.StatusFailed), string(tasks.StatusCancelled), string(tasks.StatusTimeout):
		receiptStatus = ReceiptFailed
	default:
		return DeliveryReceipt{}, false
	}
	errorMessage := ""
	if record.Error != nil {
		errorMessage = strings.TrimSpace(*record.Error)
	}
	occurredAt := record.UpdatedAt.UTC()
	return DeliveryReceipt{
		ReceiptID:    "task:" + firstNonBlank(taskID, traceID) + ":" + receiptStatus,
		TraceID:      traceID,
		ConnectorID:  firstNonBlank(options.ConnectorID, "default-send-connector"),
		Channel:      firstNonBlank(options.Channel, ChannelInternalWebhook),
		TenantID:     firstNonBlank(options.TenantID, "default"),
		MessageID:    firstNonBlank(traceID, taskID),
		Status:       receiptStatus,
		ErrorMessage: errorMessage,
		OccurredAt:   occurredAt,
		Metadata: map[string]any{
			"task_id":     taskID,
			"task_type":   strings.TrimSpace(record.TaskType),
			"task_status": status,
		},
	}, true
}

// OutgoingDeliveryUpdateFromReceipt maps connector receipts to the message delivery update shape.
func OutgoingDeliveryUpdateFromReceipt(receipt DeliveryReceipt) (tasks.OutgoingDeliveryUpdate, bool) {
	update := tasks.OutgoingDeliveryUpdate{
		TraceID: strings.TrimSpace(receipt.TraceID),
		TaskID:  metadataText(receipt.Metadata, "task_id"),
	}
	if update.TraceID == "" && update.TaskID == "" {
		update.TraceID = strings.TrimSpace(receipt.MessageID)
	}
	switch strings.TrimSpace(receipt.Status) {
	case ReceiptSent, ReceiptDelivered, ReceiptRead:
		update.SendStatus = "success"
	case ReceiptFailed, ReceiptRevoked:
		update.SendStatus = "failed"
		update.SendError = firstNonBlank(receipt.ErrorMessage, receipt.ErrorCode)
	default:
		return tasks.OutgoingDeliveryUpdate{}, false
	}
	if update.TraceID == "" && update.TaskID == "" {
		return tasks.OutgoingDeliveryUpdate{}, false
	}
	return update, true
}

func messageTypeFromTask(taskType string) (string, bool) {
	switch strings.TrimSpace(taskType) {
	case "send_text":
		return MessageTypeText, true
	case "send_image":
		return MessageTypeImage, true
	case "send_video":
		return MessageTypeVideo, true
	case "send_voice":
		return MessageTypeVoice, true
	case "send_file":
		return MessageTypeFile, true
	default:
		return "", false
	}
}

func mediaFromTask(taskID string, messageType string, payload map[string]any) []MediaAttachment {
	url := payloadText(payload, "media_url")
	if url == "" {
		return nil
	}
	metadata := map[string]any{}
	if filename := payloadText(payload, "filename"); filename != "" {
		metadata["filename"] = filename
	}
	if duration := payload["voice_duration_sec"]; duration != nil {
		metadata["voice_duration_sec"] = duration
	}
	return []MediaAttachment{{
		AttachmentID: firstNonBlank(payloadText(payload, "msg_id"), taskID),
		Type:         messageType,
		URL:          url,
		MIMEType:     payloadText(payload, "media_mime"),
		Metadata:     metadata,
	}}
}

func idempotencyKey(record tasks.Record, payload map[string]any) string {
	if batchID := payloadText(payload, "client_batch_id"); batchID != "" {
		if batchIndex := payload["client_batch_index"]; batchIndex != nil {
			return fmt.Sprintf("%s:%v", batchID, batchIndex)
		}
		return batchID
	}
	return firstNonBlank(stringPointerValue(record.TraceID), strings.TrimSpace(record.TaskID))
}

func payloadText(payload map[string]any, key string) string {
	value, ok := payload[key]
	if !ok {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case []byte:
		return strings.TrimSpace(string(typed))
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func metadataText(metadata map[string]any, key string) string {
	if metadata == nil {
		return ""
	}
	return payloadText(metadata, key)
}

func stringPointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
