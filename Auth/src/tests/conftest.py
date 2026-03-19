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

from main import app


@pytest_asyncio.fixture
async def client():
    async with AsyncClient(transport=ASGITransport(app=app), base_url="http://test") as ac:
        yield ac
