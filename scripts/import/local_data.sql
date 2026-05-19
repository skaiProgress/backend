-- Импорт данных из Supabase (data_only_dump.sql) в локальный Postgres.
-- Требует: backend/migrations уже применены (docker compose up / make migrate-up).
-- Порядок: auth.users → profiles → courses → lessons → assignments → materials.

BEGIN;

SET session_replication_role = replica;

TRUNCATE TABLE
  public.course_assignments,
  public.course_materials,
  public.lessons,
  public.courses,
  public.profiles
RESTART IDENTITY CASCADE;

DELETE FROM auth.users;

-- auth.users
COPY auth.users (instance_id, id, aud, role, email, encrypted_password, email_confirmed_at, invited_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, email_change_token_new, email_change, email_change_sent_at, last_sign_in_at, raw_app_meta_data, raw_user_meta_data, is_super_admin, created_at, updated_at, phone, phone_confirmed_at, phone_change, phone_change_token, phone_change_sent_at, email_change_token_current, email_change_confirm_status, banned_until, reauthentication_token, reauthentication_sent_at, is_sso_user, deleted_at, is_anonymous) FROM stdin;
00000000-0000-0000-0000-000000000000	e93e5ca8-02b8-4cde-97ad-997d27939b71	authenticated	authenticated	u@gmail.com	$2a$10$HQ0UFf2egsebmtdzikulROVFz5EqDy2oPMy3OBTdb7.Wmxc8Xn7xq	2026-02-25 11:02:00.60558+00	\N		\N		\N			\N	2026-02-25 11:53:32.714127+00	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-02-25 11:02:00.526305+00	2026-02-25 11:53:32.761728+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	418168da-0988-40f8-9e39-6fb9a16f7f04	authenticated	authenticated	44sadik1@gmail.com	$2a$10$Fd/X7hQDKCeoxFJdygXVDO6JMKag0vhjd.R48z5BuwY.Vl2KWj54e	2026-03-27 06:48:47.096725+00	\N		\N		\N			\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-03-27 06:48:47.047445+00	2026-03-27 06:48:47.098213+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	b2ddce86-1cc9-4116-bc75-89e046b09415	authenticated	authenticated	44sadik2@gmail.com	$2a$10$CTcj/v/81M8Lr3CVQqP0ouFRuaZYkLAMZk2MmVbGDbvQjiEpS5.Oy	2026-03-27 06:49:17.983843+00	\N		\N		\N			\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-03-27 06:49:17.98066+00	2026-03-27 06:49:17.984554+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	cf3c5924-0bb8-47a1-a868-6778753f1fdf	authenticated	authenticated	44sadik3@gmail.com	$2a$10$sZX0ogROaQMcrrIiViYJXeHFXm.wGit0orB0RqWv7JW1jElHxKvDS	2026-03-27 06:49:44.720741+00	\N		\N		\N			\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-03-27 06:49:44.718141+00	2026-03-27 06:49:44.721433+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	6621431c-4f4c-45cc-b3e7-6127ec76594e	authenticated	authenticated	antiterror@gmail.com	$2a$10$kp1Og4dV7gMA4Sl99mCckewxXuFaLciXVc.Y9VAdbfQy/1lJrHRia	2026-02-26 05:39:21.647976+00	\N		\N		\N			\N	2026-03-27 07:28:40.614461+00	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-02-26 05:39:21.62469+00	2026-03-27 09:00:27.34912+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	authenticated	authenticated	adminaq@gmail.com	$2a$06$yYa/qV4.t1okeM3blPrsIeQkxWlQjLOBuEdmyYHMSMEglL3OImSoy	2026-02-25 09:27:32.913595+00	\N		\N		\N			\N	2026-05-19 07:13:36.555998+00	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-02-25 09:27:32.900512+00	2026-05-19 07:13:36.610433+00	\N	\N			\N		0	\N		\N	f	\N	f
\.

-- public.profiles
COPY public.profiles (id, email, role, created_at, updated_at, full_name, is_active) FROM stdin;
e93e5ca8-02b8-4cde-97ad-997d27939b71	u@gmail.com	user	2026-02-25 11:02:00.521627+00	2026-02-25 11:02:00.93886+00	Серега Серый	t
809cf23d-ff95-41bc-aa6d-0c38f2d2aced	adminaq@gmail.com	super_admin	2026-02-25 09:27:32.90017+00	2026-02-25 11:39:10.18862+00	\N	t
6621431c-4f4c-45cc-b3e7-6127ec76594e	antiterror@gmail.com	user	2026-02-26 05:39:21.624334+00	2026-02-26 05:39:21.961523+00	Antiterror	t
418168da-0988-40f8-9e39-6fb9a16f7f04	44sadik1@gmail.com	user	2026-03-27 06:48:47.04712+00	2026-03-27 06:48:47.412452+00	44 Садик 1	t
b2ddce86-1cc9-4116-bc75-89e046b09415	44sadik2@gmail.com	user	2026-03-27 06:49:17.980361+00	2026-03-27 06:49:18.289806+00	44 Садик 2	t
cf3c5924-0bb8-47a1-a868-6778753f1fdf	44sadik3@gmail.com	user	2026-03-27 06:49:44.717823+00	2026-03-27 06:49:45.027196+00	44 садик 3	t
\.

