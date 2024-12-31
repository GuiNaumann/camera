-- Tabela 'states' (Estados do Brasil)
CREATE TABLE states
(
    id         SERIAL PRIMARY KEY,
    state_name VARCHAR(50) NOT NULL,
    state_code VARCHAR(2)  NOT NULL
);

-- Inserindo os estados do Brasil
INSERT INTO states (state_name, state_code)
VALUES ('Acre', 'AC');
INSERT INTO states (state_name, state_code)
VALUES ('Alagoas', 'AL');
INSERT INTO states (state_name, state_code)
VALUES ('Amapá', 'AP');
INSERT INTO states (state_name, state_code)
VALUES ('Amazonas', 'AM');
INSERT INTO states (state_name, state_code)
VALUES ('Bahia', 'BA');
INSERT INTO states (state_name, state_code)
VALUES ('Ceará', 'CE');
INSERT INTO states (state_name, state_code)
VALUES ('Distrito Federal', 'DF');
INSERT INTO states (state_name, state_code)
VALUES ('Espírito Santo', 'ES');
INSERT INTO states (state_name, state_code)
VALUES ('Goiás', 'GO');
INSERT INTO states (state_name, state_code)
VALUES ('Maranhão', 'MA');
INSERT INTO states (state_name, state_code)
VALUES ('Mato Grosso', 'MT');
INSERT INTO states (state_name, state_code)
VALUES ('Mato Grosso do Sul', 'MS');
INSERT INTO states (state_name, state_code)
VALUES ('Minas Gerais', 'MG');
INSERT INTO states (state_name, state_code)
VALUES ('Pará', 'PA');
INSERT INTO states (state_name, state_code)
VALUES ('Paraíba', 'PB');
INSERT INTO states (state_name, state_code)
VALUES ('Paraná', 'PR');
INSERT INTO states (state_name, state_code)
VALUES ('Pernambuco', 'PE');
INSERT INTO states (state_name, state_code)
VALUES ('Piauí', 'PI');
INSERT INTO states (state_name, state_code)
VALUES ('Rio de Janeiro', 'RJ');
INSERT INTO states (state_name, state_code)
VALUES ('Rio Grande do Norte', 'RN');
INSERT INTO states (state_name, state_code)
VALUES ('Rio Grande do Sul', 'RS');
INSERT INTO states (state_name, state_code)
VALUES ('Rondônia', 'RO');
INSERT INTO states (state_name, state_code)
VALUES ('Roraima', 'RR');
INSERT INTO states (state_name, state_code)
VALUES ('Santa Catarina', 'SC');
INSERT INTO states (state_name, state_code)
VALUES ('São Paulo', 'SP');
INSERT INTO states (state_name, state_code)
VALUES ('Sergipe', 'SE');
INSERT INTO states (state_name, state_code)
VALUES ('Tocantins', 'TO');
