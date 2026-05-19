-- Supabase-compatible auth.users (minimal columns for login + data import).
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.users (
    instance_id          UUID,
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aud                  TEXT,
    role                 TEXT,
    email                TEXT UNIQUE,
    encrypted_password   TEXT,
    email_confirmed_at   TIMESTAMPTZ,
    invited_at           TIMESTAMPTZ,
    confirmation_token   TEXT,
    confirmation_sent_at TIMESTAMPTZ,
    recovery_token       TEXT,
    recovery_sent_at     TIMESTAMPTZ,
    email_change_token_new TEXT,
    email_change         TEXT,
    email_change_sent_at TIMESTAMPTZ,
    last_sign_in_at      TIMESTAMPTZ,
    raw_app_meta_data    JSONB,
    raw_user_meta_data   JSONB,
    is_super_admin       BOOLEAN,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    phone                TEXT,
    phone_confirmed_at   TIMESTAMPTZ,
    phone_change         TEXT,
    phone_change_token   TEXT,
    phone_change_sent_at TIMESTAMPTZ,
    email_change_token_current TEXT,
    email_change_confirm_status SMALLINT,
    banned_until         TIMESTAMPTZ,
    reauthentication_token TEXT,
    reauthentication_sent_at TIMESTAMPTZ,
    is_sso_user          BOOLEAN DEFAULT FALSE,
    deleted_at           TIMESTAMPTZ,
    is_anonymous         BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS auth_users_email_idx ON auth.users (LOWER(email));
