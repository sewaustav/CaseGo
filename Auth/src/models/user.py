from datetime import datetime
from typing import Optional
from sqlalchemy import String, Integer, Boolean, func
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.dialects.postgresql import TIMESTAMP

from .base import MyBaseModel


class User(MyBaseModel):
	"""Модель пользователя"""
	__tablename__ = "users"

	id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
	username: Mapped[str] = mapped_column(
		String(50),
		unique=True,
		nullable=False,
		index=True
	)
	email: Mapped[str] = mapped_column(
		String(100),
		unique=True,
		nullable=False,
		index=True
	)
	hashed_password: Mapped[str] = mapped_column(String(255), nullable=False)
	
	role: Mapped[int] = mapped_column(
        Integer, 
        default=1, 
        server_default="1", 
        nullable=False
    )

	is_active: Mapped[bool] = mapped_column(Boolean, default=True, server_default="true")
	is_verified: Mapped[bool] = mapped_column(Boolean, default=False, server_default="false")

	created_at: Mapped[datetime] = mapped_column(
		TIMESTAMP(timezone=True),
		server_default=func.now(),
		nullable=False
	)
	updated_at: Mapped[datetime] = mapped_column(
		TIMESTAMP(timezone=True),
		server_default=func.now(),
		onupdate=func.now(),
		nullable=False
	)

	last_login_at: Mapped[Optional[datetime]] = mapped_column(
		TIMESTAMP(timezone=True),
		nullable=True
	)
	email_verified_at: Mapped[Optional[datetime]] = mapped_column(
		TIMESTAMP(timezone=True),
		nullable=True
	)

	deleted_at: Mapped[Optional[datetime]] = mapped_column(
		TIMESTAMP(timezone=True),
		nullable=True
	)