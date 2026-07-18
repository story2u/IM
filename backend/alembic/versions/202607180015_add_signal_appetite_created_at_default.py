"""add signal appetite created_at default

Revision ID: 202607180015
Revises: 202607180014
Create Date: 2026-07-18
"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op

revision: str = "202607180015"
down_revision: str | None = "202607180014"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.alter_column(
        "signal_appetite_events",
        "created_at",
        server_default=sa.text("now()"),
    )


def downgrade() -> None:
    op.alter_column(
        "signal_appetite_events",
        "created_at",
        server_default=None,
    )
