DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS tokens;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT,
    password TEXT
);

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    token TEXT
);

INSERT INTO
    users (email, password)
VALUES
    (
        'test@test.ru',
        '$2a$12$oHZyRyQxlrWVTjNACxFzw.udNF00.Mk4KTLjfyCZyojgzkWidKrNW'
    );