--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2
-- Dumped by pg_dump version 16.2

-- Started on 2024-07-14 11:13:30

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
-- TOC entry 218 (class 1259 OID 16509)
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    transaction_id integer NOT NULL,
    wallet_id integer,
    amount numeric(10,2),
    transaction_type character varying(20),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    wallet_id_source integer
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16508)
-- Name: transactions_transaction_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.transactions_transaction_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transactions_transaction_id_seq OWNER TO postgres;

--
-- TOC entry 4856 (class 0 OID 0)
-- Dependencies: 217
-- Name: transactions_transaction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.transactions_transaction_id_seq OWNED BY public.transactions.transaction_id;


--
-- TOC entry 216 (class 1259 OID 16499)
-- Name: wallets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wallets (
    wallet_id integer NOT NULL,
    user_id integer,
    balance numeric(10,2) DEFAULT 0.00,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.wallets OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 16498)
-- Name: wallets_wallet_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.wallets_wallet_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.wallets_wallet_id_seq OWNER TO postgres;

--
-- TOC entry 4857 (class 0 OID 0)
-- Dependencies: 215
-- Name: wallets_wallet_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.wallets_wallet_id_seq OWNED BY public.wallets.wallet_id;


--
-- TOC entry 4697 (class 2604 OID 16512)
-- Name: transactions transaction_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions ALTER COLUMN transaction_id SET DEFAULT nextval('public.transactions_transaction_id_seq'::regclass);


--
-- TOC entry 4693 (class 2604 OID 16502)
-- Name: wallets wallet_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wallets ALTER COLUMN wallet_id SET DEFAULT nextval('public.wallets_wallet_id_seq'::regclass);


--
-- TOC entry 4850 (class 0 OID 16509)
-- Dependencies: 218
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (transaction_id, wallet_id, amount, transaction_type, created_at, wallet_id_source) FROM stdin;
27	13	10000.00	in	2024-07-14 00:07:23.635399	0
28	15	10000.00	in	2024-07-14 00:07:32.508514	0
29	14	10000.00	in	2024-07-14 00:07:40.911285	0
30	13	100.00	out	2024-07-14 00:08:00.606123	15
31	15	100.00	in	2024-07-14 00:08:00.606565	13
\.


--
-- TOC entry 4848 (class 0 OID 16499)
-- Dependencies: 216
-- Data for Name: wallets; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wallets (wallet_id, user_id, balance, created_at, updated_at) FROM stdin;
14	8	10000.00	2024-07-13 17:05:52.668762	2024-07-14 00:07:40.909038
13	6	9900.00	2024-07-13 17:05:32.292502	2024-07-14 00:08:00.603432
15	9	10100.00	2024-07-13 17:07:05.263643	2024-07-14 00:08:00.605384
\.


--
-- TOC entry 4858 (class 0 OID 0)
-- Dependencies: 217
-- Name: transactions_transaction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.transactions_transaction_id_seq', 31, true);


--
-- TOC entry 4859 (class 0 OID 0)
-- Dependencies: 215
-- Name: wallets_wallet_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.wallets_wallet_id_seq', 15, true);


--
-- TOC entry 4702 (class 2606 OID 16515)
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (transaction_id);


--
-- TOC entry 4700 (class 2606 OID 16507)
-- Name: wallets wallets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wallets
    ADD CONSTRAINT wallets_pkey PRIMARY KEY (wallet_id);


--
-- TOC entry 4703 (class 2606 OID 16516)
-- Name: transactions transactions_wallet_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES public.wallets(wallet_id);


-- Completed on 2024-07-14 11:13:30

--
-- PostgreSQL database dump complete
--

