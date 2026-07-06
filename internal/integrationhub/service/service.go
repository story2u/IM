package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"im-go/internal/integrationhub/domain"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation failed")
)

type Repository interface {
	Health(ctx context.Context) error
	Overview(ctx context.Context) (domain.OverviewData, error)
	ListChannels(ctx context.Context) ([]domain.Channel, error)
	CreateChannel(ctx context.Context, channel domain.Channel) (domain.Channel, error)
	TouchChannel(ctx context.Context, id string, now time.Time, audit domain.AuditLogEntry) (domain.Channel, error)
	UpdateChannelStatus(ctx context.Context, id string, status domain.ChannelStatus, now time.Time, audit domain.AuditLogEntry) (domain.Channel, error)
	MessageFlow(ctx context.Context, filter MessageEventFilter) ([]domain.PipelineStageStats, []domain.MessageEvent, error)
	Conversations(ctx context.Context, filter ConversationFilter) ([]domain.Conversation, []domain.ConversationMessage, error)
	GetConversation(ctx context.Context, id string) (domain.Conversation, error)
	ConversationMessages(ctx context.Context, conversationID string) ([]domain.ConversationMessage, error)
	QueueOutboundMessage(ctx context.Context, message domain.ConversationMessage, item domain.OutboxItem, event domain.MessageEvent, audit domain.AuditLogEntry) (domain.OutboxItem, error)
	ListAIPolicies(ctx context.Context) ([]domain.AIPolicy, error)
	UpdateAIPolicyEnabled(ctx context.Context, id string, enabled bool, now time.Time, audit domain.AuditLogEntry) (domain.AIPolicy, error)
	ListSOPWorkflows(ctx context.Context) ([]domain.SOPWorkflow, error)
	ListOutbox(ctx context.Context, status string) ([]domain.OutboxItem, error)
	MoveOutbox(ctx context.Context, id string, status domain.OutboxStatus, incrementRetry bool, now time.Time, audit domain.AuditLogEntry) (domain.OutboxItem, error)
	Observability(ctx context.Context) (domain.ObservabilityData, error)
	AuditLog(ctx context.Context, filter AuditFilter) ([]domain.AuditLogEntry, error)
	Settings(ctx context.Context) (domain.PlatformSettings, error)
	UpdateSettings(ctx context.Context, settings domain.PlatformSettings, now time.Time, audit domain.AuditLogEntry) (domain.PlatformSettings, error)
	RecordInboundEvent(ctx context.Context, event domain.InboundEvent, audit domain.AuditLogEntry) (domain.InboundEvent, error)
}

type Clock func() time.Time
type IDGenerator func(prefix string) string

type Service struct {
	repo Repository
	now  Clock
	id   IDGenerator
}

func New(repo Repository) *Service {
	return &Service{repo: repo, now: func() time.Time { return time.Now().UTC() }, id: randomID}
}

func (s *Service) SetClock(now Clock) {
	if now != nil {
		s.now = now
	}
}

func (s *Service) SetIDGenerator(id IDGenerator) {
	if id != nil {
		s.id = id
	}
}

func (s *Service) Health(ctx context.Context) error {
	return s.repo.Health(ctx)
}

func (s *Service) Overview(ctx context.Context) (domain.OverviewData, error) {
	return s.repo.Overview(ctx)
}

func (s *Service) Channels(ctx context.Context) ([]domain.Channel, error) {
	return s.repo.ListChannels(ctx)
}

type CreateChannelInput struct {
	Kind                domain.ChannelKind         `json:"kind"`
	Name                string                     `json:"name"`
	Status              domain.ChannelStatus       `json:"status"`
	ReceiveCapabilities []domain.ReceiveCapability `json:"receiveCapabilities"`
	SendCapabilities    []domain.SendCapability    `json:"sendCapabilities"`
}

