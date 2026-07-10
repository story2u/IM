from uuid import UUID

from fastapi import APIRouter, Depends

from app.api.deps import get_message_repo, get_opportunity_repo, require_user
from app.application.dto import ChatMessageRead
from app.application.mappers import to_chat_message_read
from app.infrastructure.db.models import User
from app.infrastructure.db.repositories import MessageRepository, OpportunityRepository

router = APIRouter()


@router.get("", response_model=list[ChatMessageRead])
async def list_messages(
    opportunity_id: UUID,
    current_user: User = Depends(require_user),
    repo: MessageRepository = Depends(get_message_repo),
    opportunity_repo: OpportunityRepository = Depends(get_opportunity_repo),
) -> list[ChatMessageRead]:
    opportunity = await opportunity_repo.get(opportunity_id)
    if not opportunity or opportunity.owner_user_id != current_user.id:
        return []
    messages = await repo.list_by_opportunity(opportunity_id)
    return [to_chat_message_read(message) for message in messages]
