--
-- PostgreSQL database dump
--

\restrict FObGhe4YgCsuOAsFa24r8Tmm5dRZ9TSZvDFXMj9kKwX7DzrPqYYEcznI6nqH8vy

-- Dumped from database version 17.6
-- Dumped by pg_dump version 18.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: audit_log_entries; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.audit_log_entries (instance_id, id, payload, created_at, ip_address) FROM stdin;
\.


--
-- Data for Name: custom_oauth_providers; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.custom_oauth_providers (id, provider_type, identifier, name, client_id, client_secret, acceptable_client_ids, scopes, pkce_enabled, attribute_mapping, authorization_params, enabled, email_optional, issuer, discovery_url, skip_nonce_check, cached_discovery, discovery_cached_at, authorization_url, token_url, userinfo_url, jwks_uri, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: flow_state; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.flow_state (id, user_id, auth_code, code_challenge_method, code_challenge, provider_type, provider_access_token, provider_refresh_token, created_at, updated_at, authentication_method, auth_code_issued_at, invite_token, referrer, oauth_client_state_id, linking_target_id, email_optional) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.users (instance_id, id, aud, role, email, encrypted_password, email_confirmed_at, invited_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, email_change_token_new, email_change, email_change_sent_at, last_sign_in_at, raw_app_meta_data, raw_user_meta_data, is_super_admin, created_at, updated_at, phone, phone_confirmed_at, phone_change, phone_change_token, phone_change_sent_at, email_change_token_current, email_change_confirm_status, banned_until, reauthentication_token, reauthentication_sent_at, is_sso_user, deleted_at, is_anonymous) FROM stdin;
00000000-0000-0000-0000-000000000000	e93e5ca8-02b8-4cde-97ad-997d27939b71	authenticated	authenticated	u@gmail.com	$2a$10$HQ0UFf2egsebmtdzikulROVFz5EqDy2oPMy3OBTdb7.Wmxc8Xn7xq	2026-02-25 11:02:00.60558+00	\N		\N		\N			\N	2026-02-25 11:53:32.714127+00	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-02-25 11:02:00.526305+00	2026-02-25 11:53:32.761728+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	418168da-0988-40f8-9e39-6fb9a16f7f04	authenticated	authenticated	44sadik1@gmail.com	$2a$10$Fd/X7hQDKCeoxFJdygXVDO6JMKag0vhjd.R48z5BuwY.Vl2KWj54e	2026-03-27 06:48:47.096725+00	\N		\N		\N			\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-03-27 06:48:47.047445+00	2026-03-27 06:48:47.098213+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	b2ddce86-1cc9-4116-bc75-89e046b09415	authenticated	authenticated	44sadik2@gmail.com	$2a$10$CTcj/v/81M8Lr3CVQqP0ouFRuaZYkLAMZk2MmVbGDbvQjiEpS5.Oy	2026-03-27 06:49:17.983843+00	\N		\N		\N			\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-03-27 06:49:17.98066+00	2026-03-27 06:49:17.984554+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	cf3c5924-0bb8-47a1-a868-6778753f1fdf	authenticated	authenticated	44sadik3@gmail.com	$2a$10$sZX0ogROaQMcrrIiViYJXeHFXm.wGit0orB0RqWv7JW1jElHxKvDS	2026-03-27 06:49:44.720741+00	\N		\N		\N			\N	\N	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-03-27 06:49:44.718141+00	2026-03-27 06:49:44.721433+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	6621431c-4f4c-45cc-b3e7-6127ec76594e	authenticated	authenticated	antiterror@gmail.com	$2a$10$kp1Og4dV7gMA4Sl99mCckewxXuFaLciXVc.Y9VAdbfQy/1lJrHRia	2026-02-26 05:39:21.647976+00	\N		\N		\N			\N	2026-03-27 07:28:40.614461+00	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-02-26 05:39:21.62469+00	2026-03-27 09:00:27.34912+00	\N	\N			\N		0	\N		\N	f	\N	f
00000000-0000-0000-0000-000000000000	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	authenticated	authenticated	adminaq@gmail.com	$2a$06$yYa/qV4.t1okeM3blPrsIeQkxWlQjLOBuEdmyYHMSMEglL3OImSoy	2026-02-25 09:27:32.913595+00	\N		\N		\N			\N	2026-05-19 07:13:36.555998+00	{"provider": "email", "providers": ["email"]}	{"email_verified": true}	\N	2026-02-25 09:27:32.900512+00	2026-05-19 07:13:36.610433+00	\N	\N			\N		0	\N		\N	f	\N	f
\.


--
-- Data for Name: identities; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.identities (provider_id, user_id, identity_data, provider, last_sign_in_at, created_at, updated_at, id) FROM stdin;
809cf23d-ff95-41bc-aa6d-0c38f2d2aced	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	{"sub": "809cf23d-ff95-41bc-aa6d-0c38f2d2aced", "email": "adminaq@gmail.com", "email_verified": false, "phone_verified": false}	email	2026-02-25 09:27:32.909417+00	2026-02-25 09:27:32.909505+00	2026-02-25 09:27:32.909505+00	9d27d74f-596e-4201-853b-01c0a8e9049a
e93e5ca8-02b8-4cde-97ad-997d27939b71	e93e5ca8-02b8-4cde-97ad-997d27939b71	{"sub": "e93e5ca8-02b8-4cde-97ad-997d27939b71", "email": "u@gmail.com", "email_verified": false, "phone_verified": false}	email	2026-02-25 11:02:00.589177+00	2026-02-25 11:02:00.589235+00	2026-02-25 11:02:00.589235+00	9d40caed-8ecd-470d-a27d-ec6fbb97174b
6621431c-4f4c-45cc-b3e7-6127ec76594e	6621431c-4f4c-45cc-b3e7-6127ec76594e	{"sub": "6621431c-4f4c-45cc-b3e7-6127ec76594e", "email": "antiterror@gmail.com", "email_verified": false, "phone_verified": false}	email	2026-02-26 05:39:21.643391+00	2026-02-26 05:39:21.643473+00	2026-02-26 05:39:21.643473+00	0af7b006-9780-4830-b7bf-e128187d67e2
418168da-0988-40f8-9e39-6fb9a16f7f04	418168da-0988-40f8-9e39-6fb9a16f7f04	{"sub": "418168da-0988-40f8-9e39-6fb9a16f7f04", "email": "44sadik1@gmail.com", "email_verified": false, "phone_verified": false}	email	2026-03-27 06:48:47.085092+00	2026-03-27 06:48:47.085151+00	2026-03-27 06:48:47.085151+00	1b8bc027-8707-4370-ae04-54e023b6271f
b2ddce86-1cc9-4116-bc75-89e046b09415	b2ddce86-1cc9-4116-bc75-89e046b09415	{"sub": "b2ddce86-1cc9-4116-bc75-89e046b09415", "email": "44sadik2@gmail.com", "email_verified": false, "phone_verified": false}	email	2026-03-27 06:49:17.981976+00	2026-03-27 06:49:17.982023+00	2026-03-27 06:49:17.982023+00	e0e2b89d-f42d-462d-bef9-2490007974fc
cf3c5924-0bb8-47a1-a868-6778753f1fdf	cf3c5924-0bb8-47a1-a868-6778753f1fdf	{"sub": "cf3c5924-0bb8-47a1-a868-6778753f1fdf", "email": "44sadik3@gmail.com", "email_verified": false, "phone_verified": false}	email	2026-03-27 06:49:44.719452+00	2026-03-27 06:49:44.719497+00	2026-03-27 06:49:44.719497+00	21686bbf-4948-44ee-9413-1f08816945d5
\.


--
-- Data for Name: instances; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.instances (id, uuid, raw_base_config, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: oauth_clients; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.oauth_clients (id, client_secret_hash, registration_type, redirect_uris, grant_types, client_name, client_uri, logo_uri, created_at, updated_at, deleted_at, client_type, token_endpoint_auth_method) FROM stdin;
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.sessions (id, user_id, created_at, updated_at, factor_id, aal, not_after, refreshed_at, user_agent, ip, tag, oauth_client_id, refresh_token_hmac_key, refresh_token_counter, scopes) FROM stdin;
de43f9bd-10e7-4c8e-961d-cbfdcb09b5f8	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-03-30 06:27:01.22551+00	2026-03-30 06:27:01.22551+00	\N	aal1	\N	\N	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36	109.175.175.254	\N	\N	\N	\N	\N
867be88d-780f-4840-a78a-1f8b74700ff9	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-04-03 09:26:21.770621+00	2026-04-03 09:26:21.770621+00	\N	aal1	\N	\N	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36	109.175.175.254	\N	\N	\N	\N	\N
03c71904-5035-47fd-b66f-055fe495437a	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-03-27 09:04:36.763634+00	2026-04-11 14:45:50.343962+00	\N	aal1	\N	2026-04-11 14:45:50.343841	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36	93.190.241.212	\N	\N	\N	\N	\N
6c507989-63e2-4a9c-ad39-606de3e33381	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-05-19 07:13:36.558276+00	2026-05-19 07:13:36.558276+00	\N	aal1	\N	\N	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36	109.175.175.254	\N	\N	\N	\N	\N
\.


--
-- Data for Name: mfa_amr_claims; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.mfa_amr_claims (session_id, created_at, updated_at, authentication_method, id) FROM stdin;
03c71904-5035-47fd-b66f-055fe495437a	2026-03-27 09:04:36.79653+00	2026-03-27 09:04:36.79653+00	password	e1885452-27ff-461b-a513-e255df7696f9
de43f9bd-10e7-4c8e-961d-cbfdcb09b5f8	2026-03-30 06:27:01.31313+00	2026-03-30 06:27:01.31313+00	password	3b23c2c2-6646-4466-b3f1-1ba425373e94
867be88d-780f-4840-a78a-1f8b74700ff9	2026-04-03 09:26:21.845467+00	2026-04-03 09:26:21.845467+00	password	405d7968-c73c-4f66-8fd2-627ead63640b
6c507989-63e2-4a9c-ad39-606de3e33381	2026-05-19 07:13:36.615328+00	2026-05-19 07:13:36.615328+00	password	0001d3e4-7eac-4f68-a6c4-afefd11bfb82
\.


--
-- Data for Name: mfa_factors; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.mfa_factors (id, user_id, friendly_name, factor_type, status, created_at, updated_at, secret, phone, last_challenged_at, web_authn_credential, web_authn_aaguid, last_webauthn_challenge_data) FROM stdin;
\.


--
-- Data for Name: mfa_challenges; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.mfa_challenges (id, factor_id, created_at, verified_at, ip_address, otp_code, web_authn_session_data) FROM stdin;
\.


--
-- Data for Name: oauth_authorizations; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.oauth_authorizations (id, authorization_id, client_id, user_id, redirect_uri, scope, state, resource, code_challenge, code_challenge_method, response_type, status, authorization_code, created_at, expires_at, approved_at, nonce) FROM stdin;
\.


--
-- Data for Name: oauth_client_states; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.oauth_client_states (id, provider_type, code_verifier, created_at) FROM stdin;
\.


--
-- Data for Name: oauth_consents; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.oauth_consents (id, user_id, client_id, scopes, granted_at, revoked_at) FROM stdin;
\.


--
-- Data for Name: one_time_tokens; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.one_time_tokens (id, user_id, token_type, token_hash, relates_to, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.refresh_tokens (instance_id, id, token, user_id, revoked, created_at, updated_at, parent, session_id) FROM stdin;
00000000-0000-0000-0000-000000000000	29	fqpn4yt5bi6m	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	f	2026-03-30 06:27:01.273935+00	2026-03-30 06:27:01.273935+00	\N	de43f9bd-10e7-4c8e-961d-cbfdcb09b5f8
00000000-0000-0000-0000-000000000000	30	4tlmqatc56un	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	f	2026-04-03 09:26:21.812907+00	2026-04-03 09:26:21.812907+00	\N	867be88d-780f-4840-a78a-1f8b74700ff9
00000000-0000-0000-0000-000000000000	28	fwarhxgcayji	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	t	2026-03-27 09:04:36.785854+00	2026-04-11 14:45:50.282048+00	\N	03c71904-5035-47fd-b66f-055fe495437a
00000000-0000-0000-0000-000000000000	31	uowrdpowcgxo	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	f	2026-04-11 14:45:50.308227+00	2026-04-11 14:45:50.308227+00	fwarhxgcayji	03c71904-5035-47fd-b66f-055fe495437a
00000000-0000-0000-0000-000000000000	32	72yh64xszd3y	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	f	2026-05-19 07:13:36.590598+00	2026-05-19 07:13:36.590598+00	\N	6c507989-63e2-4a9c-ad39-606de3e33381
\.


--
-- Data for Name: sso_providers; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.sso_providers (id, resource_id, created_at, updated_at, disabled) FROM stdin;
\.


--
-- Data for Name: saml_providers; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.saml_providers (id, sso_provider_id, entity_id, metadata_xml, metadata_url, attribute_mapping, created_at, updated_at, name_id_format) FROM stdin;
\.


--
-- Data for Name: saml_relay_states; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.saml_relay_states (id, sso_provider_id, request_id, for_email, redirect_to, created_at, updated_at, flow_state_id) FROM stdin;
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.schema_migrations (version) FROM stdin;
20171026211738
20171026211808
20171026211834
20180103212743
20180108183307
20180119214651
20180125194653
00
20210710035447
20210722035447
20210730183235
20210909172000
20210927181326
20211122151130
20211124214934
20211202183645
20220114185221
20220114185340
20220224000811
20220323170000
20220429102000
20220531120530
20220614074223
20220811173540
20221003041349
20221003041400
20221011041400
20221020193600
20221021073300
20221021082433
20221027105023
20221114143122
20221114143410
20221125140132
20221208132122
20221215195500
20221215195800
20221215195900
20230116124310
20230116124412
20230131181311
20230322519590
20230402418590
20230411005111
20230508135423
20230523124323
20230818113222
20230914180801
20231027141322
20231114161723
20231117164230
20240115144230
20240214120130
20240306115329
20240314092811
20240427152123
20240612123726
20240729123726
20240802193726
20240806073726
20241009103726
20250717082212
20250731150234
20250804100000
20250901200500
20250903112500
20250904133000
20250925093508
20251007112900
20251104100000
20251111201300
20251201000000
20260115000000
20260121000000
20260219120000
20260302000000
\.


--
-- Data for Name: sso_domains; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.sso_domains (id, sso_provider_id, domain, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: webauthn_challenges; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.webauthn_challenges (id, user_id, challenge_type, session_data, created_at, expires_at) FROM stdin;
\.


--
-- Data for Name: webauthn_credentials; Type: TABLE DATA; Schema: auth; Owner: -
--

COPY auth.webauthn_credentials (id, user_id, credential_id, public_key, attestation_type, aaguid, sign_count, transports, backup_eligible, backed_up, friendly_name, created_at, updated_at, last_used_at) FROM stdin;
\.


--
-- Data for Name: courses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.courses (id, title, description, status, created_at, updated_at, cover_url, created_by) FROM stdin;
db8a9d31-e0c4-405b-9152-7af9112deaae	Пожарно-технический минимум 	Минимальный минимум знаний по пожарной безопасности	published	2026-02-25 10:19:54.573255+00	2026-02-25 10:19:54.573255+00	\N	\N
0aac07cf-5e01-41d6-9b53-d4c6b3a81777	Безопасность и охрана труда  	Безопасность и охрана труда	published	2026-02-25 11:06:43.743288+00	2026-02-25 11:06:43.743288+00	\N	\N
929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Антитеррористическая защищенность 	Согласно Закону РК № 416 и Приказу МОН № 117 \nсобственники и руководители обязаны проводить обучение персонала по\nантитеррористической защищённости на всех уязвимых объектах, включая\nшколы, больницы, торгово-развлекательные центры, театры и\nадминистративные здания \nШтраф статья 425 КоАП РК - от 100 до 1 000 МРП	published	2026-02-25 12:10:32.011907+00	2026-03-27 07:27:59.995026+00	\N	\N
\.


--
-- Data for Name: profiles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.profiles (id, email, role, created_at, updated_at, full_name, is_active) FROM stdin;
e93e5ca8-02b8-4cde-97ad-997d27939b71	u@gmail.com	user	2026-02-25 11:02:00.521627+00	2026-02-25 11:02:00.93886+00	Серега Серый	t
809cf23d-ff95-41bc-aa6d-0c38f2d2aced	adminaq@gmail.com	super_admin	2026-02-25 09:27:32.90017+00	2026-02-25 11:39:10.18862+00	\N	t
6621431c-4f4c-45cc-b3e7-6127ec76594e	antiterror@gmail.com	user	2026-02-26 05:39:21.624334+00	2026-02-26 05:39:21.961523+00	Antiterror	t
418168da-0988-40f8-9e39-6fb9a16f7f04	44sadik1@gmail.com	user	2026-03-27 06:48:47.04712+00	2026-03-27 06:48:47.412452+00	44 Садик 1	t
b2ddce86-1cc9-4116-bc75-89e046b09415	44sadik2@gmail.com	user	2026-03-27 06:49:17.980361+00	2026-03-27 06:49:18.289806+00	44 Садик 2	t
cf3c5924-0bb8-47a1-a868-6778753f1fdf	44sadik3@gmail.com	user	2026-03-27 06:49:44.717823+00	2026-03-27 06:49:45.027196+00	44 садик 3	t
\.


--
-- Data for Name: course_assignments; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.course_assignments (id, user_id, course_id, assigned_by, assigned_at, expires_at, status, revoked_at, updated_at) FROM stdin;
6f93d893-4fd8-45f0-9b13-7d706b7a251f	e93e5ca8-02b8-4cde-97ad-997d27939b71	0aac07cf-5e01-41d6-9b53-d4c6b3a81777	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-02-25 11:26:17.179+00	\N	revoked	2026-02-25 11:26:25.771+00	2026-02-25 11:26:26.084992+00
1b603a79-0a3b-4437-a423-72ed0c2e5e3e	e93e5ca8-02b8-4cde-97ad-997d27939b71	db8a9d31-e0c4-405b-9152-7af9112deaae	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-02-25 11:26:37.51+00	\N	active	\N	2026-02-25 11:26:37.809774+00
6f818e4a-a5d8-4fc2-8259-8e2209351888	6621431c-4f4c-45cc-b3e7-6127ec76594e	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	809cf23d-ff95-41bc-aa6d-0c38f2d2aced	2026-02-26 05:39:38.307+00	\N	active	\N	2026-02-26 05:39:38.667114+00
\.


--
-- Data for Name: course_materials; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.course_materials (id, course_id, name, file_url, file_type, file_size, created_at) FROM stdin;
\.


--
-- Data for Name: lessons; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.lessons (id, course_id, title, description, youtube_url, youtube_video_id, order_index, is_free, created_at, updated_at) FROM stdin;
ec5b5c7d-7032-4b63-97af-28f1c2e04423	db8a9d31-e0c4-405b-9152-7af9112deaae	Вебинар 1 	\N	https://youtu.be/h5A9eQNpQcY?si=RJaCY3eyUEC6T905	h5A9eQNpQcY	1	f	2026-02-25 10:20:23.534599+00	2026-02-25 10:20:23.534599+00
fa581763-3dba-438a-bdb5-e11f6aa88274	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Антитеррористическая защищенность урок 3	\N	https://www.youtube.com/embed/8faDx2QHHOw?si=qq-yWhkLtoZSpRw5	8faDx2QHHOw	5	f	2026-02-25 12:13:53.244202+00	2026-02-25 12:15:08.534833+00
43e0aed0-86df-4541-8fee-a4c31f9c1734	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Антитеррористическая защищенность урок 1 	\N	https://www.youtube.com/embed/8Ojarw5VpPI?si=tA9tpAN34zIRRM5U	8Ojarw5VpPI	3	f	2026-02-25 12:13:11.853755+00	2026-02-25 12:15:08.536762+00
9ab6ebf3-cb8d-4499-b7a8-6b738f32c9e6	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Паспорт антитеррористической безопасности	\N	https://youtu.be/-i31Msk6DO8?si=NXZBqmcV78yieyl2	-i31Msk6DO8	6	f	2026-02-25 12:14:20.33632+00	2026-02-25 12:15:08.538192+00
9fbae39d-1fdc-433a-8cbd-149626f9ff0a	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Вебинар по Антитеррористической защищенности на казахском языке 	\N	https://youtu.be/XDbb0z678B4	XDbb0z678B4	2	f	2026-02-25 12:14:51.474903+00	2026-02-25 12:15:08.54157+00
53c152ec-4aa0-43d0-9d89-e20e27264cf7	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Вебинар по Антитеррористической защищенности 	\N	https://youtu.be/bPw__jP2JIA	bPw__jP2JIA	1	f	2026-02-25 12:11:07.223517+00	2026-02-25 12:15:08.54458+00
38ea4a97-31a7-4172-8037-55a86fc48fa3	929afcc7-bf56-4605-a6bd-fe0869b8a1e3	Антитеррористическая защищенность урок 2 	\N	https://www.youtube.com/embed/Cq_VAsudJ8I?si=P3qXHLWJKIL8bFFr	Cq_VAsudJ8I	4	f	2026-02-25 12:13:30.406536+00	2026-02-25 12:15:08.550298+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, email, password_hash, role, name, is_active, created_at, updated_at) FROM stdin;
20d1346a-37da-4cf1-8fd7-a3aaf78f6a50	adminaq@gmail.com	\N	super_admin	Super Admin	t	2026-02-25 09:13:42.458361+00	2026-02-25 09:13:42.458361+00
0d0c45c6-4706-4bfd-bafb-95c10a80fde3	admin@aiqadam.kz	\N	admin	Администратор	t	2026-02-25 09:13:42.458361+00	2026-02-25 09:13:42.458361+00
e1716efc-1267-40bc-9a60-0a722c1e3aac	user@aiqadam.kz	\N	user	Сотрудник	t	2026-02-25 09:13:42.458361+00	2026-02-25 09:13:42.458361+00
7d07a963-23c4-47d5-b910-71879d264f6e	ivan.petrov@company.kz	\N	user	Иван Петров	t	2026-02-25 09:13:42.458361+00	2026-02-25 09:13:42.458361+00
ea87f5aa-c89f-4a49-88ef-0c239888434b	maria.smi@company.kz	\N	user	Мария Смирнова	f	2026-02-25 09:13:42.458361+00	2026-02-25 09:13:42.458361+00
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: realtime; Owner: -
--

COPY realtime.schema_migrations (version, inserted_at) FROM stdin;
20211116024918	2026-02-25 06:16:21
20211116045059	2026-02-25 06:16:21
20211116050929	2026-02-25 06:16:21
20211116051442	2026-02-25 06:16:21
20211116212300	2026-02-25 06:16:21
20211116213355	2026-02-25 06:16:21
20211116213934	2026-02-25 06:16:21
20211116214523	2026-02-25 06:16:21
20211122062447	2026-02-25 06:16:21
20211124070109	2026-02-25 06:16:21
20211202204204	2026-02-25 06:16:21
20211202204605	2026-02-25 06:16:21
20211210212804	2026-02-25 06:16:21
20211228014915	2026-02-25 06:16:21
20220107221237	2026-02-25 08:06:33
20220228202821	2026-02-25 08:06:33
20220312004840	2026-02-25 08:06:33
20220603231003	2026-02-25 08:06:33
20220603232444	2026-02-25 08:06:33
20220615214548	2026-02-25 08:06:33
20220712093339	2026-02-25 08:06:33
20220908172859	2026-02-25 08:06:33
20220916233421	2026-02-25 08:06:33
20230119133233	2026-02-25 08:06:33
20230128025114	2026-02-25 08:06:33
20230128025212	2026-02-25 08:06:33
20230227211149	2026-02-25 08:06:33
20230228184745	2026-02-25 08:06:33
20230308225145	2026-02-25 08:06:33
20230328144023	2026-02-25 08:06:33
20231018144023	2026-02-25 08:06:33
20231204144023	2026-02-25 08:06:33
20231204144024	2026-02-25 08:06:33
20231204144025	2026-02-25 08:06:33
20240108234812	2026-02-25 08:06:33
20240109165339	2026-02-25 08:06:33
20240227174441	2026-02-25 08:06:33
20240311171622	2026-02-25 08:06:33
20240321100241	2026-02-25 08:06:33
20240401105812	2026-02-25 08:06:33
20240418121054	2026-02-25 08:06:33
20240523004032	2026-02-25 08:06:33
20240618124746	2026-02-25 08:06:33
20240801235015	2026-02-25 08:06:33
20240805133720	2026-02-25 08:06:33
20240827160934	2026-02-25 08:06:33
20240919163303	2026-02-25 08:06:33
20240919163305	2026-02-25 08:06:33
20241019105805	2026-02-25 08:06:33
20241030150047	2026-02-25 08:06:33
20241108114728	2026-02-25 08:06:33
20241121104152	2026-02-25 08:06:33
20241130184212	2026-02-25 08:06:33
20241220035512	2026-02-25 08:06:33
20241220123912	2026-02-25 08:06:33
20241224161212	2026-02-25 08:06:33
20250107150512	2026-02-25 08:06:33
20250110162412	2026-02-25 08:06:33
20250123174212	2026-02-25 08:06:33
20250128220012	2026-02-25 08:06:33
20250506224012	2026-02-25 08:06:33
20250523164012	2026-02-25 08:06:33
20250714121412	2026-02-25 08:06:34
20250905041441	2026-02-25 08:06:34
20251103001201	2026-02-25 08:06:34
20251120212548	2026-02-25 08:06:34
20251120215549	2026-02-25 08:06:34
20260218120000	2026-03-02 09:00:48
20260326120000	2026-05-19 07:13:13
\.


--
-- Data for Name: subscription; Type: TABLE DATA; Schema: realtime; Owner: -
--

COPY realtime.subscription (id, subscription_id, entity, filters, claims, created_at, action_filter) FROM stdin;
\.


--
-- Data for Name: buckets; Type: TABLE DATA; Schema: storage; Owner: -
--

COPY storage.buckets (id, name, owner, created_at, updated_at, public, avif_autodetection, file_size_limit, allowed_mime_types, owner_id, type) FROM stdin;
Material_bucket	Material_bucket	\N	2026-03-27 06:56:48.079337+00	2026-03-27 06:56:48.079337+00	t	f	\N	\N	\N	STANDARD
\.


--
-- Data for Name: buckets_analytics; Type: TABLE DATA; Schema: storage; Owner: -
--

COPY storage.buckets_analytics (name, type, format, created_at, updated_at, id, deleted_at) FROM stdin;
\.


--
-- Data for Name: buckets_vectors; Type: TABLE DATA; Schema: storage; Owner: -
--

COPY storage.buckets_vectors (id, type, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: migrations; Type: TABLE DATA; Schema: storage; Owner: -
--

COPY storage.migrations (id, name, hash, executed_at) FROM stdin;
0	create-migrations-table	e18db593bcde2aca2a408c4d1100f6abba2195df	2026-02-25 06:17:29.379087
1	initialmigration	6ab16121fbaa08bbd11b712d05f358f9b555d777	2026-02-25 06:17:29.415874
2	storage-schema	f6a1fa2c93cbcd16d4e487b362e45fca157a8dbd	2026-02-25 06:17:29.421474
3	pathtoken-column	2cb1b0004b817b29d5b0a971af16bafeede4b70d	2026-02-25 06:17:29.449609
4	add-migrations-rls	427c5b63fe1c5937495d9c635c263ee7a5905058	2026-02-25 06:17:29.462435
5	add-size-functions	79e081a1455b63666c1294a440f8ad4b1e6a7f84	2026-02-25 06:17:29.466736
6	change-column-name-in-get-size	ded78e2f1b5d7e616117897e6443a925965b30d2	2026-02-25 06:17:29.472198
7	add-rls-to-buckets	e7e7f86adbc51049f341dfe8d30256c1abca17aa	2026-02-25 06:17:29.477231
8	add-public-to-buckets	fd670db39ed65f9d08b01db09d6202503ca2bab3	2026-02-25 06:17:29.481429
9	fix-search-function	af597a1b590c70519b464a4ab3be54490712796b	2026-02-25 06:17:29.485911
10	search-files-search-function	b595f05e92f7e91211af1bbfe9c6a13bb3391e16	2026-02-25 06:17:29.491753
11	add-trigger-to-auto-update-updated_at-column	7425bdb14366d1739fa8a18c83100636d74dcaa2	2026-02-25 06:17:29.496502
12	add-automatic-avif-detection-flag	8e92e1266eb29518b6a4c5313ab8f29dd0d08df9	2026-02-25 06:17:29.501055
13	add-bucket-custom-limits	cce962054138135cd9a8c4bcd531598684b25e7d	2026-02-25 06:17:29.50517
14	use-bytes-for-max-size	941c41b346f9802b411f06f30e972ad4744dad27	2026-02-25 06:17:29.510114
15	add-can-insert-object-function	934146bc38ead475f4ef4b555c524ee5d66799e5	2026-02-25 06:17:29.535517
16	add-version	76debf38d3fd07dcfc747ca49096457d95b1221b	2026-02-25 06:17:29.539902
17	drop-owner-foreign-key	f1cbb288f1b7a4c1eb8c38504b80ae2a0153d101	2026-02-25 06:17:29.544185
18	add_owner_id_column_deprecate_owner	e7a511b379110b08e2f214be852c35414749fe66	2026-02-25 06:17:29.54837
19	alter-default-value-objects-id	02e5e22a78626187e00d173dc45f58fa66a4f043	2026-02-25 06:17:29.555484
20	list-objects-with-delimiter	cd694ae708e51ba82bf012bba00caf4f3b6393b7	2026-02-25 06:17:29.560047
21	s3-multipart-uploads	8c804d4a566c40cd1e4cc5b3725a664a9303657f	2026-02-25 06:17:29.565933
22	s3-multipart-uploads-big-ints	9737dc258d2397953c9953d9b86920b8be0cdb73	2026-02-25 06:17:29.578786
23	optimize-search-function	9d7e604cddc4b56a5422dc68c9313f4a1b6f132c	2026-02-25 06:17:29.5893
24	operation-function	8312e37c2bf9e76bbe841aa5fda889206d2bf8aa	2026-02-25 06:17:29.593794
25	custom-metadata	d974c6057c3db1c1f847afa0e291e6165693b990	2026-02-25 06:17:29.59791
26	objects-prefixes	215cabcb7f78121892a5a2037a09fedf9a1ae322	2026-02-25 06:17:29.602351
27	search-v2	859ba38092ac96eb3964d83bf53ccc0b141663a6	2026-02-25 06:17:29.606107
28	object-bucket-name-sorting	c73a2b5b5d4041e39705814fd3a1b95502d38ce4	2026-02-25 06:17:29.611277
29	create-prefixes	ad2c1207f76703d11a9f9007f821620017a66c21	2026-02-25 06:17:29.615256
30	update-object-levels	2be814ff05c8252fdfdc7cfb4b7f5c7e17f0bed6	2026-02-25 06:17:29.619095
31	objects-level-index	b40367c14c3440ec75f19bbce2d71e914ddd3da0	2026-02-25 06:17:29.623967
32	backward-compatible-index-on-objects	e0c37182b0f7aee3efd823298fb3c76f1042c0f7	2026-02-25 06:17:29.628254
33	backward-compatible-index-on-prefixes	b480e99ed951e0900f033ec4eb34b5bdcb4e3d49	2026-02-25 06:17:29.631929
34	optimize-search-function-v1	ca80a3dc7bfef894df17108785ce29a7fc8ee456	2026-02-25 06:17:29.635768
35	add-insert-trigger-prefixes	458fe0ffd07ec53f5e3ce9df51bfdf4861929ccc	2026-02-25 06:17:29.639475
36	optimise-existing-functions	6ae5fca6af5c55abe95369cd4f93985d1814ca8f	2026-02-25 06:17:29.643196
37	add-bucket-name-length-trigger	3944135b4e3e8b22d6d4cbb568fe3b0b51df15c1	2026-02-25 06:17:29.64704
38	iceberg-catalog-flag-on-buckets	02716b81ceec9705aed84aa1501657095b32e5c5	2026-02-25 06:17:29.651858
39	add-search-v2-sort-support	6706c5f2928846abee18461279799ad12b279b78	2026-02-25 06:17:29.662251
40	fix-prefix-race-conditions-optimized	7ad69982ae2d372b21f48fc4829ae9752c518f6b	2026-02-25 06:17:29.665998
41	add-object-level-update-trigger	07fcf1a22165849b7a029deed059ffcde08d1ae0	2026-02-25 06:17:29.669722
42	rollback-prefix-triggers	771479077764adc09e2ea2043eb627503c034cd4	2026-02-25 06:17:29.673999
43	fix-object-level	84b35d6caca9d937478ad8a797491f38b8c2979f	2026-02-25 06:17:29.67774
44	vector-bucket-type	99c20c0ffd52bb1ff1f32fb992f3b351e3ef8fb3	2026-02-25 06:17:29.681513
45	vector-buckets	049e27196d77a7cb76497a85afae669d8b230953	2026-02-25 06:17:29.688169
46	buckets-objects-grants	fedeb96d60fefd8e02ab3ded9fbde05632f84aed	2026-02-25 06:17:29.700136
47	iceberg-table-metadata	649df56855c24d8b36dd4cc1aeb8251aa9ad42c2	2026-02-25 06:17:29.704856
48	iceberg-catalog-ids	e0e8b460c609b9999ccd0df9ad14294613eed939	2026-02-25 06:17:29.708833
49	buckets-objects-grants-postgres	072b1195d0d5a2f888af6b2302a1938dd94b8b3d	2026-02-25 06:17:29.723482
50	search-v2-optimised	6323ac4f850aa14e7387eb32102869578b5bd478	2026-02-25 06:17:29.728231
51	index-backward-compatible-search	2ee395d433f76e38bcd3856debaf6e0e5b674011	2026-02-25 06:17:30.28228
52	drop-not-used-indexes-and-functions	5cc44c8696749ac11dd0dc37f2a3802075f3a171	2026-02-25 06:17:30.284072
53	drop-index-lower-name	d0cb18777d9e2a98ebe0bc5cc7a42e57ebe41854	2026-02-25 06:17:30.293898
54	drop-index-object-level	6289e048b1472da17c31a7eba1ded625a6457e67	2026-02-25 06:17:30.296639
55	prevent-direct-deletes	262a4798d5e0f2e7c8970232e03ce8be695d5819	2026-02-25 06:17:30.29865
57	s3-multipart-uploads-metadata	f127886e00d1b374fadbc7c6b31e09336aad5287	2026-04-07 04:54:11.98252
58	operation-ergonomics	00ca5d483b3fe0d522133d9002ccc5df98365120	2026-04-07 04:54:12.014743
56	fix-optimized-search-function	b823ed1e418101032fa01374edc9a436e54e3ed4	2026-02-25 06:17:30.304327
59	drop-unused-functions	38456f13e39691c2bbb4b5151d0d1cdbabd4a8c4	2026-05-19 07:13:10.033022
60	optimize-existing-functions-again	db35e1c91a9201e59f4fef8d972c2f277d68b157	2026-05-19 07:13:10.051061
\.


--
-- Data for Name: objects; Type: TABLE DATA; Schema: storage; Owner: -
--

COPY storage.objects (id, bucket_id, name, owner, created_at, updated_at, last_accessed_at, metadata, version, owner_id, user_metadata) FROM stdin;
\.


--
-- Data for Name: s3_multipart_uploads; Type: TABLE DATA; Schema: storage; Owner: -
--

COPY storage.s3_multipart_uploads (id, in_progress_size, upload_signature, bucket_id, key, version, owner_id, created_at, user_metadata, metadata) FROM stdin;
\.


--
-- Data for Name: s3_multipart_uploads_parts; Type: TABLE DATA; Schema: storage; Owner: -
--

COPY storage.s3_multipart_uploads_parts (id, upload_id, size, part_number, bucket_id, key, etag, owner_id, version, created_at) FROM stdin;
\.


--
-- Data for Name: vector_indexes; Type: TABLE DATA; Schema: storage; Owner: -
--

COPY storage.vector_indexes (id, name, bucket_id, data_type, dimension, distance_metric, metadata_configuration, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: secrets; Type: TABLE DATA; Schema: vault; Owner: -
--

COPY vault.secrets (id, name, description, secret, key_id, nonce, created_at, updated_at) FROM stdin;
\.


--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: auth; Owner: -
--

SELECT pg_catalog.setval('auth.refresh_tokens_id_seq', 32, true);


--
-- Name: subscription_id_seq; Type: SEQUENCE SET; Schema: realtime; Owner: -
--

SELECT pg_catalog.setval('realtime.subscription_id_seq', 1, false);


--
-- PostgreSQL database dump complete
--

\unrestrict FObGhe4YgCsuOAsFa24r8Tmm5dRZ9TSZvDFXMj9kKwX7DzrPqYYEcznI6nqH8vy

