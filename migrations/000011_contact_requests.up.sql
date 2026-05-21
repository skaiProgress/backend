CREATE TABLE IF NOT EXISTS public.contact_requests (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    email       TEXT NOT NULL,
    phone       TEXT NOT NULL,
    company     TEXT,
    message     TEXT,
    status      TEXT NOT NULL DEFAULT 'new'
                    CHECK (status IN ('new', 'in_progress', 'done')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contact_requests_status
    ON public.contact_requests (status);

CREATE INDEX IF NOT EXISTS idx_contact_requests_created_at
    ON public.contact_requests (created_at DESC);

DROP TRIGGER IF EXISTS trg_contact_requests_updated_at ON public.contact_requests;
CREATE TRIGGER trg_contact_requests_updated_at
    BEFORE UPDATE ON public.contact_requests
    FOR EACH ROW
    EXECUTE FUNCTION public.set_updated_at();
