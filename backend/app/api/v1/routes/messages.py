from uuid import UUID

from fastapi import APIRouter, Depends

from app.api.deps import get_message_repo
from app.application.dto import ChatMessageRead
from app.application.mappers import to_chat_message_read
from app.infrastructure.db.repositories import MessageRepository

router = APIRouter()


@router.get("", response_model=list[ChatMessageRead])
async def list_messages(
    opportunity_id: UUID,
    repo: MessageRepository = Depends(get_message_repo),
) -> list[ChatMessageRead]:
    messages = await repo.list_by_opportunity(opportunity_id)
    return [to_chat_message_read(message) for message in messages]
