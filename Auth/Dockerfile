FROM python:3.13.2-slim AS base

WORKDIR /app

ARG UID=1000
ARG GID=1000

RUN groupadd -g "${GID}" python \
    && useradd --create-home --no-log-init -u "${UID}" -g "${GID}" python

RUN chown python:python /app

USER python

FROM base AS app
COPY --chown=python:python requirements.txt ./
RUN pip install --user -r requirements.txt
COPY --chown=python:python . .

ENV PYTHONUNBUFFERED=1 \
  PYTHONPATH="/app:/app/src" \
  PATH="${PATH}:/home/python/.local/bin" \
  USER="python"

WORKDIR /app

CMD ["python", "-m", "uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]
