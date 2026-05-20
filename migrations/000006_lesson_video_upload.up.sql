-- Support direct video uploads for lessons (in addition to YouTube).

ALTER TABLE public.lessons
    ADD COLUMN IF NOT EXISTS video_source TEXT NOT NULL DEFAULT 'youtube',
    ADD COLUMN IF NOT EXISTS video_url TEXT;

ALTER TABLE public.lessons
    ALTER COLUMN youtube_url DROP NOT NULL,
    ALTER COLUMN youtube_video_id DROP NOT NULL;

UPDATE public.lessons
SET video_source = 'youtube'
WHERE video_source IS NULL OR video_source = '';

ALTER TABLE public.lessons
    DROP CONSTRAINT IF EXISTS lessons_video_source_check;

ALTER TABLE public.lessons
    ADD CONSTRAINT lessons_video_source_check
        CHECK (video_source IN ('youtube', 'upload'));

ALTER TABLE public.lessons
    DROP CONSTRAINT IF EXISTS lessons_video_payload_check;

ALTER TABLE public.lessons
    ADD CONSTRAINT lessons_video_payload_check CHECK (
        (
            video_source = 'youtube'
            AND youtube_url IS NOT NULL
            AND youtube_video_id IS NOT NULL
            AND btrim(youtube_url) <> ''
            AND btrim(youtube_video_id) <> ''
        )
        OR (
            video_source = 'upload'
            AND video_url IS NOT NULL
            AND btrim(video_url) <> ''
        )
    );
