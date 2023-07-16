-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- Table: User
CREATE TABLE IF NOT EXISTS users
(
    id character varying(255) NOT NULL primary key
);

-- Table: Videos
CREATE TABLE IF NOT EXISTS videos
(
    id character varying(255) NOT NULL primary key,
    user_id character varying(255) NOT NULL,
    url character varying(255) NOT NULL,
    duration integer NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);


-- Table: Annotations
CREATE TABLE IF NOT EXISTS annotations
(
    id character varying(255) NOT NULL primary key,
    video_id character varying(255) NOT NULL references videos(id) on delete cascade,
    user_id character varying(255) NOT NULL,
    start_time integer NOT NULL,
    end_time integer NOT NULL,
    type character varying(255) NOT NULL,
    message text,
    url character varying(255),
    title character varying(255),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS videos CASCADE;
DROP TABLE IF EXISTS annotations CASCADE;
