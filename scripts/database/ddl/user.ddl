-- Tabela 'user'
CREATE TABLE "user"
(
    id                   SERIAL PRIMARY KEY,
    user_name            VARCHAR(60)                            NOT NULL,
    user_document        VARCHAR(11),
    user_email           VARCHAR(120) DEFAULT ''                NOT NULL,
    user_picture         TEXT,
    password_hash        TEXT                                   NOT NULL,
    password_modified_at TIMESTAMP,
    street               VARCHAR(70),
    address_number       INT,
    district             VARCHAR(70),
    zip_code             VARCHAR(8),
    federal_unit         INT,
    biography            TEXT,
    status_code          SMALLINT     DEFAULT 0                 NOT NULL,
    badge_code           INT,
    city                 VARCHAR(70),
    social_name          VARCHAR(60),
    is_active            SMALLINT     DEFAULT 0,
    modified_at          TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at           TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    terms_accepted       INT          DEFAULT 0                 NOT NULL,
    CONSTRAINT fk_user_federal_unit FOREIGN KEY (federal_unit) REFERENCES states (id)
);

CREATE INDEX idx_federal_unit ON "user" (federal_unit);
