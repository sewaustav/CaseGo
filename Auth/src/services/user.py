from fastapi import HTTPException, status
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession

from models.user import User
from schemas.user import UserCreate, UserRegister
from .auth import get_password_hash
from repositories.user import create_user, get_user_by_email, get_user_by_username


async def register_user(body: UserRegister, db: AsyncSession) -> User:
    	
    if await get_user_by_email(body.email, db):
        raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="Email already registered")
    if await get_user_by_username(body.username, db):
        raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="Username already registered")

    data = body.model_dump()
    data["hashed_password"] = get_password_hash(data.pop("password"))
    user_create = UserCreate(**data)
    try:
        user = await create_user(user_create, db)
        await db.commit()
        await db.refresh(user)
        return user
    except IntegrityError:
        await db.rollback()
        raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="User already exists")

