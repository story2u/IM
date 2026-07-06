package service_test

import (
	"context"
	"testing"
	"time"

	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/service"
	"im-go/internal/integrationhub/store"
)

func TestSendConversationMessageQueuesOutboxAndAudit(t *testing.T) {
	repo := store.NewMemory(domain.SeedSnapshot(time.Date(2026, 7, 6, 9, 0, 0, 0, time.UTC)))
	svc := service.New(repo)
	svc.SetClock(func() time.Time { return time.Date(2026, 7, 6, 10, 0, 0, 0, time.UTC) })
	var next int
	svc.SetIDGenerator(func(prefix string) string {
		next++
		return prefix + "_test_" + string(rune('a'+next))
	})

	item, err := svc.SendConversationMessage(context.Background(), "conv_101", service.SendMessageInput{
		Content: "Please review the updated quote.",
		Sender:  "Sarah Chen",
	})
	if err != nil {
		t.Fatalf("send message: %v", err)
	}
	if item.Status != domain.OutboxPending || item.ConversationID != "conv_101" || item.DeliveryMethod != domain.SendAPI {
		t.Fatalf("unexpected outbox item: %+v", item)
	}
	messages, err := svc.ConversationMessages(context.Background(), "conv_101")
	if err != nil {
		t.Fatalf("messages: %v", err)
	}
	if got := messages[len(messages)-1].Content; got != "Please review the updated quote." {
		t.Fatalf("last message = %q", got)
	}
	audit, err := svc.AuditLog(context.Background(), service.AuditFilter{Query: item.ID})
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	if len(audit) == 0 || audit[0].Action != "Queued outbound message" {
		t.Fatalf("missing queued audit entry: %+v", audit)
	}
}

func TestCreateChannelValidatesKind(t *testing.T) {
	svc := service.New(store.NewMemory(domain.SeedSnapshot(time.Now().UTC())))
	_, err := svc.CreateChannel(context.Background(), service.CreateChannelInput{
		Kind: "unsupported",
		Name: "Unsupported",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
