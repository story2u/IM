package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"im-go/internal/integrationhub/domain"
)

func (s *PostgresStore) ListAIPolicies(ctx context.Context) ([]domain.AIPolicy, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, kind, name, enabled, priority, trigger_condition, fallback_strategy,
       success_rate_7d, invocations_24h, created_at, updated_at
FROM integration_ai_policies
ORDER BY priority ASC, name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	policies := make([]domain.AIPolicy, 0)
	for rows.Next() {
		var policy domain.AIPolicy
		var kind string
		if err := rows.Scan(
			&policy.ID, &kind, &policy.Name, &policy.Enabled, &policy.Priority,
			&policy.TriggerCondition, &policy.FallbackStrategy, &policy.SuccessRate7d,
			&policy.Invocations24h, &policy.CreatedAt, &policy.UpdatedAt,
		); err != nil {
			return nil, err
		}
		policy.Kind = domain.AIPolicyKind(kind)
		policies = append(policies, policy)
	}
	return policies, rows.Err()
}

func (s *PostgresStore) UpdateAIPolicyEnabled(ctx context.Context, id string, enabled bool, now time.Time, audit domain.AuditLogEntry) (domain.AIPolicy, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.AIPolicy{}, err
	}
	defer rollback(ctx, tx)
	var policy domain.AIPolicy
	var kind string
	err = tx.QueryRow(ctx, `
UPDATE integration_ai_policies
SET enabled = $2, updated_at = $3
WHERE id = $1
RETURNING id, kind, name, enabled, priority, trigger_condition, fallback_strategy,
          success_rate_7d, invocations_24h, created_at, updated_at`,
		id, enabled, now).Scan(
		&policy.ID, &kind, &policy.Name, &policy.Enabled, &policy.Priority,
		&policy.TriggerCondition, &policy.FallbackStrategy, &policy.SuccessRate7d,
		&policy.Invocations24h, &policy.CreatedAt, &policy.UpdatedAt,
	)
	if err != nil {
		return domain.AIPolicy{}, mapNotFound(err)
	}
	policy.Kind = domain.AIPolicyKind(kind)
	if err := insertAudit(ctx, tx, audit); err != nil {
		return domain.AIPolicy{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.AIPolicy{}, err
	}
	return policy, nil
}

