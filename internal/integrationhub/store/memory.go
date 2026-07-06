package store

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/service"
)

type MemoryStore struct {
	mu       sync.RWMutex
	snapshot domain.Snapshot
}

func NewMemory(snapshot domain.Snapshot) *MemoryStore {
	return &MemoryStore{snapshot: snapshot}
}

func (s *MemoryStore) Health(context.Context) error {
	return nil
}

func (s *MemoryStore) Overview(context.Context) (domain.OverviewData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return domain.OverviewData{
		Stats:     s.overviewLocked(),
		Channels:  cloneSlice(s.snapshot.Channels),
		Incidents: cloneSlice(s.snapshot.Incidents),
		Traffic:   cloneSlice(s.snapshot.TrafficSeries),
	}, nil
}

func (s *MemoryStore) ListChannels(context.Context) ([]domain.Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneSlice(s.snapshot.Channels), nil
}

func (s *MemoryStore) CreateChannel(_ context.Context, channel domain.Channel) (domain.Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshot.Channels = append(s.snapshot.Channels, channel)
	s.snapshot.AuditLog = append([]domain.AuditLogEntry{{
		ID: "aud_" + channel.ID, Time: channel.CreatedAt, Actor: "System", ActorType: domain.ActorSystem,
		Action: "Created channel", Target: channel.ID, Channel: &channel.Kind, Result: domain.AuditSuccess,
	}}, s.snapshot.AuditLog...)
	return channel, nil
}

func (s *MemoryStore) TouchChannel(_ context.Context, id string, now time.Time, audit domain.AuditLogEntry) (domain.Channel, error) {
	return s.updateChannel(id, "", now, audit)
}

func (s *MemoryStore) UpdateChannelStatus(_ context.Context, id string, status domain.ChannelStatus, now time.Time, audit domain.AuditLogEntry) (domain.Channel, error) {
	return s.updateChannel(id, status, now, audit)
}

func (s *MemoryStore) updateChannel(id string, status domain.ChannelStatus, now time.Time, audit domain.AuditLogEntry) (domain.Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.snapshot.Channels {
		if s.snapshot.Channels[i].ID == id {
			if status != "" {
				s.snapshot.Channels[i].Status = status
			}
			s.snapshot.Channels[i].LastSyncAt = now
			s.snapshot.Channels[i].UpdatedAt = now
			audit.Channel = &s.snapshot.Channels[i].Kind
			s.appendAuditLocked(audit)
			return s.snapshot.Channels[i], nil
		}
	}
	return domain.Channel{}, service.ErrNotFound
}

func (s *MemoryStore) MessageFlow(_ context.Context, filter service.MessageEventFilter) ([]domain.PipelineStageStats, []domain.MessageEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]domain.MessageEvent, 0, len(s.snapshot.MessageEvents))
	for _, event := range s.snapshot.MessageEvents {
		if filter.Channel != "" && filter.Channel != "all" && string(event.Channel) != filter.Channel {
			continue
		}
		if filter.Status != "" && filter.Status != "all" && string(event.Status) != filter.Status {
			continue
		}
		if filter.EventType != "" && filter.EventType != "all" && event.EventType != filter.EventType {
			continue
		}
		if filter.TraceID != "" && !strings.Contains(strings.ToLower(event.TraceID), strings.ToLower(filter.TraceID)) {
			continue
		}
		events = append(events, event)
	}
	return cloneSlice(s.snapshot.PipelineStats), events, nil
}

func (s *MemoryStore) Conversations(_ context.Context, filter service.ConversationFilter) ([]domain.Conversation, []domain.ConversationMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conversations := make([]domain.Conversation, 0, len(s.snapshot.Conversations))
	for _, conversation := range s.snapshot.Conversations {
		if filter.Channel != "" && filter.Channel != "all" && string(conversation.Channel) != filter.Channel {
			continue
		}
		if filter.Query != "" && !strings.Contains(strings.ToLower(conversation.ContactName+" "+conversation.ContactHandle), strings.ToLower(filter.Query)) {
			continue
		}
		conversations = append(conversations, conversation)
	}
	return conversations, cloneSlice(s.snapshot.Messages), nil
}

