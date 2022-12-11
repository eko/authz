--
-- PostgreSQL database dump
--

-- Dumped from database version 15.1 (Debian 15.1-1.pgdg110+1)
-- Dumped by pg_dump version 15.1 (Debian 15.1-1.pgdg110+1)

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
-- Name: authz_actions; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_actions (
    id bigint NOT NULL,
    is_locked boolean,
    name text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_actions OWNER TO root;

--
-- Name: authz_actions_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.authz_actions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authz_actions_id_seq OWNER TO root;

--
-- Name: authz_actions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.authz_actions_id_seq OWNED BY public.authz_actions.id;


--
-- Name: authz_policies; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_policies (
    id bigint NOT NULL,
    name text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_policies OWNER TO root;

--
-- Name: authz_policies_actions; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_policies_actions (
    policy_id bigint NOT NULL,
    action_id bigint NOT NULL
);


ALTER TABLE public.authz_policies_actions OWNER TO root;

--
-- Name: authz_policies_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.authz_policies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authz_policies_id_seq OWNER TO root;

--
-- Name: authz_policies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.authz_policies_id_seq OWNED BY public.authz_policies.id;


--
-- Name: authz_policies_resources; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_policies_resources (
    policy_id bigint NOT NULL,
    resource_id bigint NOT NULL
);


ALTER TABLE public.authz_policies_resources OWNER TO root;

--
-- Name: authz_resources; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_resources (
    id bigint NOT NULL,
    is_locked boolean,
    kind text,
    value text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_resources OWNER TO root;

--
-- Name: authz_resources_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.authz_resources_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authz_resources_id_seq OWNER TO root;

--
-- Name: authz_resources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.authz_resources_id_seq OWNED BY public.authz_resources.id;


--
-- Name: authz_roles; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_roles (
    id bigint NOT NULL,
    name text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_roles OWNER TO root;

--
-- Name: authz_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.authz_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authz_roles_id_seq OWNER TO root;

--
-- Name: authz_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.authz_roles_id_seq OWNED BY public.authz_roles.id;


--
-- Name: authz_roles_policies; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_roles_policies (
    role_id bigint NOT NULL,
    policy_id bigint NOT NULL
);


ALTER TABLE public.authz_roles_policies OWNER TO root;

--
-- Name: authz_subjects; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_subjects (
    id bigint NOT NULL,
    is_locked boolean,
    value text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_subjects OWNER TO root;

--
-- Name: authz_subjects_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.authz_subjects_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authz_subjects_id_seq OWNER TO root;

--
-- Name: authz_subjects_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.authz_subjects_id_seq OWNED BY public.authz_subjects.id;


--
-- Name: authz_subjects_roles; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_subjects_roles (
    subject_id bigint NOT NULL,
    role_id bigint NOT NULL
);


ALTER TABLE public.authz_subjects_roles OWNER TO root;

--
-- Name: authz_actions id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_actions ALTER COLUMN id SET DEFAULT nextval('public.authz_actions_id_seq'::regclass);


--
-- Name: authz_policies id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_policies ALTER COLUMN id SET DEFAULT nextval('public.authz_policies_id_seq'::regclass);


--
-- Name: authz_resources id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_resources ALTER COLUMN id SET DEFAULT nextval('public.authz_resources_id_seq'::regclass);


--
-- Name: authz_roles id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_roles ALTER COLUMN id SET DEFAULT nextval('public.authz_roles_id_seq'::regclass);


--
-- Name: authz_subjects id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_subjects ALTER COLUMN id SET DEFAULT nextval('public.authz_subjects_id_seq'::regclass);


--
-- Name: authz_actions authz_actions_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_actions
    ADD CONSTRAINT authz_actions_pkey PRIMARY KEY (id);


--
-- Name: authz_policies_actions authz_policies_actions_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_policies_actions
    ADD CONSTRAINT authz_policies_actions_pkey PRIMARY KEY (policy_id, action_id);


--
-- Name: authz_policies authz_policies_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_policies
    ADD CONSTRAINT authz_policies_pkey PRIMARY KEY (id);


--
-- Name: authz_policies_resources authz_policies_resources_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_policies_resources
    ADD CONSTRAINT authz_policies_resources_pkey PRIMARY KEY (policy_id, resource_id);


--
-- Name: authz_resources authz_resources_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_resources
    ADD CONSTRAINT authz_resources_pkey PRIMARY KEY (id);


--
-- Name: authz_roles authz_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_roles
    ADD CONSTRAINT authz_roles_pkey PRIMARY KEY (id);


--
-- Name: authz_roles_policies authz_roles_policies_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_roles_policies
    ADD CONSTRAINT authz_roles_policies_pkey PRIMARY KEY (role_id, policy_id);


--
-- Name: authz_subjects authz_subjects_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_subjects
    ADD CONSTRAINT authz_subjects_pkey PRIMARY KEY (id);


--
-- Name: authz_subjects_roles authz_subjects_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_subjects_roles
    ADD CONSTRAINT authz_subjects_roles_pkey PRIMARY KEY (subject_id, role_id);


--
-- Name: idx_authz_actions_name; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_authz_actions_name ON public.authz_actions USING btree (name);


--
-- Name: idx_authz_policies_name; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_authz_policies_name ON public.authz_policies USING btree (name);


--
-- Name: idx_authz_roles_name; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_authz_roles_name ON public.authz_roles USING btree (name);


--
-- Name: idx_authz_subjects_value; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_authz_subjects_value ON public.authz_subjects USING btree (value);


--
-- Name: authz_policies_actions fk_authz_policies_actions_action; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_policies_actions
    ADD CONSTRAINT fk_authz_policies_actions_action FOREIGN KEY (action_id) REFERENCES public.authz_actions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_policies_actions fk_authz_policies_actions_policy; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_policies_actions
    ADD CONSTRAINT fk_authz_policies_actions_policy FOREIGN KEY (policy_id) REFERENCES public.authz_policies(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_policies_resources fk_authz_policies_resources_policy; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_policies_resources
    ADD CONSTRAINT fk_authz_policies_resources_policy FOREIGN KEY (policy_id) REFERENCES public.authz_policies(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_policies_resources fk_authz_policies_resources_resource; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_policies_resources
    ADD CONSTRAINT fk_authz_policies_resources_resource FOREIGN KEY (resource_id) REFERENCES public.authz_resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_roles_policies fk_authz_roles_policies_policy; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_roles_policies
    ADD CONSTRAINT fk_authz_roles_policies_policy FOREIGN KEY (policy_id) REFERENCES public.authz_policies(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_roles_policies fk_authz_roles_policies_role; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_roles_policies
    ADD CONSTRAINT fk_authz_roles_policies_role FOREIGN KEY (role_id) REFERENCES public.authz_roles(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_subjects_roles fk_authz_subjects_roles_role; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_subjects_roles
    ADD CONSTRAINT fk_authz_subjects_roles_role FOREIGN KEY (role_id) REFERENCES public.authz_roles(id);


--
-- Name: authz_subjects_roles fk_authz_subjects_roles_subject; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_subjects_roles
    ADD CONSTRAINT fk_authz_subjects_roles_subject FOREIGN KEY (subject_id) REFERENCES public.authz_subjects(id);


--
-- PostgreSQL database dump complete
--

