CREATE DATABASE tickets;
GRANT ALL PRIVILEGES ON DATABASE tickets TO program;

CREATE DATABASE flights;
GRANT ALL PRIVILEGES ON DATABASE flights TO program;

CREATE DATABASE privileges;
GRANT ALL PRIVILEGES ON DATABASE privileges TO program;

CREATE DATABASE notes;
GRANT ALL PRIVILEGES ON DATABASE notes TO program;

\connect notes;

CREATE TABLE scope (
    ID SERIAL PRIMARY KEY,
    scope_type INT NOT NULL
        CHECK (scope_type = 1 or scope_type = 2),
    owner_id INT NOT NULL
);

CREATE TABLE tag (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(20) NOT NULL
);

CREATE TABLE tag_note_connection (
    ID SERIAL PRIMARY KEY,
    tag_ID INT REFERENCES tag(ID)
    note_ID INT REFERENCES note(ID)
);

CREATE TABLE note
(
    ID           SERIAL PRIMARY KEY,
    scope_id    int REFERENCES scope(ID),
    -- связь с тегами через доп таблицу
    author_id   int,
    title           VARCHAR(50),
    content       VARCHAR(2000),
    created     TIMESTAMP
);

-- \connect flights;

-- CREATE TABLE airport
-- (
--     id      SERIAL PRIMARY KEY,
--     name    VARCHAR(255),
--     city    VARCHAR(255),
--     country VARCHAR(255)
-- );

-- CREATE TABLE flight
-- (
--     id              SERIAL PRIMARY KEY,
--     flight_number   VARCHAR(20)              NOT NULL,
--     datetime        TIMESTAMP WITH TIME ZONE NOT NULL,
--     from_airport_id INT REFERENCES airport (id),
--     to_airport_id   INT REFERENCES airport (id),
--     price           INT                      NOT NULL
-- );

-- INSERT INTO airport VALUES (1, 'Шереметьево', 'Москва', 'Россия');
-- INSERT INTO airport VALUES (2, 'Пулково', 'Санкт-Петербург', 'Россия');
-- INSERT INTO flight VALUES (1, 'AFL031', '2021-10-08 20:00', 2, 1, 1500);

-- \connect privileges;

-- CREATE TABLE privilege
-- (
--     id       SERIAL PRIMARY KEY,
--     username VARCHAR(80) NOT NULL UNIQUE,
--     status   VARCHAR(80) NOT NULL DEFAULT 'BRONZE'
--         CHECK (status IN ('BRONZE', 'SILVER', 'GOLD')),
--     balance  INT
-- );

-- CREATE TABLE privilege_history
-- (
--     id             SERIAL PRIMARY KEY,
--     privilege_id   INT REFERENCES privilege (id),
--     ticket_uid     uuid        NOT NULL,
--     datetime       TIMESTAMP   NOT NULL,
--     balance_diff   INT         NOT NULL,
--     operation_type VARCHAR(20) NOT NULL
--         CHECK (operation_type IN ('FILL_IN_BALANCE', 'DEBIT_THE_ACCOUNT'))
-- );

