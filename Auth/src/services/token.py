import sys
import uuid
from datetime import datetime, timedelta, timezone
from typing import Any, Dict, Optional

import jwt
from core.config import get_settings
from fastapi import HTTPException, status
from schemas.token import TokenPair, TokenPayload
from sqlalchemy.orm import Mapped

settings = get_settings()


def create_token(
    subject: str,
    token_type: str = "access",
    expires_delta: Optional[timedelta] = None,
    additional_data: Optional[Dict[str, Any]] = None,
    not_before: Optional[datetime] = None,
) -> str:
    """
    Создание JWT токена

    Args:
    subject: ID пользователя (строка)
    token_type: Тип токена ("access" или "refresh")
    expires_delta: Время жизни токена
    additional_data: Дополнительные данные для payload
    not_before: Не использовать токен до этого времени

    Returns:
    token: - токен
    """

    if not subject:
        raise ValueError("Subject cannot be empty")

    if expires_delta is None:
        if token_type == "access":
            expires_delta = timedelta(seconds=settings.ACCESS_TOKEN_EXPIRE_SECONDS)
        elif token_type == "refresh":
            expires_delta = timedelta(seconds=settings.REFRESH_TOKEN_EXPIRE_SECONDS)
        else:
            expires_delta = timedelta(hours=1)

    now = datetime.now(timezone.utc)
    expire = now + expires_delta

    if not_before:
        if not_before.tzinfo is None:
            not_before = not_before.replace(tzinfo=timezone.utc)
    else:
        not_before = now

    token_data = TokenPayload(
        sub=subject, exp=expire, iat=now, nbf=not_before, jti=str(uuid.uuid4())
    )

    payload = token_data.model_dump()
    if additional_data:
        payload.update(additional_data)

    private_key = settings.JWT_PRIVATE_KEY_PATH.read_text()


    token = jwt.encode(payload, private_key, algorithm=settings.JWT_ALG)

    return token


def decode_token(token: str) -> TokenPayload:
    """Декодирование access токена"""
    try:
        public_key = settings.JWT_PUBLIC_KEY_PATH.read_text()

        payload = jwt.decode(token, public_key, algorithms=[settings.JWT_ALG])

    except jwt.ExpiredSignatureError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token expired",
            headers={"WWW-Authenticate": "Bearer"},
        )
    except jwt.InvalidTokenError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token",
            headers={"WWW-Authenticate": "Bearer"},
        )

    required_keys = {"sub", "exp", "nbf", "iat", "jti"}
    if not required_keys <= payload.keys():
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token structure",
            headers={"WWW-Authenticate": "Bearer"},
        )

    return TokenPayload.model_validate(payload)


def create_token_pair(
    user_id: int | Mapped[int], additional_data: Optional[dict[str, Any]] = None
) -> TokenPair:
    """Создаёт access + refresh токены"""
    
    payload_data = additional_data.copy() if additional_data else {}
    payload_data.update({
        "iss": "auth", # Issuer
        "aud": "all",  # Audience
    })
        
    access_token = create_token(subject=str(user_id), additional_data=payload_data)

    refresh_token = create_token(
        subject=str(user_id), additional_data=payload_data, token_type="refresh"
    )

    return TokenPair(
        access_token=access_token,
        refresh_token=refresh_token,
        expires_in=settings.ACCESS_TOKEN_EXPIRE_SECONDS,
        token_type="Bearer",
    )


async def refresh_access_token(
    refresh_token: str,
    additional_data: Optional[Dict[str, Any]] = None,
) -> TokenPair:
    """Обновление access токена"""
    payload = decode_token(refresh_token)
    user_id = int(payload.sub)

    token_pair = create_token_pair(user_id=user_id, additional_data=additional_data)

    return token_pair
