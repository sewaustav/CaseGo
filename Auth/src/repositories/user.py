from typing import Optional

from pydantic import EmailStr
from sqlalchemy import select, or_
from sqlalchemy.ext.asyncio import AsyncSession

from models.user import User
from schemas.user import UserCreate

async def get_user_by_id(user_id: int, db: AsyncSession) -> Optional[User]:
    """Получение пользователя по id из БД"""
    stmt = select(User).where(User.id == user_id)
    result = await db.execute(stmt)
    user = result.scalar_one_or_none()
    return user

async def get_user_by_login(login: str,
                   db: AsyncSession) -> Optional[User]:
    """Получение пользователя по username(email) из БД"""
    stmt = select(User).where(or_(User.username == login, User.email == login))
    result = await db.execute(stmt)
    user = result.scalar_one_or_none()
    return user

async def get_user_by_email(email: EmailStr,
                   db: AsyncSession) -> Optional[User]:
    """Получение пользователя по email из БД"""
    stmt = select(User).where(User.email == email)
    result = await db.execute(stmt)
    user = result.scalar_one_or_none()
    return user

async def get_user_by_username(username: str,
                   db: AsyncSession) -> Optional[User]:
    """Получение пользователя по username из БД"""
    stmt = select(User).where(User.username == username)
    result = await db.execute(stmt)
    user = result.scalar_one_or_none()
    return user


async def create_user(user_create: UserCreate, db: AsyncSession) -> User:
    """Добавляет объект user в сессию, но не делает commit.
       Возвращает ORM-объект; refresh нужно сделать вне или внутри транзакции."""
    user = User(**user_create.model_dump())
    db.add(user)
    return user
