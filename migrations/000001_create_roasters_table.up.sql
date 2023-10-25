CREATE TABLE IF NOT EXISTS roasters (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    location text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);
