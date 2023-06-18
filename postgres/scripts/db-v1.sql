CREATE DATABASE statistics;
GRANT ALL PRIVILEGES ON DATABASE statistics TO program;

\connect statistics;
GRANT USAGE, CREATE ON SCHEMA public TO program;

CREATE DATABASE notes;
GRANT ALL PRIVILEGES ON DATABASE notes TO program;

\connect notes;

CREATE TABLE scope (
    ID SERIAL PRIMARY KEY,
    scope_type INT NOT NULL
        CHECK (scope_type = 1 or scope_type = 2) default 1,
    owner_id INT NOT NULL
);

CREATE TABLE tag (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(20) NOT NULL
);

CREATE TABLE note (
    ID           SERIAL PRIMARY KEY,
    scope_id    int REFERENCES scope(ID),
    -- связь с тегами через доп таблицу
    author_id   int,
    title           VARCHAR(50),
    content       VARCHAR(2000),
    CreatedAt     TIMESTAMP default current_timestamp,
    UpdatedAt   TIMESTAMP default current_timestamp
);

CREATE TABLE tag_note_connection (
    ID SERIAL PRIMARY KEY,
    tag_ID INT REFERENCES tag(ID),
    note_ID INT REFERENCES note(ID)
);


