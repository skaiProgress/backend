-- org_events: calendar events tied to an organization
CREATE TABLE IF NOT EXISTS public.org_events (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES public.organizations(id) ON DELETE CASCADE,
    employee_id   UUID REFERENCES public.profiles(id) ON DELETE SET NULL,
    title         TEXT NOT NULL,
    event_type    TEXT NOT NULL CHECK (event_type IN ('training','drill','inspection','meeting')),
    briefing_kind TEXT CHECK (briefing_kind IN ('introductory','primary','repeat','unscheduled','targeted')),
    starts_at     TIMESTAMPTZ NOT NULL,
    location      TEXT NOT NULL DEFAULT '',
    participants  INT,
    created_by    UUID REFERENCES public.profiles(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_org_events_org ON public.org_events(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_events_employee ON public.org_events(employee_id);

-- briefing_records: journal entries for fire-safety briefings
CREATE TABLE IF NOT EXISTS public.briefing_records (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES public.organizations(id) ON DELETE CASCADE,
    event_id         UUID REFERENCES public.org_events(id) ON DELETE SET NULL,
    employee_id      UUID NOT NULL REFERENCES public.profiles(id) ON DELETE CASCADE,
    employee_name    TEXT NOT NULL,
    position         TEXT NOT NULL DEFAULT '',
    briefing_kind    TEXT NOT NULL CHECK (briefing_kind IN ('introductory','primary','repeat','unscheduled','targeted')),
    instructor_name  TEXT NOT NULL DEFAULT '',
    instructor_id    UUID REFERENCES public.profiles(id) ON DELETE SET NULL,
    date_conducted   DATE NOT NULL,
    employee_signed  BOOLEAN NOT NULL DEFAULT FALSE,
    employee_signed_at TIMESTAMPTZ,
    instructor_signed BOOLEAN NOT NULL DEFAULT FALSE,
    instructor_signed_at TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_briefing_records_org ON public.briefing_records(organization_id);
CREATE INDEX IF NOT EXISTS idx_briefing_records_employee ON public.briefing_records(employee_id);
CREATE INDEX IF NOT EXISTS idx_briefing_records_event ON public.briefing_records(event_id);
