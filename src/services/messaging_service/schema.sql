CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    sender TEXT NOT NULL,
    receiver TEXT NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);