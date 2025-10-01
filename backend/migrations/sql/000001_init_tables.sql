CREATE TABLE IF NOT EXISTS "user"
(
    id             UUID    NOT NULL PRIMARY KEY,
    ostrovok_login TEXT    NOT NULL UNIQUE,
    password_hash  TEXT    NOT NULL,
    is_admin       BOOLEAN NOT NULL,
    app_limit      INTEGER NOT NULL DEFAULT 5,
    rating         INTEGER NOT NULL DEFAULT 5
);

CREATE TABLE IF NOT EXISTS location
(
    id   UUID NOT NULL PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS room
(
    id   UUID NOT NULL PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS hotel
(
    id          UUID NOT NULL PRIMARY KEY,
    name        TEXT UNIQUE,
    location_id UUID NOT NULL REFERENCES location (id)
);

CREATE TABLE IF NOT EXISTS offer
(
    id                 UUID      NOT NULL PRIMARY KEY,
    hotel_id           UUID      NOT NULL REFERENCES hotel (id),
    room_id            UUID      NOT NULL REFERENCES room (id),
    check_in_at        TIMESTAMP WITH TIME ZONE NOT NULL,
    check_out_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    status             VARCHAR(20) DEFAULT 'created',
    expiration_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    participants_limit NUMERIC   NOT NULL,
    task               TEXT,
    CONSTRAINT uq_offer UNIQUE (hotel_id, room_id, check_in_at, check_out_at)
);

CREATE TABLE IF NOT EXISTS application
(
    id       UUID NOT NULL PRIMARY KEY,
    user_id  UUID NOT NULL REFERENCES "user" (id),
    offer_id UUID NOT NULL REFERENCES offer (id),
    status   VARCHAR(16),

    CONSTRAINT uq_application UNIQUE (user_id, offer_id)
);

CREATE TABLE IF NOT EXISTS report
(
    id             UUID NOT NULL PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES application (id),
    expiration_at  TIMESTAMP WITH TIME ZONE,
    status         VARCHAR(16),
    text           TEXT
);

CREATE TABLE IF NOT EXISTS photo
(
    id        UUID NOT NULL PRIMARY KEY,
    report_id UUID REFERENCES report (id),
    s3_link   TEXT UNIQUE
);
