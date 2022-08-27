DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS tokens;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT,
    password TEXT,
    is_activated BOOLEAN
);

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    token TEXT
);

INSERT INTO
    users (email, password, is_activated)
VALUES
    (
        'test@test.ru',
        '$2a$12$oHZyRyQxlrWVTjNACxFzw.udNF00.Mk4KTLjfyCZyojgzkWidKrNW',
        FALSE
    );