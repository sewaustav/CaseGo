CREATE TABLE levels
(
    id      SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    xp      INT    NOT NULL DEFAULT 0,
    streak  INT    NOT NULL DEFAULT 0,
    level   INT    NOT NULL DEFAULT 1,
    last_active TIMESTAMP WITH TIME ZONE
)

CREATE INDEX idx_levels_user_id ON levels (user_id);