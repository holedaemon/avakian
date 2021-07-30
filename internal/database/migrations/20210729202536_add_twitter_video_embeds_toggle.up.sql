BEGIN;

ALTER TABLE guilds ADD COLUMN embed_twitter_videos boolean DEFAULT 'false' NOT NULL;

COMMIT;