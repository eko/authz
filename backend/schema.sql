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
    is_locked boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.authz_actions OWNER TO root;

--
-- Name: authz_policies; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_policies (
    id text NOT NULL,
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
-- Name: authz_principals_roles; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authz_principals_roles (
    principal_id text NOT NULL,
    role_id text NOT NULL
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
-- Name: authz_principals authz_principals_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals
    ADD CONSTRAINT authz_principals_pkey PRIMARY KEY (id);


--
-- Name: authz_principals_roles authz_principals_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_roles
    ADD CONSTRAINT authz_principals_roles_pkey PRIMARY KEY (principal_id, role_id);


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
-- Name: authz_principals_roles fk_authz_principals_roles_principal; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_roles
    ADD CONSTRAINT fk_authz_principals_roles_principal FOREIGN KEY (principal_id) REFERENCES public.authz_principals(id);


--
-- Name: authz_principals_roles fk_authz_principals_roles_role; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authz_principals_roles
    ADD CONSTRAINT fk_authz_principals_roles_role FOREIGN KEY (role_id) REFERENCES public.authz_roles(id);


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