-- public.courses
COPY public.courses (id, title, description, status, created_at, updated_at, cover_url, created_by) FROM stdin;
db8a9d31-e0c4-405b-9152-7af9112deaae	Пожарно-технический минимум 	Минимальный минимум знаний по пожарной безопасности	published	2026-02-25 10:19:54.573255+00	2026-02-25 10:19:54.573255+00	\N	\N
0aac07cf-5e01-41d6-9b53-d4c6b3a81777	Безопасность и охрана труда  	Безопасность и охрана труда	published	2026-02-25 11:06:43.743288+00	2026-02-25 11:06:43.743288+00	\N	\N
929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Антитеррористическая защищенность 	Согласно Закону РК № 416 и Приказу МОН № 117 \nсобственники и руководители обязаны проводить обучение персонала по\nантитеррористической защищённости на всех уязвимых объектах, включая\nшколы, больницы, торгово-развлекательные центры, театры и\nадминистративные здания \nШтраф статья 425 КоАП РК - от 100 до 1 000 МРП	published	2026-02-25 12:10:32.011907+00	2026-03-27 07:27:59.995026+00	\N	\N
\.

-- public.lessons
COPY public.lessons (id, course_id, title, description, youtube_url, youtube_video_id, order_index, is_free, created_at, updated_at) FROM stdin;
ec5b5c7d-7032-4b63-97af-28f1c2e04423	db8a9d31-e0c4-405b-9152-7af9112deaae	Вебинар 1 	\N	https://youtu.be/h5A9eQNpQcY?si=RJaCY3eyUEC6T905	h5A9eQNpQcY	1	f	2026-02-25 10:20:23.534599+00	2026-02-25 10:20:23.534599+00
fa581763-3dba-438a-bdb5-e11f6aa88274	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Антитеррористическая защищенность урок 3	\N	https://www.youtube.com/embed/8faDx2QHHOw?si=qq-yWhkLtoZSpRw5	8faDx2QHHOw	5	f	2026-02-25 12:13:53.244202+00	2026-02-25 12:15:08.534833+00
43e0aed0-86df-4541-8fee-a4c31f9c1734	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Антитеррористическая защищенность урок 1 	\N	https://www.youtube.com/embed/8Ojarw5VpPI?si=tA9tpAN34zIRRM5U	8Ojarw5VpPI	3	f	2026-02-25 12:13:11.853755+00	2026-02-25 12:15:08.536762+00
9ab6ebf3-cb8d-4499-b7a8-6b738f32c9e6	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Паспорт антитеррористической безопасности	\N	https://youtu.be/-i31Msk6DO8?si=NXZBqmcV78yieyl2	-i31Msk6DO8	6	f	2026-02-25 12:14:20.33632+00	2026-02-25 12:15:08.538192+00
9fbae39d-1fdc-433a-8cbd-149626f9ff0a	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Вебинар по Антитеррористической защищенности на казахском языке 	\N	https://youtu.be/XDbb0z678B4	XDbb0z678B4	2	f	2026-02-25 12:14:51.474903+00	2026-02-25 12:15:08.54157+00
53c152ec-4aa0-43d0-9d89-e20e27264cf7	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Вебинар по Антитеррористической защищенности 	\N	https://youtu.be/bPw__jP2JIA	bPw__jP2JIA	1	f	2026-02-25 12:11:07.223517+00	2026-02-25 12:15:08.54458+00
38ea4a97-31a7-4172-8037-55a86fc48fa3	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Антитеррористическая защищенность урок 2 	\N	https://www.youtube.com/embed/Cq_VAsudJ8I?si=P3qXHLWJKIL8bFFr	Cq_VAsudJ8I	4	f	2026-02-25 12:13:30.406536+00	2026-02-25 12:15:08.550298+00
\.

-- public.course_assignments
COPY public.course_assignments (id, user_id, course_id, assigned_by, assigned_at, expires_at, status, revoked_at, updated_at) FROM stdin;
6f93d893-4fd8-45f0-9b13-7d706b7a251f	e93e5ca8-02b8-4cde-97ad-997d27939b71	0aac07cf-5e01-41d6-9b53-d4c6b3a81777	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-02-25 11:26:17.179+00	\N	revoked	2026-02-25 11:26:25.771+00	2026-02-25 11:26:26.084992+00
1b603a79-0a3b-4437-a423-72ed0c2e5e3e	e93e5ca8-02b8-4cde-97ad-997d27939b71	db8a9d31-e0c4-405b-9152-7af9112deaae	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-02-25 11:26:37.51+00	\N	active	\N	2026-02-25 11:26:37.809774+00
6f818e4a-a5d8-4fc2-8259-8e2209351888	6621431c-4f4c-45cc-b3e7-6127ec76594e	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-02-26 05:39:38.307+00	\N	active	\N	2026-02-26 05:39:38.667114+00
\.

-- public.course_materials (пусто в дампе)
COPY public.course_materials (id, course_id, name, file_url, file_type, file_size, created_at) FROM stdin;
\.

SET session_replication_role = DEFAULT;

COMMIT;
