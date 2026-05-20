-- AI Analysis module: stores Gemini analysis results per employee per quiz attempt.

CREATE TABLE IF NOT EXISTS public.ai_analysis (
    id              SERIAL PRIMARY KEY,
    employee_id     UUID REFERENCES public.profiles (id) ON DELETE CASCADE,
    organization_id UUID REFERENCES public.organizations (id) ON DELETE CASCADE,
    quiz_result_id  UUID,
    course_name     VARCHAR(255),
    score           FLOAT,
    weak_topics     TEXT[],
    recommendation  TEXT,
    risk_level      VARCHAR(10),
    summary         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_analysis_org
    ON public.ai_analysis (organization_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_ai_analysis_employee
    ON public.ai_analysis (employee_id, created_at DESC);

-- Knowledge base: regulatory norms for ПТМ fire safety prompts.
CREATE TABLE IF NOT EXISTS public.knowledge_base (
    id         SERIAL PRIMARY KEY,
    topic      VARCHAR(100) NOT NULL,
    content    TEXT NOT NULL,
    article    VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO public.knowledge_base (topic, content, article) VALUES
('эвакуация',        'Руководитель обязан обеспечить наличие планов эвакуации на каждом этаже. Пути эвакуации должны быть свободны.', 'п.15'),
('огнетушители',     'Огнетушители проверяются не реже 1 раза в год. На каждые 100 кв.м. не менее 1 огнетушителя.', 'п.23'),
('инструктаж',       'Первичный инструктаж — до начала работы. Повторный — не реже 1 раза в полгода.', 'п.8'),
('сигнализация',     'Системы пожарной сигнализации проходят техобслуживание не реже 1 раза в квартал.', 'п.31'),
('действия при пожаре', 'При пожаре: позвонить 101, эвакуировать людей, тушить первичными средствами.', 'п.12');
