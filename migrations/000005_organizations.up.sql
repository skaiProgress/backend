-- Organizations (companies) for multi-tenant structure
CREATE TABLE IF NOT EXISTS public.organizations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    bin             TEXT,
    phone           TEXT,
    email           TEXT,
    address         TEXT,
    contact_person  TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_organizations_bin
    ON public.organizations (bin)
    WHERE bin IS NOT NULL AND bin <> '';

DROP TRIGGER IF EXISTS trg_organizations_updated_at ON public.organizations;
CREATE TRIGGER trg_organizations_updated_at
    BEFORE UPDATE ON public.organizations
    FOR EACH ROW
    EXECUTE FUNCTION public.set_updated_at();

-- Link profiles to organizations; add org_admin role
ALTER TABLE public.profiles
    ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES public.organizations (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_profiles_organization_id
    ON public.profiles (organization_id)
    WHERE organization_id IS NOT NULL;

ALTER TABLE public.profiles DROP CONSTRAINT IF EXISTS profiles_role_check;
ALTER TABLE public.profiles
    ADD CONSTRAINT profiles_role_check
        CHECK (role IN ('user', 'admin', 'super_admin', 'org_admin'));
