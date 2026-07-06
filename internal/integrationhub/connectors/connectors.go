package connectors

import (
	"context"
	"encoding/json"

	"im-go/internal/integrationhub/domain"
)

type Adapter interface {
	Kind() domain.ChannelKind
	ParseInbound(ctx context.Context, raw json.RawMessage) (domain.InboundEvent, error)
}
