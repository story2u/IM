"""repair oauth user schema

Revision ID: 202607100002
Revises: 202607100001
Create Date: 2026-07-10 22:40:00.000000
"""

from collections.abc import Sequence

from alembic import op

revision: str = "202607100002"
down_revision: str | None = "202607100001"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    # This migration is intentionally idempotent. It repairs production databases
    # where a previous one-shot migrate container may have left Alembic state and
    # actual OAuth tables out of sync.
    op.execute("CREATE EXTENSION IF NOT EXISTS pgcrypto")

    op.execute(
        """
        CREATE TABLE IF NOT EXISTS users (
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            id UUID NOT NULL DEFAULT gen_random_uuid(),
            email VARCHAR NOT NULL,
            display_name VARCHAR NOT NULL DEFAULT '',
            avatar_url VARCHAR NOT NULL DEFAULT '',
            password_hash VARCHAR NULL,
            is_active BOOLEAN NOT NULL DEFAULT true,
            is_admin BOOLEAN NOT NULL DEFAULT false,
            last_login_at TIMESTAMP WITH TIME ZONE NULL,
            PRIMARY KEY (id)
        )
        """
    )
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS id UUID")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS display_name VARCHAR")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url VARCHAR")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN")
    op.execute("ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP WITH TIME ZONE")
    op.execute("UPDATE users SET created_at = now() WHERE created_at IS NULL")
    op.execute("UPDATE users SET updated_at = now() WHERE updated_at IS NULL")
    op.execute("UPDATE users SET id = gen_random_uuid() WHERE id IS NULL")
    op.execute("UPDATE users SET display_name = COALESCE(NULLIF(display_name, ''), email, '') WHERE display_name IS NULL OR display_name = ''")
    op.execute("UPDATE users SET avatar_url = '' WHERE avatar_url IS NULL")
    op.execute("UPDATE users SET is_active = true WHERE is_active IS NULL")
    op.execute("UPDATE users SET is_admin = false WHERE is_admin IS NULL")
    op.execute("ALTER TABLE users ALTER COLUMN created_at SET NOT NULL")
    op.execute("ALTER TABLE users ALTER COLUMN updated_at SET NOT NULL")
    op.execute("ALTER TABLE users ALTER COLUMN id SET NOT NULL")
    op.execute("ALTER TABLE users ALTER COLUMN email SET NOT NULL")
    op.execute("ALTER TABLE users ALTER COLUMN display_name SET NOT NULL")
    op.execute("ALTER TABLE users ALTER COLUMN avatar_url SET NOT NULL")
    op.execute("ALTER TABLE users ALTER COLUMN is_active SET NOT NULL")
    op.execute("ALTER TABLE users ALTER COLUMN is_admin SET NOT NULL")
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conrelid = 'users'::regclass AND contype = 'p'
            ) THEN
                ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
            END IF;
        END $$;
        """
    )
    op.execute("CREATE UNIQUE INDEX IF NOT EXISTS ix_users_email ON users (email)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_users_is_active ON users (is_active)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_users_is_admin ON users (is_admin)")

    op.execute(
        """
        CREATE TABLE IF NOT EXISTS auth_accounts (
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            id UUID NOT NULL DEFAULT gen_random_uuid(),
            user_id UUID NOT NULL,
            provider VARCHAR NOT NULL,
            provider_subject VARCHAR NOT NULL,
            email VARCHAR NULL,
            PRIMARY KEY (id)
        )
        """
    )
    op.execute("ALTER TABLE auth_accounts ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE")
    op.execute("ALTER TABLE auth_accounts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE")
    op.execute("ALTER TABLE auth_accounts ADD COLUMN IF NOT EXISTS id UUID")
    op.execute("ALTER TABLE auth_accounts ADD COLUMN IF NOT EXISTS user_id UUID")
    op.execute("ALTER TABLE auth_accounts ADD COLUMN IF NOT EXISTS provider VARCHAR")
    op.execute("ALTER TABLE auth_accounts ADD COLUMN IF NOT EXISTS provider_subject VARCHAR")
    op.execute("ALTER TABLE auth_accounts ADD COLUMN IF NOT EXISTS email VARCHAR")
    op.execute("UPDATE auth_accounts SET created_at = now() WHERE created_at IS NULL")
    op.execute("UPDATE auth_accounts SET updated_at = now() WHERE updated_at IS NULL")
    op.execute("UPDATE auth_accounts SET id = gen_random_uuid() WHERE id IS NULL")
    op.execute("ALTER TABLE auth_accounts ALTER COLUMN created_at SET NOT NULL")
    op.execute("ALTER TABLE auth_accounts ALTER COLUMN updated_at SET NOT NULL")
    op.execute("ALTER TABLE auth_accounts ALTER COLUMN id SET NOT NULL")
    op.execute("ALTER TABLE auth_accounts ALTER COLUMN user_id SET NOT NULL")
    op.execute("ALTER TABLE auth_accounts ALTER COLUMN provider SET NOT NULL")
    op.execute("ALTER TABLE auth_accounts ALTER COLUMN provider_subject SET NOT NULL")
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conrelid = 'auth_accounts'::regclass AND contype = 'p'
            ) THEN
                ALTER TABLE auth_accounts ADD CONSTRAINT auth_accounts_pkey PRIMARY KEY (id);
            END IF;
        END $$;
        """
    )
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'fk_auth_accounts_user_id_users'
            ) THEN
                ALTER TABLE auth_accounts
                    ADD CONSTRAINT fk_auth_accounts_user_id_users
                    FOREIGN KEY (user_id) REFERENCES users (id);
            END IF;
        END $$;
        """
    )
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'uq_auth_accounts_provider_subject'
            ) THEN
                ALTER TABLE auth_accounts
                    ADD CONSTRAINT uq_auth_accounts_provider_subject
                    UNIQUE (provider, provider_subject);
            END IF;
        END $$;
        """
    )
    op.execute("CREATE INDEX IF NOT EXISTS ix_auth_accounts_email ON auth_accounts (email)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_auth_accounts_provider ON auth_accounts (provider)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_auth_accounts_provider_subject ON auth_accounts (provider_subject)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_auth_accounts_user_id ON auth_accounts (user_id)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_auth_accounts_user_provider ON auth_accounts (user_id, provider)")

    op.execute("ALTER TABLE opportunities ADD COLUMN IF NOT EXISTS owner_user_id UUID")
    op.execute("CREATE INDEX IF NOT EXISTS ix_opportunities_owner_user_id ON opportunities (owner_user_id)")
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'fk_opportunities_owner_user_id_users'
            ) THEN
                ALTER TABLE opportunities
                    ADD CONSTRAINT fk_opportunities_owner_user_id_users
                    FOREIGN KEY (owner_user_id) REFERENCES users (id);
            END IF;
        END $$;
        """
    )

    op.execute("ALTER TABLE messages ADD COLUMN IF NOT EXISTS owner_user_id UUID")
    op.execute("CREATE INDEX IF NOT EXISTS ix_messages_owner_user_id ON messages (owner_user_id)")
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'fk_messages_owner_user_id_users'
            ) THEN
                ALTER TABLE messages
                    ADD CONSTRAINT fk_messages_owner_user_id_users
                    FOREIGN KEY (owner_user_id) REFERENCES users (id);
            END IF;
        END $$;
        """
    )

    op.execute(
        """
        CREATE TABLE IF NOT EXISTS telegram_user_configs (
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            id UUID NOT NULL DEFAULT gen_random_uuid(),
            user_id UUID NOT NULL,
            enabled BOOLEAN NOT NULL DEFAULT false,
            api_id INTEGER NULL,
            api_hash_encrypted VARCHAR NULL,
            session_encrypted VARCHAR NULL,
            PRIMARY KEY (id)
        )
        """
    )
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'fk_telegram_user_configs_user_id_users'
            ) THEN
                ALTER TABLE telegram_user_configs
                    ADD CONSTRAINT fk_telegram_user_configs_user_id_users
                    FOREIGN KEY (user_id) REFERENCES users (id);
            END IF;
        END $$;
        """
    )
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'uq_telegram_user_configs_user_id'
            ) THEN
                ALTER TABLE telegram_user_configs
                    ADD CONSTRAINT uq_telegram_user_configs_user_id UNIQUE (user_id);
            END IF;
        END $$;
        """
    )
    op.execute("CREATE INDEX IF NOT EXISTS ix_telegram_user_configs_user_id ON telegram_user_configs (user_id)")

    op.execute(
        """
        CREATE TABLE IF NOT EXISTS telegram_monitors (
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            id UUID NOT NULL DEFAULT gen_random_uuid(),
            user_id UUID NOT NULL,
            telegram_config_id UUID NOT NULL,
            enabled BOOLEAN NOT NULL DEFAULT true,
            name VARCHAR NOT NULL DEFAULT 'Telegram 群监控',
            chat_id VARCHAR NOT NULL,
            chat_title VARCHAR NULL,
            backfill_limit INTEGER NOT NULL DEFAULT 30,
            last_error VARCHAR NULL,
            PRIMARY KEY (id)
        )
        """
    )
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'fk_telegram_monitors_user_id_users'
            ) THEN
                ALTER TABLE telegram_monitors
                    ADD CONSTRAINT fk_telegram_monitors_user_id_users
                    FOREIGN KEY (user_id) REFERENCES users (id);
            END IF;
        END $$;
        """
    )
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'fk_telegram_monitors_telegram_config_id'
            ) THEN
                ALTER TABLE telegram_monitors
                    ADD CONSTRAINT fk_telegram_monitors_telegram_config_id
                    FOREIGN KEY (telegram_config_id) REFERENCES telegram_user_configs (id);
            END IF;
        END $$;
        """
    )
    op.execute(
        """
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_constraint
                WHERE conname = 'uq_telegram_monitors_user_chat'
            ) THEN
                ALTER TABLE telegram_monitors
                    ADD CONSTRAINT uq_telegram_monitors_user_chat UNIQUE (user_id, chat_id);
            END IF;
        END $$;
        """
    )
    op.execute("CREATE INDEX IF NOT EXISTS ix_telegram_monitors_chat_id ON telegram_monitors (chat_id)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_telegram_monitors_enabled ON telegram_monitors (enabled)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_telegram_monitors_telegram_config_id ON telegram_monitors (telegram_config_id)")
    op.execute("CREATE INDEX IF NOT EXISTS ix_telegram_monitors_user_id ON telegram_monitors (user_id)")


def downgrade() -> None:
    # No-op: the previous migration owns these objects. This repair migration only
    # makes production schema convergence explicit.
    pass
