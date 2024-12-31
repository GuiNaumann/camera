-- Tabela 'type_user'
CREATE TABLE type_user
(
    id           SERIAL PRIMARY KEY,
    numeric_code INT                                 NOT NULL,
    label        VARCHAR(30)                         NOT NULL,
    status_code  SMALLINT  DEFAULT 0                 NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT un_user_type_numeric_code UNIQUE (numeric_code)
);