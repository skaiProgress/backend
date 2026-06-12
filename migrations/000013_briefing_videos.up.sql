-- Fire-safety briefing videos: super admin uploads one video per briefing kind
-- on a course flagged as a briefing course; org-admin links them via calendar.

ALTER TABLE public.courses
    ADD COLUMN IF NOT EXISTS is_briefing_course BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS public.briefing_videos (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id      UUID NOT NULL REFERENCES public.courses(id) ON DELETE CASCADE,
    briefing_kind  TEXT NOT NULL CHECK (briefing_kind IN ('introductory','primary','repeat','unscheduled','targeted')),
    video_url      TEXT NOT NULL,
    video_path     TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, briefing_kind)
);

CREATE INDEX IF NOT EXISTS idx_briefing_videos_course ON public.briefing_videos(course_id);

-- org_events gains a link to the source briefing course and an end-of-window time.
-- starts_at remains the start of the validity window.
ALTER TABLE public.org_events
    ADD COLUMN IF NOT EXISTS course_id UUID REFERENCES public.courses(id) ON DELETE SET NULL;

ALTER TABLE public.org_events
    ADD COLUMN IF NOT EXISTS ends_at TIMESTAMPTZ;
