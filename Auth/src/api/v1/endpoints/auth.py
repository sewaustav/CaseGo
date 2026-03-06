from fastapi import APIRouter, Depends
from fastapi.exceptions import HTTPException

from schemas.token import TokenPair, RefreshTokenRequest, AuthRequest

from Auth.src.repositories.user import get_user_by_login
from ...dependencies import get_db_session, google_oauth
from services.token import create_token_pair, refresh_access_token
from fastapi.security import OAuth2PasswordRequestForm
from services.auth import authenticate_user
from typing import Annotated

from sqlalchemy.ext.asyncio import AsyncSession

from schemas.user import UserResponse, UserRegister
from services.user import register_user

router = APIRouter()


@router.post("/auth/google")
async def google_auth(data: AuthRequest, db: Annotated[AsyncSession, Depends(get_db_session)]):
	user_data = google_oauth.verify_google_token(data.id_token)

	if not user_data:
		raise HTTPException(
			status_code=401,
			detail="Invalid Google Token"
		)

	# Теперь у тебя есть данные юзера.
	# user_data['sub'] — это уникальный ID пользователя в Google (не меняется)
	# user_data['email'] — почта

	user = UserRegister(
		email=user_data['email'],
		username=user_data['email'],
		password=None
	)

	usr = await get_user_by_login(user_data['email'], db)
	if usr is not None:
		await register_user(user, db)
	new_user = await authenticate_user(user_data['email'], db, None)
	if user is None:
		raise HTTPException(status_code=401, detail="user is none")
	return new_user

	# return {
	# 	"status": "ok",
	# 	"user": {
	# 		"email": user_data.get("email"),
	# 		"name": user_data.get("name")
	# 	},
	# 	"message": "Welcome to the club, buddy"
	# }

@router.post("/register", response_model=UserResponse, status_code=201)
async def register_user_endpoint(body: UserRegister, db: Annotated[AsyncSession, Depends(get_db_session)]):
	new_user = await register_user(body, db)
	return UserResponse.model_validate(new_user)


@router.post("/token")
async def login_for_access_token_endpoint(form_data: Annotated[OAuth2PasswordRequestForm, Depends()],
								 db: Annotated[AsyncSession, Depends(get_db_session)]) -> TokenPair:
	"""Ручка для входа пользователя по username/email и паролю с созданием access_token и обновлением информации о сессии"""
	user = await authenticate_user(form_data.username, db, form_data.password)
	if user is None:
		raise HTTPException(status_code=401, detail="user is none")
	tokens_data = create_token_pair(user_id=user.id, additional_data={"user_role": str(user.role)})

	return tokens_data


@router.post("/refresh")
async def refresh_token_endpoint(body: RefreshTokenRequest,
							   db: Annotated[AsyncSession, Depends(get_db_session)]
							   ) -> TokenPair:
	"""
	Обновление access токена с помощью refresh токена
	"""
	tokens_data = await refresh_access_token(
		refresh_token=body.refresh_token,
		additional_data=None,
	)

	return tokens_data
