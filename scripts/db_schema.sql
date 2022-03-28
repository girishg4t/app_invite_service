-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS user (
    "id" varchar(36) PRIMARY KEY,
    "username"  varchar(36) NOT NULL,
    "password"  varchar(36) NOT NULL,
    "role" varchar(10) NOT NULL,
    "created_at" timestamp default current_timestamp,
    "updated_at" timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS app_token (
    "id" varchar(36) PRIMARY KEY,
    "username"  varchar(36) NOT NULL,
    "token"  varchar(12) NOT NULL,
    "exp_date"  timestamp NOT NULL,
    "is_active" boolean NOT NULL default true,
    "created_at" timestamp default current_timestamp,
    "updated_at" timestamp default current_timestamp
);

INSERT INTO user (
        id,
        username,
        password,
        role
    )
VALUES (
        '8a2a79d1-7717-4967-97c5-e26c683fcdc6',
        'admin',
        'admin',
        'ADMIN'
    ) ON CONFLICT DO NOTHING;

