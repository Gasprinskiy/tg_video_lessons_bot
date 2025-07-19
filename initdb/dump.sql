--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9
-- Dumped by pg_dump version 16.9

--
-- Name: bot_users_profile; Type: TABLE; Schema: public; Owner: admin_user
--

CREATE TABLE public.bot_users_profile (
    u_id bigint NOT NULL,
    tg_id bigint NOT NULL,
    tg_user_name character varying(100) NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    birth_date date NOT NULL,
    phone_number character varying(150) NOT NULL,
    active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.bot_users_profile OWNER TO admin_user;

--
-- Name: bot_users_profile_u_id_seq; Type: SEQUENCE; Schema: public; Owner: admin_user
--

ALTER TABLE public.bot_users_profile ALTER COLUMN u_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.bot_users_profile_u_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);



--
-- Name: bot_users_profile_u_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin_user
--

SELECT pg_catalog.setval('public.bot_users_profile_u_id_seq', 5, true);


--
-- Name: bot_users_profile bot_users_profile_pkey; Type: CONSTRAINT; Schema: public; Owner: admin_user
--

ALTER TABLE ONLY public.bot_users_profile
    ADD CONSTRAINT bot_users_profile_pkey PRIMARY KEY (u_id);


--
-- Name: tg_id_index; Type: INDEX; Schema: public; Owner: admin_user
--

CREATE INDEX tg_id_index ON public.bot_users_profile USING btree (tg_id) INCLUDE (tg_id) WITH (deduplicate_items='true');


--
-- PostgreSQL database dump complete
--

