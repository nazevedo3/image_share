CREATE TABLE IF NOT EXISTS images (
    id bigserial PRIMARY KEY,  
    title text NOT NULL,
    description text NOT NULL,
    tags text[] NOT NULL,
    path text NOT NULL
);