func (s *Service) CreateChannel(ctx context.Context, input CreateChannelInput) (domain.Channel, error) {
	kind := input.Kind
	if !validChannelKind(kind) {
		return domain.Channel{}, validation("unsupported channel kind")
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.Channel{}, validation("name is required")
	}
	status := input.Status
	if status == "" {
		status = domain.ChannelConnected
	}
	if !validChannelStatus(status) {
		return domain.Channel{}, validation("unsupported channel status")
	}
	receive := input.ReceiveCapabilities
	if len(receive) == 0 {
		receive = []domain.ReceiveCapability{domain.ReceiveWebhook}
	}
	for _, capability := range receive {
		if !validReceiveCapability(capability) {
			return domain.Channel{}, validation("unsupported receive capability")
		}
	}
	send := input.SendCapabilities
	if len(send) == 0 {
		send = []domain.SendCapability{domain.SendAPI}
	}
	for _, capability := range send {
		if !validSendCapability(capability) {
			return domain.Channel{}, validation("unsupported send capability")
		}
	}
	now := s.now().UTC()
	connectorID := s.id("conn")
	channel := domain.Channel{
		ID:                  s.id("ch"),
		Kind:                kind,
		Name:                name,
		Status:              status,
		ReceiveCapabilities: receive,
		SendCapabilities:    send,
		LastSyncAt:          now,
		ConnectorID:         &connectorID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	return s.repo.CreateChannel(ctx, channel)
}

func (s *Service) TestChannel(ctx context.Context, id string) (domain.Channel, error) {
	now := s.now().UTC()
	return s.repo.TouchChannel(ctx, id, now, s.audit("System", domain.ActorSystem, "Tested channel connection", id, nil, domain.AuditSuccess))
}

func (s *Service) SetChannelStatus(ctx context.Context, id string, status domain.ChannelStatus) (domain.Channel, error) {
	if !validChannelStatus(status) {
		return domain.Channel{}, validation("unsupported channel status")
	}
	now := s.now().UTC()
	action := fmt.Sprintf("Set channel status to %s", status)
	return s.repo.UpdateChannelStatus(ctx, id, status, now, s.audit("System", domain.ActorSystem, action, id, nil, domain.AuditSuccess))
}

type MessageEventFilter struct {
	Channel   string
	Status    string
	EventType string
	TraceID   string
}

func (s *Service) MessageFlow(ctx context.Context, filter MessageEventFilter) ([]domain.PipelineStageStats, []domain.MessageEvent, error) {
	return s.repo.MessageFlow(ctx, filter)
}

type ConversationFilter struct {
	Channel string
	Query   string
}

func (s *Service) Conversations(ctx context.Context, filter ConversationFilter) ([]domain.Conversation, []domain.ConversationMessage, error) {
	return s.repo.Conversations(ctx, filter)
}

func (s *Service) ConversationMessages(ctx context.Context, conversationID string) ([]domain.ConversationMessage, error) {
	return s.repo.ConversationMessages(ctx, strings.TrimSpace(conversationID))
}

type SendMessageInput struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

func (s *Service) SendConversationMessage(ctx context.Context, conversationID string, input SendMessageInput) (domain.OutboxItem, error) {
	conversation, err := s.repo.GetConversation(ctx, strings.TrimSpace(conversationID))
	if err != nil {
		return domain.OutboxItem{}, err
	}
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return domain.OutboxItem{}, validation("content is required")
	}
	sender := strings.TrimSpace(input.Sender)
	if sender == "" {
		sender = "Operator"
	}
	now := s.now().UTC()
	messageID := s.id("msg")
	outboxID := s.id("out")
	traceID := s.id("trace")
	message := domain.ConversationMessage{
		ID:             messageID,
		ConversationID: conversation.ID,
		ChannelID:      conversation.ChannelID,
		Channel:        conversation.Channel,
		Direction:      domain.DirectionOutbound,
		Author:         sender,
		Content:        content,
		Time:           now,
		MessageType:    "text",
	}
	item := domain.OutboxItem{
		ID:                outboxID,
		CreatedAt:         now,
		UpdatedAt:         now,
		ChannelID:         conversation.ChannelID,
		Channel:           conversation.Channel,
		ConversationID:    conversation.ID,
		ConversationLabel: conversation.ContactName,
		MessageID:         &messageID,
		MessageType:       "Text",
		Sender:            sender,
		DeliveryMethod:    domain.SendAPI,
		Status:            domain.OutboxPending,
		IdempotencyKey:    s.id("idem"),
		Payload: map[string]any{
			"content": content,
			"type":    "text",
		},
	}
	event := domain.MessageEvent{
		ID:                s.id("evt"),
		Time:              now,
		Channel:           conversation.Channel,
		Direction:         domain.DirectionOutbound,
		ConversationID:    conversation.ID,
		ConversationLabel: conversation.ContactName,
		EventType:         "message.queued",
		Status:            domain.EventPending,
		TraceID:           traceID,
	}
	audit := s.audit(sender, domain.ActorUser, "Queued outbound message", outboxID, &conversation.Channel, domain.AuditSuccess)
	audit.TraceID = traceID
	return s.repo.QueueOutboundMessage(ctx, message, item, event, audit)
}

func (s *Service) AIPolicies(ctx context.Context) ([]domain.AIPolicy, error) {
	return s.repo.ListAIPolicies(ctx)
}

func (s *Service) SetAIPolicyEnabled(ctx context.Context, id string, enabled bool) (domain.AIPolicy, error) {
	return s.repo.UpdateAIPolicyEnabled(ctx, strings.TrimSpace(id), enabled, s.now().UTC(), s.audit("System", domain.ActorSystem, "Updated AI policy", id, nil, domain.AuditSuccess))
}

func (s *Service) SOPWorkflows(ctx context.Context) ([]domain.SOPWorkflow, error) {
	return s.repo.ListSOPWorkflows(ctx)
}

func (s *Service) Outbox(ctx context.Context, status string) ([]domain.OutboxItem, error) {
	status = strings.TrimSpace(status)
	if status != "" && status != "all" && !validOutboxStatus(domain.OutboxStatus(status)) {
		return nil, validation("unsupported outbox status")
	}
	return s.repo.ListOutbox(ctx, status)
}

