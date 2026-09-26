CREATE TABLE documents(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    filename TEXT NOT NULL,
    storage_path TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK(status IN('pending','processing','ready','failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);