ALTER TABLE public.contact_requests
    ADD COLUMN IF NOT EXISTS phone TEXT;

UPDATE public.contact_requests
SET phone = ''
WHERE phone IS NULL;

ALTER TABLE public.contact_requests
    ALTER COLUMN phone SET NOT NULL;

ALTER TABLE public.contact_requests
    ALTER COLUMN message DROP NOT NULL;
