from fastapi import APIRouter, Depends
from sqlmodel.ext.asyncio.session import AsyncSession

from app.api.deps import get_work_time_service, require_admin
from app.core.config import Settings, get_settings
from app.application.dto import ConfigRead, ConfigUpdate, WorkModeRead
from app.core.time_window import WorkTimeConfig, WorkTimeService
from app.infrastructure.db.models import AppConfig
from app.infrastructure.db.repositories import ConfigRepository
from app.infrastructure.db.session import get_session

router = APIRouter()


@router.get("", response_model=list[ConfigRead])
async def list_configs(
    settings: Settings = Depends(get_settings),
    session: AsyncSession = Depends(get_session),
) -> list[ConfigRead]:
    # MVP exposes the known config documents explicitly.
    repo = ConfigRepository(session)
    work_hours = await repo.get_value("working_hours")
    default_work_hours = WorkTimeConfig.from_settings(settings).model_dump()
    return [
        ConfigRead(
            key="working_hours",
            value=work_hours or default_work_hours,
            description="工作时间配置",
        )
    ]


@router.get("/work-mode", response_model=WorkModeRead)
async def work_mode(
    service: WorkTimeService = Depends(get_work_time_service),
) -> WorkModeRead:
    is_working_time = service.is_working_time()
    return WorkModeRead(
        mode="work" if is_working_time else "ai",
        is_working_time=is_working_time,
        work_time=service.config,
    )


@router.patch("/{key}", response_model=ConfigRead)
async def update_config(
    key: str,
    payload: ConfigUpdate,
    _: None = Depends(require_admin),
    session: AsyncSession = Depends(get_session),
) -> ConfigRead:
    if key == "working_hours":
        WorkTimeConfig.model_validate(payload.value)

    repo = ConfigRepository(session)
    config: AppConfig = await repo.set_value(
        key,
        payload.value,
        description=payload.description,
    )
    return ConfigRead(
        key=config.key,
        value=config.value,
        description=config.description,
        updated_at=config.updated_at,
    )