func (s *Service) RetryOutbox(ctx context.Context, id string) (domain.OutboxItem, error) {
	return s.repo.MoveOutbox(ctx, strings.TrimSpace(id), domain.OutboxSending, true, s.now().UTC(), s.audit("System", domain.ActorSystem, "Retried outbox message", id, nil, domain.AuditSuccess))
}

func (s *Service) ApproveOutbox(ctx context.Context, id string) (domain.OutboxItem, error) {
	return s.repo.MoveOutbox(ctx, strings.TrimSpace(id), domain.OutboxSending, false, s.now().UTC(), s.audit("System", domain.ActorSystem, "Approved outbox message", id, nil, domain.AuditSuccess))
}

func (s *Service) CancelOutbox(ctx context.Context, id string) (domain.OutboxItem, error) {
	return s.repo.MoveOutbox(ctx, strings.TrimSpace(id), domain.OutboxCanceled, false, s.now().UTC(), s.audit("System", domain.ActorSystem, "Canceled outbox message", id, nil, domain.AuditSuccess))
}

func (s *Service) Observability(ctx context.Context) (domain.ObservabilityData, error) {
	return s.repo.Observability(ctx)
}

type AuditFilter struct {
	ActorType string
	Query     string
}

func (s *Service) AuditLog(ctx context.Context, filter AuditFilter) ([]domain.AuditLogEntry, error) {
	return s.repo.AuditLog(ctx, filter)
}

func (s *Service) Settings(ctx context.Context) (domain.PlatformSettings, error) {
	return s.repo.Settings(ctx)
}

func (s *Service) UpdateSettings(ctx context.Context, settings domain.PlatformSettings) (domain.PlatformSettings, error) {
	settings.WorkspaceName = strings.TrimSpace(settings.WorkspaceName)
	settings.Timezone = strings.TrimSpace(settings.Timezone)
	settings.DefaultLanguage = strings.TrimSpace(settings.DefaultLanguage)
	settings.Environment = strings.TrimSpace(settings.Environment)
	settings.Region = strings.TrimSpace(settings.Region)
	settings.WebhookURL = strings.TrimSpace(settings.WebhookURL)
	return s.repo.UpdateSettings(ctx, settings, s.now().UTC(), s.audit("System", domain.ActorSystem, "Updated platform settings", "settings", nil, domain.AuditSuccess))
}

func (s *Service) RecordInboundEvent(ctx context.Context, event domain.InboundEvent) (domain.InboundEvent, error) {
	if !validChannelKind(event.ConnectorKind) {
		return domain.InboundEvent{}, validation("unsupported connector kind")
	}
	if strings.TrimSpace(event.EventType) == "" {
		return domain.InboundEvent{}, validation("eventType is required")
	}
	now := s.now().UTC()
	if event.ID == "" {
		event.ID = s.id("in")
	}
	if event.ReceivedAt.IsZero() {
		event.ReceivedAt = now
	}
	if event.TraceID == "" {
		event.TraceID = s.id("trace")
	}
	if event.Status == "" {
		event.Status = domain.EventPending
	}
	audit := s.audit("System", domain.ActorSystem, "Recorded inbound connector event", event.ID, &event.ConnectorKind, domain.AuditSuccess)
	audit.TraceID = event.TraceID
	return s.repo.RecordInboundEvent(ctx, event, audit)
}

func (s *Service) audit(actor string, actorType domain.AuditActorType, action, target string, channel *domain.ChannelKind, result domain.AuditResult) domain.AuditLogEntry {
	return domain.AuditLogEntry{
		ID:        s.id("aud"),
		Time:      s.now().UTC(),
		Actor:     actor,
		ActorType: actorType,
		Action:    action,
		Target:    target,
		Channel:   channel,
		Result:    result,
	}
}

func validation(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}

func validChannelKind(kind domain.ChannelKind) bool {
	switch kind {
	case domain.ChannelWeCom, domain.ChannelFeishu, domain.ChannelDingTalk, domain.ChannelWhatsApp, domain.ChannelTelegram, domain.ChannelEmail:
		return true
	default:
		return false
	}
}

func validChannelStatus(status domain.ChannelStatus) bool {
	switch status {
	case domain.ChannelConnected, domain.ChannelDegraded, domain.ChannelDisabled:
		return true
	default:
		return false
	}
}

func validReceiveCapability(capability domain.ReceiveCapability) bool {
	switch capability {
	case domain.ReceiveWebhook, domain.ReceivePolling, domain.ReceiveRPA:
		return true
	default:
		return false
	}
}

func validSendCapability(capability domain.SendCapability) bool {
	switch capability {
	case domain.SendAPI, domain.SendRPA, domain.SendManualApproval:
		return true
	default:
		return false
	}
}

func validOutboxStatus(status domain.OutboxStatus) bool {
	switch status {
	case domain.OutboxPending, domain.OutboxSending, domain.OutboxFailed, domain.OutboxSent, domain.OutboxRequiresApproval, domain.OutboxCanceled:
		return true
	default:
		return false
	}
}

func randomID(prefix string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}
