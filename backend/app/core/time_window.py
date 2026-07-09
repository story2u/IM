from datetime import datetime, time
from zoneinfo import ZoneInfo

from pydantic import BaseModel, Field

from app.core.config import Settings


class WorkTimeConfig(BaseModel):
    timezone: str = "Asia/Shanghai"
    weekdays: list[int] = Field(default_factory=lambda: [1, 2, 3, 4, 5])
    start: str = "09:00"
    end: str = "18:30"
    auto_reply_after_hours: bool = True

    @classmethod
    def from_settings(cls, settings: Settings) -> "WorkTimeConfig":
        return cls(
            timezone=settings.default_timezone,
            weekdays=settings.default_workdays,
            start=settings.default_work_start,
            end=settings.default_work_end,
        )


def parse_hhmm(value: str) -> time:
    hour, minute = value.split(":", maxsplit=1)
    return time(hour=int(hour), minute=int(minute))


class WorkTimeService:
    def __init__(self, config: WorkTimeConfig) -> None:
        self.config = config

    def now(self) -> datetime:
        return datetime.now(ZoneInfo(self.config.timezone))

    def is_working_time(self, at: datetime | None = None) -> bool:
        current = at.astimezone(ZoneInfo(self.config.timezone)) if at else self.now()
        iso_weekday = current.isoweekday()
        if iso_weekday not in self.config.weekdays:
            return False

        start = parse_hhmm(self.config.start)
        end = parse_hhmm(self.config.end)
        current_time = current.time()
        if start <= end:
            return start <= current_time <= end

        # Overnight window, for example 22:00-06:00.
        return current_time >= start or current_time <= end
