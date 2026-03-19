import os

# Должны быть установлены ДО импорта приложения, т.к. settings кэшируется при первом импорте
os.environ.setdefault("REDIS_DB", "0")
os.environ.setdefault("CORS_ORIGINS", "http://localhost")
os.environ.setdefault("JWT_ALG", "RS256")
os.environ.setdefault("GOOGLE_CLIENT_ID", "test-google-client-id")
os.environ.setdefault("ACCESS_TOKEN_EXPIRE_SECONDS", "300")
os.environ.setdefault("REFRESH_TOKEN_EXPIRE_SECONDS", "3600")

import pytest_asyncio
from httpx import AsyncClient, ASGITransport
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker

from core.config import get_settings
from main import app
from api.dependencies import get_db_session

settings = get_settings()


@pytest_asyncio.fixture
async def client():
    engine = create_async_engine(settings.postgres_async_url)
    session_factory = async_sessionmaker(engine, expire_on_commit=False)

    async def override_db():
        async with session_factory() as session:
            yield session

    app.dependency_overrides[get_db_session] = override_db

    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
        yield ac

    app.dependency_overrides.clear()
    await engine.dispose()
