CREATE TABLE users (
  id          SERIAL PRIMARY KEY,
  email       VARCHAR(254) UNIQUE,
  password    VARCHAR(128),
  facebook_id VARCHAR(128) UNIQUE,
  num         INTEGER NOT NULL DEFAULT 1,
  CONSTRAINT email_or_facebook_id
    CHECK(facebook_id IS NOT NULL OR email IS NOT NULL),
  CONSTRAINT email_and_password CHECK(
    (email IS NOT NULL AND password IS NOT NULL) OR
    (email IS NULL AND password IS NULL)
  )
)

