"""merge release and device-agent migration heads

Revision ID: 202607180013
Revises: 202607160001, 202607180012
Create Date: 2026-07-18
"""

from collections.abc import Sequence


revision: str = "202607180013"
down_revision: tuple[str, str] = ("202607160001", "202607180012")
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    pass


def downgrade() -> None:
    pass