func (s *MemoryStore) GetConversation(_ context.Context, id string) (domain.Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, conversation := range s.snapshot.Conversations {
		if conversation.ID == id {
			return conversation, nil
		}
	}
	return domain.Conversation{}, service.ErrNotFound
}

func (s *MemoryStore) ConversationMessages(_ context.Context, conversationID string) ([]domain.ConversationMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.hasConversationLocked(conversationID) {
		return nil, service.ErrNotFound
	}
	messages := make([]domain.ConversationMessage, 0)
	for _, message := range s.snapshot.Messages {
		if message.ConversationID == conversationID {
			messages = append(messages, message)
		}
	}
	return messages, nil
}

func (s *MemoryStore) QueueOutboundMessage(_ context.Context, message domain.ConversationMessage, item domain.OutboxItem, event domain.MessageEvent, audit domain.AuditLogEntry) (domain.OutboxItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.snapshot.Conversations {
		if s.snapshot.Conversations[i].ID == message.ConversationID {
			s.snapshot.Conversations[i].LastMessagePreview = message.Content
			s.snapshot.Conversations[i].LastMessageAt = message.Time
			s.snapshot.Conversations[i].UpdatedAt = message.Time
			break
		}
	}
	s.snapshot.Messages = append(s.snapshot.Messages, message)
	s.snapshot.OutboxItems = append([]domain.OutboxItem{item}, s.snapshot.OutboxItems...)
	s.snapshot.MessageEvents = append([]domain.MessageEvent{event}, s.snapshot.MessageEvents...)
	s.appendAuditLocked(audit)
	return item, nil
}

func (s *MemoryStore) ListAIPolicies(context.Context) ([]domain.AIPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneSlice(s.snapshot.AIPolicies), nil
}

func (s *MemoryStore) UpdateAIPolicyEnabled(_ context.Context, id string, enabled bool, now time.Time, audit domain.AuditLogEntry) (domain.AIPolicy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.snapshot.AIPolicies {
		if s.snapshot.AIPolicies[i].ID == id {
			s.snapshot.AIPolicies[i].Enabled = enabled
			s.snapshot.AIPolicies[i].UpdatedAt = now
			s.appendAuditLocked(audit)
			return s.snapshot.AIPolicies[i], nil
		}
	}
	return domain.AIPolicy{}, service.ErrNotFound
}

func (s *MemoryStore) ListSOPWorkflows(context.Context) ([]domain.SOPWorkflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneSlice(s.snapshot.SOPWorkflows), nil
}

