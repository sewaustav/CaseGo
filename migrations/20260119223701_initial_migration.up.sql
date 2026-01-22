CREATE TABLE profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    avatar TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT NOT NULL DEFAULT '',
    username VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50),
    sex SMALLINT, 
    profession VARCHAR(255),
    case_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_social_links (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    CONSTRAINT fk_user_profile FOREIGN KEY (user_id) REFERENCES profiles(user_id) ON DELETE CASCADE
);

CREATE TABLE user_purposes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    purpose TEXT NOT NULL,
    CONSTRAINT fk_user_profile_purpose FOREIGN KEY (user_id) REFERENCES profiles(user_id) ON DELETE CASCADE
);


CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_social_links_user_id ON user_social_links(user_id);
CREATE INDEX idx_purposes_user_id ON user_purposes(user_id);