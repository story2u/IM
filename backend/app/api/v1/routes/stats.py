from fastapi import APIRouter, Depends

from app.api.deps import get_opportunity_repo, require_user
from app.application.dto import StatsSummaryRead
from app.application.use_cases.statistics import build_stats_summary
from app.infrastructure.db.models import User
from app.infrastructure.db.repositories import OpportunityRepository

router = APIRouter()


@router.get("/summary", response_model=StatsSummaryRead)
async def summary(
    current_user: User = Depends(require_user),
    repo: OpportunityRepository = Depends(get_opportunity_repo),
) -> StatsSummaryRead:
    opportunities = await repo.list(owner_user_id=current_user.id, limit=500)
    return build_stats_summary(opportunities)
