-- 1. Таблица кейсов (Cases)
CREATE TABLE cases (
                       id BIGSERIAL PRIMARY KEY,
                       topic TEXT NOT NULL,
                       category INTEGER NOT NULL,
                       is_generated BOOLEAN NOT NULL DEFAULT FALSE,
                       description TEXT NOT NULL,
                       first_question TEXT NOT NULL,
                       creator BIGINT NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 2. Таблица диалогов (Dialog)
CREATE TABLE dialogs (
                         id BIGSERIAL PRIMARY KEY,
                         case_id BIGINT NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
                         user_id BIGINT NOT NULL,
                         model_name VARCHAR(255),
                         started_at TIMESTAMP WITH TIME ZONE,
                         ended_at TIMESTAMP WITH TIME ZONE
);

-- 3. Таблица шагов взаимодействия (Interaction)
CREATE TABLE interactions (
                              id BIGSERIAL PRIMARY KEY,
                              dialog_id BIGINT NOT NULL REFERENCES dialogs(id) ON DELETE CASCADE,
                              step INTEGER NOT NULL,
                              question TEXT NOT NULL,
                              answer TEXT NOT NULL,
                              tokens_used INTEGER NOT NULL DEFAULT 0,
                              created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 4. Таблица результатов (Result)
CREATE TABLE results (
                         id BIGSERIAL PRIMARY KEY,
                         case_id BIGINT NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
                         user_id BIGINT NOT NULL,
                         dialog_id BIGINT NOT NULL REFERENCES dialogs(id) ON DELETE CASCADE,
                         steps_count INTEGER NOT NULL DEFAULT 0,
                         tokens_used INTEGER NOT NULL DEFAULT 0,
                         finished_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                         assertiveness DOUBLE PRECISION NOT NULL DEFAULT 0,
                         empathy DOUBLE PRECISION NOT NULL DEFAULT 0,
                         clarity_communication DOUBLE PRECISION NOT NULL DEFAULT 0,
                         resistance DOUBLE PRECISION NOT NULL DEFAULT 0,
                         eloquence DOUBLE PRECISION NOT NULL DEFAULT 0,
                         initiative DOUBLE PRECISION NOT NULL DEFAULT 0
);

-- Индексы для ускорения выборок по внешним ключам
CREATE INDEX idx_dialogs_case_id ON dialogs(case_id);
CREATE INDEX idx_interactions_dialog_id ON interactions(dialog_id);
CREATE INDEX idx_results_dialog_id ON results(dialog_id);