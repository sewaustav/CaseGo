from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException, status

from ...dependencies import get_user_by_token, get_db_session
from models.user import User
from schemas.user import UserResponse
from repositories.user import get_all_users, update_user_role
from sqlalchemy.ext.asyncio import AsyncSession
from pydantic import BaseModel

router = APIRouter()


class RoleUpdate(BaseModel):
    role: int


def _require_admin(current_user: User) -> User:
    if current_user.role != 0:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Admin only")
    return current_user


@router.get("/me", response_model=UserResponse)
async def read_users_me_endpoint(
    current_user: Annotated[User, Depends(get_user_by_token)]
):
    return current_user


@router.get("/", response_model=list[UserResponse])
async def list_users(
    current_user: Annotated[User, Depends(get_user_by_token)],
    db: Annotated[AsyncSession, Depends(get_db_session)],
):
    _require_admin(current_user)
    return await get_all_users(db)


@router.patch("/{user_id}/role", response_model=UserResponse)
async def update_role(
    user_id: int,
    body: RoleUpdate,
    current_user: Annotated[User, Depends(get_user_by_token)],
    db: Annotated[AsyncSession, Depends(get_db_session)],
):
    _require_admin(current_user)
    user = await update_user_role(user_id, body.role, db)
    if user is None:
        raise HTTPException(status_code=404, detail="User not found")
    return user
