CREATE TABLE IF NOT EXISTS "user"
(
    id             UUID    NOT NULL PRIMARY KEY,
    ostrovok_login TEXT    NOT NULL,
    password_hash  TEXT    NOT NULL,
    is_admin       BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS location
(
    id   UUID NOT NULL PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS room_type
(
    id   UUID NOT NULL PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS hotel
(
    id          UUID NOT NULL PRIMARY KEY,
    name        TEXT,
    location_id UUID NOT NULL REFERENCES location (id)
);

CREATE TABLE IF NOT EXISTS slay_slot
(
    id              UUID PRIMARY KEY,
    hotel_id        UUID      NOT NULL REFERENCES hotel (id),
    room_type_id    UUID      NOT NULL REFERENCES room_type (id),
    check_in_at     TIMESTAMP NOT NULL,
    check_out_at    TIMESTAMP NOT NULL,
    available_count INTEGER DEFAULT 1,
    CONSTRAINT uq_slot UNIQUE (hotel_id, room_type_id, check_in_at, check_out_at)
);

CREATE TABLE IF NOT EXISTS offer
(
    id            UUID NOT NULL PRIMARY KEY,
    slay_slot_id  UUID NOT NULL REFERENCES slay_slot (id),
    expiration_at TIMESTAMP,
    task          TEXT
);

CREATE TABLE IF NOT EXISTS application
(
    id       UUID NOT NULL PRIMARY KEY,
    user_id  UUID NOT NULL REFERENCES "user" (id),
    offer_id UUID NOT NULL REFERENCES offer (id),
    status   VARCHAR(16)
);

CREATE TABLE IF NOT EXISTS report
(
    id            UUID NOT NULL PRIMARY KEY,
    expiration_at TIMESTAMP,
    status        VARCHAR(16),
    text          TEXT
);

CREATE TABLE IF NOT EXISTS photo
(
    id        UUID NOT NULL PRIMARY KEY,
    report_id UUID REFERENCES report (id),
    s3_link   TEXT
);
