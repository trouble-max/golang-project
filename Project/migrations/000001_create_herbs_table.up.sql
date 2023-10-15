CREATE TABLE IF NOT EXISTS herbs (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    description text NOT NULL,
    price numeric(5, 2) NOT NULL,
    culinary_uses text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);