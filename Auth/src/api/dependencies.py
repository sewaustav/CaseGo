from typing import Annotated, Optional

from fastapi import Depends, status
from fastapi.security import OAuth2PasswordBearer
from Auth.src.services.google_validation import GoogleOAUTH
from core.config import get_settings

from models.user import User
from services.token import decode_token
from repositories.user import get_user_by_id

from typing import AsyncGenerator

from fastapi import HTTPException, Request
from sqlalchemy.ext.asyncio import AsyncSession

# from rabbit.producer import RabbitMQClient
from core.redis_client import RedisClient

settings = get_settings()
oauth2_scheme = OAuth2PasswordBearer(tokenUrl=f"{settings.API_V1_PREFIX}/auth/token")
google_oauth = GoogleOAUTH(client_id=settings.GOOGLE_CLIENT_ID)


# async def get_rabbitmq(request: Request) -> RabbitMQClient:
# 	"""Зависимость для RabbitMQ клиента."""
# 	rabbitmq = request.state.rabbitmq
# 	if not rabbitmq:
# 		raise HTTPException(status_code=500, detail="RabbitMQ client not initialized")
# 	return rabbitmq


async def get_redis_client(request: Request) -> RedisClient:
	"""Зависимость для redis клиента."""
	redis_client = request.state.redis_client
	if not redis_client:
		raise HTTPException(status_code=500, detail="Redis client not initialized")
	return redis_client


async def get_db_session(request: Request) -> AsyncGenerator[AsyncSession, None]:
	"""Получение сессии для работы с бд"""
	async_session_factory = request.state.async_session_factory
	async with async_session_factory() as session:
		try:
			yield session
		finally:
			await session.close()


async def get_user_id_from_token(token: str = Depends(oauth2_scheme)) -> Optional[int]:
	"""Получение user_id из Token"""
	token_data = decode_token(token)
	return int(token_data.sub)


async def get_user_by_token(
	token: Annotated[str, Depends(oauth2_scheme)],
	db: Annotated[AsyncSession, Depends(get_db_session)],
) -> Optional[User]:
	"""Получение текущего пользователя по Token"""
	token_data = decode_token(token)

	user = await get_user_by_id(user_id=int(token_data.sub), db=db)
	if user is None:
		raise HTTPException(
			status_code=status.HTTP_401_UNAUTHORIZED,
			detail="Could not validate credentials",
			headers={"WWW-Authenticate": "Bearer"},
		)
	return user
