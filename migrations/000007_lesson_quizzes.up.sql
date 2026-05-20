-- Mini-quizzes for lessons (5 questions, uploaded as .txt by admin).

CREATE TABLE IF NOT EXISTS public.lesson_quizzes (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id        UUID NOT NULL UNIQUE REFERENCES public.lessons (id) ON DELETE CASCADE,
    source_file_name TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DROP TRIGGER IF EXISTS trg_lesson_quizzes_updated_at ON public.lesson_quizzes;
CREATE TRIGGER trg_lesson_quizzes_updated_at
    BEFORE UPDATE ON public.lesson_quizzes
    FOR EACH ROW
    EXECUTE FUNCTION public.set_updated_at();

CREATE TABLE IF NOT EXISTS public.lesson_quiz_questions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id       UUID NOT NULL REFERENCES public.lesson_quizzes (id) ON DELETE CASCADE,
    order_index   INT NOT NULL CHECK (order_index BETWEEN 1 AND 5),
    question_text TEXT NOT NULL,
    option_a      TEXT NOT NULL,
    option_b      TEXT NOT NULL,
    option_c      TEXT NOT NULL,
    correct_option CHAR(1) NOT NULL CHECK (correct_option IN ('A', 'B', 'C')),
    UNIQUE (quiz_id, order_index)
);

CREATE INDEX IF NOT EXISTS idx_lesson_quiz_questions_quiz
    ON public.lesson_quiz_questions (quiz_id, order_index);

CREATE TABLE IF NOT EXISTS public.lesson_quiz_attempts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES public.profiles (id) ON DELETE CASCADE,
    lesson_id    UUID NOT NULL REFERENCES public.lessons (id) ON DELETE CASCADE,
    quiz_id      UUID NOT NULL REFERENCES public.lesson_quizzes (id) ON DELETE CASCADE,
    answers      JSONB NOT NULL,
    score        INT NOT NULL CHECK (score >= 0),
    max_score    INT NOT NULL DEFAULT 5 CHECK (max_score = 5),
    passed       BOOLEAN NOT NULL,
    completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lesson_quiz_attempts_user_lesson
    ON public.lesson_quiz_attempts (user_id, lesson_id, completed_at DESC);
