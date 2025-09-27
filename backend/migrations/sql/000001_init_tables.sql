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

CREATE TABLE IF NOT EXISTS offer
(
    id            UUID NOT NULL PRIMARY KEY,
    user_id       UUID,
    hotel_id      UUID,
    expiration_at TIMESTAMP,
    location_id   UUID,
    used          BOOLEAN,
    task          TEXT,
    CONSTRAINT fk_application_location FOREIGN KEY (location_id) REFERENCES location (id),
    CONSTRAINT fk_application_user FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE TABLE IF NOT EXISTS report
(
    id            UUID NOT NULL PRIMARY KEY,
    expiration_at TIMESTAMP,
    status        VARCHAR(16),
    text          TEXT
);

CREATE TABLE IF NOT EXISTS application
(
    id       UUID NOT NULL PRIMARY KEY,
    user_id  UUID,
    offer_id UUID,
    status   VARCHAR(16),
    CONSTRAINT fk_application_user FOREIGN KEY (user_id) REFERENCES "user" (id),
    CONSTRAINT fk_application_offer FOREIGN KEY (offer_id) REFERENCES offer (id)
);

CREATE TABLE IF NOT EXISTS photo
(
    id UUID NOT NULL PRIMARY KEY,
    report_id UUID,
    s3_link TEXT,
    CONSTRAINT fk_application_report FOREIGN KEY (report_id) REFERENCES report(id)
)
