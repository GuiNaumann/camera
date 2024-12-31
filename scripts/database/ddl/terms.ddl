-- Tabela 'terms'
CREATE TABLE terms
(
    id          SERIAL PRIMARY KEY,
    id_author   INT,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_terms_author FOREIGN KEY (id_author) REFERENCES "user" (id)
);