func (s *MemoryStore) ListOutbox(_ context.Context, status string) ([]domain.OutboxItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]domain.OutboxItem, 0, len(s.snapshot.OutboxItems))
	for _, item := range s.snapshot.OutboxItems {
		if item.Status == domain.OutboxCanceled {
			continue
		}
		if status != "" && status != "all" && string(item.Status) != status {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *MemoryStore) MoveOutbox(_ context.Context, id string, status domain.OutboxStatus, incrementRetry bool, now time.Time, audit domain.AuditLogEntry) (domain.OutboxItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.snapshot.OutboxItems {
		if s.snapshot.OutboxItems[i].ID == id {
			item := &s.snapshot.OutboxItems[i]
			if err := validateOutboxTransition(item.Status, status, incrementRetry); err != nil {
				return domain.OutboxItem{}, err
			}
			item.Status = status
			item.UpdatedAt = now
			if incrementRetry {
				item.RetryCount++
			}
			if status == domain.OutboxSending {
				item.LastError = nil
			}
			audit.Channel = &item.Channel
			s.appendAuditLocked(audit)
			return *item, nil
		}
	}
	return domain.OutboxItem{}, service.ErrNotFound
}

func (s *MemoryStore) Observability(context.Context) (domain.ObservabilityData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return domain.ObservabilityData{
		Channels:      cloneSlice(s.snapshot.Channels),
		MessageEvents: cloneSlice(s.snapshot.MessageEvents),
		Traffic:       cloneSlice(s.snapshot.TrafficSeries),
		Stats:         s.overviewLocked(),
	}, nil
}

func (s *MemoryStore) AuditLog(_ context.Context, filter service.AuditFilter) ([]domain.AuditLogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries := make([]domain.AuditLogEntry, 0, len(s.snapshot.AuditLog))
	for _, entry := range s.snapshot.AuditLog {
		if filter.ActorType != "" && filter.ActorType != "all" && string(entry.ActorType) != filter.ActorType {
			continue
		}
		if filter.Query != "" && !strings.Contains(strings.ToLower(entry.Action+" "+entry.Target), strings.ToLower(filter.Query)) {
			continue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *MemoryStore) Settings(context.Context) (domain.PlatformSettings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshot.Settings, nil
}

func (s *MemoryStore) UpdateSettings(_ context.Context, settings domain.PlatformSettings, _ time.Time, audit domain.AuditLogEntry) (domain.PlatformSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if settings.WorkspaceName != "" {
		s.snapshot.Settings.WorkspaceName = settings.WorkspaceName
	}
	if settings.Timezone != "" {
		s.snapshot.Settings.Timezone = settings.Timezone
	}
	if settings.DefaultLanguage != "" {
		s.snapshot.Settings.DefaultLanguage = settings.DefaultLanguage
	}
	if settings.Environment != "" {
		s.snapshot.Settings.Environment = settings.Environment
	}
	if settings.Region != "" {
		s.snapshot.Settings.Region = settings.Region
	}
	if settings.RetentionDays > 0 {
		s.snapshot.Settings.RetentionDays = settings.RetentionDays
	}
	if settings.WebhookURL != "" {
		s.snapshot.Settings.WebhookURL = settings.WebhookURL
	}
	if settings.EnabledProviders != nil {
		s.snapshot.Settings.EnabledProviders = settings.EnabledProviders
	}
	s.appendAuditLocked(audit)
	return s.snapshot.Settings, nil
}

func (s *MemoryStore) RecordInboundEvent(_ context.Context, event domain.InboundEvent, audit domain.AuditLogEntry) (domain.InboundEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshot.InboundEvents = append([]domain.InboundEvent{event}, s.snapshot.InboundEvents...)
	s.snapshot.MessageEvents = append([]domain.MessageEvent{{
		ID:                "evt_" + event.ID,
		Time:              event.ReceivedAt,
		Channel:           event.ConnectorKind,
		Direction:         domain.DirectionInbound,
		ConversationID:    fmt.Sprint(event.Normalized["conversationId"]),
		ConversationLabel: fmt.Sprint(event.Normalized["conversationLabel"]),
		EventType:         event.EventType,
		Status:            event.Status,
		TraceID:           event.TraceID,
	}}, s.snapshot.MessageEvents...)
	s.appendAuditLocked(audit)
	return event, nil
}

func (s *MemoryStore) overviewLocked() domain.OverviewStats {
	activeChannels := 0
	messages := 0
	pending := 0
	for _, channel := range s.snapshot.Channels {
		if channel.Status != domain.ChannelDisabled {
			activeChannels++
		}
		messages += channel.MessagesToday
	}
	for _, item := range s.snapshot.OutboxItems {
		if item.Status == domain.OutboxPending || item.Status == domain.OutboxRequiresApproval {
			pending++
		}
	}
	return domain.OverviewStats{
		ActiveChannels: activeChannels, TotalChannels: len(s.snapshot.Channels),
		MessagesIngestedToday: messages, AIActionsToday: 6420, OutboxPending: pending,
		ErrorRate: 0.021, P95LatencyMs: 1180,
	}
}

func (s *MemoryStore) appendAuditLocked(entry domain.AuditLogEntry) {
	s.snapshot.AuditLog = append([]domain.AuditLogEntry{entry}, s.snapshot.AuditLog...)
}

func (s *MemoryStore) hasConversationLocked(id string) bool {
	return slices.ContainsFunc(s.snapshot.Conversations, func(conversation domain.Conversation) bool {
		return conversation.ID == id
	})
}

func cloneSlice[T any](in []T) []T {
	return append([]T(nil), in...)
}
