CREATE TYPE USER_ROLE AS ENUM ('admin', 'candidate', 'user');
CREATE TYPE GENDER AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS Users (
    id serial PRIMARY KEY,
    email varchar(320) NOT NULL,
    phone varchar(20) NOT NULL,
    password varchar(255) NOT NULL,
    role USER_ROLE NOT NULL DEFAULT 'user',
    gender GENDER NOT NULL DEFAULT 'male',
    date_of_birth date NOT NULL,
    description text DEFAULT '',
    facebook_id varchar(255),
    google_id varchar(255),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp
);
