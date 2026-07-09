from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status

from app.api.deps import get_rule_repo, require_admin
from app.application.dto import RuleCreate, RuleRead, RuleUpdate
from app.infrastructure.db.models import Rule
from app.infrastructure.db.repositories import RuleRepository

router = APIRouter()


@router.get("", response_model=list[RuleRead])
async def list_rules(
    repo: RuleRepository = Depends(get_rule_repo),
) -> list[Rule]:
    return await repo.list()


@router.post("", response_model=RuleRead)
async def create_rule(
    payload: RuleCreate,
    _: None = Depends(require_admin),
    repo: RuleRepository = Depends(get_rule_repo),
) -> Rule:
    return await repo.create(Rule(**payload.model_dump()))


@router.patch("/{rule_id}", response_model=RuleRead)
async def update_rule(
    rule_id: UUID,
    payload: RuleUpdate,
    _: None = Depends(require_admin),
    repo: RuleRepository = Depends(get_rule_repo),
) -> Rule:
    rule = await repo.get(rule_id)
    if not rule:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="rule not found")
    for key, value in payload.model_dump(exclude_unset=True).items():
        setattr(rule, key, value)
    return await repo.save(rule)


@router.delete("/{rule_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_rule(
    rule_id: UUID,
    _: None = Depends(require_admin),
    repo: RuleRepository = Depends(get_rule_repo),
) -> None:
    rule = await repo.get(rule_id)
    if not rule:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="rule not found")
    await repo.delete(rule)
