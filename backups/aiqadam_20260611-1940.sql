--
-- PostgreSQL database dump
--

\restrict 7nhzqGsCuAwzQSQL5qp5PgvTSNX8jURsrCfw7NaqOgaOAW6jlp1odRIeozgxRuY

-- Dumped from database version 16.14
-- Dumped by pg_dump version 16.14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.profiles DROP CONSTRAINT IF EXISTS profiles_organization_id_fkey;
ALTER TABLE IF EXISTS ONLY public.profiles DROP CONSTRAINT IF EXISTS profiles_id_fkey;
ALTER TABLE IF EXISTS ONLY public.org_events DROP CONSTRAINT IF EXISTS org_events_organization_id_fkey;
ALTER TABLE IF EXISTS ONLY public.org_events DROP CONSTRAINT IF EXISTS org_events_employee_id_fkey;
ALTER TABLE IF EXISTS ONLY public.org_events DROP CONSTRAINT IF EXISTS org_events_created_by_fkey;
ALTER TABLE IF EXISTS ONLY public.org_events DROP CONSTRAINT IF EXISTS org_events_course_id_fkey;
ALTER TABLE IF EXISTS ONLY public.lessons DROP CONSTRAINT IF EXISTS lessons_course_id_fkey;
ALTER TABLE IF EXISTS ONLY public.lesson_quizzes DROP CONSTRAINT IF EXISTS lesson_quizzes_lesson_id_fkey;
ALTER TABLE IF EXISTS ONLY public.lesson_quiz_questions DROP CONSTRAINT IF EXISTS lesson_quiz_questions_quiz_id_fkey;
ALTER TABLE IF EXISTS ONLY public.lesson_quiz_attempts DROP CONSTRAINT IF EXISTS lesson_quiz_attempts_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.lesson_quiz_attempts DROP CONSTRAINT IF EXISTS lesson_quiz_attempts_quiz_id_fkey;
ALTER TABLE IF EXISTS ONLY public.lesson_quiz_attempts DROP CONSTRAINT IF EXISTS lesson_quiz_attempts_lesson_id_fkey;
ALTER TABLE IF EXISTS ONLY public.courses DROP CONSTRAINT IF EXISTS courses_created_by_fkey;
ALTER TABLE IF EXISTS ONLY public.course_materials DROP CONSTRAINT IF EXISTS course_materials_course_id_fkey;
ALTER TABLE IF EXISTS ONLY public.course_assignments DROP CONSTRAINT IF EXISTS course_assignments_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.course_assignments DROP CONSTRAINT IF EXISTS course_assignments_course_id_fkey;
ALTER TABLE IF EXISTS ONLY public.course_assignments DROP CONSTRAINT IF EXISTS course_assignments_assigned_by_fkey;
ALTER TABLE IF EXISTS ONLY public.briefing_videos DROP CONSTRAINT IF EXISTS briefing_videos_course_id_fkey;
ALTER TABLE IF EXISTS ONLY public.briefing_records DROP CONSTRAINT IF EXISTS briefing_records_organization_id_fkey;
ALTER TABLE IF EXISTS ONLY public.briefing_records DROP CONSTRAINT IF EXISTS briefing_records_instructor_id_fkey;
ALTER TABLE IF EXISTS ONLY public.briefing_records DROP CONSTRAINT IF EXISTS briefing_records_event_id_fkey;
ALTER TABLE IF EXISTS ONLY public.briefing_records DROP CONSTRAINT IF EXISTS briefing_records_employee_id_fkey;
ALTER TABLE IF EXISTS ONLY public.ai_analysis DROP CONSTRAINT IF EXISTS ai_analysis_organization_id_fkey;
ALTER TABLE IF EXISTS ONLY public.ai_analysis DROP CONSTRAINT IF EXISTS ai_analysis_employee_id_fkey;
DROP TRIGGER IF EXISTS trg_profiles_updated_at ON public.profiles;
DROP TRIGGER IF EXISTS trg_organizations_updated_at ON public.organizations;
DROP TRIGGER IF EXISTS trg_lessons_updated_at ON public.lessons;
DROP TRIGGER IF EXISTS trg_lesson_quizzes_updated_at ON public.lesson_quizzes;
DROP TRIGGER IF EXISTS trg_courses_updated_at ON public.courses;
DROP TRIGGER IF EXISTS trg_contact_requests_updated_at ON public.contact_requests;
DROP TRIGGER IF EXISTS trg_assignments_updated_at ON public.course_assignments;
DROP INDEX IF EXISTS public.idx_profiles_organization_id;
DROP INDEX IF EXISTS public.idx_organizations_bin;
DROP INDEX IF EXISTS public.idx_org_events_org;
DROP INDEX IF EXISTS public.idx_org_events_employee;
DROP INDEX IF EXISTS public.idx_lessons_course_order;
DROP INDEX IF EXISTS public.idx_lesson_quiz_questions_quiz;
DROP INDEX IF EXISTS public.idx_lesson_quiz_attempts_user_lesson;
DROP INDEX IF EXISTS public.idx_course_materials_course_id;
DROP INDEX IF EXISTS public.idx_course_assignments_training_completed;
DROP INDEX IF EXISTS public.idx_contact_requests_status;
DROP INDEX IF EXISTS public.idx_contact_requests_created_at;
DROP INDEX IF EXISTS public.idx_briefing_videos_course;
DROP INDEX IF EXISTS public.idx_briefing_records_org;
DROP INDEX IF EXISTS public.idx_briefing_records_event;
DROP INDEX IF EXISTS public.idx_briefing_records_employee;
DROP INDEX IF EXISTS public.idx_assignments_user;
DROP INDEX IF EXISTS public.idx_assignments_course;
DROP INDEX IF EXISTS public.idx_ai_analysis_org;
DROP INDEX IF EXISTS public.idx_ai_analysis_employee;
DROP INDEX IF EXISTS auth.auth_users_email_idx;
ALTER TABLE IF EXISTS ONLY public.schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.profiles DROP CONSTRAINT IF EXISTS profiles_pkey;
ALTER TABLE IF EXISTS ONLY public.organizations DROP CONSTRAINT IF EXISTS organizations_pkey;
ALTER TABLE IF EXISTS ONLY public.org_events DROP CONSTRAINT IF EXISTS org_events_pkey;
ALTER TABLE IF EXISTS ONLY public.lessons DROP CONSTRAINT IF EXISTS lessons_pkey;
ALTER TABLE IF EXISTS ONLY public.lesson_quizzes DROP CONSTRAINT IF EXISTS lesson_quizzes_pkey;
ALTER TABLE IF EXISTS ONLY public.lesson_quizzes DROP CONSTRAINT IF EXISTS lesson_quizzes_lesson_id_key;
ALTER TABLE IF EXISTS ONLY public.lesson_quiz_questions DROP CONSTRAINT IF EXISTS lesson_quiz_questions_quiz_id_order_index_key;
ALTER TABLE IF EXISTS ONLY public.lesson_quiz_questions DROP CONSTRAINT IF EXISTS lesson_quiz_questions_pkey;
ALTER TABLE IF EXISTS ONLY public.lesson_quiz_attempts DROP CONSTRAINT IF EXISTS lesson_quiz_attempts_pkey;
ALTER TABLE IF EXISTS ONLY public.knowledge_base DROP CONSTRAINT IF EXISTS knowledge_base_pkey;
ALTER TABLE IF EXISTS ONLY public.courses DROP CONSTRAINT IF EXISTS courses_pkey;
ALTER TABLE IF EXISTS ONLY public.course_materials DROP CONSTRAINT IF EXISTS course_materials_pkey;
ALTER TABLE IF EXISTS ONLY public.course_assignments DROP CONSTRAINT IF EXISTS course_assignments_user_id_course_id_key;
ALTER TABLE IF EXISTS ONLY public.course_assignments DROP CONSTRAINT IF EXISTS course_assignments_pkey;
ALTER TABLE IF EXISTS ONLY public.contact_requests DROP CONSTRAINT IF EXISTS contact_requests_pkey;
ALTER TABLE IF EXISTS ONLY public.briefing_videos DROP CONSTRAINT IF EXISTS briefing_videos_pkey;
ALTER TABLE IF EXISTS ONLY public.briefing_videos DROP CONSTRAINT IF EXISTS briefing_videos_course_id_briefing_kind_key;
ALTER TABLE IF EXISTS ONLY public.briefing_records DROP CONSTRAINT IF EXISTS briefing_records_pkey;
ALTER TABLE IF EXISTS ONLY public.backend_meta DROP CONSTRAINT IF EXISTS backend_meta_pkey;
ALTER TABLE IF EXISTS ONLY public.ai_analysis DROP CONSTRAINT IF EXISTS ai_analysis_pkey;
ALTER TABLE IF EXISTS ONLY auth.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY auth.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS public.knowledge_base ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.ai_analysis ALTER COLUMN id DROP DEFAULT;
DROP TABLE IF EXISTS public.schema_migrations;
DROP TABLE IF EXISTS public.profiles;
DROP TABLE IF EXISTS public.organizations;
DROP TABLE IF EXISTS public.org_events;
DROP TABLE IF EXISTS public.lessons;
DROP TABLE IF EXISTS public.lesson_quizzes;
DROP TABLE IF EXISTS public.lesson_quiz_questions;
DROP TABLE IF EXISTS public.lesson_quiz_attempts;
DROP SEQUENCE IF EXISTS public.knowledge_base_id_seq;
DROP TABLE IF EXISTS public.knowledge_base;
DROP TABLE IF EXISTS public.courses;
DROP TABLE IF EXISTS public.course_materials;
DROP TABLE IF EXISTS public.course_assignments;
DROP TABLE IF EXISTS public.contact_requests;
DROP TABLE IF EXISTS public.briefing_videos;
DROP TABLE IF EXISTS public.briefing_records;
DROP TABLE IF EXISTS public.backend_meta;
DROP SEQUENCE IF EXISTS public.ai_analysis_id_seq;
DROP TABLE IF EXISTS public.ai_analysis;
DROP TABLE IF EXISTS auth.users;
DROP FUNCTION IF EXISTS public.set_updated_at();
DROP EXTENSION IF EXISTS pgcrypto;
DROP SCHEMA IF EXISTS auth;
--
-- Name: auth; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA auth;


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: set_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: auth; Owner: -
--

CREATE TABLE auth.users (
    instance_id uuid,
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    aud text,
    role text,
    email text,
    encrypted_password text,
    email_confirmed_at timestamp with time zone,
    invited_at timestamp with time zone,
    confirmation_token text,
    confirmation_sent_at timestamp with time zone,
    recovery_token text,
    recovery_sent_at timestamp with time zone,
    email_change_token_new text,
    email_change text,
    email_change_sent_at timestamp with time zone,
    last_sign_in_at timestamp with time zone,
    raw_app_meta_data jsonb,
    raw_user_meta_data jsonb,
    is_super_admin boolean,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    phone text,
    phone_confirmed_at timestamp with time zone,
    phone_change text,
    phone_change_token text,
    phone_change_sent_at timestamp with time zone,
    email_change_token_current text,
    email_change_confirm_status smallint,
    banned_until timestamp with time zone,
    reauthentication_token text,
    reauthentication_sent_at timestamp with time zone,
    is_sso_user boolean DEFAULT false,
    deleted_at timestamp with time zone,
    is_anonymous boolean DEFAULT false
);


--
-- Name: ai_analysis; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ai_analysis (
    id integer NOT NULL,
    employee_id uuid,
    organization_id uuid,
    quiz_result_id uuid,
    course_name character varying(255),
    score double precision,
    weak_topics text[],
    recommendation text,
    risk_level character varying(10),
    summary text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: ai_analysis_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ai_analysis_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ai_analysis_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ai_analysis_id_seq OWNED BY public.ai_analysis.id;


--
-- Name: backend_meta; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.backend_meta (
    id integer DEFAULT 1 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT backend_meta_id_check CHECK ((id = 1))
);


--
-- Name: briefing_records; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.briefing_records (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    organization_id uuid NOT NULL,
    event_id uuid,
    employee_id uuid NOT NULL,
    employee_name text NOT NULL,
    "position" text DEFAULT ''::text NOT NULL,
    briefing_kind text NOT NULL,
    instructor_name text DEFAULT ''::text NOT NULL,
    instructor_id uuid,
    date_conducted date NOT NULL,
    employee_signed boolean DEFAULT false NOT NULL,
    employee_signed_at timestamp with time zone,
    instructor_signed boolean DEFAULT false NOT NULL,
    instructor_signed_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT briefing_records_briefing_kind_check CHECK ((briefing_kind = ANY (ARRAY['introductory'::text, 'primary'::text, 'repeat'::text, 'unscheduled'::text, 'targeted'::text])))
);


--
-- Name: briefing_videos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.briefing_videos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    course_id uuid NOT NULL,
    briefing_kind text NOT NULL,
    video_url text NOT NULL,
    video_path text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT briefing_videos_briefing_kind_check CHECK ((briefing_kind = ANY (ARRAY['introductory'::text, 'primary'::text, 'repeat'::text, 'unscheduled'::text, 'targeted'::text])))
);


