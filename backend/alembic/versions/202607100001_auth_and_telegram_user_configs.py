"""add oauth users and telegram monitors

Revision ID: 202607100001
Revises: 202607090001
Create Date: 2026-07-10 20:55:00.000000
"""

from collections.abc import Sequence

import sqlalchemy as sa
import sqlmodel
from alembic import op
from sqlalchemy.dialects import postgresql

revision: str = "202607100001"
down_revision: str | None = "202607090001"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.create_table(
        "users",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("email", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("display_name", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("avatar_url", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("password_hash", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("is_active", sa.Boolean(), nullable=False),
        sa.Column("is_admin", sa.Boolean(), nullable=False),
        sa.Column("last_login_at", sa.DateTime(timezone=True), nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index("ix_users_email", "users", ["email"], unique=True)
    op.create_index("ix_users_is_active", "users", ["is_active"], unique=False)
    op.create_index("ix_users_is_admin", "users", ["is_admin"], unique=False)

    op.create_table(
        "auth_accounts",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("user_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("provider", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("provider_subject", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("email", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.ForeignKeyConstraint(["user_id"], ["users.id"]),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("provider", "provider_subject", name="uq_auth_accounts_provider_subject"),
    )
    op.create_index("ix_auth_accounts_email", "auth_accounts", ["email"], unique=False)
    op.create_index("ix_auth_accounts_provider", "auth_accounts", ["provider"], unique=False)
    op.create_index("ix_auth_accounts_provider_subject", "auth_accounts", ["provider_subject"], unique=False)
    op.create_index("ix_auth_accounts_user_id", "auth_accounts", ["user_id"], unique=False)
    op.create_index("ix_auth_accounts_user_provider", "auth_accounts", ["user_id", "provider"], unique=False)

    op.add_column("opportunities", sa.Column("owner_user_id", postgresql.UUID(as_uuid=True), nullable=True))
    op.create_foreign_key(
        "fk_opportunities_owner_user_id_users",
        "opportunities",
        "users",
        ["owner_user_id"],
        ["id"],
    )
    op.create_index("ix_opportunities_owner_user_id", "opportunities", ["owner_user_id"], unique=False)

    op.add_column("messages", sa.Column("owner_user_id", postgresql.UUID(as_uuid=True), nullable=True))
    op.create_foreign_key(
        "fk_messages_owner_user_id_users",
        "messages",
        "users",
        ["owner_user_id"],
        ["id"],
    )
    op.create_index("ix_messages_owner_user_id", "messages", ["owner_user_id"], unique=False)

    op.create_table(
        "telegram_user_configs",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("user_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("enabled", sa.Boolean(), nullable=False),
        sa.Column("api_id", sa.Integer(), nullable=True),
        sa.Column("api_hash_encrypted", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("session_encrypted", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.ForeignKeyConstraint(["user_id"], ["users.id"]),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("user_id", name="uq_telegram_user_configs_user_id"),
    )
    op.create_index("ix_telegram_user_configs_user_id", "telegram_user_configs", ["user_id"], unique=False)

    op.create_table(
        "telegram_monitors",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("user_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("telegram_config_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("enabled", sa.Boolean(), nullable=False),
        sa.Column("name", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("chat_id", sqlmodel.sql.sqltypes.AutoString(), nullable=False),
        sa.Column("chat_title", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.Column("backfill_limit", sa.Integer(), nullable=False),
        sa.Column("last_error", sqlmodel.sql.sqltypes.AutoString(), nullable=True),
        sa.ForeignKeyConstraint(["telegram_config_id"], ["telegram_user_configs.id"]),
        sa.ForeignKeyConstraint(["user_id"], ["users.id"]),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("user_id", "chat_id", name="uq_telegram_monitors_user_chat"),
    )
    op.create_index("ix_telegram_monitors_chat_id", "telegram_monitors", ["chat_id"], unique=False)
    op.create_index("ix_telegram_monitors_enabled", "telegram_monitors", ["enabled"], unique=False)
    op.create_index("ix_telegram_monitors_telegram_config_id", "telegram_monitors", ["telegram_config_id"], unique=False)
    op.create_index("ix_telegram_monitors_user_id", "telegram_monitors", ["user_id"], unique=False)


def downgrade() -> None:
    op.drop_index("ix_telegram_monitors_user_id", table_name="telegram_monitors")
    op.drop_index("ix_telegram_monitors_telegram_config_id", table_name="telegram_monitors")
    op.drop_index("ix_telegram_monitors_enabled", table_name="telegram_monitors")
    op.drop_index("ix_telegram_monitors_chat_id", table_name="telegram_monitors")
    op.drop_table("telegram_monitors")

    op.drop_index("ix_telegram_user_configs_user_id", table_name="telegram_user_configs")
    op.drop_table("telegram_user_configs")

    op.drop_index("ix_messages_owner_user_id", table_name="messages")
    op.drop_constraint("fk_messages_owner_user_id_users", "messages", type_="foreignkey")
    op.drop_column("messages", "owner_user_id")

    op.drop_index("ix_opportunities_owner_user_id", table_name="opportunities")
    op.drop_constraint("fk_opportunities_owner_user_id_users", "opportunities", type_="foreignkey")
    op.drop_column("opportunities", "owner_user_id")

    op.drop_index("ix_auth_accounts_user_provider", table_name="auth_accounts")
    op.drop_index("ix_auth_accounts_user_id", table_name="auth_accounts")
    op.drop_index("ix_auth_accounts_provider_subject", table_name="auth_accounts")
    op.drop_index("ix_auth_accounts_provider", table_name="auth_accounts")
    op.drop_index("ix_auth_accounts_email", table_name="auth_accounts")
    op.drop_table("auth_accounts")

    op.drop_index("ix_users_is_admin", table_name="users")
    op.drop_index("ix_users_is_active", table_name="users")
    op.drop_index("ix_users_email", table_name="users")
    op.drop_table("users")
