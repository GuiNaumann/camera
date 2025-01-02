CREATE TABLE camera
(
    id              SERIAL PRIMARY KEY,
    id_user         INT                                 NOT NULL,
    name            VARCHAR(100)                        NOT NULL,
    description     TEXT                                NOT NULL,
    ip_address      VARCHAR(50)                         NOT NULL,
    port            INT                                 NOT NULL,
    username        VARCHAR(100)                        NOT NULL,
    password        VARCHAR(100)                        NOT NULL,
    stream_path     TEXT                                NOT NULL,
    camera_type     VARCHAR(50)                         NOT NULL,
    is_active       BOOLEAN   DEFAULT TRUE,
    status_code     INT       DEFAULT 0                 NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_camera FOREIGN KEY (id_user) REFERENCES "user" (id)
);
