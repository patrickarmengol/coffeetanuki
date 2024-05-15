CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name text NOT NULL
)

CREATE TABLE IF NOT EXISTS users_roles (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    role_id bigint NOT NULL REFERENCES roles ON DELETE CASCADE;
)

INSERT INTO roles (name)
VALUES
    ('guest'),
    ('user'),
    ('moderator'),
    ('admin');
