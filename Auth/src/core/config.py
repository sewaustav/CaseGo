from functools import lru_cache
from pathlib import Path
from typing import Literal

from pydantic import Field, RedisDsn, SecretStr, computed_field
from pydantic import AliasChoices
from pydantic_settings import BaseSettings, SettingsConfigDict

CURRENT_FILE = Path(__file__).resolve()


def find_env_file(start_path: Path) -> Path:
    for parent in start_path.parents:
        if (parent / ".env").exists():
            return parent / ".env"
    return start_path.parent.parent.parent / ".env"


ENV_PATH = find_env_file(CURRENT_FILE)
BASE_DIR = ENV_PATH.parent


class Settings(BaseSettings):
    """App settings"""

    model_config = SettingsConfigDict(
        env_file=ENV_PATH,
        env_file_encoding="utf-8",
        case_sensitive=True,
        extra="allow",
    )
    

    # App
    DEBUG: bool = False
    ENVIRONMENT: Literal["development", "testing", "production"] = "development"
    API_V1_PREFIX: str = "/api/v1"
    PROJECT_NAME: str | None = "Auth service"

    # Database
    POSTGRES_HOST: str = Field(...)
    POSTGRES_PORT: int = Field(
        validation_alias=AliasChoices("POSTGRES_PORT_FORWARD", "POSTGRES_PORT")
    )
    POSTGRES_USER: str = Field(...)
    POSTGRES_PASSWORD: SecretStr = Field(...)
    POSTGRES_DB: str = Field(...)

    # Redis
    REDIS_HOST: str = Field(default="localhost")
    REDIS_PORT: int = Field(validation_alias=AliasChoices("REDIS_PORT_FORWARD", "REDIS_PORT"))
    REDIS_DB: int = Field(...)
    REDIS_PASSWORD: SecretStr | None = None

    # CORS
    CORS_ORIGINS: str = Field(...)

    # Auth
    JWT_ALG: str = Field(...)
    JWT_PRIVATE_KEY_PATH: Path = Path(BASE_DIR / "keys/private.pem")
    JWT_PUBLIC_KEY_PATH: Path = Path(BASE_DIR / "keys/public.pem")
    ACCESS_TOKEN_EXPIRE_SECONDS: int = Field(default=15)
    REFRESH_TOKEN_EXPIRE_SECONDS: int = Field(default=60)

    @property
    def postgres_async_url(self) -> str:
        return f"postgresql+asyncpg://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD.get_secret_value()}@{self.POSTGRES_HOST}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"

    @property
    def postgres_sync_url(self) -> str:
        return f"postgresql+psycopg2://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD.get_secret_value()}@{self.POSTGRES_HOST}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"

    @computed_field
    @property
    def redis_url(self) -> str:
        redis_dsn = RedisDsn.build(
            scheme="redis",
            username=None,
            password=self.REDIS_PASSWORD.get_secret_value()
            if self.REDIS_PASSWORD
            else None,
            host=self.REDIS_HOST,
            port=self.REDIS_PORT,
            path=f"{self.REDIS_DB}",
        )
        return str(redis_dsn)


@lru_cache
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
