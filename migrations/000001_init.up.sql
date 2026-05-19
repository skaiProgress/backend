-- Initial backend schema placeholder for AIQADAM backend-wrapper.
CREATE TABLE IF NOT EXISTS backend_meta (
    id         INT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO backend_meta (id) VALUES (1)
ON CONFLICT (id) DO NOTHING;
