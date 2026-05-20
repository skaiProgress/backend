ALTER TABLE public.profiles DROP CONSTRAINT IF EXISTS profiles_role_check;
ALTER TABLE public.profiles
    ADD CONSTRAINT profiles_role_check
        CHECK (role IN ('user', 'admin', 'super_admin'));

ALTER TABLE public.profiles DROP COLUMN IF EXISTS organization_id;

DROP TRIGGER IF EXISTS trg_organizations_updated_at ON public.organizations;
DROP INDEX IF EXISTS idx_organizations_bin;
DROP TABLE IF EXISTS public.organizations;
