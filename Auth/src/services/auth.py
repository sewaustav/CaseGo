import sys
from passlib.context import CryptContext
from typing import Optional

from fastapi import HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession

from models.user import User
from repositories.user import get_user_by_login

pwd_context = CryptContext(schemes=["argon2"], deprecated="auto")


# ======== Пароли ========
def verify_password(plain_password: str, hashed_password: str) -> bool:
	"""Проверка пароля"""
	return pwd_context.verify(plain_password, hashed_password)


def get_password_hash(password: str) -> str:
	"""Хэширование пароля"""
	return pwd_context.hash(password)


# ======== Аутентификация ========
async def authenticate_user(login: str, db: AsyncSession, password: str | None) -> Optional[User]:
	user = await get_user_by_login(login=login, db=db)
	if password is not None and (not user or not verify_password(password, str(user.hashed_password))):
		raise HTTPException(
			status_code=status.HTTP_401_UNAUTHORIZED,
			detail="Incorrect username/email or password",
			headers={"WWW-Authenticate": "Bearer"},
		)
	return user

