import importlib


def test_backend_entrypoints_import(monkeypatch) -> None:
    monkeypatch.setenv("DATABASE_URL", "postgresql+asyncpg://user:password@localhost:5432/im")
    monkeypatch.setenv("ADMIN_API_TOKEN", "test-token")
    monkeypatch.setenv("TELEGRAM_WEBHOOK_SECRET", "test-secret")

    importlib.import_module("app.infrastructure.db.repositories")
    importlib.import_module("app.worker.tasks")
    importlib.import_module("app.main")
