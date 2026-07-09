"""init tables

Revision ID: 202607080001
Revises:
Create Date: 2026-07-08 00:01:00.000000
"""

from collections.abc import Sequence

import sqlalchemy as sa
import sqlmodel
from alembic import op
from sqlalchemy.dialects import postgresql

revision: str = "202607080001"
down_revision: str | None = None
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.create_table(
        "opportunities",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("channel", sa.Enum("TELEGRAM", "WECOM", native_enum=False), nullable=False),
        sa.Column("conversation_id", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("customer_external_id", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("contact_name", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("contact_avatar", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("source_message_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("title", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("summary", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("matched_keywords", postgresql.JSONB(astext_type=sa.Text()), nullable=False),
        sa.Column("confidence", sa.Float(), nullable=False),
        sa.Column("priority", sa.Enum("LOW", "NORMAL", "HIGH", "URGENT", native_enum=False), nullable=False),
        sa.Column(
            "status",
            sa.Enum(
                "PENDING_HUMAN",
                "AI_AUTO_REPLY",
                "REPLIED",
                "FOLLOWING",
                "IGNORED",
                "CLOSED",
                native_enum=False,
            ),
            nullable=False,
        ),
        sa.Column("detection_reason", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("ai_reply_draft", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("final_reply", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("assigned_to", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("follow_up_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("last_message_preview", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("last_message_at", sa.DateTime(timezone=True), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index("ix_opportunities_channel", "opportunities", ["channel"], unique=False)
    op.create_index("ix_opportunities_channel_conversation", "opportunities", ["channel", "conversation_id"], unique=False)
    op.create_index("ix_opportunities_conversation_id", "opportunities", ["conversation_id"], unique=False)
    op.create_index("ix_opportunities_customer_external_id", "opportunities", ["customer_external_id"], unique=False)
    op.create_index("ix_opportunities_last_message_at", "opportunities", ["last_message_at"], unique=False)
    op.create_index("ix_opportunities_priority", "opportunities", ["priority"], unique=False)
    op.create_index("ix_opportunities_source_message_id", "opportunities", ["source_message_id"], unique=False)
    op.create_index("ix_opportunities_status", "opportunities", ["status"], unique=False)
    op.create_index("ix_opportunities_status_created", "opportunities", ["status", "created_at"], unique=False)

    op.create_table(
        "messages",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("channel", sa.Enum("TELEGRAM", "WECOM", native_enum=False), nullable=False),
        sa.Column("external_message_id", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("conversation_id", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("sender_external_id", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("sender_display_name", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("direction", sa.Enum("INCOMING", "OUTGOING", native_enum=False), nullable=False),
        sa.Column("source", sa.Enum("HUMAN", "AI", native_enum=False), nullable=True),
        sa.Column("text", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("raw_payload", postgresql.JSONB(astext_type=sa.Text()), nullable=False),
        sa.Column("opportunity_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("sent_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("processed_at", sa.DateTime(timezone=True), nullable=True),
        sa.ForeignKeyConstraint(["opportunity_id"], ["opportunities.id"]),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("channel", "external_message_id", name="uq_message_channel_external"),
    )
    op.create_index("ix_messages_channel", "messages", ["channel"], unique=False)
    op.create_index("ix_messages_conversation_created", "messages", ["conversation_id", "created_at"], unique=False)
    op.create_index("ix_messages_conversation_id", "messages", ["conversation_id"], unique=False)
    op.create_index("ix_messages_direction", "messages", ["direction"], unique=False)
    op.create_index("ix_messages_external_message_id", "messages", ["external_message_id"], unique=False)
    op.create_index("ix_messages_opportunity_id", "messages", ["opportunity_id"], unique=False)
    op.create_index("ix_messages_sender_external_id", "messages", ["sender_external_id"], unique=False)
    op.create_index("ix_messages_sent_at", "messages", ["sent_at"], unique=False)

    op.create_table(
        "rules",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("name", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("enabled", sa.Boolean(), nullable=False),
        sa.Column("priority", sa.Integer(), nullable=False),
        sa.Column("rule_type", sa.Enum("KEYWORD", "REGEX", "AI_HINT", native_enum=False), nullable=False),
        sa.Column("pattern", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("score", sa.Float(), nullable=False),
        sa.Column("extra_data", postgresql.JSONB(astext_type=sa.Text()), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index("ix_rules_enabled", "rules", ["enabled"], unique=False)
    op.create_index("ix_rules_name", "rules", ["name"], unique=False)
    op.create_index("ix_rules_priority", "rules", ["priority"], unique=False)
    op.create_index("ix_rules_rule_type", "rules", ["rule_type"], unique=False)

    op.create_table(
        "app_configs",
        sa.Column("key", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("value", postgresql.JSONB(astext_type=sa.Text()), nullable=False),
        sa.Column("description", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.PrimaryKeyConstraint("key"),
    )

    op.create_table(
        "reply_templates",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("title", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("content", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("category", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("enabled", sa.Boolean(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index("ix_reply_templates_category", "reply_templates", ["category"], unique=False)
    op.create_index("ix_reply_templates_enabled", "reply_templates", ["enabled"], unique=False)
    op.create_index("ix_reply_templates_title", "reply_templates", ["title"], unique=False)


def downgrade() -> None:
    op.drop_index("ix_reply_templates_title", table_name="reply_templates")
    op.drop_index("ix_reply_templates_enabled", table_name="reply_templates")
    op.drop_index("ix_reply_templates_category", table_name="reply_templates")
    op.drop_table("reply_templates")

    op.drop_table("app_configs")

    op.drop_index("ix_rules_rule_type", table_name="rules")
    op.drop_index("ix_rules_priority", table_name="rules")
    op.drop_index("ix_rules_name", table_name="rules")
    op.drop_index("ix_rules_enabled", table_name="rules")
    op.drop_table("rules")

    op.drop_index("ix_messages_sent_at", table_name="messages")
    op.drop_index("ix_messages_sender_external_id", table_name="messages")
    op.drop_index("ix_messages_opportunity_id", table_name="messages")
    op.drop_index("ix_messages_external_message_id", table_name="messages")
    op.drop_index("ix_messages_direction", table_name="messages")
    op.drop_index("ix_messages_conversation_id", table_name="messages")
    op.drop_index("ix_messages_conversation_created", table_name="messages")
    op.drop_index("ix_messages_channel", table_name="messages")
    op.drop_table("messages")

    op.drop_index("ix_opportunities_status_created", table_name="opportunities")
    op.drop_index("ix_opportunities_status", table_name="opportunities")
    op.drop_index("ix_opportunities_source_message_id", table_name="opportunities")
    op.drop_index("ix_opportunities_priority", table_name="opportunities")
    op.drop_index("ix_opportunities_last_message_at", table_name="opportunities")
    op.drop_index("ix_opportunities_customer_external_id", table_name="opportunities")
    op.drop_index("ix_opportunities_conversation_id", table_name="opportunities")
    op.drop_index("ix_opportunities_channel_conversation", table_name="opportunities")
    op.drop_index("ix_opportunities_channel", table_name="opportunities")
    op.drop_table("opportunities")
