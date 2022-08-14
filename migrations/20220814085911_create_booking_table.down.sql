CREATE TYPE BOOKING_STATUS AS ENUM ('PENDING', 'CONFIRMED', 'CANCELLED', 'COMPLETED');

CREATE TABLE IF NOT EXISTS Bookings (
    id serial PRIMARY KEY,
    booking_status BOOKING_STATUS NOT NULL DEFAULT 'PENDING',
    price int NOT NULL,
    user_id serial NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted_at timestamp
);
