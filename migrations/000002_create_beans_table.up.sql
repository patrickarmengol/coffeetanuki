CREATE TYPE roast_level_enum AS ENUM ('light', 'medium-light', 'medium', 'medium-dark', 'dark');

CREATE TABLE IF NOT EXISTS beans (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    roast_level roast_level_enum NOT NULL,
    roaster_id bigint NOT NULL REFERENCES roasters (id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

