from datetime import datetime
from zoneinfo import ZoneInfo

from app.core.time_window import WorkTimeConfig, WorkTimeService


def test_work_time_service_detects_working_hours() -> None:
    service = WorkTimeService(
        WorkTimeConfig(
            timezone="Asia/Shanghai",
            weekdays=[1, 2, 3, 4, 5],
            start="09:00",
            end="18:00",
        )
    )

    assert service.is_working_time(datetime(2026, 7, 8, 10, 0, tzinfo=ZoneInfo("Asia/Shanghai")))
    assert not service.is_working_time(datetime(2026, 7, 8, 20, 0, tzinfo=ZoneInfo("Asia/Shanghai")))


def test_work_time_service_detects_weekend() -> None:
    service = WorkTimeService(
        WorkTimeConfig(
            timezone="Asia/Shanghai",
            weekdays=[1, 2, 3, 4, 5],
            start="09:00",
            end="18:00",
        )
    )

    assert not service.is_working_time(datetime(2026, 7, 11, 10, 0, tzinfo=ZoneInfo("Asia/Shanghai")))
