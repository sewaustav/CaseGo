import sys
from typing import TypedDict, AsyncIterator, cast

from fastapi import FastAPI
from contextlib import asynccontextmanager

from sqlalchemy.ext.asyncio import async_sessionmaker, AsyncSession
from sqlalchemy.sql.expression import text

from .database import get_async_engine, get_async_session_factory


class AppState(TypedDict):
    async_session_factory: async_sessionmaker[AsyncSession]


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[AppState]:
    async_engine = get_async_engine()
    async_session_factory = get_async_session_factory(async_engine)

    try:
        async with async_session_factory() as test_session:
            await test_session.execute(text("SELECT 1"))
            print("✓ Database connection successful", file=sys.stderr)
    except Exception as e:
        print(f"✗ Database connection failed: {e}", file=sys.stderr)

    yield cast(
        AppState,
        {
            "async_session_factory": async_session_factory,
        },
    )

    await async_engine.dispose()
