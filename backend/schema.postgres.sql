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
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_actions OWNER TO root;

--
-- Name: authz_attributes; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_attributes (
    id bigint NOT NULL,
    key_name text,
    value text
);


ALTER TABLE public.authz_attributes OWNER TO root;

--
-- Name: authz_attributes_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.authz_attributes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authz_attributes_id_seq OWNER TO root;

--
-- Name: authz_attributes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.authz_attributes_id_seq OWNED BY public.authz_attributes.id;


--
-- Name: authz_clients; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_clients (
    id text NOT NULL,
    secret character varying(512),
    name text,
    domain character varying(512),
    data text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_clients OWNER TO root;

--
-- Name: authz_compiled_policies; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_compiled_policies (
    policy_id text,
    principal_id text,
    resource_kind text,
    resource_value text,
    action_id text,
    version bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_compiled_policies OWNER TO root;

--
-- Name: authz_oauth_tokens; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_oauth_tokens (
    id bigint NOT NULL,
    code character varying(512),
    access character varying(512),
    refresh character varying(512),
    data text,
    expired_at bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_oauth_tokens OWNER TO root;

--
-- Name: authz_oauth_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.authz_oauth_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authz_oauth_tokens_id_seq OWNER TO root;

--
-- Name: authz_oauth_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.authz_oauth_tokens_id_seq OWNED BY public.authz_oauth_tokens.id;


--
-- Name: authz_policies; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_policies (
    id text NOT NULL,
    attribute_rules jsonb,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_policies OWNER TO root;

--
-- Name: authz_policies_actions; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_policies_actions (
    policy_id text NOT NULL,
    action_id text NOT NULL
);


ALTER TABLE public.authz_policies_actions OWNER TO root;

--
-- Name: authz_policies_resources; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_policies_resources (
    policy_id text NOT NULL,
    resource_id text NOT NULL
);


ALTER TABLE public.authz_policies_resources OWNER TO root;

--
-- Name: authz_principals; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_principals (
    id text NOT NULL,
    is_locked boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_principals OWNER TO root;

--
-- Name: authz_principals_attributes; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_principals_attributes (
    principal_id text NOT NULL,
    attribute_id bigint NOT NULL
);


ALTER TABLE public.authz_principals_attributes OWNER TO root;

--
-- Name: authz_principals_roles; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_principals_roles (
    role_id text NOT NULL,
    principal_id text NOT NULL
);


ALTER TABLE public.authz_principals_roles OWNER TO root;

--
-- Name: authz_resources; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_resources (
    id text NOT NULL,
    kind text,
    value text,
    is_locked boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_resources OWNER TO root;

--
-- Name: authz_resources_attributes; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_resources_attributes (
    resource_id text NOT NULL,
    attribute_id bigint NOT NULL
);


ALTER TABLE public.authz_resources_attributes OWNER TO root;

--
-- Name: authz_roles; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_roles (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_roles OWNER TO root;

--
-- Name: authz_roles_policies; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_roles_policies (
    role_id text NOT NULL,
    policy_id text NOT NULL
);


ALTER TABLE public.authz_roles_policies OWNER TO root;

--
-- Name: authz_stats; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_stats (
    id text NOT NULL,
    date text,
    checks_allowed_number bigint,
    checks_denied_number bigint
);


ALTER TABLE public.authz_stats OWNER TO root;

--
-- Name: authz_users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_users (
    username text NOT NULL,
    password_hash text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_users OWNER TO root;

--
-- Name: authz_attributes id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_attributes ALTER COLUMN id SET DEFAULT nextval('public.authz_attributes_id_seq'::regclass);


--
-- Name: authz_oauth_tokens id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_oauth_tokens ALTER COLUMN id SET DEFAULT nextval('public.authz_oauth_tokens_id_seq'::regclass);


--
-- Name: authz_actions authz_actions_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_actions
    ADD CONSTRAINT authz_actions_pkey PRIMARY KEY (id);


--
-- Name: authz_attributes authz_attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_attributes
    ADD CONSTRAINT authz_attributes_pkey PRIMARY KEY (id);


--
-- Name: authz_clients authz_clients_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_clients
    ADD CONSTRAINT authz_clients_pkey PRIMARY KEY (id);


--
-- Name: authz_oauth_tokens authz_oauth_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_oauth_tokens
    ADD CONSTRAINT authz_oauth_tokens_pkey PRIMARY KEY (id);


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
-- Name: authz_principals_attributes authz_principals_attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_attributes
    ADD CONSTRAINT authz_principals_attributes_pkey PRIMARY KEY (principal_id, attribute_id);


--
-- Name: authz_principals authz_principals_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals
    ADD CONSTRAINT authz_principals_pkey PRIMARY KEY (id);


--
-- Name: authz_principals_roles authz_principals_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_roles
    ADD CONSTRAINT authz_principals_roles_pkey PRIMARY KEY (role_id, principal_id);


--
-- Name: authz_resources_attributes authz_resources_attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_resources_attributes
    ADD CONSTRAINT authz_resources_attributes_pkey PRIMARY KEY (resource_id, attribute_id);


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
-- Name: authz_stats authz_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_stats
    ADD CONSTRAINT authz_stats_pkey PRIMARY KEY (id);


--
-- Name: authz_users authz_users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_users
    ADD CONSTRAINT authz_users_pkey PRIMARY KEY (username);


--
-- Name: idx_authz_compiled_policies_action_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_authz_compiled_policies_action_id ON public.authz_compiled_policies USING btree (action_id);


--
-- Name: idx_authz_compiled_policies_policy_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_authz_compiled_policies_policy_id ON public.authz_compiled_policies USING btree (policy_id);


--
-- Name: idx_authz_compiled_policies_principal_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_authz_compiled_policies_principal_id ON public.authz_compiled_policies USING btree (principal_id);


--
-- Name: idx_authz_compiled_policies_resource_kind; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_authz_compiled_policies_resource_kind ON public.authz_compiled_policies USING btree (resource_kind);


--
-- Name: idx_authz_compiled_policies_resource_value; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_authz_compiled_policies_resource_value ON public.authz_compiled_policies USING btree (resource_value);


--
-- Name: idx_authz_compiled_policies_version; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_authz_compiled_policies_version ON public.authz_compiled_policies USING btree (version);


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
-- Name: authz_principals_attributes fk_authz_principals_attributes_attribute; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_attributes
    ADD CONSTRAINT fk_authz_principals_attributes_attribute FOREIGN KEY (attribute_id) REFERENCES public.authz_attributes(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_principals_attributes fk_authz_principals_attributes_principal; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_attributes
    ADD CONSTRAINT fk_authz_principals_attributes_principal FOREIGN KEY (principal_id) REFERENCES public.authz_principals(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_principals_roles fk_authz_principals_roles_principal; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_roles
    ADD CONSTRAINT fk_authz_principals_roles_principal FOREIGN KEY (principal_id) REFERENCES public.authz_principals(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_principals_roles fk_authz_principals_roles_role; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_roles
    ADD CONSTRAINT fk_authz_principals_roles_role FOREIGN KEY (role_id) REFERENCES public.authz_roles(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_resources_attributes fk_authz_resources_attributes_attribute; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_resources_attributes
    ADD CONSTRAINT fk_authz_resources_attributes_attribute FOREIGN KEY (attribute_id) REFERENCES public.authz_attributes(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: authz_resources_attributes fk_authz_resources_attributes_resource; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_resources_attributes
    ADD CONSTRAINT fk_authz_resources_attributes_resource FOREIGN KEY (resource_id) REFERENCES public.authz_resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


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
-- PostgreSQL database dump complete
--

