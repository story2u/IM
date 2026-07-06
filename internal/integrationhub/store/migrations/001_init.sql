CREATE TABLE IF NOT EXISTS integration_connectors (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS integration_channels (
    id TEXT PRIMARY KEY,
    connector_id TEXT REFERENCES integration_connectors(id) ON DELETE SET NULL,
    kind TEXT NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    receive_capabilities TEXT[] NOT NULL DEFAULT '{}',
    send_capabilities TEXT[] NOT NULL DEFAULT '{}',
    last_sync_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    error_count_24h INTEGER NOT NULL DEFAULT 0,
    messages_today INTEGER NOT NULL DEFAULT 0,
    active_conversations INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_integration_channels_kind ON integration_channels(kind);
CREATE INDEX IF NOT EXISTS idx_integration_channels_status ON integration_channels(status);

CREATE TABLE IF NOT EXISTS integration_channel_accounts (
    id TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL REFERENCES integration_channels(id) ON DELETE CASCADE,
    display_name TEXT NOT NULL,
    external_account_id TEXT NOT NULL,
    status TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_integration_channel_accounts_channel ON integration_channel_accounts(channel_id);

CREATE TABLE IF NOT EXISTS integration_sop_workflows (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    trigger_expr TEXT NOT NULL,
    channels TEXT[] NOT NULL DEFAULT '{}',
    active_conversations INTEGER NOT NULL DEFAULT 0,
    completion_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
    sla_minutes INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS integration_sop_steps (
    id TEXT PRIMARY KEY,
    workflow_id TEXT NOT NULL REFERENCES integration_sop_workflows(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    name TEXT NOT NULL,
    condition_expr TEXT NOT NULL,
    ai_action TEXT,
    human_action TEXT,
    timeout_minutes INTEGER NOT NULL DEFAULT 0,
    fallback TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_integration_sop_steps_workflow ON integration_sop_steps(workflow_id, position);

CREATE TABLE IF NOT EXISTS integration_conversations (
    id TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL REFERENCES integration_channels(id) ON DELETE RESTRICT,
    channel_kind TEXT NOT NULL,
    contact_name TEXT NOT NULL,
    contact_handle TEXT NOT NULL,
    last_message_preview TEXT NOT NULL DEFAULT '',
    last_message_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    assigned_operator TEXT,
    ai_status TEXT NOT NULL,
    sop_stage TEXT NOT NULL,
    sop_workflow_id TEXT REFERENCES integration_sop_workflows(id) ON DELETE SET NULL,
    sop_workflow_name TEXT,
    unread INTEGER NOT NULL DEFAULT 0,
    tags TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_integration_conversations_channel ON integration_conversations(channel_kind);
CREATE INDEX IF NOT EXISTS idx_integration_conversations_last_message ON integration_conversations(last_message_at DESC);

CREATE TABLE IF NOT EXISTS integration_messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL REFERENCES integration_conversations(id) ON DELETE CASCADE,
    channel_id TEXT NOT NULL REFERENCES integration_channels(id) ON DELETE RESTRICT,
    channel_kind TEXT NOT NULL,
    direction TEXT NOT NULL,
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    message_type TEXT NOT NULL DEFAULT 'text',
    is_ai_generated BOOLEAN NOT NULL DEFAULT false,
    external_message_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_integration_messages_conversation_time ON integration_messages(conversation_id, created_at);

CREATE TABLE IF NOT EXISTS integration_inbound_events (
    id TEXT PRIMARY KEY,
    channel_id TEXT REFERENCES integration_channels(id) ON DELETE SET NULL,
    connector_kind TEXT NOT NULL,
    event_type TEXT NOT NULL,
    external_event_id TEXT,
    received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    normalized_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    adapter_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    trace_id TEXT NOT NULL,
    status TEXT NOT NULL,
    error TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_integration_inbound_events_external
    ON integration_inbound_events(connector_kind, external_event_id)
    WHERE external_event_id IS NOT NULL AND external_event_id <> '';

CREATE INDEX IF NOT EXISTS idx_integration_inbound_events_received ON integration_inbound_events(received_at DESC);

CREATE TABLE IF NOT EXISTS integration_outbound_commands (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL REFERENCES integration_conversations(id) ON DELETE CASCADE,
    channel_id TEXT NOT NULL REFERENCES integration_channels(id) ON DELETE RESTRICT,
    channel_kind TEXT NOT NULL,
    conversation_label TEXT NOT NULL,
    message_id TEXT REFERENCES integration_messages(id) ON DELETE SET NULL,
    message_type TEXT NOT NULL,
    sender TEXT NOT NULL,
    delivery_method TEXT NOT NULL,
    status TEXT NOT NULL,
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    idempotency_key TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_integration_outbound_commands_status ON integration_outbound_commands(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_integration_outbound_commands_conversation ON integration_outbound_commands(conversation_id, created_at DESC);

CREATE TABLE IF NOT EXISTS integration_message_events (
    id TEXT PRIMARY KEY,
    time TIMESTAMPTZ NOT NULL DEFAULT now(),
    channel_kind TEXT NOT NULL,
    direction TEXT NOT NULL,
    conversation_id TEXT NOT NULL,
    conversation_label TEXT NOT NULL,
    event_type TEXT NOT NULL,
    status TEXT NOT NULL,
    latency_ms INTEGER NOT NULL DEFAULT 0,
    trace_id TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_integration_message_events_time ON integration_message_events(time DESC);
CREATE INDEX IF NOT EXISTS idx_integration_message_events_filters ON integration_message_events(channel_kind, status, event_type);
CREATE INDEX IF NOT EXISTS idx_integration_message_events_trace ON integration_message_events(trace_id);

CREATE TABLE IF NOT EXISTS integration_ai_policies (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL,
    name TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    priority INTEGER NOT NULL DEFAULT 1,
    trigger_condition TEXT NOT NULL,
    fallback_strategy TEXT NOT NULL,
    success_rate_7d DOUBLE PRECISION NOT NULL DEFAULT 0,
    invocations_24h INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS integration_audit_logs (
    id TEXT PRIMARY KEY,
    time TIMESTAMPTZ NOT NULL DEFAULT now(),
    actor TEXT NOT NULL,
    actor_type TEXT NOT NULL,
    action TEXT NOT NULL,
    target TEXT NOT NULL,
    channel_kind TEXT,
    result TEXT NOT NULL,
    ip TEXT,
    trace_id TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_integration_audit_logs_time ON integration_audit_logs(time DESC);
CREATE INDEX IF NOT EXISTS idx_integration_audit_logs_actor_type ON integration_audit_logs(actor_type);

CREATE TABLE IF NOT EXISTS integration_incidents (
    id TEXT PRIMARY KEY,
    time TIMESTAMPTZ NOT NULL DEFAULT now(),
    severity TEXT NOT NULL,
    summary TEXT NOT NULL,
    channel_kind TEXT
);

CREATE INDEX IF NOT EXISTS idx_integration_incidents_time ON integration_incidents(time DESC);

CREATE TABLE IF NOT EXISTS integration_pipeline_stage_stats (
    stage TEXT PRIMARY KEY,
    label TEXT NOT NULL,
    throughput_per_min INTEGER NOT NULL DEFAULT 0,
    failures_1h INTEGER NOT NULL DEFAULT 0,
    avg_latency_ms INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS integration_traffic_points (
    hour_index INTEGER PRIMARY KEY,
    hour_label TEXT NOT NULL,
    inbound INTEGER NOT NULL DEFAULT 0,
    outbound INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS integration_platform_settings (
    id TEXT PRIMARY KEY,
    workspace_name TEXT NOT NULL,
    timezone TEXT NOT NULL,
    default_language TEXT NOT NULL,
    environment TEXT NOT NULL,
    region TEXT NOT NULL,
    retention_days INTEGER NOT NULL,
    webhook_url TEXT NOT NULL,
    enabled_providers TEXT[] NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
