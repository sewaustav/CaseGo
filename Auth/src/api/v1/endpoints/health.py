from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import text
from datetime import datetime

from ...dependencies import get_db_session

router = APIRouter(tags=["health"])


@router.get("/", summary="Health check", response_model=dict)
async def health_check(
    db: AsyncSession = Depends(get_db_session),
):
    health_data = {
        "status": "healthy",
        "timestamp": datetime.utcnow().isoformat(),
        "services": {},
    }

    try:
        await db.execute(text("SELECT 1"))
        await db.commit()
        health_data["services"]["postgresql"] = {"status": "healthy"}
    except Exception as e:
        health_data["status"] = "unhealthy"
        health_data["services"]["postgresql"] = {"status": "unhealthy", "message": str(e)}

    health_data["services"]["app"] = {"status": "healthy"}

    return health_data