func (s *PostgresStore) ListSOPWorkflows(ctx context.Context) ([]domain.SOPWorkflow, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, name, trigger_expr, channels, active_conversations, completion_rate, sla_minutes, status, created_at, updated_at
FROM integration_sop_workflows
ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workflows := make([]domain.SOPWorkflow, 0)
	for rows.Next() {
		var workflow domain.SOPWorkflow
		var channels []string
		var status string
		if err := rows.Scan(
			&workflow.ID, &workflow.Name, &workflow.Trigger, &channels, &workflow.ActiveConversations,
			&workflow.CompletionRate, &workflow.SLAMinutes, &status, &workflow.CreatedAt, &workflow.UpdatedAt,
		); err != nil {
			return nil, err
		}
		workflow.Channels = channelKindsFromStrings(channels)
		workflow.Status = domain.WorkflowStatus(status)
		workflows = append(workflows, workflow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	steps, err := s.sopSteps(ctx)
	if err != nil {
		return nil, err
	}
	for i := range workflows {
		workflows[i].Steps = steps[workflows[i].ID]
	}
	return workflows, nil
}

func (s *PostgresStore) sopSteps(ctx context.Context) (map[string][]domain.WorkflowStep, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, workflow_id, position, name, condition_expr, ai_action, human_action, timeout_minutes, fallback
FROM integration_sop_steps
ORDER BY workflow_id, position`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	steps := map[string][]domain.WorkflowStep{}
	for rows.Next() {
		var workflowID string
		var step domain.WorkflowStep
		if err := rows.Scan(
			&step.ID, &workflowID, &step.Position, &step.Name, &step.Condition,
			&step.AIAction, &step.HumanAction, &step.TimeoutMinutes, &step.Fallback,
		); err != nil {
			return nil, err
		}
		steps[workflowID] = append(steps[workflowID], step)
	}
	return steps, rows.Err()
}

func (s *PostgresStore) Settings(ctx context.Context) (domain.PlatformSettings, error) {
	var settings domain.PlatformSettings
	err := s.pool.QueryRow(ctx, `
SELECT workspace_name, timezone, default_language, environment, region, retention_days, webhook_url, enabled_providers
FROM integration_platform_settings
WHERE id = 'default'`).Scan(
		&settings.WorkspaceName, &settings.Timezone, &settings.DefaultLanguage, &settings.Environment,
		&settings.Region, &settings.RetentionDays, &settings.WebhookURL, &settings.EnabledProviders,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return defaultSettings(), nil
		}
		return domain.PlatformSettings{}, mapNotFound(err)
	}
	return settings, nil
}

func (s *PostgresStore) UpdateSettings(ctx context.Context, patch domain.PlatformSettings, now time.Time, audit domain.AuditLogEntry) (domain.PlatformSettings, error) {
	current, err := s.Settings(ctx)
	if err != nil {
		return domain.PlatformSettings{}, err
	}
	if patch.WorkspaceName != "" {
		current.WorkspaceName = patch.WorkspaceName
	}
	if patch.Timezone != "" {
		current.Timezone = patch.Timezone
	}
	if patch.DefaultLanguage != "" {
		current.DefaultLanguage = patch.DefaultLanguage
	}
	if patch.Environment != "" {
		current.Environment = patch.Environment
	}
	if patch.Region != "" {
		current.Region = patch.Region
	}
	if patch.RetentionDays > 0 {
		current.RetentionDays = patch.RetentionDays
	}
	if patch.WebhookURL != "" {
		current.WebhookURL = patch.WebhookURL
	}
	if patch.EnabledProviders != nil {
		current.EnabledProviders = patch.EnabledProviders
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.PlatformSettings{}, err
	}
	defer rollback(ctx, tx)
	if _, err := tx.Exec(ctx, `
INSERT INTO integration_platform_settings (
  id, workspace_name, timezone, default_language, environment, region, retention_days, webhook_url, enabled_providers, updated_at
) VALUES ('default',$1,$2,$3,$4,$5,$6,$7,$8,$9)
ON CONFLICT (id) DO UPDATE SET
  workspace_name = EXCLUDED.workspace_name,
  timezone = EXCLUDED.timezone,
  default_language = EXCLUDED.default_language,
  environment = EXCLUDED.environment,
  region = EXCLUDED.region,
  retention_days = EXCLUDED.retention_days,
  webhook_url = EXCLUDED.webhook_url,
  enabled_providers = EXCLUDED.enabled_providers,
  updated_at = EXCLUDED.updated_at`,
		current.WorkspaceName, current.Timezone, current.DefaultLanguage, current.Environment, current.Region,
		current.RetentionDays, current.WebhookURL, current.EnabledProviders, now); err != nil {
		return domain.PlatformSettings{}, err
	}
	if err := insertAudit(ctx, tx, audit); err != nil {
		return domain.PlatformSettings{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.PlatformSettings{}, err
	}
	return current, nil
}

func defaultSettings() domain.PlatformSettings {
	return domain.PlatformSettings{
		WorkspaceName:    "IM Integration Platform",
		Timezone:         "UTC",
		DefaultLanguage:  "en",
		Environment:      "development",
		Region:           "local",
		RetentionDays:    90,
		WebhookURL:       "",
		EnabledProviders: []string{},
	}
}

func (s *PostgresStore) recentIncidents(ctx context.Context) ([]domain.RecentIncident, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, time, severity, summary, channel_kind
FROM integration_incidents
ORDER BY time DESC
LIMIT 20`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	incidents := make([]domain.RecentIncident, 0)
	for rows.Next() {
		var incident domain.RecentIncident
		var severity string
		var channel *string
		if err := rows.Scan(&incident.ID, &incident.Time, &severity, &incident.Summary, &channel); err != nil {
			return nil, err
		}
		incident.Severity = domain.IncidentSeverity(severity)
		if channel != nil {
			value := domain.ChannelKind(*channel)
			incident.Channel = &value
		}
		incidents = append(incidents, incident)
	}
	return incidents, rows.Err()
}

func (s *PostgresStore) trafficSeries(ctx context.Context) ([]domain.TrafficPoint, error) {
	rows, err := s.pool.Query(ctx, `
SELECT hour_label, inbound, outbound
FROM integration_traffic_points
ORDER BY hour_index ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	points := make([]domain.TrafficPoint, 0)
	for rows.Next() {
		var point domain.TrafficPoint
		if err := rows.Scan(&point.Hour, &point.Inbound, &point.Outbound); err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, rows.Err()
}
