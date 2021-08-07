BEGIN;

CREATE TABLE custom_commands (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    guild_snowflake text NOT NULL,

    trigger text NOT NULL,
    body text NOT NULL,

    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),

    UNIQUE (trigger, guild_snowflake)
);

COMMIT;