--
-- Name: contact_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.contact_requests (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    company text,
    message text,
    status text DEFAULT 'new'::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    phone text NOT NULL,
    CONSTRAINT contact_requests_status_check CHECK ((status = ANY (ARRAY['new'::text, 'in_progress'::text, 'done'::text])))
);


--
-- Name: course_assignments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.course_assignments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    course_id uuid NOT NULL,
    assigned_by uuid,
    assigned_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone,
    status text DEFAULT 'active'::text NOT NULL,
    revoked_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    training_completed_at timestamp with time zone,
    CONSTRAINT course_assignments_status_check CHECK ((status = ANY (ARRAY['active'::text, 'revoked'::text])))
);


--
-- Name: course_materials; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.course_materials (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    course_id uuid NOT NULL,
    name text NOT NULL,
    file_url text NOT NULL,
    file_type text,
    file_size bigint,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: courses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.courses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    description text DEFAULT ''::text NOT NULL,
    status text DEFAULT 'draft'::text NOT NULL,
    cover_url text,
    created_by uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    is_briefing_course boolean DEFAULT false NOT NULL,
    CONSTRAINT courses_status_check CHECK ((status = ANY (ARRAY['draft'::text, 'published'::text])))
);


--
-- Name: knowledge_base; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.knowledge_base (
    id integer NOT NULL,
    topic character varying(100) NOT NULL,
    content text NOT NULL,
    article character varying(50),
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: knowledge_base_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.knowledge_base_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: knowledge_base_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.knowledge_base_id_seq OWNED BY public.knowledge_base.id;


--
-- Name: lesson_quiz_attempts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.lesson_quiz_attempts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    lesson_id uuid NOT NULL,
    quiz_id uuid NOT NULL,
    answers jsonb NOT NULL,
    score integer NOT NULL,
    max_score integer DEFAULT 5 NOT NULL,
    passed boolean NOT NULL,
    completed_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT lesson_quiz_attempts_max_score_check CHECK ((max_score = 5)),
    CONSTRAINT lesson_quiz_attempts_score_check CHECK ((score >= 0))
);


--
-- Name: lesson_quiz_questions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.lesson_quiz_questions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    quiz_id uuid NOT NULL,
    order_index integer NOT NULL,
    question_text text NOT NULL,
    option_a text NOT NULL,
    option_b text NOT NULL,
    option_c text NOT NULL,
    correct_option character(1) NOT NULL,
    CONSTRAINT lesson_quiz_questions_correct_option_check CHECK ((correct_option = ANY (ARRAY['A'::bpchar, 'B'::bpchar, 'C'::bpchar]))),
    CONSTRAINT lesson_quiz_questions_order_index_check CHECK (((order_index >= 1) AND (order_index <= 5)))
);


--
-- Name: lesson_quizzes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.lesson_quizzes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    lesson_id uuid NOT NULL,
    source_file_name text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: lessons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.lessons (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    course_id uuid NOT NULL,
    title text NOT NULL,
    description text,
    youtube_url text,
    youtube_video_id text,
    order_index integer DEFAULT 1 NOT NULL,
    is_free boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    video_source text DEFAULT 'youtube'::text NOT NULL,
    video_url text,
    CONSTRAINT lessons_video_payload_check CHECK ((((video_source = 'youtube'::text) AND (youtube_url IS NOT NULL) AND (youtube_video_id IS NOT NULL) AND (btrim(youtube_url) <> ''::text) AND (btrim(youtube_video_id) <> ''::text)) OR ((video_source = 'upload'::text) AND (video_url IS NOT NULL) AND (btrim(video_url) <> ''::text)))),
    CONSTRAINT lessons_video_source_check CHECK ((video_source = ANY (ARRAY['youtube'::text, 'upload'::text])))
);


--
-- Name: org_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.org_events (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    organization_id uuid NOT NULL,
    employee_id uuid,
    title text NOT NULL,
    event_type text NOT NULL,
    briefing_kind text,
    starts_at timestamp with time zone NOT NULL,
    location text DEFAULT ''::text NOT NULL,
    participants integer,
    created_by uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    course_id uuid,
    ends_at timestamp with time zone,
    CONSTRAINT org_events_briefing_kind_check CHECK ((briefing_kind = ANY (ARRAY['introductory'::text, 'primary'::text, 'repeat'::text, 'unscheduled'::text, 'targeted'::text]))),
    CONSTRAINT org_events_event_type_check CHECK ((event_type = ANY (ARRAY['training'::text, 'drill'::text, 'inspection'::text, 'meeting'::text])))
);


