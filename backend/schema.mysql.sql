-- MySQL dump 10.13  Distrib 8.0.31, for Linux (x86_64)
--
-- Host: localhost    Database: root
-- ------------------------------------------------------
-- Server version	8.0.31

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `authz_actions`
--

DROP TABLE IF EXISTS `authz_actions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_actions` (
  `id` varchar(191) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_attributes`
--

DROP TABLE IF EXISTS `authz_attributes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_attributes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `key_name` longtext,
  `value` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_audit`
--

DROP TABLE IF EXISTS `authz_audit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_audit` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `date` datetime(3) DEFAULT NULL,
  `principal` longtext,
  `resource_kind` longtext,
  `resource_value` longtext,
  `action` longtext,
  `is_allowed` tinyint(1) DEFAULT NULL,
  `policy_id` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_clients`
--

DROP TABLE IF EXISTS `authz_clients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_clients` (
  `id` varchar(191) NOT NULL,
  `secret` varchar(512) DEFAULT NULL,
  `name` longtext,
  `domain` varchar(512) DEFAULT NULL,
  `data` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_compiled_policies`
--

DROP TABLE IF EXISTS `authz_compiled_policies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_compiled_policies` (
  `policy_id` varchar(191) DEFAULT NULL,
  `principal_id` varchar(191) DEFAULT NULL,
  `resource_kind` varchar(191) DEFAULT NULL,
  `resource_value` varchar(191) DEFAULT NULL,
  `action_id` varchar(191) DEFAULT NULL,
  `version` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  KEY `idx_authz_compiled_policies_resource_value` (`resource_value`),
  KEY `idx_authz_compiled_policies_action_id` (`action_id`),
  KEY `idx_authz_compiled_policies_version` (`version`),
  KEY `idx_authz_compiled_policies_policy_id` (`policy_id`),
  KEY `idx_authz_compiled_policies_principal_id` (`principal_id`),
  KEY `idx_authz_compiled_policies_resource_kind` (`resource_kind`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_oauth_tokens`
--

DROP TABLE IF EXISTS `authz_oauth_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_oauth_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(512) DEFAULT NULL,
  `access` varchar(512) DEFAULT NULL,
  `refresh` varchar(512) DEFAULT NULL,
  `data` text,
  `expired_at` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_policies`
--

DROP TABLE IF EXISTS `authz_policies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_policies` (
  `id` varchar(191) NOT NULL,
  `attribute_rules` json DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_policies_actions`
--

DROP TABLE IF EXISTS `authz_policies_actions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_policies_actions` (
  `policy_id` varchar(191) NOT NULL,
  `action_id` varchar(191) NOT NULL,
  PRIMARY KEY (`policy_id`,`action_id`),
  KEY `fk_authz_policies_actions_action` (`action_id`),
  CONSTRAINT `fk_authz_policies_actions_action` FOREIGN KEY (`action_id`) REFERENCES `authz_actions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_authz_policies_actions_policy` FOREIGN KEY (`policy_id`) REFERENCES `authz_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_policies_resources`
--

DROP TABLE IF EXISTS `authz_policies_resources`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_policies_resources` (
  `policy_id` varchar(191) NOT NULL,
  `resource_id` varchar(191) NOT NULL,
  PRIMARY KEY (`policy_id`,`resource_id`),
  KEY `fk_authz_policies_resources_resource` (`resource_id`),
  CONSTRAINT `fk_authz_policies_resources_policy` FOREIGN KEY (`policy_id`) REFERENCES `authz_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_authz_policies_resources_resource` FOREIGN KEY (`resource_id`) REFERENCES `authz_resources` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_principals`
--

DROP TABLE IF EXISTS `authz_principals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_principals` (
  `id` varchar(191) NOT NULL,
  `is_locked` tinyint(1) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_principals_attributes`
--

DROP TABLE IF EXISTS `authz_principals_attributes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_principals_attributes` (
  `principal_id` varchar(191) NOT NULL,
  `attribute_id` bigint NOT NULL,
  PRIMARY KEY (`principal_id`,`attribute_id`),
  KEY `fk_authz_principals_attributes_attribute` (`attribute_id`),
  CONSTRAINT `fk_authz_principals_attributes_attribute` FOREIGN KEY (`attribute_id`) REFERENCES `authz_attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_authz_principals_attributes_principal` FOREIGN KEY (`principal_id`) REFERENCES `authz_principals` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_principals_roles`
--

DROP TABLE IF EXISTS `authz_principals_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_principals_roles` (
  `role_id` varchar(191) NOT NULL,
  `principal_id` varchar(191) NOT NULL,
  PRIMARY KEY (`role_id`,`principal_id`),
  KEY `fk_authz_principals_roles_principal` (`principal_id`),
  CONSTRAINT `fk_authz_principals_roles_principal` FOREIGN KEY (`principal_id`) REFERENCES `authz_principals` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_authz_principals_roles_role` FOREIGN KEY (`role_id`) REFERENCES `authz_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_resources`
--

DROP TABLE IF EXISTS `authz_resources`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_resources` (
  `id` varchar(191) NOT NULL,
  `kind` longtext,
  `value` longtext,
  `is_locked` tinyint(1) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_resources_attributes`
--

DROP TABLE IF EXISTS `authz_resources_attributes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_resources_attributes` (
  `resource_id` varchar(191) NOT NULL,
  `attribute_id` bigint NOT NULL,
  PRIMARY KEY (`resource_id`,`attribute_id`),
  KEY `fk_authz_resources_attributes_attribute` (`attribute_id`),
  CONSTRAINT `fk_authz_resources_attributes_attribute` FOREIGN KEY (`attribute_id`) REFERENCES `authz_attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_authz_resources_attributes_resource` FOREIGN KEY (`resource_id`) REFERENCES `authz_resources` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_roles`
--

DROP TABLE IF EXISTS `authz_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_roles` (
  `id` varchar(191) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_roles_policies`
--

DROP TABLE IF EXISTS `authz_roles_policies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_roles_policies` (
  `role_id` varchar(191) NOT NULL,
  `policy_id` varchar(191) NOT NULL,
  PRIMARY KEY (`role_id`,`policy_id`),
  KEY `fk_authz_roles_policies_policy` (`policy_id`),
  CONSTRAINT `fk_authz_roles_policies_policy` FOREIGN KEY (`policy_id`) REFERENCES `authz_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_authz_roles_policies_role` FOREIGN KEY (`role_id`) REFERENCES `authz_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_stats`
--

DROP TABLE IF EXISTS `authz_stats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_stats` (
  `id` varchar(191) NOT NULL,
  `date` datetime(3) DEFAULT NULL,
  `checks_allowed_number` bigint DEFAULT NULL,
  `checks_denied_number` bigint DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `authz_users`
--

DROP TABLE IF EXISTS `authz_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authz_users` (
  `username` varchar(191) NOT NULL,
  `password_hash` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-01-16 19:38:26
