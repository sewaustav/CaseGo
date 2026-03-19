"""Интеграционные тесты для Auth/User эндпоинтов.

Каждый тест генерирует уникальные данные через UUID, поэтому тесты
не конфликтуют друг с другом и не требуют очистки БД.
"""
import uuid

from httpx import AsyncClient


def make_user() -> dict:
    uid = uuid.uuid4().hex[:8]
    return {
        "username": f"user_{uid}",
        "email": f"user_{uid}@test.com",
        "password": "TestPass123!",
    }


# ---------------------------------------------------------------------------
# Register
# ---------------------------------------------------------------------------


async def test_register_success(client: AsyncClient):
    data = make_user()
    r = await client.post("/api/v1/auth/register", json=data)
    assert r.status_code == 201
    body = r.json()
    assert body["email"] == data["email"]
    assert body["username"] == data["username"]
    assert "id" in body


async def test_register_duplicate_email_returns_409(client: AsyncClient):
    data = make_user()
    await client.post("/api/v1/auth/register", json=data)

    duplicate = {**data, "username": f"other_{uuid.uuid4().hex[:8]}"}
    r = await client.post("/api/v1/auth/register", json=duplicate)
    assert r.status_code == 409


async def test_register_duplicate_username_returns_409(client: AsyncClient):
    data = make_user()
    await client.post("/api/v1/auth/register", json=data)

    duplicate = {**data, "email": f"other_{uuid.uuid4().hex[:8]}@test.com"}
    r = await client.post("/api/v1/auth/register", json=duplicate)
    assert r.status_code == 409


# ---------------------------------------------------------------------------
# Login
# ---------------------------------------------------------------------------


async def test_login_returns_token_pair(client: AsyncClient):
    data = make_user()
    await client.post("/api/v1/auth/register", json=data)

    r = await client.post(
        "/api/v1/auth/token",
        data={"username": data["username"], "password": data["password"]},
    )
    assert r.status_code == 200
    body = r.json()
    assert "access_token" in body
    assert "refresh_token" in body
    assert body["token_type"] == "Bearer"


async def test_login_wrong_password_returns_401(client: AsyncClient):
    data = make_user()
    await client.post("/api/v1/auth/register", json=data)

    r = await client.post(
        "/api/v1/auth/token",
        data={"username": data["username"], "password": "wrongpassword"},
    )
    assert r.status_code == 401


async def test_login_nonexistent_user_returns_401(client: AsyncClient):
    r = await client.post(
        "/api/v1/auth/token",
        data={"username": "ghost_user_xyz123", "password": "anypassword"},
    )
    assert r.status_code == 401


# ---------------------------------------------------------------------------
# Refresh
# ---------------------------------------------------------------------------


async def test_refresh_token_returns_new_access_token(client: AsyncClient):
    data = make_user()
    await client.post("/api/v1/auth/register", json=data)
    login = await client.post(
        "/api/v1/auth/token",
        data={"username": data["username"], "password": data["password"]},
    )
    tokens = login.json()

    r = await client.post("/api/v1/auth/refresh", json={"refresh_token": tokens["refresh_token"]})
    assert r.status_code == 200
    assert "access_token" in r.json()


async def test_refresh_invalid_token_returns_401(client: AsyncClient):
    r = await client.post("/api/v1/auth/refresh", json={"refresh_token": "not.a.valid.token"})
    assert r.status_code == 401


# ---------------------------------------------------------------------------
# /users/me
# ---------------------------------------------------------------------------


async def test_get_me_returns_current_user(client: AsyncClient):
    data = make_user()
    await client.post("/api/v1/auth/register", json=data)
    login = await client.post(
        "/api/v1/auth/token",
        data={"username": data["username"], "password": data["password"]},
    )
    tokens = login.json()

    r = await client.get(
        "/api/v1/users/me",
        headers={"Authorization": f"Bearer {tokens['access_token']}"},
    )
    assert r.status_code == 200
    assert r.json()["email"] == data["email"]


async def test_get_me_without_token_returns_401(client: AsyncClient):
    r = await client.get("/api/v1/users/me")
    assert r.status_code == 401
