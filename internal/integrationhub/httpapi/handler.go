package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"im-go/internal/integrationhub/connectors"
	"im-go/internal/integrationhub/domain"
	"im-go/internal/integrationhub/service"
)

type Handler struct {
	service *service.Service
	mux     *http.ServeMux
}

func New(service *service.Service) http.Handler {
	h := &Handler{service: service, mux: http.NewServeMux()}
	h.routes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) routes() {
	h.mux.HandleFunc("GET /healthz", h.health)
	h.mux.HandleFunc("GET /api/v1/health", h.health)
	h.mux.HandleFunc("GET /api/v1/overview", h.overview)
	h.mux.HandleFunc("GET /api/v1/channels", h.channels)
	h.mux.HandleFunc("POST /api/v1/channels", h.createChannel)
	h.mux.HandleFunc("POST /api/v1/channels/{id}/test", h.testChannel)
	h.mux.HandleFunc("POST /api/v1/channels/{id}/disable", h.disableChannel)
	h.mux.HandleFunc("POST /api/v1/channels/{id}/enable", h.enableChannel)
	h.mux.HandleFunc("GET /api/v1/message-flow", h.messageFlow)
	h.mux.HandleFunc("GET /api/v1/conversations", h.conversations)
	h.mux.HandleFunc("GET /api/v1/conversations/{id}/messages", h.conversationMessages)
	h.mux.HandleFunc("POST /api/v1/conversations/{id}/messages", h.sendConversationMessage)
	h.mux.HandleFunc("GET /api/v1/ai/policies", h.aiPolicies)
	h.mux.HandleFunc("PATCH /api/v1/ai/policies/{id}", h.updateAIPolicy)
	h.mux.HandleFunc("GET /api/v1/sop/workflows", h.sopWorkflows)
	h.mux.HandleFunc("GET /api/v1/outbox", h.outbox)
	h.mux.HandleFunc("POST /api/v1/outbox/{id}/retry", h.retryOutbox)
	h.mux.HandleFunc("POST /api/v1/outbox/{id}/approve", h.approveOutbox)
	h.mux.HandleFunc("POST /api/v1/outbox/{id}/cancel", h.cancelOutbox)
	h.mux.HandleFunc("GET /api/v1/observability", h.observability)
	h.mux.HandleFunc("GET /api/v1/audit-logs", h.auditLog)
	h.mux.HandleFunc("GET /api/v1/settings", h.settings)
	h.mux.HandleFunc("PATCH /api/v1/settings", h.updateSettings)
	h.mux.HandleFunc("POST /api/v1/connectors/wecom/callback", h.wecomCallback)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Health(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"ok":      false,
			"service": "im-integration-api",
			"error":   err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"service": "im-integration-api",
	})
}

func (h *Handler) overview(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.Overview(r.Context())
	writeResult(w, map[string]any{
		"overviewStats":   data.Stats,
		"channels":        data.Channels,
		"recentIncidents": data.Incidents,
		"trafficSeries":   data.Traffic,
	}, err)
}

func (h *Handler) channels(w http.ResponseWriter, r *http.Request) {
	channels, err := h.service.Channels(r.Context())
	writeResult(w, map[string]any{"channels": channels}, err)
}

func (h *Handler) createChannel(w http.ResponseWriter, r *http.Request) {
	var input service.CreateChannelInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	channel, err := h.service.CreateChannel(r.Context(), input)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"channel": channel})
}

func (h *Handler) testChannel(w http.ResponseWriter, r *http.Request) {
	channel, err := h.service.TestChannel(r.Context(), r.PathValue("id"))
	writeResult(w, map[string]any{"channel": channel, "ok": true}, err)
}

func (h *Handler) disableChannel(w http.ResponseWriter, r *http.Request) {
	channel, err := h.service.SetChannelStatus(r.Context(), r.PathValue("id"), domain.ChannelDisabled)
	writeResult(w, map[string]any{"channel": channel}, err)
}

func (h *Handler) enableChannel(w http.ResponseWriter, r *http.Request) {
	channel, err := h.service.SetChannelStatus(r.Context(), r.PathValue("id"), domain.ChannelConnected)
	writeResult(w, map[string]any{"channel": channel}, err)
}

func (h *Handler) messageFlow(w http.ResponseWriter, r *http.Request) {
	stats, events, err := h.service.MessageFlow(r.Context(), service.MessageEventFilter{
		Channel:   strings.TrimSpace(r.URL.Query().Get("channel")),
		Status:    strings.TrimSpace(r.URL.Query().Get("status")),
		EventType: strings.TrimSpace(r.URL.Query().Get("eventType")),
		TraceID:   strings.TrimSpace(r.URL.Query().Get("traceId")),
	})
	writeResult(w, map[string]any{"pipelineStats": stats, "messageEvents": events}, err)
}

