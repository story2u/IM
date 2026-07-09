"""add opportunity source fields

Revision ID: 202607090001
Revises: 202607080001
Create Date: 2026-07-09 00:01:00.000000
"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

revision: str = "202607090001"
down_revision: str | None = "202607080001"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.add_column(
        "opportunities",
        sa.Column("source_type", sa.String(), nullable=False, server_default="private"),
    )
    op.add_column("opportunities", sa.Column("group_name", sa.String(), nullable=True))
    op.add_column(
        "opportunities",
        sa.Column(
            "raw_message_links",
            postgresql.JSONB(astext_type=sa.Text()),
            nullable=False,
            server_default=sa.text("'[]'::jsonb"),
        ),
    )
    op.add_column(
        "opportunities",
        sa.Column("trust_score", sa.Integer(), nullable=False, server_default="70"),
    )
    op.create_index("ix_opportunities_source_type", "opportunities", ["source_type"], unique=False)

    op.alter_column("opportunities", "source_type", server_default=None)
    op.alter_column("opportunities", "raw_message_links", server_default=None)
    op.alter_column("opportunities", "trust_score", server_default=None)


def downgrade() -> None:
    op.drop_index("ix_opportunities_source_type", table_name="opportunities")
    op.drop_column("opportunities", "trust_score")
    op.drop_column("opportunities", "raw_message_links")
    op.drop_column("opportunities", "group_name")
    op.drop_column("opportunities", "source_type")
