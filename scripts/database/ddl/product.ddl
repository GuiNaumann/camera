CREATE TABLE product
(
    id              SERIAL PRIMARY KEY,
    id_user         INT                                 NOT NULL,
    name            VARCHAR(100)                        NOT NULL,
    description     TEXT                                NOT NULL,
    quantidade      INT                                 NOT NULL,
    preco           DECIMAL(10, 2)                     NOT NULL,
    tamanho         VARCHAR(10)                         NOT NULL,
    image_url       TEXT                                NULL,
    is_active       BOOLEAN   DEFAULT TRUE,
    parameter       BOOLEAN   DEFAULT FALSE,
    status_code     INT       DEFAULT 0                 NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_product FOREIGN KEY (id_user) REFERENCES "user" (id)
);
