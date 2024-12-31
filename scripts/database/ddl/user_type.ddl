-- Tabela 'user_type'
CREATE TABLE user_type
(
    id_user      INT                                 NOT NULL,
    id_user_type INT                                 NOT NULL,
    status_code  SMALLINT  DEFAULT 0                 NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_type_user FOREIGN KEY (id_user) REFERENCES "user" (id),
    CONSTRAINT fk_user_type_user_type FOREIGN KEY (id_user_type) REFERENCES type_user (id)
);

CREATE INDEX idx_user_type_user ON user_type (id_user);
CREATE INDEX idx_user_type_user_type ON user_type (id_user_type);