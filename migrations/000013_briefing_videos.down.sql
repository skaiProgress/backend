ALTER TABLE public.org_events DROP COLUMN IF EXISTS ends_at;
ALTER TABLE public.org_events DROP COLUMN IF EXISTS course_id;

DROP TABLE IF EXISTS public.briefing_videos;

ALTER TABLE public.courses DROP COLUMN IF EXISTS is_briefing_course;