--
-- Name: organizations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.organizations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    bin text,
    phone text,
    email text,
    address text,
    contact_person text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: profiles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.profiles (
    id uuid NOT NULL,
    email text,
    role text DEFAULT 'user'::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    full_name text,
    is_active boolean DEFAULT true NOT NULL,
    phone text,
    "position" text,
    department text,
    bio text,
    avatar_url text,
    organization_id uuid,
    CONSTRAINT profiles_role_check CHECK ((role = ANY (ARRAY['user'::text, 'admin'::text, 'super_admin'::text, 'org_admin'::text])))
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: ai_analysis id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_analysis ALTER COLUMN id SET DEFAULT nextval('public.ai_analysis_id_seq'::regclass);


--
-- Name: knowledge_base id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.knowledge_base ALTER COLUMN id SET DEFAULT nextval('public.knowledge_base_id_seq'::regclass);


--
-- Data for Name: users; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.users (instance_id, id, aud, role, email, encrypted_password, email_confirmed_at, invited_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, email_change_token_new, email_change, email_change_sent_at, last_sign_in_at, raw_app_meta_data, raw_user_meta_data, is_super_admin, created_at, updated_at, phone, phone_confirmed_at, phone_change, phone_change_token, phone_change_sent_at, email_change_token_current, email_change_confirm_status, banned_until, reauthentication_token, reauthentication_sent_at, is_sso_user, deleted_at, is_anonymous) FROM stdin;
00000000-0000-0000-0000-000000000000	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	authenticated	authenticated	adminaq@gmail.com	$2a$06$yYa/qV4.t1okeM3blPrsIeQkxWlQjLOBuEdmyYHMSMEglL3OImSoy	2026-02-25 09:27:32.913595+00	\N		\N		\N			\N	2026-05-19 07:13:36.555998+00	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-02-25 09:27:32.900512+00	2026-05-19 07:13:36.610433+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	a209cb79-9286-42cd-b635-636ef4bd27cb	authenticated	authenticated	sky@gmail.com	$2a$10$rzoS10QG7fVHHxnUJz48pOiDhPQzYfTeM32Yz2GJDpHPaXmhD/5fi	2026-05-19 11:08:56.042618+00	\N	\N	\N	\N	\N	\N	\N	\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-05-19 11:08:56.042618+00	2026-05-19 11:08:56.042618+00	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	f	\N	f
00000000-0000-0000-0000-000000000000	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	authenticated	authenticated	sky1@gmail.com	$2a$10$s11Oxpfrfeb0GxpeXhbEU.GmwTpPdLfqUrCTMrtlxdIUxQy32krAq	2026-05-19 11:10:59.34605+00	\N	\N	\N	\N	\N	\N	\N	\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-05-19 11:10:59.34605+00	2026-05-19 11:10:59.34605+00	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	f	\N	f
00000000-0000-0000-0000-000000000000	7352c1fa-5021-41e8-89ed-7d2d70a98e01	authenticated	authenticated	de@gmail.com	$2a$10$nV/YQuNpWrazoHVVf0YBK.lOPppHEXsQJJfiP3J1jHH7x1wApAh36	2026-05-20 06:59:46.106618+00	\N	\N	\N	\N	\N	\N	\N	\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-05-20 06:59:46.106618+00	2026-05-20 06:59:46.106618+00	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	f	\N	f
00000000-0000-0000-0000-000000000000	7fda9576-3492-4578-b09c-1fb7bc03bcad	authenticated	authenticated	er@gmail.com	$2a$10$Riz2pU/J8LCSiDDPQVYV8ecjcm2hnqb.4ddk34gy6.KaKmty4aL5G	2026-05-20 10:05:07.677823+00	\N	\N	\N	\N	\N	\N	\N	\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-05-20 10:05:07.677823+00	2026-05-20 10:05:07.677823+00	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	f	\N	f
\.


--
-- Data for Name: ai_analysis; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ai_analysis (id, employee_id, organization_id, quiz_result_id, course_name, score, weak_topics, recommendation, risk_level, summary, created_at) FROM stdin;
1	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	5dfcd98f-9a4b-49a5-9840-0a622764522a	Обучения по ПТМ (Пожарно-технический-минимум)	60	{"Итоговая проверка знаний по ПТМ","Протокол заседания комиссии"}	Рекомендуется пройти дополнительное обучение по процедурам итоговой проверки знаний и составлению протоколов заседаний комиссии. Также полезно ознакомиться с нормативными документами, регулирующими эти процессы.	medium	Сотрудник продемонстрировал средний уровень знаний по пожарной безопасности, допустив ошибки в ключевых аспектах итоговой проверки и протоколирования. Рекомендуется дополнительное обучение для повышения квалификации.	2026-05-20 10:09:15.440607+00
2	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	bd9c3871-392a-491c-8d3c-91f329b8a759	Обучения по ПТМ (Пожарно-технический-минимум)	40	{"программа ПТМ для руководителей","организационные основы пожарной безопасности","удостоверение по пожарной безопасности"}	Рекомендуется пройти дополнительное обучение по программе ПТМ, уделяя внимание организационным основам пожарной безопасности и срокам действия удостоверений. Также полезно ознакомиться с официальными документами и нормативами в области пожарной безопасности.	high	Сотрудник продемонстрировал низкий уровень знаний по ключевым аспектам пожарной безопасности, что требует немедленного внимания и дополнительного обучения.	2026-05-20 10:09:19.045887+00
3	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	249e34ad-b60e-459d-9496-7614aaf4554d	Обучения по ПТМ (Пожарно-технический-минимум)	60	{"противопожарный инструктаж","методы обучения по ПТМ"}	Рекомендуется пройти дополнительное обучение по противопожарному инструктажу и методам проведения обучения по ПТМ. Также полезно ознакомиться с нормативными документами в области пожарной безопасности.	medium	Сотрудник продемонстрировал средний уровень знаний по пожарной безопасности, допустив ошибки в понимании задач противопожарного инструктажа и методов обучения.	2026-05-20 10:09:23.244881+00
4	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	441db675-afbb-4824-9f17-6c560f14d9bd	Обучения по ПТМ (Пожарно-технический-минимум)	20	{"ПТМ без отрыва от производства","Частота обучения ПТМ","Проверка знаний после ПТМ","Выдача удостоверений после ПТМ"}	Рекомендуется пройти дополнительное обучение по основам пожарной безопасности, уделяя особое внимание процедурам ПТМ без отрыва от производства. Также стоит ознакомиться с нормативными документами, касающимися частоты обучения и проверки знаний.	high	Сотрудник продемонстрировал низкий уровень знаний по пожарной безопасности, что требует немедленного внимания и дополнительного обучения.	2026-05-20 10:09:28.009353+00
5	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	b1fb9222-d586-4cf4-90ab-e8602365304a	Обучения по ПТМ (Пожарно-технический-минимум)	0	{инструктажи,"экзамены по ПТМ","повторные инструктажи","основания для внепланового инструктажа"}	Рекомендуется пройти повторное обучение по пожарной безопасности, уделяя особое внимание правилам проведения инструктажей, требованиям к вводному инструктажу и процедурам, связанным с экзаменами по ПТМ. Также полезно ознакомиться с нормативными документами, касающимися пожарной безопасности на рабочем месте.	high	Сотрудник не продемонстрировал знаний в области пожарной безопасности, что указывает на высокий уровень риска. Необходимо срочно провести дополнительное обучение.	2026-05-20 10:09:31.624519+00
6	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	254df416-767c-44a4-96fe-424e4b4e421b	Обучения по ПТМ (Пожарно-технический-минимум)	60	{"инструктаж для вновь принятых работников","процедура повторного инструктажа"}	Рекомендуется пройти дополнительное обучение по основам пожарной безопасности, уделяя особое внимание инструктажам и их процедурам. Также полезно ознакомиться с внутренними регламентами компании по этому вопросу.	medium	Сотрудник продемонстрировал средний уровень знаний по пожарной безопасности с результатом 60%. Необходимо улучшить понимание процедур инструктажей.	2026-05-20 10:09:34.781827+00
7	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	16117ef4-3d54-4f42-9d28-0c92eacff379	Обучения по ПТМ (Пожарно-технический-минимум)	40	{"Обязанности по ПТМ","Частота повторного обучения ПТМ","Документы по результатам обучения ПТМ"}	Рекомендуется пройти дополнительное обучение по обязанностям сотрудников в области ПТМ, частоте повторного обучения и документам, выдаваемым по результатам обучения. Также стоит ознакомиться с актуальными нормативными актами и инструкциями.	high	Сотрудник продемонстрировал низкий уровень знаний по пожарной безопасности, что требует немедленного внимания и дополнительного обучения.	2026-05-20 10:09:37.754877+00
8	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	9cb24b92-cdd7-4e72-9c79-ae6b9063f5a3	Обучения по ПТМ (Пожарно-технический-минимум)	40	{"Обучение мерам пожарной безопасности","Формы обучения во внеурочное время","Сезонные факторы пожарной опасности"}	Рекомендуется пройти дополнительные занятия по основам пожарной безопасности, включая лекции и практические занятия. Также полезно изучить материалы о сезонных изменениях в пожарной опасности и формах обучения.	high	Сотрудник продемонстрировал низкий уровень знаний по пожарной безопасности, правильно ответив только на 40% вопросов. Необходимо усилить обучение по ключевым темам.	2026-05-20 10:09:42.97704+00
9	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	76164654-ad61-4f8e-8006-2defb8af7215	Обучения по ПТМ (Пожарно-технический-минимум)	60	{"противопожарный режим","пожароопасные свойства сырья и материалов"}	Рекомендуется пройти дополнительное обучение по противопожарному режиму и пожароопасным свойствам материалов. Также стоит изучить внутренние правила и требования пожарной безопасности, чтобы лучше понимать риски и меры предосторожности.	medium	Сотрудник продемонстрировал недостаточные знания в области противопожарного режима и пожароопасных свойств материалов, что требует дополнительного обучения.	2026-05-20 10:09:46.194765+00
10	7fda9576-3492-4578-b09c-1fb7bc03bcad	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	18535f0a-5b44-4a02-84bb-941fb6e503b0	Обучения по ПТМ (Пожарно-технический-минимум)	20	{"различия программ ПТМ","горение и пожаровзрывоопасные свойства веществ","организационные меры в ПТМ","расширенная программа ПТМ для ответственных лиц"}	Рекомендуется пройти дополнительное обучение по всем темам, связанным с пожарной безопасностью, включая теоретические и практические занятия. Также полезно ознакомиться с нормативными документами и инструкциями по ПТМ.	high	Сотрудник продемонстрировал низкий уровень знаний по пожарной безопасности, что требует немедленного внимания и дополнительного обучения.	2026-05-20 10:09:49.737491+00
\.


--
-- Data for Name: backend_meta; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.backend_meta (id, created_at) FROM stdin;
1	2026-05-19 09:32:00.779288+00
\.


--
-- Data for Name: briefing_records; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.briefing_records (id, organization_id, event_id, employee_id, employee_name, "position", briefing_kind, instructor_name, instructor_id, date_conducted, employee_signed, employee_signed_at, instructor_signed, instructor_signed_at, created_at, updated_at) FROM stdin;
19d9a6aa-edad-4267-9aec-1e42082d3c10	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	aa171e8c-cae2-4d11-94e7-14a239ba6312	7352c1fa-5021-41e8-89ed-7d2d70a98e01	Дэриэль	Сварщик	introductory	Инструктор по ПБ	\N	2026-05-20	t	2026-05-20 07:07:23.993751+00	t	2026-05-20 07:10:38.962982+00	2026-05-20 07:07:23.986169+00	2026-05-20 07:10:38.962982+00
a753e319-af26-4aaf-818b-170bd587b205	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	b160c9ef-db70-469c-9448-2b2ab955fec0	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	Миронов Антон	Не указана	repeat	Инструктор по ПБ	\N	2026-06-10	t	2026-06-10 06:55:32.949448+00	f	\N	2026-06-10 06:55:32.931942+00	2026-06-10 06:55:32.949448+00
\.


--
-- Data for Name: briefing_videos; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.briefing_videos (id, course_id, briefing_kind, video_url, video_path, created_at, updated_at) FROM stdin;
aa00ceaa-0398-49e3-9703-f7e68c7f3227	09443e7f-eb91-4838-9622-c7c389a04960	introductory	http://localhost:8080/files/videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/introductory/1781074059877_bd1c36d71fd64a91bd44eb7ccb6cf6d0.mp4	videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/introductory/1781074059877_bd1c36d71fd64a91bd44eb7ccb6cf6d0.mp4	2026-06-10 06:47:39.940632+00	2026-06-10 06:47:39.940632+00
61923db2-1a1c-4250-a311-0496ae4be392	09443e7f-eb91-4838-9622-c7c389a04960	primary	http://localhost:8080/files/videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/primary/1781074065808_07088e5cb20146769d7c6699affdddb4.mp4	videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/primary/1781074065808_07088e5cb20146769d7c6699affdddb4.mp4	2026-06-10 06:47:45.862743+00	2026-06-10 06:47:45.862743+00
507b20c3-42ce-47f0-8a92-a02ff49460f9	09443e7f-eb91-4838-9622-c7c389a04960	repeat	http://localhost:8080/files/videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/repeat/1781074086082_db00b30bf45d49359f27ea80ef663a48.mp4	videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/repeat/1781074086082_db00b30bf45d49359f27ea80ef663a48.mp4	2026-06-10 06:48:06.117431+00	2026-06-10 06:48:06.117431+00
464fc1ea-4cf3-4699-b98a-d7aacccdfb02	09443e7f-eb91-4838-9622-c7c389a04960	targeted	http://localhost:8080/files/videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/targeted/1781074100151_57102baf7a484007b1588667b5593e43.mp4	videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/targeted/1781074100151_57102baf7a484007b1588667b5593e43.mp4	2026-06-10 06:48:20.185994+00	2026-06-10 06:48:20.185994+00
90d7e2bc-49e1-4c63-9a38-a6eabc19793e	09443e7f-eb91-4838-9622-c7c389a04960	unscheduled	http://localhost:8080/files/videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/unscheduled/1781074119012_689e4a2632484440b83e9cf1b1ff853b.mp4	videos/briefings/09443e7f-eb91-4838-9622-c7c389a04960/unscheduled/1781074119012_689e4a2632484440b83e9cf1b1ff853b.mp4	2026-06-10 06:48:39.039046+00	2026-06-10 06:48:39.039046+00
\.


--
-- Data for Name: contact_requests; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.contact_requests (id, name, email, company, message, status, created_at, updated_at, phone) FROM stdin;
b81e26a7-96d3-47ea-9ad9-bd78de2ad20e	Иванов Иван	ivan@gmail.com	TOO BIG BAG	Пройти обучения по Пожарной безопасности	new	2026-05-21 07:14:29.277536+00	2026-05-21 07:19:49.369022+00	
682f7028-9b49-44fe-a267-2b80cc04706f	Иванов Иван Иван	iv@gmailc.com	TOO PONCHIK	\N	new	2026-05-21 07:24:12.786394+00	2026-05-21 07:24:12.786394+00	+7 777 888 99 99
\.


--
-- Data for Name: course_assignments; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.course_assignments (id, user_id, course_id, assigned_by, assigned_at, expires_at, status, revoked_at, updated_at, training_completed_at) FROM stdin;
77857500-2bc5-48b9-ad1b-c9d82439e612	a209cb79-9286-42cd-b635-636ef4bd27cb	09443e7f-eb91-4838-9622-c7c389a04960	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-05-19 11:10:11.823519+00	\N	active	\N	2026-05-19 11:10:11.824708+00	\N
d4ccb404-3d2b-4843-b53e-6e1b923e4d6c	7352c1fa-5021-41e8-89ed-7d2d70a98e01	09443e7f-eb91-4838-9622-c7c389a04960	a209cb79-9286-42cd-b635-636ef4bd27cb	2026-05-20 08:59:32.135859+00	\N	active	\N	2026-05-20 09:43:38.676276+00	2026-05-20 09:43:38.676276+00
ab92b61a-20e8-41fc-967e-cc99865feaf5	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	09443e7f-eb91-4838-9622-c7c389a04960	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-05-20 05:43:57.403621+00	\N	active	\N	2026-05-20 09:58:37.430733+00	2026-05-20 09:58:37.430733+00
b7eea3a4-668a-4a86-abb6-77bad665c047	7fda9576-3492-4578-b09c-1fb7bc03bcad	09443e7f-eb91-4838-9622-c7c389a04960	a209cb79-9286-42cd-b635-636ef4bd27cb	2026-05-20 10:05:37.79936+00	\N	active	\N	2026-05-20 10:09:49.740311+00	2026-05-20 10:09:49.740311+00
\.


--
-- Data for Name: course_materials; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.course_materials (id, course_id, name, file_url, file_type, file_size, created_at) FROM stdin;
\.


--
-- Data for Name: courses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.courses (id, title, description, status, cover_url, created_by, created_at, updated_at, is_briefing_course) FROM stdin;
09443e7f-eb91-4838-9622-c7c389a04960	Обучения по ПТМ (Пожарно-технический-минимум)	Приказ МВД РК № 777 \nОбучение пожарно-техническому минимуму проводится независимо \n\nот направления деятельности  \nШтраф статья 410 КоАП РК - от 5 до 300 МРП	published	\N	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-05-19 10:27:08.485952+00	2026-06-10 06:48:42.833513+00	t
\.


--
-- Data for Name: knowledge_base; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.knowledge_base (id, topic, content, article, created_at) FROM stdin;
1	эвакуация	Руководитель обязан обеспечить наличие планов эвакуации на каждом этаже. Пути эвакуации должны быть свободны.	п.15	2026-05-20 08:57:36.599271+00
2	огнетушители	Огнетушители проверяются не реже 1 раза в год. На каждые 100 кв.м. не менее 1 огнетушителя.	п.23	2026-05-20 08:57:36.599271+00
3	инструктаж	Первичный инструктаж — до начала работы. Повторный — не реже 1 раза в полгода.	п.8	2026-05-20 08:57:36.599271+00
4	сигнализация	Системы пожарной сигнализации проходят техобслуживание не реже 1 раза в квартал.	п.31	2026-05-20 08:57:36.599271+00
5	действия при пожаре	При пожаре: позвонить 101, эвакуировать людей, тушить первичными средствами.	п.12	2026-05-20 08:57:36.599271+00
\.


--
-- Data for Name: lesson_quiz_attempts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.lesson_quiz_attempts (id, user_id, lesson_id, quiz_id, answers, score, max_score, passed, completed_at) FROM stdin;
668d27e4-a081-43fa-b750-cce31c5c0e9d	7352c1fa-5021-41e8-89ed-7d2d70a98e01	107b5ff3-ac0b-4be3-a2c1-4ad73983d23b	8db2aee3-9501-48e9-93a0-c5cceed00ac8	[{"answer": "B", "question_id": "c4d195bd-24ad-4281-b83f-04a1e4c3acd0"}, {"answer": "C", "question_id": "33236276-00a0-4097-8d93-b0bde00368b0"}, {"answer": "A", "question_id": "a6e9b646-5a36-4c23-b917-773b79b4d8fb"}, {"answer": "B", "question_id": "242db456-13ce-43d9-819e-a011a883828d"}, {"answer": "C", "question_id": "d4039a2b-8425-4deb-8d6d-a898c7b0646f"}]	5	5	t	2026-05-20 09:11:13.040634+00
2c863afe-9913-4496-8f96-13f86014f2c4	7352c1fa-5021-41e8-89ed-7d2d70a98e01	7612594a-faf7-4cd5-b6d0-2219e98d8ec8	d2e3058b-343b-4116-a6b9-e522e5f7450c	[{"answer": "B", "question_id": "6f882a04-450a-4e32-bcec-137205acb6d2"}, {"answer": "C", "question_id": "fa403c28-1167-46c1-b63e-9e0fa694820d"}, {"answer": "B", "question_id": "5a7d28a0-f1bc-4064-9ac5-e7e32f515fa6"}, {"answer": "B", "question_id": "6fdf28cd-5a8f-4887-ae98-8197a01fb19b"}, {"answer": "C", "question_id": "706ad8b6-a85a-4e41-8bfb-5d3ac3f02117"}]	5	5	t	2026-05-20 09:12:47.714618+00
e6ac0cae-4368-45aa-98b8-4b9fb40f2514	7352c1fa-5021-41e8-89ed-7d2d70a98e01	88f0e9b7-ea8c-4d92-af5a-4b7a03ac1592	2ba232c8-dc09-4752-8a2f-9bb33e902059	[{"answer": "B", "question_id": "48f50d24-d075-418e-a8ec-d5e337e02da7"}, {"answer": "A", "question_id": "d45143a9-2eca-4183-8d39-2b2062ca2e0f"}, {"answer": "A", "question_id": "bdefc34d-af47-4c4b-96e6-1be17f8e577b"}, {"answer": "C", "question_id": "0d848c57-2853-44ee-a883-d6947a09d460"}, {"answer": "B", "question_id": "79fea14a-4003-4d40-9260-cbf6a5f73f65"}]	5	5	t	2026-05-20 09:13:45.484039+00
2f2481f1-8d91-45a7-a7ba-946524cab206	7352c1fa-5021-41e8-89ed-7d2d70a98e01	23068a0a-5c34-4919-805c-c9c529901ce5	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	[{"answer": "B", "question_id": "27296a0f-ee27-4c1c-9743-24e55bcd2a97"}, {"answer": "B", "question_id": "8a8c4e45-f00b-444d-b31a-fe961374a928"}, {"answer": "B", "question_id": "bdab9067-3dad-4b1d-9161-38e547dbd8f8"}, {"answer": "A", "question_id": "cbdeb458-04eb-4ff8-9400-ac060b1f1309"}, {"answer": "A", "question_id": "24620ccc-f01b-4c0a-b0d5-30077e102076"}]	3	5	f	2026-05-20 09:16:01.619645+00
33d1e22a-73b3-4df0-a40a-c846c0114922	7352c1fa-5021-41e8-89ed-7d2d70a98e01	23068a0a-5c34-4919-805c-c9c529901ce5	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	[{"answer": "B", "question_id": "27296a0f-ee27-4c1c-9743-24e55bcd2a97"}, {"answer": "B", "question_id": "8a8c4e45-f00b-444d-b31a-fe961374a928"}, {"answer": "B", "question_id": "bdab9067-3dad-4b1d-9161-38e547dbd8f8"}, {"answer": "A", "question_id": "cbdeb458-04eb-4ff8-9400-ac060b1f1309"}, {"answer": "A", "question_id": "24620ccc-f01b-4c0a-b0d5-30077e102076"}]	3	5	f	2026-05-20 09:16:04.013717+00
452d264b-0422-46a0-8cfa-297c88751145	7352c1fa-5021-41e8-89ed-7d2d70a98e01	c252c5c6-c596-427c-9654-1a2d6a35fe97	5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	[{"answer": "A", "question_id": "19f46318-c417-45e6-bf37-a2727b1b1513"}, {"answer": "B", "question_id": "2444ff29-0d4a-40bd-9ab9-577b01f33297"}, {"answer": "C", "question_id": "0e80d4f6-a566-4c9e-ba48-7a159123ff60"}, {"answer": "C", "question_id": "d2d25ac5-251e-48e0-b6ea-ba2785e9431d"}, {"answer": "B", "question_id": "3c6fce19-eb06-4f0d-8a4b-423c3ba276fc"}]	0	5	f	2026-05-20 09:20:47.138467+00
53906bd4-9f76-48e2-bd44-3c6877916006	7352c1fa-5021-41e8-89ed-7d2d70a98e01	0defd165-3d67-4bfa-ac95-552e407afb0a	25046c2b-6fed-4240-bf24-bb66a7fc8028	[{"answer": "A", "question_id": "5504f969-aa9b-4234-a290-8789b342e157"}, {"answer": "B", "question_id": "5b08dda8-d511-4efe-8dd2-cb89b9074098"}, {"answer": "C", "question_id": "cea05028-ee7c-40ee-9ac9-16b7123330ae"}, {"answer": "B", "question_id": "b92cbad0-8562-4b2f-908c-c2d3b0757698"}, {"answer": "C", "question_id": "a4735f76-2777-492a-80d6-3ead034dae28"}]	2	5	f	2026-05-20 09:24:08.049894+00
f19d3289-8714-4fd4-90ec-21b24e14cf7a	7352c1fa-5021-41e8-89ed-7d2d70a98e01	e4ff4692-a1cb-41c5-9b12-3c5ce1846f84	8e48c271-d7cd-4a01-82ec-71480211dd46	[{"answer": "A", "question_id": "a789b3ad-2120-4adf-99bf-7d14672883d0"}, {"answer": "A", "question_id": "4482ee64-1801-4b36-980c-a3801483b0da"}, {"answer": "A", "question_id": "806cdfa0-ae54-412b-8800-836c50cb8cf0"}, {"answer": "B", "question_id": "637a0bcb-91ac-4478-9399-911af5cd124f"}, {"answer": "C", "question_id": "e586c105-6e51-4238-9efe-7cb0b5b09970"}]	1	5	f	2026-05-20 09:30:16.399716+00
3e2595e4-a9bc-4b8f-b0e0-0269bef5f9f6	7352c1fa-5021-41e8-89ed-7d2d70a98e01	411ea650-90b6-496e-9cbe-70f61fc7aa5a	0abcd72c-bf6a-4f08-a909-de61520c4da1	[{"answer": "A", "question_id": "9e9f38ac-ec2c-4993-a224-7091dbde606b"}, {"answer": "B", "question_id": "cf3c5548-0ad0-4ac7-8758-07f774eab220"}, {"answer": "B", "question_id": "0d356a18-59d6-4f8c-84b5-58e7f2216370"}, {"answer": "C", "question_id": "b40d0b79-4841-4daf-b91f-a2ab2ed8ac2c"}, {"answer": "A", "question_id": "9a4b8e2e-cc02-4aca-b995-f0bf92f406e8"}]	2	5	f	2026-05-20 09:30:29.517932+00
cc00e926-9541-4d75-88d8-9e3dc60b555f	7352c1fa-5021-41e8-89ed-7d2d70a98e01	a394628f-a661-44a0-a721-a6774dbea1f6	88880461-dcb7-44eb-b7aa-b085bf7e6468	[{"answer": "B", "question_id": "e57cdf83-28d0-45ac-ab40-1b5b0932f095"}, {"answer": "B", "question_id": "7ec6a75a-f717-4623-8f93-ea7527050f8c"}, {"answer": "B", "question_id": "5f99a63d-648b-42d2-85c4-14a055c1de51"}, {"answer": "C", "question_id": "e4bbbdd2-f7fd-4949-b9d0-c0b5789d1c2b"}, {"answer": "C", "question_id": "bab5cc99-79ec-4b2e-bd7f-40a2f963e754"}]	0	5	f	2026-05-20 09:32:01.191583+00
6a177714-bee7-461d-bce6-8ffc1fbd079c	7352c1fa-5021-41e8-89ed-7d2d70a98e01	04b035e0-508e-4e0a-8079-ca7aa6b9a729	d6d459f8-5143-488a-b588-d55ae9de0938	[{"answer": "A", "question_id": "1783a9fb-5956-49f5-9160-2ecc921903d2"}, {"answer": "B", "question_id": "584edd81-8dff-4b82-b91a-e301ce65b01b"}, {"answer": "C", "question_id": "ae7d0247-596e-47a1-8898-8d4a5c6bac62"}, {"answer": "B", "question_id": "2a1ac6cf-6134-4630-bdfa-69ad3c65e8cd"}, {"answer": "C", "question_id": "452fcf12-7d4e-48fa-9406-db490536037a"}]	3	5	f	2026-05-20 09:32:12.016717+00
780a83dd-b0c4-4a26-bb4b-a3b244da8e0b	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	107b5ff3-ac0b-4be3-a2c1-4ad73983d23b	8db2aee3-9501-48e9-93a0-c5cceed00ac8	[{"answer": "A", "question_id": "c4d195bd-24ad-4281-b83f-04a1e4c3acd0"}, {"answer": "B", "question_id": "33236276-00a0-4097-8d93-b0bde00368b0"}, {"answer": "B", "question_id": "a6e9b646-5a36-4c23-b917-773b79b4d8fb"}, {"answer": "B", "question_id": "242db456-13ce-43d9-819e-a011a883828d"}, {"answer": "A", "question_id": "d4039a2b-8425-4deb-8d6d-a898c7b0646f"}]	1	5	f	2026-05-20 09:53:36.509231+00
d6d1f61b-a6bf-4bdf-9a7f-f3ac206972ef	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	7612594a-faf7-4cd5-b6d0-2219e98d8ec8	d2e3058b-343b-4116-a6b9-e522e5f7450c	[{"answer": "A", "question_id": "6f882a04-450a-4e32-bcec-137205acb6d2"}, {"answer": "A", "question_id": "fa403c28-1167-46c1-b63e-9e0fa694820d"}, {"answer": "C", "question_id": "5a7d28a0-f1bc-4064-9ac5-e7e32f515fa6"}, {"answer": "B", "question_id": "6fdf28cd-5a8f-4887-ae98-8197a01fb19b"}, {"answer": "C", "question_id": "706ad8b6-a85a-4e41-8bfb-5d3ac3f02117"}]	2	5	f	2026-05-20 09:53:52.026609+00
e1373bc2-db6f-4d02-bbf0-381a90bdb4cc	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	88f0e9b7-ea8c-4d92-af5a-4b7a03ac1592	2ba232c8-dc09-4752-8a2f-9bb33e902059	[{"answer": "A", "question_id": "48f50d24-d075-418e-a8ec-d5e337e02da7"}, {"answer": "A", "question_id": "d45143a9-2eca-4183-8d39-2b2062ca2e0f"}, {"answer": "C", "question_id": "bdefc34d-af47-4c4b-96e6-1be17f8e577b"}, {"answer": "C", "question_id": "0d848c57-2853-44ee-a883-d6947a09d460"}, {"answer": "C", "question_id": "79fea14a-4003-4d40-9260-cbf6a5f73f65"}]	2	5	f	2026-05-20 09:54:06.834025+00
7c6e192d-3c72-45c9-b8d6-9eaeef85ecda	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	23068a0a-5c34-4919-805c-c9c529901ce5	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	[{"answer": "C", "question_id": "27296a0f-ee27-4c1c-9743-24e55bcd2a97"}, {"answer": "B", "question_id": "8a8c4e45-f00b-444d-b31a-fe961374a928"}, {"answer": "C", "question_id": "bdab9067-3dad-4b1d-9161-38e547dbd8f8"}, {"answer": "C", "question_id": "cbdeb458-04eb-4ff8-9400-ac060b1f1309"}, {"answer": "C", "question_id": "24620ccc-f01b-4c0a-b0d5-30077e102076"}]	1	5	f	2026-05-20 09:54:34.529228+00
823c83e5-2334-4f8c-8215-243c0b4ea2bb	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	c252c5c6-c596-427c-9654-1a2d6a35fe97	5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	[{"answer": "B", "question_id": "19f46318-c417-45e6-bf37-a2727b1b1513"}, {"answer": "A", "question_id": "2444ff29-0d4a-40bd-9ab9-577b01f33297"}, {"answer": "C", "question_id": "0e80d4f6-a566-4c9e-ba48-7a159123ff60"}, {"answer": "C", "question_id": "d2d25ac5-251e-48e0-b6ea-ba2785e9431d"}, {"answer": "A", "question_id": "3c6fce19-eb06-4f0d-8a4b-423c3ba276fc"}]	3	5	f	2026-05-20 09:54:50.264877+00
68947cdc-03b3-42f3-9398-08d674c7aaf1	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	0defd165-3d67-4bfa-ac95-552e407afb0a	25046c2b-6fed-4240-bf24-bb66a7fc8028	[{"answer": "C", "question_id": "5504f969-aa9b-4234-a290-8789b342e157"}, {"answer": "A", "question_id": "5b08dda8-d511-4efe-8dd2-cb89b9074098"}, {"answer": "A", "question_id": "cea05028-ee7c-40ee-9ac9-16b7123330ae"}, {"answer": "B", "question_id": "b92cbad0-8562-4b2f-908c-c2d3b0757698"}, {"answer": "B", "question_id": "a4735f76-2777-492a-80d6-3ead034dae28"}]	1	5	f	2026-05-20 09:55:27.93739+00
bf537893-03eb-4c15-8572-0b19967aba9a	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	e4ff4692-a1cb-41c5-9b12-3c5ce1846f84	8e48c271-d7cd-4a01-82ec-71480211dd46	[{"answer": "C", "question_id": "a789b3ad-2120-4adf-99bf-7d14672883d0"}, {"answer": "B", "question_id": "4482ee64-1801-4b36-980c-a3801483b0da"}, {"answer": "A", "question_id": "806cdfa0-ae54-412b-8800-836c50cb8cf0"}, {"answer": "B", "question_id": "637a0bcb-91ac-4478-9399-911af5cd124f"}, {"answer": "B", "question_id": "e586c105-6e51-4238-9efe-7cb0b5b09970"}]	2	5	f	2026-05-20 09:55:40.274791+00
76005fe2-11e0-4f1d-b839-e3c528fbf28a	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	411ea650-90b6-496e-9cbe-70f61fc7aa5a	0abcd72c-bf6a-4f08-a909-de61520c4da1	[{"answer": "A", "question_id": "9e9f38ac-ec2c-4993-a224-7091dbde606b"}, {"answer": "A", "question_id": "cf3c5548-0ad0-4ac7-8758-07f774eab220"}, {"answer": "B", "question_id": "0d356a18-59d6-4f8c-84b5-58e7f2216370"}, {"answer": "B", "question_id": "b40d0b79-4841-4daf-b91f-a2ab2ed8ac2c"}, {"answer": "B", "question_id": "9a4b8e2e-cc02-4aca-b995-f0bf92f406e8"}]	0	5	f	2026-05-20 09:55:57.013833+00
9f37fe73-39a4-4c24-8557-782075e958e8	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	a394628f-a661-44a0-a721-a6774dbea1f6	88880461-dcb7-44eb-b7aa-b085bf7e6468	[{"answer": "A", "question_id": "e57cdf83-28d0-45ac-ab40-1b5b0932f095"}, {"answer": "C", "question_id": "7ec6a75a-f717-4623-8f93-ea7527050f8c"}, {"answer": "B", "question_id": "5f99a63d-648b-42d2-85c4-14a055c1de51"}, {"answer": "A", "question_id": "e4bbbdd2-f7fd-4949-b9d0-c0b5789d1c2b"}, {"answer": "A", "question_id": "bab5cc99-79ec-4b2e-bd7f-40a2f963e754"}]	1	5	f	2026-05-20 09:56:12.080797+00
4f78c1b7-6c9d-4c34-babc-7c54a785206e	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	04b035e0-508e-4e0a-8079-ca7aa6b9a729	d6d459f8-5143-488a-b588-d55ae9de0938	[{"answer": "B", "question_id": "1783a9fb-5956-49f5-9160-2ecc921903d2"}, {"answer": "B", "question_id": "584edd81-8dff-4b82-b91a-e301ce65b01b"}, {"answer": "A", "question_id": "ae7d0247-596e-47a1-8898-8d4a5c6bac62"}, {"answer": "A", "question_id": "2a1ac6cf-6134-4630-bdfa-69ad3c65e8cd"}, {"answer": "B", "question_id": "452fcf12-7d4e-48fa-9406-db490536037a"}]	2	5	f	2026-05-20 09:56:26.93776+00
249e34ad-b60e-459d-9496-7614aaf4554d	7fda9576-3492-4578-b09c-1fb7bc03bcad	107b5ff3-ac0b-4be3-a2c1-4ad73983d23b	8db2aee3-9501-48e9-93a0-c5cceed00ac8	[{"answer": "B", "question_id": "c4d195bd-24ad-4281-b83f-04a1e4c3acd0"}, {"answer": "C", "question_id": "33236276-00a0-4097-8d93-b0bde00368b0"}, {"answer": "B", "question_id": "a6e9b646-5a36-4c23-b917-773b79b4d8fb"}, {"answer": "B", "question_id": "242db456-13ce-43d9-819e-a011a883828d"}, {"answer": "A", "question_id": "d4039a2b-8425-4deb-8d6d-a898c7b0646f"}]	3	5	f	2026-05-20 10:06:07.551916+00
254df416-767c-44a4-96fe-424e4b4e421b	7fda9576-3492-4578-b09c-1fb7bc03bcad	7612594a-faf7-4cd5-b6d0-2219e98d8ec8	d2e3058b-343b-4116-a6b9-e522e5f7450c	[{"answer": "A", "question_id": "6f882a04-450a-4e32-bcec-137205acb6d2"}, {"answer": "C", "question_id": "fa403c28-1167-46c1-b63e-9e0fa694820d"}, {"answer": "B", "question_id": "5a7d28a0-f1bc-4064-9ac5-e7e32f515fa6"}, {"answer": "B", "question_id": "6fdf28cd-5a8f-4887-ae98-8197a01fb19b"}, {"answer": "A", "question_id": "706ad8b6-a85a-4e41-8bfb-5d3ac3f02117"}]	3	5	f	2026-05-20 10:06:27.497095+00
16117ef4-3d54-4f42-9d28-0c92eacff379	7fda9576-3492-4578-b09c-1fb7bc03bcad	88f0e9b7-ea8c-4d92-af5a-4b7a03ac1592	2ba232c8-dc09-4752-8a2f-9bb33e902059	[{"answer": "B", "question_id": "48f50d24-d075-418e-a8ec-d5e337e02da7"}, {"answer": "C", "question_id": "d45143a9-2eca-4183-8d39-2b2062ca2e0f"}, {"answer": "A", "question_id": "bdefc34d-af47-4c4b-96e6-1be17f8e577b"}, {"answer": "A", "question_id": "0d848c57-2853-44ee-a883-d6947a09d460"}, {"answer": "C", "question_id": "79fea14a-4003-4d40-9260-cbf6a5f73f65"}]	2	5	f	2026-05-20 10:06:40.424725+00
441db675-afbb-4824-9f17-6c560f14d9bd	7fda9576-3492-4578-b09c-1fb7bc03bcad	23068a0a-5c34-4919-805c-c9c529901ce5	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	[{"answer": "A", "question_id": "27296a0f-ee27-4c1c-9743-24e55bcd2a97"}, {"answer": "B", "question_id": "8a8c4e45-f00b-444d-b31a-fe961374a928"}, {"answer": "B", "question_id": "bdab9067-3dad-4b1d-9161-38e547dbd8f8"}, {"answer": "B", "question_id": "cbdeb458-04eb-4ff8-9400-ac060b1f1309"}, {"answer": "C", "question_id": "24620ccc-f01b-4c0a-b0d5-30077e102076"}]	1	5	f	2026-05-20 10:06:55.277071+00
76164654-ad61-4f8e-8006-2defb8af7215	7fda9576-3492-4578-b09c-1fb7bc03bcad	c252c5c6-c596-427c-9654-1a2d6a35fe97	5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	[{"answer": "B", "question_id": "19f46318-c417-45e6-bf37-a2727b1b1513"}, {"answer": "B", "question_id": "2444ff29-0d4a-40bd-9ab9-577b01f33297"}, {"answer": "A", "question_id": "0e80d4f6-a566-4c9e-ba48-7a159123ff60"}, {"answer": "B", "question_id": "d2d25ac5-251e-48e0-b6ea-ba2785e9431d"}, {"answer": "B", "question_id": "3c6fce19-eb06-4f0d-8a4b-423c3ba276fc"}]	3	5	f	2026-05-20 10:07:08.207087+00
bd9c3871-392a-491c-8d3c-91f329b8a759	7fda9576-3492-4578-b09c-1fb7bc03bcad	0defd165-3d67-4bfa-ac95-552e407afb0a	25046c2b-6fed-4240-bf24-bb66a7fc8028	[{"answer": "A", "question_id": "5504f969-aa9b-4234-a290-8789b342e157"}, {"answer": "A", "question_id": "5b08dda8-d511-4efe-8dd2-cb89b9074098"}, {"answer": "B", "question_id": "cea05028-ee7c-40ee-9ac9-16b7123330ae"}, {"answer": "A", "question_id": "b92cbad0-8562-4b2f-908c-c2d3b0757698"}, {"answer": "B", "question_id": "a4735f76-2777-492a-80d6-3ead034dae28"}]	2	5	f	2026-05-20 10:08:15.024644+00
18535f0a-5b44-4a02-84bb-941fb6e503b0	7fda9576-3492-4578-b09c-1fb7bc03bcad	e4ff4692-a1cb-41c5-9b12-3c5ce1846f84	8e48c271-d7cd-4a01-82ec-71480211dd46	[{"answer": "B", "question_id": "a789b3ad-2120-4adf-99bf-7d14672883d0"}, {"answer": "A", "question_id": "4482ee64-1801-4b36-980c-a3801483b0da"}, {"answer": "C", "question_id": "806cdfa0-ae54-412b-8800-836c50cb8cf0"}, {"answer": "C", "question_id": "637a0bcb-91ac-4478-9399-911af5cd124f"}, {"answer": "B", "question_id": "e586c105-6e51-4238-9efe-7cb0b5b09970"}]	1	5	f	2026-05-20 10:08:28.099853+00
b1fb9222-d586-4cf4-90ab-e8602365304a	7fda9576-3492-4578-b09c-1fb7bc03bcad	411ea650-90b6-496e-9cbe-70f61fc7aa5a	0abcd72c-bf6a-4f08-a909-de61520c4da1	[{"answer": "A", "question_id": "9e9f38ac-ec2c-4993-a224-7091dbde606b"}, {"answer": "A", "question_id": "cf3c5548-0ad0-4ac7-8758-07f774eab220"}, {"answer": "A", "question_id": "0d356a18-59d6-4f8c-84b5-58e7f2216370"}, {"answer": "C", "question_id": "b40d0b79-4841-4daf-b91f-a2ab2ed8ac2c"}, {"answer": "B", "question_id": "9a4b8e2e-cc02-4aca-b995-f0bf92f406e8"}]	0	5	f	2026-05-20 10:08:44.309186+00
9cb24b92-cdd7-4e72-9c79-ae6b9063f5a3	7fda9576-3492-4578-b09c-1fb7bc03bcad	a394628f-a661-44a0-a721-a6774dbea1f6	88880461-dcb7-44eb-b7aa-b085bf7e6468	[{"answer": "A", "question_id": "e57cdf83-28d0-45ac-ab40-1b5b0932f095"}, {"answer": "A", "question_id": "7ec6a75a-f717-4623-8f93-ea7527050f8c"}, {"answer": "B", "question_id": "5f99a63d-648b-42d2-85c4-14a055c1de51"}, {"answer": "B", "question_id": "e4bbbdd2-f7fd-4949-b9d0-c0b5789d1c2b"}, {"answer": "B", "question_id": "bab5cc99-79ec-4b2e-bd7f-40a2f963e754"}]	2	5	f	2026-05-20 10:08:56.716083+00
5dfcd98f-9a4b-49a5-9840-0a622764522a	7fda9576-3492-4578-b09c-1fb7bc03bcad	04b035e0-508e-4e0a-8079-ca7aa6b9a729	d6d459f8-5143-488a-b588-d55ae9de0938	[{"answer": "B", "question_id": "1783a9fb-5956-49f5-9160-2ecc921903d2"}, {"answer": "B", "question_id": "584edd81-8dff-4b82-b91a-e301ce65b01b"}, {"answer": "A", "question_id": "ae7d0247-596e-47a1-8898-8d4a5c6bac62"}, {"answer": "B", "question_id": "2a1ac6cf-6134-4630-bdfa-69ad3c65e8cd"}, {"answer": "B", "question_id": "452fcf12-7d4e-48fa-9406-db490536037a"}]	3	5	f	2026-05-20 10:09:09.377484+00
\.


--
-- Data for Name: lesson_quiz_questions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.lesson_quiz_questions (id, quiz_id, order_index, question_text, option_a, option_b, option_c, correct_option) FROM stdin;
c4d195bd-24ad-4281-b83f-04a1e4c3acd0	8db2aee3-9501-48e9-93a0-c5cceed00ac8	1	В чем основное отличие противопожарного инструктажа от пожарно-технического минимума?	Инструктаж проводится только в учебном центре	ПТМ предполагает более углубленное изучение норм и практических навыков	Между ними нет различий	B
33236276-00a0-4097-8d93-b0bde00368b0	8db2aee3-9501-48e9-93a0-c5cceed00ac8	2	Кто проводит противопожарный инструктаж в организации?	Любой сотрудник организации	Только представители МЧС	Руководитель организации или ответственное лицо по пожарной безопасности	C
a6e9b646-5a36-4c23-b917-773b79b4d8fb	8db2aee3-9501-48e9-93a0-c5cceed00ac8	3	Что входит в задачи противопожарного инструктажа?	Изучение требований пожарной безопасности и действий при пожаре	Только изучение эвакуационных выходов	Только заполнение журнала инструктажа	A
242db456-13ce-43d9-819e-a011a883828d	8db2aee3-9501-48e9-93a0-c5cceed00ac8	4	Для кого обязателен противопожарный инструктаж?	Только для руководителей	Для всех работников организации	Только для работников производственных объектов	B
d4039a2b-8425-4deb-8d6d-a898c7b0646f	8db2aee3-9501-48e9-93a0-c5cceed00ac8	5	Как может проводиться обучение по пожарно-техническому минимуму?	Только с отрывом от производства	Только дистанционно	С отрывом и без отрыва от производства	C
6f882a04-450a-4e32-bcec-137205acb6d2	d2e3058b-343b-4116-a6b9-e522e5f7450c	1	Какой инструктаж обязателен для всех вновь принятых работников?	Повторный	Вводный	Целевой	B
fa403c28-1167-46c1-b63e-9e0fa694820d	d2e3058b-343b-4116-a6b9-e522e5f7450c	2	Где проводится первичный противопожарный инструктаж?	В учебном центре	В кабинете руководителя	Непосредственно на рабочем месте	C
5a7d28a0-f1bc-4064-9ac5-e7e32f515fa6	d2e3058b-343b-4116-a6b9-e522e5f7450c	3	Когда проводится внеплановый инструктаж?	Только один раз в год	При изменении технологического процесса или выявленных нарушениях	Только перед отпуском	B
6fdf28cd-5a8f-4887-ae98-8197a01fb19b	d2e3058b-343b-4116-a6b9-e522e5f7450c	4	Для каких работ требуется целевой инструктаж?	Для обычной офисной работы	Для разовых работ с повышенной пожарной опасностью	Только для новых сотрудников	B
706ad8b6-a85a-4e41-8bfb-5d3ac3f02117	d2e3058b-343b-4116-a6b9-e522e5f7450c	5	Что происходит, если работник показал неудовлетворительные знания после инструктажа?	Допускается к работе без повторной проверки	Освобождается от обучения	Проходит инструктаж повторно	C
48f50d24-d075-418e-a8ec-d5e337e02da7	2ba232c8-dc09-4752-8a2f-9bb33e902059	1	Где проходит ПТМ с отрывом от производства?	Только внутри организации	В учебном центре	Только онлайн	B
d45143a9-2eca-4183-8d39-2b2062ca2e0f	2ba232c8-dc09-4752-8a2f-9bb33e902059	2	Кто из перечисленных обязан проходить ПТМ с отрывом от производства?	Руководители организаций	Все временные работники	Только стажеры	A
bdefc34d-af47-4c4b-96e6-1be17f8e577b	2ba232c8-dc09-4752-8a2f-9bb33e902059	3	В какой срок проводится первичное обучение ПТМ после приема на работу?	В течение месяца	Через 6 месяцев	Через 3 года	A
0d848c57-2853-44ee-a883-d6947a09d460	2ba232c8-dc09-4752-8a2f-9bb33e902059	4	Как часто проводится повторное обучение ПТМ с отрывом от производства?	Каждый год	Один раз в 5 лет	Не реже одного раза в 3 года	C
79fea14a-4003-4d40-9260-cbf6a5f73f65	2ba232c8-dc09-4752-8a2f-9bb33e902059	5	Что получает работник после успешной сдачи проверки знаний ПТМ?	Допуск к отпуску	Удостоверение установленной формы	Только устное подтверждение	B
27296a0f-ee27-4c1c-9743-24e55bcd2a97	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	1	Где проводится ПТМ без отрыва от производства?	Только в учебном центре	Непосредственно в организации	Только в подразделениях МЧС	B
8a8c4e45-f00b-444d-b31a-fe961374a928	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	2	Кто может проводить обучение ПТМ внутри организации?	Любой сотрудник отдела кадров	Руководитель или ответственное лицо, прошедшее обучение в учебном центре	Только внешний инспектор	B
bdab9067-3dad-4b1d-9161-38e547dbd8f8	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	3	Как часто проводится обучение ПТМ без отрыва от производства?	Не реже одного раза в год	Один раз в 3 года	Один раз в 5 лет	A
cbdeb458-04eb-4ff8-9400-ac060b1f1309	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	4	Что является особенностью проверки знаний после ПТМ без отрыва от производства?	Проверка проводится комиссией организации	Экзамен проводится только в МЧС	Проверка отсутствует	A
24620ccc-f01b-4c0a-b0d5-30077e102076	7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	5	Выдается ли удостоверение после прохождения ПТМ без отрыва от производства?	Да, обязательно	Нет, удостоверение не выдается	Только руководителям	B
19f46318-c417-45e6-bf37-a2727b1b1513	5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	1	Что должно учитываться при проведении противопожарного инструктажа?	Только возраст работников	Специфика деятельности конкретной организации	Только стаж сотрудника	B
2444ff29-0d4a-40bd-9ab9-577b01f33297	5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	2	Что включает ознакомление с противопожарным режимом?	Изучение внутренних правил и требований пожарной безопасности	Только эвакуацию	Только работу с документацией	A
0e80d4f6-a566-4c9e-ba48-7a159123ff60	5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	3	Для чего работникам объясняются виды огнетушителей?	Для выбора по классу пожара	Только для хранения	Для бухгалтерского учета	A
d2d25ac5-251e-48e0-b6ea-ba2785e9431d	5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	4	Что входит в практическую часть инструктажа?	Только просмотр видео	Практическое занятие	Только заполнение журнала	B
3c6fce19-eb06-4f0d-8a4b-423c3ba276fc	5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	5	Зачем работникам изучать пожароопасные свойства сырья и материалов?	Чтобы понимать риски технологического процесса	Только для отчетности	Это касается только пожарных служб	A
a789b3ad-2120-4adf-99bf-7d14672883d0	8e48c271-d7cd-4a01-82ec-71480211dd46	1	Какова продолжительность программы ПТМ для лиц, ответственных за пожарную безопасность?	8 учебных часов	12 учебных часов	16 учебных часов	B
4482ee64-1801-4b36-980c-a3801483b0da	8e48c271-d7cd-4a01-82ec-71480211dd46	2	Чем эта программа отличается от программы для руководителей?	В ней отсутствуют практические занятия	Более глубоко изучается пожарная опасность производства и пожароопасные работы	Она проводится только дистанционно	B
806cdfa0-ae54-412b-8800-836c50cb8cf0	8e48c271-d7cd-4a01-82ec-71480211dd46	3	Что входит в изучение темы о горении и пожаровзрывоопасных свойствах веществ?	Категорирование помещений и показатели взрывопожароопасности	Только порядок эвакуации	Только оформление журнала	A
637a0bcb-91ac-4478-9399-911af5cd124f	8e48c271-d7cd-4a01-82ec-71480211dd46	4	Какие организационные меры рассматриваются в программе?	Пожарно-технические комиссии и противопожарная пропаганда	Только закупка оборудования	Только кадровый учет	A
e586c105-6e51-4238-9efe-7cb0b5b09970	8e48c271-d7cd-4a01-82ec-71480211dd46	5	Почему ответственное лицо проходит расширенную программу ПТМ?	Потому что затем может проводить ПТМ внутри организации	Только для работы с документами	Для получения должности руководителя	A
9e9f38ac-ec2c-4993-a224-7091dbde606b	0abcd72c-bf6a-4f08-a909-de61520c4da1	1	Что считается нарушением при проведении инструктажа?	Запись в журнале учета	Проведение практического занятия	Отсутствие фиксации инструктажа в журнале	C
cf3c5548-0ad0-4ac7-8758-07f774eab220	0abcd72c-bf6a-4f08-a909-de61520c4da1	2	Может ли новый сотрудник приступить к работе без вводного инструктажа?	Да, если имеет опыт работы	Нет, к работе он не допускается	Да, с разрешения коллег	B
0d356a18-59d6-4f8c-84b5-58e7f2216370	0abcd72c-bf6a-4f08-a909-de61520c4da1	3	Что происходит, если работник не сдал экзамен ПТМ?	Допускается к самостоятельной работе	Освобождается от повторной проверки	Должен пройти повторную сдачу не позднее одного месяца	C
b40d0b79-4841-4daf-b91f-a2ab2ed8ac2c	0abcd72c-bf6a-4f08-a909-de61520c4da1	4	Почему отсутствие повторного инструктажа на производственном объекте в течение двух лет является нарушением?	Потому что он должен проводиться ежегодно	Потому что он проводится раз в пять лет	Потому что он необязателен	A
9a4b8e2e-cc02-4aca-b995-f0bf92f406e8	0abcd72c-bf6a-4f08-a909-de61520c4da1	5	Что может стать основанием для внепланового инструктажа?	Выявленные нарушения и результаты обследования организации	Только смена директора	Только прием нового сотрудника	A
1783a9fb-5956-49f5-9160-2ecc921903d2	d6d459f8-5143-488a-b588-d55ae9de0938	1	Кто проводит итоговую проверку знаний после ПТМ с отрывом от производства?	Квалификационная комиссия	Любой сотрудник организации	Только отдел кадров	A
584edd81-8dff-4b82-b91a-e301ce65b01b	d6d459f8-5143-488a-b588-d55ae9de0938	2	Какое минимальное количество человек должно быть в квалификационной комиссии?	Два	Три	Пять	B
ae7d0247-596e-47a1-8898-8d4a5c6bac62	d6d459f8-5143-488a-b588-d55ae9de0938	3	Что фиксируется в протоколе заседания комиссии?	Только дата обучения	Фамилия, должность, организация, причина обучения и результат	Только итог экзамена	B
2a1ac6cf-6134-4630-bdfa-69ad3c65e8cd	d6d459f8-5143-488a-b588-d55ae9de0938	4	Что получает работник при успешной сдаче экзамена?	Устное подтверждение	Удостоверение по проверке знаний ПТМ	Освобождение от всех инструктажей	B
452fcf12-7d4e-48fa-9406-db490536037a	d6d459f8-5143-488a-b588-d55ae9de0938	5	Что происходит, если экзамен не сдан?	Работник может продолжать самостоятельную работу	Повторная сдача проводится не позднее одного месяца	Обучение автоматически аннулируется	B
5504f969-aa9b-4234-a290-8789b342e157	25046c2b-6fed-4240-bf24-bb66a7fc8028	1	Где руководители организаций проходят ПТМ?	В учебном центре	Только внутри компании	Только дистанционно	A
5b08dda8-d511-4efe-8dd2-cb89b9074098	25046c2b-6fed-4240-bf24-bb66a7fc8028	2	На сколько учебных часов рассчитана программа ПТМ для руководителей?	6 часов	8 часов	12 часов	B
cea05028-ee7c-40ee-9ac9-16b7123330ae	25046c2b-6fed-4240-bf24-bb66a7fc8028	3	Что входит в организационные основы пожарной безопасности для руководителей?	Обучение работников и ведение документации	Только эвакуация	Только работа с огнетушителем	A
b92cbad0-8562-4b2f-908c-c2d3b0757698	25046c2b-6fed-4240-bf24-bb66a7fc8028	4	Что входит в практические занятия программы ПТМ для руководителей?	Работа с огнетушителем и тренировка эвакуации	Только лекция по законодательству	Только письменный тест	A
a4735f76-2777-492a-80d6-3ead034dae28	25046c2b-6fed-4240-bf24-bb66a7fc8028	5	Что получает руководитель после успешной сдачи экзамена?	Удостоверение сроком действия три года	Постоянный допуск без повторного обучения	Только запись в журнале	A
e57cdf83-28d0-45ac-ab40-1b5b0932f095	88880461-dcb7-44eb-b7aa-b085bf7e6468	1	На кого распространяется обучение мерам пожарной безопасности по главе 3 Правил?	Только на работников организаций	Только на школьников	На население в целом	C
7ec6a75a-f717-4623-8f93-ea7527050f8c	88880461-dcb7-44eb-b7aa-b085bf7e6468	2	Как проводится обучение в учебное время?	Через занятия в образовательных учреждениях	Только через телевидение	Только через инструктаж на предприятиях	A
5f99a63d-648b-42d2-85c4-14a055c1de51	88880461-dcb7-44eb-b7aa-b085bf7e6468	3	Что относится к формам обучения во внеурочное время?	Лекции, беседы, учебные фильмы и тематические мероприятия	Только экзамены	Только онлайн-курсы	A
e4bbbdd2-f7fd-4949-b9d0-c0b5789d1c2b	88880461-dcb7-44eb-b7aa-b085bf7e6468	4	Почему обучение населения зависит от времени года?	Из-за изменения уровня пожарной опасности	Только из-за расписания школ	Потому что зимой обучение запрещено	A
bab5cc99-79ec-4b2e-bd7f-40a2f963e754	88880461-dcb7-44eb-b7aa-b085bf7e6468	5	О чем акцентируется обучение в летний период?	Безопасность отопительных приборов	Требования пожарной безопасности в лесных и степных массивах	Только эвакуация из зданий	B
\.


--
-- Data for Name: lesson_quizzes; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.lesson_quizzes (id, lesson_id, source_file_name, created_at, updated_at) FROM stdin;
8db2aee3-9501-48e9-93a0-c5cceed00ac8	107b5ff3-ac0b-4be3-a2c1-4ad73983d23b	тест1.txt	2026-05-20 08:02:49.265721+00	2026-05-20 08:02:49.265721+00
d2e3058b-343b-4116-a6b9-e522e5f7450c	7612594a-faf7-4cd5-b6d0-2219e98d8ec8	тест2.txt	2026-05-20 08:03:45.073733+00	2026-05-20 08:03:45.073733+00
2ba232c8-dc09-4752-8a2f-9bb33e902059	88f0e9b7-ea8c-4d92-af5a-4b7a03ac1592	тест3.txt	2026-05-20 08:04:30.363125+00	2026-05-20 08:04:30.363125+00
7d29b1b6-c9ca-4cc3-b817-e37e2e46c7bf	23068a0a-5c34-4919-805c-c9c529901ce5	тест4.txt	2026-05-20 08:05:18.189784+00	2026-05-20 08:05:18.189784+00
5a4937f0-5bd5-475a-a0cf-d3e33df0dc79	c252c5c6-c596-427c-9654-1a2d6a35fe97	тест5.txt	2026-05-20 08:06:03.635249+00	2026-05-20 08:06:03.635249+00
25046c2b-6fed-4240-bf24-bb66a7fc8028	0defd165-3d67-4bfa-ac95-552e407afb0a	тест6.txt	2026-05-20 08:07:02.297354+00	2026-05-20 08:07:02.297354+00
8e48c271-d7cd-4a01-82ec-71480211dd46	e4ff4692-a1cb-41c5-9b12-3c5ce1846f84	тест7.txt	2026-05-20 08:07:45.837626+00	2026-05-20 08:07:45.837626+00
0abcd72c-bf6a-4f08-a909-de61520c4da1	411ea650-90b6-496e-9cbe-70f61fc7aa5a	тест8.txt	2026-05-20 08:08:30.701606+00	2026-05-20 08:08:30.701606+00
88880461-dcb7-44eb-b7aa-b085bf7e6468	a394628f-a661-44a0-a721-a6774dbea1f6	тест9.txt	2026-05-20 08:09:11.173065+00	2026-05-20 08:09:11.173065+00
d6d459f8-5143-488a-b588-d55ae9de0938	04b035e0-508e-4e0a-8079-ca7aa6b9a729	тест10.txt	2026-05-20 08:10:05.467307+00	2026-05-20 08:10:05.467307+00
\.


--
-- Data for Name: lessons; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.lessons (id, course_id, title, description, youtube_url, youtube_video_id, order_index, is_free, created_at, updated_at, video_source, video_url) FROM stdin;
107b5ff3-ac0b-4be3-a2c1-4ad73983d23b	09443e7f-eb91-4838-9622-c7c389a04960	Приказ 276 ПТМ	\N	\N	\N	1	f	2026-05-20 08:01:57.676729+00	2026-05-20 08:01:57.676729+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264117660_d02417e7354c49939e25ccb676d1d4a0.mp4
7612594a-faf7-4cd5-b6d0-2219e98d8ec8	09443e7f-eb91-4838-9622-c7c389a04960	Обучения населения	\N	\N	\N	2	f	2026-05-20 08:03:24.632911+00	2026-05-20 08:03:24.632911+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264204606_f686c9fc08684a73b8e8afa05cd2e876.mp4
88f0e9b7-ea8c-4d92-af5a-4b7a03ac1592	09443e7f-eb91-4838-9622-c7c389a04960	Нарушения при инструктаже	\N	\N	\N	3	f	2026-05-20 08:04:10.84792+00	2026-05-20 08:04:10.84792+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264250834_5775865a6ae3433d895a36be49c4b23d.mp4
23068a0a-5c34-4919-805c-c9c529901ce5	09443e7f-eb91-4838-9622-c7c389a04960	Курс ПТМ	\N	\N	\N	4	f	2026-05-20 08:04:58.681578+00	2026-05-20 08:04:58.681578+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264298655_2eeb12d7657543fd8ce575ebc5e14423.mp4
c252c5c6-c596-427c-9654-1a2d6a35fe97	09443e7f-eb91-4838-9622-c7c389a04960	Программа обучения ПТМ	\N	\N	\N	5	f	2026-05-20 08:05:48.410091+00	2026-05-20 08:05:48.410091+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264348384_04ef44a711714be2b6b2f98185417ec8.mp4
0defd165-3d67-4bfa-ac95-552e407afb0a	09443e7f-eb91-4838-9622-c7c389a04960	Пирамида выживания ПТМ	\N	\N	\N	6	f	2026-05-20 08:06:30.709575+00	2026-05-20 08:06:30.709575+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264390700_7117997fda8e4f9da8426074099d04a1.mp4
e4ff4692-a1cb-41c5-9b12-3c5ce1846f84	09443e7f-eb91-4838-9622-c7c389a04960	ПТМ на рабочем месте	\N	\N	\N	7	f	2026-05-20 08:07:31.581029+00	2026-05-20 08:07:31.581029+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264451561_61170dd5e9374f4ba3c749417a38a0f5.mp4
411ea650-90b6-496e-9cbe-70f61fc7aa5a	09443e7f-eb91-4838-9622-c7c389a04960	ПТМ с отрывом от производства	\N	\N	\N	8	f	2026-05-20 08:08:16.868834+00	2026-05-20 08:08:16.868834+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264496856_9fbc76bbf6c94623998c3d10bf26302e.mp4
a394628f-a661-44a0-a721-a6774dbea1f6	09443e7f-eb91-4838-9622-c7c389a04960	5 видов инструктажей	\N	\N	\N	9	f	2026-05-20 08:08:57.309915+00	2026-05-20 08:08:57.309915+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264537297_6e36f99feaed4fc48a9ee43937146cf9.mp4
04b035e0-508e-4e0a-8079-ca7aa6b9a729	09443e7f-eb91-4838-9622-c7c389a04960	Пожарная безопасность	\N	\N	\N	10	f	2026-05-20 08:09:50.24996+00	2026-05-20 08:09:50.24996+00	upload	http://localhost:8080/files/videos/09443e7f-eb91-4838-9622-c7c389a04960/1779264590237_b5444a3f560b4a76a4dbc5a5e1566101.mp4
\.


--
-- Data for Name: org_events; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.org_events (id, organization_id, employee_id, title, event_type, briefing_kind, starts_at, location, participants, created_by, created_at, updated_at, course_id, ends_at) FROM stdin;
aa171e8c-cae2-4d11-94e7-14a239ba6312	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	7352c1fa-5021-41e8-89ed-7d2d70a98e01	Вводный инструктаж — Дэриэль	training	introductory	2026-05-20 07:04:00+00	Переговорная / Учебный класс	\N	a209cb79-9286-42cd-b635-636ef4bd27cb	2026-05-20 06:59:46.115557+00	2026-05-20 07:02:26.251434+00	\N	\N
7ce8a195-e952-44e9-a2b0-fb896c0d3532	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	7352c1fa-5021-41e8-89ed-7d2d70a98e01	Первичный инструктаж — Дэриэль	training	primary	2026-05-21 14:00:00+00	Переговорная / Учебный класс	\N	\N	2026-05-20 07:07:23.995596+00	2026-05-20 07:07:23.995596+00	\N	\N
b5cb1acf-53cf-485e-a23f-42ee611d5749	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	7352c1fa-5021-41e8-89ed-7d2d70a98e01	Целевой инструктаж — Дэриэль	training	targeted	2026-05-20 09:00:00+00	Переговорная / Учебный класс	\N	a209cb79-9286-42cd-b635-636ef4bd27cb	2026-05-20 07:27:54.95749+00	2026-05-20 07:27:54.95749+00	\N	\N
ad1b954d-e9fc-47da-be49-565d05274c90	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	7fda9576-3492-4578-b09c-1fb7bc03bcad	Вводный инструктаж — ER	training	introductory	2026-05-20 14:00:00+00	Переговорная / Учебный класс	\N	a209cb79-9286-42cd-b635-636ef4bd27cb	2026-05-20 10:05:07.687888+00	2026-05-20 10:05:07.687888+00	\N	\N
b160c9ef-db70-469c-9448-2b2ab955fec0	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	Повторный инструктаж — Миронов Антон	training	repeat	2026-06-10 06:00:00+00	Онлайн (видео-инструктаж)	\N	a209cb79-9286-42cd-b635-636ef4bd27cb	2026-06-10 06:54:39.149653+00	2026-06-10 06:54:39.149653+00	09443e7f-eb91-4838-9622-c7c389a04960	2026-06-10 08:00:00+00
\.


--
-- Data for Name: organizations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.organizations (id, name, bin, phone, email, address, contact_person, is_active, created_at, updated_at) FROM stdin;
2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80	ТОО SkyTech	070105777999	+7 777 777 88 99	skytech@gmail.com	Муканова 34	\N	t	2026-05-19 11:05:32.350379+00	2026-05-19 11:05:32.350379+00
\.


--
-- Data for Name: profiles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.profiles (id, email, role, created_at, updated_at, full_name, is_active, phone, "position", department, bio, avatar_url, organization_id) FROM stdin;
a209cb79-9286-42cd-b635-636ef4bd27cb	sky@gmail.com	org_admin	2026-05-19 11:08:56.045152+00	2026-05-19 11:08:56.045152+00	Шамшук Адэлиева	t	\N	\N	\N	\N	\N	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80
ea88ad51-94ee-4ae8-b7cc-1c8c5af815c9	sky1@gmail.com	user	2026-05-19 11:10:59.355972+00	2026-05-19 11:10:59.355972+00	Миронов Антон	t	\N	\N	\N	\N	\N	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80
7352c1fa-5021-41e8-89ed-7d2d70a98e01	de@gmail.com	user	2026-05-20 06:59:46.110693+00	2026-05-20 06:59:46.110693+00	Дэриэль	t	\N	\N	\N	\N	\N	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80
7fda9576-3492-4578-b09c-1fb7bc03bcad	er@gmail.com	user	2026-05-20 10:05:07.683007+00	2026-05-20 10:05:07.683007+00	ER	t	\N	\N	\N	\N	\N	2ae2c4c9-b8e7-4f16-ac57-4f2ba1de2e80
809cf23d-ff95-41bc-aa6d-0c38f2d2aced	adminaq@gmail.com	super_admin	2026-02-25 09:27:32.90017+00	2026-05-21 07:17:25.507003+00	Фарида Дарья Захировна	t	+7 777 777 88 99	Администратор платформы	\N	\N	\N	\N
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.schema_migrations (version, dirty) FROM stdin;
13	f
\.


--
-- Name: ai_analysis_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.ai_analysis_id_seq', 10, true);


--
-- Name: knowledge_base_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.knowledge_base_id_seq', 5, true);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: auth; Owner: -
--

ALTER TABLE ONLY auth.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: ai_analysis ai_analysis_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_analysis
    ADD CONSTRAINT ai_analysis_pkey PRIMARY KEY (id);


--
-- Name: backend_meta backend_meta_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.backend_meta
    ADD CONSTRAINT backend_meta_pkey PRIMARY KEY (id);


--
-- Name: briefing_records briefing_records_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.briefing_records
    ADD CONSTRAINT briefing_records_pkey PRIMARY KEY (id);


--
-- Name: briefing_videos briefing_videos_course_id_briefing_kind_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.briefing_videos
    ADD CONSTRAINT briefing_videos_course_id_briefing_kind_key UNIQUE (course_id, briefing_kind);


--
-- Name: briefing_videos briefing_videos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.briefing_videos
    ADD CONSTRAINT briefing_videos_pkey PRIMARY KEY (id);


--
-- Name: contact_requests contact_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact_requests
    ADD CONSTRAINT contact_requests_pkey PRIMARY KEY (id);


--
-- Name: course_assignments course_assignments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_assignments
    ADD CONSTRAINT course_assignments_pkey PRIMARY KEY (id);


--
-- Name: course_assignments course_assignments_user_id_course_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_assignments
    ADD CONSTRAINT course_assignments_user_id_course_id_key UNIQUE (user_id, course_id);


--
-- Name: course_materials course_materials_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_materials
    ADD CONSTRAINT course_materials_pkey PRIMARY KEY (id);


--
-- Name: courses courses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT courses_pkey PRIMARY KEY (id);


--
-- Name: knowledge_base knowledge_base_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.knowledge_base
    ADD CONSTRAINT knowledge_base_pkey PRIMARY KEY (id);


--
-- Name: lesson_quiz_attempts lesson_quiz_attempts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quiz_attempts
    ADD CONSTRAINT lesson_quiz_attempts_pkey PRIMARY KEY (id);


--
-- Name: lesson_quiz_questions lesson_quiz_questions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quiz_questions
    ADD CONSTRAINT lesson_quiz_questions_pkey PRIMARY KEY (id);


--
-- Name: lesson_quiz_questions lesson_quiz_questions_quiz_id_order_index_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quiz_questions
    ADD CONSTRAINT lesson_quiz_questions_quiz_id_order_index_key UNIQUE (quiz_id, order_index);


--
-- Name: lesson_quizzes lesson_quizzes_lesson_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quizzes
    ADD CONSTRAINT lesson_quizzes_lesson_id_key UNIQUE (lesson_id);


--
-- Name: lesson_quizzes lesson_quizzes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quizzes
    ADD CONSTRAINT lesson_quizzes_pkey PRIMARY KEY (id);


--
-- Name: lessons lessons_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lessons
    ADD CONSTRAINT lessons_pkey PRIMARY KEY (id);


--
-- Name: org_events org_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.org_events
    ADD CONSTRAINT org_events_pkey PRIMARY KEY (id);


--
-- Name: organizations organizations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_pkey PRIMARY KEY (id);


--
-- Name: profiles profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: auth_users_email_idx; Type: INDEX; Schema: auth; Owner: -
--

CREATE INDEX auth_users_email_idx ON auth.users USING btree (lower(email));


--
-- Name: idx_ai_analysis_employee; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ai_analysis_employee ON public.ai_analysis USING btree (employee_id, created_at DESC);


--
-- Name: idx_ai_analysis_org; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ai_analysis_org ON public.ai_analysis USING btree (organization_id, created_at DESC);


--
-- Name: idx_assignments_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_assignments_course ON public.course_assignments USING btree (course_id);


--
-- Name: idx_assignments_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_assignments_user ON public.course_assignments USING btree (user_id);


--
-- Name: idx_briefing_records_employee; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_briefing_records_employee ON public.briefing_records USING btree (employee_id);


--
-- Name: idx_briefing_records_event; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_briefing_records_event ON public.briefing_records USING btree (event_id);


--
-- Name: idx_briefing_records_org; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_briefing_records_org ON public.briefing_records USING btree (organization_id);


--
-- Name: idx_briefing_videos_course; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_briefing_videos_course ON public.briefing_videos USING btree (course_id);


--
-- Name: idx_contact_requests_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_requests_created_at ON public.contact_requests USING btree (created_at DESC);


--
-- Name: idx_contact_requests_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_requests_status ON public.contact_requests USING btree (status);


--
-- Name: idx_course_assignments_training_completed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_course_assignments_training_completed ON public.course_assignments USING btree (user_id, course_id) WHERE (training_completed_at IS NOT NULL);


--
-- Name: idx_course_materials_course_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_course_materials_course_id ON public.course_materials USING btree (course_id);


--
-- Name: idx_lesson_quiz_attempts_user_lesson; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_lesson_quiz_attempts_user_lesson ON public.lesson_quiz_attempts USING btree (user_id, lesson_id, completed_at DESC);


--
-- Name: idx_lesson_quiz_questions_quiz; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_lesson_quiz_questions_quiz ON public.lesson_quiz_questions USING btree (quiz_id, order_index);


--
-- Name: idx_lessons_course_order; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_lessons_course_order ON public.lessons USING btree (course_id, order_index);


--
-- Name: idx_org_events_employee; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_org_events_employee ON public.org_events USING btree (employee_id);


--
-- Name: idx_org_events_org; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_org_events_org ON public.org_events USING btree (organization_id);


--
-- Name: idx_organizations_bin; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_organizations_bin ON public.organizations USING btree (bin) WHERE ((bin IS NOT NULL) AND (bin <> ''::text));


--
-- Name: idx_profiles_organization_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_profiles_organization_id ON public.profiles USING btree (organization_id) WHERE (organization_id IS NOT NULL);


--
-- Name: course_assignments trg_assignments_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_assignments_updated_at BEFORE UPDATE ON public.course_assignments FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: contact_requests trg_contact_requests_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_contact_requests_updated_at BEFORE UPDATE ON public.contact_requests FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: courses trg_courses_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_courses_updated_at BEFORE UPDATE ON public.courses FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: lesson_quizzes trg_lesson_quizzes_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_lesson_quizzes_updated_at BEFORE UPDATE ON public.lesson_quizzes FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: lessons trg_lessons_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_lessons_updated_at BEFORE UPDATE ON public.lessons FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: organizations trg_organizations_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_organizations_updated_at BEFORE UPDATE ON public.organizations FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: profiles trg_profiles_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_profiles_updated_at BEFORE UPDATE ON public.profiles FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: ai_analysis ai_analysis_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_analysis
    ADD CONSTRAINT ai_analysis_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.profiles(id) ON DELETE CASCADE;


--
-- Name: ai_analysis ai_analysis_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_analysis
    ADD CONSTRAINT ai_analysis_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organizations(id) ON DELETE CASCADE;


--
-- Name: briefing_records briefing_records_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.briefing_records
    ADD CONSTRAINT briefing_records_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.profiles(id) ON DELETE CASCADE;


--
-- Name: briefing_records briefing_records_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.briefing_records
    ADD CONSTRAINT briefing_records_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.org_events(id) ON DELETE SET NULL;


--
-- Name: briefing_records briefing_records_instructor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.briefing_records
    ADD CONSTRAINT briefing_records_instructor_id_fkey FOREIGN KEY (instructor_id) REFERENCES public.profiles(id) ON DELETE SET NULL;


--
-- Name: briefing_records briefing_records_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.briefing_records
    ADD CONSTRAINT briefing_records_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organizations(id) ON DELETE CASCADE;


--
-- Name: briefing_videos briefing_videos_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.briefing_videos
    ADD CONSTRAINT briefing_videos_course_id_fkey FOREIGN KEY (course_id) REFERENCES public.courses(id) ON DELETE CASCADE;


--
-- Name: course_assignments course_assignments_assigned_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_assignments
    ADD CONSTRAINT course_assignments_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES public.profiles(id);


--
-- Name: course_assignments course_assignments_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_assignments
    ADD CONSTRAINT course_assignments_course_id_fkey FOREIGN KEY (course_id) REFERENCES public.courses(id) ON DELETE CASCADE;


--
-- Name: course_assignments course_assignments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_assignments
    ADD CONSTRAINT course_assignments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.profiles(id) ON DELETE CASCADE;


--
-- Name: course_materials course_materials_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.course_materials
    ADD CONSTRAINT course_materials_course_id_fkey FOREIGN KEY (course_id) REFERENCES public.courses(id) ON DELETE CASCADE;


--
-- Name: courses courses_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT courses_created_by_fkey FOREIGN KEY (created_by) REFERENCES auth.users(id) ON DELETE SET NULL;


--
-- Name: lesson_quiz_attempts lesson_quiz_attempts_lesson_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quiz_attempts
    ADD CONSTRAINT lesson_quiz_attempts_lesson_id_fkey FOREIGN KEY (lesson_id) REFERENCES public.lessons(id) ON DELETE CASCADE;


--
-- Name: lesson_quiz_attempts lesson_quiz_attempts_quiz_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quiz_attempts
    ADD CONSTRAINT lesson_quiz_attempts_quiz_id_fkey FOREIGN KEY (quiz_id) REFERENCES public.lesson_quizzes(id) ON DELETE CASCADE;


--
-- Name: lesson_quiz_attempts lesson_quiz_attempts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quiz_attempts
    ADD CONSTRAINT lesson_quiz_attempts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.profiles(id) ON DELETE CASCADE;


--
-- Name: lesson_quiz_questions lesson_quiz_questions_quiz_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quiz_questions
    ADD CONSTRAINT lesson_quiz_questions_quiz_id_fkey FOREIGN KEY (quiz_id) REFERENCES public.lesson_quizzes(id) ON DELETE CASCADE;


--
-- Name: lesson_quizzes lesson_quizzes_lesson_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lesson_quizzes
    ADD CONSTRAINT lesson_quizzes_lesson_id_fkey FOREIGN KEY (lesson_id) REFERENCES public.lessons(id) ON DELETE CASCADE;


--
-- Name: lessons lessons_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lessons
    ADD CONSTRAINT lessons_course_id_fkey FOREIGN KEY (course_id) REFERENCES public.courses(id) ON DELETE CASCADE;


--
-- Name: org_events org_events_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.org_events
    ADD CONSTRAINT org_events_course_id_fkey FOREIGN KEY (course_id) REFERENCES public.courses(id) ON DELETE SET NULL;


--
-- Name: org_events org_events_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.org_events
    ADD CONSTRAINT org_events_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.profiles(id) ON DELETE SET NULL;


--
-- Name: org_events org_events_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.org_events
    ADD CONSTRAINT org_events_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.profiles(id) ON DELETE SET NULL;


--
-- Name: org_events org_events_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.org_events
    ADD CONSTRAINT org_events_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organizations(id) ON DELETE CASCADE;


--
-- Name: profiles profiles_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_id_fkey FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE;


--
-- Name: profiles profiles_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organizations(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

\unrestrict 7nhzqGsCuAwzQSQL5qp5PgvTSNX8jURsrCfw7NaqOgaOAW6jlp1odRIeozgxRuY

