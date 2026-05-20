-- Employee marks course training complete; triggers AI analytics for org-admin.

ALTER TABLE public.course_assignments
    ADD COLUMN IF NOT EXISTS training_completed_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_course_assignments_training_completed
    ON public.course_assignments (user_id, course_id)
    WHERE training_completed_at IS NOT NULL;
