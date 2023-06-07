--
-- PostgreSQL database dump
--

-- Dumped from database version 13.0 (Debian 13.0-1.pgdg100+1)
-- Dumped by pg_dump version 13.0 (Debian 13.0-1.pgdg100+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: privilege; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.privilege (
    id integer NOT NULL,
    username character varying(80) NOT NULL,
    status character varying(80) DEFAULT 'BRONZE'::character varying NOT NULL,
    balance integer,
    CONSTRAINT privilege_status_check CHECK (((status)::text = ANY ((ARRAY['BRONZE'::character varying, 'SILVER'::character varying, 'GOLD'::character varying])::text[])))
);


ALTER TABLE public.privilege OWNER TO postgres;

--
-- Name: privilege_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.privilege_history (
    id integer NOT NULL,
    privilege_id integer,
    ticket_uid uuid NOT NULL,
    datetime timestamp without time zone NOT NULL,
    balance_diff integer NOT NULL,
    operation_type character varying(20) NOT NULL,
    CONSTRAINT privilege_history_operation_type_check CHECK (((operation_type)::text = ANY ((ARRAY['FILL_IN_BALANCE'::character varying, 'DEBIT_THE_ACCOUNT'::character varying])::text[])))
);


ALTER TABLE public.privilege_history OWNER TO postgres;

--
-- Name: privilege_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.privilege_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.privilege_history_id_seq OWNER TO postgres;

--
-- Name: privilege_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.privilege_history_id_seq OWNED BY public.privilege_history.id;


--
-- Name: privilege_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.privilege_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.privilege_id_seq OWNER TO postgres;

--
-- Name: privilege_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.privilege_id_seq OWNED BY public.privilege.id;


--
-- Name: privilege id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.privilege ALTER COLUMN id SET DEFAULT nextval('public.privilege_id_seq'::regclass);


--
-- Name: privilege_history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.privilege_history ALTER COLUMN id SET DEFAULT nextval('public.privilege_history_id_seq'::regclass);


--
-- Data for Name: privilege; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.privilege (id, username, status, balance) FROM stdin;
\.


--
-- Data for Name: privilege_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.privilege_history (id, privilege_id, ticket_uid, datetime, balance_diff, operation_type) FROM stdin;
\.


--
-- Name: privilege_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.privilege_history_id_seq', 1, false);


--
-- Name: privilege_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.privilege_id_seq', 1, false);


--
-- Name: privilege_history privilege_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.privilege_history
    ADD CONSTRAINT privilege_history_pkey PRIMARY KEY (id);


--
-- Name: privilege privilege_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.privilege
    ADD CONSTRAINT privilege_pkey PRIMARY KEY (id);


--
-- Name: privilege privilege_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.privilege
    ADD CONSTRAINT privilege_username_key UNIQUE (username);


--
-- Name: privilege_history privilege_history_privilege_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.privilege_history
    ADD CONSTRAINT privilege_history_privilege_id_fkey FOREIGN KEY (privilege_id) REFERENCES public.privilege(id);


--
-- PostgreSQL database dump complete
--

