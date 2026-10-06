DROP TABLE IF EXISTS platform.outbox_events;
DROP FUNCTION IF EXISTS platform.set_updated_at();
DROP SCHEMA IF EXISTS platform;
DROP EXTENSION IF EXISTS citext;