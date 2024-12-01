CREATE TABLE
    IF NOT EXISTS costumers (
        id SERIAL PRIMARY KEY,
        name VARCHAR(128) NOT NULL,
        cpf VARCHAR(11) NOT NULL,
        birthdate DATE NOT NULL,
        phone VARCHAR(15) NOT NULL,
        email VARCHAR(256) NOT NULL,
        phone_validated BOOLEAN NOT NULL,
        email_validated BOOLEAN NOT NULL,
        deleted BOOLEAN NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ
    );

CREATE TABLE
    IF NOT EXISTS outbox (
        id BIGSERIAL PRIMARY KEY,
        event_name TEXT NOT NULL,
        message JSONB NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        processed_at TIMESTAMPTZ DEFAULT NULL
    );