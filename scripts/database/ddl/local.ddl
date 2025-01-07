CREATE TABLE local
(
    id          SERIAL PRIMARY KEY,
    id_user     INT                                 NOT NULL,
    name        VARCHAR(100)                        NOT NULL,
    description TEXT                                NOT NULL,
    state       VARCHAR(50)                         NOT NULL,
    city        VARCHAR(50)                         NOT NULL,
    street      VARCHAR(100)                        NOT NULL,
    is_active   BOOLEAN   DEFAULT TRUE,
    status_code     INT       DEFAULT 0                 NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_local FOREIGN KEY (id_user) REFERENCES "user" (id)
);
