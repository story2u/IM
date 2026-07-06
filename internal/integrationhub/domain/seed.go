package domain

import (
	"fmt"
	"math"
	"time"
)

func SeedSnapshot(now time.Time) Snapshot {
	now = now.UTC()
	ago := func(minutes int) time.Time { return now.Add(-time.Duration(minutes) * time.Minute).UTC() }
	operatorSarah := "Sarah Chen"
	operatorWei := "Wei Zhao"
	quoteWorkflow := "Enterprise Quote Approval"
	shippingWorkflow := "Shipping Escalation"
	contractWorkflow := "Contract Signing"
	partnerWorkflow := "Partner Onboarding"
	refundWorkflow := "Refund Processing"
	orderWorkflow := "Order Confirmation"
	emailChannel := ChannelEmail
	whatsAppChannel := ChannelWhatsApp
	dingTalkChannel := ChannelDingTalk

	connectors := []Connector{
		{ID: "conn_wecom", Kind: ChannelWeCom, Name: "WeCom connector", Status: ChannelConnected, CreatedAt: ago(1440), UpdatedAt: ago(1)},
		{ID: "conn_feishu", Kind: ChannelFeishu, Name: "Feishu connector", Status: ChannelConnected, CreatedAt: ago(1440), UpdatedAt: ago(2)},
		{ID: "conn_dingtalk", Kind: ChannelDingTalk, Name: "DingTalk connector", Status: ChannelDegraded, CreatedAt: ago(1440), UpdatedAt: ago(18)},
		{ID: "conn_whatsapp", Kind: ChannelWhatsApp, Name: "WhatsApp connector", Status: ChannelConnected, CreatedAt: ago(1440), UpdatedAt: ago(1)},
		{ID: "conn_telegram", Kind: ChannelTelegram, Name: "Telegram connector", Status: ChannelConnected, CreatedAt: ago(1440), UpdatedAt: ago(3)},
		{ID: "conn_email", Kind: ChannelEmail, Name: "Email connector", Status: ChannelDisabled, CreatedAt: ago(1440), UpdatedAt: ago(540)},
	}

	channels := []Channel{
		{
			ID: "ch_wecom", Kind: ChannelWeCom, Name: "WeCom - Sales Workspace", Status: ChannelConnected,
			ReceiveCapabilities: []ReceiveCapability{ReceiveWebhook}, SendCapabilities: []SendCapability{SendAPI},
			LastSyncAt: ago(1), ErrorCount24h: 2, MessagesToday: 4218, ActiveConversations: 312,
			ConnectorID: ptr("conn_wecom"), CreatedAt: ago(1440), UpdatedAt: ago(1),
		},
		{
			ID: "ch_feishu", Kind: ChannelFeishu, Name: "Feishu - Customer Success", Status: ChannelConnected,
			ReceiveCapabilities: []ReceiveCapability{ReceiveWebhook}, SendCapabilities: []SendCapability{SendAPI},
			LastSyncAt: ago(2), ErrorCount24h: 0, MessagesToday: 2854, ActiveConversations: 201,
			ConnectorID: ptr("conn_feishu"), CreatedAt: ago(1440), UpdatedAt: ago(2),
		},
		{
			ID: "ch_dingtalk", Kind: ChannelDingTalk, Name: "DingTalk - Partner Channel", Status: ChannelDegraded,
			ReceiveCapabilities: []ReceiveCapability{ReceiveWebhook, ReceivePolling}, SendCapabilities: []SendCapability{SendAPI, SendManualApproval},
			LastSyncAt: ago(18), ErrorCount24h: 47, MessagesToday: 963, ActiveConversations: 84,
			ConnectorID: ptr("conn_dingtalk"), CreatedAt: ago(1440), UpdatedAt: ago(18),
		},
		{
			ID: "ch_whatsapp", Kind: ChannelWhatsApp, Name: "WhatsApp Business", Status: ChannelConnected,
			ReceiveCapabilities: []ReceiveCapability{ReceiveWebhook}, SendCapabilities: []SendCapability{SendAPI},
			LastSyncAt: ago(1), ErrorCount24h: 5, MessagesToday: 3109, ActiveConversations: 265,
			ConnectorID: ptr("conn_whatsapp"), CreatedAt: ago(1440), UpdatedAt: ago(1),
		},
		{
			ID: "ch_telegram", Kind: ChannelTelegram, Name: "Telegram Bot - Community", Status: ChannelConnected,
			ReceiveCapabilities: []ReceiveCapability{ReceiveWebhook}, SendCapabilities: []SendCapability{SendAPI},
			LastSyncAt: ago(3), ErrorCount24h: 1, MessagesToday: 1542, ActiveConversations: 118,
			ConnectorID: ptr("conn_telegram"), CreatedAt: ago(1440), UpdatedAt: ago(3),
		},
		{
			ID: "ch_email_rpa", Kind: ChannelEmail, Name: "Email - Support Mailbox (RPA)", Status: ChannelDisabled,
			ReceiveCapabilities: []ReceiveCapability{ReceivePolling, ReceiveRPA}, SendCapabilities: []SendCapability{SendRPA, SendManualApproval},
			LastSyncAt: ago(540), ErrorCount24h: 0, MessagesToday: 0, ActiveConversations: 0,
			ConnectorID: ptr("conn_email"), CreatedAt: ago(1440), UpdatedAt: ago(540),
		},
	}

	accounts := []ChannelAccount{
		{ID: "acct_wecom_sales", ChannelID: "ch_wecom", DisplayName: "Sales Workspace", ExternalAccountID: "sales-workspace", Status: ChannelConnected, CreatedAt: ago(1440), UpdatedAt: ago(1)},
		{ID: "acct_feishu_success", ChannelID: "ch_feishu", DisplayName: "Customer Success", ExternalAccountID: "customer-success", Status: ChannelConnected, CreatedAt: ago(1440), UpdatedAt: ago(2)},
		{ID: "acct_dingtalk_partner", ChannelID: "ch_dingtalk", DisplayName: "Partner Channel", ExternalAccountID: "partner-channel", Status: ChannelDegraded, CreatedAt: ago(1440), UpdatedAt: ago(18)},
		{ID: "acct_whatsapp_business", ChannelID: "ch_whatsapp", DisplayName: "WhatsApp Business", ExternalAccountID: "business-primary", Status: ChannelConnected, CreatedAt: ago(1440), UpdatedAt: ago(1)},
		{ID: "acct_telegram_community", ChannelID: "ch_telegram", DisplayName: "Community Bot", ExternalAccountID: "community-bot", Status: ChannelConnected, CreatedAt: ago(1440), UpdatedAt: ago(3)},
		{ID: "acct_email_support", ChannelID: "ch_email_rpa", DisplayName: "Support Mailbox", ExternalAccountID: "support@example.com", Status: ChannelDisabled, CreatedAt: ago(1440), UpdatedAt: ago(540)},
	}

	messageEvents := make([]MessageEvent, 0, 64)
	eventTypes := []string{"message.received", "message.sent", "message.delivered", "message.failed", "sop.step.completed", "ai.reply_drafted", "ai.handoff_triggered"}
	channelKinds := []ChannelKind{ChannelWeCom, ChannelFeishu, ChannelDingTalk, ChannelWhatsApp, ChannelTelegram, ChannelEmail}
	statuses := []MessageEventStatus{EventSuccess, EventSuccess, EventSuccess, EventSuccess, EventPending, EventRetrying, EventFailed}
	for i := 0; i < 64; i++ {
		channel := channelKinds[i%len(channelKinds)]
		messageEvents = append(messageEvents, MessageEvent{
			ID:                fmt.Sprintf("evt_%d", 1000+i),
			Time:              ago(3 + i*4),
			Channel:           channel,
			Direction:         []MessageDirection{DirectionInbound, DirectionOutbound}[i%2],
			ConversationID:    fmt.Sprintf("conv_%d", 100+i%40),
			ConversationLabel: fmt.Sprintf("Conversation #%d", 100+i%40),
			EventType:         eventTypes[i%len(eventTypes)],
			Status:            statuses[i%len(statuses)],
			LatencyMs:         40 + (i*37)%900,
			TraceID:           fmt.Sprintf("trace-%06x", 100000+i*7919),
		})
	}

	conversations := []Conversation{
		{ID: "conv_101", ChannelID: "ch_wecom", Channel: ChannelWeCom, ContactName: "Li Wei", ContactHandle: "@liwei_procurement", LastMessagePreview: "Can you send the latest quote and payment terms?", LastMessageAt: ago(4), AssignedOperator: &operatorSarah, AIStatus: AIAutoReplying, SOPStage: SOPInProgress, SOPWorkflowID: ptr("wf_1"), SOPWorkflowName: &quoteWorkflow, Unread: 2, Tags: []string{"enterprise", "quote"}, CreatedAt: ago(60), UpdatedAt: ago(4)},
		{ID: "conv_102", ChannelID: "ch_whatsapp", Channel: ChannelWhatsApp, ContactName: "Maria Gonzalez", ContactHandle: "+52 55 1234 0098", LastMessagePreview: "The tracking number has not updated in 3 days.", LastMessageAt: ago(9), AIStatus: AIHandedOff, SOPStage: SOPWaitingHuman, SOPWorkflowID: ptr("wf_2"), SOPWorkflowName: &shippingWorkflow, Unread: 1, Tags: []string{"logistics", "escalation"}, CreatedAt: ago(120), UpdatedAt: ago(9)},
		{ID: "conv_103", ChannelID: "ch_feishu", Channel: ChannelFeishu, ContactName: "Zhang Min", ContactHandle: "@zhangmin", LastMessagePreview: "Contract terms confirmed; waiting for seal.", LastMessageAt: ago(21), AssignedOperator: &operatorWei, AIStatus: AIMonitoring, SOPStage: SOPCompleted, SOPWorkflowName: &contractWorkflow, Tags: []string{"contract"}, CreatedAt: ago(240), UpdatedAt: ago(21)},
		{ID: "conv_104", ChannelID: "ch_telegram", Channel: ChannelTelegram, ContactName: "Alex Petrov", ContactHandle: "@alexp", LastMessagePreview: "Is there a way to integrate with our own CRM?", LastMessageAt: ago(33), AssignedOperator: &operatorSarah, AIStatus: AIAutoReplying, SOPStage: SOPNone, Tags: []string{"pre-sales"}, CreatedAt: ago(260), UpdatedAt: ago(33)},
		{ID: "conv_105", ChannelID: "ch_dingtalk", Channel: ChannelDingTalk, ContactName: "Chen Hao", ContactHandle: "@chenhao_partner", LastMessagePreview: "Question about authentication in the integration docs.", LastMessageAt: ago(58), AIStatus: AIIdle, SOPStage: SOPFailed, SOPWorkflowID: ptr("wf_3"), SOPWorkflowName: &partnerWorkflow, Unread: 3, Tags: []string{"partner", "technical"}, CreatedAt: ago(300), UpdatedAt: ago(58)},
		{ID: "conv_106", ChannelID: "ch_email_rpa", Channel: ChannelEmail, ContactName: "Support Ticket #8821", ContactHandle: "j.turner@acme-corp.com", LastMessagePreview: "Following up on the refund request submitted last week.", LastMessageAt: ago(72), AssignedOperator: &operatorWei, AIStatus: AIHandedOff, SOPStage: SOPWaitingHuman, SOPWorkflowName: &refundWorkflow, Tags: []string{"billing"}, CreatedAt: ago(360), UpdatedAt: ago(72)},
		{ID: "conv_107", ChannelID: "ch_whatsapp", Channel: ChannelWhatsApp, ContactName: "Fatima Al-Sayed", ContactHandle: "+971 50 220 3344", LastMessagePreview: "Perfect, thank you for the quick response!", LastMessageAt: ago(96), AssignedOperator: &operatorSarah, AIStatus: AIMonitoring, SOPStage: SOPCompleted, SOPWorkflowName: &orderWorkflow, Tags: []string{"order"}, CreatedAt: ago(380), UpdatedAt: ago(96)},
	}

	messages := []ConversationMessage{
		{ID: "msg_1", ConversationID: "conv_101", ChannelID: "ch_wecom", Channel: ChannelWeCom, Direction: DirectionInbound, Author: "Li Wei", Content: "Hello, we would like to learn more about the enterprise pricing plan.", Time: ago(40), MessageType: "text"},
		{ID: "msg_2", ConversationID: "conv_101", ChannelID: "ch_wecom", Channel: ChannelWeCom, Direction: DirectionOutbound, Author: "AI Assistant", Content: "Hi Li Wei, I prepared the enterprise pricing overview with seats and annual fees. I can send the PDF shortly.", Time: ago(38), MessageType: "text", IsAIGenerated: true},
		{ID: "msg_3", ConversationID: "conv_101", ChannelID: "ch_wecom", Channel: ChannelWeCom, Direction: DirectionInbound, Author: "Li Wei", Content: "Can you send the latest quote and payment terms?", Time: ago(4), MessageType: "text"},
	}

	snapshot := Snapshot{
		Connectors:      connectors,
		ChannelAccounts: accounts,
		Channels:        channels,
		PipelineStats: []PipelineStageStats{
			{Stage: StageConnector, Label: "Connector", ThroughputPerMin: 214, Failures1h: 3, AvgLatencyMs: 82},
			{Stage: StageIngest, Label: "Ingest", ThroughputPerMin: 211, Failures1h: 2, AvgLatencyMs: 46},
			{Stage: StageNormalize, Label: "Normalize", ThroughputPerMin: 209, Failures1h: 1, AvgLatencyMs: 34},
			{Stage: StageStore, Label: "Store", ThroughputPerMin: 209, Failures1h: 0, AvgLatencyMs: 21},
			{Stage: StageSOPAI, Label: "SOP / AI", ThroughputPerMin: 187, Failures1h: 6, AvgLatencyMs: 640},
			{Stage: StageOutbox, Label: "Outbox", ThroughputPerMin: 164, Failures1h: 4, AvgLatencyMs: 118},
			{Stage: StageDelivery, Label: "Delivery", ThroughputPerMin: 159, Failures1h: 5, AvgLatencyMs: 245},
		},
		MessageEvents: messageEvents,
		Conversations: conversations,
		Messages:      messages,
		AIPolicies: []AIPolicy{
			{ID: "pol_1", Kind: PolicyIntentClassification, Name: "Inbound Intent Classifier", Enabled: true, Priority: 1, TriggerCondition: "On every inbound message", FallbackStrategy: "Route to default queue", SuccessRate7d: 0.97, Invocations24h: 8420, CreatedAt: ago(1440), UpdatedAt: ago(30)},
			{ID: "pol_2", Kind: PolicyRiskDetection, Name: "Compliance & Risk Screening", Enabled: true, Priority: 1, TriggerCondition: "Message contains financial or legal terms", FallbackStrategy: "Flag for human review and block auto-reply", SuccessRate7d: 0.99, Invocations24h: 612, CreatedAt: ago(1440), UpdatedAt: ago(30)},
			{ID: "pol_3", Kind: PolicyReplyDrafting, Name: "Sales Reply Drafting", Enabled: true, Priority: 2, TriggerCondition: "Intent = pre-sales and confidence > 0.8", FallbackStrategy: "Draft only, require operator approval", SuccessRate7d: 0.91, Invocations24h: 2140, CreatedAt: ago(1440), UpdatedAt: ago(30)},
			{ID: "pol_4", Kind: PolicyKnowledgeRetrieval, Name: "Product Knowledge Retrieval", Enabled: true, Priority: 2, TriggerCondition: "Question matches knowledge base topics", FallbackStrategy: "Fall back to human handoff", SuccessRate7d: 0.94, Invocations24h: 3305, CreatedAt: ago(1440), UpdatedAt: ago(30)},
			{ID: "pol_5", Kind: PolicyToolCalling, Name: "Order Status Lookup", Enabled: true, Priority: 3, TriggerCondition: "Intent = order_status", FallbackStrategy: "Retry once, then escalate", SuccessRate7d: 0.88, Invocations24h: 1876, CreatedAt: ago(1440), UpdatedAt: ago(30)},
			{ID: "pol_6", Kind: PolicyHumanHandoff, Name: "Escalation Handoff Policy", Enabled: true, Priority: 1, TriggerCondition: "Risk flag or SLA breach imminent", FallbackStrategy: "Assign to on-call operator", SuccessRate7d: 0.99, Invocations24h: 214, CreatedAt: ago(1440), UpdatedAt: ago(30)},
			{ID: "pol_7", Kind: PolicyAutoReply, Name: "After-hours Auto Reply", Enabled: false, Priority: 4, TriggerCondition: "Outside business hours and no operator online", FallbackStrategy: "Queue for next business day", SuccessRate7d: 0.95, Invocations24h: 0, CreatedAt: ago(1440), UpdatedAt: ago(30)},
		},
		SOPWorkflows: []SOPWorkflow{
			{ID: "wf_1", Name: quoteWorkflow, Trigger: "Intent = pricing_request AND deal_size > 50k", Channels: []ChannelKind{ChannelWeCom, ChannelFeishu}, ActiveConversations: 18, CompletionRate: 0.86, SLAMinutes: 240, Status: WorkflowActive, CreatedAt: ago(1440), UpdatedAt: ago(120), Steps: []WorkflowStep{
				step("wf_1_s1", "Classify request", "Inbound message received", ptr("Classify intent and extract deal size"), nil, 2, "Route to manual triage", 1),
				step("wf_1_s2", "Generate quote draft", "Intent confirmed as pricing_request", ptr("Draft quote from pricing table"), nil, 5, "Notify sales ops", 2),
				step("wf_1_s3", "Sales approval", "Quote draft ready", nil, ptr("Sales manager reviews and approves discount"), 120, "Escalate to regional director", 3),
				step("wf_1_s4", "Send to contact", "Quote approved", ptr("Send PDF via original channel"), nil, 5, "Retry delivery, then alert operator", 4),
			}},
			{ID: "wf_2", Name: shippingWorkflow, Trigger: "Intent = shipping_delay AND days_delayed > 2", Channels: []ChannelKind{ChannelWhatsApp, ChannelEmail}, ActiveConversations: 7, CompletionRate: 0.72, SLAMinutes: 60, Status: WorkflowActive, CreatedAt: ago(1440), UpdatedAt: ago(120), Steps: []WorkflowStep{
				step("wf_2_s1", "Verify tracking status", "Escalation triggered", ptr("Call logistics API for latest tracking event"), nil, 3, "Escalate directly to human", 1),
				step("wf_2_s2", "Support review", "Tracking confirms delay", nil, ptr("Support agent contacts carrier"), 45, "Escalate to logistics manager", 2),
			}},
			{ID: "wf_3", Name: partnerWorkflow, Trigger: "New partner application submitted", Channels: []ChannelKind{ChannelDingTalk}, ActiveConversations: 4, CompletionRate: 0.64, SLAMinutes: 1440, Status: WorkflowPaused, CreatedAt: ago(1440), UpdatedAt: ago(120), Steps: []WorkflowStep{
				step("wf_3_s1", "Document collection", "Application received", ptr("Send required document checklist"), nil, 60, "Remind after 24h", 1),
				step("wf_3_s2", "Technical review", "Documents complete", nil, ptr("Solutions engineer reviews integration plan"), 720, "Escalate to partner manager", 2),
			}},
		},
		OutboxItems: []OutboxItem{
			{ID: "out_1", CreatedAt: ago(2), ChannelID: "ch_wecom", Channel: ChannelWeCom, ConversationID: "conv_101", ConversationLabel: "Li Wei - Enterprise Quote", MessageType: "Document", Sender: "AI Assistant", DeliveryMethod: SendAPI, Status: OutboxSending, IdempotencyKey: "seed-out-1"},
			{ID: "out_2", CreatedAt: ago(6), ChannelID: "ch_dingtalk", Channel: ChannelDingTalk, ConversationID: "conv_105", ConversationLabel: "Chen Hao - Partner Onboarding", MessageType: "Text", Sender: "System", DeliveryMethod: SendRPA, Status: OutboxFailed, RetryCount: 3, LastError: ptr("RPA session timeout after 30s"), IdempotencyKey: "seed-out-2"},
			{ID: "out_3", CreatedAt: ago(11), ChannelID: "ch_email_rpa", Channel: ChannelEmail, ConversationID: "conv_106", ConversationLabel: "Ticket #8821 - Refund", MessageType: "Email", Sender: "Wei Zhao", DeliveryMethod: SendManualApproval, Status: OutboxRequiresApproval, IdempotencyKey: "seed-out-3"},
			{ID: "out_4", CreatedAt: ago(14), ChannelID: "ch_whatsapp", Channel: ChannelWhatsApp, ConversationID: "conv_102", ConversationLabel: "Maria Gonzalez - Shipping", MessageType: "Text", Sender: "AI Assistant", DeliveryMethod: SendAPI, Status: OutboxSent, IdempotencyKey: "seed-out-4"},
			{ID: "out_5", CreatedAt: ago(19), ChannelID: "ch_telegram", Channel: ChannelTelegram, ConversationID: "conv_104", ConversationLabel: "Alex Petrov - Pre-sales", MessageType: "Text", Sender: "AI Assistant", DeliveryMethod: SendAPI, Status: OutboxPending, IdempotencyKey: "seed-out-5"},
			{ID: "out_6", CreatedAt: ago(26), ChannelID: "ch_feishu", Channel: ChannelFeishu, ConversationID: "conv_103", ConversationLabel: "Zhang Min - Contract", MessageType: "Card", Sender: "System", DeliveryMethod: SendAPI, Status: OutboxSent, RetryCount: 1, IdempotencyKey: "seed-out-6"},
			{ID: "out_7", CreatedAt: ago(31), ChannelID: "ch_whatsapp", Channel: ChannelWhatsApp, ConversationID: "conv_107", ConversationLabel: "Fatima Al-Sayed - Order", MessageType: "Text", Sender: "AI Assistant", DeliveryMethod: SendAPI, Status: OutboxFailed, RetryCount: 2, LastError: ptr("Recipient number opted out of messages"), IdempotencyKey: "seed-out-7"},
		},
		AuditLog: []AuditLogEntry{
			{ID: "aud_1", Time: ago(3), Actor: "Sarah Chen", ActorType: ActorUser, Action: "Approved outbound message", Target: "out_3", Channel: &emailChannel, Result: AuditSuccess, IP: ptr("10.20.4.18"), TraceID: "trace-audit-1"},
			{ID: "aud_2", Time: ago(12), Actor: "AI Assistant", ActorType: ActorAI, Action: "Drafted reply", Target: "conv_101", Channel: ptr(ChannelWeCom), Result: AuditSuccess, TraceID: "trace-audit-2"},
			{ID: "aud_3", Time: ago(24), Actor: "System", ActorType: ActorSystem, Action: "Disabled channel after repeated failures", Target: "ch_email_rpa", Channel: &emailChannel, Result: AuditSuccess, TraceID: "trace-audit-3"},
			{ID: "aud_4", Time: ago(40), Actor: "Wei Zhao", ActorType: ActorUser, Action: "Updated SOP workflow", Target: "wf_3", Result: AuditSuccess, IP: ptr("10.20.4.42"), TraceID: "trace-audit-4"},
			{ID: "aud_5", Time: ago(58), Actor: "System", ActorType: ActorSystem, Action: "API key rotation failed", Target: "ch_dingtalk", Channel: &dingTalkChannel, Result: AuditFailure, TraceID: "trace-audit-5"},
		},
		Incidents: []RecentIncident{
			{ID: "inc_1", Time: ago(18), Severity: IncidentWarning, Summary: "DingTalk webhook latency exceeded 5s threshold", Channel: &dingTalkChannel},
			{ID: "inc_2", Time: ago(52), Severity: IncidentCritical, Summary: "Email RPA connector disabled after 12 consecutive failures", Channel: &emailChannel},
			{ID: "inc_3", Time: ago(95), Severity: IncidentInfo, Summary: "WhatsApp template message rejected - outdated template ID", Channel: &whatsAppChannel},
		},
		TrafficSeries: seedTrafficSeries(),
		Settings: PlatformSettings{
			WorkspaceName:    "Acme Growth Ops",
			Timezone:         "Asia/Singapore",
			DefaultLanguage:  "en",
			Environment:      "production",
			Region:           "ap-southeast-1",
			RetentionDays:    90,
			WebhookURL:       "https://hooks.imintegration.local/events",
			EnabledProviders: []string{"openai", "internal-rules"},
		},
	}
	return snapshot
}

func seedTrafficSeries() []TrafficPoint {
	points := make([]TrafficPoint, 0, 24)
	for i := 0; i < 24; i++ {
		base := 300 + math.Sin(float64(i)/2.4)*120
		points = append(points, TrafficPoint{
			Hour:     fmt.Sprintf("%02d:00", i),
			Inbound:  max(20, int(math.Round(base+float64((i*37)%60)))),
			Outbound: max(15, int(math.Round(base*0.82+float64((i*29)%50)))),
		})
	}
	return points
}

func step(id, name, condition string, aiAction, humanAction *string, timeoutMinutes int, fallback string, position int) WorkflowStep {
	return WorkflowStep{
		ID: id, Name: name, Condition: condition, AIAction: aiAction, HumanAction: humanAction,
		TimeoutMinutes: timeoutMinutes, Fallback: fallback, Position: position,
	}
}

func ptr[T any](value T) *T {
	return &value
}
