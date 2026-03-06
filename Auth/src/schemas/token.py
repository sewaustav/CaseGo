from datetime import datetime
from pydantic import BaseModel

class AuthRequest(BaseModel):
    id_token: str

class TokenPair(BaseModel):
	"""
	Ответ при логине / refresh
	"""
	access_token: str
	refresh_token: str
	expires_in: int
	token_type: str


class TokenPayload(BaseModel):
	"""
	JWT data
	"""
	sub: str  # user_id
	jti: str  # token id
	exp: datetime
	nbf: datetime
	iat: datetime
# role
# created_by service


class RefreshTokenPayload(BaseModel):
	"""
	Payload refresh токена
	"""
	sub: str
	exp: datetime


class RefreshTokenRequest(BaseModel):
	refresh_token: str
