"""add password reset challenges and JWT auth version

Revision ID: 202607150001
Revises: 202607140002
Create Date: 2026-07-15
"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op

revision: str = "202607150001"
down_revision: str | None = "202607140002"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.add_column(
        "users",
        sa.Column("auth_version", sa.Integer(), nullable=False, server_default="0"),
    )
    op.alter_column("users", "auth_version", server_default=None)
    op.create_table(
        "password_reset_challenges",
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("id", sa.Uuid(), nullable=False),
        sa.Column("user_id", sa.Uuid(), nullable=False),
        sa.Column("token_digest", sa.String(length=64), nullable=False),
        sa.Column("code_digest", sa.String(length=64), nullable=False),
        sa.Column("expires_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("used_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("failed_attempts", sa.Integer(), nullable=False),
        sa.CheckConstraint(
            "failed_attempts >= 0", name="ck_password_reset_failed_attempts"
        ),
        sa.ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("token_digest"),
    )
    op.create_index(
        "ix_password_reset_challenges_user_id",
        "password_reset_challenges",
        ["user_id"],
    )
    op.create_index(
        "ix_password_reset_challenges_expires_at",
        "password_reset_challenges",
        ["expires_at"],
    )
    op.create_index(
        "ix_password_reset_user_expires",
        "password_reset_challenges",
        ["user_id", "expires_at"],
    )


def downgrade() -> None:
    op.drop_table("password_reset_challenges")
    op.drop_column("users", "auth_version")
