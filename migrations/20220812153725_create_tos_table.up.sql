CREATE TABLE IF NOT EXISTS TermsOfService (
    id serial PRIMARY KEY,
    content TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp
);
