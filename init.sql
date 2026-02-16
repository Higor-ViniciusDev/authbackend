CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(), --Generator radom uuid in criation
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE, --Unique email validation for login and logout
    password_hash VARCHAR(255) NOT NULL,
    verified      BOOLEAN      NOT NULL DEFAULT FALSE, /* Start false for verification */
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    verified_at   TIMESTAMPTZ
);