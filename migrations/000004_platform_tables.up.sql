-- Courses, lessons, materials, assignments (migrated from Supabase schema).

CREATE TABLE IF NOT EXISTS public.courses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published')),
    cover_url   TEXT,
    created_by  UUID REFERENCES auth.users (id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DROP TRIGGER IF EXISTS trg_courses_updated_at ON public.courses;
CREATE TRIGGER trg_courses_updated_at
    BEFORE UPDATE ON public.courses
    FOR EACH ROW
    EXECUTE FUNCTION public.set_updated_at();

CREATE TABLE IF NOT EXISTS public.lessons (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id        UUID NOT NULL REFERENCES public.courses (id) ON DELETE CASCADE,
    title            TEXT NOT NULL,
    description      TEXT,
    youtube_url      TEXT NOT NULL,
    youtube_video_id TEXT NOT NULL,
    order_index      INT NOT NULL DEFAULT 1,
    is_free          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lessons_course_order ON public.lessons (course_id, order_index);

DROP TRIGGER IF EXISTS trg_lessons_updated_at ON public.lessons;
CREATE TRIGGER trg_lessons_updated_at
    BEFORE UPDATE ON public.lessons
    FOR EACH ROW
    EXECUTE FUNCTION public.set_updated_at();

CREATE TABLE IF NOT EXISTS public.course_materials (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id  UUID NOT NULL REFERENCES public.courses (id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    file_url   TEXT NOT NULL,
    file_type  TEXT,
    file_size  BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_materials_course_id ON public.course_materials (course_id);

CREATE TABLE IF NOT EXISTS public.course_assignments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES public.profiles (id) ON DELETE CASCADE,
    course_id   UUID NOT NULL REFERENCES public.courses (id) ON DELETE CASCADE,
    assigned_by UUID REFERENCES public.profiles (id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ,
    status      TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'revoked')),
    revoked_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_assignments_user ON public.course_assignments (user_id);
CREATE INDEX IF NOT EXISTS idx_assignments_course ON public.course_assignments (course_id);

DROP TRIGGER IF EXISTS trg_assignments_updated_at ON public.course_assignments;
CREATE TRIGGER trg_assignments_updated_at
    BEFORE UPDATE ON public.course_assignments
    FOR EACH ROW
    EXECUTE FUNCTION public.set_updated_at();