func (h *Handler) conversations(w http.ResponseWriter, r *http.Request) {
	conversations, messages, err := h.service.Conversations(r.Context(), service.ConversationFilter{
		Channel: strings.TrimSpace(r.URL.Query().Get("channel")),
		Query:   strings.TrimSpace(r.URL.Query().Get("q")),
	})
	writeResult(w, map[string]any{"conversations": conversations, "messages": messages}, err)
}

func (h *Handler) conversationMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.service.ConversationMessages(r.Context(), r.PathValue("id"))
	writeResult(w, map[string]any{"messages": messages}, err)
}

func (h *Handler) sendConversationMessage(w http.ResponseWriter, r *http.Request) {
	var input service.SendMessageInput
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.service.SendConversationMessage(r.Context(), r.PathValue("id"), input)
	writeResult(w, map[string]any{"outboxItem": item}, err)
}

func (h *Handler) aiPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := h.service.AIPolicies(r.Context())
	writeResult(w, map[string]any{"aiPolicies": policies}, err)
}

func (h *Handler) updateAIPolicy(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Enabled bool `json:"enabled"`
	}
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	policy, err := h.service.SetAIPolicyEnabled(r.Context(), r.PathValue("id"), input.Enabled)
	writeResult(w, map[string]any{"aiPolicy": policy}, err)
}

func (h *Handler) sopWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := h.service.SOPWorkflows(r.Context())
	writeResult(w, map[string]any{"sopWorkflows": workflows}, err)
}

func (h *Handler) outbox(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.Outbox(r.Context(), r.URL.Query().Get("status"))
	writeResult(w, map[string]any{"outboxItems": items}, err)
}

func (h *Handler) retryOutbox(w http.ResponseWriter, r *http.Request) {
	item, err := h.service.RetryOutbox(r.Context(), r.PathValue("id"))
	writeResult(w, map[string]any{"outboxItem": item}, err)
}

func (h *Handler) approveOutbox(w http.ResponseWriter, r *http.Request) {
	item, err := h.service.ApproveOutbox(r.Context(), r.PathValue("id"))
	writeResult(w, map[string]any{"outboxItem": item}, err)
}

func (h *Handler) cancelOutbox(w http.ResponseWriter, r *http.Request) {
	item, err := h.service.CancelOutbox(r.Context(), r.PathValue("id"))
	writeResult(w, map[string]any{"outboxItem": item}, err)
}

func (h *Handler) observability(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.Observability(r.Context())
	writeResult(w, map[string]any{
		"channels":      data.Channels,
		"messageEvents": data.MessageEvents,
		"trafficSeries": data.Traffic,
		"overviewStats": data.Stats,
	}, err)
}

func (h *Handler) auditLog(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.AuditLog(r.Context(), service.AuditFilter{
		ActorType: strings.TrimSpace(r.URL.Query().Get("actorType")),
		Query:     strings.TrimSpace(r.URL.Query().Get("q")),
	})
	writeResult(w, map[string]any{"auditLog": entries}, err)
}

func (h *Handler) settings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.Settings(r.Context())
	writeResult(w, map[string]any{"settings": settings}, err)
}

func (h *Handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.PlatformSettings
	if err := readJSON(r, &settings); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, err := h.service.UpdateSettings(r.Context(), settings)
	writeResult(w, map[string]any{"settings": updated}, err)
}

func (h *Handler) wecomCallback(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	adapter := connectors.WeComAdapter{}
	event, err := adapter.ParseInbound(r.Context(), raw)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if channelID := strings.TrimSpace(r.URL.Query().Get("channelId")); channelID != "" {
		event.ChannelID = channelID
	}
	recorded, err := h.service.RecordInboundEvent(r.Context(), event)
	writeResult(w, map[string]any{"inboundEvent": recorded}, err)
}

func writeResult(w http.ResponseWriter, body map[string]any, err error) {
	if err == nil {
		if body == nil {
			body = map[string]any{}
		}
		writeJSON(w, http.StatusOK, body)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if errors.Is(err, service.ErrValidation) {
		writeError(w, http.StatusBadRequest, strings.TrimPrefix(err.Error(), service.ErrValidation.Error()+": "))
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}

func readJSON(r *http.Request, out any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{"message": message},
	})
}
