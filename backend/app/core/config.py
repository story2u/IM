from functools import lru_cache
from typing import Literal

from pydantic import AnyHttpUrl, Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",
    )

    app_env: Literal["local", "dev", "staging", "prod"] = "local"
    app_name: str = "Opportunity IM Assistant API"
    debug: bool = False
    api_v1_prefix: str = "/api/v1"

    database_url: str
    redis_url: str = "redis://redis:6379/0"
    celery_broker_url: str = "redis://redis:6379/1"
    celery_result_backend: str = "redis://redis:6379/2"

    admin_api_token: str
    jwt_secret_key: str = ""
    access_token_expire_minutes: int = 60 * 24 * 7
    frontend_base_url: str = "http://localhost:3000"
    cors_origins: list[AnyHttpUrl | str] = Field(default_factory=list)

    google_oauth_client_id: str = ""
    google_oauth_client_secret: str = ""
    google_oauth_redirect_uri: str = ""
    apple_oauth_client_id: str = ""
    apple_oauth_client_secret: str = ""
    apple_oauth_redirect_uri: str = ""
    apple_oauth_team_id: str = ""
    apple_oauth_key_id: str = ""
    apple_oauth_private_key: str = ""

    telegram_free_monitor_limit: int = 1

    default_timezone: str = "Asia/Shanghai"
    default_workdays: list[int] = Field(default_factory=lambda: [1, 2, 3, 4, 5])
    default_work_start: str = "09:00"
    default_work_end: str = "18:30"
    pending_human_sla_minutes: int = 30

    im_send_enabled: bool = False
    telegram_bot_token: str = ""
    telegram_webhook_secret: str = ""

    wecom_corp_id: str = ""
    wecom_agent_id: str = ""
    wecom_secret: str = ""
    wecom_token: str = ""
    wecom_aes_key: str = ""

    ai_enabled: bool = False
    litellm_model: str = "openai/gpt-4o-mini"
    openai_api_key: str = ""


@lru_cache
def get_settings() -> Settings:
    return Settings()
