-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- Table: User
CREATE TABLE IF NOT EXISTS users
(
    id character varying(255) NOT NULL
);

-- Table: Videos
CREATE TABLE IF NOT EXISTS videos
(
    id character varying(255) NOT NULL,
    user_id character varying(255) NOT NULL,
    url character varying(255) NOT NULL,
    duration integer NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS videos_created_at
    ON videos USING btree (created_at);


-- Table: Annotations
CREATE TABLE IF NOT EXISTS annotations
(
    id character varying(255) NOT NULL,
    video_id character varying(255) NOT NULL,
    user_id character varying(255) NOT NULL,
    start_time integer NOT NULL,
    end_time integer NOT NULL,
    type integer NOT NULL,
    notes text,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS videos CASCADE;
DROP TABLE IF EXISTS annotations CASCADE;
