package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/service"
	"im-go/internal/integrationhub/store"
)

func TestWeComCallbackRecordsGenericInboundEvent(t *testing.T) {
	repo := store.NewMemory(domain.SeedSnapshot(time.Date(2026, 7, 6, 9, 0, 0, 0, time.UTC)))
	svc := service.New(repo)
	svc.SetIDGenerator(func(prefix string) string { return prefix + "_fixed" })
	handler := New(svc)

	payload := bytes.NewBufferString(`{
		"eventId":"wecom-event-1",
		"externalUserId":"external-1",
		"senderName":"Li Wei",
		"conversationId":"conv_external_1",
		"content":"hello",
		"msgType":"text",
		"createTime":1783332000
	}`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/connectors/wecom/callback?channelId=ch_wecom", payload)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", response.Code, response.Body.String())
	}
	var body struct {
		InboundEvent domain.InboundEvent `json:"inboundEvent"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.InboundEvent.ConnectorKind != domain.ChannelWeCom {
		t.Fatalf("connector kind = %q", body.InboundEvent.ConnectorKind)
	}
	if body.InboundEvent.AdapterPayload["provider"] != "wecom" {
		t.Fatalf("adapter payload should keep provider-specific payload: %+v", body.InboundEvent.AdapterPayload)
	}
	if body.InboundEvent.Normalized["conversationId"] != "conv_external_1" {
		t.Fatalf("normalized payload = %+v", body.InboundEvent.Normalized)
	}
}
