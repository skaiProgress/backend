ALTER TABLE public.contact_requests ALTER COLUMN message SET NOT NULL;
ALTER TABLE public.contact_requests DROP COLUMN IF EXISTS phone;
