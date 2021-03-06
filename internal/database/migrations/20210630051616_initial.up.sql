BEGIN;

CREATE TABLE guilds (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    guild_snowflake text NOT NULL UNIQUE,

    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE prefixes (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    guild_snowflake text NOT NULL REFERENCES guilds (guild_snowflake),
    prefix text NOT NULL,

    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),

    UNIQUE (guild_snowflake, prefix)
);

CREATE TABLE pronouns (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    guild_snowflake text NOT NULL REFERENCES guilds (guild_snowflake),

    pronoun text NOT NULL,
    role_snowflake text NOT NULL UNIQUE,
    
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),

    UNIQUE (guild_snowflake, pronoun)
);

COMMIT;