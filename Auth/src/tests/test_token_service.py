"""Unit-тесты для services/token.py — не требуют БД."""
from datetime import timedelta

import pytest
from fastapi import HTTPException

from services.token import create_token, create_token_pair, decode_token, refresh_access_token

# iss/aud нужны, т.к. decode_token их проверяет
_PAYLOAD = {"iss": "auth", "aud": "all"}


def test_create_and_decode_token():
    token = create_token("42", expires_delta=timedelta(minutes=5), additional_data=_PAYLOAD)
    payload = decode_token(token)
    assert payload.sub == "42"


def test_create_token_empty_subject_raises():
    with pytest.raises(ValueError):
        create_token("")


def test_decode_expired_token_returns_401():
    token = create_token("1", expires_delta=timedelta(seconds=-10), additional_data=_PAYLOAD)
    with pytest.raises(HTTPException) as exc_info:
        decode_token(token)
    assert exc_info.value.status_code == 401


def test_decode_invalid_token_returns_401():
    with pytest.raises(HTTPException) as exc_info:
        decode_token("this.is.not.valid")
    assert exc_info.value.status_code == 401


def test_create_token_pair_has_both_tokens():
    pair = create_token_pair(user_id=1)
    assert pair.access_token
    assert pair.refresh_token
    assert pair.token_type == "Bearer"
    assert pair.expires_in > 0


def test_token_pair_encodes_correct_user_id():
    pair = create_token_pair(user_id=99)
    payload = decode_token(pair.access_token)
    assert payload.sub == "99"


async def test_refresh_access_token_returns_new_pair():
    pair = create_token_pair(user_id=7)
    new_pair = await refresh_access_token(pair.refresh_token)
    payload = decode_token(new_pair.access_token)
    assert payload.sub == "7"


async def test_refresh_with_invalid_token_returns_401():
    with pytest.raises(HTTPException) as exc_info:
        await refresh_access_token("not.a.token")
    assert exc_info.value.status_code == 401
