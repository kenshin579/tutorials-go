-- MySQL dump 10.13  Distrib 8.0.33, for macos13.3 (arm64)
--
-- Host: localhost    Database: keycloak-test
-- ------------------------------------------------------
-- Server version	8.0.25

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
-- Table structure for table `ADMIN_EVENT_ENTITY`
--

DROP TABLE IF EXISTS `ADMIN_EVENT_ENTITY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ADMIN_EVENT_ENTITY` (
  `ID` varchar(36) NOT NULL,
  `ADMIN_EVENT_TIME` bigint DEFAULT NULL,
  `REALM_ID` varchar(255) DEFAULT NULL,
  `OPERATION_TYPE` varchar(255) DEFAULT NULL,
  `AUTH_REALM_ID` varchar(255) DEFAULT NULL,
  `AUTH_CLIENT_ID` varchar(255) DEFAULT NULL,
  `AUTH_USER_ID` varchar(255) DEFAULT NULL,
  `IP_ADDRESS` varchar(255) DEFAULT NULL,
  `RESOURCE_PATH` text,
  `REPRESENTATION` text,
  `ERROR` varchar(255) DEFAULT NULL,
  `RESOURCE_TYPE` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_ADMIN_EVENT_TIME` (`REALM_ID`,`ADMIN_EVENT_TIME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ADMIN_EVENT_ENTITY`
--

LOCK TABLES `ADMIN_EVENT_ENTITY` WRITE;
/*!40000 ALTER TABLE `ADMIN_EVENT_ENTITY` DISABLE KEYS */;
/*!40000 ALTER TABLE `ADMIN_EVENT_ENTITY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ASSOCIATED_POLICY`
--

DROP TABLE IF EXISTS `ASSOCIATED_POLICY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ASSOCIATED_POLICY` (
  `POLICY_ID` varchar(36) NOT NULL,
  `ASSOCIATED_POLICY_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`POLICY_ID`,`ASSOCIATED_POLICY_ID`),
  KEY `IDX_ASSOC_POL_ASSOC_POL_ID` (`ASSOCIATED_POLICY_ID`),
  CONSTRAINT `FK_FRSR5S213XCX4WNKOG82SSRFY` FOREIGN KEY (`ASSOCIATED_POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`),
  CONSTRAINT `FK_FRSRPAS14XCX4WNKOG82SSRFY` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ASSOCIATED_POLICY`
--

LOCK TABLES `ASSOCIATED_POLICY` WRITE;
/*!40000 ALTER TABLE `ASSOCIATED_POLICY` DISABLE KEYS */;
/*!40000 ALTER TABLE `ASSOCIATED_POLICY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `AUTHENTICATION_EXECUTION`
--

DROP TABLE IF EXISTS `AUTHENTICATION_EXECUTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AUTHENTICATION_EXECUTION` (
  `ID` varchar(36) NOT NULL,
  `ALIAS` varchar(255) DEFAULT NULL,
  `AUTHENTICATOR` varchar(36) DEFAULT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  `FLOW_ID` varchar(36) DEFAULT NULL,
  `REQUIREMENT` int DEFAULT NULL,
  `PRIORITY` int DEFAULT NULL,
  `AUTHENTICATOR_FLOW` tinyint NOT NULL DEFAULT '0',
  `AUTH_FLOW_ID` varchar(36) DEFAULT NULL,
  `AUTH_CONFIG` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_AUTH_EXEC_REALM_FLOW` (`REALM_ID`,`FLOW_ID`),
  KEY `IDX_AUTH_EXEC_FLOW` (`FLOW_ID`),
  CONSTRAINT `FK_AUTH_EXEC_FLOW` FOREIGN KEY (`FLOW_ID`) REFERENCES `AUTHENTICATION_FLOW` (`ID`),
  CONSTRAINT `FK_AUTH_EXEC_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `AUTHENTICATION_EXECUTION`
--

LOCK TABLES `AUTHENTICATION_EXECUTION` WRITE;
/*!40000 ALTER TABLE `AUTHENTICATION_EXECUTION` DISABLE KEYS */;
INSERT INTO `AUTHENTICATION_EXECUTION` (`ID`, `ALIAS`, `AUTHENTICATOR`, `REALM_ID`, `FLOW_ID`, `REQUIREMENT`, `PRIORITY`, `AUTHENTICATOR_FLOW`, `AUTH_FLOW_ID`, `AUTH_CONFIG`) VALUES ('008b1a1b-3eae-4c31-9acb-359b6523f18a',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','c4fe3a48-6c0f-4c9e-b28d-ffaf9b9ddfdc',2,20,1,'3484027b-59da-4989-93dc-46251490ec58',NULL),('05d74038-41a7-4372-8396-0a48b989b191',NULL,'auth-otp-form','4327ba47-4116-44ea-9c4d-02907dca81e7','abb9c598-7f17-410f-a338-82f1f2e165fa',0,20,0,NULL,NULL),('06a22dec-8525-4ff9-8410-34b4fe01a9bd',NULL,'client-x509','dcc080c5-aede-4fd3-8b01-bd0928b730a2','f84a1f13-ef43-4a16-b9d6-63fe787b86dd',2,40,0,NULL,NULL),('0a8f92c0-d923-4b3c-8d23-1e59b0a85c4a',NULL,'direct-grant-validate-username','40ae881c-f4e4-4b07-b097-a67d2bf515e6','73f81fba-ce8f-4131-8637-3c64b1323d86',0,10,0,NULL,NULL),('0ab42766-5c19-4f91-a486-bf244d5dd2f4',NULL,'client-jwt','dcc080c5-aede-4fd3-8b01-bd0928b730a2','f84a1f13-ef43-4a16-b9d6-63fe787b86dd',2,20,0,NULL,NULL),('0b5a8001-9afa-4ee9-b18c-a43f40614f35',NULL,'conditional-user-configured','40ae881c-f4e4-4b07-b097-a67d2bf515e6','a86ff042-5592-47dd-bf01-7865d992e796',0,10,0,NULL,NULL),('0d6a67c5-e394-44c6-be2b-13ed86d92b65',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','b6c69e0d-d61b-45ca-acfc-1dc79b81bb62',1,20,1,'1dc10e9b-26a9-4420-aed4-cef7ce1cc0fc',NULL),('0e24e361-e617-4cd7-9ffa-0381c0c54ccf',NULL,'idp-username-password-form','dcc080c5-aede-4fd3-8b01-bd0928b730a2','0cfc44d0-dcec-4963-92be-00797e95107b',0,10,0,NULL,NULL),('0fa266df-dbad-44ba-b098-b27a93fc5f04',NULL,'idp-username-password-form','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5fca5bf2-7017-4f4a-865d-1c3db899423c',0,10,0,NULL,NULL),('108cf4c1-063b-4e16-b4b2-f35bab98cfd1',NULL,'idp-email-verification','4327ba47-4116-44ea-9c4d-02907dca81e7','34156d67-bf62-439e-aca9-5da894adc723',2,10,0,NULL,NULL),('116ca944-643d-4fd4-898d-eb21f250858e',NULL,'auth-otp-form','40ae881c-f4e4-4b07-b097-a67d2bf515e6','3f87f84a-0a1e-4919-a380-10eabd214d58',0,20,0,NULL,NULL),('15aedda6-8576-4a61-b538-beec311f8492',NULL,'registration-page-form','dcc080c5-aede-4fd3-8b01-bd0928b730a2','f77f953b-5e25-4601-b047-748724d795b3',0,10,1,'be564502-8c4b-4eb6-a9e7-2927a798abf0',NULL),('1687e228-3deb-4223-9db0-be484d11907d',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','6bfae19d-d2c6-4dd5-836c-6042247853cc',2,30,1,'b6c69e0d-d61b-45ca-acfc-1dc79b81bb62',NULL),('17c0ca2a-cc3f-4906-b5e8-cc1fd521eede',NULL,'idp-review-profile','40ae881c-f4e4-4b07-b097-a67d2bf515e6','409007a9-4804-45b2-a701-79b9f94cea54',0,10,0,NULL,'d33fb1bf-6714-46fd-a0c5-0ad1c91aa4f2'),('182359ee-550a-4a1f-a599-09cecf48162d',NULL,'conditional-user-configured','40ae881c-f4e4-4b07-b097-a67d2bf515e6','15d69b0b-32fe-4c04-ba61-c6f2c92d1ae1',0,10,0,NULL,NULL),('1cf6c75f-af09-459e-9e5b-7f9e28b6581f',NULL,'registration-user-creation','dcc080c5-aede-4fd3-8b01-bd0928b730a2','be564502-8c4b-4eb6-a9e7-2927a798abf0',0,20,0,NULL,NULL),('1d7f6895-9c24-4e3a-bd28-57c44d26ad83',NULL,'direct-grant-validate-otp','dcc080c5-aede-4fd3-8b01-bd0928b730a2','7bd6316d-df7b-43c9-b72f-af3f62def6c1',0,20,0,NULL,NULL),('1e126a03-7fa4-4d69-b58d-3d1ae2eeb56c',NULL,'conditional-user-configured','4327ba47-4116-44ea-9c4d-02907dca81e7','6153dc60-22f4-4624-a6a2-a5d0a8ee1ad0',0,10,0,NULL,NULL),('1facb9bb-7ba2-49f4-923e-0e3f84f4d41f',NULL,'direct-grant-validate-password','4327ba47-4116-44ea-9c4d-02907dca81e7','3af9837f-1ac4-4fe2-a3a7-579df7d70cbd',0,20,0,NULL,NULL),('207f4e63-f5a8-4f6a-8d82-f8a2f351b443',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','e4be6c41-61ae-4230-a920-d66ab72fd11a',0,20,1,'c4fe3a48-6c0f-4c9e-b28d-ffaf9b9ddfdc',NULL),('212b1f7d-d9e1-41f6-86c3-ee3f10e68f12',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','80314021-add0-43af-8e10-dbf41e5fe483',2,20,1,'9589b999-1269-439f-93bc-eebce417f7ad',NULL),('22e69780-fc79-4b15-ac2c-486d1f6ed7a0',NULL,'client-jwt','4327ba47-4116-44ea-9c4d-02907dca81e7','f378b286-fd54-4b8f-8615-df0118dacd24',2,20,0,NULL,NULL),('25151c39-ba9f-455c-91bc-2159cc653ed8',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','4a21f571-4b69-4ef8-92af-09fe24026b64',2,30,1,'ac2d291b-2581-4830-9cc6-c123412d3072',NULL),('26174320-9649-41c5-9d98-6b25cac8a7fb',NULL,'conditional-user-configured','dcc080c5-aede-4fd3-8b01-bd0928b730a2','62e95aab-0f29-402c-9d4e-1c062c966af7',0,10,0,NULL,NULL),('26195f4b-866f-4028-a3cb-7c9820a1e671',NULL,'client-jwt','40ae881c-f4e4-4b07-b097-a67d2bf515e6','a9a6d8d1-afbc-447a-8685-cd855fb3e627',2,20,0,NULL,NULL),('266e8c58-d2a4-4ad9-a38d-ae889f0e1765',NULL,'auth-username-password-form','dcc080c5-aede-4fd3-8b01-bd0928b730a2','b6c69e0d-d61b-45ca-acfc-1dc79b81bb62',0,10,0,NULL,NULL),('2875bbda-5a14-471f-9163-a7e81dc74103',NULL,'conditional-user-configured','4327ba47-4116-44ea-9c4d-02907dca81e7','9c471d7d-c49c-4912-a921-7e8fbfa8c845',0,10,0,NULL,NULL),('2907b4d8-08da-40eb-827e-96d269210bc0',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','ac2d291b-2581-4830-9cc6-c123412d3072',1,20,1,'3f87f84a-0a1e-4919-a380-10eabd214d58',NULL),('2eb08d06-fa66-4dbe-b313-0293bbd9351e',NULL,'idp-review-profile','40ae881c-f4e4-4b07-b097-a67d2bf515e6','e4be6c41-61ae-4230-a920-d66ab72fd11a',0,10,0,NULL,'fa8ff9b0-26bb-442c-bfdb-505699cf6443'),('2f725fc9-bc47-4680-b04f-4de867a624e3',NULL,'auth-otp-form','40ae881c-f4e4-4b07-b097-a67d2bf515e6','a86ff042-5592-47dd-bf01-7865d992e796',0,20,0,NULL,NULL),('2f839bf6-a600-4884-96f1-72919f854f0d',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','0cfc44d0-dcec-4963-92be-00797e95107b',1,20,1,'9e16df50-ccc9-4bd6-842d-7349838e3168',NULL),('316a33f4-9a29-4d72-83db-2fcbd2bc8fa7',NULL,'client-secret-jwt','dcc080c5-aede-4fd3-8b01-bd0928b730a2','f84a1f13-ef43-4a16-b9d6-63fe787b86dd',2,30,0,NULL,NULL),('32d2f255-0719-47b5-bc5a-b9632e705981',NULL,'http-basic-authenticator','4327ba47-4116-44ea-9c4d-02907dca81e7','330e968f-77f6-4862-9ba6-8b37217c3255',0,10,0,NULL,NULL),('33f2e6ae-5caf-4149-b816-b8546a36ef29',NULL,'reset-otp','dcc080c5-aede-4fd3-8b01-bd0928b730a2','62e95aab-0f29-402c-9d4e-1c062c966af7',0,20,0,NULL,NULL),('34b8a8d0-c504-4f24-b507-5d775584452b',NULL,'registration-terms-and-conditions','4327ba47-4116-44ea-9c4d-02907dca81e7','76074602-02dc-43ec-a571-34efecba9953',3,70,0,NULL,NULL),('372166a6-78da-48e0-8b8f-29526161cc8d',NULL,'direct-grant-validate-otp','4327ba47-4116-44ea-9c4d-02907dca81e7','9c471d7d-c49c-4912-a921-7e8fbfa8c845',0,20,0,NULL,NULL),('3a7a298b-99a2-449b-974d-671503ea4284',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','73f81fba-ce8f-4131-8637-3c64b1323d86',1,30,1,'4e206e20-7402-45d2-ab7b-87d6d4b138e0',NULL),('3ce4260a-a7ef-476c-80d8-e93a17fcf8c0',NULL,'registration-page-form','40ae881c-f4e4-4b07-b097-a67d2bf515e6','08b0543c-51a6-4200-bb72-38404a916954',0,10,1,'e9314680-3395-49c7-a314-82fcff06249d',NULL),('3f2d0cbd-e59a-4744-9d07-3f370dae5feb',NULL,'reset-credential-email','40ae881c-f4e4-4b07-b097-a67d2bf515e6','d5b7af02-aba4-4b1c-bcca-e068b83dfee1',0,20,0,NULL,NULL),('4155641a-5a3e-4825-b790-cb0eeebbf368',NULL,'identity-provider-redirector','dcc080c5-aede-4fd3-8b01-bd0928b730a2','6bfae19d-d2c6-4dd5-836c-6042247853cc',2,25,0,NULL,NULL),('4421af11-0dce-4cae-83ef-613a3b40831b',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','616a6379-b1b2-4e14-96ea-db9d45b16538',1,40,1,'6153dc60-22f4-4624-a6a2-a5d0a8ee1ad0',NULL),('47c859bb-4139-474d-aea9-70fb85e29340',NULL,'auth-spnego','40ae881c-f4e4-4b07-b097-a67d2bf515e6','4a21f571-4b69-4ef8-92af-09fe24026b64',3,20,0,NULL,NULL),('4a971568-ec5d-40e3-9f43-7c19642c33b0',NULL,'direct-grant-validate-username','4327ba47-4116-44ea-9c4d-02907dca81e7','3af9837f-1ac4-4fe2-a3a7-579df7d70cbd',0,10,0,NULL,NULL),('4aa94ee5-4347-4c54-a296-208d001b9627',NULL,'idp-username-password-form','4327ba47-4116-44ea-9c4d-02907dca81e7','2ef963b2-037d-4185-8c73-503fc6e878da',0,10,0,NULL,NULL),('4ae38037-02dd-4aec-b0ab-6fc8b108fac1',NULL,'registration-password-action','dcc080c5-aede-4fd3-8b01-bd0928b730a2','be564502-8c4b-4eb6-a9e7-2927a798abf0',0,50,0,NULL,NULL),('4b7e0899-11b9-4777-bdbf-6a5b6ec4baa3',NULL,'registration-recaptcha-action','4327ba47-4116-44ea-9c4d-02907dca81e7','76074602-02dc-43ec-a571-34efecba9953',3,60,0,NULL,NULL),('4d02dd27-d115-4311-8323-abd8a4c02db6',NULL,'identity-provider-redirector','4327ba47-4116-44ea-9c4d-02907dca81e7','7bf8b92b-a6be-4924-a010-8750f64b92e9',2,25,0,NULL,NULL),('4d1963e9-5253-4709-8387-165b9c8bb4d0',NULL,'idp-email-verification','40ae881c-f4e4-4b07-b097-a67d2bf515e6','b241fbb5-0c09-4d3d-9ae0-9245a055ff6f',2,10,0,NULL,NULL),('4e3eefa4-5683-4e3d-8c53-af6a60f97b36',NULL,'client-secret','4327ba47-4116-44ea-9c4d-02907dca81e7','f378b286-fd54-4b8f-8615-df0118dacd24',2,10,0,NULL,NULL),('4e97871f-d816-4b09-9d5e-ca449705ca10',NULL,'idp-email-verification','dcc080c5-aede-4fd3-8b01-bd0928b730a2','c73fcc66-40d1-4495-bc79-5f10c9a1b722',2,10,0,NULL,NULL),('4e9e949f-3854-4ded-bae1-1196c92cfc52',NULL,'auth-otp-form','dcc080c5-aede-4fd3-8b01-bd0928b730a2','1dc10e9b-26a9-4420-aed4-cef7ce1cc0fc',0,20,0,NULL,NULL),('52deb635-94f2-4e6a-9d56-76286dcf755e',NULL,'docker-http-basic-authenticator','4327ba47-4116-44ea-9c4d-02907dca81e7','692a7ac4-e4aa-47ac-8bb4-5267c7930f6d',0,10,0,NULL,NULL),('54645524-3b7a-4226-8086-e03aad58e7cd',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','3484027b-59da-4989-93dc-46251490ec58',0,20,1,'80314021-add0-43af-8e10-dbf41e5fe483',NULL),('57868366-1c1c-4ccc-bd94-d4747912256b',NULL,'auth-cookie','40ae881c-f4e4-4b07-b097-a67d2bf515e6','4a21f571-4b69-4ef8-92af-09fe24026b64',2,10,0,NULL,NULL),('5b491eab-b013-46f7-bb9d-e293146db0b5',NULL,'direct-grant-validate-password','dcc080c5-aede-4fd3-8b01-bd0928b730a2','b639104c-892c-4843-a2ed-ca4324be5788',0,20,0,NULL,NULL),('5ca55432-c97a-463d-a462-572839c0bc5b',NULL,'conditional-user-configured','dcc080c5-aede-4fd3-8b01-bd0928b730a2','1dc10e9b-26a9-4420-aed4-cef7ce1cc0fc',0,10,0,NULL,NULL),('5eabdef6-5996-40ba-a6b9-1e736f52c0d3',NULL,'reset-credentials-choose-user','4327ba47-4116-44ea-9c4d-02907dca81e7','616a6379-b1b2-4e14-96ea-db9d45b16538',0,10,0,NULL,NULL),('61f94cf7-7eb7-4115-a06d-5888bf742899',NULL,'conditional-user-configured','40ae881c-f4e4-4b07-b097-a67d2bf515e6','3f87f84a-0a1e-4919-a380-10eabd214d58',0,10,0,NULL,NULL),('6332306d-fe8c-434e-a64e-d2168f8d3fef',NULL,'idp-confirm-link','4327ba47-4116-44ea-9c4d-02907dca81e7','2d986fa8-a731-473e-ad1d-a3bcfe9f6724',0,10,0,NULL,NULL),('688f7d8e-5ec5-4f74-8cfb-e9a44afa55fb',NULL,'reset-credential-email','dcc080c5-aede-4fd3-8b01-bd0928b730a2','06d58040-a014-4f13-a9ea-72835ad3c576',0,20,0,NULL,NULL),('6a0a50e3-0a84-4e11-b8af-8b4524005eef',NULL,'idp-create-user-if-unique','4327ba47-4116-44ea-9c4d-02907dca81e7','77637b0f-15f0-4e0f-9e54-4294535f01c8',2,10,0,NULL,'bee1c985-9a19-48b5-9e0a-19b813f9a415'),('6c96ee85-1907-4a1b-8c24-33a5d38aba66',NULL,'docker-http-basic-authenticator','dcc080c5-aede-4fd3-8b01-bd0928b730a2','1250dad6-ebc6-4a7e-a847-a83a67967a51',0,10,0,NULL,NULL),('6d71cb0f-4645-46f7-b707-84719f120806',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','b639104c-892c-4843-a2ed-ca4324be5788',1,30,1,'7bd6316d-df7b-43c9-b72f-af3f62def6c1',NULL),('71c94ee0-ba63-48be-bc91-74acacab657c',NULL,'reset-password','40ae881c-f4e4-4b07-b097-a67d2bf515e6','d5b7af02-aba4-4b1c-bcca-e068b83dfee1',0,30,0,NULL,NULL),('78108f01-ed9e-4114-8666-fb7a293bf787',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','65417fbf-411b-4301-8a33-234ae9d9ee5e',0,20,1,'b241fbb5-0c09-4d3d-9ae0-9245a055ff6f',NULL),('7b375cbc-0e47-4415-8e9c-ff10633de338',NULL,'direct-grant-validate-password','40ae881c-f4e4-4b07-b097-a67d2bf515e6','73f81fba-ce8f-4131-8637-3c64b1323d86',0,20,0,NULL,NULL),('7bfa79e2-c832-4f90-bb81-5f728fc017da',NULL,'auth-cookie','4327ba47-4116-44ea-9c4d-02907dca81e7','7bf8b92b-a6be-4924-a010-8750f64b92e9',2,10,0,NULL,NULL),('7c27152f-881d-47c9-ba41-d305c88b4cf9',NULL,'auth-username-password-form','40ae881c-f4e4-4b07-b097-a67d2bf515e6','ac2d291b-2581-4830-9cc6-c123412d3072',0,10,0,NULL,NULL),('7ce623e2-7cca-43fc-b8a3-4fbf471410ce',NULL,'direct-grant-validate-username','dcc080c5-aede-4fd3-8b01-bd0928b730a2','b639104c-892c-4843-a2ed-ca4324be5788',0,10,0,NULL,NULL),('7ec9fff0-4ac6-4865-a7ce-05ba3fd09f11',NULL,'idp-email-verification','40ae881c-f4e4-4b07-b097-a67d2bf515e6','80314021-add0-43af-8e10-dbf41e5fe483',2,10,0,NULL,NULL),('808b2a91-9c0d-4076-9581-ab00d52765e6',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','77637b0f-15f0-4e0f-9e54-4294535f01c8',2,20,1,'2d986fa8-a731-473e-ad1d-a3bcfe9f6724',NULL),('82665914-644e-4d67-b382-ac6f6e539231',NULL,'client-secret','dcc080c5-aede-4fd3-8b01-bd0928b730a2','f84a1f13-ef43-4a16-b9d6-63fe787b86dd',2,10,0,NULL,NULL),('8630225d-f60b-482d-8049-6508920ea66d',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','06d58040-a014-4f13-a9ea-72835ad3c576',1,40,1,'62e95aab-0f29-402c-9d4e-1c062c966af7',NULL),('8ae1957f-c838-4588-861a-466587da4ae0',NULL,'idp-confirm-link','40ae881c-f4e4-4b07-b097-a67d2bf515e6','65417fbf-411b-4301-8a33-234ae9d9ee5e',0,10,0,NULL,NULL),('8b715299-e093-4fc7-9a1c-cee0662eaea6',NULL,'conditional-user-configured','4327ba47-4116-44ea-9c4d-02907dca81e7','223dc23f-861c-4e10-89b1-a5fede09a82a',0,10,0,NULL,NULL),('8bea830b-5438-45b0-b982-79d165d08c44',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','c73fcc66-40d1-4495-bc79-5f10c9a1b722',2,20,1,'0cfc44d0-dcec-4963-92be-00797e95107b',NULL),('8ca771da-f738-4735-a7c4-b0266fed0308',NULL,'reset-otp','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1a26f0eb-5055-4ece-a233-22ae767e5e7a',0,20,0,NULL,NULL),('8d90a7e8-d4b5-4c82-8993-a7651a8cf844',NULL,'auth-otp-form','dcc080c5-aede-4fd3-8b01-bd0928b730a2','9e16df50-ccc9-4bd6-842d-7349838e3168',0,20,0,NULL,NULL),('8ef29272-77f6-4091-87a8-b1c64da944b8',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','3af9837f-1ac4-4fe2-a3a7-579df7d70cbd',1,30,1,'9c471d7d-c49c-4912-a921-7e8fbfa8c845',NULL),('9491cf10-e102-4aab-9dfd-0c447a3318a5',NULL,'reset-credential-email','4327ba47-4116-44ea-9c4d-02907dca81e7','616a6379-b1b2-4e14-96ea-db9d45b16538',0,20,0,NULL,NULL),('95db3c95-4d2b-4ed7-bb0a-35b632102dd3',NULL,'idp-review-profile','4327ba47-4116-44ea-9c4d-02907dca81e7','48fb5d06-357a-463e-9263-1de5c87f092c',0,10,0,NULL,'24523cd1-73bc-4810-932d-6438b849af6a'),('967b075d-f210-4711-96cf-4fca84f5c54f',NULL,'idp-review-profile','dcc080c5-aede-4fd3-8b01-bd0928b730a2','f97a7702-55e0-4f57-b4a1-ae4b09239fa8',0,10,0,NULL,'0b74b1d6-2788-464b-8b28-59e468d0b488'),('96ab61e9-4049-4007-9e6b-ce016d5085d0',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','d5b7af02-aba4-4b1c-bcca-e068b83dfee1',1,40,1,'1a26f0eb-5055-4ece-a233-22ae767e5e7a',NULL),('98121572-1d88-490a-bd7c-3d6667406808',NULL,'http-basic-authenticator','dcc080c5-aede-4fd3-8b01-bd0928b730a2','543dfaac-77e4-499e-aab0-851fd15d7db7',0,10,0,NULL,NULL),('98b46119-5cff-45a6-af38-57060ddfc390',NULL,'client-secret-jwt','40ae881c-f4e4-4b07-b097-a67d2bf515e6','a9a6d8d1-afbc-447a-8685-cd855fb3e627',2,30,0,NULL,NULL),('99040ffb-bb29-4fef-93f8-43bfff067ab2',NULL,'client-x509','4327ba47-4116-44ea-9c4d-02907dca81e7','f378b286-fd54-4b8f-8615-df0118dacd24',2,40,0,NULL,NULL),('9ab63d89-1793-4c41-a971-276c763d7433',NULL,'idp-create-user-if-unique','40ae881c-f4e4-4b07-b097-a67d2bf515e6','c4fe3a48-6c0f-4c9e-b28d-ffaf9b9ddfdc',2,10,0,NULL,'ba96f2c7-fdca-4963-a250-2901eaa68ac8'),('9cdb1d1d-0c28-4e15-ab18-b8510a257149',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','2ef963b2-037d-4185-8c73-503fc6e878da',1,20,1,'223dc23f-861c-4e10-89b1-a5fede09a82a',NULL),('a07986f8-8d66-419f-a9de-f50ec49c799e',NULL,'docker-http-basic-authenticator','40ae881c-f4e4-4b07-b097-a67d2bf515e6','ea3b8b38-9a74-403e-8058-be535ccea217',0,10,0,NULL,NULL),('a2b729a0-6c68-4583-b17a-b1327b733e79',NULL,'identity-provider-redirector','40ae881c-f4e4-4b07-b097-a67d2bf515e6','4a21f571-4b69-4ef8-92af-09fe24026b64',2,25,0,NULL,NULL),('a2cb18c7-f024-4c81-85aa-54c5121fd8bb',NULL,'idp-confirm-link','40ae881c-f4e4-4b07-b097-a67d2bf515e6','3484027b-59da-4989-93dc-46251490ec58',0,10,0,NULL,NULL),('a4b9b3cd-dbd6-472d-814f-e5426dafd2d0',NULL,'registration-terms-and-conditions','dcc080c5-aede-4fd3-8b01-bd0928b730a2','be564502-8c4b-4eb6-a9e7-2927a798abf0',3,70,0,NULL,NULL),('a7cc6b58-8578-44c3-b9b4-053a93539164',NULL,'registration-password-action','4327ba47-4116-44ea-9c4d-02907dca81e7','76074602-02dc-43ec-a571-34efecba9953',0,50,0,NULL,NULL),('abe1e3df-9606-4d14-b653-f4c8466b0445',NULL,'http-basic-authenticator','40ae881c-f4e4-4b07-b097-a67d2bf515e6','962742cd-4fa7-4654-a84f-10cda8492259',0,10,0,NULL,NULL),('abf09b50-6fd6-4e31-acaa-847c78d2d38f',NULL,'auth-otp-form','40ae881c-f4e4-4b07-b097-a67d2bf515e6','15d69b0b-32fe-4c04-ba61-c6f2c92d1ae1',0,20,0,NULL,NULL),('aceaebf8-1a7b-4564-b0d6-19bee0530e60',NULL,'registration-recaptcha-action','dcc080c5-aede-4fd3-8b01-bd0928b730a2','be564502-8c4b-4eb6-a9e7-2927a798abf0',3,60,0,NULL,NULL),('ad68ca07-80d4-4578-8f6f-2674e935e7bc',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','9589b999-1269-439f-93bc-eebce417f7ad',1,20,1,'15d69b0b-32fe-4c04-ba61-c6f2c92d1ae1',NULL),('afdda164-75e9-4744-91d0-6f798d35a15e',NULL,'registration-user-creation','40ae881c-f4e4-4b07-b097-a67d2bf515e6','e9314680-3395-49c7-a314-82fcff06249d',0,20,0,NULL,NULL),('b44f032e-210f-4e53-ace1-2bf8c5edbd89',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','409007a9-4804-45b2-a701-79b9f94cea54',0,20,1,'74c902dc-6c87-4332-b633-f5fee94e6286',NULL),('b5f78b49-cf43-42cf-bcfa-0b2969f69982',NULL,'registration-password-action','40ae881c-f4e4-4b07-b097-a67d2bf515e6','e9314680-3395-49c7-a314-82fcff06249d',0,50,0,NULL,NULL),('bae7316a-e609-4c0e-9c93-e90904ab7a6f',NULL,'registration-user-creation','4327ba47-4116-44ea-9c4d-02907dca81e7','76074602-02dc-43ec-a571-34efecba9953',0,20,0,NULL,NULL),('bc0ada30-7e39-4374-a781-3864b3bfd01d',NULL,'auth-otp-form','4327ba47-4116-44ea-9c4d-02907dca81e7','223dc23f-861c-4e10-89b1-a5fede09a82a',0,20,0,NULL,NULL),('be98b3f3-2d91-4cc9-9fbd-37fe1edede43',NULL,'idp-create-user-if-unique','dcc080c5-aede-4fd3-8b01-bd0928b730a2','96f3f037-c48e-451a-bc28-4bcd33f6ae3b',2,10,0,NULL,'1f88cc4f-af88-488e-8304-3febe6c52f97'),('c03f26a1-85f2-40dd-9859-0d643cc33857',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','7bf8b92b-a6be-4924-a010-8750f64b92e9',2,30,1,'7dc4a4fa-cefd-4f10-b196-f48d39cfa0db',NULL),('c19c2aad-8ca2-4162-a174-f0f394c9e6d6',NULL,'conditional-user-configured','dcc080c5-aede-4fd3-8b01-bd0928b730a2','9e16df50-ccc9-4bd6-842d-7349838e3168',0,10,0,NULL,NULL),('c3925de5-3007-46ad-b64f-a7ee49eb6c09',NULL,'direct-grant-validate-otp','40ae881c-f4e4-4b07-b097-a67d2bf515e6','4e206e20-7402-45d2-ab7b-87d6d4b138e0',0,20,0,NULL,NULL),('c46c9194-c6dc-433c-ab6a-9d7eae6d67fd',NULL,'auth-spnego','4327ba47-4116-44ea-9c4d-02907dca81e7','7bf8b92b-a6be-4924-a010-8750f64b92e9',3,20,0,NULL,NULL),('c8615722-86c9-491b-ad1f-71a9b43ace41',NULL,'auth-spnego','dcc080c5-aede-4fd3-8b01-bd0928b730a2','6bfae19d-d2c6-4dd5-836c-6042247853cc',3,20,0,NULL,NULL),('ca342a0a-1ef1-445d-bba0-83c23b8cc4d1',NULL,'client-x509','40ae881c-f4e4-4b07-b097-a67d2bf515e6','a9a6d8d1-afbc-447a-8685-cd855fb3e627',2,40,0,NULL,NULL),('cc8dfd6f-dad5-4918-9e8f-072655503d96',NULL,'idp-username-password-form','40ae881c-f4e4-4b07-b097-a67d2bf515e6','9589b999-1269-439f-93bc-eebce417f7ad',0,10,0,NULL,NULL),('cc9a6f58-4a89-4dc1-8d77-fec66c117625',NULL,'conditional-user-configured','40ae881c-f4e4-4b07-b097-a67d2bf515e6','4e206e20-7402-45d2-ab7b-87d6d4b138e0',0,10,0,NULL,NULL),('cd5fc4b2-5fb5-41fe-b1c0-69066c364717',NULL,'registration-terms-and-conditions','40ae881c-f4e4-4b07-b097-a67d2bf515e6','e9314680-3395-49c7-a314-82fcff06249d',3,70,0,NULL,NULL),('ce4017c0-5db6-4855-9304-32d4ed3379ee',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','5fca5bf2-7017-4f4a-865d-1c3db899423c',1,20,1,'a86ff042-5592-47dd-bf01-7865d992e796',NULL),('cfb27ce5-b077-4fa2-ba22-ad9331c7e986',NULL,'reset-credentials-choose-user','dcc080c5-aede-4fd3-8b01-bd0928b730a2','06d58040-a014-4f13-a9ea-72835ad3c576',0,10,0,NULL,NULL),('d0d2d901-3b0d-45e0-ae28-b9f035245d7c',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','b241fbb5-0c09-4d3d-9ae0-9245a055ff6f',2,20,1,'5fca5bf2-7017-4f4a-865d-1c3db899423c',NULL),('d25ee1e0-3ff6-4b18-89a0-7f62a4d42398',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','5dff23ef-4265-41b9-a482-a2749930bd22',0,20,1,'c73fcc66-40d1-4495-bc79-5f10c9a1b722',NULL),('d32a005c-ff91-4458-b62a-9f27814eb932',NULL,'reset-otp','4327ba47-4116-44ea-9c4d-02907dca81e7','6153dc60-22f4-4624-a6a2-a5d0a8ee1ad0',0,20,0,NULL,NULL),('d3b99a05-8426-42ad-83e1-4b6dc9b36e1e',NULL,'reset-credentials-choose-user','40ae881c-f4e4-4b07-b097-a67d2bf515e6','d5b7af02-aba4-4b1c-bcca-e068b83dfee1',0,10,0,NULL,NULL),('d4a410e2-e5cc-48fd-bb5e-1e92178bb3a6',NULL,'client-secret-jwt','4327ba47-4116-44ea-9c4d-02907dca81e7','f378b286-fd54-4b8f-8615-df0118dacd24',2,30,0,NULL,NULL),('e02dff5c-7deb-4730-a2aa-3d61f187d7a1',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','7dc4a4fa-cefd-4f10-b196-f48d39cfa0db',1,20,1,'abb9c598-7f17-410f-a338-82f1f2e165fa',NULL),('e0b6dfd0-5236-4f3e-8d5d-aa7af7ccf109',NULL,'client-secret','40ae881c-f4e4-4b07-b097-a67d2bf515e6','a9a6d8d1-afbc-447a-8685-cd855fb3e627',2,10,0,NULL,NULL),('e4de4027-90c7-41b4-9f24-55aa53077258',NULL,'reset-password','dcc080c5-aede-4fd3-8b01-bd0928b730a2','06d58040-a014-4f13-a9ea-72835ad3c576',0,30,0,NULL,NULL),('e80458c8-c39c-4db4-a3c4-70097b887e5f',NULL,'auth-cookie','dcc080c5-aede-4fd3-8b01-bd0928b730a2','6bfae19d-d2c6-4dd5-836c-6042247853cc',2,10,0,NULL,NULL),('e858215b-bb62-4c5a-87e0-b8da424ab372',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','34156d67-bf62-439e-aca9-5da894adc723',2,20,1,'2ef963b2-037d-4185-8c73-503fc6e878da',NULL),('e9120b25-b912-4bce-83fd-b5a8dfd71589',NULL,'registration-page-form','4327ba47-4116-44ea-9c4d-02907dca81e7','8b55c406-d1d6-4a0c-bffa-4f2bf77f97f3',0,10,1,'76074602-02dc-43ec-a571-34efecba9953',NULL),('ebfb72cc-4b31-473f-af64-348c9676d5e2',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','2d986fa8-a731-473e-ad1d-a3bcfe9f6724',0,20,1,'34156d67-bf62-439e-aca9-5da894adc723',NULL),('eeba8010-3d7c-4d75-9354-0415e2ebe65a',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','f97a7702-55e0-4f57-b4a1-ae4b09239fa8',0,20,1,'96f3f037-c48e-451a-bc28-4bcd33f6ae3b',NULL),('f09def32-f13d-46f1-9aa6-fd91b4f3bd28',NULL,'reset-password','4327ba47-4116-44ea-9c4d-02907dca81e7','616a6379-b1b2-4e14-96ea-db9d45b16538',0,30,0,NULL,NULL),('f38f3577-101f-4dc7-8345-d133a7f102f2',NULL,'registration-recaptcha-action','40ae881c-f4e4-4b07-b097-a67d2bf515e6','e9314680-3395-49c7-a314-82fcff06249d',3,60,0,NULL,NULL),('f535fc3b-183f-40b6-93e2-345565803fb1',NULL,NULL,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','96f3f037-c48e-451a-bc28-4bcd33f6ae3b',2,20,1,'5dff23ef-4265-41b9-a482-a2749930bd22',NULL),('f61121ef-1af2-40a2-9bf7-55e0f41c76b2',NULL,'conditional-user-configured','4327ba47-4116-44ea-9c4d-02907dca81e7','abb9c598-7f17-410f-a338-82f1f2e165fa',0,10,0,NULL,NULL),('f70e17d4-e46e-4174-9d16-328ae08d2e0b',NULL,NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','74c902dc-6c87-4332-b633-f5fee94e6286',2,20,1,'65417fbf-411b-4301-8a33-234ae9d9ee5e',NULL),('f713cda0-53b2-4b19-81d2-67b24b08846a',NULL,'idp-create-user-if-unique','40ae881c-f4e4-4b07-b097-a67d2bf515e6','74c902dc-6c87-4332-b633-f5fee94e6286',2,10,0,NULL,'bf72caf8-1b4d-4f61-873f-5ef4fbff9a7b'),('f7bc9091-48c6-412e-a081-d81130336757',NULL,NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','48fb5d06-357a-463e-9263-1de5c87f092c',0,20,1,'77637b0f-15f0-4e0f-9e54-4294535f01c8',NULL),('fa4c1762-eca0-4914-a4ff-a0a3c7ad42bc',NULL,'auth-username-password-form','4327ba47-4116-44ea-9c4d-02907dca81e7','7dc4a4fa-cefd-4f10-b196-f48d39cfa0db',0,10,0,NULL,NULL),('fb0ac9d2-e7f8-4451-b67e-fbe637c8fa8e',NULL,'idp-confirm-link','dcc080c5-aede-4fd3-8b01-bd0928b730a2','5dff23ef-4265-41b9-a482-a2749930bd22',0,10,0,NULL,NULL),('fd50c8e4-0b5e-47b6-9091-ae9369f2d11e',NULL,'conditional-user-configured','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1a26f0eb-5055-4ece-a233-22ae767e5e7a',0,10,0,NULL,NULL),('fe6ee2c5-5b47-4798-a7f1-b3e69e9a6dfb',NULL,'conditional-user-configured','dcc080c5-aede-4fd3-8b01-bd0928b730a2','7bd6316d-df7b-43c9-b72f-af3f62def6c1',0,10,0,NULL,NULL);
/*!40000 ALTER TABLE `AUTHENTICATION_EXECUTION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `AUTHENTICATION_FLOW`
--

DROP TABLE IF EXISTS `AUTHENTICATION_FLOW`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AUTHENTICATION_FLOW` (
  `ID` varchar(36) NOT NULL,
  `ALIAS` varchar(255) DEFAULT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  `PROVIDER_ID` varchar(36) NOT NULL DEFAULT 'basic-flow',
  `TOP_LEVEL` tinyint NOT NULL DEFAULT '0',
  `BUILT_IN` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  KEY `IDX_AUTH_FLOW_REALM` (`REALM_ID`),
  CONSTRAINT `FK_AUTH_FLOW_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `AUTHENTICATION_FLOW`
--

LOCK TABLES `AUTHENTICATION_FLOW` WRITE;
/*!40000 ALTER TABLE `AUTHENTICATION_FLOW` DISABLE KEYS */;
INSERT INTO `AUTHENTICATION_FLOW` (`ID`, `ALIAS`, `DESCRIPTION`, `REALM_ID`, `PROVIDER_ID`, `TOP_LEVEL`, `BUILT_IN`) VALUES ('06d58040-a014-4f13-a9ea-72835ad3c576','reset credentials','Reset credentials for a user if they forgot their password or something','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',1,1),('08b0543c-51a6-4200-bb72-38404a916954','registration','registration flow','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',1,1),('0cfc44d0-dcec-4963-92be-00797e95107b','Verify Existing Account by Re-authentication','Reauthentication of existing account','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('1250dad6-ebc6-4a7e-a847-a83a67967a51','docker auth','Used by Docker clients to authenticate against the IDP','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',1,1),('15d69b0b-32fe-4c04-ba61-c6f2c92d1ae1','Copy of first broker login First broker login - Conditional OTP','Flow to determine if the OTP is required for the authentication','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,0),('1a26f0eb-5055-4ece-a233-22ae767e5e7a','Reset - Conditional OTP','Flow to determine if the OTP should be reset or not. Set to REQUIRED to force.','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('1dc10e9b-26a9-4420-aed4-cef7ce1cc0fc','Browser - Conditional OTP','Flow to determine if the OTP is required for the authentication','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('223dc23f-861c-4e10-89b1-a5fede09a82a','First broker login - Conditional OTP','Flow to determine if the OTP is required for the authentication','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('2d986fa8-a731-473e-ad1d-a3bcfe9f6724','Handle Existing Account','Handle what to do if there is existing account with same email/username like authenticated identity provider','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('2ef963b2-037d-4185-8c73-503fc6e878da','Verify Existing Account by Re-authentication','Reauthentication of existing account','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('330e968f-77f6-4862-9ba6-8b37217c3255','saml ecp','SAML ECP Profile Authentication Flow','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',1,1),('34156d67-bf62-439e-aca9-5da894adc723','Account verification options','Method with which to verity the existing account','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('3484027b-59da-4989-93dc-46251490ec58','Copy of first broker login Handle Existing Account','Handle what to do if there is existing account with same email/username like authenticated identity provider','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,0),('3af9837f-1ac4-4fe2-a3a7-579df7d70cbd','direct grant','OpenID Connect Resource Owner Grant','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',1,1),('3f87f84a-0a1e-4919-a380-10eabd214d58','Browser - Conditional OTP','Flow to determine if the OTP is required for the authentication','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('409007a9-4804-45b2-a701-79b9f94cea54','first broker login','Actions taken after first broker login with identity provider account, which is not yet linked to any Keycloak account','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',1,1),('48fb5d06-357a-463e-9263-1de5c87f092c','first broker login','Actions taken after first broker login with identity provider account, which is not yet linked to any Keycloak account','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',1,1),('4a21f571-4b69-4ef8-92af-09fe24026b64','browser','browser based authentication','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',1,1),('4e206e20-7402-45d2-ab7b-87d6d4b138e0','Direct Grant - Conditional OTP','Flow to determine if the OTP is required for the authentication','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('543dfaac-77e4-499e-aab0-851fd15d7db7','saml ecp','SAML ECP Profile Authentication Flow','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',1,1),('5dff23ef-4265-41b9-a482-a2749930bd22','Handle Existing Account','Handle what to do if there is existing account with same email/username like authenticated identity provider','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('5fca5bf2-7017-4f4a-865d-1c3db899423c','Verify Existing Account by Re-authentication','Reauthentication of existing account','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('6153dc60-22f4-4624-a6a2-a5d0a8ee1ad0','Reset - Conditional OTP','Flow to determine if the OTP should be reset or not. Set to REQUIRED to force.','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('616a6379-b1b2-4e14-96ea-db9d45b16538','reset credentials','Reset credentials for a user if they forgot their password or something','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',1,1),('62e95aab-0f29-402c-9d4e-1c062c966af7','Reset - Conditional OTP','Flow to determine if the OTP should be reset or not. Set to REQUIRED to force.','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('65417fbf-411b-4301-8a33-234ae9d9ee5e','Handle Existing Account','Handle what to do if there is existing account with same email/username like authenticated identity provider','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('692a7ac4-e4aa-47ac-8bb4-5267c7930f6d','docker auth','Used by Docker clients to authenticate against the IDP','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',1,1),('6bfae19d-d2c6-4dd5-836c-6042247853cc','browser','browser based authentication','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',1,1),('73f81fba-ce8f-4131-8637-3c64b1323d86','direct grant','OpenID Connect Resource Owner Grant','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',1,1),('74c902dc-6c87-4332-b633-f5fee94e6286','User creation or linking','Flow for the existing/non-existing user alternatives','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('76074602-02dc-43ec-a571-34efecba9953','registration form','registration form','4327ba47-4116-44ea-9c4d-02907dca81e7','form-flow',0,1),('77637b0f-15f0-4e0f-9e54-4294535f01c8','User creation or linking','Flow for the existing/non-existing user alternatives','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('7bd6316d-df7b-43c9-b72f-af3f62def6c1','Direct Grant - Conditional OTP','Flow to determine if the OTP is required for the authentication','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('7bf8b92b-a6be-4924-a010-8750f64b92e9','browser','browser based authentication','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',1,1),('7dc4a4fa-cefd-4f10-b196-f48d39cfa0db','forms','Username, password, otp and other auth forms.','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('80314021-add0-43af-8e10-dbf41e5fe483','Copy of first broker login Account verification options','Method with which to verity the existing account','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,0),('8b55c406-d1d6-4a0c-bffa-4f2bf77f97f3','registration','registration flow','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',1,1),('9589b999-1269-439f-93bc-eebce417f7ad','Copy of first broker login Verify Existing Account by Re-authentication','Reauthentication of existing account','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,0),('962742cd-4fa7-4654-a84f-10cda8492259','saml ecp','SAML ECP Profile Authentication Flow','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',1,1),('96f3f037-c48e-451a-bc28-4bcd33f6ae3b','User creation or linking','Flow for the existing/non-existing user alternatives','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('9c471d7d-c49c-4912-a921-7e8fbfa8c845','Direct Grant - Conditional OTP','Flow to determine if the OTP is required for the authentication','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('9e16df50-ccc9-4bd6-842d-7349838e3168','First broker login - Conditional OTP','Flow to determine if the OTP is required for the authentication','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('a86ff042-5592-47dd-bf01-7865d992e796','First broker login - Conditional OTP','Flow to determine if the OTP is required for the authentication','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('a9a6d8d1-afbc-447a-8685-cd855fb3e627','clients','Base authentication for clients','40ae881c-f4e4-4b07-b097-a67d2bf515e6','client-flow',1,1),('abb9c598-7f17-410f-a338-82f1f2e165fa','Browser - Conditional OTP','Flow to determine if the OTP is required for the authentication','4327ba47-4116-44ea-9c4d-02907dca81e7','basic-flow',0,1),('ac2d291b-2581-4830-9cc6-c123412d3072','forms','Username, password, otp and other auth forms.','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('b241fbb5-0c09-4d3d-9ae0-9245a055ff6f','Account verification options','Method with which to verity the existing account','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,1),('b639104c-892c-4843-a2ed-ca4324be5788','direct grant','OpenID Connect Resource Owner Grant','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',1,1),('b6c69e0d-d61b-45ca-acfc-1dc79b81bb62','forms','Username, password, otp and other auth forms.','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('be564502-8c4b-4eb6-a9e7-2927a798abf0','registration form','registration form','dcc080c5-aede-4fd3-8b01-bd0928b730a2','form-flow',0,1),('c4fe3a48-6c0f-4c9e-b28d-ffaf9b9ddfdc','Copy of first broker login User creation or linking','Flow for the existing/non-existing user alternatives','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',0,0),('c73fcc66-40d1-4495-bc79-5f10c9a1b722','Account verification options','Method with which to verity the existing account','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',0,1),('d5b7af02-aba4-4b1c-bcca-e068b83dfee1','reset credentials','Reset credentials for a user if they forgot their password or something','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',1,1),('e4be6c41-61ae-4230-a920-d66ab72fd11a','Copy of first broker login','Actions taken after first broker login with identity provider account, which is not yet linked to any Keycloak account','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',1,0),('e9314680-3395-49c7-a314-82fcff06249d','registration form','registration form','40ae881c-f4e4-4b07-b097-a67d2bf515e6','form-flow',0,1),('ea3b8b38-9a74-403e-8058-be535ccea217','docker auth','Used by Docker clients to authenticate against the IDP','40ae881c-f4e4-4b07-b097-a67d2bf515e6','basic-flow',1,1),('f378b286-fd54-4b8f-8615-df0118dacd24','clients','Base authentication for clients','4327ba47-4116-44ea-9c4d-02907dca81e7','client-flow',1,1),('f77f953b-5e25-4601-b047-748724d795b3','registration','registration flow','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',1,1),('f84a1f13-ef43-4a16-b9d6-63fe787b86dd','clients','Base authentication for clients','dcc080c5-aede-4fd3-8b01-bd0928b730a2','client-flow',1,1),('f97a7702-55e0-4f57-b4a1-ae4b09239fa8','first broker login','Actions taken after first broker login with identity provider account, which is not yet linked to any Keycloak account','dcc080c5-aede-4fd3-8b01-bd0928b730a2','basic-flow',1,1);
/*!40000 ALTER TABLE `AUTHENTICATION_FLOW` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `AUTHENTICATOR_CONFIG`
--

DROP TABLE IF EXISTS `AUTHENTICATOR_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AUTHENTICATOR_CONFIG` (
  `ID` varchar(36) NOT NULL,
  `ALIAS` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_AUTH_CONFIG_REALM` (`REALM_ID`),
  CONSTRAINT `FK_AUTH_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `AUTHENTICATOR_CONFIG`
--

LOCK TABLES `AUTHENTICATOR_CONFIG` WRITE;
/*!40000 ALTER TABLE `AUTHENTICATOR_CONFIG` DISABLE KEYS */;
INSERT INTO `AUTHENTICATOR_CONFIG` (`ID`, `ALIAS`, `REALM_ID`) VALUES ('0b74b1d6-2788-464b-8b28-59e468d0b488','review profile config','dcc080c5-aede-4fd3-8b01-bd0928b730a2'),('1f88cc4f-af88-488e-8304-3febe6c52f97','create unique user config','dcc080c5-aede-4fd3-8b01-bd0928b730a2'),('24523cd1-73bc-4810-932d-6438b849af6a','review profile config','4327ba47-4116-44ea-9c4d-02907dca81e7'),('3039100f-06c4-4330-b82b-a0a4b1bbd526','Copy of first broker login create unique user config','40ae881c-f4e4-4b07-b097-a67d2bf515e6'),('a2d51910-2b56-4570-9321-a74392ca0ec1','Copy of first broker login review profile config','40ae881c-f4e4-4b07-b097-a67d2bf515e6'),('ba96f2c7-fdca-4963-a250-2901eaa68ac8','Copy of first broker login create unique user config','40ae881c-f4e4-4b07-b097-a67d2bf515e6'),('bee1c985-9a19-48b5-9e0a-19b813f9a415','create unique user config','4327ba47-4116-44ea-9c4d-02907dca81e7'),('bf72caf8-1b4d-4f61-873f-5ef4fbff9a7b','create unique user config','40ae881c-f4e4-4b07-b097-a67d2bf515e6'),('d33fb1bf-6714-46fd-a0c5-0ad1c91aa4f2','review profile config','40ae881c-f4e4-4b07-b097-a67d2bf515e6'),('fa8ff9b0-26bb-442c-bfdb-505699cf6443','Copy of first broker login review profile config','40ae881c-f4e4-4b07-b097-a67d2bf515e6');
/*!40000 ALTER TABLE `AUTHENTICATOR_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `AUTHENTICATOR_CONFIG_ENTRY`
--

DROP TABLE IF EXISTS `AUTHENTICATOR_CONFIG_ENTRY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AUTHENTICATOR_CONFIG_ENTRY` (
  `AUTHENTICATOR_ID` varchar(36) NOT NULL,
  `VALUE` longtext,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`AUTHENTICATOR_ID`,`NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `AUTHENTICATOR_CONFIG_ENTRY`
--

LOCK TABLES `AUTHENTICATOR_CONFIG_ENTRY` WRITE;
/*!40000 ALTER TABLE `AUTHENTICATOR_CONFIG_ENTRY` DISABLE KEYS */;
INSERT INTO `AUTHENTICATOR_CONFIG_ENTRY` (`AUTHENTICATOR_ID`, `VALUE`, `NAME`) VALUES ('0b74b1d6-2788-464b-8b28-59e468d0b488','missing','update.profile.on.first.login'),('1f88cc4f-af88-488e-8304-3febe6c52f97','false','require.password.update.after.registration'),('24523cd1-73bc-4810-932d-6438b849af6a','missing','update.profile.on.first.login'),('3039100f-06c4-4330-b82b-a0a4b1bbd526','false','require.password.update.after.registration'),('a2d51910-2b56-4570-9321-a74392ca0ec1','missing','update.profile.on.first.login'),('ba96f2c7-fdca-4963-a250-2901eaa68ac8','false','require.password.update.after.registration'),('bee1c985-9a19-48b5-9e0a-19b813f9a415','false','require.password.update.after.registration'),('bf72caf8-1b4d-4f61-873f-5ef4fbff9a7b','false','require.password.update.after.registration'),('d33fb1bf-6714-46fd-a0c5-0ad1c91aa4f2','missing','update.profile.on.first.login'),('fa8ff9b0-26bb-442c-bfdb-505699cf6443','off','update.profile.on.first.login');
/*!40000 ALTER TABLE `AUTHENTICATOR_CONFIG_ENTRY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `BROKER_LINK`
--

DROP TABLE IF EXISTS `BROKER_LINK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `BROKER_LINK` (
  `IDENTITY_PROVIDER` varchar(255) NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `BROKER_USER_ID` varchar(255) DEFAULT NULL,
  `BROKER_USERNAME` varchar(255) DEFAULT NULL,
  `TOKEN` text,
  `USER_ID` varchar(255) NOT NULL,
  PRIMARY KEY (`IDENTITY_PROVIDER`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `BROKER_LINK`
--

LOCK TABLES `BROKER_LINK` WRITE;
/*!40000 ALTER TABLE `BROKER_LINK` DISABLE KEYS */;
/*!40000 ALTER TABLE `BROKER_LINK` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT`
--

DROP TABLE IF EXISTS `CLIENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT` (
  `ID` varchar(36) NOT NULL,
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `FULL_SCOPE_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `CLIENT_ID` varchar(255) DEFAULT NULL,
  `NOT_BEFORE` int DEFAULT NULL,
  `PUBLIC_CLIENT` tinyint NOT NULL DEFAULT '0',
  `SECRET` varchar(255) DEFAULT NULL,
  `BASE_URL` varchar(255) DEFAULT NULL,
  `BEARER_ONLY` tinyint NOT NULL DEFAULT '0',
  `MANAGEMENT_URL` varchar(255) DEFAULT NULL,
  `SURROGATE_AUTH_REQUIRED` tinyint NOT NULL DEFAULT '0',
  `REALM_ID` varchar(36) DEFAULT NULL,
  `PROTOCOL` varchar(255) DEFAULT NULL,
  `NODE_REREG_TIMEOUT` int DEFAULT '0',
  `FRONTCHANNEL_LOGOUT` tinyint NOT NULL DEFAULT '0',
  `CONSENT_REQUIRED` tinyint NOT NULL DEFAULT '0',
  `NAME` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `SERVICE_ACCOUNTS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `CLIENT_AUTHENTICATOR_TYPE` varchar(255) DEFAULT NULL,
  `ROOT_URL` varchar(255) DEFAULT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `REGISTRATION_TOKEN` varchar(255) DEFAULT NULL,
  `STANDARD_FLOW_ENABLED` tinyint NOT NULL DEFAULT '1',
  `IMPLICIT_FLOW_ENABLED` tinyint NOT NULL DEFAULT '0',
  `DIRECT_ACCESS_GRANTS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `ALWAYS_DISPLAY_IN_CONSOLE` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_B71CJLBENV945RB6GCON438AT` (`REALM_ID`,`CLIENT_ID`),
  KEY `IDX_CLIENT_ID` (`CLIENT_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT`
--

LOCK TABLES `CLIENT` WRITE;
/*!40000 ALTER TABLE `CLIENT` DISABLE KEYS */;
INSERT INTO `CLIENT` (`ID`, `ENABLED`, `FULL_SCOPE_ALLOWED`, `CLIENT_ID`, `NOT_BEFORE`, `PUBLIC_CLIENT`, `SECRET`, `BASE_URL`, `BEARER_ONLY`, `MANAGEMENT_URL`, `SURROGATE_AUTH_REQUIRED`, `REALM_ID`, `PROTOCOL`, `NODE_REREG_TIMEOUT`, `FRONTCHANNEL_LOGOUT`, `CONSENT_REQUIRED`, `NAME`, `SERVICE_ACCOUNTS_ENABLED`, `CLIENT_AUTHENTICATOR_TYPE`, `ROOT_URL`, `DESCRIPTION`, `REGISTRATION_TOKEN`, `STANDARD_FLOW_ENABLED`, `IMPLICIT_FLOW_ENABLED`, `DIRECT_ACCESS_GRANTS_ENABLED`, `ALWAYS_DISPLAY_IN_CONSOLE`) VALUES ('05021082-bfbc-4ec6-872a-e6c0916922c1',1,0,'security-admin-console',0,1,NULL,'/admin/master/console/',0,NULL,0,'4327ba47-4116-44ea-9c4d-02907dca81e7','openid-connect',0,0,0,'${client_security-admin-console}',0,'client-secret','${authAdminUrl}',NULL,NULL,1,0,0,0),('163f0b9a-c049-4744-8958-8a1e9efe572a',0,0,'admin-cli',0,1,NULL,'',0,'',0,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','openid-connect',0,0,0,'${client_admin-cli}',0,'client-secret','','',NULL,0,0,1,0),('1947f168-b049-4b78-8031-afcf98eae08d',0,0,'account',0,1,NULL,'/realms/test-realm/account/',0,'',0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','openid-connect',0,0,0,'${client_account}',0,'client-secret','${authBaseUrl}','',NULL,1,0,0,0),('1a0f720d-df73-4056-926d-10d520c81992',0,0,'account-console',0,1,NULL,'/realms/test-realm/account/',0,'',0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','openid-connect',0,0,0,'${client_account-console}',0,'client-secret','${authBaseUrl}','',NULL,1,0,0,0),('1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,0,'test-realm-realm',0,0,NULL,NULL,1,NULL,0,'4327ba47-4116-44ea-9c4d-02907dca81e7',NULL,0,0,0,'test-realm Realm',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5',1,0,'admin-cli',0,1,NULL,'',0,'',0,'4327ba47-4116-44ea-9c4d-02907dca81e7','openid-connect',0,0,0,'${client_admin-cli}',0,'client-secret','','',NULL,0,0,1,0),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0',0,0,'admin-cli',0,1,NULL,'',0,'',0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','openid-connect',0,0,0,'${client_admin-cli}',0,'client-secret','','',NULL,0,0,1,0),('320d776e-9e8e-4263-abf6-ebe3d862d547',0,0,'account-console',0,1,NULL,'/realms/TEST/account/',0,'',0,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','openid-connect',0,0,0,'${client_account-console}',0,'client-secret','${authBaseUrl}','',NULL,1,0,0,0),('33ce64f8-6ffd-430d-a260-c5c8f6d92308',0,0,'security-admin-console',0,1,NULL,'/admin/TEST/console/',0,'',0,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','openid-connect',0,0,0,'${client_security-admin-console}',0,'client-secret','${authAdminUrl}','',NULL,1,0,0,0),('41fe1d7e-3717-4965-8c77-e9868b3d98d8',0,0,'account',0,1,NULL,'/realms/TEST/account/',0,'',0,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','openid-connect',0,0,0,'${client_account}',0,'client-secret','${authBaseUrl}','',NULL,1,0,0,0),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,0,'realm-management',0,0,NULL,NULL,1,NULL,0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','openid-connect',0,0,0,'${client_realm-management}',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('657f6bd0-aa09-4703-99f7-e3f48e2de466',1,0,'account',0,1,NULL,'/realms/master/account/',0,'',0,'4327ba47-4116-44ea-9c4d-02907dca81e7','openid-connect',0,0,0,'${client_account}',0,'client-secret','${authBaseUrl}','',NULL,1,0,0,0),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f',1,1,'test',0,1,NULL,'http://localhost:3000/',0,'http://localhost:3000/',0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','openid-connect',-1,0,0,'',0,'client-secret','http://localhost:3000/','',NULL,1,0,1,0),('72c6029f-f8d2-4256-a326-2642c15f3a1e',1,0,'TEST-realm',0,0,NULL,NULL,1,NULL,0,'4327ba47-4116-44ea-9c4d-02907dca81e7',NULL,0,0,0,'TEST Realm',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('95e70707-5b66-41b0-a127-590f786b2fba',1,1,'TEST',0,1,NULL,'https://example.com',0,'https://example.com',0,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','openid-connect',-1,0,0,'ARC Brain dev',0,'client-secret','https://example.com','https://example.com/',NULL,1,0,0,0),('9e5fcd3f-0e4b-4627-8292-011015df3de6',1,0,'broker',0,0,NULL,NULL,1,NULL,0,'4327ba47-4116-44ea-9c4d-02907dca81e7','openid-connect',0,0,0,'${client_broker}',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f',0,0,'broker',0,0,NULL,NULL,1,'',0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','openid-connect',0,0,0,'${client_broker}',0,'client-secret',NULL,'',NULL,1,0,0,0),('a6b5e3ba-b72b-4028-92be-acb75bd541c1',0,0,'broker',0,0,NULL,NULL,1,'',0,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','openid-connect',0,0,0,'${client_broker}',0,'client-secret',NULL,'',NULL,1,0,0,0),('b1c791ca-b006-4f60-af37-20225d5876f5',1,0,'account-console',0,1,NULL,'/realms/master/account/',0,'',0,'4327ba47-4116-44ea-9c4d-02907dca81e7','openid-connect',0,0,0,'${client_account-console}',0,'client-secret','${authBaseUrl}','',NULL,1,0,0,0),('b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,0,'master-realm',0,0,NULL,NULL,1,NULL,0,'4327ba47-4116-44ea-9c4d-02907dca81e7',NULL,0,0,0,'master Realm',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('c12ebc78-3392-401a-8328-5dbb4cddf222',1,0,'security-admin-console',0,1,NULL,'/admin/test-realm/console/',0,NULL,0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','openid-connect',0,0,0,'${client_security-admin-console}',0,'client-secret','${authAdminUrl}',NULL,NULL,1,0,0,0),('ecabf8ab-e548-4b7a-a158-4d3e774afd77',1,1,'sample-app',0,0,'gjsghsDGpwmkxtTaEUrF24fa2352jRHrocmE8','https://example.com/',0,'https://example.com/',0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','openid-connect',-1,1,0,'sample-app',0,'client-secret','https://example.com/','sample-app for testing keycloak',NULL,1,0,0,1),('ff653271-7bc7-410a-99e0-cfeff6180f7b',0,0,'realm-management',0,0,NULL,NULL,1,'',0,'dcc080c5-aede-4fd3-8b01-bd0928b730a2','openid-connect',0,0,0,'${client_realm-management}',0,'client-secret',NULL,'',NULL,1,0,0,0);
/*!40000 ALTER TABLE `CLIENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_ATTRIBUTES`
--

DROP TABLE IF EXISTS `CLIENT_ATTRIBUTES`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_ATTRIBUTES` (
  `CLIENT_ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `VALUE` longtext CHARACTER SET utf8 COLLATE utf8_general_ci,
  PRIMARY KEY (`CLIENT_ID`,`NAME`),
  KEY `IDX_CLIENT_ATT_BY_NAME_VALUE` (`NAME`,`VALUE`(255)),
  CONSTRAINT `FK3C47C64BEACCA966` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_ATTRIBUTES`
--

LOCK TABLES `CLIENT_ATTRIBUTES` WRITE;
/*!40000 ALTER TABLE `CLIENT_ATTRIBUTES` DISABLE KEYS */;
INSERT INTO `CLIENT_ATTRIBUTES` (`CLIENT_ID`, `NAME`, `VALUE`) VALUES ('05021082-bfbc-4ec6-872a-e6c0916922c1','pkce.code.challenge.method','S256'),('05021082-bfbc-4ec6-872a-e6c0916922c1','post.logout.redirect.uris','+'),('163f0b9a-c049-4744-8958-8a1e9efe572a','backchannel.logout.revoke.offline.tokens','false'),('163f0b9a-c049-4744-8958-8a1e9efe572a','backchannel.logout.session.required','true'),('163f0b9a-c049-4744-8958-8a1e9efe572a','display.on.consent.screen','false'),('163f0b9a-c049-4744-8958-8a1e9efe572a','oauth2.device.authorization.grant.enabled','false'),('163f0b9a-c049-4744-8958-8a1e9efe572a','oidc.ciba.grant.enabled','false'),('1947f168-b049-4b78-8031-afcf98eae08d','backchannel.logout.revoke.offline.tokens','false'),('1947f168-b049-4b78-8031-afcf98eae08d','backchannel.logout.session.required','true'),('1947f168-b049-4b78-8031-afcf98eae08d','display.on.consent.screen','false'),('1947f168-b049-4b78-8031-afcf98eae08d','oauth2.device.authorization.grant.enabled','false'),('1947f168-b049-4b78-8031-afcf98eae08d','oidc.ciba.grant.enabled','false'),('1947f168-b049-4b78-8031-afcf98eae08d','post.logout.redirect.uris','+'),('1a0f720d-df73-4056-926d-10d520c81992','backchannel.logout.revoke.offline.tokens','false'),('1a0f720d-df73-4056-926d-10d520c81992','backchannel.logout.session.required','true'),('1a0f720d-df73-4056-926d-10d520c81992','display.on.consent.screen','false'),('1a0f720d-df73-4056-926d-10d520c81992','oauth2.device.authorization.grant.enabled','false'),('1a0f720d-df73-4056-926d-10d520c81992','oidc.ciba.grant.enabled','false'),('1a0f720d-df73-4056-926d-10d520c81992','pkce.code.challenge.method','S256'),('1a0f720d-df73-4056-926d-10d520c81992','post.logout.redirect.uris','+'),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','backchannel.logout.revoke.offline.tokens','false'),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','backchannel.logout.session.required','true'),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','display.on.consent.screen','false'),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','oauth2.device.authorization.grant.enabled','false'),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','oidc.ciba.grant.enabled','false'),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','backchannel.logout.revoke.offline.tokens','false'),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','backchannel.logout.session.required','true'),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','display.on.consent.screen','false'),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','oauth2.device.authorization.grant.enabled','false'),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','oidc.ciba.grant.enabled','false'),('320d776e-9e8e-4263-abf6-ebe3d862d547','backchannel.logout.revoke.offline.tokens','false'),('320d776e-9e8e-4263-abf6-ebe3d862d547','backchannel.logout.session.required','true'),('320d776e-9e8e-4263-abf6-ebe3d862d547','display.on.consent.screen','false'),('320d776e-9e8e-4263-abf6-ebe3d862d547','oauth2.device.authorization.grant.enabled','false'),('320d776e-9e8e-4263-abf6-ebe3d862d547','oidc.ciba.grant.enabled','false'),('320d776e-9e8e-4263-abf6-ebe3d862d547','pkce.code.challenge.method','S256'),('320d776e-9e8e-4263-abf6-ebe3d862d547','post.logout.redirect.uris','+'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','backchannel.logout.revoke.offline.tokens','false'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','backchannel.logout.session.required','true'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','display.on.consent.screen','false'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','oauth2.device.authorization.grant.enabled','false'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','oidc.ciba.grant.enabled','false'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','pkce.code.challenge.method','S256'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','post.logout.redirect.uris','+'),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','backchannel.logout.revoke.offline.tokens','false'),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','backchannel.logout.session.required','true'),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','display.on.consent.screen','false'),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','oauth2.device.authorization.grant.enabled','false'),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','oidc.ciba.grant.enabled','false'),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','post.logout.redirect.uris','+'),('657f6bd0-aa09-4703-99f7-e3f48e2de466','backchannel.logout.revoke.offline.tokens','false'),('657f6bd0-aa09-4703-99f7-e3f48e2de466','backchannel.logout.session.required','true'),('657f6bd0-aa09-4703-99f7-e3f48e2de466','display.on.consent.screen','false'),('657f6bd0-aa09-4703-99f7-e3f48e2de466','oauth2.device.authorization.grant.enabled','false'),('657f6bd0-aa09-4703-99f7-e3f48e2de466','oidc.ciba.grant.enabled','false'),('657f6bd0-aa09-4703-99f7-e3f48e2de466','post.logout.redirect.uris','+'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','acr.loa.map','{}'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','backchannel.logout.revoke.offline.tokens','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','backchannel.logout.session.required','true'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','client_credentials.use_refresh_token','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','client.use.lightweight.access.token.enabled','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','display.on.consent.screen','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','login_theme','custom'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','oauth2.device.authorization.grant.enabled','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','oidc.ciba.grant.enabled','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','post.logout.redirect.uris','http://localhost:3000/*'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','require.pushed.authorization.requests','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','tls.client.certificate.bound.access.tokens','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','token.response.type.bearer.lower-case','false'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','use.refresh.tokens','false'),('95e70707-5b66-41b0-a127-590f786b2fba','acr.loa.map','{}'),('95e70707-5b66-41b0-a127-590f786b2fba','backchannel.logout.revoke.offline.tokens','false'),('95e70707-5b66-41b0-a127-590f786b2fba','backchannel.logout.session.required','true'),('95e70707-5b66-41b0-a127-590f786b2fba','client_credentials.use_refresh_token','false'),('95e70707-5b66-41b0-a127-590f786b2fba','client.use.lightweight.access.token.enabled','false'),('95e70707-5b66-41b0-a127-590f786b2fba','display.on.consent.screen','false'),('95e70707-5b66-41b0-a127-590f786b2fba','oauth2.device.authorization.grant.enabled','false'),('95e70707-5b66-41b0-a127-590f786b2fba','oidc.ciba.grant.enabled','false'),('95e70707-5b66-41b0-a127-590f786b2fba','post.logout.redirect.uris','https://example.com/*'),('95e70707-5b66-41b0-a127-590f786b2fba','require.pushed.authorization.requests','false'),('95e70707-5b66-41b0-a127-590f786b2fba','tls.client.certificate.bound.access.tokens','false'),('95e70707-5b66-41b0-a127-590f786b2fba','token.response.type.bearer.lower-case','false'),('95e70707-5b66-41b0-a127-590f786b2fba','use.refresh.tokens','false'),('b1c791ca-b006-4f60-af37-20225d5876f5','backchannel.logout.revoke.offline.tokens','false'),('b1c791ca-b006-4f60-af37-20225d5876f5','backchannel.logout.session.required','true'),('b1c791ca-b006-4f60-af37-20225d5876f5','display.on.consent.screen','false'),('b1c791ca-b006-4f60-af37-20225d5876f5','oauth2.device.authorization.grant.enabled','false'),('b1c791ca-b006-4f60-af37-20225d5876f5','oidc.ciba.grant.enabled','false'),('b1c791ca-b006-4f60-af37-20225d5876f5','pkce.code.challenge.method','S256'),('b1c791ca-b006-4f60-af37-20225d5876f5','post.logout.redirect.uris','+'),('c12ebc78-3392-401a-8328-5dbb4cddf222','pkce.code.challenge.method','S256'),('c12ebc78-3392-401a-8328-5dbb4cddf222','post.logout.redirect.uris','+'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','acr.loa.map','{}'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','backchannel.logout.revoke.offline.tokens','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','backchannel.logout.session.required','true'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','client_credentials.use_refresh_token','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','client.secret.creation.time','1713428888'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','client.use.lightweight.access.token.enabled','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','display.on.consent.screen','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','oauth2.device.authorization.grant.enabled','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','oidc.ciba.grant.enabled','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','post.logout.redirect.uris','https://example.com/*'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','require.pushed.authorization.requests','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','tls.client.certificate.bound.access.tokens','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','token.response.type.bearer.lower-case','false'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','use.refresh.tokens','true');
/*!40000 ALTER TABLE `CLIENT_ATTRIBUTES` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_AUTH_FLOW_BINDINGS`
--

DROP TABLE IF EXISTS `CLIENT_AUTH_FLOW_BINDINGS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_AUTH_FLOW_BINDINGS` (
  `CLIENT_ID` varchar(36) NOT NULL,
  `FLOW_ID` varchar(36) DEFAULT NULL,
  `BINDING_NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`BINDING_NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_AUTH_FLOW_BINDINGS`
--

LOCK TABLES `CLIENT_AUTH_FLOW_BINDINGS` WRITE;
/*!40000 ALTER TABLE `CLIENT_AUTH_FLOW_BINDINGS` DISABLE KEYS */;
INSERT INTO `CLIENT_AUTH_FLOW_BINDINGS` (`CLIENT_ID`, `FLOW_ID`, `BINDING_NAME`) VALUES ('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','4a21f571-4b69-4ef8-92af-09fe24026b64','browser'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','73f81fba-ce8f-4131-8637-3c64b1323d86','direct_grant'),('95e70707-5b66-41b0-a127-590f786b2fba','6bfae19d-d2c6-4dd5-836c-6042247853cc','browser'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','4a21f571-4b69-4ef8-92af-09fe24026b64','browser');
/*!40000 ALTER TABLE `CLIENT_AUTH_FLOW_BINDINGS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_INITIAL_ACCESS`
--

DROP TABLE IF EXISTS `CLIENT_INITIAL_ACCESS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_INITIAL_ACCESS` (
  `ID` varchar(36) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `TIMESTAMP` int DEFAULT NULL,
  `EXPIRATION` int DEFAULT NULL,
  `COUNT` int DEFAULT NULL,
  `REMAINING_COUNT` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_CLIENT_INIT_ACC_REALM` (`REALM_ID`),
  CONSTRAINT `FK_CLIENT_INIT_ACC_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_INITIAL_ACCESS`
--

LOCK TABLES `CLIENT_INITIAL_ACCESS` WRITE;
/*!40000 ALTER TABLE `CLIENT_INITIAL_ACCESS` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_INITIAL_ACCESS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_NODE_REGISTRATIONS`
--

DROP TABLE IF EXISTS `CLIENT_NODE_REGISTRATIONS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_NODE_REGISTRATIONS` (
  `CLIENT_ID` varchar(36) NOT NULL,
  `VALUE` int DEFAULT NULL,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`NAME`),
  CONSTRAINT `FK4129723BA992F594` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_NODE_REGISTRATIONS`
--

LOCK TABLES `CLIENT_NODE_REGISTRATIONS` WRITE;
/*!40000 ALTER TABLE `CLIENT_NODE_REGISTRATIONS` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_NODE_REGISTRATIONS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SCOPE`
--

DROP TABLE IF EXISTS `CLIENT_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SCOPE` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `PROTOCOL` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_CLI_SCOPE` (`REALM_ID`,`NAME`),
  KEY `IDX_REALM_CLSCOPE` (`REALM_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SCOPE`
--

LOCK TABLES `CLIENT_SCOPE` WRITE;
/*!40000 ALTER TABLE `CLIENT_SCOPE` DISABLE KEYS */;
INSERT INTO `CLIENT_SCOPE` (`ID`, `NAME`, `REALM_ID`, `DESCRIPTION`, `PROTOCOL`) VALUES ('066554e5-ed72-4255-a59e-37bce592656c','address','4327ba47-4116-44ea-9c4d-02907dca81e7','OpenID Connect built-in scope: address','openid-connect'),('091b2057-69d2-4a0e-a76c-7cbd76698850','role_list','4327ba47-4116-44ea-9c4d-02907dca81e7','SAML role list','saml'),('14d13053-c1b9-43f8-b993-9d7a8aef8069','acr','40ae881c-f4e4-4b07-b097-a67d2bf515e6','OpenID Connect scope for add acr (authentication context class reference) to the token','openid-connect'),('15afd007-2320-4439-853b-2ddaa8b2ff71','microprofile-jwt','4327ba47-4116-44ea-9c4d-02907dca81e7','Microprofile - JWT built-in scope','openid-connect'),('192eb290-9206-4ec7-94ce-f62e6ab82179','microprofile-jwt','40ae881c-f4e4-4b07-b097-a67d2bf515e6','Microprofile - JWT built-in scope','openid-connect'),('282b6071-c081-4ebf-bb56-5241153f1811','offline_access','40ae881c-f4e4-4b07-b097-a67d2bf515e6','OpenID Connect built-in scope: offline_access','openid-connect'),('29003608-50c2-4693-953b-316f0c71fa25','roles','dcc080c5-aede-4fd3-8b01-bd0928b730a2','OpenID Connect scope for add user roles to the access token','openid-connect'),('310fbe9a-4f1f-4683-aca4-516d640a5d9b','address','40ae881c-f4e4-4b07-b097-a67d2bf515e6','OpenID Connect built-in scope: address','openid-connect'),('39474b38-4002-4590-a90e-7e77c25785b5','acr','4327ba47-4116-44ea-9c4d-02907dca81e7','OpenID Connect scope for add acr (authentication context class reference) to the token','openid-connect'),('3d4aaffa-8399-4bff-9dac-000e43b45ca9','email','dcc080c5-aede-4fd3-8b01-bd0928b730a2','OpenID Connect built-in scope: email','openid-connect'),('3e744bf6-ec28-46cf-af95-5cf3f4ee8d58','email','4327ba47-4116-44ea-9c4d-02907dca81e7','OpenID Connect built-in scope: email','openid-connect'),('40f0b7c4-e882-42e9-8cc2-00055b119ca8','profile','40ae881c-f4e4-4b07-b097-a67d2bf515e6','OpenID Connect built-in scope: profile','openid-connect'),('65b407fe-a29a-4caa-90c1-060588438771','address','dcc080c5-aede-4fd3-8b01-bd0928b730a2','OpenID Connect built-in scope: address','openid-connect'),('6bffc256-3149-4d67-a5d2-888076a43a46','role_list','40ae881c-f4e4-4b07-b097-a67d2bf515e6','SAML role list','saml'),('6eed7521-b56d-4af5-8927-771c01723b21','phone','40ae881c-f4e4-4b07-b097-a67d2bf515e6','OpenID Connect built-in scope: phone','openid-connect'),('74beed12-a014-4842-91e7-5334c22a5bc0','web-origins','dcc080c5-aede-4fd3-8b01-bd0928b730a2','OpenID Connect scope for add allowed web origins to the access token','openid-connect'),('7b474ad1-7a28-4e53-970d-e9a5a14bfbab','acr','dcc080c5-aede-4fd3-8b01-bd0928b730a2','OpenID Connect scope for add acr (authentication context class reference) to the token','openid-connect'),('815032a2-e511-4d5f-a543-e9a3ee4f648c','email','40ae881c-f4e4-4b07-b097-a67d2bf515e6','OpenID Connect built-in scope: email','openid-connect'),('90977ab5-207e-4b07-836d-e8d59805927b','roles','4327ba47-4116-44ea-9c4d-02907dca81e7','OpenID Connect scope for add user roles to the access token','openid-connect'),('9c93a802-d9c8-468e-8983-129fc287b26c','profile','dcc080c5-aede-4fd3-8b01-bd0928b730a2','OpenID Connect built-in scope: profile','openid-connect'),('aaa43c0a-2c37-42b7-a80c-f98eef322343','web-origins','40ae881c-f4e4-4b07-b097-a67d2bf515e6','OpenID Connect scope for add allowed web origins to the access token','openid-connect'),('b9388bd3-40d7-4f4a-87e8-2d78275ee434','roles','40ae881c-f4e4-4b07-b097-a67d2bf515e6','OpenID Connect scope for add user roles to the access token','openid-connect'),('bd4873a9-1ad5-40ba-bec3-5506355835e3','offline_access','dcc080c5-aede-4fd3-8b01-bd0928b730a2','OpenID Connect built-in scope: offline_access','openid-connect'),('bdb21a59-753f-498c-8be2-9fae1f633cd3','role_list','dcc080c5-aede-4fd3-8b01-bd0928b730a2','SAML role list','saml'),('c294c451-1979-493f-ad49-15220c9bd6f8','web-origins','4327ba47-4116-44ea-9c4d-02907dca81e7','OpenID Connect scope for add allowed web origins to the access token','openid-connect'),('c36c9c18-f3c5-451b-bf26-2887df9d803d','offline_access','4327ba47-4116-44ea-9c4d-02907dca81e7','OpenID Connect built-in scope: offline_access','openid-connect'),('d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70','phone','dcc080c5-aede-4fd3-8b01-bd0928b730a2','OpenID Connect built-in scope: phone','openid-connect'),('ef9560b3-735c-4882-9aa1-dce8fe0697af','profile','4327ba47-4116-44ea-9c4d-02907dca81e7','OpenID Connect built-in scope: profile','openid-connect'),('fafceb53-8f99-4a18-81bf-59c2b83bcd6c','microprofile-jwt','dcc080c5-aede-4fd3-8b01-bd0928b730a2','Microprofile - JWT built-in scope','openid-connect'),('fd09aae5-67c6-4162-a03e-95f8cc07d053','phone','4327ba47-4116-44ea-9c4d-02907dca81e7','OpenID Connect built-in scope: phone','openid-connect');
/*!40000 ALTER TABLE `CLIENT_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SCOPE_ATTRIBUTES`
--

DROP TABLE IF EXISTS `CLIENT_SCOPE_ATTRIBUTES`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SCOPE_ATTRIBUTES` (
  `SCOPE_ID` varchar(36) NOT NULL,
  `VALUE` text,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`SCOPE_ID`,`NAME`),
  KEY `IDX_CLSCOPE_ATTRS` (`SCOPE_ID`),
  CONSTRAINT `FK_CL_SCOPE_ATTR_SCOPE` FOREIGN KEY (`SCOPE_ID`) REFERENCES `CLIENT_SCOPE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SCOPE_ATTRIBUTES`
--

LOCK TABLES `CLIENT_SCOPE_ATTRIBUTES` WRITE;
/*!40000 ALTER TABLE `CLIENT_SCOPE_ATTRIBUTES` DISABLE KEYS */;
INSERT INTO `CLIENT_SCOPE_ATTRIBUTES` (`SCOPE_ID`, `VALUE`, `NAME`) VALUES ('066554e5-ed72-4255-a59e-37bce592656c','${addressScopeConsentText}','consent.screen.text'),('066554e5-ed72-4255-a59e-37bce592656c','true','display.on.consent.screen'),('066554e5-ed72-4255-a59e-37bce592656c','true','include.in.token.scope'),('091b2057-69d2-4a0e-a76c-7cbd76698850','${samlRoleListScopeConsentText}','consent.screen.text'),('091b2057-69d2-4a0e-a76c-7cbd76698850','true','display.on.consent.screen'),('14d13053-c1b9-43f8-b993-9d7a8aef8069','false','display.on.consent.screen'),('14d13053-c1b9-43f8-b993-9d7a8aef8069','false','include.in.token.scope'),('15afd007-2320-4439-853b-2ddaa8b2ff71','false','display.on.consent.screen'),('15afd007-2320-4439-853b-2ddaa8b2ff71','true','include.in.token.scope'),('192eb290-9206-4ec7-94ce-f62e6ab82179','false','display.on.consent.screen'),('192eb290-9206-4ec7-94ce-f62e6ab82179','true','include.in.token.scope'),('282b6071-c081-4ebf-bb56-5241153f1811','${offlineAccessScopeConsentText}','consent.screen.text'),('282b6071-c081-4ebf-bb56-5241153f1811','true','display.on.consent.screen'),('29003608-50c2-4693-953b-316f0c71fa25','${rolesScopeConsentText}','consent.screen.text'),('29003608-50c2-4693-953b-316f0c71fa25','true','display.on.consent.screen'),('29003608-50c2-4693-953b-316f0c71fa25','false','include.in.token.scope'),('310fbe9a-4f1f-4683-aca4-516d640a5d9b','${addressScopeConsentText}','consent.screen.text'),('310fbe9a-4f1f-4683-aca4-516d640a5d9b','true','display.on.consent.screen'),('310fbe9a-4f1f-4683-aca4-516d640a5d9b','true','include.in.token.scope'),('39474b38-4002-4590-a90e-7e77c25785b5','false','display.on.consent.screen'),('39474b38-4002-4590-a90e-7e77c25785b5','false','include.in.token.scope'),('3d4aaffa-8399-4bff-9dac-000e43b45ca9','${emailScopeConsentText}','consent.screen.text'),('3d4aaffa-8399-4bff-9dac-000e43b45ca9','true','display.on.consent.screen'),('3d4aaffa-8399-4bff-9dac-000e43b45ca9','true','include.in.token.scope'),('3e744bf6-ec28-46cf-af95-5cf3f4ee8d58','${emailScopeConsentText}','consent.screen.text'),('3e744bf6-ec28-46cf-af95-5cf3f4ee8d58','true','display.on.consent.screen'),('3e744bf6-ec28-46cf-af95-5cf3f4ee8d58','true','include.in.token.scope'),('40f0b7c4-e882-42e9-8cc2-00055b119ca8','${profileScopeConsentText}','consent.screen.text'),('40f0b7c4-e882-42e9-8cc2-00055b119ca8','true','display.on.consent.screen'),('40f0b7c4-e882-42e9-8cc2-00055b119ca8','true','include.in.token.scope'),('65b407fe-a29a-4caa-90c1-060588438771','${addressScopeConsentText}','consent.screen.text'),('65b407fe-a29a-4caa-90c1-060588438771','true','display.on.consent.screen'),('65b407fe-a29a-4caa-90c1-060588438771','true','include.in.token.scope'),('6bffc256-3149-4d67-a5d2-888076a43a46','${samlRoleListScopeConsentText}','consent.screen.text'),('6bffc256-3149-4d67-a5d2-888076a43a46','true','display.on.consent.screen'),('6eed7521-b56d-4af5-8927-771c01723b21','${phoneScopeConsentText}','consent.screen.text'),('6eed7521-b56d-4af5-8927-771c01723b21','true','display.on.consent.screen'),('6eed7521-b56d-4af5-8927-771c01723b21','true','include.in.token.scope'),('74beed12-a014-4842-91e7-5334c22a5bc0','','consent.screen.text'),('74beed12-a014-4842-91e7-5334c22a5bc0','false','display.on.consent.screen'),('74beed12-a014-4842-91e7-5334c22a5bc0','false','include.in.token.scope'),('7b474ad1-7a28-4e53-970d-e9a5a14bfbab','false','display.on.consent.screen'),('7b474ad1-7a28-4e53-970d-e9a5a14bfbab','false','include.in.token.scope'),('815032a2-e511-4d5f-a543-e9a3ee4f648c','${emailScopeConsentText}','consent.screen.text'),('815032a2-e511-4d5f-a543-e9a3ee4f648c','true','display.on.consent.screen'),('815032a2-e511-4d5f-a543-e9a3ee4f648c','true','include.in.token.scope'),('90977ab5-207e-4b07-836d-e8d59805927b','${rolesScopeConsentText}','consent.screen.text'),('90977ab5-207e-4b07-836d-e8d59805927b','true','display.on.consent.screen'),('90977ab5-207e-4b07-836d-e8d59805927b','false','include.in.token.scope'),('9c93a802-d9c8-468e-8983-129fc287b26c','${profileScopeConsentText}','consent.screen.text'),('9c93a802-d9c8-468e-8983-129fc287b26c','true','display.on.consent.screen'),('9c93a802-d9c8-468e-8983-129fc287b26c','true','include.in.token.scope'),('aaa43c0a-2c37-42b7-a80c-f98eef322343','','consent.screen.text'),('aaa43c0a-2c37-42b7-a80c-f98eef322343','false','display.on.consent.screen'),('aaa43c0a-2c37-42b7-a80c-f98eef322343','false','include.in.token.scope'),('b9388bd3-40d7-4f4a-87e8-2d78275ee434','${rolesScopeConsentText}','consent.screen.text'),('b9388bd3-40d7-4f4a-87e8-2d78275ee434','true','display.on.consent.screen'),('b9388bd3-40d7-4f4a-87e8-2d78275ee434','false','include.in.token.scope'),('bd4873a9-1ad5-40ba-bec3-5506355835e3','${offlineAccessScopeConsentText}','consent.screen.text'),('bd4873a9-1ad5-40ba-bec3-5506355835e3','true','display.on.consent.screen'),('bdb21a59-753f-498c-8be2-9fae1f633cd3','${samlRoleListScopeConsentText}','consent.screen.text'),('bdb21a59-753f-498c-8be2-9fae1f633cd3','true','display.on.consent.screen'),('c294c451-1979-493f-ad49-15220c9bd6f8','','consent.screen.text'),('c294c451-1979-493f-ad49-15220c9bd6f8','false','display.on.consent.screen'),('c294c451-1979-493f-ad49-15220c9bd6f8','false','include.in.token.scope'),('c36c9c18-f3c5-451b-bf26-2887df9d803d','${offlineAccessScopeConsentText}','consent.screen.text'),('c36c9c18-f3c5-451b-bf26-2887df9d803d','true','display.on.consent.screen'),('d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70','${phoneScopeConsentText}','consent.screen.text'),('d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70','true','display.on.consent.screen'),('d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70','true','include.in.token.scope'),('ef9560b3-735c-4882-9aa1-dce8fe0697af','${profileScopeConsentText}','consent.screen.text'),('ef9560b3-735c-4882-9aa1-dce8fe0697af','true','display.on.consent.screen'),('ef9560b3-735c-4882-9aa1-dce8fe0697af','true','include.in.token.scope'),('fafceb53-8f99-4a18-81bf-59c2b83bcd6c','false','display.on.consent.screen'),('fafceb53-8f99-4a18-81bf-59c2b83bcd6c','true','include.in.token.scope'),('fd09aae5-67c6-4162-a03e-95f8cc07d053','${phoneScopeConsentText}','consent.screen.text'),('fd09aae5-67c6-4162-a03e-95f8cc07d053','true','display.on.consent.screen'),('fd09aae5-67c6-4162-a03e-95f8cc07d053','true','include.in.token.scope');
/*!40000 ALTER TABLE `CLIENT_SCOPE_ATTRIBUTES` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SCOPE_CLIENT`
--

DROP TABLE IF EXISTS `CLIENT_SCOPE_CLIENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SCOPE_CLIENT` (
  `CLIENT_ID` varchar(255) NOT NULL,
  `SCOPE_ID` varchar(255) NOT NULL,
  `DEFAULT_SCOPE` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`CLIENT_ID`,`SCOPE_ID`),
  KEY `IDX_CLSCOPE_CL` (`CLIENT_ID`),
  KEY `IDX_CL_CLSCOPE` (`SCOPE_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SCOPE_CLIENT`
--

LOCK TABLES `CLIENT_SCOPE_CLIENT` WRITE;
/*!40000 ALTER TABLE `CLIENT_SCOPE_CLIENT` DISABLE KEYS */;
INSERT INTO `CLIENT_SCOPE_CLIENT` (`CLIENT_ID`, `SCOPE_ID`, `DEFAULT_SCOPE`) VALUES ('05021082-bfbc-4ec6-872a-e6c0916922c1','066554e5-ed72-4255-a59e-37bce592656c',0),('05021082-bfbc-4ec6-872a-e6c0916922c1','15afd007-2320-4439-853b-2ddaa8b2ff71',0),('05021082-bfbc-4ec6-872a-e6c0916922c1','39474b38-4002-4590-a90e-7e77c25785b5',1),('05021082-bfbc-4ec6-872a-e6c0916922c1','3e744bf6-ec28-46cf-af95-5cf3f4ee8d58',1),('05021082-bfbc-4ec6-872a-e6c0916922c1','90977ab5-207e-4b07-836d-e8d59805927b',1),('05021082-bfbc-4ec6-872a-e6c0916922c1','c294c451-1979-493f-ad49-15220c9bd6f8',1),('05021082-bfbc-4ec6-872a-e6c0916922c1','c36c9c18-f3c5-451b-bf26-2887df9d803d',0),('05021082-bfbc-4ec6-872a-e6c0916922c1','ef9560b3-735c-4882-9aa1-dce8fe0697af',1),('05021082-bfbc-4ec6-872a-e6c0916922c1','fd09aae5-67c6-4162-a03e-95f8cc07d053',0),('163f0b9a-c049-4744-8958-8a1e9efe572a','29003608-50c2-4693-953b-316f0c71fa25',1),('163f0b9a-c049-4744-8958-8a1e9efe572a','3d4aaffa-8399-4bff-9dac-000e43b45ca9',1),('163f0b9a-c049-4744-8958-8a1e9efe572a','65b407fe-a29a-4caa-90c1-060588438771',0),('163f0b9a-c049-4744-8958-8a1e9efe572a','74beed12-a014-4842-91e7-5334c22a5bc0',1),('163f0b9a-c049-4744-8958-8a1e9efe572a','7b474ad1-7a28-4e53-970d-e9a5a14bfbab',1),('163f0b9a-c049-4744-8958-8a1e9efe572a','9c93a802-d9c8-468e-8983-129fc287b26c',1),('163f0b9a-c049-4744-8958-8a1e9efe572a','bd4873a9-1ad5-40ba-bec3-5506355835e3',0),('163f0b9a-c049-4744-8958-8a1e9efe572a','d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70',0),('163f0b9a-c049-4744-8958-8a1e9efe572a','fafceb53-8f99-4a18-81bf-59c2b83bcd6c',0),('1947f168-b049-4b78-8031-afcf98eae08d','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('1947f168-b049-4b78-8031-afcf98eae08d','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('1947f168-b049-4b78-8031-afcf98eae08d','282b6071-c081-4ebf-bb56-5241153f1811',0),('1947f168-b049-4b78-8031-afcf98eae08d','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('1947f168-b049-4b78-8031-afcf98eae08d','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('1947f168-b049-4b78-8031-afcf98eae08d','6eed7521-b56d-4af5-8927-771c01723b21',0),('1947f168-b049-4b78-8031-afcf98eae08d','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('1947f168-b049-4b78-8031-afcf98eae08d','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('1947f168-b049-4b78-8031-afcf98eae08d','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('1a0f720d-df73-4056-926d-10d520c81992','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('1a0f720d-df73-4056-926d-10d520c81992','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('1a0f720d-df73-4056-926d-10d520c81992','282b6071-c081-4ebf-bb56-5241153f1811',0),('1a0f720d-df73-4056-926d-10d520c81992','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('1a0f720d-df73-4056-926d-10d520c81992','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('1a0f720d-df73-4056-926d-10d520c81992','6eed7521-b56d-4af5-8927-771c01723b21',0),('1a0f720d-df73-4056-926d-10d520c81992','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('1a0f720d-df73-4056-926d-10d520c81992','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('1a0f720d-df73-4056-926d-10d520c81992','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','066554e5-ed72-4255-a59e-37bce592656c',0),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','15afd007-2320-4439-853b-2ddaa8b2ff71',0),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','39474b38-4002-4590-a90e-7e77c25785b5',1),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','3e744bf6-ec28-46cf-af95-5cf3f4ee8d58',1),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','90977ab5-207e-4b07-836d-e8d59805927b',1),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','c294c451-1979-493f-ad49-15220c9bd6f8',1),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','c36c9c18-f3c5-451b-bf26-2887df9d803d',0),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','ef9560b3-735c-4882-9aa1-dce8fe0697af',1),('207b5fc8-fd83-4b5e-aca7-4e095d4906b5','fd09aae5-67c6-4162-a03e-95f8cc07d053',0),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','282b6071-c081-4ebf-bb56-5241153f1811',0),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','6eed7521-b56d-4af5-8927-771c01723b21',0),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('2d1e69ab-66d5-4fe5-be9c-27b709a9f7c0','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('320d776e-9e8e-4263-abf6-ebe3d862d547','29003608-50c2-4693-953b-316f0c71fa25',1),('320d776e-9e8e-4263-abf6-ebe3d862d547','3d4aaffa-8399-4bff-9dac-000e43b45ca9',1),('320d776e-9e8e-4263-abf6-ebe3d862d547','65b407fe-a29a-4caa-90c1-060588438771',0),('320d776e-9e8e-4263-abf6-ebe3d862d547','74beed12-a014-4842-91e7-5334c22a5bc0',1),('320d776e-9e8e-4263-abf6-ebe3d862d547','7b474ad1-7a28-4e53-970d-e9a5a14bfbab',1),('320d776e-9e8e-4263-abf6-ebe3d862d547','9c93a802-d9c8-468e-8983-129fc287b26c',1),('320d776e-9e8e-4263-abf6-ebe3d862d547','bd4873a9-1ad5-40ba-bec3-5506355835e3',0),('320d776e-9e8e-4263-abf6-ebe3d862d547','d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70',0),('320d776e-9e8e-4263-abf6-ebe3d862d547','fafceb53-8f99-4a18-81bf-59c2b83bcd6c',0),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','29003608-50c2-4693-953b-316f0c71fa25',1),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','3d4aaffa-8399-4bff-9dac-000e43b45ca9',1),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','65b407fe-a29a-4caa-90c1-060588438771',0),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','74beed12-a014-4842-91e7-5334c22a5bc0',1),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','7b474ad1-7a28-4e53-970d-e9a5a14bfbab',1),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','9c93a802-d9c8-468e-8983-129fc287b26c',1),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','bd4873a9-1ad5-40ba-bec3-5506355835e3',0),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70',0),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','fafceb53-8f99-4a18-81bf-59c2b83bcd6c',0),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','29003608-50c2-4693-953b-316f0c71fa25',1),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','3d4aaffa-8399-4bff-9dac-000e43b45ca9',1),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','65b407fe-a29a-4caa-90c1-060588438771',0),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','74beed12-a014-4842-91e7-5334c22a5bc0',1),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','7b474ad1-7a28-4e53-970d-e9a5a14bfbab',1),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','9c93a802-d9c8-468e-8983-129fc287b26c',1),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','bd4873a9-1ad5-40ba-bec3-5506355835e3',0),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70',0),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','fafceb53-8f99-4a18-81bf-59c2b83bcd6c',0),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','282b6071-c081-4ebf-bb56-5241153f1811',0),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','6eed7521-b56d-4af5-8927-771c01723b21',0),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('5873dcb9-9412-4db8-b4d1-e57a1d76dbd9','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('657f6bd0-aa09-4703-99f7-e3f48e2de466','066554e5-ed72-4255-a59e-37bce592656c',0),('657f6bd0-aa09-4703-99f7-e3f48e2de466','15afd007-2320-4439-853b-2ddaa8b2ff71',0),('657f6bd0-aa09-4703-99f7-e3f48e2de466','39474b38-4002-4590-a90e-7e77c25785b5',1),('657f6bd0-aa09-4703-99f7-e3f48e2de466','3e744bf6-ec28-46cf-af95-5cf3f4ee8d58',1),('657f6bd0-aa09-4703-99f7-e3f48e2de466','90977ab5-207e-4b07-836d-e8d59805927b',1),('657f6bd0-aa09-4703-99f7-e3f48e2de466','c294c451-1979-493f-ad49-15220c9bd6f8',1),('657f6bd0-aa09-4703-99f7-e3f48e2de466','c36c9c18-f3c5-451b-bf26-2887df9d803d',0),('657f6bd0-aa09-4703-99f7-e3f48e2de466','ef9560b3-735c-4882-9aa1-dce8fe0697af',1),('657f6bd0-aa09-4703-99f7-e3f48e2de466','fd09aae5-67c6-4162-a03e-95f8cc07d053',0),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','282b6071-c081-4ebf-bb56-5241153f1811',0),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','6eed7521-b56d-4af5-8927-771c01723b21',0),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('95e70707-5b66-41b0-a127-590f786b2fba','29003608-50c2-4693-953b-316f0c71fa25',1),('95e70707-5b66-41b0-a127-590f786b2fba','3d4aaffa-8399-4bff-9dac-000e43b45ca9',1),('95e70707-5b66-41b0-a127-590f786b2fba','65b407fe-a29a-4caa-90c1-060588438771',0),('95e70707-5b66-41b0-a127-590f786b2fba','74beed12-a014-4842-91e7-5334c22a5bc0',1),('95e70707-5b66-41b0-a127-590f786b2fba','7b474ad1-7a28-4e53-970d-e9a5a14bfbab',1),('95e70707-5b66-41b0-a127-590f786b2fba','9c93a802-d9c8-468e-8983-129fc287b26c',1),('95e70707-5b66-41b0-a127-590f786b2fba','bd4873a9-1ad5-40ba-bec3-5506355835e3',0),('95e70707-5b66-41b0-a127-590f786b2fba','d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70',0),('95e70707-5b66-41b0-a127-590f786b2fba','fafceb53-8f99-4a18-81bf-59c2b83bcd6c',0),('9e5fcd3f-0e4b-4627-8292-011015df3de6','066554e5-ed72-4255-a59e-37bce592656c',0),('9e5fcd3f-0e4b-4627-8292-011015df3de6','15afd007-2320-4439-853b-2ddaa8b2ff71',0),('9e5fcd3f-0e4b-4627-8292-011015df3de6','39474b38-4002-4590-a90e-7e77c25785b5',1),('9e5fcd3f-0e4b-4627-8292-011015df3de6','3e744bf6-ec28-46cf-af95-5cf3f4ee8d58',1),('9e5fcd3f-0e4b-4627-8292-011015df3de6','90977ab5-207e-4b07-836d-e8d59805927b',1),('9e5fcd3f-0e4b-4627-8292-011015df3de6','c294c451-1979-493f-ad49-15220c9bd6f8',1),('9e5fcd3f-0e4b-4627-8292-011015df3de6','c36c9c18-f3c5-451b-bf26-2887df9d803d',0),('9e5fcd3f-0e4b-4627-8292-011015df3de6','ef9560b3-735c-4882-9aa1-dce8fe0697af',1),('9e5fcd3f-0e4b-4627-8292-011015df3de6','fd09aae5-67c6-4162-a03e-95f8cc07d053',0),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','282b6071-c081-4ebf-bb56-5241153f1811',0),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','6eed7521-b56d-4af5-8927-771c01723b21',0),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('a3e6a779-77c2-43c3-bc2a-3eaf4324385f','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','29003608-50c2-4693-953b-316f0c71fa25',1),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','3d4aaffa-8399-4bff-9dac-000e43b45ca9',1),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','65b407fe-a29a-4caa-90c1-060588438771',0),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','74beed12-a014-4842-91e7-5334c22a5bc0',1),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','7b474ad1-7a28-4e53-970d-e9a5a14bfbab',1),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','9c93a802-d9c8-468e-8983-129fc287b26c',1),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','bd4873a9-1ad5-40ba-bec3-5506355835e3',0),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70',0),('a6b5e3ba-b72b-4028-92be-acb75bd541c1','fafceb53-8f99-4a18-81bf-59c2b83bcd6c',0),('b1c791ca-b006-4f60-af37-20225d5876f5','066554e5-ed72-4255-a59e-37bce592656c',0),('b1c791ca-b006-4f60-af37-20225d5876f5','15afd007-2320-4439-853b-2ddaa8b2ff71',0),('b1c791ca-b006-4f60-af37-20225d5876f5','39474b38-4002-4590-a90e-7e77c25785b5',1),('b1c791ca-b006-4f60-af37-20225d5876f5','3e744bf6-ec28-46cf-af95-5cf3f4ee8d58',1),('b1c791ca-b006-4f60-af37-20225d5876f5','90977ab5-207e-4b07-836d-e8d59805927b',1),('b1c791ca-b006-4f60-af37-20225d5876f5','c294c451-1979-493f-ad49-15220c9bd6f8',1),('b1c791ca-b006-4f60-af37-20225d5876f5','c36c9c18-f3c5-451b-bf26-2887df9d803d',0),('b1c791ca-b006-4f60-af37-20225d5876f5','ef9560b3-735c-4882-9aa1-dce8fe0697af',1),('b1c791ca-b006-4f60-af37-20225d5876f5','fd09aae5-67c6-4162-a03e-95f8cc07d053',0),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','066554e5-ed72-4255-a59e-37bce592656c',0),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','15afd007-2320-4439-853b-2ddaa8b2ff71',0),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','39474b38-4002-4590-a90e-7e77c25785b5',1),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','3e744bf6-ec28-46cf-af95-5cf3f4ee8d58',1),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','90977ab5-207e-4b07-836d-e8d59805927b',1),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','c294c451-1979-493f-ad49-15220c9bd6f8',1),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','c36c9c18-f3c5-451b-bf26-2887df9d803d',0),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','ef9560b3-735c-4882-9aa1-dce8fe0697af',1),('b79d0e1c-6d44-4f01-b111-278fe3db31ee','fd09aae5-67c6-4162-a03e-95f8cc07d053',0),('c12ebc78-3392-401a-8328-5dbb4cddf222','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('c12ebc78-3392-401a-8328-5dbb4cddf222','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('c12ebc78-3392-401a-8328-5dbb4cddf222','282b6071-c081-4ebf-bb56-5241153f1811',0),('c12ebc78-3392-401a-8328-5dbb4cddf222','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('c12ebc78-3392-401a-8328-5dbb4cddf222','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('c12ebc78-3392-401a-8328-5dbb4cddf222','6eed7521-b56d-4af5-8927-771c01723b21',0),('c12ebc78-3392-401a-8328-5dbb4cddf222','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('c12ebc78-3392-401a-8328-5dbb4cddf222','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('c12ebc78-3392-401a-8328-5dbb4cddf222','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','282b6071-c081-4ebf-bb56-5241153f1811',0),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','6eed7521-b56d-4af5-8927-771c01723b21',0),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('ff653271-7bc7-410a-99e0-cfeff6180f7b','29003608-50c2-4693-953b-316f0c71fa25',1),('ff653271-7bc7-410a-99e0-cfeff6180f7b','3d4aaffa-8399-4bff-9dac-000e43b45ca9',1),('ff653271-7bc7-410a-99e0-cfeff6180f7b','65b407fe-a29a-4caa-90c1-060588438771',0),('ff653271-7bc7-410a-99e0-cfeff6180f7b','74beed12-a014-4842-91e7-5334c22a5bc0',1),('ff653271-7bc7-410a-99e0-cfeff6180f7b','7b474ad1-7a28-4e53-970d-e9a5a14bfbab',1),('ff653271-7bc7-410a-99e0-cfeff6180f7b','9c93a802-d9c8-468e-8983-129fc287b26c',1),('ff653271-7bc7-410a-99e0-cfeff6180f7b','bd4873a9-1ad5-40ba-bec3-5506355835e3',0),('ff653271-7bc7-410a-99e0-cfeff6180f7b','d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70',0),('ff653271-7bc7-410a-99e0-cfeff6180f7b','fafceb53-8f99-4a18-81bf-59c2b83bcd6c',0);
/*!40000 ALTER TABLE `CLIENT_SCOPE_CLIENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SCOPE_ROLE_MAPPING`
--

DROP TABLE IF EXISTS `CLIENT_SCOPE_ROLE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SCOPE_ROLE_MAPPING` (
  `SCOPE_ID` varchar(36) NOT NULL,
  `ROLE_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`SCOPE_ID`,`ROLE_ID`),
  KEY `IDX_CLSCOPE_ROLE` (`SCOPE_ID`),
  KEY `IDX_ROLE_CLSCOPE` (`ROLE_ID`),
  CONSTRAINT `FK_CL_SCOPE_RM_SCOPE` FOREIGN KEY (`SCOPE_ID`) REFERENCES `CLIENT_SCOPE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SCOPE_ROLE_MAPPING`
--

LOCK TABLES `CLIENT_SCOPE_ROLE_MAPPING` WRITE;
/*!40000 ALTER TABLE `CLIENT_SCOPE_ROLE_MAPPING` DISABLE KEYS */;
INSERT INTO `CLIENT_SCOPE_ROLE_MAPPING` (`SCOPE_ID`, `ROLE_ID`) VALUES ('282b6071-c081-4ebf-bb56-5241153f1811','4d890335-e9e1-45ab-8c8c-98c9e8db8bc4'),('bd4873a9-1ad5-40ba-bec3-5506355835e3','6bbfee03-d9ef-43c5-85b6-59627959df8d'),('c36c9c18-f3c5-451b-bf26-2887df9d803d','b197dcfd-a869-436d-8746-8d288ae56ef7');
/*!40000 ALTER TABLE `CLIENT_SCOPE_ROLE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SESSION`
--

DROP TABLE IF EXISTS `CLIENT_SESSION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SESSION` (
  `ID` varchar(36) NOT NULL,
  `CLIENT_ID` varchar(36) DEFAULT NULL,
  `REDIRECT_URI` varchar(255) DEFAULT NULL,
  `STATE` varchar(255) DEFAULT NULL,
  `TIMESTAMP` int DEFAULT NULL,
  `SESSION_ID` varchar(36) DEFAULT NULL,
  `AUTH_METHOD` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(255) DEFAULT NULL,
  `AUTH_USER_ID` varchar(36) DEFAULT NULL,
  `CURRENT_ACTION` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_CLIENT_SESSION_SESSION` (`SESSION_ID`),
  CONSTRAINT `FK_B4AO2VCVAT6UKAU74WBWTFQO1` FOREIGN KEY (`SESSION_ID`) REFERENCES `USER_SESSION` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SESSION`
--

LOCK TABLES `CLIENT_SESSION` WRITE;
/*!40000 ALTER TABLE `CLIENT_SESSION` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_SESSION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SESSION_AUTH_STATUS`
--

DROP TABLE IF EXISTS `CLIENT_SESSION_AUTH_STATUS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SESSION_AUTH_STATUS` (
  `AUTHENTICATOR` varchar(36) NOT NULL,
  `STATUS` int DEFAULT NULL,
  `CLIENT_SESSION` varchar(36) NOT NULL,
  PRIMARY KEY (`CLIENT_SESSION`,`AUTHENTICATOR`),
  CONSTRAINT `AUTH_STATUS_CONSTRAINT` FOREIGN KEY (`CLIENT_SESSION`) REFERENCES `CLIENT_SESSION` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SESSION_AUTH_STATUS`
--

LOCK TABLES `CLIENT_SESSION_AUTH_STATUS` WRITE;
/*!40000 ALTER TABLE `CLIENT_SESSION_AUTH_STATUS` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_SESSION_AUTH_STATUS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SESSION_NOTE`
--

DROP TABLE IF EXISTS `CLIENT_SESSION_NOTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SESSION_NOTE` (
  `NAME` varchar(255) NOT NULL,
  `VALUE` varchar(255) DEFAULT NULL,
  `CLIENT_SESSION` varchar(36) NOT NULL,
  PRIMARY KEY (`CLIENT_SESSION`,`NAME`),
  CONSTRAINT `FK5EDFB00FF51C2736` FOREIGN KEY (`CLIENT_SESSION`) REFERENCES `CLIENT_SESSION` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SESSION_NOTE`
--

LOCK TABLES `CLIENT_SESSION_NOTE` WRITE;
/*!40000 ALTER TABLE `CLIENT_SESSION_NOTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_SESSION_NOTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SESSION_PROT_MAPPER`
--

DROP TABLE IF EXISTS `CLIENT_SESSION_PROT_MAPPER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SESSION_PROT_MAPPER` (
  `PROTOCOL_MAPPER_ID` varchar(36) NOT NULL,
  `CLIENT_SESSION` varchar(36) NOT NULL,
  PRIMARY KEY (`CLIENT_SESSION`,`PROTOCOL_MAPPER_ID`),
  CONSTRAINT `FK_33A8SGQW18I532811V7O2DK89` FOREIGN KEY (`CLIENT_SESSION`) REFERENCES `CLIENT_SESSION` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SESSION_PROT_MAPPER`
--

LOCK TABLES `CLIENT_SESSION_PROT_MAPPER` WRITE;
/*!40000 ALTER TABLE `CLIENT_SESSION_PROT_MAPPER` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_SESSION_PROT_MAPPER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_SESSION_ROLE`
--

DROP TABLE IF EXISTS `CLIENT_SESSION_ROLE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SESSION_ROLE` (
  `ROLE_ID` varchar(255) NOT NULL,
  `CLIENT_SESSION` varchar(36) NOT NULL,
  PRIMARY KEY (`CLIENT_SESSION`,`ROLE_ID`),
  CONSTRAINT `FK_11B7SGQW18I532811V7O2DV76` FOREIGN KEY (`CLIENT_SESSION`) REFERENCES `CLIENT_SESSION` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_SESSION_ROLE`
--

LOCK TABLES `CLIENT_SESSION_ROLE` WRITE;
/*!40000 ALTER TABLE `CLIENT_SESSION_ROLE` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_SESSION_ROLE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CLIENT_USER_SESSION_NOTE`
--

DROP TABLE IF EXISTS `CLIENT_USER_SESSION_NOTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_USER_SESSION_NOTE` (
  `NAME` varchar(255) NOT NULL,
  `VALUE` text,
  `CLIENT_SESSION` varchar(36) NOT NULL,
  PRIMARY KEY (`CLIENT_SESSION`,`NAME`),
  CONSTRAINT `FK_CL_USR_SES_NOTE` FOREIGN KEY (`CLIENT_SESSION`) REFERENCES `CLIENT_SESSION` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CLIENT_USER_SESSION_NOTE`
--

LOCK TABLES `CLIENT_USER_SESSION_NOTE` WRITE;
/*!40000 ALTER TABLE `CLIENT_USER_SESSION_NOTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_USER_SESSION_NOTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `COMPONENT`
--

DROP TABLE IF EXISTS `COMPONENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `COMPONENT` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) DEFAULT NULL,
  `PARENT_ID` varchar(36) DEFAULT NULL,
  `PROVIDER_ID` varchar(36) DEFAULT NULL,
  `PROVIDER_TYPE` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  `SUB_TYPE` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_COMPONENT_REALM` (`REALM_ID`),
  KEY `IDX_COMPONENT_PROVIDER_TYPE` (`PROVIDER_TYPE`),
  CONSTRAINT `FK_COMPONENT_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `COMPONENT`
--

LOCK TABLES `COMPONENT` WRITE;
/*!40000 ALTER TABLE `COMPONENT` DISABLE KEYS */;
INSERT INTO `COMPONENT` (`ID`, `NAME`, `PARENT_ID`, `PROVIDER_ID`, `PROVIDER_TYPE`, `REALM_ID`, `SUB_TYPE`) VALUES ('049eba05-ff13-48e0-a1ea-c63bd8cd808a','rsa-enc-generated','dcc080c5-aede-4fd3-8b01-bd0928b730a2','rsa-enc-generated','org.keycloak.keys.KeyProvider','dcc080c5-aede-4fd3-8b01-bd0928b730a2',NULL),('109e27cd-2b7c-4f16-a5bc-be11f305748b','Allowed Protocol Mapper Types','40ae881c-f4e4-4b07-b097-a67d2bf515e6','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','anonymous'),('15a17a5b-0182-434b-9bea-1d4f9b590a87','Allowed Client Scopes','4327ba47-4116-44ea-9c4d-02907dca81e7','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','anonymous'),('27b6183a-7d2f-4572-b1a2-1c470a446ce9','rsa-generated','4327ba47-4116-44ea-9c4d-02907dca81e7','rsa-generated','org.keycloak.keys.KeyProvider','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL),('2c8ac756-c844-46ce-bc98-7dabf49cea79','rsa-enc-generated','40ae881c-f4e4-4b07-b097-a67d2bf515e6','rsa-enc-generated','org.keycloak.keys.KeyProvider','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL),('30254afd-b60c-439f-893d-71e5c6045845','Full Scope Disabled','4327ba47-4116-44ea-9c4d-02907dca81e7','scope','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','anonymous'),('3200da3d-a43f-47e2-89b9-dfd396feb0af','Allowed Client Scopes','dcc080c5-aede-4fd3-8b01-bd0928b730a2','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','authenticated'),('3293c342-c07d-4706-9af5-30ab7cc3ee2a','Trusted Hosts','40ae881c-f4e4-4b07-b097-a67d2bf515e6','trusted-hosts','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','anonymous'),('3bb59e96-8933-44f6-a3d5-dbd560f450c3','hmac-generated-hs512','dcc080c5-aede-4fd3-8b01-bd0928b730a2','hmac-generated','org.keycloak.keys.KeyProvider','dcc080c5-aede-4fd3-8b01-bd0928b730a2',NULL),('41ca9748-0e9d-426b-8920-cbabcb302931','Trusted Hosts','dcc080c5-aede-4fd3-8b01-bd0928b730a2','trusted-hosts','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','anonymous'),('4666d7df-5fbd-4d9d-a994-3bf556430900','Allowed Client Scopes','dcc080c5-aede-4fd3-8b01-bd0928b730a2','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','anonymous'),('4f7ea068-f73d-46b8-802f-019464083c96','aes-generated','40ae881c-f4e4-4b07-b097-a67d2bf515e6','aes-generated','org.keycloak.keys.KeyProvider','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL),('57666d8a-3a28-4467-8977-dab5a8405bf2','Allowed Protocol Mapper Types','dcc080c5-aede-4fd3-8b01-bd0928b730a2','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','anonymous'),('589b00b8-e403-41fa-beef-753843c27c28','Trusted Hosts','4327ba47-4116-44ea-9c4d-02907dca81e7','trusted-hosts','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','anonymous'),('5b29d81d-533a-4f90-a19a-f230c8212d6e','Max Clients Limit','4327ba47-4116-44ea-9c4d-02907dca81e7','max-clients','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','anonymous'),('68d14db1-396c-4d64-97cf-f006667d48fe','Allowed Protocol Mapper Types','40ae881c-f4e4-4b07-b097-a67d2bf515e6','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','authenticated'),('783dc8ee-496e-4a68-b354-7883c9504c0b','Consent Required','dcc080c5-aede-4fd3-8b01-bd0928b730a2','consent-required','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','anonymous'),('83773dcd-2802-4ae7-9fea-4b24fc6ba2d3','hmac-generated-hs512','4327ba47-4116-44ea-9c4d-02907dca81e7','hmac-generated','org.keycloak.keys.KeyProvider','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL),('8a1b4b19-5104-429a-b262-f5f59d2b5606','rsa-generated','dcc080c5-aede-4fd3-8b01-bd0928b730a2','rsa-generated','org.keycloak.keys.KeyProvider','dcc080c5-aede-4fd3-8b01-bd0928b730a2',NULL),('8ccc4751-d08c-4bab-93e5-f3af367c141e','Allowed Protocol Mapper Types','4327ba47-4116-44ea-9c4d-02907dca81e7','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','authenticated'),('9061afd2-4711-4908-b689-4141d0013d69','rsa-enc-generated','4327ba47-4116-44ea-9c4d-02907dca81e7','rsa-enc-generated','org.keycloak.keys.KeyProvider','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL),('94334b3b-e187-49b6-9ea3-618928ac1a54','Consent Required','40ae881c-f4e4-4b07-b097-a67d2bf515e6','consent-required','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','anonymous'),('97fd0a05-d85b-4468-966d-f7005f2a9d0d','Allowed Protocol Mapper Types','dcc080c5-aede-4fd3-8b01-bd0928b730a2','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','authenticated'),('a2d35538-96dc-440f-811a-9ab5af1120c7',NULL,'40ae881c-f4e4-4b07-b097-a67d2bf515e6','declarative-user-profile','org.keycloak.userprofile.UserProfileProvider','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL),('ae5d694e-8d58-4c7c-b26d-102cba5e0c7d','hmac-generated-hs512','40ae881c-f4e4-4b07-b097-a67d2bf515e6','hmac-generated','org.keycloak.keys.KeyProvider','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL),('b0590b57-ec62-4f1b-82df-7c7171efa612','Allowed Client Scopes','40ae881c-f4e4-4b07-b097-a67d2bf515e6','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','anonymous'),('b7798925-8e72-464d-b76e-951e5ea2a688','aes-generated','dcc080c5-aede-4fd3-8b01-bd0928b730a2','aes-generated','org.keycloak.keys.KeyProvider','dcc080c5-aede-4fd3-8b01-bd0928b730a2',NULL),('bc9cb75c-4512-4d5b-9284-27e900a79763','Allowed Client Scopes','4327ba47-4116-44ea-9c4d-02907dca81e7','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','authenticated'),('bf431d4f-712a-4417-bc9b-881c566afd0b','Allowed Client Scopes','40ae881c-f4e4-4b07-b097-a67d2bf515e6','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','authenticated'),('c88fa1ed-9cab-4a1f-9e71-9c9befd2c02c','Full Scope Disabled','dcc080c5-aede-4fd3-8b01-bd0928b730a2','scope','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','anonymous'),('d2bb4848-001e-4e99-bc50-aa183de69051','aes-generated','4327ba47-4116-44ea-9c4d-02907dca81e7','aes-generated','org.keycloak.keys.KeyProvider','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL),('d3bb1a07-b3bb-40ff-9281-69b7808af565','Consent Required','4327ba47-4116-44ea-9c4d-02907dca81e7','consent-required','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','anonymous'),('e61a547c-a280-4802-8c67-50ac2f0f859f','Full Scope Disabled','40ae881c-f4e4-4b07-b097-a67d2bf515e6','scope','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','anonymous'),('e85c5fac-bbe2-47e0-b99a-91b9fbadf966','Max Clients Limit','dcc080c5-aede-4fd3-8b01-bd0928b730a2','max-clients','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','anonymous'),('edd5efdf-b0fb-4644-a49b-4ff570dc37f0','Allowed Protocol Mapper Types','4327ba47-4116-44ea-9c4d-02907dca81e7','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','anonymous'),('effebba1-3ba0-4f1a-8ff6-4e3b32262881',NULL,'4327ba47-4116-44ea-9c4d-02907dca81e7','declarative-user-profile','org.keycloak.userprofile.UserProfileProvider','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL),('f9b42840-4ce9-483e-b1d1-8465f046266d','Max Clients Limit','40ae881c-f4e4-4b07-b097-a67d2bf515e6','max-clients','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','anonymous'),('fdd1d4a0-a99d-4ece-9fb5-8dca87d28f5e','rsa-generated','40ae881c-f4e4-4b07-b097-a67d2bf515e6','rsa-generated','org.keycloak.keys.KeyProvider','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL);
/*!40000 ALTER TABLE `COMPONENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `COMPONENT_CONFIG`
--

DROP TABLE IF EXISTS `COMPONENT_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `COMPONENT_CONFIG` (
  `ID` varchar(36) NOT NULL,
  `COMPONENT_ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `VALUE` longtext CHARACTER SET utf8 COLLATE utf8_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_COMPO_CONFIG_COMPO` (`COMPONENT_ID`),
  CONSTRAINT `FK_COMPONENT_CONFIG` FOREIGN KEY (`COMPONENT_ID`) REFERENCES `COMPONENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `COMPONENT_CONFIG`
--

LOCK TABLES `COMPONENT_CONFIG` WRITE;
/*!40000 ALTER TABLE `COMPONENT_CONFIG` DISABLE KEYS */;
INSERT INTO `COMPONENT_CONFIG` (`ID`, `COMPONENT_ID`, `NAME`, `VALUE`) VALUES ('03c17acd-4360-44a4-8bc6-aad45e210365','109e27cd-2b7c-4f16-a5bc-be11f305748b','allowed-protocol-mapper-types','oidc-address-mapper'),('0b8d3a4d-eac3-4540-85b0-e3da7cbe0fca','4f7ea068-f73d-46b8-802f-019464083c96','kid','9a119feb-2998-4925-b7fd-628c539635c8'),('0db6fe5a-81a6-42ea-a461-4dbed13394f7','15a17a5b-0182-434b-9bea-1d4f9b590a87','allow-default-scopes','true'),('0dd12a88-1755-4beb-b1e8-28856241d955','b0590b57-ec62-4f1b-82df-7c7171efa612','allow-default-scopes','true'),('10aa765f-9187-4f06-b0b9-24eaddda70b6','fdd1d4a0-a99d-4ece-9fb5-8dca87d28f5e','privateKey','MIIEowIBAAKCAQEAw8Xm9XriD6gspfXotV4cm0VgO8qIGHUsC2ma15EogBeil7RZ+HEalzkO4bfIjy1An7V6uWjpv26erqJVr1+4sil6j95y2luu06/upp02gFMR4ijvNQ3pIOzGeOIbAh2FQKvtXfZO1s9y+/uxoOHJmi6hYizLs7ZGBtYYV9v7JPXtyxUy3nAP1AK3AjrPF9HA5vmBdP6rpXw2+s/lthTjSLJRkykeOYdrcka8n6ykfSVv61GS80QfSEGE2UigcGLrT4PsOijcgM/P85vFwyPgjdtjJ7aJjJSu1v13B1sAvATYSM0EcQacg8pFdBtf2EsuM5vjsc552+/5R4z7cQCGbQIDAQABAoIBAAENcFIY3gbthYftjC8QGXaa5zXgad46dbPTVoTFgAKS80nTgLmoNKy5dCNchJSDJBoiX02IizEp8WpzmaL+uxR3dUmOuEQE7X9ag8hGMhkl2S5uapLOHKxo+XbZkwGY57LS+cy2hUf8XNw5Ry7b5+SPTvN4gC6FfVeOERrFJlLViTqS1oQYm4KLP46WRaPOxwp5XUUQuWKpgJKlIkSdvPvyOyoOPyvqy15PkTiy/0akS2qRMlqqG6NTrKxd0qK2vsZPJ1GcyOdApf7Hwz/EwH2eIQoAooutFKHs2AwlxZLtzuGh8tsUR3C2NKs8CdQZE1/+wvh7GTHm70xyDAGcM0ECgYEA6gXO6uIrGRZNg8X3QV3wE4Lzd0jwhA/s/G0E+05dix6D1lUM9fVMoELU+Q3nBPAtZ8hYLjpN0YFFlfkw3MvPfLFk5RVVozq019JP1iTPWLupAGaI2ojhI0OmguyTMAKcBFUQ0kIZX78hNPNOEqGAjO1RDXJ6aDju1pt56LS3whECgYEA1iiGo9M/rfj6owHG5o29s0NGaFDfkWKJus9l+8hY78QPpHMYHfmoiu1fNmsaptYue8OWopVEEh+kgCmaEEVb21LZycvdQ2s+olLhWVBXu5/at+lgoXe9+kEhsMfWyqUHmnlO2n0OxGduLT8jPClMWFKLhx2c9oKywrqY67/EYp0CgYBba71zoLr5Z+8MJU/8JzhcRvHZjZL35EjOK8Cgc/KzIE6ccklH5HX0vWb5jGbNVQ5H1sor9Pblezy1480k1DHQInSp0XXM+GghT7WEkIi3v0e8MlIQHRzma80mpEiznrFYN+sEWHIVJ1NPniTHvnO7mhHp0Ojkwij7iW1MbRzEIQKBgQCDqSiVYVNd/psASiBhL7T3l52allXMSNtJ+SXGtHZQ+aVwQb4K96jxuFt31kLtXPH47tcWH1RZHBEDJhBsfmepn/b7BBWE1FMOcovOYAF//Rf3R0g8HKS0TQSMbV+U9/6Drp6W3pmMj9YBhTaBTxZ7mkvjhoYwW0vRQeyvyxbOVQKBgBRLL7o36O/arYWAkIp8kStDrLCKHXw49kUH/IM+pLla9ju1amr/VZs7JV42wH30paBzlTEv9xYTDQGYAuFhMmccIozCRyzhSRHpaB2xqF2txvKl9K43o1HvF3fn82/Bt6Ir3PO+SjLJ4CzgQXWezojhnf0/RpXQPbnjSXXsbTEn'),('1ad2341a-911c-4c84-a778-f0fbae66c2ba','8ccc4751-d08c-4bab-93e5-f3af367c141e','allowed-protocol-mapper-types','oidc-full-name-mapper'),('1b88bdb1-951b-43c7-bf47-6c01b242985f','8ccc4751-d08c-4bab-93e5-f3af367c141e','allowed-protocol-mapper-types','oidc-address-mapper'),('2d9cf98d-2f74-4047-b1d6-bd13417d0831','27b6183a-7d2f-4572-b1a2-1c470a446ce9','keyUse','SIG'),('2fa16895-4b67-4f17-9802-1a72c8c23a5a','049eba05-ff13-48e0-a1ea-c63bd8cd808a','keyUse','ENC'),('31e4aefb-131b-49e1-9880-ca7f3a1b8f1b','049eba05-ff13-48e0-a1ea-c63bd8cd808a','certificate','MIICpzCCAY8CBgGPqY8YhTANBgkqhkiG9w0BAQsFADAXMRUwEwYDVQQDDAxBUkNCUkFJTl9ERVYwHhcNMjQwNTI0MDc0MjU4WhcNMzQwNTI0MDc0NDM4WjAXMRUwEwYDVQQDDAxBUkNCUkFJTl9ERVYwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7rQHTfxxPly3rDSwL0CymftfUm1njV7Go9jOUiNyiC70vx1OkILu10+TI3ofodLDcfBBnaryn0RaLSsl7loy+OxLHtq7WxgepyhVl9kCfy3aJdc8zqsMSBugmrw/IXaZ0BLDEhtGS2aW6Aeh7DXmQqhvKmwppJegc7MoQ2t+iB0m0ya//QNf6WxlUUJhXGCY2CNToSKTNvUtNiacxlqCWSELYwiAb2OBa+CqOfIaEyQnIn2KeBMANgR3t3dZNpNRqpNNwdEkeFypnisWUN8waQmhl712NRGFQuXOCoXGMDs4Mx3Qq/jogD3zjl+nZbb3O2Of5PHn+Kjg/lGNP9ckvAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAAbNjJTAE0hI5jw7YLFYmF6Kd3V2IhWhRmwK8wP7NEAJnQz0I3cRAZ01OthM0wgOEp1dvCzhnU26bidwq02BwpQrpFwGlEJ80WGusfPbR3UKvERDNFdJkPZ7BTIyoM0czpSb0PBZJuP/M2qw1DWqpg1SC7lSJqmmDyTLVcXnLetc5yq+uJ8ovLl4DxVV0HgyykCCbCzR/w0lCmXiyAITwMsGOhYKLr2zhVJd40hYbuTMmTJbfWsSncC8SnHSGCMB+ayzsJYlhUqb86jJujaxmVOyDiYPf6sgVY14HP/VsouGZ+LTWQqflyMf1vfE4jSx+kAs+eWJrGlTlnnWJVpvVz4='),('32b136df-0fe8-4323-ac4c-4bb36a65b51c','ae5d694e-8d58-4c7c-b26d-102cba5e0c7d','priority','100'),('37c49b14-c28c-4c8b-aab4-65b23fdff698','68d14db1-396c-4d64-97cf-f006667d48fe','allowed-protocol-mapper-types','saml-role-list-mapper'),('386a8722-c93a-4933-a958-989c7228557c','83773dcd-2802-4ae7-9fea-4b24fc6ba2d3','kid','6a6f88b0-7619-4dfa-971b-649d05ab4282'),('38b98ea3-24ef-438c-b774-a32b8126abab','57666d8a-3a28-4467-8977-dab5a8405bf2','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('3905e374-c91c-4bc2-8e7f-614b39d460d7','2c8ac756-c844-46ce-bc98-7dabf49cea79','privateKey','MIIEowIBAAKCAQEAsfkMWSmAYPdsv9a+osQk2TLpAmL/9WH7x+T7dFw0wFBYEaUJwF6e7KQ3D/ENLw1uUnnXoLUYJ1h2zJw8zL76IfmxZpaq1fBRYD//77sarFt48Fg3fd2MTsBpoGI7XKnFud/IjyGTmtLM53YvN+5x9gbKjAgs53gDDju8sVotYBnMfU7ufZeqLxe4RFDxIcII/okjf3Fm0kr+0nc04qfzmLhpUnLVrAzKXtMsnwrAdTxo38hTw3EJl9uvqOZx6J1fRJlxgVmZdlJ6Nb5nlu2xSF2NjT/mXjm79zbUIcpO0XjT6/w297rGbnu4BLwudT0XrfYFxlxcB/IIFjZog6E0QQIDAQABAoIBABXPHhtDAa9gQpxTlXLetGFFYtZfCVypkDATuFiu/+vdLJ2k6sf/EFyRVwoGEQaGD9HdinXwGzCfH5e9QZoZgQ9uyLV4myCyXVZ9IUDgXclnQC+7r7kl0A7Kd2cnAn09dLRtXudjRvI/CMsIaYriFmb+uA1m5xHKXB1/ZoUZNtxCuEILzNmdZO8xYn+jcp8jv/frAHYxT2GEY0FACEie/MJHHDhDMwXYVHV1KFmKQxI3ute1+Zj6Zm/dehMl0ArnBryVuE8O8EmarJXKn3TPSTPm5Uu2fgCqhDD+RE7NmHQXOp6idAZrwobVfa4yNHzTQ2isHmup8BsrtwL5Jh6j41ECgYEA9TKYUNVJf9msr2Y/Y2cdVse1D/BxmKobqI3Na7bRuP5F7D4mmmUwr0JthAczR1hMbhH4DDdU7bSA7ha/efJ3ugcv7vkDLZZOVfGRx437QYh64IGpHxK0pS9ITXIn2PWso7pnZBmF9QxKH5M8J3PzA9VB49YfXvjL+Aprdi+r/zUCgYEAudBGIIdu0j2AJY44Pk5Y28qsJ6qcjjfgIZaDSz0Cdhu15Dg9t/nmEtRup2Q20/pDHJQ39r3uzJpUCDidYjcDdBlIEynIuCuIZCHF4Rzo79rp36126T0ICsd4iqscebGacDS8GztK4Mr4naYyEUzYNRBX2OJsnoq9T43P0huJRl0CgYEA6cSM9zQB2QSc4LTo5kpe7GeGNmYUx2oW6IIZQQKRvH+gBnXStAyZd94rL172AqUqaR50kKEINYCME3JYp1kF0LQdfQangwT0NngTdl/lFjuaewTjSp6432vO0Rdu6ih/qbGD9SloT5Kh8KydAzhGjnb5VNDGI/Or8xVsEsesgakCgYBMKFYlzPyaBQT535Gjz97Rrv7ifyYNkE831QXZF5djqzXg5UA+oJkxDIqR5xwlw8Qv+Xv6kJxIldRtTi8LazrqIdaNrCmNeqI7UDBXdM7wSAxlViaPsCkUqe74/ur74dRHuwWCL8of2nENxGlu980B6sHmrd8RGBJBggE0v36DLQKBgA4RMpJqF/2zFwe6R16NQ7YL61jHIZ2undksOpCAn//F+bNd4ibxdtJPpJ0GaoDSs7Mw9gleODPxMoA5VwL+cwJRISM+5WEr8i/wPwBqc8+rKslrc6Cp7VJ9gLiD7PN5U8cbGVqKqvqtLi7BC922psQAAMTIf6AiRIsIZr7cw4s3'),('39da9567-8fa7-41e8-85da-9ffec7378869','4f7ea068-f73d-46b8-802f-019464083c96','priority','100'),('3a90ebb1-1bf9-4728-a2c5-e6c096ad8ae2','41ca9748-0e9d-426b-8920-cbabcb302931','client-uris-must-match','true'),('3aa942d4-2bcd-4382-b39b-c99b2dce9387','8ccc4751-d08c-4bab-93e5-f3af367c141e','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('3d14df9e-012a-42f1-abb0-a5456fcf9f4c','8a1b4b19-5104-429a-b262-f5f59d2b5606','priority','100'),('3e9d3f38-9ca5-458d-b7bf-d5eede131f43','8ccc4751-d08c-4bab-93e5-f3af367c141e','allowed-protocol-mapper-types','saml-user-property-mapper'),('41f0b1e2-5bd3-4107-ae22-ee536128b72f','68d14db1-396c-4d64-97cf-f006667d48fe','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('48c7b6c4-21b7-423d-be36-4fc875e59a29','fdd1d4a0-a99d-4ece-9fb5-8dca87d28f5e','priority','100'),('49470538-d452-40e3-a98e-e28fe2d2f1c7','2c8ac756-c844-46ce-bc98-7dabf49cea79','keyUse','ENC'),('4a74e256-ba2a-4118-a173-09d19645812f','57666d8a-3a28-4467-8977-dab5a8405bf2','allowed-protocol-mapper-types','oidc-full-name-mapper'),('4a9f563e-3e3f-41c8-95f5-44ed3909d858','8ccc4751-d08c-4bab-93e5-f3af367c141e','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('4b2c5602-b6a8-4633-a3b9-3b77c4c9507f','97fd0a05-d85b-4468-966d-f7005f2a9d0d','allowed-protocol-mapper-types','oidc-address-mapper'),('4b4c8d68-d082-4b2c-9c7b-59d0724c1d3e','049eba05-ff13-48e0-a1ea-c63bd8cd808a','privateKey','MIIEpAIBAAKCAQEAu60B038cT5ct6w0sC9Aspn7X1JtZ41exqPYzlIjcogu9L8dTpCC7tdPkyN6H6HSw3HwQZ2q8p9EWi0rJe5aMvjsSx7au1sYHqcoVZfZAn8t2iXXPM6rDEgboJq8PyF2mdASwxIbRktmlugHoew15kKobypsKaSXoHOzKENrfogdJtMmv/0DX+lsZVFCYVxgmNgjU6Eikzb1LTYmnMZaglkhC2MIgG9jgWvgqjnyGhMkJyJ9ingTADYEd7d3WTaTUaqTTcHRJHhcqZ4rFlDfMGkJoZe9djURhULlzgqFxjA7ODMd0Kv46IA9845fp2W29ztjn+Tx5/io4P5RjT/XJLwIDAQABAoIBAAE09/yf654Z72r8kovrbGNNdIVkNszVzaU5dBnEoTOlJFVMCvdUpyrP2i8U8eH59wDFBwtIbh0RuxcZ7JCxZk++Febw7spty0lM3u+0jjQj/81CMlK6JmwrzIwtZijd5rPDeW/NPQWRgC60+qbdq7AB2GSPcUAnSn3y3eGJTWg3LfaC8jvEo3F8PxUXHS1qdvjAmYWydpwiHQcoB/IssOp2bVEmYDuao8rRKXfja35zc+jqUSWBdCUGvbGYuQTH4tes15CdZbtK4/aui25n0rG2E6OWgMHeIRppmHoA+EG0vyTimkENTL6Zahj+oHS2RWHA0BW4wKCbYW3CX9mIS3kCgYEA32k3BD4JmG35XApYveIgLfF4qGA6vMzAeo3NqIcnRYYXdqtMo4vSn3JfL7UW5LghEUKe+5p+FgE4bJ9b5UxmtaxMoKwHaxUL3MNjPapKGJEyZGPGsZAtOVWLyXmYB+m3uRlCAiPEkYwnq85rIP1t2GjIn3oPW/a59VcBS7vkDHcCgYEA1w1XZ5FBPBD8vCWaGHzfA9/+BAYg62Dgbtf+xmGUmqgCIkdJ/tN/qv7BqXbdfy0EV9lo/GNAkx+efsqoBRZLxzLX5pAGgVCys9QK/ip6KlgKYDI+VETcabBbxsZA4syvifUdTr5odfj4ag/WAbzUpXx/wA7MFpZ3nwF3WUn0rwkCgYEAoUbE7o0SRDVvg1/8u+aXMFNWtMXy4QQ2FsJzKiuWz/uCyKnUQ2PWgkAAMuJZSncZd+pN9neKebwbzV4k6pyCsLdXAc3t9QFWdOGfrI5XuvBmHk5gyyG3Y+I7bRAYDe5MMJTpL74+UouIv9/dOg141HagDXAB0nODvnY2e0OB1vUCgYBOFJ7+kIPB/lz1JyXq7DPA4WwGI/+B1rvGIxBzEOz2tjeIIKAiMMJy4GqKUAkd6sdf6iUvNg41HM7cNFKK/kxnN/Oh1/s0qosntb6ECAIxK6Qgxz1QNWxdx1WbN1JJxo2ZpnCMrZ+Z7dqsD7HhEaGXsDtmDTlWni0yg0LFHZAEmQKBgQCFk6quWPwNPoneqxqJmwBJcS398qZIaKpMuLP0LsyS34mVw8fp3PfD1nm50Pix6eaCbnAylSpQniTBJb6Tq73XxX6Qo06tcp4wZM0BTXIRWAUPLVsnKq3Z6RB3Nfzqt4HoJMlpbw5aEkDtuBoGEASGkM6T7DgY7xPrXbu/XuVsQA=='),('4dca65be-dcb8-481e-bafa-cbaf86d2d8e6','97fd0a05-d85b-4468-966d-f7005f2a9d0d','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('52723bfb-4af6-49a5-ba7e-611671aa5f76','109e27cd-2b7c-4f16-a5bc-be11f305748b','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('530670b4-86b8-43cc-b0db-812725713be3','4666d7df-5fbd-4d9d-a994-3bf556430900','allow-default-scopes','true'),('54ea44f5-2f36-42e6-aba5-4857835d3798','e85c5fac-bbe2-47e0-b99a-91b9fbadf966','max-clients','200'),('557ac210-a1ac-44ec-993c-92102705e1cb','a2d35538-96dc-440f-811a-9ab5af1120c7','kc.user.profile.config','{\"attributes\":[{\"name\":\"username\",\"displayName\":\"${username}\",\"validations\":{\"length\":{\"min\":3,\"max\":255},\"username-prohibited-characters\":{},\"up-username-not-idn-homograph\":{}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"email\",\"displayName\":\"${email}\",\"validations\":{\"email\":{},\"length\":{\"max\":255}},\"annotations\":{},\"required\":{\"roles\":[\"user\"]},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"firstName\",\"displayName\":\"${firstName}\",\"validations\":{\"length\":{\"max\":255}},\"annotations\":{},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"lastName\",\"displayName\":\"${lastName}\",\"validations\":{},\"annotations\":{},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"test_role\",\"displayName\":\"test_role\",\"validations\":{},\"annotations\":{},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false}],\"groups\":[{\"name\":\"user-metadata\",\"displayHeader\":\"User metadata\",\"displayDescription\":\"Attributes, which refer to user metadata\"}]}'),('5b21fea8-be9d-4049-b24a-81ceec1fb95c','68d14db1-396c-4d64-97cf-f006667d48fe','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('5d38795a-feee-4939-a33b-9a37d8b7e18d','589b00b8-e403-41fa-beef-753843c27c28','client-uris-must-match','true'),('5fdf0f0b-9069-4a8c-af8c-6aac4bc40182','68d14db1-396c-4d64-97cf-f006667d48fe','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('66da32b9-8b0e-4a5e-bd70-bf8a3aecaede','ae5d694e-8d58-4c7c-b26d-102cba5e0c7d','algorithm','HS512'),('6807d96d-3654-488d-9073-5bfa31e80ee5','049eba05-ff13-48e0-a1ea-c63bd8cd808a','algorithm','RSA-OAEP'),('6d6b405f-0794-4cf2-b288-7355f2d49fa2','57666d8a-3a28-4467-8977-dab5a8405bf2','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('6d73c1b1-cf7e-40da-b7b1-065ec4a45c5d','97fd0a05-d85b-4468-966d-f7005f2a9d0d','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('6d80f27e-9f45-4a80-bcd8-af6a456b9fc3','109e27cd-2b7c-4f16-a5bc-be11f305748b','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('6f55df71-1c5b-4903-9404-baef246cd949','edd5efdf-b0fb-4644-a49b-4ff570dc37f0','allowed-protocol-mapper-types','oidc-full-name-mapper'),('6fdbd747-e7d4-42f3-869a-79501dc58dff','fdd1d4a0-a99d-4ece-9fb5-8dca87d28f5e','keyUse','SIG'),('738db0c6-300e-47cd-94b3-b77f5ac490a1','3293c342-c07d-4706-9af5-30ab7cc3ee2a','host-sending-registration-request-must-match','true'),('7e4e47af-0996-40de-837c-87f3fa1bc257','9061afd2-4711-4908-b689-4141d0013d69','privateKey','MIIEpQIBAAKCAQEArDluze87aqXIvF7lxK03wrCz5ACEQ7ULqx7GwhvNV4Xc+QwBlT9qotGhQlneKfUaqex7s1du6RsM3fsm9/AUrVM6pPd7jl01xsckPWL2wnJBxMSex2ABXiHag3MmrSd8DD1WRvahkdFmRhSrmLbeDhe4CVDqRi9Mfmg9ZyVfOqaJWZRa+rTRKPC5JIJGWQaad1iDDWOkPc32Iuv/xwJBrYcHnR71nGyOgs0k5WLAyXftDHjoSF+Srqtgyf9ZAYDItC5LCTxxFqHR/6UL/ahKXA1fB/pQ6yVkNEGou2J+hEwMp4TMPXZRWbmC2Ocej/vP++ctoNL2pwSg6cYxpNNr0wIDAQABAoIBAC/Nfy5k1SKcYnO6iV8GT0BYSI4kAJJEebklATkWe2/sJXHw/a193SzVL8PZattNf3mjvJACWDQWgINKtz3BYxPa0p/AW4if4ZHLa4koY2KEUTH6+zFOugJ1NhAfUaBlIb7J9Z0bzqvJKfCQwzJPq6HQvpHDZktVzI8XgteS3XRLGTHaUYN9tUWqO76nAHSMNO5ogDPGAjWm3Q0xGVLnTrdO08GBVYC7eWcowP+TzmyWj396joZqjEu7w663R2frG7czTrXxDE40CnzwSRzOY2vxaalwJo4AQMFd980WmLyRPl5El5GzJsFJKGEx1Vf74l4rQln7YxXDn6QjNX3Zt/kCgYEA5cs7F4Bl8NQcKgw78+4pIxftLjtBTDDolBklchZISfwS9ApgiDS4f1KzCb4XZqRtPT6Sd7vIHscHLpv0iARJb9v8LzU6NeGuQopWh2IQy0A4H8y3jdapmlN9j1wk3pIKPufJi6X9e/pY9X8ofX44trY7P7rdpPqcyCGgqN0M6vkCgYEAv915u4aI4ChAzD6tXW+l0iMA7cIjDD/atmnsEMo5sejA9wVZLVJMWCm7Sss+Z5mpAopiAMP4DHu/bAV0mLF3cSFyY8DQ76ahZ6jGoC0ROFuQ8NwHkBjwwUqmdVqs5cv2xBD9utBjZohEr0pU1vyXuhnRRb9RFToCECnd1g5QlCsCgYEAhkVoM+VMq5nlExSu+2uQEfdrGXZ6QyDY5aCD5tPqTYDDOmHN4gicPZl/EBRI7CrrwbuMLfZFiuZU5gEus0N9/aSXCKabatyBzFO1F8pPcdQGd15RasmhtJFQbrIywPKolfUuKTFGmmjADkLMz/cvAGQe3rA7zqi5b59mlwUDBVECgYEAplVrY5eOskYzZt4vjgFs0HBoLbdQix4QjnPlo+ite/88yupoFVJzvTrwlDb/MvGUBA/URrWeJbRij1NrcbTox7snYDOY7yqWYzd+ev6XHzTCwwz9wX7dubTt+m9R1SauF+xjC8H9arZYVV8rjfYN09juBHmH4c8YWsrw9tISKWECgYEAkKMSjL2J/uzGA1bPNAXshPgnOMk04qerHw2RbN6rpfTaZXlG7SlJ+prHTdqjRkkEJCFbhOqhJ0LkQeHgRX4nmGUs12uExlxoB19avBPwFlaqi6dDORg3cX7TOSXPTsIJa0PQ8lmut+0CKU2f1sK5RLnSqpV6hsNmoMy9NSpEbNQ='),('7f809e84-f1e5-46a1-a976-8f2423a80576','109e27cd-2b7c-4f16-a5bc-be11f305748b','allowed-protocol-mapper-types','oidc-full-name-mapper'),('82288863-446f-49c9-8bb2-5e1981a1f019','109e27cd-2b7c-4f16-a5bc-be11f305748b','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('85103c60-f751-43ea-a705-1c45585f0995','97fd0a05-d85b-4468-966d-f7005f2a9d0d','allowed-protocol-mapper-types','saml-user-property-mapper'),('85c186a2-635d-4d80-a933-999a8841c9f1','d2bb4848-001e-4e99-bc50-aa183de69051','secret','LlLTm0b-Xy0JTxffxddtaQ'),('85c628d8-59a4-4141-8491-38e4f5bfc26b','ae5d694e-8d58-4c7c-b26d-102cba5e0c7d','secret','3rtTVrYNsWATt1HG2XtMIkQeag1ejovOuwmwxjf3LxPuLHPUY5Q8qOXi6UunCuRVi7v4ZoupDNTOnYjYXsJkxwLYbVg5V7mBZXkLVmOmeC8BOGPE-ON9_WTaoY-PpVY9Q-0LgMoXnvFHM96Wn9WPzbmraeTzjctQT1bfbo7J5t0'),('86300286-42c5-4b02-86cd-629b81ae5b29','83773dcd-2802-4ae7-9fea-4b24fc6ba2d3','secret','8qGKdOovsPdn-YpffllKxImg5URqmM-3AgyOJ0XJ8TY4Ql7sBs1yK4S7W8ctQhzZ9qJHYTsgb8XUqFCdkGlK-m1DHPbjYCXGDIdXRj8bROWarsRIrx_I-5IWKBOHbfaZyI-KBiI6zPn7qdxrtWRb1s20QyOZRg3VahpcnrGsn68'),('8707df67-7730-4054-976e-559b6a93f0bf','68d14db1-396c-4d64-97cf-f006667d48fe','allowed-protocol-mapper-types','oidc-full-name-mapper'),('870ab3c4-aad8-49a1-92ca-d337ca65d79c','ae5d694e-8d58-4c7c-b26d-102cba5e0c7d','kid','ab7a75b1-bcae-4d17-abf7-744888565203'),('88f41479-d005-49b1-89e5-c5fd62724210','edd5efdf-b0fb-4644-a49b-4ff570dc37f0','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('8d32cae9-f3aa-4b6d-85fe-f622f450a31d','8a1b4b19-5104-429a-b262-f5f59d2b5606','certificate','MIICpzCCAY8CBgGPqY8XLDANBgkqhkiG9w0BAQsFADAXMRUwEwYDVQQDDAxBUkNCUkFJTl9ERVYwHhcNMjQwNTI0MDc0MjU4WhcNMzQwNTI0MDc0NDM4WjAXMRUwEwYDVQQDDAxBUkNCUkFJTl9ERVYwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDtoa6EHSqxZkhdGuJOeVD09UVQf3SpZiaYS5ebNKQ4D1hULaZyJYWv9XEbFnR5iieVbsQlQnEd9KDhgZpH8PXXQjeSZ++Id2G1IJE/XsJLR0hj3V4nAVRHb6i0+9HrAFZVtiDg1R04IQvD6fUr+XTX4TpXiT3LeG/r9MhWR8U6vfid0E0364rjofrviTfN4cuU8tC0jd6xPsn+ZM0g35+G1moxq1Q6+qXaZmtDYarUY92XS700YGVAcDnHFDpqWRaq8ZJjyVyqMwGN/fUq/VfXU79xeTHWw6VRutmQzuAMgZCmRJXysDblidWNwQaRkzVn37gDqWBR2WX7I2H2KxcfAgMBAAEwDQYJKoZIhvcNAQELBQADggEBANoZTRH4C59hLDI6Wwax1TzPhrnT3T6qz4Xe8WQbcFY4oXzLXUm2qymMtt5gK8rNuE9/gSrVNNA1QXfPTOLkXuKy1ezNMVelVD13rXzzYSRkcnEiCqPjl/E41sZi+fvCksCDGoi+7hGH38N+EKi8PAkS7FvBIOTcI572yymq1Pf1ilw8gvLz4mgov00Ehyt6BFcejKlooD5CeedGd2WlXPm7Il4gUOY6z+JlcWK9FFdgapeeXGZ5rU9PMl73PBHkQo6ge6UNC9nia2cJQi7AOB19zO7enBGAI3WQCRzYk+/8B72bS4M1uxialjrZUqsgQDeHWBN9K6wAWNuMBd4ps/w='),('8dded0c9-0f3c-4aea-a626-1fadab865a96','3bb59e96-8933-44f6-a3d5-dbd560f450c3','secret','hlAsHH0THLihoe1vBgxo_F_taxZKaK_uKwpWLwzkudo3YE2bpNmf4PHbAZVTBnorzpvaZxcr051qSV8hIuOhXY4fWoAR4EW-q-iB-5CjR3_Ve_l6ZSQnL_9lvqiFbIiUzAUqKYx8VFVz3Xh0YPQcucOVB8xGdMA6VWTMw8h706w'),('8f0f5820-7d48-4238-88d9-ecdbd7c9019a','9061afd2-4711-4908-b689-4141d0013d69','certificate','MIICmzCCAYMCBgGO7xoIMzANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjQwNDE4MDI0NTQ0WhcNMzQwNDE4MDI0NzI0WjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCsOW7N7ztqpci8XuXErTfCsLPkAIRDtQurHsbCG81Xhdz5DAGVP2qi0aFCWd4p9Rqp7HuzV27pGwzd+yb38BStUzqk93uOXTXGxyQ9YvbCckHExJ7HYAFeIdqDcyatJ3wMPVZG9qGR0WZGFKuYtt4OF7gJUOpGL0x+aD1nJV86polZlFr6tNEo8LkkgkZZBpp3WIMNY6Q9zfYi6//HAkGthwedHvWcbI6CzSTlYsDJd+0MeOhIX5Kuq2DJ/1kBgMi0LksJPHEWodH/pQv9qEpcDV8H+lDrJWQ0Qai7Yn6ETAynhMw9dlFZuYLY5x6P+8/75y2g0vanBKDpxjGk02vTAgMBAAEwDQYJKoZIhvcNAQELBQADggEBACQBK0hmTpNThoDmm+mqOKLzO4YXVyw3DHVVbbL39vt+3M/Flw1YyiQa6UGVR+VhT/vsEseTDFGOt6FJnK4rdG6CSXYrFYXoK1kY90jYMB1EXJFrYfN/08GvEwKwIOaBvsK55HugDjG4vayDUndiuC4XKb6IgjYAycH/CsPNSplKrrUcX+ACe4deXxa8STlGAFHlHEYMgwxYZKwdg7i1SMB9tq/mhFX9DnjW5dUlGyUztTN9w7lQ5I82GXxoPY9+dyb7N2uX0grN9VWRAloZr+uelEJVcUwU2huFB0EepDE/9uHODP0OLxggp4VTv8SpMcBJ6FzeUErPCa0dZaTF9zo='),('8f91694b-253d-423f-a01b-2fd575f99911','109e27cd-2b7c-4f16-a5bc-be11f305748b','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('90344f49-253c-43df-93bb-5a515cc201f4','57666d8a-3a28-4467-8977-dab5a8405bf2','allowed-protocol-mapper-types','saml-user-property-mapper'),('91b52e03-41f8-45b4-b0c0-86e8fba70701','8a1b4b19-5104-429a-b262-f5f59d2b5606','privateKey','MIIEogIBAAKCAQEA7aGuhB0qsWZIXRriTnlQ9PVFUH90qWYmmEuXmzSkOA9YVC2mciWFr/VxGxZ0eYonlW7EJUJxHfSg4YGaR/D110I3kmfviHdhtSCRP17CS0dIY91eJwFUR2+otPvR6wBWVbYg4NUdOCELw+n1K/l01+E6V4k9y3hv6/TIVkfFOr34ndBNN+uK46H674k3zeHLlPLQtI3esT7J/mTNIN+fhtZqMatUOvql2mZrQ2Gq1GPdl0u9NGBlQHA5xxQ6alkWqvGSY8lcqjMBjf31Kv1X11O/cXkx1sOlUbrZkM7gDIGQpkSV8rA25YnVjcEGkZM1Z9+4A6lgUdll+yNh9isXHwIDAQABAoIBAAMAk+/fXna3UBq9CDbtmXQ+1Q1bTuIFBF/QNK4lUB3VM8x8dyniB3pdo3wwIHCYqj9dEuVVnZA/UvkGZiiahsCdeY5D+ebDe1yKeMtQKyxATk3UsifUAHiWlI8Uz2YkrvKORoQQrTnUUQ28mvhfQnanxdDtTvaPIvrcMNvGNFu1nH7+j32KmCGffUQekyQTtlmHEuqoILCeNHt6iwR+ee4kfrQzDU2d1LsGxv+I7qCTMye+OkZkjaQQSzssyFdIUwQFBj62y3yRtRNfeUK9XZ1lEMoNhIgxotQ0LGs5vFTVJUBl8h/sPYrlUmCjV4EWuk+vnHIPVGc0GOuEtVOD3EECgYEA+CTsSIxI2GNbdABk1MtgoqRQvEdrXJSn6QNMYzA0QwntFs9jgusTO6/cH3/QR9x36EI5qfXzFOI4z7THzd+hwEctLAKPnoLgwG5IB+tOC6+K/x4OA162Dc/fnd4E9p/QeZEu57MiwML93zalmu8nsGMSuFvdfyR4nf0iuEg/lqkCgYEA9SePJ0RX+inSno19qK7DMTavi9pC3WXlsDp4GxaLsaFDWB235OST1EIwoGHjk/iX8fLxS/DgUWn8RSJfH9hsbItUZTlURSbzwWiBye+5FOLvpqCeUwqNRo1LlO18IsGRhyVo/0swQTkooaNAeleG6a7fgjsI9C51xGaes4EIBIcCgYBYaheCdfGSoDw6pdVHeLuS18cofq4DS7hULuetw8QrSsMeSIClagrtTmi5FHpq0GQs7kPLiiW1gDFL2JcEhLUGZJX0w2jWyr0381NY9v9U7JQy+Et36ewmGbYMcsoD8cd/OTqkNdChLBj2ymrRPHtAvtwxshHGtoD/ke0oCh4WGQKBgHZXRk9aOm3USounGYWX9zmGgG8uSSC+04Wic9//nid9BRlAH0kq3gGUggyptEX1tsNg1wIloOMsGk71LJ4exxgOmgrTBc7r4rb2T9C/G8TtV2pEFqlXaqMoCdg2W7obXriyYmD6VqzlfquoNrPYFJQxiks/XC8jTk4ki29zVf95AoGAVu82z9puDRxafMC73zmC5GB6cJT75yA4CLAzZ+Lk5uKEX/XcbzZQo29ZPZwHYOKhkC8xiObNWMW5+IgyqVVCEnAuDBCqZSI1N7Mduhrw/HM31C6OkRxQcKO7IC+V6QeFbOQWaz78JqzS+kAsleAUD9dvnRcLaYdf0mD3QZuRagU='),('95964240-dcda-4fbb-8b2f-605b2cf6be58','2c8ac756-c844-46ce-bc98-7dabf49cea79','priority','100'),('95980a63-41f4-4dba-8278-008f7cf7227d','27b6183a-7d2f-4572-b1a2-1c470a446ce9','certificate','MIICmzCCAYMCBgGO7xoGyjANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjQwNDE4MDI0NTQ0WhcNMzQwNDE4MDI0NzI0WjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC6Gf5n37L/cwmaxP0wTVtE4ObVeu2inkufG2xB/5FefvRjq3OeCMPJKFoqnkXis9COgc86vd+qb9Cz4N5ZYl/2Cp8hoVkCqbYLLjCLst2G2VxZnVcyJnEOm49FCO/OmZSGb/IF9suA8M6ySF6jIkEBxwGbrxW/oSQG4AuDoMgSmV6XC0C8MElKnDpbqgl2t2FQOzdFAsSXfcSpyAnQIynEN2Whj/Oyfdsldum4KqqW3VSoFAQjfyVuI0oFotfaodPp73JJAZXsQQA789bDLrYCYO2xbyHwTumGLxqsCPiDtS2jR12OzyFFNpVbp+dRSexxbdZ3OJsiXCipHrZJvJmjAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAAH72BgrDdoAF7ALl/h+OsUhUcvHk8o2a/nU2/qwW85XvBWB00lZegJu13xqBJGIg1siUD1lCPE5f7X5CxeKOPhedS6CvljotjMRi8we6+2Jz0Y/eqyrD9wWLVlKxWsZulS44WCloiAqouKV6wSUyAL99v16hqwfzFh62bJO2f1akv3nsPLOq5dBaeOf47Dtz7PECeIm+MiWxvPrT8064SAQZXY+qqgCWzp7bcmzzt1tiAhRHFo509wjvXSUY4nzu30NFX7agDE2IqmaOMoPUwUuMgwaMWaykwfdXaDmqnbtzYhkAapjT7u+2vDu8vTc2OXVIqFDmohuyI8rhLHJPFE='),('9f2e6466-cef7-47ab-8999-5af9608ed36f','27b6183a-7d2f-4572-b1a2-1c470a446ce9','privateKey','MIIEogIBAAKCAQEAuhn+Z9+y/3MJmsT9ME1bRODm1Xrtop5LnxtsQf+RXn70Y6tzngjDyShaKp5F4rPQjoHPOr3fqm/Qs+DeWWJf9gqfIaFZAqm2Cy4wi7LdhtlcWZ1XMiZxDpuPRQjvzpmUhm/yBfbLgPDOskheoyJBAccBm68Vv6EkBuALg6DIEplelwtAvDBJSpw6W6oJdrdhUDs3RQLEl33EqcgJ0CMpxDdloY/zsn3bJXbpuCqqlt1UqBQEI38lbiNKBaLX2qHT6e9ySQGV7EEAO/PWwy62AmDtsW8h8E7phi8arAj4g7Uto0ddjs8hRTaVW6fnUUnscW3WdzibIlwoqR62SbyZowIDAQABAoH/H0I+eTQZ3M51zolH/dUAf5FHNJRxVd9A1HjGtk83G6g3ZIZn5SOTvzYH7yCmnG2Xmlgzdqpt5zLg21XqAxmDBNDD6eIwEUJwUXBwtwZUL1JTMLo86y+JzepTJkVaAlnJmMs028uYYxAPp+KAAua+BI5otr2x6X4njwAS7wkSkgzna23kOz+QHG6EeKU6TE+oCBE6LZ78BuLdbWZE8mXIbXXitTKQQ2P7CmHBpAmlcyNnRWjdzqEmBHo8mhq9jNMjvhMAh9RBWGQyqoy4mtnmo8NygzVFoDzVHkJ15X4r/qOhjxvVyMP/nhOHX9sM0FTkX0u1hE/3FYTrIfgl5x1xAoGBAOWJY5eHN3Svx98rMdkCaR3xvGEba3jCDOg6h+9s5wKc9Xn5rzfaIkPA4E11L+qg3Q+/hofaPoNnw0O6VXOG3+VpgqfVQgrAJNdzn5faUEe8BYfIn7I7L6LHPrwzglhJEQwe91S9p5KTcJPG05UXoZMmMX8X4SO+1U+WbDnRsxMJAoGBAM+Opc7yprSyrOSyefs1kAXLTasGxfSgoSICg2W4kRxJh399tXGpANZQkd2dgfiYBxjL5fgEz8fY+U1DUU6R4YveaFuwSb+3JNUXOiA4z4Mqlmo4QdXnq25NRjnSGw5xWPaplSFEpd/DMExLGTxhD08mZYgduc5UvGoHdUemGVZLAoGBAK1yHdF4mY0Q7uVSDH7A9ZFtAz4VrCY56+rnn6RqFD3DQTMEW0THDjFIY7XkIWsYVzL5NL0fSzGcjM521O9RMYp3KgWMfjUFcFFly2jBzPwHtd4e2Z2iX6KPEHbCLXJs8/bGx3o/PYri4qSoD2WPz8YTjD9PWg7auvFC51DhhoGRAoGBAJVhBeaGbwJ34cUKyo3Iw+nXowNm8YuZG779bhIEnrNsHA/eqnqt9oNII2MLXCJNUDZBcTZqSBx+BRbdTyQsC9b9b8AlhT5skh6nA9dy665sNSsnaAKKJLBS/yrYE072tt93t3PlEziTIyyVlAkHldwSN1qagkKEa8InsfmqCmw3AoGAcqkELawft4GWpXvWY3PjXDa+34laS1V3L+sV+AIKamSnDOj/otMUUJBcfb2Yo4FmAYzgQWr4hY22AfgVfENDzVxQDb9vdcmiD/2C+gFyZJGb/WN3Wem9bo8bmK8pkITc0xkY7nuqrnYjAxu94A/8PzSOPXT+hKy+GvdEa5560+E='),('a0ce5bea-7a12-45fe-9640-4d0192baa724','d2bb4848-001e-4e99-bc50-aa183de69051','kid','4916b317-ba0d-44f1-b180-27d57862c9cd'),('a0de10ca-aec5-4eb4-8f90-b194ce91c983','57666d8a-3a28-4467-8977-dab5a8405bf2','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('a5cd2e1c-3ba6-4abf-a2e3-7095ec3dd19f','57666d8a-3a28-4467-8977-dab5a8405bf2','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('a5d594cb-6f97-4677-8ff5-d9e606757de2','97fd0a05-d85b-4468-966d-f7005f2a9d0d','allowed-protocol-mapper-types','oidc-full-name-mapper'),('a686e847-0d02-472b-a3cb-eea31df8181e','3bb59e96-8933-44f6-a3d5-dbd560f450c3','kid','3d1da024-7a97-4a9b-8be0-3c66344089ef'),('a992a812-4e92-4548-9587-5e0a3154ef19','effebba1-3ba0-4f1a-8ff6-4e3b32262881','kc.user.profile.config','{\"attributes\":[{\"name\":\"username\",\"displayName\":\"${username}\",\"validations\":{\"length\":{\"min\":3,\"max\":255},\"username-prohibited-characters\":{},\"up-username-not-idn-homograph\":{}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"email\",\"displayName\":\"${email}\",\"validations\":{\"email\":{},\"length\":{\"max\":255}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"firstName\",\"displayName\":\"${firstName}\",\"validations\":{\"length\":{\"max\":255},\"person-name-prohibited-characters\":{}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"lastName\",\"displayName\":\"${lastName}\",\"validations\":{\"length\":{\"max\":255},\"person-name-prohibited-characters\":{}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false}],\"groups\":[{\"name\":\"user-metadata\",\"displayHeader\":\"User metadata\",\"displayDescription\":\"Attributes, which refer to user metadata\"}]}'),('ae7d14c3-452a-4a18-abbf-5137f644ec3d','edd5efdf-b0fb-4644-a49b-4ff570dc37f0','allowed-protocol-mapper-types','saml-user-property-mapper'),('b00d367c-7d84-48ec-8535-f2d3765d8524','edd5efdf-b0fb-4644-a49b-4ff570dc37f0','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('b2b57bfd-1ef2-4963-8618-cb5a4fe3d693','68d14db1-396c-4d64-97cf-f006667d48fe','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('b4e7f44d-fc4f-45c4-b102-0fb041ea6c1f','9061afd2-4711-4908-b689-4141d0013d69','keyUse','ENC'),('b5e06310-1ba1-4678-b8a0-004be529077f','8ccc4751-d08c-4bab-93e5-f3af367c141e','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('b7ab90f9-011d-41bd-a5ee-20965eeafe61','f9b42840-4ce9-483e-b1d1-8465f046266d','max-clients','200'),('b7c4f1f1-b79c-47d0-ac70-9d5f519702d5','8ccc4751-d08c-4bab-93e5-f3af367c141e','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('ba1c615a-5612-439d-ac24-97dbedcd8823','57666d8a-3a28-4467-8977-dab5a8405bf2','allowed-protocol-mapper-types','oidc-address-mapper'),('bb64846e-de50-4030-b5b6-665b31f52b88','9061afd2-4711-4908-b689-4141d0013d69','priority','100'),('bdd4c53a-02f3-4251-bba1-f1a110fa85cb','589b00b8-e403-41fa-beef-753843c27c28','host-sending-registration-request-must-match','true'),('be83c4e0-34a1-4098-9657-531272c3ff23','97fd0a05-d85b-4468-966d-f7005f2a9d0d','allowed-protocol-mapper-types','saml-role-list-mapper'),('bf1c9053-2396-4c53-b546-07693259c791','3bb59e96-8933-44f6-a3d5-dbd560f450c3','priority','100'),('c01627bb-8f08-4c0a-af25-234a4bc1f0e2','2c8ac756-c844-46ce-bc98-7dabf49cea79','certificate','MIICozCCAYsCBgGO8DkXNTANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDDAp0ZXN0LXJlYWxtMB4XDTI0MDQxODA3NTkxN1oXDTM0MDQxODA4MDA1N1owFTETMBEGA1UEAwwKdGVzdC1yZWFsbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALH5DFkpgGD3bL/WvqLEJNky6QJi//Vh+8fk+3RcNMBQWBGlCcBenuykNw/xDS8NblJ516C1GCdYdsycPMy++iH5sWaWqtXwUWA//++7GqxbePBYN33djE7AaaBiO1ypxbnfyI8hk5rSzOd2LzfucfYGyowILOd4Aw47vLFaLWAZzH1O7n2Xqi8XuERQ8SHCCP6JI39xZtJK/tJ3NOKn85i4aVJy1awMyl7TLJ8KwHU8aN/IU8NxCZfbr6jmceidX0SZcYFZmXZSejW+Z5btsUhdjY0/5l45u/c21CHKTtF40+v8Nve6xm57uAS8LnU9F632BcZcXAfyCBY2aIOhNEECAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAdOkFgLwAfdZSvSDcc6pbxAbzVjq1X24R145iGuHGjkm3NfYSj3c1GrVWn8V+j92X3ZqfYb7/UmXXElp4Tl6uPAg3xYj5K1DVDKjX6YCklXyIv+dxvf8X+ZA6Ck9K5xKgcLAmGkgSh3RzI7s5Z4DRqvv/teBKvGi0tFVfufBtneMIS988JzX2RftnpDPx3sWB8Za0ZQiGLVs0AZJgbnIJ9oEdhUdG/L6lfnNLlHFPJSKFq+p+pBvYNr4AU3zcYm4k+RX0WpAl/h+/PWnB+MyL5GT+DSLvdDbpjXMXrvLT4Z0+bGlWtfd1aqBSdQdeC9trDAvzHf8jfNhyR8VBYMaU1A=='),('c1053275-a9c7-4182-b276-271e8b326357','2c8ac756-c844-46ce-bc98-7dabf49cea79','algorithm','RSA-OAEP'),('c21e94b3-65cb-466b-a89a-b1a6fca1b106','9061afd2-4711-4908-b689-4141d0013d69','algorithm','RSA-OAEP'),('c53efb5e-3824-475a-a249-f7076f2b6f8b','109e27cd-2b7c-4f16-a5bc-be11f305748b','allowed-protocol-mapper-types','saml-user-property-mapper'),('c54c3373-1b05-4575-902c-0460c083fa72','bf431d4f-712a-4417-bc9b-881c566afd0b','allow-default-scopes','true'),('c55e6f46-ceae-4077-993e-aea732f3cf44','4f7ea068-f73d-46b8-802f-019464083c96','secret','HSBv8q_oWyvqqJ6tqsd6ww'),('c6b60c9a-d872-45b1-a1fd-f33981098df8','5b29d81d-533a-4f90-a19a-f230c8212d6e','max-clients','200'),('c6f2fa0e-2871-431d-8c23-77ceba398fff','97fd0a05-d85b-4468-966d-f7005f2a9d0d','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('cb8ed7fa-01a3-48fc-9811-e82511dd5a22','41ca9748-0e9d-426b-8920-cbabcb302931','host-sending-registration-request-must-match','true'),('cd43fa2e-ad50-4ffe-9253-d80f8aac949b','68d14db1-396c-4d64-97cf-f006667d48fe','allowed-protocol-mapper-types','saml-user-property-mapper'),('cd9661c7-a0aa-4cec-9a57-1af67767bf25','b7798925-8e72-464d-b76e-951e5ea2a688','kid','f89eeb0a-357b-49d5-8ab0-1254c9dd361d'),('cfbe0a0d-e049-4e13-8000-a5185dbbab4d','3200da3d-a43f-47e2-89b9-dfd396feb0af','allow-default-scopes','true'),('d14993ab-8894-4e94-92c2-9ac248e78a0f','edd5efdf-b0fb-4644-a49b-4ff570dc37f0','allowed-protocol-mapper-types','oidc-address-mapper'),('d32a65a7-1aa5-4d8c-a682-faeef1507f85','27b6183a-7d2f-4572-b1a2-1c470a446ce9','priority','100'),('d5806cf2-0856-4319-a525-1a98193eb312','83773dcd-2802-4ae7-9fea-4b24fc6ba2d3','algorithm','HS512'),('d63fa052-364a-43dd-8b6c-63804dcb9511','8ccc4751-d08c-4bab-93e5-f3af367c141e','allowed-protocol-mapper-types','saml-role-list-mapper'),('dc5c61d8-55a6-4ee3-8755-94cef5455a1f','8a1b4b19-5104-429a-b262-f5f59d2b5606','keyUse','SIG'),('dd8fbe48-6032-452a-a11b-7209774f1128','68d14db1-396c-4d64-97cf-f006667d48fe','allowed-protocol-mapper-types','oidc-address-mapper'),('ddbdc267-e477-494c-ac35-c83947b442a3','fdd1d4a0-a99d-4ece-9fb5-8dca87d28f5e','certificate','MIICozCCAYsCBgGO8DkVmTANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDDAp0ZXN0LXJlYWxtMB4XDTI0MDQxODA3NTkxNloXDTM0MDQxODA4MDA1NlowFTETMBEGA1UEAwwKdGVzdC1yZWFsbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMPF5vV64g+oLKX16LVeHJtFYDvKiBh1LAtpmteRKIAXope0WfhxGpc5DuG3yI8tQJ+1erlo6b9unq6iVa9fuLIpeo/ectpbrtOv7qadNoBTEeIo7zUN6SDsxnjiGwIdhUCr7V32TtbPcvv7saDhyZouoWIsy7O2RgbWGFfb+yT17csVMt5wD9QCtwI6zxfRwOb5gXT+q6V8NvrP5bYU40iyUZMpHjmHa3JGvJ+spH0lb+tRkvNEH0hBhNlIoHBi60+D7Doo3IDPz/ObxcMj4I3bYye2iYyUrtb9dwdbALwE2EjNBHEGnIPKRXQbX9hLLjOb47HOedvv+UeM+3EAhm0CAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAA+dEy39ZD72DOoGxBoCyl1aqwziAi5P84VOWay5bgMwewi7p4lfUYhbxHwv0B5YAvUMihpfId58CjT+wdpmMtuVf0HK9AB48v1x+arTBtgDCYbYO9w5DJUvop8x929iB1bimZ9nqOMwCs9HpVzO4vr5DGGpmeORsrPVtVDupGg1yQ1YaTfUk2Mts8XxNZtXsHEw9xhHlmj6R4IkC56IAq/iBd7Wf9ZgMfJ1iiRQfefHzlHFOgXvit5LKy7mjhNl1BMlt6DRc4FI6da3ZY/wTHs41+e2rT8IwJNYDbZnTclNJyZlpHLje260D6+6AVBFAp0IfiE8TmJ8+/qwbhq21Ew=='),('ddcdf50d-3785-4997-a00c-0ab288c2f116','edd5efdf-b0fb-4644-a49b-4ff570dc37f0','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('df1770d4-e5e1-4f12-93ae-ff00d67b012c','3293c342-c07d-4706-9af5-30ab7cc3ee2a','client-uris-must-match','true'),('e00dbc06-d21b-419d-9dc5-b790b29f374d','83773dcd-2802-4ae7-9fea-4b24fc6ba2d3','priority','100'),('e47a571a-bc94-4621-bdfb-ca5e55590ca4','d2bb4848-001e-4e99-bc50-aa183de69051','priority','100'),('e4a5f52f-1bda-4dc6-b1f6-42520b7f06dc','97fd0a05-d85b-4468-966d-f7005f2a9d0d','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('e7338bac-afc6-4e12-8861-a61d10bd75fb','049eba05-ff13-48e0-a1ea-c63bd8cd808a','priority','100'),('f1a246c0-5426-4b29-a2a5-b51806cc46f9','edd5efdf-b0fb-4644-a49b-4ff570dc37f0','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('f2afff45-5544-45a7-bec2-a54368a655ad','57666d8a-3a28-4467-8977-dab5a8405bf2','allowed-protocol-mapper-types','saml-role-list-mapper'),('f307dae8-8123-4bc7-845f-13758058b6e9','b7798925-8e72-464d-b76e-951e5ea2a688','secret','ASLw48UqwFW8lJZW9OiHnw'),('f3bb2ec6-9ddd-4d44-a460-20f54b45319e','b7798925-8e72-464d-b76e-951e5ea2a688','priority','100'),('f70159b5-4d94-48c7-8d4d-27e6bc79bb95','bc9cb75c-4512-4d5b-9284-27e900a79763','allow-default-scopes','true'),('fb1b1e3a-ca56-4154-aceb-d9a0db519ade','edd5efdf-b0fb-4644-a49b-4ff570dc37f0','allowed-protocol-mapper-types','saml-role-list-mapper'),('fbb0e2a6-6c51-43ed-9e70-badbc5e7931d','109e27cd-2b7c-4f16-a5bc-be11f305748b','allowed-protocol-mapper-types','saml-role-list-mapper'),('fc6adb28-d792-4c09-92b6-87e765d73055','3bb59e96-8933-44f6-a3d5-dbd560f450c3','algorithm','HS512');
/*!40000 ALTER TABLE `COMPONENT_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `COMPOSITE_ROLE`
--

DROP TABLE IF EXISTS `COMPOSITE_ROLE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `COMPOSITE_ROLE` (
  `COMPOSITE` varchar(36) NOT NULL,
  `CHILD_ROLE` varchar(36) NOT NULL,
  PRIMARY KEY (`COMPOSITE`,`CHILD_ROLE`),
  KEY `IDX_COMPOSITE` (`COMPOSITE`),
  KEY `IDX_COMPOSITE_CHILD` (`CHILD_ROLE`),
  CONSTRAINT `FK_A63WVEKFTU8JO1PNJ81E7MCE2` FOREIGN KEY (`COMPOSITE`) REFERENCES `KEYCLOAK_ROLE` (`ID`),
  CONSTRAINT `FK_GR7THLLB9LU8Q4VQA4524JJY8` FOREIGN KEY (`CHILD_ROLE`) REFERENCES `KEYCLOAK_ROLE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `COMPOSITE_ROLE`
--

LOCK TABLES `COMPOSITE_ROLE` WRITE;
/*!40000 ALTER TABLE `COMPOSITE_ROLE` DISABLE KEYS */;
INSERT INTO `COMPOSITE_ROLE` (`COMPOSITE`, `CHILD_ROLE`) VALUES ('0daf1344-2778-40fe-a0de-edfed0959aa8','a118a201-c733-428b-bf9c-121bcccfb8b5'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','175e4ef5-829d-4b42-a98e-dbebadb49f04'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','22990ac5-82f5-4deb-97d1-342a47945b59'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','6bbfee03-d9ef-43c5-85b6-59627959df8d'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','8c549c79-801f-4d45-892d-22fe10111eb2'),('27381875-2dde-4add-93ab-9da6abf300ce','0fe3ce64-452d-4c89-99f8-af1921093c45'),('27381875-2dde-4add-93ab-9da6abf300ce','1492ec55-0066-41d1-bfe2-bff7e7ac5259'),('27381875-2dde-4add-93ab-9da6abf300ce','1512286b-232a-48da-a722-79bfa293bd1d'),('27381875-2dde-4add-93ab-9da6abf300ce','2cf04f97-ea0a-4c30-bcd2-00a25d0700cf'),('27381875-2dde-4add-93ab-9da6abf300ce','2d58e9c6-470c-4884-bcf5-c6f92e437a2e'),('27381875-2dde-4add-93ab-9da6abf300ce','609ad758-5a8f-42f5-8777-456c380c014f'),('27381875-2dde-4add-93ab-9da6abf300ce','65374982-5c23-4b66-a06a-d969c0cd154a'),('27381875-2dde-4add-93ab-9da6abf300ce','71e7629d-894b-4aff-825a-765cf58a6f6b'),('27381875-2dde-4add-93ab-9da6abf300ce','76a88fee-aa85-44de-a9de-81c6823d603f'),('27381875-2dde-4add-93ab-9da6abf300ce','7717b2f2-fcf1-45d4-86db-0d4c7f9e91a2'),('27381875-2dde-4add-93ab-9da6abf300ce','88d1e80b-8996-4609-a820-bd3d3efc57ab'),('27381875-2dde-4add-93ab-9da6abf300ce','8b146ddf-257f-4390-90ea-2d03c7ab7ce1'),('27381875-2dde-4add-93ab-9da6abf300ce','a4947d89-da1b-4fb1-ae4b-323921695182'),('27381875-2dde-4add-93ab-9da6abf300ce','aea97389-2ee0-4b61-96a7-ff09a73d7996'),('27381875-2dde-4add-93ab-9da6abf300ce','bef3ad5c-c6ed-4710-9800-ad029e4632f3'),('27381875-2dde-4add-93ab-9da6abf300ce','c7f14c0f-92f7-4a83-947c-9191b4958168'),('27381875-2dde-4add-93ab-9da6abf300ce','d872e92a-c9f5-482c-9e10-de35d8ed80a3'),('27381875-2dde-4add-93ab-9da6abf300ce','dee91445-1e60-47a2-9f1a-51a8fdf925ab'),('2bbbedd2-62e8-4965-a5cd-6664c7b100a5','245ce92a-a642-4653-bd40-4bb0b179a387'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','06ad34dc-2164-4067-896b-933e9dadb08c'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','07178d4e-bdd9-412f-924b-62977dc1d23b'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','0746f171-6572-4ee4-8cab-2141dcd70f13'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','07755138-7efe-46df-9ae5-4f7869d5c570'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','0cd165e8-0b76-44c6-8c78-23f2f8b96e18'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','0edbe58b-31c3-4503-81e6-3426b88107d7'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','0f5abc1b-e9de-4f60-94fa-7a9f537eb90a'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','10114928-9ec5-4ce4-baf0-696dd614ea04'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','173fcaa5-a382-4ae8-819a-317f8d671bec'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','1ad87cfd-6f78-42f8-9b1e-220a2a813fe1'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','1b146798-0c31-4fe4-8d8d-8341e5b463ec'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','241a9760-03f0-41fa-911c-0a023532161e'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','245ce92a-a642-4653-bd40-4bb0b179a387'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','24e2a886-cc15-44c9-a176-1b12c038f88d'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','2bbbedd2-62e8-4965-a5cd-6664c7b100a5'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','2d0f3e90-9421-4164-b396-ec0ec75b7099'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','2e299e39-7ae8-4394-82af-0488b8ba850c'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','31102399-1273-491f-8abf-99c630871105'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','39d9f8a7-97c4-459c-9a52-78e597492a95'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','39e1e0ae-01d2-40c5-82b7-66727d7542bb'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','41a1c8c2-a5f7-433d-a283-f71f08f7779e'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','43087049-6e35-422e-a082-97d32212cfcf'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','4f90bb4f-92f2-45dc-ad8a-f97723d1d828'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','57e92e5d-9b6b-4782-9045-d5ac3b1cd017'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','581d6c1d-f6d2-4d8c-bfcc-56d8e5f8de3d'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','623ac08d-4ff5-49c7-9392-c6f8a49ef26b'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','6552d7cc-3d28-424a-b9ab-c0719ad9d1f1'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','76ab719c-a862-441e-a66b-329547590e1d'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','7a15b2f4-66f6-4ab9-b90e-1d3054723dd1'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','84ec0238-dc71-4066-b99e-49e59cc6bc7a'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','888d64fb-b137-4f1e-b530-b4aef00ab967'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','8a2220b6-356e-4d7e-995f-83ca8da2e9fe'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','8b6e8ddd-250e-4bac-8bdd-de3a9d4ab9e6'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','8bcad746-bec2-4e86-ba9a-e85979488a4f'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','94526ddc-9e84-4c0b-958b-07184a18ab32'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','a2fa72c7-f176-4d55-aab3-9c6e18a6c3e3'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','a6381b94-b388-4ac0-94a4-a9c24384f9ba'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','a880a908-88f7-45fc-bf1e-4f6fa51fff11'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','ad5932af-bcfc-4464-962e-5490c41eda23'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','b717dbd5-35d8-4865-94c2-bf876c0146f4'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','bad4e73b-8479-4bfc-a203-68300ea9aaf4'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','c1887edf-e0fe-43f9-9857-97e256a6df6e'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','c73afc5a-5da9-416e-a0d9-d50931ad6c0e'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','c8ee3a69-7a2d-48f8-9373-6677e3724523'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','ce73d6bd-9432-4068-a9e6-02f000f8b4d3'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','d0f5422f-9597-4ea6-abbb-7c346c6a8096'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','d5f7c5c8-33c4-483c-802e-6ed2c1bc2f29'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','d69ea56d-5246-4e89-a45f-4effdf6d9bcf'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','d892644d-7123-4207-a8c5-74fcc7da7653'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','db24697d-b20d-4a2a-99b6-f41d26e73b7e'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','df5897b4-53cf-44d4-94ac-c8d1651af181'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','e293b015-eada-4709-8d6c-da9c98a5ec04'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','ea082079-5ed5-4688-8d02-1a0a2dbed4fa'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','ee7eb7c2-62cf-4526-9560-a1e1c4bd27f2'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','f9e0e15e-14f9-4327-b781-04bc3078317e'),('4b5fb1ff-dc8c-4f63-bfaf-0772d809fd00','ec9a917b-c79f-422b-941d-40214a7e84cf'),('623ac08d-4ff5-49c7-9392-c6f8a49ef26b','db24697d-b20d-4a2a-99b6-f41d26e73b7e'),('6552d7cc-3d28-424a-b9ab-c0719ad9d1f1','888d64fb-b137-4f1e-b530-b4aef00ab967'),('76a88fee-aa85-44de-a9de-81c6823d603f','d872e92a-c9f5-482c-9e10-de35d8ed80a3'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','29fbe3d8-6cde-4599-a684-c054a2abb2f4'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','3323ca6f-35d6-4ef6-bba9-d3a65e39d42d'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','346d7236-020b-4f7d-8922-a51bf71a9482'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','4b5fb1ff-dc8c-4f63-bfaf-0772d809fd00'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','612b43cd-8762-4e1e-862d-4c632a35b038'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','661dd437-4d3d-4e80-a862-83a2ae44ca5e'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','669e6b20-0a78-4732-a482-fc68d855d1df'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','7252897e-111a-4b0e-b87c-c13d52e42696'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','84b5254f-5b61-45e9-9081-7a93b3eba444'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','8abdd60d-ccc5-4315-85b6-9b927eaca4a7'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','95575514-531e-4f81-91a9-f0c35b743a81'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','9b9a62d1-f35a-4d79-823e-13957aff5203'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','c8fc11ef-66b7-438a-a52e-8668bf8ac1b4'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','cb163866-350b-404d-9bbf-13407643d098'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','cde65e2d-9505-48a9-8a5b-4831d65773bf'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','d305da17-3117-49ce-90eb-91fc9bb929af'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','ec9a917b-c79f-422b-941d-40214a7e84cf'),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','f6fafe29-29a1-4124-9a53-c98fd932885d'),('84b5254f-5b61-45e9-9081-7a93b3eba444','661dd437-4d3d-4e80-a862-83a2ae44ca5e'),('84b5254f-5b61-45e9-9081-7a93b3eba444','8abdd60d-ccc5-4315-85b6-9b927eaca4a7'),('8b6e8ddd-250e-4bac-8bdd-de3a9d4ab9e6','c8ee3a69-7a2d-48f8-9373-6677e3724523'),('8b6e8ddd-250e-4bac-8bdd-de3a9d4ab9e6','ea082079-5ed5-4688-8d02-1a0a2dbed4fa'),('8c549c79-801f-4d45-892d-22fe10111eb2','81666137-6815-4158-9161-116220fffe4a'),('8e5f68f7-bcb4-43a5-ae98-6dfc29e18395','a621e3b3-8f5f-460b-ad97-f3ab58c894cb'),('91e47ad5-5b4d-4440-b982-3dd69c10b2fb','dcc9c139-ec68-483f-9b2e-6db54c81fad4'),('94526ddc-9e84-4c0b-958b-07184a18ab32','b717dbd5-35d8-4865-94c2-bf876c0146f4'),('94526ddc-9e84-4c0b-958b-07184a18ab32','d5f7c5c8-33c4-483c-802e-6ed2c1bc2f29'),('bef3ad5c-c6ed-4710-9800-ad029e4632f3','65374982-5c23-4b66-a06a-d969c0cd154a'),('bef3ad5c-c6ed-4710-9800-ad029e4632f3','88d1e80b-8996-4609-a820-bd3d3efc57ab'),('c4fb1f60-f60e-4022-9b59-36c1eef2f82a','2a6a0506-9553-4b69-822f-6be7668d3b91'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','4d890335-e9e1-45ab-8c8c-98c9e8db8bc4'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','5d573d8e-9b92-4b95-a6e7-84f8b1eb3cef'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','91e47ad5-5b4d-4440-b982-3dd69c10b2fb'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','d62ac468-1487-4e80-afb8-d12035ef8d24'),('d0f5422f-9597-4ea6-abbb-7c346c6a8096','41a1c8c2-a5f7-433d-a283-f71f08f7779e'),('d0f5422f-9597-4ea6-abbb-7c346c6a8096','4f90bb4f-92f2-45dc-ad8a-f97723d1d828'),('e36aa23d-1f2a-417d-9f68-90234cd0703e','2c0e4faf-3427-4063-ae6d-7c8b65a292ba'),('e7b5b42b-09a2-46bf-b6b5-87f71892c9dd','22ae141b-4a27-4414-9d66-d95ebc1c264c'),('e7b5b42b-09a2-46bf-b6b5-87f71892c9dd','4590b97a-ecbe-42a6-a0fc-85fc8f576d5b'),('e7b5b42b-09a2-46bf-b6b5-87f71892c9dd','8e5f68f7-bcb4-43a5-ae98-6dfc29e18395'),('e7b5b42b-09a2-46bf-b6b5-87f71892c9dd','b197dcfd-a869-436d-8746-8d288ae56ef7');
/*!40000 ALTER TABLE `COMPOSITE_ROLE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CREDENTIAL`
--

DROP TABLE IF EXISTS `CREDENTIAL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CREDENTIAL` (
  `ID` varchar(36) NOT NULL,
  `SALT` tinyblob,
  `TYPE` varchar(255) DEFAULT NULL,
  `USER_ID` varchar(36) DEFAULT NULL,
  `CREATED_DATE` bigint DEFAULT NULL,
  `USER_LABEL` varchar(255) DEFAULT NULL,
  `SECRET_DATA` longtext,
  `CREDENTIAL_DATA` longtext,
  `PRIORITY` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_USER_CREDENTIAL` (`USER_ID`),
  CONSTRAINT `FK_PFYR0GLASQYL0DEI3KL69R6V0` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CREDENTIAL`
--

LOCK TABLES `CREDENTIAL` WRITE;
/*!40000 ALTER TABLE `CREDENTIAL` DISABLE KEYS */;
INSERT INTO `CREDENTIAL` (`ID`, `SALT`, `TYPE`, `USER_ID`, `CREATED_DATE`, `USER_LABEL`, `SECRET_DATA`, `CREDENTIAL_DATA`, `PRIORITY`) VALUES ('b9c9505b-b6c0-4c0d-b87a-67e45cf749d5',NULL,'password','491fe75e-534f-48ae-ada1-36ea0ee54a1f',1713408446025,NULL,'{\"value\":\"xR7mK9pL2nQ5vW8jY4cB1fA6hD3sE0tU9iO2wN5rM7kP4qJ8zX6gV1bC3yH0mL5n==\",\"salt\":\"aB3cD5eF7gH9iJ1kL3mN5o==\",\"additionalParameters\":{}}','{\"hashIterations\":210000,\"algorithm\":\"pbkdf2-sha512\",\"additionalParameters\":{}}',10);
/*!40000 ALTER TABLE `CREDENTIAL` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `DATABASECHANGELOG`
--

DROP TABLE IF EXISTS `DATABASECHANGELOG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DATABASECHANGELOG` (
  `ID` varchar(255) NOT NULL,
  `AUTHOR` varchar(255) NOT NULL,
  `FILENAME` varchar(255) NOT NULL,
  `DATEEXECUTED` datetime NOT NULL,
  `ORDEREXECUTED` int NOT NULL,
  `EXECTYPE` varchar(10) NOT NULL,
  `MD5SUM` varchar(35) DEFAULT NULL,
  `DESCRIPTION` varchar(255) DEFAULT NULL,
  `COMMENTS` varchar(255) DEFAULT NULL,
  `TAG` varchar(255) DEFAULT NULL,
  `LIQUIBASE` varchar(20) DEFAULT NULL,
  `CONTEXTS` varchar(255) DEFAULT NULL,
  `LABELS` varchar(255) DEFAULT NULL,
  `DEPLOYMENT_ID` varchar(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `DATABASECHANGELOG`
--

LOCK TABLES `DATABASECHANGELOG` WRITE;
/*!40000 ALTER TABLE `DATABASECHANGELOG` DISABLE KEYS */;
INSERT INTO `DATABASECHANGELOG` (`ID`, `AUTHOR`, `FILENAME`, `DATEEXECUTED`, `ORDEREXECUTED`, `EXECTYPE`, `MD5SUM`, `DESCRIPTION`, `COMMENTS`, `TAG`, `LIQUIBASE`, `CONTEXTS`, `LABELS`, `DEPLOYMENT_ID`) VALUES ('1.0.0.Final-KEYCLOAK-5461','sthorger@redhat.com','META-INF/jpa-changelog-1.0.0.Final.xml','2024-04-18 11:46:58',1,'EXECUTED','9:6f1016664e21e16d26517a4418f5e3df','createTable tableName=APPLICATION_DEFAULT_ROLES; createTable tableName=CLIENT; createTable tableName=CLIENT_SESSION; createTable tableName=CLIENT_SESSION_ROLE; createTable tableName=COMPOSITE_ROLE; createTable tableName=CREDENTIAL; createTable tab...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.0.0.Final-KEYCLOAK-5461','sthorger@redhat.com','META-INF/db2-jpa-changelog-1.0.0.Final.xml','2024-04-18 11:46:58',2,'MARK_RAN','9:828775b1596a07d1200ba1d49e5e3941','createTable tableName=APPLICATION_DEFAULT_ROLES; createTable tableName=CLIENT; createTable tableName=CLIENT_SESSION; createTable tableName=CLIENT_SESSION_ROLE; createTable tableName=COMPOSITE_ROLE; createTable tableName=CREDENTIAL; createTable tab...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.1.0.Beta1','sthorger@redhat.com','META-INF/jpa-changelog-1.1.0.Beta1.xml','2024-04-18 11:46:59',3,'EXECUTED','9:5f090e44a7d595883c1fb61f4b41fd38','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION; createTable tableName=CLIENT_ATTRIBUTES; createTable tableName=CLIENT_SESSION_NOTE; createTable tableName=APP_NODE_REGISTRATIONS; addColumn table...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.1.0.Final','sthorger@redhat.com','META-INF/jpa-changelog-1.1.0.Final.xml','2024-04-18 11:46:59',4,'EXECUTED','9:c07e577387a3d2c04d1adc9aaad8730e','renameColumn newColumnName=EVENT_TIME, oldColumnName=TIME, tableName=EVENT_ENTITY','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.2.0.Beta1','psilva@redhat.com','META-INF/jpa-changelog-1.2.0.Beta1.xml','2024-04-18 11:47:00',5,'EXECUTED','9:b68ce996c655922dbcd2fe6b6ae72686','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION; createTable tableName=PROTOCOL_MAPPER; createTable tableName=PROTOCOL_MAPPER_CONFIG; createTable tableName=...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.2.0.Beta1','psilva@redhat.com','META-INF/db2-jpa-changelog-1.2.0.Beta1.xml','2024-04-18 11:47:00',6,'MARK_RAN','9:543b5c9989f024fe35c6f6c5a97de88e','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION; createTable tableName=PROTOCOL_MAPPER; createTable tableName=PROTOCOL_MAPPER_CONFIG; createTable tableName=...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.2.0.RC1','bburke@redhat.com','META-INF/jpa-changelog-1.2.0.CR1.xml','2024-04-18 11:47:02',7,'EXECUTED','9:765afebbe21cf5bbca048e632df38336','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete tableName=USER_SESSION; createTable tableName=MIGRATION_MODEL; createTable tableName=IDENTITY_P...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.2.0.RC1','bburke@redhat.com','META-INF/db2-jpa-changelog-1.2.0.CR1.xml','2024-04-18 11:47:02',8,'MARK_RAN','9:db4a145ba11a6fdaefb397f6dbf829a1','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete tableName=USER_SESSION; createTable tableName=MIGRATION_MODEL; createTable tableName=IDENTITY_P...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.2.0.Final','keycloak','META-INF/jpa-changelog-1.2.0.Final.xml','2024-04-18 11:47:02',9,'EXECUTED','9:9d05c7be10cdb873f8bcb41bc3a8ab23','update tableName=CLIENT; update tableName=CLIENT; update tableName=CLIENT','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.3.0','bburke@redhat.com','META-INF/jpa-changelog-1.3.0.xml','2024-04-18 11:47:03',10,'EXECUTED','9:18593702353128d53111f9b1ff0b82b8','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_PROT_MAPPER; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete tableName=USER_SESSION; createTable tableName=ADMI...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.4.0','bburke@redhat.com','META-INF/jpa-changelog-1.4.0.xml','2024-04-18 11:47:04',11,'EXECUTED','9:6122efe5f090e41a85c0f1c9e52cbb62','delete tableName=CLIENT_SESSION_AUTH_STATUS; delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_PROT_MAPPER; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete table...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.4.0','bburke@redhat.com','META-INF/db2-jpa-changelog-1.4.0.xml','2024-04-18 11:47:04',12,'MARK_RAN','9:e1ff28bf7568451453f844c5d54bb0b5','delete tableName=CLIENT_SESSION_AUTH_STATUS; delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_PROT_MAPPER; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete table...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.5.0','bburke@redhat.com','META-INF/jpa-changelog-1.5.0.xml','2024-04-18 11:47:04',13,'EXECUTED','9:7af32cd8957fbc069f796b61217483fd','delete tableName=CLIENT_SESSION_AUTH_STATUS; delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_PROT_MAPPER; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete table...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.6.1_from15','mposolda@redhat.com','META-INF/jpa-changelog-1.6.1.xml','2024-04-18 11:47:04',14,'EXECUTED','9:6005e15e84714cd83226bf7879f54190','addColumn tableName=REALM; addColumn tableName=KEYCLOAK_ROLE; addColumn tableName=CLIENT; createTable tableName=OFFLINE_USER_SESSION; createTable tableName=OFFLINE_CLIENT_SESSION; addPrimaryKey constraintName=CONSTRAINT_OFFL_US_SES_PK2, tableName=...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.6.1_from16-pre','mposolda@redhat.com','META-INF/jpa-changelog-1.6.1.xml','2024-04-18 11:47:04',15,'MARK_RAN','9:bf656f5a2b055d07f314431cae76f06c','delete tableName=OFFLINE_CLIENT_SESSION; delete tableName=OFFLINE_USER_SESSION','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.6.1_from16','mposolda@redhat.com','META-INF/jpa-changelog-1.6.1.xml','2024-04-18 11:47:04',16,'MARK_RAN','9:f8dadc9284440469dcf71e25ca6ab99b','dropPrimaryKey constraintName=CONSTRAINT_OFFLINE_US_SES_PK, tableName=OFFLINE_USER_SESSION; dropPrimaryKey constraintName=CONSTRAINT_OFFLINE_CL_SES_PK, tableName=OFFLINE_CLIENT_SESSION; addColumn tableName=OFFLINE_USER_SESSION; update tableName=OF...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.6.1','mposolda@redhat.com','META-INF/jpa-changelog-1.6.1.xml','2024-04-18 11:47:04',17,'EXECUTED','9:d41d8cd98f00b204e9800998ecf8427e','empty','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.7.0','bburke@redhat.com','META-INF/jpa-changelog-1.7.0.xml','2024-04-18 11:47:05',18,'EXECUTED','9:3368ff0be4c2855ee2dd9ca813b38d8e','createTable tableName=KEYCLOAK_GROUP; createTable tableName=GROUP_ROLE_MAPPING; createTable tableName=GROUP_ATTRIBUTE; createTable tableName=USER_GROUP_MEMBERSHIP; createTable tableName=REALM_DEFAULT_GROUPS; addColumn tableName=IDENTITY_PROVIDER; ...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.8.0','mposolda@redhat.com','META-INF/jpa-changelog-1.8.0.xml','2024-04-18 11:47:06',19,'EXECUTED','9:8ac2fb5dd030b24c0570a763ed75ed20','addColumn tableName=IDENTITY_PROVIDER; createTable tableName=CLIENT_TEMPLATE; createTable tableName=CLIENT_TEMPLATE_ATTRIBUTES; createTable tableName=TEMPLATE_SCOPE_MAPPING; dropNotNullConstraint columnName=CLIENT_ID, tableName=PROTOCOL_MAPPER; ad...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.8.0-2','keycloak','META-INF/jpa-changelog-1.8.0.xml','2024-04-18 11:47:06',20,'EXECUTED','9:f91ddca9b19743db60e3057679810e6c','dropDefaultValue columnName=ALGORITHM, tableName=CREDENTIAL; update tableName=CREDENTIAL','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.8.0','mposolda@redhat.com','META-INF/db2-jpa-changelog-1.8.0.xml','2024-04-18 11:47:06',21,'MARK_RAN','9:831e82914316dc8a57dc09d755f23c51','addColumn tableName=IDENTITY_PROVIDER; createTable tableName=CLIENT_TEMPLATE; createTable tableName=CLIENT_TEMPLATE_ATTRIBUTES; createTable tableName=TEMPLATE_SCOPE_MAPPING; dropNotNullConstraint columnName=CLIENT_ID, tableName=PROTOCOL_MAPPER; ad...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.8.0-2','keycloak','META-INF/db2-jpa-changelog-1.8.0.xml','2024-04-18 11:47:06',22,'MARK_RAN','9:f91ddca9b19743db60e3057679810e6c','dropDefaultValue columnName=ALGORITHM, tableName=CREDENTIAL; update tableName=CREDENTIAL','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.9.0','mposolda@redhat.com','META-INF/jpa-changelog-1.9.0.xml','2024-04-18 11:47:06',23,'EXECUTED','9:bc3d0f9e823a69dc21e23e94c7a94bb1','update tableName=REALM; update tableName=REALM; update tableName=REALM; update tableName=REALM; update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=REALM; update tableName=REALM; customChange; dr...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.9.1','keycloak','META-INF/jpa-changelog-1.9.1.xml','2024-04-18 11:47:06',24,'EXECUTED','9:c9999da42f543575ab790e76439a2679','modifyDataType columnName=PRIVATE_KEY, tableName=REALM; modifyDataType columnName=PUBLIC_KEY, tableName=REALM; modifyDataType columnName=CERTIFICATE, tableName=REALM','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.9.1','keycloak','META-INF/db2-jpa-changelog-1.9.1.xml','2024-04-18 11:47:06',25,'MARK_RAN','9:0d6c65c6f58732d81569e77b10ba301d','modifyDataType columnName=PRIVATE_KEY, tableName=REALM; modifyDataType columnName=CERTIFICATE, tableName=REALM','',NULL,'4.25.1',NULL,NULL,'3408415468'),('1.9.2','keycloak','META-INF/jpa-changelog-1.9.2.xml','2024-04-18 11:47:06',26,'EXECUTED','9:fc576660fc016ae53d2d4778d84d86d0','createIndex indexName=IDX_USER_EMAIL, tableName=USER_ENTITY; createIndex indexName=IDX_USER_ROLE_MAPPING, tableName=USER_ROLE_MAPPING; createIndex indexName=IDX_USER_GROUP_MAPPING, tableName=USER_GROUP_MEMBERSHIP; createIndex indexName=IDX_USER_CO...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-2.0.0','psilva@redhat.com','META-INF/jpa-changelog-authz-2.0.0.xml','2024-04-18 11:47:07',27,'EXECUTED','9:43ed6b0da89ff77206289e87eaa9c024','createTable tableName=RESOURCE_SERVER; addPrimaryKey constraintName=CONSTRAINT_FARS, tableName=RESOURCE_SERVER; addUniqueConstraint constraintName=UK_AU8TT6T700S9V50BU18WS5HA6, tableName=RESOURCE_SERVER; createTable tableName=RESOURCE_SERVER_RESOU...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-2.5.1','psilva@redhat.com','META-INF/jpa-changelog-authz-2.5.1.xml','2024-04-18 11:47:07',28,'EXECUTED','9:44bae577f551b3738740281eceb4ea70','update tableName=RESOURCE_SERVER_POLICY','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.1.0-KEYCLOAK-5461','bburke@redhat.com','META-INF/jpa-changelog-2.1.0.xml','2024-04-18 11:47:08',29,'EXECUTED','9:bd88e1f833df0420b01e114533aee5e8','createTable tableName=BROKER_LINK; createTable tableName=FED_USER_ATTRIBUTE; createTable tableName=FED_USER_CONSENT; createTable tableName=FED_USER_CONSENT_ROLE; createTable tableName=FED_USER_CONSENT_PROT_MAPPER; createTable tableName=FED_USER_CR...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.2.0','bburke@redhat.com','META-INF/jpa-changelog-2.2.0.xml','2024-04-18 11:47:08',30,'EXECUTED','9:a7022af5267f019d020edfe316ef4371','addColumn tableName=ADMIN_EVENT_ENTITY; createTable tableName=CREDENTIAL_ATTRIBUTE; createTable tableName=FED_CREDENTIAL_ATTRIBUTE; modifyDataType columnName=VALUE, tableName=CREDENTIAL; addForeignKeyConstraint baseTableName=FED_CREDENTIAL_ATTRIBU...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.3.0','bburke@redhat.com','META-INF/jpa-changelog-2.3.0.xml','2024-04-18 11:47:09',31,'EXECUTED','9:fc155c394040654d6a79227e56f5e25a','createTable tableName=FEDERATED_USER; addPrimaryKey constraintName=CONSTR_FEDERATED_USER, tableName=FEDERATED_USER; dropDefaultValue columnName=TOTP, tableName=USER_ENTITY; dropColumn columnName=TOTP, tableName=USER_ENTITY; addColumn tableName=IDE...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.4.0','bburke@redhat.com','META-INF/jpa-changelog-2.4.0.xml','2024-04-18 11:47:09',32,'EXECUTED','9:eac4ffb2a14795e5dc7b426063e54d88','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.5.0','bburke@redhat.com','META-INF/jpa-changelog-2.5.0.xml','2024-04-18 11:47:09',33,'EXECUTED','9:54937c05672568c4c64fc9524c1e9462','customChange; modifyDataType columnName=USER_ID, tableName=OFFLINE_USER_SESSION','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.5.0-unicode-oracle','hmlnarik@redhat.com','META-INF/jpa-changelog-2.5.0.xml','2024-04-18 11:47:09',34,'MARK_RAN','9:3a32bace77c84d7678d035a7f5a8084e','modifyDataType columnName=DESCRIPTION, tableName=AUTHENTICATION_FLOW; modifyDataType columnName=DESCRIPTION, tableName=CLIENT_TEMPLATE; modifyDataType columnName=DESCRIPTION, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=DESCRIPTION,...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.5.0-unicode-other-dbs','hmlnarik@redhat.com','META-INF/jpa-changelog-2.5.0.xml','2024-04-18 11:47:10',35,'EXECUTED','9:33d72168746f81f98ae3a1e8e0ca3554','modifyDataType columnName=DESCRIPTION, tableName=AUTHENTICATION_FLOW; modifyDataType columnName=DESCRIPTION, tableName=CLIENT_TEMPLATE; modifyDataType columnName=DESCRIPTION, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=DESCRIPTION,...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.5.0-duplicate-email-support','slawomir@dabek.name','META-INF/jpa-changelog-2.5.0.xml','2024-04-18 11:47:10',36,'EXECUTED','9:61b6d3d7a4c0e0024b0c839da283da0c','addColumn tableName=REALM','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.5.0-unique-group-names','hmlnarik@redhat.com','META-INF/jpa-changelog-2.5.0.xml','2024-04-18 11:47:10',37,'EXECUTED','9:8dcac7bdf7378e7d823cdfddebf72fda','addUniqueConstraint constraintName=SIBLING_NAMES, tableName=KEYCLOAK_GROUP','',NULL,'4.25.1',NULL,NULL,'3408415468'),('2.5.1','bburke@redhat.com','META-INF/jpa-changelog-2.5.1.xml','2024-04-18 11:47:10',38,'EXECUTED','9:a2b870802540cb3faa72098db5388af3','addColumn tableName=FED_USER_CONSENT','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.0.0','bburke@redhat.com','META-INF/jpa-changelog-3.0.0.xml','2024-04-18 11:47:10',39,'EXECUTED','9:132a67499ba24bcc54fb5cbdcfe7e4c0','addColumn tableName=IDENTITY_PROVIDER','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.2.0-fix','keycloak','META-INF/jpa-changelog-3.2.0.xml','2024-04-18 11:47:10',40,'MARK_RAN','9:938f894c032f5430f2b0fafb1a243462','addNotNullConstraint columnName=REALM_ID, tableName=CLIENT_INITIAL_ACCESS','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.2.0-fix-with-keycloak-5416','keycloak','META-INF/jpa-changelog-3.2.0.xml','2024-04-18 11:47:10',41,'MARK_RAN','9:845c332ff1874dc5d35974b0babf3006','dropIndex indexName=IDX_CLIENT_INIT_ACC_REALM, tableName=CLIENT_INITIAL_ACCESS; addNotNullConstraint columnName=REALM_ID, tableName=CLIENT_INITIAL_ACCESS; createIndex indexName=IDX_CLIENT_INIT_ACC_REALM, tableName=CLIENT_INITIAL_ACCESS','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.2.0-fix-offline-sessions','hmlnarik','META-INF/jpa-changelog-3.2.0.xml','2024-04-18 11:47:10',42,'EXECUTED','9:fc86359c079781adc577c5a217e4d04c','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.2.0-fixed','keycloak','META-INF/jpa-changelog-3.2.0.xml','2024-04-18 11:47:12',43,'EXECUTED','9:59a64800e3c0d09b825f8a3b444fa8f4','addColumn tableName=REALM; dropPrimaryKey constraintName=CONSTRAINT_OFFL_CL_SES_PK2, tableName=OFFLINE_CLIENT_SESSION; dropColumn columnName=CLIENT_SESSION_ID, tableName=OFFLINE_CLIENT_SESSION; addPrimaryKey constraintName=CONSTRAINT_OFFL_CL_SES_P...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.3.0','keycloak','META-INF/jpa-changelog-3.3.0.xml','2024-04-18 11:47:12',44,'EXECUTED','9:d48d6da5c6ccf667807f633fe489ce88','addColumn tableName=USER_ENTITY','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-3.4.0.CR1-resource-server-pk-change-part1','glavoie@gmail.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2024-04-18 11:47:12',45,'EXECUTED','9:dde36f7973e80d71fceee683bc5d2951','addColumn tableName=RESOURCE_SERVER_POLICY; addColumn tableName=RESOURCE_SERVER_RESOURCE; addColumn tableName=RESOURCE_SERVER_SCOPE','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-3.4.0.CR1-resource-server-pk-change-part2-KEYCLOAK-6095','hmlnarik@redhat.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2024-04-18 11:47:12',46,'EXECUTED','9:b855e9b0a406b34fa323235a0cf4f640','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-3.4.0.CR1-resource-server-pk-change-part3-fixed','glavoie@gmail.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2024-04-18 11:47:12',47,'MARK_RAN','9:51abbacd7b416c50c4421a8cabf7927e','dropIndex indexName=IDX_RES_SERV_POL_RES_SERV, tableName=RESOURCE_SERVER_POLICY; dropIndex indexName=IDX_RES_SRV_RES_RES_SRV, tableName=RESOURCE_SERVER_RESOURCE; dropIndex indexName=IDX_RES_SRV_SCOPE_RES_SRV, tableName=RESOURCE_SERVER_SCOPE','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-3.4.0.CR1-resource-server-pk-change-part3-fixed-nodropindex','glavoie@gmail.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2024-04-18 11:47:13',48,'EXECUTED','9:bdc99e567b3398bac83263d375aad143','addNotNullConstraint columnName=RESOURCE_SERVER_CLIENT_ID, tableName=RESOURCE_SERVER_POLICY; addNotNullConstraint columnName=RESOURCE_SERVER_CLIENT_ID, tableName=RESOURCE_SERVER_RESOURCE; addNotNullConstraint columnName=RESOURCE_SERVER_CLIENT_ID, ...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authn-3.4.0.CR1-refresh-token-max-reuse','glavoie@gmail.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2024-04-18 11:47:13',49,'EXECUTED','9:d198654156881c46bfba39abd7769e69','addColumn tableName=REALM','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.4.0','keycloak','META-INF/jpa-changelog-3.4.0.xml','2024-04-18 11:47:14',50,'EXECUTED','9:cfdd8736332ccdd72c5256ccb42335db','addPrimaryKey constraintName=CONSTRAINT_REALM_DEFAULT_ROLES, tableName=REALM_DEFAULT_ROLES; addPrimaryKey constraintName=CONSTRAINT_COMPOSITE_ROLE, tableName=COMPOSITE_ROLE; addPrimaryKey constraintName=CONSTR_REALM_DEFAULT_GROUPS, tableName=REALM...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.4.0-KEYCLOAK-5230','hmlnarik@redhat.com','META-INF/jpa-changelog-3.4.0.xml','2024-04-18 11:47:14',51,'EXECUTED','9:7c84de3d9bd84d7f077607c1a4dcb714','createIndex indexName=IDX_FU_ATTRIBUTE, tableName=FED_USER_ATTRIBUTE; createIndex indexName=IDX_FU_CONSENT, tableName=FED_USER_CONSENT; createIndex indexName=IDX_FU_CONSENT_RU, tableName=FED_USER_CONSENT; createIndex indexName=IDX_FU_CREDENTIAL, t...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.4.1','psilva@redhat.com','META-INF/jpa-changelog-3.4.1.xml','2024-04-18 11:47:14',52,'EXECUTED','9:5a6bb36cbefb6a9d6928452c0852af2d','modifyDataType columnName=VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.4.2','keycloak','META-INF/jpa-changelog-3.4.2.xml','2024-04-18 11:47:14',53,'EXECUTED','9:8f23e334dbc59f82e0a328373ca6ced0','update tableName=REALM','',NULL,'4.25.1',NULL,NULL,'3408415468'),('3.4.2-KEYCLOAK-5172','mkanis@redhat.com','META-INF/jpa-changelog-3.4.2.xml','2024-04-18 11:47:14',54,'EXECUTED','9:9156214268f09d970cdf0e1564d866af','update tableName=CLIENT','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.0.0-KEYCLOAK-6335','bburke@redhat.com','META-INF/jpa-changelog-4.0.0.xml','2024-04-18 11:47:14',55,'EXECUTED','9:db806613b1ed154826c02610b7dbdf74','createTable tableName=CLIENT_AUTH_FLOW_BINDINGS; addPrimaryKey constraintName=C_CLI_FLOW_BIND, tableName=CLIENT_AUTH_FLOW_BINDINGS','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.0.0-CLEANUP-UNUSED-TABLE','bburke@redhat.com','META-INF/jpa-changelog-4.0.0.xml','2024-04-18 11:47:14',56,'EXECUTED','9:229a041fb72d5beac76bb94a5fa709de','dropTable tableName=CLIENT_IDENTITY_PROV_MAPPING','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.0.0-KEYCLOAK-6228','bburke@redhat.com','META-INF/jpa-changelog-4.0.0.xml','2024-04-18 11:47:14',57,'EXECUTED','9:079899dade9c1e683f26b2aa9ca6ff04','dropUniqueConstraint constraintName=UK_JKUWUVD56ONTGSUHOGM8UEWRT, tableName=USER_CONSENT; dropNotNullConstraint columnName=CLIENT_ID, tableName=USER_CONSENT; addColumn tableName=USER_CONSENT; addUniqueConstraint constraintName=UK_JKUWUVD56ONTGSUHO...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.0.0-KEYCLOAK-5579-fixed','mposolda@redhat.com','META-INF/jpa-changelog-4.0.0.xml','2024-04-18 11:47:16',58,'EXECUTED','9:139b79bcbbfe903bb1c2d2a4dbf001d9','dropForeignKeyConstraint baseTableName=CLIENT_TEMPLATE_ATTRIBUTES, constraintName=FK_CL_TEMPL_ATTR_TEMPL; renameTable newTableName=CLIENT_SCOPE_ATTRIBUTES, oldTableName=CLIENT_TEMPLATE_ATTRIBUTES; renameColumn newColumnName=SCOPE_ID, oldColumnName...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-4.0.0.CR1','psilva@redhat.com','META-INF/jpa-changelog-authz-4.0.0.CR1.xml','2024-04-18 11:47:17',59,'EXECUTED','9:b55738ad889860c625ba2bf483495a04','createTable tableName=RESOURCE_SERVER_PERM_TICKET; addPrimaryKey constraintName=CONSTRAINT_FAPMT, tableName=RESOURCE_SERVER_PERM_TICKET; addForeignKeyConstraint baseTableName=RESOURCE_SERVER_PERM_TICKET, constraintName=FK_FRSRHO213XCX4WNKOG82SSPMT...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-4.0.0.Beta3','psilva@redhat.com','META-INF/jpa-changelog-authz-4.0.0.Beta3.xml','2024-04-18 11:47:17',60,'EXECUTED','9:e0057eac39aa8fc8e09ac6cfa4ae15fe','addColumn tableName=RESOURCE_SERVER_POLICY; addColumn tableName=RESOURCE_SERVER_PERM_TICKET; addForeignKeyConstraint baseTableName=RESOURCE_SERVER_PERM_TICKET, constraintName=FK_FRSRPO2128CX4WNKOG82SSRFY, referencedTableName=RESOURCE_SERVER_POLICY','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-4.2.0.Final','mhajas@redhat.com','META-INF/jpa-changelog-authz-4.2.0.Final.xml','2024-04-18 11:47:17',61,'EXECUTED','9:42a33806f3a0443fe0e7feeec821326c','createTable tableName=RESOURCE_URIS; addForeignKeyConstraint baseTableName=RESOURCE_URIS, constraintName=FK_RESOURCE_SERVER_URIS, referencedTableName=RESOURCE_SERVER_RESOURCE; customChange; dropColumn columnName=URI, tableName=RESOURCE_SERVER_RESO...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-4.2.0.Final-KEYCLOAK-9944','hmlnarik@redhat.com','META-INF/jpa-changelog-authz-4.2.0.Final.xml','2024-04-18 11:47:17',62,'EXECUTED','9:9968206fca46eecc1f51db9c024bfe56','addPrimaryKey constraintName=CONSTRAINT_RESOUR_URIS_PK, tableName=RESOURCE_URIS','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.2.0-KEYCLOAK-6313','wadahiro@gmail.com','META-INF/jpa-changelog-4.2.0.xml','2024-04-18 11:47:17',63,'EXECUTED','9:92143a6daea0a3f3b8f598c97ce55c3d','addColumn tableName=REQUIRED_ACTION_PROVIDER','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.3.0-KEYCLOAK-7984','wadahiro@gmail.com','META-INF/jpa-changelog-4.3.0.xml','2024-04-18 11:47:17',64,'EXECUTED','9:82bab26a27195d889fb0429003b18f40','update tableName=REQUIRED_ACTION_PROVIDER','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.6.0-KEYCLOAK-7950','psilva@redhat.com','META-INF/jpa-changelog-4.6.0.xml','2024-04-18 11:47:17',65,'EXECUTED','9:e590c88ddc0b38b0ae4249bbfcb5abc3','update tableName=RESOURCE_SERVER_RESOURCE','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.6.0-KEYCLOAK-8377','keycloak','META-INF/jpa-changelog-4.6.0.xml','2024-04-18 11:47:17',66,'EXECUTED','9:5c1f475536118dbdc38d5d7977950cc0','createTable tableName=ROLE_ATTRIBUTE; addPrimaryKey constraintName=CONSTRAINT_ROLE_ATTRIBUTE_PK, tableName=ROLE_ATTRIBUTE; addForeignKeyConstraint baseTableName=ROLE_ATTRIBUTE, constraintName=FK_ROLE_ATTRIBUTE_ID, referencedTableName=KEYCLOAK_ROLE...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.6.0-KEYCLOAK-8555','gideonray@gmail.com','META-INF/jpa-changelog-4.6.0.xml','2024-04-18 11:47:17',67,'EXECUTED','9:e7c9f5f9c4d67ccbbcc215440c718a17','createIndex indexName=IDX_COMPONENT_PROVIDER_TYPE, tableName=COMPONENT','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.7.0-KEYCLOAK-1267','sguilhen@redhat.com','META-INF/jpa-changelog-4.7.0.xml','2024-04-18 11:47:17',68,'EXECUTED','9:88e0bfdda924690d6f4e430c53447dd5','addColumn tableName=REALM','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.7.0-KEYCLOAK-7275','keycloak','META-INF/jpa-changelog-4.7.0.xml','2024-04-18 11:47:17',69,'EXECUTED','9:f53177f137e1c46b6a88c59ec1cb5218','renameColumn newColumnName=CREATED_ON, oldColumnName=LAST_SESSION_REFRESH, tableName=OFFLINE_USER_SESSION; addNotNullConstraint columnName=CREATED_ON, tableName=OFFLINE_USER_SESSION; addColumn tableName=OFFLINE_USER_SESSION; customChange; createIn...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('4.8.0-KEYCLOAK-8835','sguilhen@redhat.com','META-INF/jpa-changelog-4.8.0.xml','2024-04-18 11:47:18',70,'EXECUTED','9:a74d33da4dc42a37ec27121580d1459f','addNotNullConstraint columnName=SSO_MAX_LIFESPAN_REMEMBER_ME, tableName=REALM; addNotNullConstraint columnName=SSO_IDLE_TIMEOUT_REMEMBER_ME, tableName=REALM','',NULL,'4.25.1',NULL,NULL,'3408415468'),('authz-7.0.0-KEYCLOAK-10443','psilva@redhat.com','META-INF/jpa-changelog-authz-7.0.0.xml','2024-04-18 11:47:18',71,'EXECUTED','9:fd4ade7b90c3b67fae0bfcfcb42dfb5f','addColumn tableName=RESOURCE_SERVER','',NULL,'4.25.1',NULL,NULL,'3408415468'),('8.0.0-adding-credential-columns','keycloak','META-INF/jpa-changelog-8.0.0.xml','2024-04-18 11:47:18',72,'EXECUTED','9:aa072ad090bbba210d8f18781b8cebf4','addColumn tableName=CREDENTIAL; addColumn tableName=FED_USER_CREDENTIAL','',NULL,'4.25.1',NULL,NULL,'3408415468'),('8.0.0-updating-credential-data-not-oracle-fixed','keycloak','META-INF/jpa-changelog-8.0.0.xml','2024-04-18 11:47:18',73,'EXECUTED','9:1ae6be29bab7c2aa376f6983b932be37','update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=FED_USER_CREDENTIAL; update tableName=FED_USER_CREDENTIAL; update tableName=FED_USER_CREDENTIAL','',NULL,'4.25.1',NULL,NULL,'3408415468'),('8.0.0-updating-credential-data-oracle-fixed','keycloak','META-INF/jpa-changelog-8.0.0.xml','2024-04-18 11:47:18',74,'MARK_RAN','9:14706f286953fc9a25286dbd8fb30d97','update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=FED_USER_CREDENTIAL; update tableName=FED_USER_CREDENTIAL; update tableName=FED_USER_CREDENTIAL','',NULL,'4.25.1',NULL,NULL,'3408415468'),('8.0.0-credential-cleanup-fixed','keycloak','META-INF/jpa-changelog-8.0.0.xml','2024-04-18 11:47:18',75,'EXECUTED','9:2b9cc12779be32c5b40e2e67711a218b','dropDefaultValue columnName=COUNTER, tableName=CREDENTIAL; dropDefaultValue columnName=DIGITS, tableName=CREDENTIAL; dropDefaultValue columnName=PERIOD, tableName=CREDENTIAL; dropDefaultValue columnName=ALGORITHM, tableName=CREDENTIAL; dropColumn ...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('8.0.0-resource-tag-support','keycloak','META-INF/jpa-changelog-8.0.0.xml','2024-04-18 11:47:18',76,'EXECUTED','9:91fa186ce7a5af127a2d7a91ee083cc5','addColumn tableName=MIGRATION_MODEL; createIndex indexName=IDX_UPDATE_TIME, tableName=MIGRATION_MODEL','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.0-always-display-client','keycloak','META-INF/jpa-changelog-9.0.0.xml','2024-04-18 11:47:18',77,'EXECUTED','9:6335e5c94e83a2639ccd68dd24e2e5ad','addColumn tableName=CLIENT','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.0-drop-constraints-for-column-increase','keycloak','META-INF/jpa-changelog-9.0.0.xml','2024-04-18 11:47:18',78,'MARK_RAN','9:6bdb5658951e028bfe16fa0a8228b530','dropUniqueConstraint constraintName=UK_FRSR6T700S9V50BU18WS5PMT, tableName=RESOURCE_SERVER_PERM_TICKET; dropUniqueConstraint constraintName=UK_FRSR6T700S9V50BU18WS5HA6, tableName=RESOURCE_SERVER_RESOURCE; dropPrimaryKey constraintName=CONSTRAINT_O...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.0-increase-column-size-federated-fk','keycloak','META-INF/jpa-changelog-9.0.0.xml','2024-04-18 11:47:19',79,'EXECUTED','9:d5bc15a64117ccad481ce8792d4c608f','modifyDataType columnName=CLIENT_ID, tableName=FED_USER_CONSENT; modifyDataType columnName=CLIENT_REALM_CONSTRAINT, tableName=KEYCLOAK_ROLE; modifyDataType columnName=OWNER, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=CLIENT_ID, ta...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.0-recreate-constraints-after-column-increase','keycloak','META-INF/jpa-changelog-9.0.0.xml','2024-04-18 11:47:19',80,'MARK_RAN','9:077cba51999515f4d3e7ad5619ab592c','addNotNullConstraint columnName=CLIENT_ID, tableName=OFFLINE_CLIENT_SESSION; addNotNullConstraint columnName=OWNER, tableName=RESOURCE_SERVER_PERM_TICKET; addNotNullConstraint columnName=REQUESTER, tableName=RESOURCE_SERVER_PERM_TICKET; addNotNull...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.1-add-index-to-client.client_id','keycloak','META-INF/jpa-changelog-9.0.1.xml','2024-04-18 11:47:19',81,'EXECUTED','9:be969f08a163bf47c6b9e9ead8ac2afb','createIndex indexName=IDX_CLIENT_ID, tableName=CLIENT','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.1-KEYCLOAK-12579-drop-constraints','keycloak','META-INF/jpa-changelog-9.0.1.xml','2024-04-18 11:47:19',82,'MARK_RAN','9:6d3bb4408ba5a72f39bd8a0b301ec6e3','dropUniqueConstraint constraintName=SIBLING_NAMES, tableName=KEYCLOAK_GROUP','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.1-KEYCLOAK-12579-add-not-null-constraint','keycloak','META-INF/jpa-changelog-9.0.1.xml','2024-04-18 11:47:19',83,'EXECUTED','9:966bda61e46bebf3cc39518fbed52fa7','addNotNullConstraint columnName=PARENT_GROUP, tableName=KEYCLOAK_GROUP','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.1-KEYCLOAK-12579-recreate-constraints','keycloak','META-INF/jpa-changelog-9.0.1.xml','2024-04-18 11:47:19',84,'MARK_RAN','9:8dcac7bdf7378e7d823cdfddebf72fda','addUniqueConstraint constraintName=SIBLING_NAMES, tableName=KEYCLOAK_GROUP','',NULL,'4.25.1',NULL,NULL,'3408415468'),('9.0.1-add-index-to-events','keycloak','META-INF/jpa-changelog-9.0.1.xml','2024-04-18 11:47:19',85,'EXECUTED','9:7d93d602352a30c0c317e6a609b56599','createIndex indexName=IDX_EVENT_TIME, tableName=EVENT_ENTITY','',NULL,'4.25.1',NULL,NULL,'3408415468'),('map-remove-ri','keycloak','META-INF/jpa-changelog-11.0.0.xml','2024-04-18 11:47:19',86,'EXECUTED','9:71c5969e6cdd8d7b6f47cebc86d37627','dropForeignKeyConstraint baseTableName=REALM, constraintName=FK_TRAF444KK6QRKMS7N56AIWQ5Y; dropForeignKeyConstraint baseTableName=KEYCLOAK_ROLE, constraintName=FK_KJHO5LE2C0RAL09FL8CM9WFW9','',NULL,'4.25.1',NULL,NULL,'3408415468'),('map-remove-ri','keycloak','META-INF/jpa-changelog-12.0.0.xml','2024-04-18 11:47:19',87,'EXECUTED','9:a9ba7d47f065f041b7da856a81762021','dropForeignKeyConstraint baseTableName=REALM_DEFAULT_GROUPS, constraintName=FK_DEF_GROUPS_GROUP; dropForeignKeyConstraint baseTableName=REALM_DEFAULT_ROLES, constraintName=FK_H4WPD7W4HSOOLNI3H0SW7BTJE; dropForeignKeyConstraint baseTableName=CLIENT...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('12.1.0-add-realm-localization-table','keycloak','META-INF/jpa-changelog-12.0.0.xml','2024-04-18 11:47:19',88,'EXECUTED','9:fffabce2bc01e1a8f5110d5278500065','createTable tableName=REALM_LOCALIZATIONS; addPrimaryKey tableName=REALM_LOCALIZATIONS','',NULL,'4.25.1',NULL,NULL,'3408415468'),('default-roles','keycloak','META-INF/jpa-changelog-13.0.0.xml','2024-04-18 11:47:19',89,'EXECUTED','9:fa8a5b5445e3857f4b010bafb5009957','addColumn tableName=REALM; customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('default-roles-cleanup','keycloak','META-INF/jpa-changelog-13.0.0.xml','2024-04-18 11:47:19',90,'EXECUTED','9:67ac3241df9a8582d591c5ed87125f39','dropTable tableName=REALM_DEFAULT_ROLES; dropTable tableName=CLIENT_DEFAULT_ROLES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('13.0.0-KEYCLOAK-16844','keycloak','META-INF/jpa-changelog-13.0.0.xml','2024-04-18 11:47:19',91,'EXECUTED','9:ad1194d66c937e3ffc82386c050ba089','createIndex indexName=IDX_OFFLINE_USS_PRELOAD, tableName=OFFLINE_USER_SESSION','',NULL,'4.25.1',NULL,NULL,'3408415468'),('map-remove-ri-13.0.0','keycloak','META-INF/jpa-changelog-13.0.0.xml','2024-04-18 11:47:20',92,'EXECUTED','9:d9be619d94af5a2f5d07b9f003543b91','dropForeignKeyConstraint baseTableName=DEFAULT_CLIENT_SCOPE, constraintName=FK_R_DEF_CLI_SCOPE_SCOPE; dropForeignKeyConstraint baseTableName=CLIENT_SCOPE_CLIENT, constraintName=FK_C_CLI_SCOPE_SCOPE; dropForeignKeyConstraint baseTableName=CLIENT_SC...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('13.0.0-KEYCLOAK-17992-drop-constraints','keycloak','META-INF/jpa-changelog-13.0.0.xml','2024-04-18 11:47:20',93,'MARK_RAN','9:544d201116a0fcc5a5da0925fbbc3bde','dropPrimaryKey constraintName=C_CLI_SCOPE_BIND, tableName=CLIENT_SCOPE_CLIENT; dropIndex indexName=IDX_CLSCOPE_CL, tableName=CLIENT_SCOPE_CLIENT; dropIndex indexName=IDX_CL_CLSCOPE, tableName=CLIENT_SCOPE_CLIENT','',NULL,'4.25.1',NULL,NULL,'3408415468'),('13.0.0-increase-column-size-federated','keycloak','META-INF/jpa-changelog-13.0.0.xml','2024-04-18 11:47:20',94,'EXECUTED','9:43c0c1055b6761b4b3e89de76d612ccf','modifyDataType columnName=CLIENT_ID, tableName=CLIENT_SCOPE_CLIENT; modifyDataType columnName=SCOPE_ID, tableName=CLIENT_SCOPE_CLIENT','',NULL,'4.25.1',NULL,NULL,'3408415468'),('13.0.0-KEYCLOAK-17992-recreate-constraints','keycloak','META-INF/jpa-changelog-13.0.0.xml','2024-04-18 11:47:20',95,'MARK_RAN','9:8bd711fd0330f4fe980494ca43ab1139','addNotNullConstraint columnName=CLIENT_ID, tableName=CLIENT_SCOPE_CLIENT; addNotNullConstraint columnName=SCOPE_ID, tableName=CLIENT_SCOPE_CLIENT; addPrimaryKey constraintName=C_CLI_SCOPE_BIND, tableName=CLIENT_SCOPE_CLIENT; createIndex indexName=...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('json-string-accomodation-fixed','keycloak','META-INF/jpa-changelog-13.0.0.xml','2024-04-18 11:47:20',96,'EXECUTED','9:e07d2bc0970c348bb06fb63b1f82ddbf','addColumn tableName=REALM_ATTRIBUTE; update tableName=REALM_ATTRIBUTE; dropColumn columnName=VALUE, tableName=REALM_ATTRIBUTE; renameColumn newColumnName=VALUE, oldColumnName=VALUE_NEW, tableName=REALM_ATTRIBUTE','',NULL,'4.25.1',NULL,NULL,'3408415468'),('14.0.0-KEYCLOAK-11019','keycloak','META-INF/jpa-changelog-14.0.0.xml','2024-04-18 11:47:20',97,'EXECUTED','9:24fb8611e97f29989bea412aa38d12b7','createIndex indexName=IDX_OFFLINE_CSS_PRELOAD, tableName=OFFLINE_CLIENT_SESSION; createIndex indexName=IDX_OFFLINE_USS_BY_USER, tableName=OFFLINE_USER_SESSION; createIndex indexName=IDX_OFFLINE_USS_BY_USERSESS, tableName=OFFLINE_USER_SESSION','',NULL,'4.25.1',NULL,NULL,'3408415468'),('14.0.0-KEYCLOAK-18286','keycloak','META-INF/jpa-changelog-14.0.0.xml','2024-04-18 11:47:20',98,'MARK_RAN','9:259f89014ce2506ee84740cbf7163aa7','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('14.0.0-KEYCLOAK-18286-revert','keycloak','META-INF/jpa-changelog-14.0.0.xml','2024-04-18 11:47:20',99,'MARK_RAN','9:04baaf56c116ed19951cbc2cca584022','dropIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('14.0.0-KEYCLOAK-18286-supported-dbs','keycloak','META-INF/jpa-changelog-14.0.0.xml','2024-04-18 11:47:20',100,'EXECUTED','9:bd2bd0fc7768cf0845ac96a8786fa735','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('14.0.0-KEYCLOAK-18286-unsupported-dbs','keycloak','META-INF/jpa-changelog-14.0.0.xml','2024-04-18 11:47:20',101,'MARK_RAN','9:d3d977031d431db16e2c181ce49d73e9','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('KEYCLOAK-17267-add-index-to-user-attributes','keycloak','META-INF/jpa-changelog-14.0.0.xml','2024-04-18 11:47:20',102,'EXECUTED','9:0b305d8d1277f3a89a0a53a659ad274c','createIndex indexName=IDX_USER_ATTRIBUTE_NAME, tableName=USER_ATTRIBUTE','',NULL,'4.25.1',NULL,NULL,'3408415468'),('KEYCLOAK-18146-add-saml-art-binding-identifier','keycloak','META-INF/jpa-changelog-14.0.0.xml','2024-04-18 11:47:20',103,'EXECUTED','9:2c374ad2cdfe20e2905a84c8fac48460','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('15.0.0-KEYCLOAK-18467','keycloak','META-INF/jpa-changelog-15.0.0.xml','2024-04-18 11:47:20',104,'EXECUTED','9:47a760639ac597360a8219f5b768b4de','addColumn tableName=REALM_LOCALIZATIONS; update tableName=REALM_LOCALIZATIONS; dropColumn columnName=TEXTS, tableName=REALM_LOCALIZATIONS; renameColumn newColumnName=TEXTS, oldColumnName=TEXTS_NEW, tableName=REALM_LOCALIZATIONS; addNotNullConstrai...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('17.0.0-9562','keycloak','META-INF/jpa-changelog-17.0.0.xml','2024-04-18 11:47:20',105,'EXECUTED','9:a6272f0576727dd8cad2522335f5d99e','createIndex indexName=IDX_USER_SERVICE_ACCOUNT, tableName=USER_ENTITY','',NULL,'4.25.1',NULL,NULL,'3408415468'),('18.0.0-10625-IDX_ADMIN_EVENT_TIME','keycloak','META-INF/jpa-changelog-18.0.0.xml','2024-04-18 11:47:20',106,'EXECUTED','9:015479dbd691d9cc8669282f4828c41d','createIndex indexName=IDX_ADMIN_EVENT_TIME, tableName=ADMIN_EVENT_ENTITY','',NULL,'4.25.1',NULL,NULL,'3408415468'),('19.0.0-10135','keycloak','META-INF/jpa-changelog-19.0.0.xml','2024-04-18 11:47:20',107,'EXECUTED','9:9518e495fdd22f78ad6425cc30630221','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('20.0.0-12964-supported-dbs','keycloak','META-INF/jpa-changelog-20.0.0.xml','2024-04-18 11:47:20',108,'EXECUTED','9:f2e1331a71e0aa85e5608fe42f7f681c','createIndex indexName=IDX_GROUP_ATT_BY_NAME_VALUE, tableName=GROUP_ATTRIBUTE','',NULL,'4.25.1',NULL,NULL,'3408415468'),('20.0.0-12964-unsupported-dbs','keycloak','META-INF/jpa-changelog-20.0.0.xml','2024-04-18 11:47:20',109,'MARK_RAN','9:1a6fcaa85e20bdeae0a9ce49b41946a5','createIndex indexName=IDX_GROUP_ATT_BY_NAME_VALUE, tableName=GROUP_ATTRIBUTE','',NULL,'4.25.1',NULL,NULL,'3408415468'),('client-attributes-string-accomodation-fixed','keycloak','META-INF/jpa-changelog-20.0.0.xml','2024-04-18 11:47:20',110,'EXECUTED','9:3f332e13e90739ed0c35b0b25b7822ca','addColumn tableName=CLIENT_ATTRIBUTES; update tableName=CLIENT_ATTRIBUTES; dropColumn columnName=VALUE, tableName=CLIENT_ATTRIBUTES; renameColumn newColumnName=VALUE, oldColumnName=VALUE_NEW, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('21.0.2-17277','keycloak','META-INF/jpa-changelog-21.0.2.xml','2024-04-18 11:47:20',111,'EXECUTED','9:7ee1f7a3fb8f5588f171fb9a6ab623c0','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('21.1.0-19404','keycloak','META-INF/jpa-changelog-21.1.0.xml','2024-04-18 11:47:20',112,'EXECUTED','9:3d7e830b52f33676b9d64f7f2b2ea634','modifyDataType columnName=DECISION_STRATEGY, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=LOGIC, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=POLICY_ENFORCE_MODE, tableName=RESOURCE_SERVER','',NULL,'4.25.1',NULL,NULL,'3408415468'),('21.1.0-19404-2','keycloak','META-INF/jpa-changelog-21.1.0.xml','2024-04-18 11:47:20',113,'MARK_RAN','9:627d032e3ef2c06c0e1f73d2ae25c26c','addColumn tableName=RESOURCE_SERVER_POLICY; update tableName=RESOURCE_SERVER_POLICY; dropColumn columnName=DECISION_STRATEGY, tableName=RESOURCE_SERVER_POLICY; renameColumn newColumnName=DECISION_STRATEGY, oldColumnName=DECISION_STRATEGY_NEW, tabl...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('22.0.0-17484-updated','keycloak','META-INF/jpa-changelog-22.0.0.xml','2024-04-18 11:47:20',114,'EXECUTED','9:90af0bfd30cafc17b9f4d6eccd92b8b3','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('22.0.5-24031','keycloak','META-INF/jpa-changelog-22.0.0.xml','2024-04-18 11:47:20',115,'MARK_RAN','9:a60d2d7b315ec2d3eba9e2f145f9df28','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('23.0.0-12062','keycloak','META-INF/jpa-changelog-23.0.0.xml','2024-04-18 11:47:20',116,'EXECUTED','9:2168fbe728fec46ae9baf15bf80927b8','addColumn tableName=COMPONENT_CONFIG; update tableName=COMPONENT_CONFIG; dropColumn columnName=VALUE, tableName=COMPONENT_CONFIG; renameColumn newColumnName=VALUE, oldColumnName=VALUE_NEW, tableName=COMPONENT_CONFIG','',NULL,'4.25.1',NULL,NULL,'3408415468'),('23.0.0-17258','keycloak','META-INF/jpa-changelog-23.0.0.xml','2024-04-18 11:47:20',117,'EXECUTED','9:36506d679a83bbfda85a27ea1864dca8','addColumn tableName=EVENT_ENTITY','',NULL,'4.25.1',NULL,NULL,'3408415468'),('24.0.0-9758','keycloak','META-INF/jpa-changelog-24.0.0.xml','2024-04-18 11:47:21',118,'EXECUTED','9:502c557a5189f600f0f445a9b49ebbce','addColumn tableName=USER_ATTRIBUTE; addColumn tableName=FED_USER_ATTRIBUTE; createIndex indexName=USER_ATTR_LONG_VALUES, tableName=USER_ATTRIBUTE; createIndex indexName=FED_USER_ATTR_LONG_VALUES, tableName=FED_USER_ATTRIBUTE; createIndex indexName...','',NULL,'4.25.1',NULL,NULL,'3408415468'),('24.0.0-9758-2','keycloak','META-INF/jpa-changelog-24.0.0.xml','2024-04-18 11:47:21',119,'EXECUTED','9:bf0fdee10afdf597a987adbf291db7b2','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('24.0.0-26618-drop-index-if-present','keycloak','META-INF/jpa-changelog-24.0.0.xml','2024-04-18 11:47:21',120,'EXECUTED','9:04baaf56c116ed19951cbc2cca584022','dropIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('24.0.0-26618-reindex','keycloak','META-INF/jpa-changelog-24.0.0.xml','2024-04-18 11:47:21',121,'EXECUTED','9:bd2bd0fc7768cf0845ac96a8786fa735','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('24.0.2-27228','keycloak','META-INF/jpa-changelog-24.0.2.xml','2024-04-18 11:47:21',122,'EXECUTED','9:eaee11f6b8aa25d2cc6a84fb86fc6238','customChange','',NULL,'4.25.1',NULL,NULL,'3408415468'),('24.0.2-27967-drop-index-if-present','keycloak','META-INF/jpa-changelog-24.0.2.xml','2024-04-18 11:47:21',123,'MARK_RAN','9:04baaf56c116ed19951cbc2cca584022','dropIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468'),('24.0.2-27967-reindex','keycloak','META-INF/jpa-changelog-24.0.2.xml','2024-04-18 11:47:21',124,'MARK_RAN','9:d3d977031d431db16e2c181ce49d73e9','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.25.1',NULL,NULL,'3408415468');
/*!40000 ALTER TABLE `DATABASECHANGELOG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `DATABASECHANGELOGLOCK`
--

DROP TABLE IF EXISTS `DATABASECHANGELOGLOCK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DATABASECHANGELOGLOCK` (
  `ID` int NOT NULL,
  `LOCKED` tinyint NOT NULL,
  `LOCKGRANTED` datetime DEFAULT NULL,
  `LOCKEDBY` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `DATABASECHANGELOGLOCK`
--

LOCK TABLES `DATABASECHANGELOGLOCK` WRITE;
/*!40000 ALTER TABLE `DATABASECHANGELOGLOCK` DISABLE KEYS */;
INSERT INTO `DATABASECHANGELOGLOCK` (`ID`, `LOCKED`, `LOCKGRANTED`, `LOCKEDBY`) VALUES (1,0,NULL,NULL),(1000,0,NULL,NULL),(1001,0,NULL,NULL);
/*!40000 ALTER TABLE `DATABASECHANGELOGLOCK` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `DEFAULT_CLIENT_SCOPE`
--

DROP TABLE IF EXISTS `DEFAULT_CLIENT_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DEFAULT_CLIENT_SCOPE` (
  `REALM_ID` varchar(36) NOT NULL,
  `SCOPE_ID` varchar(36) NOT NULL,
  `DEFAULT_SCOPE` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`REALM_ID`,`SCOPE_ID`),
  KEY `IDX_DEFCLS_REALM` (`REALM_ID`),
  KEY `IDX_DEFCLS_SCOPE` (`SCOPE_ID`),
  CONSTRAINT `FK_R_DEF_CLI_SCOPE_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `DEFAULT_CLIENT_SCOPE`
--

LOCK TABLES `DEFAULT_CLIENT_SCOPE` WRITE;
/*!40000 ALTER TABLE `DEFAULT_CLIENT_SCOPE` DISABLE KEYS */;
INSERT INTO `DEFAULT_CLIENT_SCOPE` (`REALM_ID`, `SCOPE_ID`, `DEFAULT_SCOPE`) VALUES ('40ae881c-f4e4-4b07-b097-a67d2bf515e6','14d13053-c1b9-43f8-b993-9d7a8aef8069',1),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','192eb290-9206-4ec7-94ce-f62e6ab82179',0),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','282b6071-c081-4ebf-bb56-5241153f1811',0),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','310fbe9a-4f1f-4683-aca4-516d640a5d9b',0),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','40f0b7c4-e882-42e9-8cc2-00055b119ca8',1),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','6bffc256-3149-4d67-a5d2-888076a43a46',1),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','6eed7521-b56d-4af5-8927-771c01723b21',0),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','815032a2-e511-4d5f-a543-e9a3ee4f648c',1),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','aaa43c0a-2c37-42b7-a80c-f98eef322343',1),('40ae881c-f4e4-4b07-b097-a67d2bf515e6','b9388bd3-40d7-4f4a-87e8-2d78275ee434',1),('4327ba47-4116-44ea-9c4d-02907dca81e7','066554e5-ed72-4255-a59e-37bce592656c',0),('4327ba47-4116-44ea-9c4d-02907dca81e7','091b2057-69d2-4a0e-a76c-7cbd76698850',1),('4327ba47-4116-44ea-9c4d-02907dca81e7','15afd007-2320-4439-853b-2ddaa8b2ff71',0),('4327ba47-4116-44ea-9c4d-02907dca81e7','39474b38-4002-4590-a90e-7e77c25785b5',1),('4327ba47-4116-44ea-9c4d-02907dca81e7','3e744bf6-ec28-46cf-af95-5cf3f4ee8d58',1),('4327ba47-4116-44ea-9c4d-02907dca81e7','90977ab5-207e-4b07-836d-e8d59805927b',1),('4327ba47-4116-44ea-9c4d-02907dca81e7','c294c451-1979-493f-ad49-15220c9bd6f8',1),('4327ba47-4116-44ea-9c4d-02907dca81e7','c36c9c18-f3c5-451b-bf26-2887df9d803d',0),('4327ba47-4116-44ea-9c4d-02907dca81e7','ef9560b3-735c-4882-9aa1-dce8fe0697af',1),('4327ba47-4116-44ea-9c4d-02907dca81e7','fd09aae5-67c6-4162-a03e-95f8cc07d053',0),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','29003608-50c2-4693-953b-316f0c71fa25',1),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','3d4aaffa-8399-4bff-9dac-000e43b45ca9',1),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','65b407fe-a29a-4caa-90c1-060588438771',0),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','74beed12-a014-4842-91e7-5334c22a5bc0',1),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','7b474ad1-7a28-4e53-970d-e9a5a14bfbab',1),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','9c93a802-d9c8-468e-8983-129fc287b26c',1),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','bd4873a9-1ad5-40ba-bec3-5506355835e3',0),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','bdb21a59-753f-498c-8be2-9fae1f633cd3',1),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70',0),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','fafceb53-8f99-4a18-81bf-59c2b83bcd6c',0);
/*!40000 ALTER TABLE `DEFAULT_CLIENT_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `EVENT_ENTITY`
--

DROP TABLE IF EXISTS `EVENT_ENTITY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `EVENT_ENTITY` (
  `ID` varchar(36) NOT NULL,
  `CLIENT_ID` varchar(255) DEFAULT NULL,
  `DETAILS_JSON` text,
  `ERROR` varchar(255) DEFAULT NULL,
  `IP_ADDRESS` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(255) DEFAULT NULL,
  `SESSION_ID` varchar(255) DEFAULT NULL,
  `EVENT_TIME` bigint DEFAULT NULL,
  `TYPE` varchar(255) DEFAULT NULL,
  `USER_ID` varchar(255) DEFAULT NULL,
  `DETAILS_JSON_LONG_VALUE` longtext CHARACTER SET utf8 COLLATE utf8_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_EVENT_TIME` (`REALM_ID`,`EVENT_TIME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `EVENT_ENTITY`
--

LOCK TABLES `EVENT_ENTITY` WRITE;
/*!40000 ALTER TABLE `EVENT_ENTITY` DISABLE KEYS */;
/*!40000 ALTER TABLE `EVENT_ENTITY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FEDERATED_IDENTITY`
--

DROP TABLE IF EXISTS `FEDERATED_IDENTITY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FEDERATED_IDENTITY` (
  `IDENTITY_PROVIDER` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  `FEDERATED_USER_ID` varchar(255) DEFAULT NULL,
  `FEDERATED_USERNAME` varchar(255) DEFAULT NULL,
  `TOKEN` text,
  `USER_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`IDENTITY_PROVIDER`,`USER_ID`),
  KEY `IDX_FEDIDENTITY_USER` (`USER_ID`),
  KEY `IDX_FEDIDENTITY_FEDUSER` (`FEDERATED_USER_ID`),
  CONSTRAINT `FK404288B92EF007A6` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FEDERATED_IDENTITY`
--

LOCK TABLES `FEDERATED_IDENTITY` WRITE;
/*!40000 ALTER TABLE `FEDERATED_IDENTITY` DISABLE KEYS */;
INSERT INTO `FEDERATED_IDENTITY` (`IDENTITY_PROVIDER`, `REALM_ID`, `FEDERATED_USER_ID`, `FEDERATED_USERNAME`, `TOKEN`, `USER_ID`) VALUES ('github','40ae881c-f4e4-4b07-b097-a67d2bf515e6','3415','hwi-han',NULL,'a80e2adf-05a7-4caf-9411-82efd36d23db');
/*!40000 ALTER TABLE `FEDERATED_IDENTITY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FEDERATED_USER`
--

DROP TABLE IF EXISTS `FEDERATED_USER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FEDERATED_USER` (
  `ID` varchar(255) NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FEDERATED_USER`
--

LOCK TABLES `FEDERATED_USER` WRITE;
/*!40000 ALTER TABLE `FEDERATED_USER` DISABLE KEYS */;
/*!40000 ALTER TABLE `FEDERATED_USER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FED_USER_ATTRIBUTE`
--

DROP TABLE IF EXISTS `FED_USER_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_ATTRIBUTE` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `USER_ID` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) DEFAULT NULL,
  `VALUE` text,
  `LONG_VALUE_HASH` binary(64) DEFAULT NULL,
  `LONG_VALUE_HASH_LOWER_CASE` binary(64) DEFAULT NULL,
  `LONG_VALUE` longtext CHARACTER SET utf8 COLLATE utf8_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_FU_ATTRIBUTE` (`USER_ID`,`REALM_ID`,`NAME`),
  KEY `FED_USER_ATTR_LONG_VALUES` (`LONG_VALUE_HASH`,`NAME`),
  KEY `FED_USER_ATTR_LONG_VALUES_LOWER_CASE` (`LONG_VALUE_HASH_LOWER_CASE`,`NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FED_USER_ATTRIBUTE`
--

LOCK TABLES `FED_USER_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `FED_USER_ATTRIBUTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FED_USER_CONSENT`
--

DROP TABLE IF EXISTS `FED_USER_CONSENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_CONSENT` (
  `ID` varchar(36) NOT NULL,
  `CLIENT_ID` varchar(255) DEFAULT NULL,
  `USER_ID` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) DEFAULT NULL,
  `CREATED_DATE` bigint DEFAULT NULL,
  `LAST_UPDATED_DATE` bigint DEFAULT NULL,
  `CLIENT_STORAGE_PROVIDER` varchar(36) DEFAULT NULL,
  `EXTERNAL_CLIENT_ID` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_FU_CONSENT` (`USER_ID`,`CLIENT_ID`),
  KEY `IDX_FU_CONSENT_RU` (`REALM_ID`,`USER_ID`),
  KEY `IDX_FU_CNSNT_EXT` (`USER_ID`,`CLIENT_STORAGE_PROVIDER`,`EXTERNAL_CLIENT_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FED_USER_CONSENT`
--

LOCK TABLES `FED_USER_CONSENT` WRITE;
/*!40000 ALTER TABLE `FED_USER_CONSENT` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_CONSENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FED_USER_CONSENT_CL_SCOPE`
--

DROP TABLE IF EXISTS `FED_USER_CONSENT_CL_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_CONSENT_CL_SCOPE` (
  `USER_CONSENT_ID` varchar(36) NOT NULL,
  `SCOPE_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`USER_CONSENT_ID`,`SCOPE_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FED_USER_CONSENT_CL_SCOPE`
--

LOCK TABLES `FED_USER_CONSENT_CL_SCOPE` WRITE;
/*!40000 ALTER TABLE `FED_USER_CONSENT_CL_SCOPE` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_CONSENT_CL_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FED_USER_CREDENTIAL`
--

DROP TABLE IF EXISTS `FED_USER_CREDENTIAL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_CREDENTIAL` (
  `ID` varchar(36) NOT NULL,
  `SALT` tinyblob,
  `TYPE` varchar(255) DEFAULT NULL,
  `CREATED_DATE` bigint DEFAULT NULL,
  `USER_ID` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) DEFAULT NULL,
  `USER_LABEL` varchar(255) DEFAULT NULL,
  `SECRET_DATA` longtext,
  `CREDENTIAL_DATA` longtext,
  `PRIORITY` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_FU_CREDENTIAL` (`USER_ID`,`TYPE`),
  KEY `IDX_FU_CREDENTIAL_RU` (`REALM_ID`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FED_USER_CREDENTIAL`
--

LOCK TABLES `FED_USER_CREDENTIAL` WRITE;
/*!40000 ALTER TABLE `FED_USER_CREDENTIAL` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_CREDENTIAL` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FED_USER_GROUP_MEMBERSHIP`
--

DROP TABLE IF EXISTS `FED_USER_GROUP_MEMBERSHIP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_GROUP_MEMBERSHIP` (
  `GROUP_ID` varchar(36) NOT NULL,
  `USER_ID` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`GROUP_ID`,`USER_ID`),
  KEY `IDX_FU_GROUP_MEMBERSHIP` (`USER_ID`,`GROUP_ID`),
  KEY `IDX_FU_GROUP_MEMBERSHIP_RU` (`REALM_ID`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FED_USER_GROUP_MEMBERSHIP`
--

LOCK TABLES `FED_USER_GROUP_MEMBERSHIP` WRITE;
/*!40000 ALTER TABLE `FED_USER_GROUP_MEMBERSHIP` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_GROUP_MEMBERSHIP` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FED_USER_REQUIRED_ACTION`
--

DROP TABLE IF EXISTS `FED_USER_REQUIRED_ACTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_REQUIRED_ACTION` (
  `REQUIRED_ACTION` varchar(255) NOT NULL DEFAULT ' ',
  `USER_ID` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`REQUIRED_ACTION`,`USER_ID`),
  KEY `IDX_FU_REQUIRED_ACTION` (`USER_ID`,`REQUIRED_ACTION`),
  KEY `IDX_FU_REQUIRED_ACTION_RU` (`REALM_ID`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FED_USER_REQUIRED_ACTION`
--

LOCK TABLES `FED_USER_REQUIRED_ACTION` WRITE;
/*!40000 ALTER TABLE `FED_USER_REQUIRED_ACTION` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_REQUIRED_ACTION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FED_USER_ROLE_MAPPING`
--

DROP TABLE IF EXISTS `FED_USER_ROLE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_ROLE_MAPPING` (
  `ROLE_ID` varchar(36) NOT NULL,
  `USER_ID` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ROLE_ID`,`USER_ID`),
  KEY `IDX_FU_ROLE_MAPPING` (`USER_ID`,`ROLE_ID`),
  KEY `IDX_FU_ROLE_MAPPING_RU` (`REALM_ID`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FED_USER_ROLE_MAPPING`
--

LOCK TABLES `FED_USER_ROLE_MAPPING` WRITE;
/*!40000 ALTER TABLE `FED_USER_ROLE_MAPPING` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_ROLE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `GROUP_ATTRIBUTE`
--

DROP TABLE IF EXISTS `GROUP_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `GROUP_ATTRIBUTE` (
  `ID` varchar(36) NOT NULL DEFAULT 'sybase-needs-something-here',
  `NAME` varchar(255) NOT NULL,
  `VALUE` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `GROUP_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_GROUP_ATTR_GROUP` (`GROUP_ID`),
  KEY `IDX_GROUP_ATT_BY_NAME_VALUE` (`NAME`,`VALUE`),
  CONSTRAINT `FK_GROUP_ATTRIBUTE_GROUP` FOREIGN KEY (`GROUP_ID`) REFERENCES `KEYCLOAK_GROUP` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `GROUP_ATTRIBUTE`
--

LOCK TABLES `GROUP_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `GROUP_ATTRIBUTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `GROUP_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `GROUP_ROLE_MAPPING`
--

DROP TABLE IF EXISTS `GROUP_ROLE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `GROUP_ROLE_MAPPING` (
  `ROLE_ID` varchar(36) NOT NULL,
  `GROUP_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`ROLE_ID`,`GROUP_ID`),
  KEY `IDX_GROUP_ROLE_MAPP_GROUP` (`GROUP_ID`),
  CONSTRAINT `FK_GROUP_ROLE_GROUP` FOREIGN KEY (`GROUP_ID`) REFERENCES `KEYCLOAK_GROUP` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `GROUP_ROLE_MAPPING`
--

LOCK TABLES `GROUP_ROLE_MAPPING` WRITE;
/*!40000 ALTER TABLE `GROUP_ROLE_MAPPING` DISABLE KEYS */;
/*!40000 ALTER TABLE `GROUP_ROLE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `IDENTITY_PROVIDER`
--

DROP TABLE IF EXISTS `IDENTITY_PROVIDER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `IDENTITY_PROVIDER` (
  `INTERNAL_ID` varchar(36) NOT NULL,
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `PROVIDER_ALIAS` varchar(255) DEFAULT NULL,
  `PROVIDER_ID` varchar(255) DEFAULT NULL,
  `STORE_TOKEN` tinyint NOT NULL DEFAULT '0',
  `AUTHENTICATE_BY_DEFAULT` tinyint NOT NULL DEFAULT '0',
  `REALM_ID` varchar(36) DEFAULT NULL,
  `ADD_TOKEN_ROLE` tinyint NOT NULL DEFAULT '1',
  `TRUST_EMAIL` tinyint NOT NULL DEFAULT '0',
  `FIRST_BROKER_LOGIN_FLOW_ID` varchar(36) DEFAULT NULL,
  `POST_BROKER_LOGIN_FLOW_ID` varchar(36) DEFAULT NULL,
  `PROVIDER_DISPLAY_NAME` varchar(255) DEFAULT NULL,
  `LINK_ONLY` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`INTERNAL_ID`),
  UNIQUE KEY `UK_2DAELWNIBJI49AVXSRTUF6XJ33` (`PROVIDER_ALIAS`,`REALM_ID`),
  KEY `IDX_IDENT_PROV_REALM` (`REALM_ID`),
  CONSTRAINT `FK2B4EBC52AE5C3B34` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `IDENTITY_PROVIDER`
--

LOCK TABLES `IDENTITY_PROVIDER` WRITE;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER` DISABLE KEYS */;
INSERT INTO `IDENTITY_PROVIDER` (`INTERNAL_ID`, `ENABLED`, `PROVIDER_ALIAS`, `PROVIDER_ID`, `STORE_TOKEN`, `AUTHENTICATE_BY_DEFAULT`, `REALM_ID`, `ADD_TOKEN_ROLE`, `TRUST_EMAIL`, `FIRST_BROKER_LOGIN_FLOW_ID`, `POST_BROKER_LOGIN_FLOW_ID`, `PROVIDER_DISPLAY_NAME`, `LINK_ONLY`) VALUES ('9c52c271-06a9-4133-92a2-779e0c011b91',1,'github','github',0,0,'40ae881c-f4e4-4b07-b097-a67d2bf515e6',0,0,'e4be6c41-61ae-4230-a920-d66ab72fd11a',NULL,NULL,0);
/*!40000 ALTER TABLE `IDENTITY_PROVIDER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `IDENTITY_PROVIDER_CONFIG`
--

DROP TABLE IF EXISTS `IDENTITY_PROVIDER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `IDENTITY_PROVIDER_CONFIG` (
  `IDENTITY_PROVIDER_ID` varchar(36) NOT NULL,
  `VALUE` longtext,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`IDENTITY_PROVIDER_ID`,`NAME`),
  CONSTRAINT `FKDC4897CF864C4E43` FOREIGN KEY (`IDENTITY_PROVIDER_ID`) REFERENCES `IDENTITY_PROVIDER` (`INTERNAL_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `IDENTITY_PROVIDER_CONFIG`
--

LOCK TABLES `IDENTITY_PROVIDER_CONFIG` WRITE;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER_CONFIG` DISABLE KEYS */;
INSERT INTO `IDENTITY_PROVIDER_CONFIG` (`IDENTITY_PROVIDER_ID`, `VALUE`, `NAME`) VALUES ('9c52c271-06a9-4133-92a2-779e0c011b91','false','acceptsPromptNoneForwardFromClient'),('9c52c271-06a9-4133-92a2-779e0c011b91','https://example.com/api/v3/','apiUrl'),('9c52c271-06a9-4133-92a2-779e0c011b91','https://example.com/','baseUrl'),('9c52c271-06a9-4133-92a2-779e0c011b91','MASKED_CLIENT_ID','clientId'),('9c52c271-06a9-4133-92a2-779e0c011b91','MASKED_SECRET','clientSecret'),('9c52c271-06a9-4133-92a2-779e0c011b91','openid email','defaultScope'),('9c52c271-06a9-4133-92a2-779e0c011b91','true','disableUserInfo'),('9c52c271-06a9-4133-92a2-779e0c011b91','false','filteredByClaim'),('9c52c271-06a9-4133-92a2-779e0c011b91','false','hideOnLoginPage'),('9c52c271-06a9-4133-92a2-779e0c011b91','IMPORT','syncMode');
/*!40000 ALTER TABLE `IDENTITY_PROVIDER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `IDENTITY_PROVIDER_MAPPER`
--

DROP TABLE IF EXISTS `IDENTITY_PROVIDER_MAPPER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `IDENTITY_PROVIDER_MAPPER` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `IDP_ALIAS` varchar(255) NOT NULL,
  `IDP_MAPPER_NAME` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_ID_PROV_MAPP_REALM` (`REALM_ID`),
  CONSTRAINT `FK_IDPM_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `IDENTITY_PROVIDER_MAPPER`
--

LOCK TABLES `IDENTITY_PROVIDER_MAPPER` WRITE;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER_MAPPER` DISABLE KEYS */;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER_MAPPER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `IDP_MAPPER_CONFIG`
--

DROP TABLE IF EXISTS `IDP_MAPPER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `IDP_MAPPER_CONFIG` (
  `IDP_MAPPER_ID` varchar(36) NOT NULL,
  `VALUE` longtext,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`IDP_MAPPER_ID`,`NAME`),
  CONSTRAINT `FK_IDPMCONFIG` FOREIGN KEY (`IDP_MAPPER_ID`) REFERENCES `IDENTITY_PROVIDER_MAPPER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `IDP_MAPPER_CONFIG`
--

LOCK TABLES `IDP_MAPPER_CONFIG` WRITE;
/*!40000 ALTER TABLE `IDP_MAPPER_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `IDP_MAPPER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `KEYCLOAK_GROUP`
--

DROP TABLE IF EXISTS `KEYCLOAK_GROUP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `KEYCLOAK_GROUP` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `PARENT_GROUP` varchar(36) NOT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `SIBLING_NAMES` (`REALM_ID`,`PARENT_GROUP`,`NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `KEYCLOAK_GROUP`
--

LOCK TABLES `KEYCLOAK_GROUP` WRITE;
/*!40000 ALTER TABLE `KEYCLOAK_GROUP` DISABLE KEYS */;
/*!40000 ALTER TABLE `KEYCLOAK_GROUP` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `KEYCLOAK_ROLE`
--

DROP TABLE IF EXISTS `KEYCLOAK_ROLE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `KEYCLOAK_ROLE` (
  `ID` varchar(36) NOT NULL,
  `CLIENT_REALM_CONSTRAINT` varchar(255) DEFAULT NULL,
  `CLIENT_ROLE` tinyint DEFAULT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `NAME` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `REALM_ID` varchar(255) DEFAULT NULL,
  `CLIENT` varchar(36) DEFAULT NULL,
  `REALM` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_J3RWUVD56ONTGSUHOGM184WW2-2` (`NAME`,`CLIENT_REALM_CONSTRAINT`),
  KEY `IDX_KEYCLOAK_ROLE_CLIENT` (`CLIENT`),
  KEY `IDX_KEYCLOAK_ROLE_REALM` (`REALM`),
  CONSTRAINT `FK_6VYQFE4CN4WLQ8R6KT5VDSJ5C` FOREIGN KEY (`REALM`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `KEYCLOAK_ROLE`
--

LOCK TABLES `KEYCLOAK_ROLE` WRITE;
/*!40000 ALTER TABLE `KEYCLOAK_ROLE` DISABLE KEYS */;
INSERT INTO `KEYCLOAK_ROLE` (`ID`, `CLIENT_REALM_CONSTRAINT`, `CLIENT_ROLE`, `DESCRIPTION`, `NAME`, `REALM_ID`, `CLIENT`, `REALM`) VALUES ('004aa900-cdc2-45a6-9330-06bc9cf45a41','657f6bd0-aa09-4703-99f7-e3f48e2de466',1,'${role_delete-account}','delete-account','4327ba47-4116-44ea-9c4d-02907dca81e7','657f6bd0-aa09-4703-99f7-e3f48e2de466',NULL),('06ad34dc-2164-4067-896b-933e9dadb08c','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_view-events}','view-events','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('07178d4e-bdd9-412f-924b-62977dc1d23b','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_manage-identity-providers}','manage-identity-providers','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('0746f171-6572-4ee4-8cab-2141dcd70f13','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_manage-realm}','manage-realm','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('07755138-7efe-46df-9ae5-4f7869d5c570','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_view-realm}','view-realm','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('084bf744-b03c-4316-be84-5593d9af0691','1947f168-b049-4b78-8031-afcf98eae08d',1,'${role_view-groups}','view-groups','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1947f168-b049-4b78-8031-afcf98eae08d',NULL),('0a6571b7-a3b7-42e9-894d-ab7799bb7e2d','1947f168-b049-4b78-8031-afcf98eae08d',1,'${role_delete-account}','delete-account','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1947f168-b049-4b78-8031-afcf98eae08d',NULL),('0cd165e8-0b76-44c6-8c78-23f2f8b96e18','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_impersonation}','impersonation','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('0daf1344-2778-40fe-a0de-edfed0959aa8','657f6bd0-aa09-4703-99f7-e3f48e2de466',1,'${role_manage-consent}','manage-consent','4327ba47-4116-44ea-9c4d-02907dca81e7','657f6bd0-aa09-4703-99f7-e3f48e2de466',NULL),('0edbe58b-31c3-4503-81e6-3426b88107d7','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_query-realms}','query-realms','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('0f5abc1b-e9de-4f60-94fa-7a9f537eb90a','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_manage-events}','manage-events','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('0fe3ce64-452d-4c89-99f8-af1921093c45','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_query-realms}','query-realms','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('10114928-9ec5-4ce4-baf0-696dd614ea04','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_manage-authorization}','manage-authorization','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('1492ec55-0066-41d1-bfe2-bff7e7ac5259','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_manage-clients}','manage-clients','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','dcc080c5-aede-4fd3-8b01-bd0928b730a2',0,'${role_default-roles}','default-roles-test','dcc080c5-aede-4fd3-8b01-bd0928b730a2',NULL,NULL),('1512286b-232a-48da-a722-79bfa293bd1d','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_view-authorization}','view-authorization','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('173fcaa5-a382-4ae8-819a-317f8d671bec','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_manage-authorization}','manage-authorization','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('175e4ef5-829d-4b42-a98e-dbebadb49f04','41fe1d7e-3717-4965-8c77-e9868b3d98d8',1,'${role_view-profile}','view-profile','dcc080c5-aede-4fd3-8b01-bd0928b730a2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',NULL),('1ad87cfd-6f78-42f8-9b1e-220a2a813fe1','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_view-authorization}','view-authorization','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('1b146798-0c31-4fe4-8d8d-8341e5b463ec','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_view-identity-providers}','view-identity-providers','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('1d41c4d8-3377-488e-8563-f93dfd895d39','657f6bd0-aa09-4703-99f7-e3f48e2de466',1,'${role_view-groups}','view-groups','4327ba47-4116-44ea-9c4d-02907dca81e7','657f6bd0-aa09-4703-99f7-e3f48e2de466',NULL),('22990ac5-82f5-4deb-97d1-342a47945b59','dcc080c5-aede-4fd3-8b01-bd0928b730a2',0,'${role_uma_authorization}','uma_authorization','dcc080c5-aede-4fd3-8b01-bd0928b730a2',NULL,NULL),('22ae141b-4a27-4414-9d66-d95ebc1c264c','4327ba47-4116-44ea-9c4d-02907dca81e7',0,'${role_uma_authorization}','uma_authorization','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL,NULL),('241a9760-03f0-41fa-911c-0a023532161e','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_manage-events}','manage-events','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('245ce92a-a642-4653-bd40-4bb0b179a387','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_query-clients}','query-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('24e2a886-cc15-44c9-a176-1b12c038f88d','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_manage-authorization}','manage-authorization','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('27381875-2dde-4add-93ab-9da6abf300ce','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_realm-admin}','realm-admin','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('280200b5-a94e-4e84-9800-0d152fdd8ba8','9e5fcd3f-0e4b-4627-8292-011015df3de6',1,'${role_read-token}','read-token','4327ba47-4116-44ea-9c4d-02907dca81e7','9e5fcd3f-0e4b-4627-8292-011015df3de6',NULL),('29fbe3d8-6cde-4599-a684-c054a2abb2f4','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_manage-realm}','manage-realm','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('2a6a0506-9553-4b69-822f-6be7668d3b91','1947f168-b049-4b78-8031-afcf98eae08d',1,'${role_view-consent}','view-consent','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1947f168-b049-4b78-8031-afcf98eae08d',NULL),('2bbbedd2-62e8-4965-a5cd-6664c7b100a5','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_view-clients}','view-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('2c0e4faf-3427-4063-ae6d-7c8b65a292ba','41fe1d7e-3717-4965-8c77-e9868b3d98d8',1,'${role_view-consent}','view-consent','dcc080c5-aede-4fd3-8b01-bd0928b730a2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',NULL),('2cf04f97-ea0a-4c30-bcd2-00a25d0700cf','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_view-events}','view-events','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('2d0f3e90-9421-4164-b396-ec0ec75b7099','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_view-realm}','view-realm','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('2d58e9c6-470c-4884-bcf5-c6f92e437a2e','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_create-client}','create-client','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('2e299e39-7ae8-4394-82af-0488b8ba850c','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_view-events}','view-events','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('31102399-1273-491f-8abf-99c630871105','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_impersonation}','impersonation','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('3323ca6f-35d6-4ef6-bba9-d3a65e39d42d','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_manage-identity-providers}','manage-identity-providers','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('346d7236-020b-4f7d-8922-a51bf71a9482','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_manage-clients}','manage-clients','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('39d9f8a7-97c4-459c-9a52-78e597492a95','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_view-events}','view-events','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('39e1e0ae-01d2-40c5-82b7-66727d7542bb','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_create-client}','create-client','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('41a1c8c2-a5f7-433d-a283-f71f08f7779e','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_query-groups}','query-groups','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('43087049-6e35-422e-a082-97d32212cfcf','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_manage-clients}','manage-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('4590b97a-ecbe-42a6-a0fc-85fc8f576d5b','657f6bd0-aa09-4703-99f7-e3f48e2de466',1,'${role_view-profile}','view-profile','4327ba47-4116-44ea-9c4d-02907dca81e7','657f6bd0-aa09-4703-99f7-e3f48e2de466',NULL),('45d90c7d-7aff-453e-b90a-616eef1740b7','1947f168-b049-4b78-8031-afcf98eae08d',1,'${role_view-applications}','view-applications','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1947f168-b049-4b78-8031-afcf98eae08d',NULL),('48eeadab-d1c8-403f-ae89-2270450ae4f6','4327ba47-4116-44ea-9c4d-02907dca81e7',0,'${role_admin}','admin','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL,NULL),('4b5fb1ff-dc8c-4f63-bfaf-0772d809fd00','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_view-clients}','view-clients','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('4d890335-e9e1-45ab-8c8c-98c9e8db8bc4','40ae881c-f4e4-4b07-b097-a67d2bf515e6',0,'${role_offline-access}','offline_access','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL,NULL),('4f90bb4f-92f2-45dc-ad8a-f97723d1d828','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_query-users}','query-users','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('54b52e8c-4a10-4ef6-83ef-daee6dff61f7','a3e6a779-77c2-43c3-bc2a-3eaf4324385f',1,'${role_read-token}','read-token','40ae881c-f4e4-4b07-b097-a67d2bf515e6','a3e6a779-77c2-43c3-bc2a-3eaf4324385f',NULL),('57e92e5d-9b6b-4782-9045-d5ac3b1cd017','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_manage-identity-providers}','manage-identity-providers','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('581d6c1d-f6d2-4d8c-bfcc-56d8e5f8de3d','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_query-realms}','query-realms','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('5d573d8e-9b92-4b95-a6e7-84f8b1eb3cef','1947f168-b049-4b78-8031-afcf98eae08d',1,'${role_view-profile}','view-profile','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1947f168-b049-4b78-8031-afcf98eae08d',NULL),('609ad758-5a8f-42f5-8777-456c380c014f','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_view-realm}','view-realm','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('612b43cd-8762-4e1e-862d-4c632a35b038','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_create-client}','create-client','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('623ac08d-4ff5-49c7-9392-c6f8a49ef26b','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_view-clients}','view-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('65374982-5c23-4b66-a06a-d969c0cd154a','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_query-users}','query-users','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('6552d7cc-3d28-424a-b9ab-c0719ad9d1f1','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_view-clients}','view-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('661dd437-4d3d-4e80-a862-83a2ae44ca5e','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_query-groups}','query-groups','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('669e6b20-0a78-4732-a482-fc68d855d1df','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_manage-users}','manage-users','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('6bbfee03-d9ef-43c5-85b6-59627959df8d','dcc080c5-aede-4fd3-8b01-bd0928b730a2',0,'${role_offline-access}','offline_access','dcc080c5-aede-4fd3-8b01-bd0928b730a2',NULL,NULL),('71e7629d-894b-4aff-825a-765cf58a6f6b','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_impersonation}','impersonation','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('7252897e-111a-4b0e-b87c-c13d52e42696','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_manage-events}','manage-events','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('76a88fee-aa85-44de-a9de-81c6823d603f','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_view-clients}','view-clients','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('76ab719c-a862-441e-a66b-329547590e1d','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_view-identity-providers}','view-identity-providers','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('7717b2f2-fcf1-45d4-86db-0d4c7f9e91a2','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_manage-realm}','manage-realm','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('7a15b2f4-66f6-4ab9-b90e-1d3054723dd1','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_create-client}','create-client','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('7bca4b96-fc2d-48ee-8a5e-bdb2ebe303cc','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_realm-admin}','realm-admin','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('81666137-6815-4158-9161-116220fffe4a','41fe1d7e-3717-4965-8c77-e9868b3d98d8',1,'${role_manage-account-links}','manage-account-links','dcc080c5-aede-4fd3-8b01-bd0928b730a2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',NULL),('81bc3485-f523-4764-89ea-e5bca3b9d2cf','657f6bd0-aa09-4703-99f7-e3f48e2de466',1,'${role_view-applications}','view-applications','4327ba47-4116-44ea-9c4d-02907dca81e7','657f6bd0-aa09-4703-99f7-e3f48e2de466',NULL),('84b5254f-5b61-45e9-9081-7a93b3eba444','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_view-users}','view-users','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('84ec0238-dc71-4066-b99e-49e59cc6bc7a','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_manage-realm}','manage-realm','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('888d64fb-b137-4f1e-b530-b4aef00ab967','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_query-clients}','query-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('88d1e80b-8996-4609-a820-bd3d3efc57ab','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_query-groups}','query-groups','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('8a2220b6-356e-4d7e-995f-83ca8da2e9fe','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_manage-realm}','manage-realm','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('8a32586f-9583-40d7-9cd2-a548f8220c80','40ae881c-f4e4-4b07-b097-a67d2bf515e6',0,'','dev','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL,NULL),('8abdd60d-ccc5-4315-85b6-9b927eaca4a7','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_query-users}','query-users','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('8b146ddf-257f-4390-90ea-2d03c7ab7ce1','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_view-identity-providers}','view-identity-providers','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('8b6e8ddd-250e-4bac-8bdd-de3a9d4ab9e6','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_view-users}','view-users','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('8bcad746-bec2-4e86-ba9a-e85979488a4f','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_impersonation}','impersonation','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('8c549c79-801f-4d45-892d-22fe10111eb2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',1,'${role_manage-account}','manage-account','dcc080c5-aede-4fd3-8b01-bd0928b730a2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',NULL),('8e5f68f7-bcb4-43a5-ae98-6dfc29e18395','657f6bd0-aa09-4703-99f7-e3f48e2de466',1,'${role_manage-account}','manage-account','4327ba47-4116-44ea-9c4d-02907dca81e7','657f6bd0-aa09-4703-99f7-e3f48e2de466',NULL),('91e47ad5-5b4d-4440-b982-3dd69c10b2fb','1947f168-b049-4b78-8031-afcf98eae08d',1,'${role_manage-account}','manage-account','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1947f168-b049-4b78-8031-afcf98eae08d',NULL),('94526ddc-9e84-4c0b-958b-07184a18ab32','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_view-users}','view-users','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('95575514-531e-4f81-91a9-f0c35b743a81','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_view-realm}','view-realm','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('9b9a62d1-f35a-4d79-823e-13957aff5203','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_impersonation}','impersonation','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('a118a201-c733-428b-bf9c-121bcccfb8b5','657f6bd0-aa09-4703-99f7-e3f48e2de466',1,'${role_view-consent}','view-consent','4327ba47-4116-44ea-9c4d-02907dca81e7','657f6bd0-aa09-4703-99f7-e3f48e2de466',NULL),('a2fa72c7-f176-4d55-aab3-9c6e18a6c3e3','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_view-realm}','view-realm','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('a4947d89-da1b-4fb1-ae4b-323921695182','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_manage-identity-providers}','manage-identity-providers','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('a621e3b3-8f5f-460b-ad97-f3ab58c894cb','657f6bd0-aa09-4703-99f7-e3f48e2de466',1,'${role_manage-account-links}','manage-account-links','4327ba47-4116-44ea-9c4d-02907dca81e7','657f6bd0-aa09-4703-99f7-e3f48e2de466',NULL),('a6381b94-b388-4ac0-94a4-a9c24384f9ba','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_view-authorization}','view-authorization','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('a880a908-88f7-45fc-bf1e-4f6fa51fff11','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_manage-clients}','manage-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('ad5932af-bcfc-4464-962e-5490c41eda23','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_manage-events}','manage-events','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('ade6b585-3b0b-4698-b78d-d5b534a6e099','41fe1d7e-3717-4965-8c77-e9868b3d98d8',1,'${role_view-applications}','view-applications','dcc080c5-aede-4fd3-8b01-bd0928b730a2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',NULL),('aea97389-2ee0-4b61-96a7-ff09a73d7996','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_manage-events}','manage-events','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('b197dcfd-a869-436d-8746-8d288ae56ef7','4327ba47-4116-44ea-9c4d-02907dca81e7',0,'${role_offline-access}','offline_access','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL,NULL),('b717dbd5-35d8-4865-94c2-bf876c0146f4','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_query-users}','query-users','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('bad4e73b-8479-4bfc-a203-68300ea9aaf4','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_view-identity-providers}','view-identity-providers','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('bef3ad5c-c6ed-4710-9800-ad029e4632f3','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_view-users}','view-users','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('c1887edf-e0fe-43f9-9857-97e256a6df6e','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_manage-identity-providers}','manage-identity-providers','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('c4fb1f60-f60e-4022-9b59-36c1eef2f82a','1947f168-b049-4b78-8031-afcf98eae08d',1,'${role_manage-consent}','manage-consent','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1947f168-b049-4b78-8031-afcf98eae08d',NULL),('c73afc5a-5da9-416e-a0d9-d50931ad6c0e','4327ba47-4116-44ea-9c4d-02907dca81e7',0,'${role_create-realm}','create-realm','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL,NULL),('c7f14c0f-92f7-4a83-947c-9191b4958168','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_manage-users}','manage-users','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','40ae881c-f4e4-4b07-b097-a67d2bf515e6',0,'${role_default-roles}','default-roles-test-realm','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL,NULL),('c8ee3a69-7a2d-48f8-9373-6677e3724523','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_query-groups}','query-groups','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('c8fc11ef-66b7-438a-a52e-8668bf8ac1b4','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_query-realms}','query-realms','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('cb163866-350b-404d-9bbf-13407643d098','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_manage-authorization}','manage-authorization','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('cde65e2d-9505-48a9-8a5b-4831d65773bf','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_view-identity-providers}','view-identity-providers','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('ce73d6bd-9432-4068-a9e6-02f000f8b4d3','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_view-authorization}','view-authorization','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('d0f5422f-9597-4ea6-abbb-7c346c6a8096','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_view-users}','view-users','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('d305da17-3117-49ce-90eb-91fc9bb929af','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_view-events}','view-events','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('d5f7c5c8-33c4-483c-802e-6ed2c1bc2f29','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_query-groups}','query-groups','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('d62ac468-1487-4e80-afb8-d12035ef8d24','40ae881c-f4e4-4b07-b097-a67d2bf515e6',0,'${role_uma_authorization}','uma_authorization','40ae881c-f4e4-4b07-b097-a67d2bf515e6',NULL,NULL),('d69ea56d-5246-4e89-a45f-4effdf6d9bcf','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_manage-clients}','manage-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('d872e92a-c9f5-482c-9e10-de35d8ed80a3','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_query-clients}','query-clients','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('d892644d-7123-4207-a8c5-74fcc7da7653','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_query-realms}','query-realms','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('db24697d-b20d-4a2a-99b6-f41d26e73b7e','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_query-clients}','query-clients','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('dcc9c139-ec68-483f-9b2e-6db54c81fad4','1947f168-b049-4b78-8031-afcf98eae08d',1,'${role_manage-account-links}','manage-account-links','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1947f168-b049-4b78-8031-afcf98eae08d',NULL),('dee91445-1e60-47a2-9f1a-51a8fdf925ab','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',1,'${role_manage-authorization}','manage-authorization','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5873dcb9-9412-4db8-b4d1-e57a1d76dbd9',NULL),('df5897b4-53cf-44d4-94ac-c8d1651af181','b79d0e1c-6d44-4f01-b111-278fe3db31ee',1,'${role_manage-users}','manage-users','4327ba47-4116-44ea-9c4d-02907dca81e7','b79d0e1c-6d44-4f01-b111-278fe3db31ee',NULL),('e081abd2-551d-469d-bbd7-e522600639f5','41fe1d7e-3717-4965-8c77-e9868b3d98d8',1,'${role_view-groups}','view-groups','dcc080c5-aede-4fd3-8b01-bd0928b730a2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',NULL),('e293b015-eada-4709-8d6c-da9c98a5ec04','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_manage-users}','manage-users','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('e36aa23d-1f2a-417d-9f68-90234cd0703e','41fe1d7e-3717-4965-8c77-e9868b3d98d8',1,'${role_manage-consent}','manage-consent','dcc080c5-aede-4fd3-8b01-bd0928b730a2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',NULL),('e7b5b42b-09a2-46bf-b6b5-87f71892c9dd','4327ba47-4116-44ea-9c4d-02907dca81e7',0,'${role_default-roles}','default-roles-master','4327ba47-4116-44ea-9c4d-02907dca81e7',NULL,NULL),('ea082079-5ed5-4688-8d02-1a0a2dbed4fa','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_query-users}','query-users','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('ec02850f-1ba7-4af0-bed3-09cc17ab1067','41fe1d7e-3717-4965-8c77-e9868b3d98d8',1,'${role_delete-account}','delete-account','dcc080c5-aede-4fd3-8b01-bd0928b730a2','41fe1d7e-3717-4965-8c77-e9868b3d98d8',NULL),('ec9a917b-c79f-422b-941d-40214a7e84cf','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_query-clients}','query-clients','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('ee7eb7c2-62cf-4526-9560-a1e1c4bd27f2','72c6029f-f8d2-4256-a326-2642c15f3a1e',1,'${role_create-client}','create-client','4327ba47-4116-44ea-9c4d-02907dca81e7','72c6029f-f8d2-4256-a326-2642c15f3a1e',NULL),('f6fafe29-29a1-4124-9a53-c98fd932885d','ff653271-7bc7-410a-99e0-cfeff6180f7b',1,'${role_view-authorization}','view-authorization','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ff653271-7bc7-410a-99e0-cfeff6180f7b',NULL),('f9e0e15e-14f9-4327-b781-04bc3078317e','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1,'${role_manage-users}','manage-users','4327ba47-4116-44ea-9c4d-02907dca81e7','1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',NULL),('ff3e2fb7-3987-4e89-a46a-67bb55f9c759','a6b5e3ba-b72b-4028-92be-acb75bd541c1',1,'${role_read-token}','read-token','dcc080c5-aede-4fd3-8b01-bd0928b730a2','a6b5e3ba-b72b-4028-92be-acb75bd541c1',NULL);
/*!40000 ALTER TABLE `KEYCLOAK_ROLE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `MIGRATION_MODEL`
--

DROP TABLE IF EXISTS `MIGRATION_MODEL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `MIGRATION_MODEL` (
  `ID` varchar(36) NOT NULL,
  `VERSION` varchar(36) DEFAULT NULL,
  `UPDATE_TIME` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  KEY `IDX_UPDATE_TIME` (`UPDATE_TIME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `MIGRATION_MODEL`
--

LOCK TABLES `MIGRATION_MODEL` WRITE;
/*!40000 ALTER TABLE `MIGRATION_MODEL` DISABLE KEYS */;
INSERT INTO `MIGRATION_MODEL` (`ID`, `VERSION`, `UPDATE_TIME`) VALUES ('nvun5','24.0.2',1713408442);
/*!40000 ALTER TABLE `MIGRATION_MODEL` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `OFFLINE_CLIENT_SESSION`
--

DROP TABLE IF EXISTS `OFFLINE_CLIENT_SESSION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `OFFLINE_CLIENT_SESSION` (
  `USER_SESSION_ID` varchar(36) NOT NULL,
  `CLIENT_ID` varchar(255) NOT NULL,
  `OFFLINE_FLAG` varchar(4) NOT NULL,
  `TIMESTAMP` int DEFAULT NULL,
  `DATA` longtext,
  `CLIENT_STORAGE_PROVIDER` varchar(36) NOT NULL DEFAULT 'local',
  `EXTERNAL_CLIENT_ID` varchar(255) NOT NULL DEFAULT 'local',
  PRIMARY KEY (`USER_SESSION_ID`,`CLIENT_ID`,`CLIENT_STORAGE_PROVIDER`,`EXTERNAL_CLIENT_ID`,`OFFLINE_FLAG`),
  KEY `IDX_US_SESS_ID_ON_CL_SESS` (`USER_SESSION_ID`),
  KEY `IDX_OFFLINE_CSS_PRELOAD` (`CLIENT_ID`,`OFFLINE_FLAG`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `OFFLINE_CLIENT_SESSION`
--

LOCK TABLES `OFFLINE_CLIENT_SESSION` WRITE;
/*!40000 ALTER TABLE `OFFLINE_CLIENT_SESSION` DISABLE KEYS */;
/*!40000 ALTER TABLE `OFFLINE_CLIENT_SESSION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `OFFLINE_USER_SESSION`
--

DROP TABLE IF EXISTS `OFFLINE_USER_SESSION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `OFFLINE_USER_SESSION` (
  `USER_SESSION_ID` varchar(36) NOT NULL,
  `USER_ID` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `CREATED_ON` int NOT NULL,
  `OFFLINE_FLAG` varchar(4) NOT NULL,
  `DATA` longtext,
  `LAST_SESSION_REFRESH` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`USER_SESSION_ID`,`OFFLINE_FLAG`),
  KEY `IDX_OFFLINE_USS_CREATEDON` (`CREATED_ON`),
  KEY `IDX_OFFLINE_USS_PRELOAD` (`OFFLINE_FLAG`,`CREATED_ON`,`USER_SESSION_ID`),
  KEY `IDX_OFFLINE_USS_BY_USER` (`USER_ID`,`REALM_ID`,`OFFLINE_FLAG`),
  KEY `IDX_OFFLINE_USS_BY_USERSESS` (`REALM_ID`,`OFFLINE_FLAG`,`USER_SESSION_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `OFFLINE_USER_SESSION`
--

LOCK TABLES `OFFLINE_USER_SESSION` WRITE;
/*!40000 ALTER TABLE `OFFLINE_USER_SESSION` DISABLE KEYS */;
/*!40000 ALTER TABLE `OFFLINE_USER_SESSION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `POLICY_CONFIG`
--

DROP TABLE IF EXISTS `POLICY_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `POLICY_CONFIG` (
  `POLICY_ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `VALUE` longtext,
  PRIMARY KEY (`POLICY_ID`,`NAME`),
  CONSTRAINT `FKDC34197CF864C4E43` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `POLICY_CONFIG`
--

LOCK TABLES `POLICY_CONFIG` WRITE;
/*!40000 ALTER TABLE `POLICY_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `POLICY_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `PROTOCOL_MAPPER`
--

DROP TABLE IF EXISTS `PROTOCOL_MAPPER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PROTOCOL_MAPPER` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `PROTOCOL` varchar(255) NOT NULL,
  `PROTOCOL_MAPPER_NAME` varchar(255) NOT NULL,
  `CLIENT_ID` varchar(36) DEFAULT NULL,
  `CLIENT_SCOPE_ID` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_PROTOCOL_MAPPER_CLIENT` (`CLIENT_ID`),
  KEY `IDX_CLSCOPE_PROTMAP` (`CLIENT_SCOPE_ID`),
  CONSTRAINT `FK_CLI_SCOPE_MAPPER` FOREIGN KEY (`CLIENT_SCOPE_ID`) REFERENCES `CLIENT_SCOPE` (`ID`),
  CONSTRAINT `FK_PCM_REALM` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `PROTOCOL_MAPPER`
--

LOCK TABLES `PROTOCOL_MAPPER` WRITE;
/*!40000 ALTER TABLE `PROTOCOL_MAPPER` DISABLE KEYS */;
INSERT INTO `PROTOCOL_MAPPER` (`ID`, `NAME`, `PROTOCOL`, `PROTOCOL_MAPPER_NAME`, `CLIENT_ID`, `CLIENT_SCOPE_ID`) VALUES ('021728c8-703d-4e77-9392-162c0bebad66','phone number','openid-connect','oidc-usermodel-attribute-mapper',NULL,'6eed7521-b56d-4af5-8927-771c01723b21'),('07f3ea11-1024-403f-b318-8ee318f1518a','audience resolve','openid-connect','oidc-audience-resolve-mapper','320d776e-9e8e-4263-abf6-ebe3d862d547',NULL),('0bb13c84-ba08-4f59-ae4c-0cd2707b9bce','given name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('0cab1b6a-1e12-4ce7-82d6-a4fb633d4a0f','realm roles','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'29003608-50c2-4693-953b-316f0c71fa25'),('0e5d7c54-7fcf-43c0-b65f-a4ce338f3052','email verified','openid-connect','oidc-usermodel-property-mapper',NULL,'3d4aaffa-8399-4bff-9dac-000e43b45ca9'),('0fffd0d0-af8f-4cde-adf2-a9d3e5b44a20','realm roles','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'90977ab5-207e-4b07-836d-e8d59805927b'),('138ed530-8f1b-47a7-9e1e-1734cd8d77e8','middle name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('18054524-c85d-448a-ac47-d621f9bc8a02','phone number verified','openid-connect','oidc-usermodel-attribute-mapper',NULL,'fd09aae5-67c6-4162-a03e-95f8cc07d053'),('191b56f5-d2b6-4a6d-8cc5-61fad158ec49','website','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('1bb5e7d4-4207-46f1-adcb-9ba968790a5e','upn','openid-connect','oidc-usermodel-attribute-mapper',NULL,'15afd007-2320-4439-853b-2ddaa8b2ff71'),('1c0b579c-a3c5-439f-a50c-4a4175a77e2b','website','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('25f42fd4-0985-44e0-aa09-313b5dc70b6f','acr loa level','openid-connect','oidc-acr-mapper',NULL,'39474b38-4002-4590-a90e-7e77c25785b5'),('2942d028-1edc-475f-937b-82df960d8ffa','allowed web origins','openid-connect','oidc-allowed-origins-mapper',NULL,'c294c451-1979-493f-ad49-15220c9bd6f8'),('2c294f2f-cc03-462a-a319-f3112592ebeb','picture','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('2cb96ec5-bfe3-47ed-be66-4e5c2ec61330','zoneinfo','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('2ddbc658-9c79-4222-8987-f5b7376cbf59','username','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('30442d51-b460-4110-ba76-6001bf5f9037','realm roles','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'b9388bd3-40d7-4f4a-87e8-2d78275ee434'),('314a4226-8563-4771-87df-9fce92b7dbf0','phone number verified','openid-connect','oidc-usermodel-attribute-mapper',NULL,'d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70'),('3187378e-c8ef-4fd8-b02d-0da7896f5d46','zoneinfo','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('3393a33f-12b4-4bbb-a8ee-56b36c99b81b','role list','saml','saml-role-list-mapper',NULL,'6bffc256-3149-4d67-a5d2-888076a43a46'),('33ef93f2-1653-48a1-a8df-263acc8f986d','username','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('34184086-7ba4-406f-92a0-31c9dd253d51','groups','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'fafceb53-8f99-4a18-81bf-59c2b83bcd6c'),('34284a3d-4d55-4df3-b21a-ad1fa1207c94','acr loa level','openid-connect','oidc-acr-mapper',NULL,'7b474ad1-7a28-4e53-970d-e9a5a14bfbab'),('34ffbda3-f11e-4fd5-8391-453a496d8fe8','audience resolve','openid-connect','oidc-audience-resolve-mapper','b1c791ca-b006-4f60-af37-20225d5876f5',NULL),('353182cb-5c75-4717-bf2a-1231cb8b6ea9','given name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('3831e91f-261e-445d-954e-025a291203d2','audience resolve','openid-connect','oidc-audience-resolve-mapper',NULL,'b9388bd3-40d7-4f4a-87e8-2d78275ee434'),('4a7f56e7-2fcd-4fa6-a71b-6a0247cffd31','email verified','openid-connect','oidc-usermodel-property-mapper',NULL,'815032a2-e511-4d5f-a543-e9a3ee4f648c'),('4d423f68-ac94-4601-ab33-ea1af09c09f8','family name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('4f3b98c9-3ed1-4f8f-abe1-68eaac9845cc','upn','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('50aec82a-e103-4eff-a2ab-00f4316c1fc7','upn','openid-connect','oidc-usermodel-attribute-mapper',NULL,'fafceb53-8f99-4a18-81bf-59c2b83bcd6c'),('516c0ba4-8437-4037-9c1f-0b825063a210','birthdate','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('52aaaf1e-58fa-40b2-a0cb-77f30de33c33','given name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('55f41b5a-b9d6-4275-a588-f81c66d0db7c','website','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('5932fcff-e9ad-430c-bb1e-ef3e37e3b7fe','profile','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('5be3b8a2-1415-4744-88ea-385db06317b5','gender','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('5fe4ee9e-601b-43b7-bcaa-00dd6bd065bd','groups','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'15afd007-2320-4439-853b-2ddaa8b2ff71'),('636043ee-d8d9-402d-bbd0-3e5d4de60d04','birthdate','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('666983e3-bbe6-4e9b-8ab3-d074ef8d1d93','profile','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('6699de21-5d88-4f1a-8870-df11d413e229','full name','openid-connect','oidc-full-name-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('6aeee5d4-b388-49dc-94ea-fa6faf92e2c2','picture','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('6ca63e2f-5170-495a-945e-ac8b46a2defb','picture','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('71996efe-e23c-4298-a1b0-3059b0a7b64b','email verified','openid-connect','oidc-usermodel-property-mapper',NULL,'3e744bf6-ec28-46cf-af95-5cf3f4ee8d58'),('7a791e19-b81b-48c7-b2c8-d721388b5af3','phone number','openid-connect','oidc-usermodel-attribute-mapper',NULL,'d9ca344b-e3ea-41e3-9d4c-2d5a17ddaa70'),('7da7b324-e9da-4cc6-a556-f5eea971c363','client roles','openid-connect','oidc-usermodel-client-role-mapper',NULL,'90977ab5-207e-4b07-836d-e8d59805927b'),('7f25bda8-415c-40a2-91b5-f284b0acb168','nickname','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('8196aa7f-8f19-435b-81ca-03ed26b8469a','client roles','openid-connect','oidc-usermodel-client-role-mapper',NULL,'29003608-50c2-4693-953b-316f0c71fa25'),('829dc1e5-b82f-469b-9b69-4f50f62405ea','updated at','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('84f2be97-83cc-4cbb-96c6-732478dce549','phone number verified','openid-connect','oidc-usermodel-attribute-mapper',NULL,'6eed7521-b56d-4af5-8927-771c01723b21'),('888803ba-8cb8-4d0a-920b-f0f73f8bd71c','birthdate','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('8bafd9bd-84ee-47f1-9040-cdca10d3ad7d','audience resolve','openid-connect','oidc-audience-resolve-mapper',NULL,'90977ab5-207e-4b07-836d-e8d59805927b'),('8d357cda-3d2b-4339-abfc-8954e05b8a65','family name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('9164f05c-fc62-4050-a4dc-d22071898bd5','audience resolve','openid-connect','oidc-audience-resolve-mapper','1a0f720d-df73-4056-926d-10d520c81992',NULL),('971205e3-82a7-4d69-83a2-1ac233a45871','acr loa level','openid-connect','oidc-acr-mapper',NULL,'14d13053-c1b9-43f8-b993-9d7a8aef8069'),('97916051-fbc2-4e33-8ecf-7cb991a21494','profile','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('99dd5b9a-eb8e-4e4a-923b-28384c9b949b','locale','openid-connect','oidc-usermodel-attribute-mapper','05021082-bfbc-4ec6-872a-e6c0916922c1',NULL),('9a8f917f-06e2-42da-b18c-165685ad6765','gender','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('9b962df6-ac03-416e-af7a-cfb3952a40b1','nickname','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('9e0c037e-cd75-4bfc-9df4-437b01de3684','email','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3d4aaffa-8399-4bff-9dac-000e43b45ca9'),('a10531b0-988e-4d54-8843-2c057423fad2','updated at','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('a2cd5654-d33c-477a-b698-27abb94173df','zoneinfo','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('aa7ca020-862d-47bb-91c5-e623954b13d1','address','openid-connect','oidc-address-mapper',NULL,'65b407fe-a29a-4caa-90c1-060588438771'),('ac444e38-613f-4371-bd99-04c6db5599de','locale','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('b0386818-20bb-4fc8-a716-da9a851aac0f','audience resolve','openid-connect','oidc-audience-resolve-mapper',NULL,'29003608-50c2-4693-953b-316f0c71fa25'),('b7478b1b-b796-4202-9306-44c064d8653f','role list','saml','saml-role-list-mapper',NULL,'bdb21a59-753f-498c-8be2-9fae1f633cd3'),('b7570251-8cbc-4405-bff3-c3e04e7490b1','email','openid-connect','oidc-usermodel-attribute-mapper',NULL,'815032a2-e511-4d5f-a543-e9a3ee4f648c'),('c1689a38-4bb9-4b58-b755-de6bb7ab9ea4','locale','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','address','openid-connect','oidc-address-mapper',NULL,'310fbe9a-4f1f-4683-aca4-516d640a5d9b'),('cff060aa-842f-456d-ace5-db03c5215b37','address','openid-connect','oidc-address-mapper',NULL,'066554e5-ed72-4255-a59e-37bce592656c'),('d21c25d8-c749-49f6-a8b6-d77fb0d9719c','middle name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('d8556613-0e26-4ce1-a3a0-21e4d595ab02','username','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('dac1e8be-fcad-4a6c-b751-e96396c173b7','upn','openid-connect','oidc-usermodel-attribute-mapper',NULL,'192eb290-9206-4ec7-94ce-f62e6ab82179'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','test_role','openid-connect','oidc-usermodel-attribute-mapper','ecabf8ab-e548-4b7a-a158-4d3e774afd77',NULL),('e1497be5-5c93-45cc-9e43-d2a26eeb567a','allowed web origins','openid-connect','oidc-allowed-origins-mapper',NULL,'74beed12-a014-4842-91e7-5334c22a5bc0'),('e1e61070-5173-494a-9432-5b38cc69cf96','family name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('e4fe6688-9f00-4cff-9bd9-0610c77ede27','locale','openid-connect','oidc-usermodel-attribute-mapper','33ce64f8-6ffd-430d-a260-c5c8f6d92308',NULL),('e97a6c0d-41f2-4d65-a6a5-0d2802370f47','updated at','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('eb0979ca-36e1-4b3d-93ac-6d84cdfb0a63','locale','openid-connect','oidc-usermodel-attribute-mapper','c12ebc78-3392-401a-8328-5dbb4cddf222',NULL),('ec272038-96c7-4eec-b0a4-139653d08cfe','gender','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('ee94a117-33ff-45f4-a2ec-56c15719d2ab','nickname','openid-connect','oidc-usermodel-attribute-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('ef665f48-8d78-4e5d-b6aa-fde9fbe8dbaf','phone number','openid-connect','oidc-usermodel-attribute-mapper',NULL,'fd09aae5-67c6-4162-a03e-95f8cc07d053'),('f1b0cdf5-e7d2-4e4d-897a-e6e41a2d46ec','middle name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'40f0b7c4-e882-42e9-8cc2-00055b119ca8'),('f5802976-75fa-43c1-a376-89b4fb745969','client roles','openid-connect','oidc-usermodel-client-role-mapper',NULL,'b9388bd3-40d7-4f4a-87e8-2d78275ee434'),('f5f6bffa-a5c4-475a-abed-55038e3c3896','full name','openid-connect','oidc-full-name-mapper',NULL,'ef9560b3-735c-4882-9aa1-dce8fe0697af'),('f6ea5b8b-104d-4e3f-a142-73cb07f9c824','groups','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'192eb290-9206-4ec7-94ce-f62e6ab82179'),('f717b364-1349-451d-b084-5bc291c85208','email','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3e744bf6-ec28-46cf-af95-5cf3f4ee8d58'),('f7889d19-a15c-4f01-9c8c-37d528740726','role list','saml','saml-role-list-mapper',NULL,'091b2057-69d2-4a0e-a76c-7cbd76698850'),('faa279e3-db18-4766-8ebd-15d3700d1759','locale','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c'),('fb6dfb7d-1a8c-4227-b3d9-ba952b697162','allowed web origins','openid-connect','oidc-allowed-origins-mapper',NULL,'aaa43c0a-2c37-42b7-a80c-f98eef322343'),('ff6772e5-ec2e-4af7-9265-16326bee2c6b','full name','openid-connect','oidc-full-name-mapper',NULL,'9c93a802-d9c8-468e-8983-129fc287b26c');
/*!40000 ALTER TABLE `PROTOCOL_MAPPER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `PROTOCOL_MAPPER_CONFIG`
--

DROP TABLE IF EXISTS `PROTOCOL_MAPPER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PROTOCOL_MAPPER_CONFIG` (
  `PROTOCOL_MAPPER_ID` varchar(36) NOT NULL,
  `VALUE` longtext,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`PROTOCOL_MAPPER_ID`,`NAME`),
  CONSTRAINT `FK_PMCONFIG` FOREIGN KEY (`PROTOCOL_MAPPER_ID`) REFERENCES `PROTOCOL_MAPPER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `PROTOCOL_MAPPER_CONFIG`
--

LOCK TABLES `PROTOCOL_MAPPER_CONFIG` WRITE;
/*!40000 ALTER TABLE `PROTOCOL_MAPPER_CONFIG` DISABLE KEYS */;
INSERT INTO `PROTOCOL_MAPPER_CONFIG` (`PROTOCOL_MAPPER_ID`, `VALUE`, `NAME`) VALUES ('021728c8-703d-4e77-9392-162c0bebad66','true','access.token.claim'),('021728c8-703d-4e77-9392-162c0bebad66','phone_number','claim.name'),('021728c8-703d-4e77-9392-162c0bebad66','true','id.token.claim'),('021728c8-703d-4e77-9392-162c0bebad66','true','introspection.token.claim'),('021728c8-703d-4e77-9392-162c0bebad66','String','jsonType.label'),('021728c8-703d-4e77-9392-162c0bebad66','phoneNumber','user.attribute'),('021728c8-703d-4e77-9392-162c0bebad66','true','userinfo.token.claim'),('0bb13c84-ba08-4f59-ae4c-0cd2707b9bce','true','access.token.claim'),('0bb13c84-ba08-4f59-ae4c-0cd2707b9bce','given_name','claim.name'),('0bb13c84-ba08-4f59-ae4c-0cd2707b9bce','true','id.token.claim'),('0bb13c84-ba08-4f59-ae4c-0cd2707b9bce','true','introspection.token.claim'),('0bb13c84-ba08-4f59-ae4c-0cd2707b9bce','String','jsonType.label'),('0bb13c84-ba08-4f59-ae4c-0cd2707b9bce','firstName','user.attribute'),('0bb13c84-ba08-4f59-ae4c-0cd2707b9bce','true','userinfo.token.claim'),('0cab1b6a-1e12-4ce7-82d6-a4fb633d4a0f','true','access.token.claim'),('0cab1b6a-1e12-4ce7-82d6-a4fb633d4a0f','realm_access.roles','claim.name'),('0cab1b6a-1e12-4ce7-82d6-a4fb633d4a0f','true','introspection.token.claim'),('0cab1b6a-1e12-4ce7-82d6-a4fb633d4a0f','String','jsonType.label'),('0cab1b6a-1e12-4ce7-82d6-a4fb633d4a0f','true','multivalued'),('0cab1b6a-1e12-4ce7-82d6-a4fb633d4a0f','foo','user.attribute'),('0e5d7c54-7fcf-43c0-b65f-a4ce338f3052','true','access.token.claim'),('0e5d7c54-7fcf-43c0-b65f-a4ce338f3052','email_verified','claim.name'),('0e5d7c54-7fcf-43c0-b65f-a4ce338f3052','true','id.token.claim'),('0e5d7c54-7fcf-43c0-b65f-a4ce338f3052','true','introspection.token.claim'),('0e5d7c54-7fcf-43c0-b65f-a4ce338f3052','boolean','jsonType.label'),('0e5d7c54-7fcf-43c0-b65f-a4ce338f3052','emailVerified','user.attribute'),('0e5d7c54-7fcf-43c0-b65f-a4ce338f3052','true','userinfo.token.claim'),('0fffd0d0-af8f-4cde-adf2-a9d3e5b44a20','true','access.token.claim'),('0fffd0d0-af8f-4cde-adf2-a9d3e5b44a20','realm_access.roles','claim.name'),('0fffd0d0-af8f-4cde-adf2-a9d3e5b44a20','true','introspection.token.claim'),('0fffd0d0-af8f-4cde-adf2-a9d3e5b44a20','String','jsonType.label'),('0fffd0d0-af8f-4cde-adf2-a9d3e5b44a20','true','multivalued'),('0fffd0d0-af8f-4cde-adf2-a9d3e5b44a20','foo','user.attribute'),('138ed530-8f1b-47a7-9e1e-1734cd8d77e8','true','access.token.claim'),('138ed530-8f1b-47a7-9e1e-1734cd8d77e8','middle_name','claim.name'),('138ed530-8f1b-47a7-9e1e-1734cd8d77e8','true','id.token.claim'),('138ed530-8f1b-47a7-9e1e-1734cd8d77e8','true','introspection.token.claim'),('138ed530-8f1b-47a7-9e1e-1734cd8d77e8','String','jsonType.label'),('138ed530-8f1b-47a7-9e1e-1734cd8d77e8','middleName','user.attribute'),('138ed530-8f1b-47a7-9e1e-1734cd8d77e8','true','userinfo.token.claim'),('18054524-c85d-448a-ac47-d621f9bc8a02','true','access.token.claim'),('18054524-c85d-448a-ac47-d621f9bc8a02','phone_number_verified','claim.name'),('18054524-c85d-448a-ac47-d621f9bc8a02','true','id.token.claim'),('18054524-c85d-448a-ac47-d621f9bc8a02','true','introspection.token.claim'),('18054524-c85d-448a-ac47-d621f9bc8a02','boolean','jsonType.label'),('18054524-c85d-448a-ac47-d621f9bc8a02','phoneNumberVerified','user.attribute'),('18054524-c85d-448a-ac47-d621f9bc8a02','true','userinfo.token.claim'),('191b56f5-d2b6-4a6d-8cc5-61fad158ec49','true','access.token.claim'),('191b56f5-d2b6-4a6d-8cc5-61fad158ec49','website','claim.name'),('191b56f5-d2b6-4a6d-8cc5-61fad158ec49','true','id.token.claim'),('191b56f5-d2b6-4a6d-8cc5-61fad158ec49','true','introspection.token.claim'),('191b56f5-d2b6-4a6d-8cc5-61fad158ec49','String','jsonType.label'),('191b56f5-d2b6-4a6d-8cc5-61fad158ec49','website','user.attribute'),('191b56f5-d2b6-4a6d-8cc5-61fad158ec49','true','userinfo.token.claim'),('1bb5e7d4-4207-46f1-adcb-9ba968790a5e','true','access.token.claim'),('1bb5e7d4-4207-46f1-adcb-9ba968790a5e','upn','claim.name'),('1bb5e7d4-4207-46f1-adcb-9ba968790a5e','true','id.token.claim'),('1bb5e7d4-4207-46f1-adcb-9ba968790a5e','true','introspection.token.claim'),('1bb5e7d4-4207-46f1-adcb-9ba968790a5e','String','jsonType.label'),('1bb5e7d4-4207-46f1-adcb-9ba968790a5e','username','user.attribute'),('1bb5e7d4-4207-46f1-adcb-9ba968790a5e','true','userinfo.token.claim'),('1c0b579c-a3c5-439f-a50c-4a4175a77e2b','true','access.token.claim'),('1c0b579c-a3c5-439f-a50c-4a4175a77e2b','website','claim.name'),('1c0b579c-a3c5-439f-a50c-4a4175a77e2b','true','id.token.claim'),('1c0b579c-a3c5-439f-a50c-4a4175a77e2b','true','introspection.token.claim'),('1c0b579c-a3c5-439f-a50c-4a4175a77e2b','String','jsonType.label'),('1c0b579c-a3c5-439f-a50c-4a4175a77e2b','website','user.attribute'),('1c0b579c-a3c5-439f-a50c-4a4175a77e2b','true','userinfo.token.claim'),('25f42fd4-0985-44e0-aa09-313b5dc70b6f','true','access.token.claim'),('25f42fd4-0985-44e0-aa09-313b5dc70b6f','true','id.token.claim'),('25f42fd4-0985-44e0-aa09-313b5dc70b6f','true','introspection.token.claim'),('2942d028-1edc-475f-937b-82df960d8ffa','true','access.token.claim'),('2942d028-1edc-475f-937b-82df960d8ffa','true','introspection.token.claim'),('2c294f2f-cc03-462a-a319-f3112592ebeb','true','access.token.claim'),('2c294f2f-cc03-462a-a319-f3112592ebeb','picture','claim.name'),('2c294f2f-cc03-462a-a319-f3112592ebeb','true','id.token.claim'),('2c294f2f-cc03-462a-a319-f3112592ebeb','true','introspection.token.claim'),('2c294f2f-cc03-462a-a319-f3112592ebeb','String','jsonType.label'),('2c294f2f-cc03-462a-a319-f3112592ebeb','picture','user.attribute'),('2c294f2f-cc03-462a-a319-f3112592ebeb','true','userinfo.token.claim'),('2cb96ec5-bfe3-47ed-be66-4e5c2ec61330','true','access.token.claim'),('2cb96ec5-bfe3-47ed-be66-4e5c2ec61330','zoneinfo','claim.name'),('2cb96ec5-bfe3-47ed-be66-4e5c2ec61330','true','id.token.claim'),('2cb96ec5-bfe3-47ed-be66-4e5c2ec61330','true','introspection.token.claim'),('2cb96ec5-bfe3-47ed-be66-4e5c2ec61330','String','jsonType.label'),('2cb96ec5-bfe3-47ed-be66-4e5c2ec61330','zoneinfo','user.attribute'),('2cb96ec5-bfe3-47ed-be66-4e5c2ec61330','true','userinfo.token.claim'),('2ddbc658-9c79-4222-8987-f5b7376cbf59','true','access.token.claim'),('2ddbc658-9c79-4222-8987-f5b7376cbf59','preferred_username','claim.name'),('2ddbc658-9c79-4222-8987-f5b7376cbf59','true','id.token.claim'),('2ddbc658-9c79-4222-8987-f5b7376cbf59','true','introspection.token.claim'),('2ddbc658-9c79-4222-8987-f5b7376cbf59','String','jsonType.label'),('2ddbc658-9c79-4222-8987-f5b7376cbf59','username','user.attribute'),('2ddbc658-9c79-4222-8987-f5b7376cbf59','true','userinfo.token.claim'),('30442d51-b460-4110-ba76-6001bf5f9037','true','access.token.claim'),('30442d51-b460-4110-ba76-6001bf5f9037','realm_access.roles','claim.name'),('30442d51-b460-4110-ba76-6001bf5f9037','true','introspection.token.claim'),('30442d51-b460-4110-ba76-6001bf5f9037','String','jsonType.label'),('30442d51-b460-4110-ba76-6001bf5f9037','true','multivalued'),('30442d51-b460-4110-ba76-6001bf5f9037','foo','user.attribute'),('314a4226-8563-4771-87df-9fce92b7dbf0','true','access.token.claim'),('314a4226-8563-4771-87df-9fce92b7dbf0','phone_number_verified','claim.name'),('314a4226-8563-4771-87df-9fce92b7dbf0','true','id.token.claim'),('314a4226-8563-4771-87df-9fce92b7dbf0','true','introspection.token.claim'),('314a4226-8563-4771-87df-9fce92b7dbf0','boolean','jsonType.label'),('314a4226-8563-4771-87df-9fce92b7dbf0','phoneNumberVerified','user.attribute'),('314a4226-8563-4771-87df-9fce92b7dbf0','true','userinfo.token.claim'),('3187378e-c8ef-4fd8-b02d-0da7896f5d46','true','access.token.claim'),('3187378e-c8ef-4fd8-b02d-0da7896f5d46','zoneinfo','claim.name'),('3187378e-c8ef-4fd8-b02d-0da7896f5d46','true','id.token.claim'),('3187378e-c8ef-4fd8-b02d-0da7896f5d46','true','introspection.token.claim'),('3187378e-c8ef-4fd8-b02d-0da7896f5d46','String','jsonType.label'),('3187378e-c8ef-4fd8-b02d-0da7896f5d46','zoneinfo','user.attribute'),('3187378e-c8ef-4fd8-b02d-0da7896f5d46','true','userinfo.token.claim'),('3393a33f-12b4-4bbb-a8ee-56b36c99b81b','Role','attribute.name'),('3393a33f-12b4-4bbb-a8ee-56b36c99b81b','Basic','attribute.nameformat'),('3393a33f-12b4-4bbb-a8ee-56b36c99b81b','false','single'),('33ef93f2-1653-48a1-a8df-263acc8f986d','true','access.token.claim'),('33ef93f2-1653-48a1-a8df-263acc8f986d','preferred_username','claim.name'),('33ef93f2-1653-48a1-a8df-263acc8f986d','true','id.token.claim'),('33ef93f2-1653-48a1-a8df-263acc8f986d','true','introspection.token.claim'),('33ef93f2-1653-48a1-a8df-263acc8f986d','String','jsonType.label'),('33ef93f2-1653-48a1-a8df-263acc8f986d','username','user.attribute'),('33ef93f2-1653-48a1-a8df-263acc8f986d','true','userinfo.token.claim'),('34184086-7ba4-406f-92a0-31c9dd253d51','true','access.token.claim'),('34184086-7ba4-406f-92a0-31c9dd253d51','groups','claim.name'),('34184086-7ba4-406f-92a0-31c9dd253d51','true','id.token.claim'),('34184086-7ba4-406f-92a0-31c9dd253d51','true','introspection.token.claim'),('34184086-7ba4-406f-92a0-31c9dd253d51','String','jsonType.label'),('34184086-7ba4-406f-92a0-31c9dd253d51','true','multivalued'),('34184086-7ba4-406f-92a0-31c9dd253d51','foo','user.attribute'),('34284a3d-4d55-4df3-b21a-ad1fa1207c94','true','access.token.claim'),('34284a3d-4d55-4df3-b21a-ad1fa1207c94','true','id.token.claim'),('34284a3d-4d55-4df3-b21a-ad1fa1207c94','true','introspection.token.claim'),('353182cb-5c75-4717-bf2a-1231cb8b6ea9','true','access.token.claim'),('353182cb-5c75-4717-bf2a-1231cb8b6ea9','given_name','claim.name'),('353182cb-5c75-4717-bf2a-1231cb8b6ea9','true','id.token.claim'),('353182cb-5c75-4717-bf2a-1231cb8b6ea9','true','introspection.token.claim'),('353182cb-5c75-4717-bf2a-1231cb8b6ea9','String','jsonType.label'),('353182cb-5c75-4717-bf2a-1231cb8b6ea9','firstName','user.attribute'),('353182cb-5c75-4717-bf2a-1231cb8b6ea9','true','userinfo.token.claim'),('3831e91f-261e-445d-954e-025a291203d2','true','access.token.claim'),('3831e91f-261e-445d-954e-025a291203d2','true','introspection.token.claim'),('4a7f56e7-2fcd-4fa6-a71b-6a0247cffd31','true','access.token.claim'),('4a7f56e7-2fcd-4fa6-a71b-6a0247cffd31','email_verified','claim.name'),('4a7f56e7-2fcd-4fa6-a71b-6a0247cffd31','true','id.token.claim'),('4a7f56e7-2fcd-4fa6-a71b-6a0247cffd31','true','introspection.token.claim'),('4a7f56e7-2fcd-4fa6-a71b-6a0247cffd31','boolean','jsonType.label'),('4a7f56e7-2fcd-4fa6-a71b-6a0247cffd31','emailVerified','user.attribute'),('4a7f56e7-2fcd-4fa6-a71b-6a0247cffd31','true','userinfo.token.claim'),('4d423f68-ac94-4601-ab33-ea1af09c09f8','true','access.token.claim'),('4d423f68-ac94-4601-ab33-ea1af09c09f8','family_name','claim.name'),('4d423f68-ac94-4601-ab33-ea1af09c09f8','true','id.token.claim'),('4d423f68-ac94-4601-ab33-ea1af09c09f8','true','introspection.token.claim'),('4d423f68-ac94-4601-ab33-ea1af09c09f8','String','jsonType.label'),('4d423f68-ac94-4601-ab33-ea1af09c09f8','lastName','user.attribute'),('4d423f68-ac94-4601-ab33-ea1af09c09f8','true','userinfo.token.claim'),('4f3b98c9-3ed1-4f8f-abe1-68eaac9845cc','true','access.token.claim'),('4f3b98c9-3ed1-4f8f-abe1-68eaac9845cc','upn','claim.name'),('4f3b98c9-3ed1-4f8f-abe1-68eaac9845cc','true','id.token.claim'),('4f3b98c9-3ed1-4f8f-abe1-68eaac9845cc','true','introspection.token.claim'),('4f3b98c9-3ed1-4f8f-abe1-68eaac9845cc','String','jsonType.label'),('4f3b98c9-3ed1-4f8f-abe1-68eaac9845cc','username','user.attribute'),('4f3b98c9-3ed1-4f8f-abe1-68eaac9845cc','true','userinfo.token.claim'),('50aec82a-e103-4eff-a2ab-00f4316c1fc7','true','access.token.claim'),('50aec82a-e103-4eff-a2ab-00f4316c1fc7','upn','claim.name'),('50aec82a-e103-4eff-a2ab-00f4316c1fc7','true','id.token.claim'),('50aec82a-e103-4eff-a2ab-00f4316c1fc7','true','introspection.token.claim'),('50aec82a-e103-4eff-a2ab-00f4316c1fc7','String','jsonType.label'),('50aec82a-e103-4eff-a2ab-00f4316c1fc7','username','user.attribute'),('50aec82a-e103-4eff-a2ab-00f4316c1fc7','true','userinfo.token.claim'),('516c0ba4-8437-4037-9c1f-0b825063a210','true','access.token.claim'),('516c0ba4-8437-4037-9c1f-0b825063a210','birthdate','claim.name'),('516c0ba4-8437-4037-9c1f-0b825063a210','true','id.token.claim'),('516c0ba4-8437-4037-9c1f-0b825063a210','true','introspection.token.claim'),('516c0ba4-8437-4037-9c1f-0b825063a210','String','jsonType.label'),('516c0ba4-8437-4037-9c1f-0b825063a210','birthdate','user.attribute'),('516c0ba4-8437-4037-9c1f-0b825063a210','true','userinfo.token.claim'),('52aaaf1e-58fa-40b2-a0cb-77f30de33c33','true','access.token.claim'),('52aaaf1e-58fa-40b2-a0cb-77f30de33c33','given_name','claim.name'),('52aaaf1e-58fa-40b2-a0cb-77f30de33c33','true','id.token.claim'),('52aaaf1e-58fa-40b2-a0cb-77f30de33c33','true','introspection.token.claim'),('52aaaf1e-58fa-40b2-a0cb-77f30de33c33','String','jsonType.label'),('52aaaf1e-58fa-40b2-a0cb-77f30de33c33','firstName','user.attribute'),('52aaaf1e-58fa-40b2-a0cb-77f30de33c33','true','userinfo.token.claim'),('55f41b5a-b9d6-4275-a588-f81c66d0db7c','true','access.token.claim'),('55f41b5a-b9d6-4275-a588-f81c66d0db7c','website','claim.name'),('55f41b5a-b9d6-4275-a588-f81c66d0db7c','true','id.token.claim'),('55f41b5a-b9d6-4275-a588-f81c66d0db7c','true','introspection.token.claim'),('55f41b5a-b9d6-4275-a588-f81c66d0db7c','String','jsonType.label'),('55f41b5a-b9d6-4275-a588-f81c66d0db7c','website','user.attribute'),('55f41b5a-b9d6-4275-a588-f81c66d0db7c','true','userinfo.token.claim'),('5932fcff-e9ad-430c-bb1e-ef3e37e3b7fe','true','access.token.claim'),('5932fcff-e9ad-430c-bb1e-ef3e37e3b7fe','profile','claim.name'),('5932fcff-e9ad-430c-bb1e-ef3e37e3b7fe','true','id.token.claim'),('5932fcff-e9ad-430c-bb1e-ef3e37e3b7fe','true','introspection.token.claim'),('5932fcff-e9ad-430c-bb1e-ef3e37e3b7fe','String','jsonType.label'),('5932fcff-e9ad-430c-bb1e-ef3e37e3b7fe','profile','user.attribute'),('5932fcff-e9ad-430c-bb1e-ef3e37e3b7fe','true','userinfo.token.claim'),('5be3b8a2-1415-4744-88ea-385db06317b5','true','access.token.claim'),('5be3b8a2-1415-4744-88ea-385db06317b5','gender','claim.name'),('5be3b8a2-1415-4744-88ea-385db06317b5','true','id.token.claim'),('5be3b8a2-1415-4744-88ea-385db06317b5','true','introspection.token.claim'),('5be3b8a2-1415-4744-88ea-385db06317b5','String','jsonType.label'),('5be3b8a2-1415-4744-88ea-385db06317b5','gender','user.attribute'),('5be3b8a2-1415-4744-88ea-385db06317b5','true','userinfo.token.claim'),('5fe4ee9e-601b-43b7-bcaa-00dd6bd065bd','true','access.token.claim'),('5fe4ee9e-601b-43b7-bcaa-00dd6bd065bd','groups','claim.name'),('5fe4ee9e-601b-43b7-bcaa-00dd6bd065bd','true','id.token.claim'),('5fe4ee9e-601b-43b7-bcaa-00dd6bd065bd','true','introspection.token.claim'),('5fe4ee9e-601b-43b7-bcaa-00dd6bd065bd','String','jsonType.label'),('5fe4ee9e-601b-43b7-bcaa-00dd6bd065bd','true','multivalued'),('5fe4ee9e-601b-43b7-bcaa-00dd6bd065bd','foo','user.attribute'),('636043ee-d8d9-402d-bbd0-3e5d4de60d04','true','access.token.claim'),('636043ee-d8d9-402d-bbd0-3e5d4de60d04','birthdate','claim.name'),('636043ee-d8d9-402d-bbd0-3e5d4de60d04','true','id.token.claim'),('636043ee-d8d9-402d-bbd0-3e5d4de60d04','true','introspection.token.claim'),('636043ee-d8d9-402d-bbd0-3e5d4de60d04','String','jsonType.label'),('636043ee-d8d9-402d-bbd0-3e5d4de60d04','birthdate','user.attribute'),('636043ee-d8d9-402d-bbd0-3e5d4de60d04','true','userinfo.token.claim'),('666983e3-bbe6-4e9b-8ab3-d074ef8d1d93','true','access.token.claim'),('666983e3-bbe6-4e9b-8ab3-d074ef8d1d93','profile','claim.name'),('666983e3-bbe6-4e9b-8ab3-d074ef8d1d93','true','id.token.claim'),('666983e3-bbe6-4e9b-8ab3-d074ef8d1d93','true','introspection.token.claim'),('666983e3-bbe6-4e9b-8ab3-d074ef8d1d93','String','jsonType.label'),('666983e3-bbe6-4e9b-8ab3-d074ef8d1d93','profile','user.attribute'),('666983e3-bbe6-4e9b-8ab3-d074ef8d1d93','true','userinfo.token.claim'),('6699de21-5d88-4f1a-8870-df11d413e229','true','access.token.claim'),('6699de21-5d88-4f1a-8870-df11d413e229','true','id.token.claim'),('6699de21-5d88-4f1a-8870-df11d413e229','true','introspection.token.claim'),('6699de21-5d88-4f1a-8870-df11d413e229','true','userinfo.token.claim'),('6aeee5d4-b388-49dc-94ea-fa6faf92e2c2','true','access.token.claim'),('6aeee5d4-b388-49dc-94ea-fa6faf92e2c2','picture','claim.name'),('6aeee5d4-b388-49dc-94ea-fa6faf92e2c2','true','id.token.claim'),('6aeee5d4-b388-49dc-94ea-fa6faf92e2c2','true','introspection.token.claim'),('6aeee5d4-b388-49dc-94ea-fa6faf92e2c2','String','jsonType.label'),('6aeee5d4-b388-49dc-94ea-fa6faf92e2c2','picture','user.attribute'),('6aeee5d4-b388-49dc-94ea-fa6faf92e2c2','true','userinfo.token.claim'),('6ca63e2f-5170-495a-945e-ac8b46a2defb','true','access.token.claim'),('6ca63e2f-5170-495a-945e-ac8b46a2defb','picture','claim.name'),('6ca63e2f-5170-495a-945e-ac8b46a2defb','true','id.token.claim'),('6ca63e2f-5170-495a-945e-ac8b46a2defb','true','introspection.token.claim'),('6ca63e2f-5170-495a-945e-ac8b46a2defb','String','jsonType.label'),('6ca63e2f-5170-495a-945e-ac8b46a2defb','picture','user.attribute'),('6ca63e2f-5170-495a-945e-ac8b46a2defb','true','userinfo.token.claim'),('71996efe-e23c-4298-a1b0-3059b0a7b64b','true','access.token.claim'),('71996efe-e23c-4298-a1b0-3059b0a7b64b','email_verified','claim.name'),('71996efe-e23c-4298-a1b0-3059b0a7b64b','true','id.token.claim'),('71996efe-e23c-4298-a1b0-3059b0a7b64b','true','introspection.token.claim'),('71996efe-e23c-4298-a1b0-3059b0a7b64b','boolean','jsonType.label'),('71996efe-e23c-4298-a1b0-3059b0a7b64b','emailVerified','user.attribute'),('71996efe-e23c-4298-a1b0-3059b0a7b64b','true','userinfo.token.claim'),('7a791e19-b81b-48c7-b2c8-d721388b5af3','true','access.token.claim'),('7a791e19-b81b-48c7-b2c8-d721388b5af3','phone_number','claim.name'),('7a791e19-b81b-48c7-b2c8-d721388b5af3','true','id.token.claim'),('7a791e19-b81b-48c7-b2c8-d721388b5af3','true','introspection.token.claim'),('7a791e19-b81b-48c7-b2c8-d721388b5af3','String','jsonType.label'),('7a791e19-b81b-48c7-b2c8-d721388b5af3','phoneNumber','user.attribute'),('7a791e19-b81b-48c7-b2c8-d721388b5af3','true','userinfo.token.claim'),('7da7b324-e9da-4cc6-a556-f5eea971c363','true','access.token.claim'),('7da7b324-e9da-4cc6-a556-f5eea971c363','resource_access.${client_id}.roles','claim.name'),('7da7b324-e9da-4cc6-a556-f5eea971c363','true','introspection.token.claim'),('7da7b324-e9da-4cc6-a556-f5eea971c363','String','jsonType.label'),('7da7b324-e9da-4cc6-a556-f5eea971c363','true','multivalued'),('7da7b324-e9da-4cc6-a556-f5eea971c363','foo','user.attribute'),('7f25bda8-415c-40a2-91b5-f284b0acb168','true','access.token.claim'),('7f25bda8-415c-40a2-91b5-f284b0acb168','nickname','claim.name'),('7f25bda8-415c-40a2-91b5-f284b0acb168','true','id.token.claim'),('7f25bda8-415c-40a2-91b5-f284b0acb168','true','introspection.token.claim'),('7f25bda8-415c-40a2-91b5-f284b0acb168','String','jsonType.label'),('7f25bda8-415c-40a2-91b5-f284b0acb168','nickname','user.attribute'),('7f25bda8-415c-40a2-91b5-f284b0acb168','true','userinfo.token.claim'),('8196aa7f-8f19-435b-81ca-03ed26b8469a','true','access.token.claim'),('8196aa7f-8f19-435b-81ca-03ed26b8469a','resource_access.${client_id}.roles','claim.name'),('8196aa7f-8f19-435b-81ca-03ed26b8469a','true','introspection.token.claim'),('8196aa7f-8f19-435b-81ca-03ed26b8469a','String','jsonType.label'),('8196aa7f-8f19-435b-81ca-03ed26b8469a','true','multivalued'),('8196aa7f-8f19-435b-81ca-03ed26b8469a','foo','user.attribute'),('829dc1e5-b82f-469b-9b69-4f50f62405ea','true','access.token.claim'),('829dc1e5-b82f-469b-9b69-4f50f62405ea','updated_at','claim.name'),('829dc1e5-b82f-469b-9b69-4f50f62405ea','true','id.token.claim'),('829dc1e5-b82f-469b-9b69-4f50f62405ea','true','introspection.token.claim'),('829dc1e5-b82f-469b-9b69-4f50f62405ea','long','jsonType.label'),('829dc1e5-b82f-469b-9b69-4f50f62405ea','updatedAt','user.attribute'),('829dc1e5-b82f-469b-9b69-4f50f62405ea','true','userinfo.token.claim'),('84f2be97-83cc-4cbb-96c6-732478dce549','true','access.token.claim'),('84f2be97-83cc-4cbb-96c6-732478dce549','phone_number_verified','claim.name'),('84f2be97-83cc-4cbb-96c6-732478dce549','true','id.token.claim'),('84f2be97-83cc-4cbb-96c6-732478dce549','true','introspection.token.claim'),('84f2be97-83cc-4cbb-96c6-732478dce549','boolean','jsonType.label'),('84f2be97-83cc-4cbb-96c6-732478dce549','phoneNumberVerified','user.attribute'),('84f2be97-83cc-4cbb-96c6-732478dce549','true','userinfo.token.claim'),('888803ba-8cb8-4d0a-920b-f0f73f8bd71c','true','access.token.claim'),('888803ba-8cb8-4d0a-920b-f0f73f8bd71c','birthdate','claim.name'),('888803ba-8cb8-4d0a-920b-f0f73f8bd71c','true','id.token.claim'),('888803ba-8cb8-4d0a-920b-f0f73f8bd71c','true','introspection.token.claim'),('888803ba-8cb8-4d0a-920b-f0f73f8bd71c','String','jsonType.label'),('888803ba-8cb8-4d0a-920b-f0f73f8bd71c','birthdate','user.attribute'),('888803ba-8cb8-4d0a-920b-f0f73f8bd71c','true','userinfo.token.claim'),('8bafd9bd-84ee-47f1-9040-cdca10d3ad7d','true','access.token.claim'),('8bafd9bd-84ee-47f1-9040-cdca10d3ad7d','true','introspection.token.claim'),('8d357cda-3d2b-4339-abfc-8954e05b8a65','true','access.token.claim'),('8d357cda-3d2b-4339-abfc-8954e05b8a65','family_name','claim.name'),('8d357cda-3d2b-4339-abfc-8954e05b8a65','true','id.token.claim'),('8d357cda-3d2b-4339-abfc-8954e05b8a65','true','introspection.token.claim'),('8d357cda-3d2b-4339-abfc-8954e05b8a65','String','jsonType.label'),('8d357cda-3d2b-4339-abfc-8954e05b8a65','lastName','user.attribute'),('8d357cda-3d2b-4339-abfc-8954e05b8a65','true','userinfo.token.claim'),('971205e3-82a7-4d69-83a2-1ac233a45871','true','access.token.claim'),('971205e3-82a7-4d69-83a2-1ac233a45871','true','id.token.claim'),('971205e3-82a7-4d69-83a2-1ac233a45871','true','introspection.token.claim'),('97916051-fbc2-4e33-8ecf-7cb991a21494','true','access.token.claim'),('97916051-fbc2-4e33-8ecf-7cb991a21494','profile','claim.name'),('97916051-fbc2-4e33-8ecf-7cb991a21494','true','id.token.claim'),('97916051-fbc2-4e33-8ecf-7cb991a21494','true','introspection.token.claim'),('97916051-fbc2-4e33-8ecf-7cb991a21494','String','jsonType.label'),('97916051-fbc2-4e33-8ecf-7cb991a21494','profile','user.attribute'),('97916051-fbc2-4e33-8ecf-7cb991a21494','true','userinfo.token.claim'),('99dd5b9a-eb8e-4e4a-923b-28384c9b949b','true','access.token.claim'),('99dd5b9a-eb8e-4e4a-923b-28384c9b949b','locale','claim.name'),('99dd5b9a-eb8e-4e4a-923b-28384c9b949b','true','id.token.claim'),('99dd5b9a-eb8e-4e4a-923b-28384c9b949b','true','introspection.token.claim'),('99dd5b9a-eb8e-4e4a-923b-28384c9b949b','String','jsonType.label'),('99dd5b9a-eb8e-4e4a-923b-28384c9b949b','locale','user.attribute'),('99dd5b9a-eb8e-4e4a-923b-28384c9b949b','true','userinfo.token.claim'),('9a8f917f-06e2-42da-b18c-165685ad6765','true','access.token.claim'),('9a8f917f-06e2-42da-b18c-165685ad6765','gender','claim.name'),('9a8f917f-06e2-42da-b18c-165685ad6765','true','id.token.claim'),('9a8f917f-06e2-42da-b18c-165685ad6765','true','introspection.token.claim'),('9a8f917f-06e2-42da-b18c-165685ad6765','String','jsonType.label'),('9a8f917f-06e2-42da-b18c-165685ad6765','gender','user.attribute'),('9a8f917f-06e2-42da-b18c-165685ad6765','true','userinfo.token.claim'),('9b962df6-ac03-416e-af7a-cfb3952a40b1','true','access.token.claim'),('9b962df6-ac03-416e-af7a-cfb3952a40b1','nickname','claim.name'),('9b962df6-ac03-416e-af7a-cfb3952a40b1','true','id.token.claim'),('9b962df6-ac03-416e-af7a-cfb3952a40b1','true','introspection.token.claim'),('9b962df6-ac03-416e-af7a-cfb3952a40b1','String','jsonType.label'),('9b962df6-ac03-416e-af7a-cfb3952a40b1','nickname','user.attribute'),('9b962df6-ac03-416e-af7a-cfb3952a40b1','true','userinfo.token.claim'),('9e0c037e-cd75-4bfc-9df4-437b01de3684','true','access.token.claim'),('9e0c037e-cd75-4bfc-9df4-437b01de3684','email','claim.name'),('9e0c037e-cd75-4bfc-9df4-437b01de3684','true','id.token.claim'),('9e0c037e-cd75-4bfc-9df4-437b01de3684','true','introspection.token.claim'),('9e0c037e-cd75-4bfc-9df4-437b01de3684','String','jsonType.label'),('9e0c037e-cd75-4bfc-9df4-437b01de3684','email','user.attribute'),('9e0c037e-cd75-4bfc-9df4-437b01de3684','true','userinfo.token.claim'),('a10531b0-988e-4d54-8843-2c057423fad2','true','access.token.claim'),('a10531b0-988e-4d54-8843-2c057423fad2','updated_at','claim.name'),('a10531b0-988e-4d54-8843-2c057423fad2','true','id.token.claim'),('a10531b0-988e-4d54-8843-2c057423fad2','true','introspection.token.claim'),('a10531b0-988e-4d54-8843-2c057423fad2','long','jsonType.label'),('a10531b0-988e-4d54-8843-2c057423fad2','updatedAt','user.attribute'),('a10531b0-988e-4d54-8843-2c057423fad2','true','userinfo.token.claim'),('a2cd5654-d33c-477a-b698-27abb94173df','true','access.token.claim'),('a2cd5654-d33c-477a-b698-27abb94173df','zoneinfo','claim.name'),('a2cd5654-d33c-477a-b698-27abb94173df','true','id.token.claim'),('a2cd5654-d33c-477a-b698-27abb94173df','true','introspection.token.claim'),('a2cd5654-d33c-477a-b698-27abb94173df','String','jsonType.label'),('a2cd5654-d33c-477a-b698-27abb94173df','zoneinfo','user.attribute'),('a2cd5654-d33c-477a-b698-27abb94173df','true','userinfo.token.claim'),('aa7ca020-862d-47bb-91c5-e623954b13d1','true','access.token.claim'),('aa7ca020-862d-47bb-91c5-e623954b13d1','true','id.token.claim'),('aa7ca020-862d-47bb-91c5-e623954b13d1','true','introspection.token.claim'),('aa7ca020-862d-47bb-91c5-e623954b13d1','country','user.attribute.country'),('aa7ca020-862d-47bb-91c5-e623954b13d1','formatted','user.attribute.formatted'),('aa7ca020-862d-47bb-91c5-e623954b13d1','locality','user.attribute.locality'),('aa7ca020-862d-47bb-91c5-e623954b13d1','postal_code','user.attribute.postal_code'),('aa7ca020-862d-47bb-91c5-e623954b13d1','region','user.attribute.region'),('aa7ca020-862d-47bb-91c5-e623954b13d1','street','user.attribute.street'),('aa7ca020-862d-47bb-91c5-e623954b13d1','true','userinfo.token.claim'),('ac444e38-613f-4371-bd99-04c6db5599de','true','access.token.claim'),('ac444e38-613f-4371-bd99-04c6db5599de','locale','claim.name'),('ac444e38-613f-4371-bd99-04c6db5599de','true','id.token.claim'),('ac444e38-613f-4371-bd99-04c6db5599de','true','introspection.token.claim'),('ac444e38-613f-4371-bd99-04c6db5599de','String','jsonType.label'),('ac444e38-613f-4371-bd99-04c6db5599de','locale','user.attribute'),('ac444e38-613f-4371-bd99-04c6db5599de','true','userinfo.token.claim'),('b0386818-20bb-4fc8-a716-da9a851aac0f','true','access.token.claim'),('b0386818-20bb-4fc8-a716-da9a851aac0f','true','introspection.token.claim'),('b7478b1b-b796-4202-9306-44c064d8653f','Role','attribute.name'),('b7478b1b-b796-4202-9306-44c064d8653f','Basic','attribute.nameformat'),('b7478b1b-b796-4202-9306-44c064d8653f','false','single'),('b7570251-8cbc-4405-bff3-c3e04e7490b1','true','access.token.claim'),('b7570251-8cbc-4405-bff3-c3e04e7490b1','email','claim.name'),('b7570251-8cbc-4405-bff3-c3e04e7490b1','true','id.token.claim'),('b7570251-8cbc-4405-bff3-c3e04e7490b1','true','introspection.token.claim'),('b7570251-8cbc-4405-bff3-c3e04e7490b1','String','jsonType.label'),('b7570251-8cbc-4405-bff3-c3e04e7490b1','email','user.attribute'),('b7570251-8cbc-4405-bff3-c3e04e7490b1','true','userinfo.token.claim'),('c1689a38-4bb9-4b58-b755-de6bb7ab9ea4','true','access.token.claim'),('c1689a38-4bb9-4b58-b755-de6bb7ab9ea4','locale','claim.name'),('c1689a38-4bb9-4b58-b755-de6bb7ab9ea4','true','id.token.claim'),('c1689a38-4bb9-4b58-b755-de6bb7ab9ea4','true','introspection.token.claim'),('c1689a38-4bb9-4b58-b755-de6bb7ab9ea4','String','jsonType.label'),('c1689a38-4bb9-4b58-b755-de6bb7ab9ea4','locale','user.attribute'),('c1689a38-4bb9-4b58-b755-de6bb7ab9ea4','true','userinfo.token.claim'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','true','access.token.claim'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','true','id.token.claim'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','true','introspection.token.claim'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','country','user.attribute.country'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','formatted','user.attribute.formatted'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','locality','user.attribute.locality'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','postal_code','user.attribute.postal_code'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','region','user.attribute.region'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','street','user.attribute.street'),('c3e3fa43-a510-4ad8-93c7-49d91b553a0d','true','userinfo.token.claim'),('cff060aa-842f-456d-ace5-db03c5215b37','true','access.token.claim'),('cff060aa-842f-456d-ace5-db03c5215b37','true','id.token.claim'),('cff060aa-842f-456d-ace5-db03c5215b37','true','introspection.token.claim'),('cff060aa-842f-456d-ace5-db03c5215b37','country','user.attribute.country'),('cff060aa-842f-456d-ace5-db03c5215b37','formatted','user.attribute.formatted'),('cff060aa-842f-456d-ace5-db03c5215b37','locality','user.attribute.locality'),('cff060aa-842f-456d-ace5-db03c5215b37','postal_code','user.attribute.postal_code'),('cff060aa-842f-456d-ace5-db03c5215b37','region','user.attribute.region'),('cff060aa-842f-456d-ace5-db03c5215b37','street','user.attribute.street'),('cff060aa-842f-456d-ace5-db03c5215b37','true','userinfo.token.claim'),('d21c25d8-c749-49f6-a8b6-d77fb0d9719c','true','access.token.claim'),('d21c25d8-c749-49f6-a8b6-d77fb0d9719c','middle_name','claim.name'),('d21c25d8-c749-49f6-a8b6-d77fb0d9719c','true','id.token.claim'),('d21c25d8-c749-49f6-a8b6-d77fb0d9719c','true','introspection.token.claim'),('d21c25d8-c749-49f6-a8b6-d77fb0d9719c','String','jsonType.label'),('d21c25d8-c749-49f6-a8b6-d77fb0d9719c','middleName','user.attribute'),('d21c25d8-c749-49f6-a8b6-d77fb0d9719c','true','userinfo.token.claim'),('d8556613-0e26-4ce1-a3a0-21e4d595ab02','true','access.token.claim'),('d8556613-0e26-4ce1-a3a0-21e4d595ab02','preferred_username','claim.name'),('d8556613-0e26-4ce1-a3a0-21e4d595ab02','true','id.token.claim'),('d8556613-0e26-4ce1-a3a0-21e4d595ab02','true','introspection.token.claim'),('d8556613-0e26-4ce1-a3a0-21e4d595ab02','String','jsonType.label'),('d8556613-0e26-4ce1-a3a0-21e4d595ab02','username','user.attribute'),('d8556613-0e26-4ce1-a3a0-21e4d595ab02','true','userinfo.token.claim'),('dac1e8be-fcad-4a6c-b751-e96396c173b7','true','access.token.claim'),('dac1e8be-fcad-4a6c-b751-e96396c173b7','upn','claim.name'),('dac1e8be-fcad-4a6c-b751-e96396c173b7','true','id.token.claim'),('dac1e8be-fcad-4a6c-b751-e96396c173b7','true','introspection.token.claim'),('dac1e8be-fcad-4a6c-b751-e96396c173b7','String','jsonType.label'),('dac1e8be-fcad-4a6c-b751-e96396c173b7','username','user.attribute'),('dac1e8be-fcad-4a6c-b751-e96396c173b7','true','userinfo.token.claim'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','true','access.token.claim'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','test_role','claim.name'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','true','id.token.claim'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','true','introspection.token.claim'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','String','jsonType.label'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','false','lightweight.claim'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','test_role','user.attribute'),('de013a82-70d6-4fd2-9093-f0a0dea05dba','true','userinfo.token.claim'),('e1497be5-5c93-45cc-9e43-d2a26eeb567a','true','access.token.claim'),('e1497be5-5c93-45cc-9e43-d2a26eeb567a','true','introspection.token.claim'),('e1e61070-5173-494a-9432-5b38cc69cf96','true','access.token.claim'),('e1e61070-5173-494a-9432-5b38cc69cf96','family_name','claim.name'),('e1e61070-5173-494a-9432-5b38cc69cf96','true','id.token.claim'),('e1e61070-5173-494a-9432-5b38cc69cf96','true','introspection.token.claim'),('e1e61070-5173-494a-9432-5b38cc69cf96','String','jsonType.label'),('e1e61070-5173-494a-9432-5b38cc69cf96','lastName','user.attribute'),('e1e61070-5173-494a-9432-5b38cc69cf96','true','userinfo.token.claim'),('e4fe6688-9f00-4cff-9bd9-0610c77ede27','true','access.token.claim'),('e4fe6688-9f00-4cff-9bd9-0610c77ede27','locale','claim.name'),('e4fe6688-9f00-4cff-9bd9-0610c77ede27','true','id.token.claim'),('e4fe6688-9f00-4cff-9bd9-0610c77ede27','true','introspection.token.claim'),('e4fe6688-9f00-4cff-9bd9-0610c77ede27','String','jsonType.label'),('e4fe6688-9f00-4cff-9bd9-0610c77ede27','locale','user.attribute'),('e4fe6688-9f00-4cff-9bd9-0610c77ede27','true','userinfo.token.claim'),('e97a6c0d-41f2-4d65-a6a5-0d2802370f47','true','access.token.claim'),('e97a6c0d-41f2-4d65-a6a5-0d2802370f47','updated_at','claim.name'),('e97a6c0d-41f2-4d65-a6a5-0d2802370f47','true','id.token.claim'),('e97a6c0d-41f2-4d65-a6a5-0d2802370f47','true','introspection.token.claim'),('e97a6c0d-41f2-4d65-a6a5-0d2802370f47','long','jsonType.label'),('e97a6c0d-41f2-4d65-a6a5-0d2802370f47','updatedAt','user.attribute'),('e97a6c0d-41f2-4d65-a6a5-0d2802370f47','true','userinfo.token.claim'),('eb0979ca-36e1-4b3d-93ac-6d84cdfb0a63','true','access.token.claim'),('eb0979ca-36e1-4b3d-93ac-6d84cdfb0a63','locale','claim.name'),('eb0979ca-36e1-4b3d-93ac-6d84cdfb0a63','true','id.token.claim'),('eb0979ca-36e1-4b3d-93ac-6d84cdfb0a63','true','introspection.token.claim'),('eb0979ca-36e1-4b3d-93ac-6d84cdfb0a63','String','jsonType.label'),('eb0979ca-36e1-4b3d-93ac-6d84cdfb0a63','locale','user.attribute'),('eb0979ca-36e1-4b3d-93ac-6d84cdfb0a63','true','userinfo.token.claim'),('ec272038-96c7-4eec-b0a4-139653d08cfe','true','access.token.claim'),('ec272038-96c7-4eec-b0a4-139653d08cfe','gender','claim.name'),('ec272038-96c7-4eec-b0a4-139653d08cfe','true','id.token.claim'),('ec272038-96c7-4eec-b0a4-139653d08cfe','true','introspection.token.claim'),('ec272038-96c7-4eec-b0a4-139653d08cfe','String','jsonType.label'),('ec272038-96c7-4eec-b0a4-139653d08cfe','gender','user.attribute'),('ec272038-96c7-4eec-b0a4-139653d08cfe','true','userinfo.token.claim'),('ee94a117-33ff-45f4-a2ec-56c15719d2ab','true','access.token.claim'),('ee94a117-33ff-45f4-a2ec-56c15719d2ab','nickname','claim.name'),('ee94a117-33ff-45f4-a2ec-56c15719d2ab','true','id.token.claim'),('ee94a117-33ff-45f4-a2ec-56c15719d2ab','true','introspection.token.claim'),('ee94a117-33ff-45f4-a2ec-56c15719d2ab','String','jsonType.label'),('ee94a117-33ff-45f4-a2ec-56c15719d2ab','nickname','user.attribute'),('ee94a117-33ff-45f4-a2ec-56c15719d2ab','true','userinfo.token.claim'),('ef665f48-8d78-4e5d-b6aa-fde9fbe8dbaf','true','access.token.claim'),('ef665f48-8d78-4e5d-b6aa-fde9fbe8dbaf','phone_number','claim.name'),('ef665f48-8d78-4e5d-b6aa-fde9fbe8dbaf','true','id.token.claim'),('ef665f48-8d78-4e5d-b6aa-fde9fbe8dbaf','true','introspection.token.claim'),('ef665f48-8d78-4e5d-b6aa-fde9fbe8dbaf','String','jsonType.label'),('ef665f48-8d78-4e5d-b6aa-fde9fbe8dbaf','phoneNumber','user.attribute'),('ef665f48-8d78-4e5d-b6aa-fde9fbe8dbaf','true','userinfo.token.claim'),('f1b0cdf5-e7d2-4e4d-897a-e6e41a2d46ec','true','access.token.claim'),('f1b0cdf5-e7d2-4e4d-897a-e6e41a2d46ec','middle_name','claim.name'),('f1b0cdf5-e7d2-4e4d-897a-e6e41a2d46ec','true','id.token.claim'),('f1b0cdf5-e7d2-4e4d-897a-e6e41a2d46ec','true','introspection.token.claim'),('f1b0cdf5-e7d2-4e4d-897a-e6e41a2d46ec','String','jsonType.label'),('f1b0cdf5-e7d2-4e4d-897a-e6e41a2d46ec','middleName','user.attribute'),('f1b0cdf5-e7d2-4e4d-897a-e6e41a2d46ec','true','userinfo.token.claim'),('f5802976-75fa-43c1-a376-89b4fb745969','true','access.token.claim'),('f5802976-75fa-43c1-a376-89b4fb745969','resource_access.${client_id}.roles','claim.name'),('f5802976-75fa-43c1-a376-89b4fb745969','true','introspection.token.claim'),('f5802976-75fa-43c1-a376-89b4fb745969','String','jsonType.label'),('f5802976-75fa-43c1-a376-89b4fb745969','true','multivalued'),('f5802976-75fa-43c1-a376-89b4fb745969','foo','user.attribute'),('f5f6bffa-a5c4-475a-abed-55038e3c3896','true','access.token.claim'),('f5f6bffa-a5c4-475a-abed-55038e3c3896','true','id.token.claim'),('f5f6bffa-a5c4-475a-abed-55038e3c3896','true','introspection.token.claim'),('f5f6bffa-a5c4-475a-abed-55038e3c3896','true','userinfo.token.claim'),('f6ea5b8b-104d-4e3f-a142-73cb07f9c824','true','access.token.claim'),('f6ea5b8b-104d-4e3f-a142-73cb07f9c824','groups','claim.name'),('f6ea5b8b-104d-4e3f-a142-73cb07f9c824','true','id.token.claim'),('f6ea5b8b-104d-4e3f-a142-73cb07f9c824','true','introspection.token.claim'),('f6ea5b8b-104d-4e3f-a142-73cb07f9c824','String','jsonType.label'),('f6ea5b8b-104d-4e3f-a142-73cb07f9c824','true','multivalued'),('f6ea5b8b-104d-4e3f-a142-73cb07f9c824','foo','user.attribute'),('f717b364-1349-451d-b084-5bc291c85208','true','access.token.claim'),('f717b364-1349-451d-b084-5bc291c85208','email','claim.name'),('f717b364-1349-451d-b084-5bc291c85208','true','id.token.claim'),('f717b364-1349-451d-b084-5bc291c85208','true','introspection.token.claim'),('f717b364-1349-451d-b084-5bc291c85208','String','jsonType.label'),('f717b364-1349-451d-b084-5bc291c85208','email','user.attribute'),('f717b364-1349-451d-b084-5bc291c85208','true','userinfo.token.claim'),('f7889d19-a15c-4f01-9c8c-37d528740726','Role','attribute.name'),('f7889d19-a15c-4f01-9c8c-37d528740726','Basic','attribute.nameformat'),('f7889d19-a15c-4f01-9c8c-37d528740726','false','single'),('faa279e3-db18-4766-8ebd-15d3700d1759','true','access.token.claim'),('faa279e3-db18-4766-8ebd-15d3700d1759','locale','claim.name'),('faa279e3-db18-4766-8ebd-15d3700d1759','true','id.token.claim'),('faa279e3-db18-4766-8ebd-15d3700d1759','true','introspection.token.claim'),('faa279e3-db18-4766-8ebd-15d3700d1759','String','jsonType.label'),('faa279e3-db18-4766-8ebd-15d3700d1759','locale','user.attribute'),('faa279e3-db18-4766-8ebd-15d3700d1759','true','userinfo.token.claim'),('fb6dfb7d-1a8c-4227-b3d9-ba952b697162','true','access.token.claim'),('fb6dfb7d-1a8c-4227-b3d9-ba952b697162','true','introspection.token.claim'),('ff6772e5-ec2e-4af7-9265-16326bee2c6b','true','access.token.claim'),('ff6772e5-ec2e-4af7-9265-16326bee2c6b','true','id.token.claim'),('ff6772e5-ec2e-4af7-9265-16326bee2c6b','true','introspection.token.claim'),('ff6772e5-ec2e-4af7-9265-16326bee2c6b','true','userinfo.token.claim');
/*!40000 ALTER TABLE `PROTOCOL_MAPPER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM`
--

DROP TABLE IF EXISTS `REALM`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM` (
  `ID` varchar(36) NOT NULL,
  `ACCESS_CODE_LIFESPAN` int DEFAULT NULL,
  `USER_ACTION_LIFESPAN` int DEFAULT NULL,
  `ACCESS_TOKEN_LIFESPAN` int DEFAULT NULL,
  `ACCOUNT_THEME` varchar(255) DEFAULT NULL,
  `ADMIN_THEME` varchar(255) DEFAULT NULL,
  `EMAIL_THEME` varchar(255) DEFAULT NULL,
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `EVENTS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `EVENTS_EXPIRATION` bigint DEFAULT NULL,
  `LOGIN_THEME` varchar(255) DEFAULT NULL,
  `NAME` varchar(255) DEFAULT NULL,
  `NOT_BEFORE` int DEFAULT NULL,
  `PASSWORD_POLICY` text,
  `REGISTRATION_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `REMEMBER_ME` tinyint NOT NULL DEFAULT '0',
  `RESET_PASSWORD_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `SOCIAL` tinyint NOT NULL DEFAULT '0',
  `SSL_REQUIRED` varchar(255) DEFAULT NULL,
  `SSO_IDLE_TIMEOUT` int DEFAULT NULL,
  `SSO_MAX_LIFESPAN` int DEFAULT NULL,
  `UPDATE_PROFILE_ON_SOC_LOGIN` tinyint NOT NULL DEFAULT '0',
  `VERIFY_EMAIL` tinyint NOT NULL DEFAULT '0',
  `MASTER_ADMIN_CLIENT` varchar(36) DEFAULT NULL,
  `LOGIN_LIFESPAN` int DEFAULT NULL,
  `INTERNATIONALIZATION_ENABLED` tinyint NOT NULL DEFAULT '0',
  `DEFAULT_LOCALE` varchar(255) DEFAULT NULL,
  `REG_EMAIL_AS_USERNAME` tinyint NOT NULL DEFAULT '0',
  `ADMIN_EVENTS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `ADMIN_EVENTS_DETAILS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `EDIT_USERNAME_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `OTP_POLICY_COUNTER` int DEFAULT '0',
  `OTP_POLICY_WINDOW` int DEFAULT '1',
  `OTP_POLICY_PERIOD` int DEFAULT '30',
  `OTP_POLICY_DIGITS` int DEFAULT '6',
  `OTP_POLICY_ALG` varchar(36) DEFAULT 'HmacSHA1',
  `OTP_POLICY_TYPE` varchar(36) DEFAULT 'totp',
  `BROWSER_FLOW` varchar(36) DEFAULT NULL,
  `REGISTRATION_FLOW` varchar(36) DEFAULT NULL,
  `DIRECT_GRANT_FLOW` varchar(36) DEFAULT NULL,
  `RESET_CREDENTIALS_FLOW` varchar(36) DEFAULT NULL,
  `CLIENT_AUTH_FLOW` varchar(36) DEFAULT NULL,
  `OFFLINE_SESSION_IDLE_TIMEOUT` int DEFAULT '0',
  `REVOKE_REFRESH_TOKEN` tinyint NOT NULL DEFAULT '0',
  `ACCESS_TOKEN_LIFE_IMPLICIT` int DEFAULT '0',
  `LOGIN_WITH_EMAIL_ALLOWED` tinyint NOT NULL DEFAULT '1',
  `DUPLICATE_EMAILS_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `DOCKER_AUTH_FLOW` varchar(36) DEFAULT NULL,
  `REFRESH_TOKEN_MAX_REUSE` int DEFAULT '0',
  `ALLOW_USER_MANAGED_ACCESS` tinyint NOT NULL DEFAULT '0',
  `SSO_MAX_LIFESPAN_REMEMBER_ME` int NOT NULL,
  `SSO_IDLE_TIMEOUT_REMEMBER_ME` int NOT NULL,
  `DEFAULT_ROLE` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_ORVSDMLA56612EAEFIQ6WL5OI` (`NAME`),
  KEY `IDX_REALM_MASTER_ADM_CLI` (`MASTER_ADMIN_CLIENT`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM`
--

LOCK TABLES `REALM` WRITE;
/*!40000 ALTER TABLE `REALM` DISABLE KEYS */;
INSERT INTO `REALM` (`ID`, `ACCESS_CODE_LIFESPAN`, `USER_ACTION_LIFESPAN`, `ACCESS_TOKEN_LIFESPAN`, `ACCOUNT_THEME`, `ADMIN_THEME`, `EMAIL_THEME`, `ENABLED`, `EVENTS_ENABLED`, `EVENTS_EXPIRATION`, `LOGIN_THEME`, `NAME`, `NOT_BEFORE`, `PASSWORD_POLICY`, `REGISTRATION_ALLOWED`, `REMEMBER_ME`, `RESET_PASSWORD_ALLOWED`, `SOCIAL`, `SSL_REQUIRED`, `SSO_IDLE_TIMEOUT`, `SSO_MAX_LIFESPAN`, `UPDATE_PROFILE_ON_SOC_LOGIN`, `VERIFY_EMAIL`, `MASTER_ADMIN_CLIENT`, `LOGIN_LIFESPAN`, `INTERNATIONALIZATION_ENABLED`, `DEFAULT_LOCALE`, `REG_EMAIL_AS_USERNAME`, `ADMIN_EVENTS_ENABLED`, `ADMIN_EVENTS_DETAILS_ENABLED`, `EDIT_USERNAME_ALLOWED`, `OTP_POLICY_COUNTER`, `OTP_POLICY_WINDOW`, `OTP_POLICY_PERIOD`, `OTP_POLICY_DIGITS`, `OTP_POLICY_ALG`, `OTP_POLICY_TYPE`, `BROWSER_FLOW`, `REGISTRATION_FLOW`, `DIRECT_GRANT_FLOW`, `RESET_CREDENTIALS_FLOW`, `CLIENT_AUTH_FLOW`, `OFFLINE_SESSION_IDLE_TIMEOUT`, `REVOKE_REFRESH_TOKEN`, `ACCESS_TOKEN_LIFE_IMPLICIT`, `LOGIN_WITH_EMAIL_ALLOWED`, `DUPLICATE_EMAILS_ALLOWED`, `DOCKER_AUTH_FLOW`, `REFRESH_TOKEN_MAX_REUSE`, `ALLOW_USER_MANAGED_ACCESS`, `SSO_MAX_LIFESPAN_REMEMBER_ME`, `SSO_IDLE_TIMEOUT_REMEMBER_ME`, `DEFAULT_ROLE`) VALUES ('40ae881c-f4e4-4b07-b097-a67d2bf515e6',60,300,300,'','','custom',1,0,0,'custom','test-realm',0,NULL,1,1,0,0,'EXTERNAL',1800,36000,0,0,'1aaa179d-deb7-4fea-9f6b-bae7092c8ba7',1800,0,NULL,1,0,0,0,0,1,30,6,'HmacSHA1','totp','4a21f571-4b69-4ef8-92af-09fe24026b64','08b0543c-51a6-4200-bb72-38404a916954','73f81fba-ce8f-4131-8637-3c64b1323d86','d5b7af02-aba4-4b1c-bcca-e068b83dfee1','a9a6d8d1-afbc-447a-8685-cd855fb3e627',2592000,0,900,1,0,'ea3b8b38-9a74-403e-8058-be535ccea217',0,0,0,0,'c8496af8-af9d-49ea-9bb2-26ba4bcbef2f'),('4327ba47-4116-44ea-9c4d-02907dca81e7',60,300,60,'keycloak.v3','keycloak.v2','custom',1,0,0,'custom','master',0,NULL,0,0,0,0,'EXTERNAL',1800,36000,0,0,'b79d0e1c-6d44-4f01-b111-278fe3db31ee',1800,0,NULL,0,0,0,0,0,1,30,6,'HmacSHA1','totp','7bf8b92b-a6be-4924-a010-8750f64b92e9','8b55c406-d1d6-4a0c-bffa-4f2bf77f97f3','3af9837f-1ac4-4fe2-a3a7-579df7d70cbd','616a6379-b1b2-4e14-96ea-db9d45b16538','f378b286-fd54-4b8f-8615-df0118dacd24',2592000,0,900,1,0,'692a7ac4-e4aa-47ac-8bb4-5267c7930f6d',0,0,0,0,'e7b5b42b-09a2-46bf-b6b5-87f71892c9dd'),('dcc080c5-aede-4fd3-8b01-bd0928b730a2',60,300,300,'keycloak.v3','keycloak.v2','custom',1,0,0,'custom','TEST',0,NULL,1,1,0,0,'EXTERNAL',1800,36000,0,0,'72c6029f-f8d2-4256-a326-2642c15f3a1e',1800,0,NULL,1,0,0,0,0,1,30,6,'HmacSHA1','totp','6bfae19d-d2c6-4dd5-836c-6042247853cc','f77f953b-5e25-4601-b047-748724d795b3','b639104c-892c-4843-a2ed-ca4324be5788','06d58040-a014-4f13-a9ea-72835ad3c576','f84a1f13-ef43-4a16-b9d6-63fe787b86dd',2592000,0,900,1,0,'1250dad6-ebc6-4a7e-a847-a83a67967a51',0,0,0,0,'14fce4d6-46b4-4c14-bd14-4ac7227ea3b0');
/*!40000 ALTER TABLE `REALM` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM_ATTRIBUTE`
--

DROP TABLE IF EXISTS `REALM_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_ATTRIBUTE` (
  `NAME` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  `VALUE` longtext CHARACTER SET utf8 COLLATE utf8_general_ci,
  PRIMARY KEY (`NAME`,`REALM_ID`),
  KEY `IDX_REALM_ATTR_REALM` (`REALM_ID`),
  CONSTRAINT `FK_8SHXD6L3E9ATQUKACXGPFFPTW` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM_ATTRIBUTE`
--

LOCK TABLES `REALM_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `REALM_ATTRIBUTE` DISABLE KEYS */;
INSERT INTO `REALM_ATTRIBUTE` (`NAME`, `REALM_ID`, `VALUE`) VALUES ('_browser_header.contentSecurityPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','frame-src \'self\'; frame-ancestors \'self\'; object-src \'none\';'),('_browser_header.contentSecurityPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','frame-src \'self\'; frame-ancestors \'self\'; object-src \'none\';'),('_browser_header.contentSecurityPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','frame-src \'self\'; frame-ancestors \'self\'; object-src \'none\';'),('_browser_header.contentSecurityPolicyReportOnly','40ae881c-f4e4-4b07-b097-a67d2bf515e6',''),('_browser_header.contentSecurityPolicyReportOnly','4327ba47-4116-44ea-9c4d-02907dca81e7',''),('_browser_header.contentSecurityPolicyReportOnly','dcc080c5-aede-4fd3-8b01-bd0928b730a2',''),('_browser_header.referrerPolicy','40ae881c-f4e4-4b07-b097-a67d2bf515e6','no-referrer'),('_browser_header.referrerPolicy','4327ba47-4116-44ea-9c4d-02907dca81e7','no-referrer'),('_browser_header.referrerPolicy','dcc080c5-aede-4fd3-8b01-bd0928b730a2','no-referrer'),('_browser_header.strictTransportSecurity','40ae881c-f4e4-4b07-b097-a67d2bf515e6','max-age=31536000; includeSubDomains'),('_browser_header.strictTransportSecurity','4327ba47-4116-44ea-9c4d-02907dca81e7','max-age=31536000; includeSubDomains'),('_browser_header.strictTransportSecurity','dcc080c5-aede-4fd3-8b01-bd0928b730a2','max-age=31536000; includeSubDomains'),('_browser_header.xContentTypeOptions','40ae881c-f4e4-4b07-b097-a67d2bf515e6','nosniff'),('_browser_header.xContentTypeOptions','4327ba47-4116-44ea-9c4d-02907dca81e7','nosniff'),('_browser_header.xContentTypeOptions','dcc080c5-aede-4fd3-8b01-bd0928b730a2','nosniff'),('_browser_header.xFrameOptions','40ae881c-f4e4-4b07-b097-a67d2bf515e6','SAMEORIGIN'),('_browser_header.xFrameOptions','4327ba47-4116-44ea-9c4d-02907dca81e7','SAMEORIGIN'),('_browser_header.xFrameOptions','dcc080c5-aede-4fd3-8b01-bd0928b730a2','SAMEORIGIN'),('_browser_header.xRobotsTag','40ae881c-f4e4-4b07-b097-a67d2bf515e6','none'),('_browser_header.xRobotsTag','4327ba47-4116-44ea-9c4d-02907dca81e7','none'),('_browser_header.xRobotsTag','dcc080c5-aede-4fd3-8b01-bd0928b730a2','none'),('_browser_header.xXSSProtection','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1; mode=block'),('_browser_header.xXSSProtection','4327ba47-4116-44ea-9c4d-02907dca81e7','1; mode=block'),('_browser_header.xXSSProtection','dcc080c5-aede-4fd3-8b01-bd0928b730a2','1; mode=block'),('acr.loa.map','40ae881c-f4e4-4b07-b097-a67d2bf515e6','{}'),('acr.loa.map','dcc080c5-aede-4fd3-8b01-bd0928b730a2','{}'),('actionTokenGeneratedByAdminLifespan','40ae881c-f4e4-4b07-b097-a67d2bf515e6','43200'),('actionTokenGeneratedByAdminLifespan','4327ba47-4116-44ea-9c4d-02907dca81e7','43200'),('actionTokenGeneratedByAdminLifespan','dcc080c5-aede-4fd3-8b01-bd0928b730a2','43200'),('actionTokenGeneratedByUserLifespan','40ae881c-f4e4-4b07-b097-a67d2bf515e6','300'),('actionTokenGeneratedByUserLifespan','4327ba47-4116-44ea-9c4d-02907dca81e7','300'),('actionTokenGeneratedByUserLifespan','dcc080c5-aede-4fd3-8b01-bd0928b730a2','300'),('bruteForceProtected','40ae881c-f4e4-4b07-b097-a67d2bf515e6','true'),('bruteForceProtected','4327ba47-4116-44ea-9c4d-02907dca81e7','false'),('bruteForceProtected','dcc080c5-aede-4fd3-8b01-bd0928b730a2','true'),('cibaAuthRequestedUserHint','40ae881c-f4e4-4b07-b097-a67d2bf515e6','login_hint'),('cibaAuthRequestedUserHint','4327ba47-4116-44ea-9c4d-02907dca81e7','login_hint'),('cibaAuthRequestedUserHint','dcc080c5-aede-4fd3-8b01-bd0928b730a2','login_hint'),('cibaBackchannelTokenDeliveryMode','40ae881c-f4e4-4b07-b097-a67d2bf515e6','poll'),('cibaBackchannelTokenDeliveryMode','4327ba47-4116-44ea-9c4d-02907dca81e7','poll'),('cibaBackchannelTokenDeliveryMode','dcc080c5-aede-4fd3-8b01-bd0928b730a2','poll'),('cibaExpiresIn','40ae881c-f4e4-4b07-b097-a67d2bf515e6','120'),('cibaExpiresIn','4327ba47-4116-44ea-9c4d-02907dca81e7','120'),('cibaExpiresIn','dcc080c5-aede-4fd3-8b01-bd0928b730a2','120'),('cibaInterval','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5'),('cibaInterval','4327ba47-4116-44ea-9c4d-02907dca81e7','5'),('cibaInterval','dcc080c5-aede-4fd3-8b01-bd0928b730a2','5'),('client-policies.policies','40ae881c-f4e4-4b07-b097-a67d2bf515e6','{\"policies\":[]}'),('client-policies.policies','4327ba47-4116-44ea-9c4d-02907dca81e7','{\"policies\":[]}'),('client-policies.policies','dcc080c5-aede-4fd3-8b01-bd0928b730a2','{\"policies\":[]}'),('client-policies.profiles','40ae881c-f4e4-4b07-b097-a67d2bf515e6','{\"profiles\":[]}'),('client-policies.profiles','4327ba47-4116-44ea-9c4d-02907dca81e7','{\"profiles\":[]}'),('client-policies.profiles','dcc080c5-aede-4fd3-8b01-bd0928b730a2','{\"profiles\":[]}'),('clientOfflineSessionIdleTimeout','40ae881c-f4e4-4b07-b097-a67d2bf515e6','0'),('clientOfflineSessionIdleTimeout','4327ba47-4116-44ea-9c4d-02907dca81e7','0'),('clientOfflineSessionIdleTimeout','dcc080c5-aede-4fd3-8b01-bd0928b730a2','0'),('clientOfflineSessionMaxLifespan','40ae881c-f4e4-4b07-b097-a67d2bf515e6','0'),('clientOfflineSessionMaxLifespan','4327ba47-4116-44ea-9c4d-02907dca81e7','0'),('clientOfflineSessionMaxLifespan','dcc080c5-aede-4fd3-8b01-bd0928b730a2','0'),('clientSessionIdleTimeout','40ae881c-f4e4-4b07-b097-a67d2bf515e6','0'),('clientSessionIdleTimeout','4327ba47-4116-44ea-9c4d-02907dca81e7','0'),('clientSessionIdleTimeout','dcc080c5-aede-4fd3-8b01-bd0928b730a2','0'),('clientSessionMaxLifespan','40ae881c-f4e4-4b07-b097-a67d2bf515e6','0'),('clientSessionMaxLifespan','4327ba47-4116-44ea-9c4d-02907dca81e7','0'),('clientSessionMaxLifespan','dcc080c5-aede-4fd3-8b01-bd0928b730a2','0'),('defaultSignatureAlgorithm','40ae881c-f4e4-4b07-b097-a67d2bf515e6','RS256'),('defaultSignatureAlgorithm','4327ba47-4116-44ea-9c4d-02907dca81e7','RS256'),('defaultSignatureAlgorithm','dcc080c5-aede-4fd3-8b01-bd0928b730a2','RS256'),('displayName','40ae881c-f4e4-4b07-b097-a67d2bf515e6','시험삼아 만들어본 Realm'),('displayName','4327ba47-4116-44ea-9c4d-02907dca81e7','Keycloak'),('displayName','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ARC Brain 개발환경(dev)'),('displayNameHtml','40ae881c-f4e4-4b07-b097-a67d2bf515e6',''),('displayNameHtml','4327ba47-4116-44ea-9c4d-02907dca81e7','<div class=\"kc-logo-text\"><span>Keycloak</span></div>'),('displayNameHtml','dcc080c5-aede-4fd3-8b01-bd0928b730a2',''),('failureFactor','40ae881c-f4e4-4b07-b097-a67d2bf515e6','30'),('failureFactor','4327ba47-4116-44ea-9c4d-02907dca81e7','30'),('failureFactor','dcc080c5-aede-4fd3-8b01-bd0928b730a2','5'),('firstBrokerLoginFlowId','40ae881c-f4e4-4b07-b097-a67d2bf515e6','409007a9-4804-45b2-a701-79b9f94cea54'),('firstBrokerLoginFlowId','4327ba47-4116-44ea-9c4d-02907dca81e7','48fb5d06-357a-463e-9263-1de5c87f092c'),('firstBrokerLoginFlowId','dcc080c5-aede-4fd3-8b01-bd0928b730a2','f97a7702-55e0-4f57-b4a1-ae4b09239fa8'),('frontendUrl','40ae881c-f4e4-4b07-b097-a67d2bf515e6',''),('frontendUrl','dcc080c5-aede-4fd3-8b01-bd0928b730a2',''),('maxDeltaTimeSeconds','40ae881c-f4e4-4b07-b097-a67d2bf515e6','43200'),('maxDeltaTimeSeconds','4327ba47-4116-44ea-9c4d-02907dca81e7','43200'),('maxDeltaTimeSeconds','dcc080c5-aede-4fd3-8b01-bd0928b730a2','43200'),('maxFailureWaitSeconds','40ae881c-f4e4-4b07-b097-a67d2bf515e6','900'),('maxFailureWaitSeconds','4327ba47-4116-44ea-9c4d-02907dca81e7','900'),('maxFailureWaitSeconds','dcc080c5-aede-4fd3-8b01-bd0928b730a2','900'),('maxTemporaryLockouts','40ae881c-f4e4-4b07-b097-a67d2bf515e6','0'),('maxTemporaryLockouts','4327ba47-4116-44ea-9c4d-02907dca81e7','0'),('maxTemporaryLockouts','dcc080c5-aede-4fd3-8b01-bd0928b730a2','0'),('minimumQuickLoginWaitSeconds','40ae881c-f4e4-4b07-b097-a67d2bf515e6','60'),('minimumQuickLoginWaitSeconds','4327ba47-4116-44ea-9c4d-02907dca81e7','60'),('minimumQuickLoginWaitSeconds','dcc080c5-aede-4fd3-8b01-bd0928b730a2','60'),('oauth2DeviceCodeLifespan','40ae881c-f4e4-4b07-b097-a67d2bf515e6','600'),('oauth2DeviceCodeLifespan','4327ba47-4116-44ea-9c4d-02907dca81e7','600'),('oauth2DeviceCodeLifespan','dcc080c5-aede-4fd3-8b01-bd0928b730a2','600'),('oauth2DevicePollingInterval','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5'),('oauth2DevicePollingInterval','4327ba47-4116-44ea-9c4d-02907dca81e7','5'),('oauth2DevicePollingInterval','dcc080c5-aede-4fd3-8b01-bd0928b730a2','5'),('offlineSessionMaxLifespan','40ae881c-f4e4-4b07-b097-a67d2bf515e6','5184000'),('offlineSessionMaxLifespan','4327ba47-4116-44ea-9c4d-02907dca81e7','5184000'),('offlineSessionMaxLifespan','dcc080c5-aede-4fd3-8b01-bd0928b730a2','5184000'),('offlineSessionMaxLifespanEnabled','40ae881c-f4e4-4b07-b097-a67d2bf515e6','false'),('offlineSessionMaxLifespanEnabled','4327ba47-4116-44ea-9c4d-02907dca81e7','false'),('offlineSessionMaxLifespanEnabled','dcc080c5-aede-4fd3-8b01-bd0928b730a2','false'),('parRequestUriLifespan','40ae881c-f4e4-4b07-b097-a67d2bf515e6','60'),('parRequestUriLifespan','4327ba47-4116-44ea-9c4d-02907dca81e7','60'),('parRequestUriLifespan','dcc080c5-aede-4fd3-8b01-bd0928b730a2','60'),('permanentLockout','40ae881c-f4e4-4b07-b097-a67d2bf515e6','false'),('permanentLockout','4327ba47-4116-44ea-9c4d-02907dca81e7','false'),('permanentLockout','dcc080c5-aede-4fd3-8b01-bd0928b730a2','false'),('quickLoginCheckMilliSeconds','40ae881c-f4e4-4b07-b097-a67d2bf515e6','1000'),('quickLoginCheckMilliSeconds','4327ba47-4116-44ea-9c4d-02907dca81e7','1000'),('quickLoginCheckMilliSeconds','dcc080c5-aede-4fd3-8b01-bd0928b730a2','1000'),('realmReusableOtpCode','40ae881c-f4e4-4b07-b097-a67d2bf515e6','false'),('realmReusableOtpCode','4327ba47-4116-44ea-9c4d-02907dca81e7','false'),('realmReusableOtpCode','dcc080c5-aede-4fd3-8b01-bd0928b730a2','false'),('waitIncrementSeconds','40ae881c-f4e4-4b07-b097-a67d2bf515e6','60'),('waitIncrementSeconds','4327ba47-4116-44ea-9c4d-02907dca81e7','60'),('waitIncrementSeconds','dcc080c5-aede-4fd3-8b01-bd0928b730a2','60'),('webAuthnPolicyAttestationConveyancePreference','40ae881c-f4e4-4b07-b097-a67d2bf515e6','not specified'),('webAuthnPolicyAttestationConveyancePreference','4327ba47-4116-44ea-9c4d-02907dca81e7','not specified'),('webAuthnPolicyAttestationConveyancePreference','dcc080c5-aede-4fd3-8b01-bd0928b730a2','not specified'),('webAuthnPolicyAttestationConveyancePreferencePasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6','not specified'),('webAuthnPolicyAttestationConveyancePreferencePasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7','not specified'),('webAuthnPolicyAttestationConveyancePreferencePasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2','not specified'),('webAuthnPolicyAuthenticatorAttachment','40ae881c-f4e4-4b07-b097-a67d2bf515e6','not specified'),('webAuthnPolicyAuthenticatorAttachment','4327ba47-4116-44ea-9c4d-02907dca81e7','not specified'),('webAuthnPolicyAuthenticatorAttachment','dcc080c5-aede-4fd3-8b01-bd0928b730a2','not specified'),('webAuthnPolicyAuthenticatorAttachmentPasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6','not specified'),('webAuthnPolicyAuthenticatorAttachmentPasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7','not specified'),('webAuthnPolicyAuthenticatorAttachmentPasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2','not specified'),('webAuthnPolicyAvoidSameAuthenticatorRegister','40ae881c-f4e4-4b07-b097-a67d2bf515e6','false'),('webAuthnPolicyAvoidSameAuthenticatorRegister','4327ba47-4116-44ea-9c4d-02907dca81e7','false'),('webAuthnPolicyAvoidSameAuthenticatorRegister','dcc080c5-aede-4fd3-8b01-bd0928b730a2','false'),('webAuthnPolicyAvoidSameAuthenticatorRegisterPasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6','false'),('webAuthnPolicyAvoidSameAuthenticatorRegisterPasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7','false'),('webAuthnPolicyAvoidSameAuthenticatorRegisterPasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2','false'),('webAuthnPolicyCreateTimeout','40ae881c-f4e4-4b07-b097-a67d2bf515e6','0'),('webAuthnPolicyCreateTimeout','4327ba47-4116-44ea-9c4d-02907dca81e7','0'),('webAuthnPolicyCreateTimeout','dcc080c5-aede-4fd3-8b01-bd0928b730a2','0'),('webAuthnPolicyCreateTimeoutPasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6','0'),('webAuthnPolicyCreateTimeoutPasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7','0'),('webAuthnPolicyCreateTimeoutPasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2','0'),('webAuthnPolicyRequireResidentKey','40ae881c-f4e4-4b07-b097-a67d2bf515e6','not specified'),('webAuthnPolicyRequireResidentKey','4327ba47-4116-44ea-9c4d-02907dca81e7','not specified'),('webAuthnPolicyRequireResidentKey','dcc080c5-aede-4fd3-8b01-bd0928b730a2','not specified'),('webAuthnPolicyRequireResidentKeyPasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6','not specified'),('webAuthnPolicyRequireResidentKeyPasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7','not specified'),('webAuthnPolicyRequireResidentKeyPasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2','not specified'),('webAuthnPolicyRpEntityName','40ae881c-f4e4-4b07-b097-a67d2bf515e6','keycloak'),('webAuthnPolicyRpEntityName','4327ba47-4116-44ea-9c4d-02907dca81e7','keycloak'),('webAuthnPolicyRpEntityName','dcc080c5-aede-4fd3-8b01-bd0928b730a2','keycloak'),('webAuthnPolicyRpEntityNamePasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6','keycloak'),('webAuthnPolicyRpEntityNamePasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7','keycloak'),('webAuthnPolicyRpEntityNamePasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2','keycloak'),('webAuthnPolicyRpId','40ae881c-f4e4-4b07-b097-a67d2bf515e6',''),('webAuthnPolicyRpId','4327ba47-4116-44ea-9c4d-02907dca81e7',''),('webAuthnPolicyRpId','dcc080c5-aede-4fd3-8b01-bd0928b730a2',''),('webAuthnPolicyRpIdPasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6',''),('webAuthnPolicyRpIdPasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7',''),('webAuthnPolicyRpIdPasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2',''),('webAuthnPolicySignatureAlgorithms','40ae881c-f4e4-4b07-b097-a67d2bf515e6','ES256'),('webAuthnPolicySignatureAlgorithms','4327ba47-4116-44ea-9c4d-02907dca81e7','ES256'),('webAuthnPolicySignatureAlgorithms','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ES256'),('webAuthnPolicySignatureAlgorithmsPasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6','ES256'),('webAuthnPolicySignatureAlgorithmsPasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7','ES256'),('webAuthnPolicySignatureAlgorithmsPasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2','ES256'),('webAuthnPolicyUserVerificationRequirement','40ae881c-f4e4-4b07-b097-a67d2bf515e6','not specified'),('webAuthnPolicyUserVerificationRequirement','4327ba47-4116-44ea-9c4d-02907dca81e7','not specified'),('webAuthnPolicyUserVerificationRequirement','dcc080c5-aede-4fd3-8b01-bd0928b730a2','not specified'),('webAuthnPolicyUserVerificationRequirementPasswordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6','not specified'),('webAuthnPolicyUserVerificationRequirementPasswordless','4327ba47-4116-44ea-9c4d-02907dca81e7','not specified'),('webAuthnPolicyUserVerificationRequirementPasswordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2','not specified');
/*!40000 ALTER TABLE `REALM_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM_DEFAULT_GROUPS`
--

DROP TABLE IF EXISTS `REALM_DEFAULT_GROUPS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_DEFAULT_GROUPS` (
  `REALM_ID` varchar(36) NOT NULL,
  `GROUP_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`REALM_ID`,`GROUP_ID`),
  UNIQUE KEY `CON_GROUP_ID_DEF_GROUPS` (`GROUP_ID`),
  KEY `IDX_REALM_DEF_GRP_REALM` (`REALM_ID`),
  CONSTRAINT `FK_DEF_GROUPS_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM_DEFAULT_GROUPS`
--

LOCK TABLES `REALM_DEFAULT_GROUPS` WRITE;
/*!40000 ALTER TABLE `REALM_DEFAULT_GROUPS` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_DEFAULT_GROUPS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM_ENABLED_EVENT_TYPES`
--

DROP TABLE IF EXISTS `REALM_ENABLED_EVENT_TYPES`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_ENABLED_EVENT_TYPES` (
  `REALM_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) NOT NULL,
  PRIMARY KEY (`REALM_ID`,`VALUE`),
  KEY `IDX_REALM_EVT_TYPES_REALM` (`REALM_ID`),
  CONSTRAINT `FK_H846O4H0W8EPX5NWEDRF5Y69J` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM_ENABLED_EVENT_TYPES`
--

LOCK TABLES `REALM_ENABLED_EVENT_TYPES` WRITE;
/*!40000 ALTER TABLE `REALM_ENABLED_EVENT_TYPES` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_ENABLED_EVENT_TYPES` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM_EVENTS_LISTENERS`
--

DROP TABLE IF EXISTS `REALM_EVENTS_LISTENERS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_EVENTS_LISTENERS` (
  `REALM_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) NOT NULL,
  PRIMARY KEY (`REALM_ID`,`VALUE`),
  KEY `IDX_REALM_EVT_LIST_REALM` (`REALM_ID`),
  CONSTRAINT `FK_H846O4H0W8EPX5NXEV9F5Y69J` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM_EVENTS_LISTENERS`
--

LOCK TABLES `REALM_EVENTS_LISTENERS` WRITE;
/*!40000 ALTER TABLE `REALM_EVENTS_LISTENERS` DISABLE KEYS */;
INSERT INTO `REALM_EVENTS_LISTENERS` (`REALM_ID`, `VALUE`) VALUES ('40ae881c-f4e4-4b07-b097-a67d2bf515e6','jboss-logging'),('4327ba47-4116-44ea-9c4d-02907dca81e7','jboss-logging'),('dcc080c5-aede-4fd3-8b01-bd0928b730a2','jboss-logging');
/*!40000 ALTER TABLE `REALM_EVENTS_LISTENERS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM_LOCALIZATIONS`
--

DROP TABLE IF EXISTS `REALM_LOCALIZATIONS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_LOCALIZATIONS` (
  `REALM_ID` varchar(255) NOT NULL,
  `LOCALE` varchar(255) NOT NULL,
  `TEXTS` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`REALM_ID`,`LOCALE`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM_LOCALIZATIONS`
--

LOCK TABLES `REALM_LOCALIZATIONS` WRITE;
/*!40000 ALTER TABLE `REALM_LOCALIZATIONS` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_LOCALIZATIONS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM_REQUIRED_CREDENTIAL`
--

DROP TABLE IF EXISTS `REALM_REQUIRED_CREDENTIAL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_REQUIRED_CREDENTIAL` (
  `TYPE` varchar(255) NOT NULL,
  `FORM_LABEL` varchar(255) DEFAULT NULL,
  `INPUT` tinyint NOT NULL DEFAULT '0',
  `SECRET` tinyint NOT NULL DEFAULT '0',
  `REALM_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`REALM_ID`,`TYPE`),
  CONSTRAINT `FK_5HG65LYBEVAVKQFKI3KPONH9V` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM_REQUIRED_CREDENTIAL`
--

LOCK TABLES `REALM_REQUIRED_CREDENTIAL` WRITE;
/*!40000 ALTER TABLE `REALM_REQUIRED_CREDENTIAL` DISABLE KEYS */;
INSERT INTO `REALM_REQUIRED_CREDENTIAL` (`TYPE`, `FORM_LABEL`, `INPUT`, `SECRET`, `REALM_ID`) VALUES ('password','password',1,1,'40ae881c-f4e4-4b07-b097-a67d2bf515e6'),('password','password',1,1,'4327ba47-4116-44ea-9c4d-02907dca81e7'),('password','password',1,1,'dcc080c5-aede-4fd3-8b01-bd0928b730a2');
/*!40000 ALTER TABLE `REALM_REQUIRED_CREDENTIAL` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM_SMTP_CONFIG`
--

DROP TABLE IF EXISTS `REALM_SMTP_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_SMTP_CONFIG` (
  `REALM_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) DEFAULT NULL,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`REALM_ID`,`NAME`),
  CONSTRAINT `FK_70EJ8XDXGXD0B9HH6180IRR0O` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM_SMTP_CONFIG`
--

LOCK TABLES `REALM_SMTP_CONFIG` WRITE;
/*!40000 ALTER TABLE `REALM_SMTP_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_SMTP_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REALM_SUPPORTED_LOCALES`
--

DROP TABLE IF EXISTS `REALM_SUPPORTED_LOCALES`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_SUPPORTED_LOCALES` (
  `REALM_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) NOT NULL,
  PRIMARY KEY (`REALM_ID`,`VALUE`),
  KEY `IDX_REALM_SUPP_LOCAL_REALM` (`REALM_ID`),
  CONSTRAINT `FK_SUPPORTED_LOCALES_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REALM_SUPPORTED_LOCALES`
--

LOCK TABLES `REALM_SUPPORTED_LOCALES` WRITE;
/*!40000 ALTER TABLE `REALM_SUPPORTED_LOCALES` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_SUPPORTED_LOCALES` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REDIRECT_URIS`
--

DROP TABLE IF EXISTS `REDIRECT_URIS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REDIRECT_URIS` (
  `CLIENT_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`VALUE`),
  KEY `IDX_REDIR_URI_CLIENT` (`CLIENT_ID`),
  CONSTRAINT `FK_1BURS8PB4OUJ97H5WUPPAHV9F` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REDIRECT_URIS`
--

LOCK TABLES `REDIRECT_URIS` WRITE;
/*!40000 ALTER TABLE `REDIRECT_URIS` DISABLE KEYS */;
INSERT INTO `REDIRECT_URIS` (`CLIENT_ID`, `VALUE`) VALUES ('05021082-bfbc-4ec6-872a-e6c0916922c1','/admin/master/console/*'),('1947f168-b049-4b78-8031-afcf98eae08d','/realms/test-realm/account/*'),('1a0f720d-df73-4056-926d-10d520c81992','/realms/test-realm/account/*'),('320d776e-9e8e-4263-abf6-ebe3d862d547','/realms/TEST/account/*'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','/admin/TEST/console/*'),('41fe1d7e-3717-4965-8c77-e9868b3d98d8','/realms/TEST/account/*'),('657f6bd0-aa09-4703-99f7-e3f48e2de466','/realms/master/account/*'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','http://localhost:3000/*'),('95e70707-5b66-41b0-a127-590f786b2fba','https://example.com/*'),('b1c791ca-b006-4f60-af37-20225d5876f5','/realms/master/account/*'),('c12ebc78-3392-401a-8328-5dbb4cddf222','/admin/test-realm/console/*'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','https://example.com/*');
/*!40000 ALTER TABLE `REDIRECT_URIS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REQUIRED_ACTION_CONFIG`
--

DROP TABLE IF EXISTS `REQUIRED_ACTION_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REQUIRED_ACTION_CONFIG` (
  `REQUIRED_ACTION_ID` varchar(36) NOT NULL,
  `VALUE` longtext,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`REQUIRED_ACTION_ID`,`NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REQUIRED_ACTION_CONFIG`
--

LOCK TABLES `REQUIRED_ACTION_CONFIG` WRITE;
/*!40000 ALTER TABLE `REQUIRED_ACTION_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `REQUIRED_ACTION_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `REQUIRED_ACTION_PROVIDER`
--

DROP TABLE IF EXISTS `REQUIRED_ACTION_PROVIDER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REQUIRED_ACTION_PROVIDER` (
  `ID` varchar(36) NOT NULL,
  `ALIAS` varchar(255) DEFAULT NULL,
  `NAME` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `DEFAULT_ACTION` tinyint NOT NULL DEFAULT '0',
  `PROVIDER_ID` varchar(255) DEFAULT NULL,
  `PRIORITY` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_REQ_ACT_PROV_REALM` (`REALM_ID`),
  CONSTRAINT `FK_REQ_ACT_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `REQUIRED_ACTION_PROVIDER`
--

LOCK TABLES `REQUIRED_ACTION_PROVIDER` WRITE;
/*!40000 ALTER TABLE `REQUIRED_ACTION_PROVIDER` DISABLE KEYS */;
INSERT INTO `REQUIRED_ACTION_PROVIDER` (`ID`, `ALIAS`, `NAME`, `REALM_ID`, `ENABLED`, `DEFAULT_ACTION`, `PROVIDER_ID`, `PRIORITY`) VALUES ('029baed7-8eec-48d7-aaeb-9d832ae8758d','UPDATE_PASSWORD','Update Password','dcc080c5-aede-4fd3-8b01-bd0928b730a2',1,0,'UPDATE_PASSWORD',30),('03fee2c0-5ba4-4117-955b-5c6ce6cc46a4','update_user_locale','Update User Locale','dcc080c5-aede-4fd3-8b01-bd0928b730a2',1,0,'update_user_locale',1000),('0661ade8-f103-42f4-b3bd-1221ee1dca28','VERIFY_PROFILE','Verify Profile','4327ba47-4116-44ea-9c4d-02907dca81e7',1,0,'VERIFY_PROFILE',90),('24abb218-f8ff-4965-95ef-9bb5bf63446c','VERIFY_EMAIL','Verify Email','4327ba47-4116-44ea-9c4d-02907dca81e7',1,0,'VERIFY_EMAIL',50),('29ccbe75-f59e-4981-ba8a-63768124115b','webauthn-register','Webauthn Register','dcc080c5-aede-4fd3-8b01-bd0928b730a2',1,0,'webauthn-register',70),('2f7a2ab3-e910-4663-b6a7-0678dce74d09','CONFIGURE_TOTP','Configure OTP','dcc080c5-aede-4fd3-8b01-bd0928b730a2',1,0,'CONFIGURE_TOTP',10),('338cbac9-e636-4258-8277-e8a6062b323b','webauthn-register','Webauthn Register','40ae881c-f4e4-4b07-b097-a67d2bf515e6',1,0,'webauthn-register',70),('3879c642-c9cb-4424-9fa0-b1db481020ba','webauthn-register-passwordless','Webauthn Register Passwordless','40ae881c-f4e4-4b07-b097-a67d2bf515e6',1,0,'webauthn-register-passwordless',80),('42750b17-883d-4c22-b758-3db4530f4731','delete_account','Delete Account','4327ba47-4116-44ea-9c4d-02907dca81e7',0,0,'delete_account',60),('583f9cc2-35b8-4358-b76c-4a67fae53429','update_user_locale','Update User Locale','40ae881c-f4e4-4b07-b097-a67d2bf515e6',1,0,'update_user_locale',1000),('604b4757-2dfb-420d-9478-71e8960fd0be','UPDATE_PROFILE','Update Profile','dcc080c5-aede-4fd3-8b01-bd0928b730a2',1,0,'UPDATE_PROFILE',40),('795eccf1-08b8-4cd8-b0af-0167fa29d3c5','VERIFY_PROFILE','Verify Profile','40ae881c-f4e4-4b07-b097-a67d2bf515e6',1,0,'VERIFY_PROFILE',90),('7a69ad9e-ab64-4864-87f3-c72222c59e67','update_user_locale','Update User Locale','4327ba47-4116-44ea-9c4d-02907dca81e7',1,0,'update_user_locale',1000),('7c0ddad7-45cc-42fb-a8a1-df1cd8e955b9','TERMS_AND_CONDITIONS','Terms and Conditions','4327ba47-4116-44ea-9c4d-02907dca81e7',0,0,'TERMS_AND_CONDITIONS',20),('851dc502-226a-4817-bc7d-5f6728423a99','TERMS_AND_CONDITIONS','Terms and Conditions','40ae881c-f4e4-4b07-b097-a67d2bf515e6',0,0,'TERMS_AND_CONDITIONS',20),('8b114f95-5486-4727-88f5-c6081a018121','UPDATE_PROFILE','Update Profile','40ae881c-f4e4-4b07-b097-a67d2bf515e6',1,0,'UPDATE_PROFILE',40),('957c4f7c-8bd1-410f-8804-3a1e267fcb2c','delete_account','Delete Account','40ae881c-f4e4-4b07-b097-a67d2bf515e6',0,0,'delete_account',60),('985a032c-732f-426d-a357-ce4a33f61405','webauthn-register','Webauthn Register','4327ba47-4116-44ea-9c4d-02907dca81e7',1,0,'webauthn-register',70),('a4d9ffcc-569b-4f12-9750-d989ebf3bb8c','webauthn-register-passwordless','Webauthn Register Passwordless','4327ba47-4116-44ea-9c4d-02907dca81e7',1,0,'webauthn-register-passwordless',80),('a5f5746b-da4c-4306-b352-57122328616e','TERMS_AND_CONDITIONS','Terms and Conditions','dcc080c5-aede-4fd3-8b01-bd0928b730a2',0,0,'TERMS_AND_CONDITIONS',20),('ac839830-3c79-42f3-9c0c-d725e0362fdb','CONFIGURE_TOTP','Configure OTP','40ae881c-f4e4-4b07-b097-a67d2bf515e6',1,0,'CONFIGURE_TOTP',10),('b4af0ff2-51e4-4e90-bd93-87f72735893e','webauthn-register-passwordless','Webauthn Register Passwordless','dcc080c5-aede-4fd3-8b01-bd0928b730a2',1,0,'webauthn-register-passwordless',80),('b70add62-3e4b-4a89-a5d1-464d8f11c68c','VERIFY_EMAIL','Verify Email','dcc080c5-aede-4fd3-8b01-bd0928b730a2',1,0,'VERIFY_EMAIL',50),('cb4d128a-8a31-46be-b4da-9e08013795ef','VERIFY_EMAIL','Verify Email','40ae881c-f4e4-4b07-b097-a67d2bf515e6',1,0,'VERIFY_EMAIL',50),('dee6bc1b-34c3-4605-9bc1-7140e3e5a750','UPDATE_PROFILE','Update Profile','4327ba47-4116-44ea-9c4d-02907dca81e7',1,0,'UPDATE_PROFILE',40),('e1653f77-40ce-47e9-8fbc-ee268ca7f721','UPDATE_PASSWORD','Update Password','4327ba47-4116-44ea-9c4d-02907dca81e7',1,0,'UPDATE_PASSWORD',30),('ecc065e7-c4c8-428b-9d4c-11a89c8ef80a','CONFIGURE_TOTP','Configure OTP','4327ba47-4116-44ea-9c4d-02907dca81e7',1,0,'CONFIGURE_TOTP',10),('eff33713-2ff1-4306-965c-4013dc426d63','VERIFY_PROFILE','Verify Profile','dcc080c5-aede-4fd3-8b01-bd0928b730a2',1,0,'VERIFY_PROFILE',90),('f4cfca02-2f64-49e3-a801-58fc147e1ed1','UPDATE_PASSWORD','Update Password','40ae881c-f4e4-4b07-b097-a67d2bf515e6',1,0,'UPDATE_PASSWORD',30),('fbc14221-b2d8-45bf-986c-f04418a45109','delete_account','Delete Account','dcc080c5-aede-4fd3-8b01-bd0928b730a2',0,0,'delete_account',60);
/*!40000 ALTER TABLE `REQUIRED_ACTION_PROVIDER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_ATTRIBUTE`
--

DROP TABLE IF EXISTS `RESOURCE_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_ATTRIBUTE` (
  `ID` varchar(36) NOT NULL DEFAULT 'sybase-needs-something-here',
  `NAME` varchar(255) NOT NULL,
  `VALUE` varchar(255) DEFAULT NULL,
  `RESOURCE_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK_5HRM2VLF9QL5FU022KQEPOVBR` (`RESOURCE_ID`),
  CONSTRAINT `FK_5HRM2VLF9QL5FU022KQEPOVBR` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_ATTRIBUTE`
--

LOCK TABLES `RESOURCE_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `RESOURCE_ATTRIBUTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_POLICY`
--

DROP TABLE IF EXISTS `RESOURCE_POLICY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_POLICY` (
  `RESOURCE_ID` varchar(36) NOT NULL,
  `POLICY_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`RESOURCE_ID`,`POLICY_ID`),
  KEY `IDX_RES_POLICY_POLICY` (`POLICY_ID`),
  CONSTRAINT `FK_FRSRPOS53XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`),
  CONSTRAINT `FK_FRSRPP213XCX4WNKOG82SSRFY` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_POLICY`
--

LOCK TABLES `RESOURCE_POLICY` WRITE;
/*!40000 ALTER TABLE `RESOURCE_POLICY` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_POLICY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_SCOPE`
--

DROP TABLE IF EXISTS `RESOURCE_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SCOPE` (
  `RESOURCE_ID` varchar(36) NOT NULL,
  `SCOPE_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`RESOURCE_ID`,`SCOPE_ID`),
  KEY `IDX_RES_SCOPE_SCOPE` (`SCOPE_ID`),
  CONSTRAINT `FK_FRSRPOS13XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`),
  CONSTRAINT `FK_FRSRPS213XCX4WNKOG82SSRFY` FOREIGN KEY (`SCOPE_ID`) REFERENCES `RESOURCE_SERVER_SCOPE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_SCOPE`
--

LOCK TABLES `RESOURCE_SCOPE` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SCOPE` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_SERVER`
--

DROP TABLE IF EXISTS `RESOURCE_SERVER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER` (
  `ID` varchar(36) NOT NULL,
  `ALLOW_RS_REMOTE_MGMT` tinyint NOT NULL DEFAULT '0',
  `POLICY_ENFORCE_MODE` tinyint DEFAULT NULL,
  `DECISION_STRATEGY` tinyint NOT NULL DEFAULT '1',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_SERVER`
--

LOCK TABLES `RESOURCE_SERVER` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_SERVER_PERM_TICKET`
--

DROP TABLE IF EXISTS `RESOURCE_SERVER_PERM_TICKET`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER_PERM_TICKET` (
  `ID` varchar(36) NOT NULL,
  `OWNER` varchar(255) DEFAULT NULL,
  `REQUESTER` varchar(255) DEFAULT NULL,
  `CREATED_TIMESTAMP` bigint NOT NULL,
  `GRANTED_TIMESTAMP` bigint DEFAULT NULL,
  `RESOURCE_ID` varchar(36) NOT NULL,
  `SCOPE_ID` varchar(36) DEFAULT NULL,
  `RESOURCE_SERVER_ID` varchar(36) NOT NULL,
  `POLICY_ID` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_FRSR6T700S9V50BU18WS5PMT` (`OWNER`,`REQUESTER`,`RESOURCE_SERVER_ID`,`RESOURCE_ID`,`SCOPE_ID`),
  KEY `FK_FRSRHO213XCX4WNKOG82SSPMT` (`RESOURCE_SERVER_ID`),
  KEY `FK_FRSRHO213XCX4WNKOG83SSPMT` (`RESOURCE_ID`),
  KEY `FK_FRSRHO213XCX4WNKOG84SSPMT` (`SCOPE_ID`),
  KEY `FK_FRSRPO2128CX4WNKOG82SSRFY` (`POLICY_ID`),
  CONSTRAINT `FK_FRSRHO213XCX4WNKOG82SSPMT` FOREIGN KEY (`RESOURCE_SERVER_ID`) REFERENCES `RESOURCE_SERVER` (`ID`),
  CONSTRAINT `FK_FRSRHO213XCX4WNKOG83SSPMT` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`),
  CONSTRAINT `FK_FRSRHO213XCX4WNKOG84SSPMT` FOREIGN KEY (`SCOPE_ID`) REFERENCES `RESOURCE_SERVER_SCOPE` (`ID`),
  CONSTRAINT `FK_FRSRPO2128CX4WNKOG82SSRFY` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_SERVER_PERM_TICKET`
--

LOCK TABLES `RESOURCE_SERVER_PERM_TICKET` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER_PERM_TICKET` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER_PERM_TICKET` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_SERVER_POLICY`
--

DROP TABLE IF EXISTS `RESOURCE_SERVER_POLICY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER_POLICY` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `TYPE` varchar(255) NOT NULL,
  `DECISION_STRATEGY` tinyint DEFAULT NULL,
  `LOGIC` tinyint DEFAULT NULL,
  `RESOURCE_SERVER_ID` varchar(36) DEFAULT NULL,
  `OWNER` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_FRSRPT700S9V50BU18WS5HA6` (`NAME`,`RESOURCE_SERVER_ID`),
  KEY `IDX_RES_SERV_POL_RES_SERV` (`RESOURCE_SERVER_ID`),
  CONSTRAINT `FK_FRSRPO213XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_SERVER_ID`) REFERENCES `RESOURCE_SERVER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_SERVER_POLICY`
--

LOCK TABLES `RESOURCE_SERVER_POLICY` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER_POLICY` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER_POLICY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_SERVER_RESOURCE`
--

DROP TABLE IF EXISTS `RESOURCE_SERVER_RESOURCE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER_RESOURCE` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `TYPE` varchar(255) DEFAULT NULL,
  `ICON_URI` varchar(255) DEFAULT NULL,
  `OWNER` varchar(255) DEFAULT NULL,
  `RESOURCE_SERVER_ID` varchar(36) DEFAULT NULL,
  `OWNER_MANAGED_ACCESS` tinyint NOT NULL DEFAULT '0',
  `DISPLAY_NAME` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_FRSR6T700S9V50BU18WS5HA6` (`NAME`,`OWNER`,`RESOURCE_SERVER_ID`),
  KEY `IDX_RES_SRV_RES_RES_SRV` (`RESOURCE_SERVER_ID`),
  CONSTRAINT `FK_FRSRHO213XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_SERVER_ID`) REFERENCES `RESOURCE_SERVER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_SERVER_RESOURCE`
--

LOCK TABLES `RESOURCE_SERVER_RESOURCE` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER_RESOURCE` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER_RESOURCE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_SERVER_SCOPE`
--

DROP TABLE IF EXISTS `RESOURCE_SERVER_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER_SCOPE` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `ICON_URI` varchar(255) DEFAULT NULL,
  `RESOURCE_SERVER_ID` varchar(36) DEFAULT NULL,
  `DISPLAY_NAME` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_FRSRST700S9V50BU18WS5HA6` (`NAME`,`RESOURCE_SERVER_ID`),
  KEY `IDX_RES_SRV_SCOPE_RES_SRV` (`RESOURCE_SERVER_ID`),
  CONSTRAINT `FK_FRSRSO213XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_SERVER_ID`) REFERENCES `RESOURCE_SERVER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_SERVER_SCOPE`
--

LOCK TABLES `RESOURCE_SERVER_SCOPE` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER_SCOPE` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RESOURCE_URIS`
--

DROP TABLE IF EXISTS `RESOURCE_URIS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_URIS` (
  `RESOURCE_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) NOT NULL,
  PRIMARY KEY (`RESOURCE_ID`,`VALUE`),
  CONSTRAINT `FK_RESOURCE_SERVER_URIS` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RESOURCE_URIS`
--

LOCK TABLES `RESOURCE_URIS` WRITE;
/*!40000 ALTER TABLE `RESOURCE_URIS` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_URIS` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ROLE_ATTRIBUTE`
--

DROP TABLE IF EXISTS `ROLE_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ROLE_ATTRIBUTE` (
  `ID` varchar(36) NOT NULL,
  `ROLE_ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `VALUE` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_ROLE_ATTRIBUTE` (`ROLE_ID`),
  CONSTRAINT `FK_ROLE_ATTRIBUTE_ID` FOREIGN KEY (`ROLE_ID`) REFERENCES `KEYCLOAK_ROLE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ROLE_ATTRIBUTE`
--

LOCK TABLES `ROLE_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `ROLE_ATTRIBUTE` DISABLE KEYS */;
INSERT INTO `ROLE_ATTRIBUTE` (`ID`, `ROLE_ID`, `NAME`, `VALUE`) VALUES ('5954db91-0680-49d2-b6b5-8225e73d66f7','8a32586f-9583-40d7-9cd2-a548f8220c80','map_edit','true'),('6ec3d1e2-2995-4510-b5fe-e985f632c3d8','8a32586f-9583-40d7-9cd2-a548f8220c80','map_view','true');
/*!40000 ALTER TABLE `ROLE_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `SCOPE_MAPPING`
--

DROP TABLE IF EXISTS `SCOPE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `SCOPE_MAPPING` (
  `CLIENT_ID` varchar(36) NOT NULL,
  `ROLE_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`ROLE_ID`),
  KEY `IDX_SCOPE_MAPPING_ROLE` (`ROLE_ID`),
  CONSTRAINT `FK_OUSE064PLMLR732LXJCN1Q5F1` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `SCOPE_MAPPING`
--

LOCK TABLES `SCOPE_MAPPING` WRITE;
/*!40000 ALTER TABLE `SCOPE_MAPPING` DISABLE KEYS */;
INSERT INTO `SCOPE_MAPPING` (`CLIENT_ID`, `ROLE_ID`) VALUES ('1a0f720d-df73-4056-926d-10d520c81992','084bf744-b03c-4316-be84-5593d9af0691'),('b1c791ca-b006-4f60-af37-20225d5876f5','1d41c4d8-3377-488e-8563-f93dfd895d39'),('320d776e-9e8e-4263-abf6-ebe3d862d547','8c549c79-801f-4d45-892d-22fe10111eb2'),('b1c791ca-b006-4f60-af37-20225d5876f5','8e5f68f7-bcb4-43a5-ae98-6dfc29e18395'),('1a0f720d-df73-4056-926d-10d520c81992','91e47ad5-5b4d-4440-b982-3dd69c10b2fb'),('320d776e-9e8e-4263-abf6-ebe3d862d547','e081abd2-551d-469d-bbd7-e522600639f5');
/*!40000 ALTER TABLE `SCOPE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `SCOPE_POLICY`
--

DROP TABLE IF EXISTS `SCOPE_POLICY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `SCOPE_POLICY` (
  `SCOPE_ID` varchar(36) NOT NULL,
  `POLICY_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`SCOPE_ID`,`POLICY_ID`),
  KEY `IDX_SCOPE_POLICY_POLICY` (`POLICY_ID`),
  CONSTRAINT `FK_FRSRASP13XCX4WNKOG82SSRFY` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`),
  CONSTRAINT `FK_FRSRPASS3XCX4WNKOG82SSRFY` FOREIGN KEY (`SCOPE_ID`) REFERENCES `RESOURCE_SERVER_SCOPE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `SCOPE_POLICY`
--

LOCK TABLES `SCOPE_POLICY` WRITE;
/*!40000 ALTER TABLE `SCOPE_POLICY` DISABLE KEYS */;
/*!40000 ALTER TABLE `SCOPE_POLICY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USERNAME_LOGIN_FAILURE`
--

DROP TABLE IF EXISTS `USERNAME_LOGIN_FAILURE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USERNAME_LOGIN_FAILURE` (
  `REALM_ID` varchar(36) NOT NULL,
  `USERNAME` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `FAILED_LOGIN_NOT_BEFORE` int DEFAULT NULL,
  `LAST_FAILURE` bigint DEFAULT NULL,
  `LAST_IP_FAILURE` varchar(255) DEFAULT NULL,
  `NUM_FAILURES` int DEFAULT NULL,
  PRIMARY KEY (`REALM_ID`,`USERNAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USERNAME_LOGIN_FAILURE`
--

LOCK TABLES `USERNAME_LOGIN_FAILURE` WRITE;
/*!40000 ALTER TABLE `USERNAME_LOGIN_FAILURE` DISABLE KEYS */;
/*!40000 ALTER TABLE `USERNAME_LOGIN_FAILURE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_ATTRIBUTE`
--

DROP TABLE IF EXISTS `USER_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_ATTRIBUTE` (
  `NAME` varchar(255) NOT NULL,
  `VALUE` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `USER_ID` varchar(36) NOT NULL,
  `ID` varchar(36) NOT NULL DEFAULT 'sybase-needs-something-here',
  `LONG_VALUE_HASH` binary(64) DEFAULT NULL,
  `LONG_VALUE_HASH_LOWER_CASE` binary(64) DEFAULT NULL,
  `LONG_VALUE` longtext CHARACTER SET utf8 COLLATE utf8_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_USER_ATTRIBUTE` (`USER_ID`),
  KEY `IDX_USER_ATTRIBUTE_NAME` (`NAME`,`VALUE`),
  KEY `USER_ATTR_LONG_VALUES` (`LONG_VALUE_HASH`,`NAME`),
  KEY `USER_ATTR_LONG_VALUES_LOWER_CASE` (`LONG_VALUE_HASH_LOWER_CASE`,`NAME`),
  CONSTRAINT `FK_5HRM2VLF9QL5FU043KQEPOVBR` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_ATTRIBUTE`
--

LOCK TABLES `USER_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `USER_ATTRIBUTE` DISABLE KEYS */;
INSERT INTO `USER_ATTRIBUTE` (`NAME`, `VALUE`, `USER_ID`, `ID`, `LONG_VALUE_HASH`, `LONG_VALUE_HASH_LOWER_CASE`, `LONG_VALUE`) VALUES ('test_role','test role','8a36ba4d-3a1a-4536-9d99-5f21304735c4','0768d1c0-a612-4423-bb20-4d7530bfce96',NULL,NULL,NULL),('test_role','yeap','f9416272-2f08-49b3-a877-0c65127a2f03','c33b3e62-7e7d-4f4b-8056-8e802e2f4820',NULL,NULL,NULL),('test_role','test','ffce733f-1bac-4e35-be5a-00eba5a8483f','fe03066c-addd-4e39-bf31-8ecb348da807',NULL,NULL,NULL);
/*!40000 ALTER TABLE `USER_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_CONSENT`
--

DROP TABLE IF EXISTS `USER_CONSENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_CONSENT` (
  `ID` varchar(36) NOT NULL,
  `CLIENT_ID` varchar(255) DEFAULT NULL,
  `USER_ID` varchar(36) NOT NULL,
  `CREATED_DATE` bigint DEFAULT NULL,
  `LAST_UPDATED_DATE` bigint DEFAULT NULL,
  `CLIENT_STORAGE_PROVIDER` varchar(36) DEFAULT NULL,
  `EXTERNAL_CLIENT_ID` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_JKUWUVD56ONTGSUHOGM8UEWRT` (`CLIENT_ID`,`CLIENT_STORAGE_PROVIDER`,`EXTERNAL_CLIENT_ID`,`USER_ID`),
  KEY `IDX_USER_CONSENT` (`USER_ID`),
  CONSTRAINT `FK_GRNTCSNT_USER` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_CONSENT`
--

LOCK TABLES `USER_CONSENT` WRITE;
/*!40000 ALTER TABLE `USER_CONSENT` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_CONSENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_CONSENT_CLIENT_SCOPE`
--

DROP TABLE IF EXISTS `USER_CONSENT_CLIENT_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_CONSENT_CLIENT_SCOPE` (
  `USER_CONSENT_ID` varchar(36) NOT NULL,
  `SCOPE_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`USER_CONSENT_ID`,`SCOPE_ID`),
  KEY `IDX_USCONSENT_CLSCOPE` (`USER_CONSENT_ID`),
  CONSTRAINT `FK_GRNTCSNT_CLSC_USC` FOREIGN KEY (`USER_CONSENT_ID`) REFERENCES `USER_CONSENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_CONSENT_CLIENT_SCOPE`
--

LOCK TABLES `USER_CONSENT_CLIENT_SCOPE` WRITE;
/*!40000 ALTER TABLE `USER_CONSENT_CLIENT_SCOPE` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_CONSENT_CLIENT_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_ENTITY`
--

DROP TABLE IF EXISTS `USER_ENTITY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_ENTITY` (
  `ID` varchar(36) NOT NULL,
  `EMAIL` varchar(255) DEFAULT NULL,
  `EMAIL_CONSTRAINT` varchar(255) DEFAULT NULL,
  `EMAIL_VERIFIED` tinyint NOT NULL DEFAULT '0',
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `FEDERATION_LINK` varchar(255) DEFAULT NULL,
  `FIRST_NAME` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `LAST_NAME` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `REALM_ID` varchar(255) DEFAULT NULL,
  `USERNAME` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  `CREATED_TIMESTAMP` bigint DEFAULT NULL,
  `SERVICE_ACCOUNT_CLIENT_LINK` varchar(255) DEFAULT NULL,
  `NOT_BEFORE` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_DYKN684SL8UP1CRFEI6ECKHD7` (`REALM_ID`,`EMAIL_CONSTRAINT`),
  UNIQUE KEY `UK_RU8TT6T700S9V50BU18WS5HA6` (`REALM_ID`,`USERNAME`),
  KEY `IDX_USER_EMAIL` (`EMAIL`),
  KEY `IDX_USER_SERVICE_ACCOUNT` (`REALM_ID`,`SERVICE_ACCOUNT_CLIENT_LINK`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_ENTITY`
--

LOCK TABLES `USER_ENTITY` WRITE;
/*!40000 ALTER TABLE `USER_ENTITY` DISABLE KEYS */;
INSERT INTO `USER_ENTITY` (`ID`, `EMAIL`, `EMAIL_CONSTRAINT`, `EMAIL_VERIFIED`, `ENABLED`, `FEDERATION_LINK`, `FIRST_NAME`, `LAST_NAME`, `REALM_ID`, `USERNAME`, `CREATED_TIMESTAMP`, `SERVICE_ACCOUNT_CLIENT_LINK`, `NOT_BEFORE`) VALUES ('491fe75e-534f-48ae-ada1-36ea0ee54a1f','kenshin579@hotmail.com','kenshin579@hotmail.com',1,1,NULL,'Admin','User','4327ba47-4116-44ea-9c4d-02907dca81e7','admin',1713408445188,NULL,0);
/*!40000 ALTER TABLE `USER_ENTITY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_FEDERATION_CONFIG`
--

DROP TABLE IF EXISTS `USER_FEDERATION_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_FEDERATION_CONFIG` (
  `USER_FEDERATION_PROVIDER_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) DEFAULT NULL,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`USER_FEDERATION_PROVIDER_ID`,`NAME`),
  CONSTRAINT `FK_T13HPU1J94R2EBPEKR39X5EU5` FOREIGN KEY (`USER_FEDERATION_PROVIDER_ID`) REFERENCES `USER_FEDERATION_PROVIDER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_FEDERATION_CONFIG`
--

LOCK TABLES `USER_FEDERATION_CONFIG` WRITE;
/*!40000 ALTER TABLE `USER_FEDERATION_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_FEDERATION_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_FEDERATION_MAPPER`
--

DROP TABLE IF EXISTS `USER_FEDERATION_MAPPER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_FEDERATION_MAPPER` (
  `ID` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `FEDERATION_PROVIDER_ID` varchar(36) NOT NULL,
  `FEDERATION_MAPPER_TYPE` varchar(255) NOT NULL,
  `REALM_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_USR_FED_MAP_FED_PRV` (`FEDERATION_PROVIDER_ID`),
  KEY `IDX_USR_FED_MAP_REALM` (`REALM_ID`),
  CONSTRAINT `FK_FEDMAPPERPM_FEDPRV` FOREIGN KEY (`FEDERATION_PROVIDER_ID`) REFERENCES `USER_FEDERATION_PROVIDER` (`ID`),
  CONSTRAINT `FK_FEDMAPPERPM_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_FEDERATION_MAPPER`
--

LOCK TABLES `USER_FEDERATION_MAPPER` WRITE;
/*!40000 ALTER TABLE `USER_FEDERATION_MAPPER` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_FEDERATION_MAPPER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_FEDERATION_MAPPER_CONFIG`
--

DROP TABLE IF EXISTS `USER_FEDERATION_MAPPER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_FEDERATION_MAPPER_CONFIG` (
  `USER_FEDERATION_MAPPER_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) DEFAULT NULL,
  `NAME` varchar(255) NOT NULL,
  PRIMARY KEY (`USER_FEDERATION_MAPPER_ID`,`NAME`),
  CONSTRAINT `FK_FEDMAPPER_CFG` FOREIGN KEY (`USER_FEDERATION_MAPPER_ID`) REFERENCES `USER_FEDERATION_MAPPER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_FEDERATION_MAPPER_CONFIG`
--

LOCK TABLES `USER_FEDERATION_MAPPER_CONFIG` WRITE;
/*!40000 ALTER TABLE `USER_FEDERATION_MAPPER_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_FEDERATION_MAPPER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_FEDERATION_PROVIDER`
--

DROP TABLE IF EXISTS `USER_FEDERATION_PROVIDER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_FEDERATION_PROVIDER` (
  `ID` varchar(36) NOT NULL,
  `CHANGED_SYNC_PERIOD` int DEFAULT NULL,
  `DISPLAY_NAME` varchar(255) DEFAULT NULL,
  `FULL_SYNC_PERIOD` int DEFAULT NULL,
  `LAST_SYNC` int DEFAULT NULL,
  `PRIORITY` int DEFAULT NULL,
  `PROVIDER_NAME` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_USR_FED_PRV_REALM` (`REALM_ID`),
  CONSTRAINT `FK_1FJ32F6PTOLW2QY60CD8N01E8` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_FEDERATION_PROVIDER`
--

LOCK TABLES `USER_FEDERATION_PROVIDER` WRITE;
/*!40000 ALTER TABLE `USER_FEDERATION_PROVIDER` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_FEDERATION_PROVIDER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_GROUP_MEMBERSHIP`
--

DROP TABLE IF EXISTS `USER_GROUP_MEMBERSHIP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_GROUP_MEMBERSHIP` (
  `GROUP_ID` varchar(36) NOT NULL,
  `USER_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`GROUP_ID`,`USER_ID`),
  KEY `IDX_USER_GROUP_MAPPING` (`USER_ID`),
  CONSTRAINT `FK_USER_GROUP_USER` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_GROUP_MEMBERSHIP`
--

LOCK TABLES `USER_GROUP_MEMBERSHIP` WRITE;
/*!40000 ALTER TABLE `USER_GROUP_MEMBERSHIP` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_GROUP_MEMBERSHIP` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_REQUIRED_ACTION`
--

DROP TABLE IF EXISTS `USER_REQUIRED_ACTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_REQUIRED_ACTION` (
  `USER_ID` varchar(36) NOT NULL,
  `REQUIRED_ACTION` varchar(255) NOT NULL DEFAULT ' ',
  PRIMARY KEY (`REQUIRED_ACTION`,`USER_ID`),
  KEY `IDX_USER_REQACTIONS` (`USER_ID`),
  CONSTRAINT `FK_6QJ3W1JW9CVAFHE19BWSIUVMD` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_REQUIRED_ACTION`
--

LOCK TABLES `USER_REQUIRED_ACTION` WRITE;
/*!40000 ALTER TABLE `USER_REQUIRED_ACTION` DISABLE KEYS */;
INSERT INTO `USER_REQUIRED_ACTION` (`USER_ID`, `REQUIRED_ACTION`) VALUES ('0d368ccc-e005-4bfd-b0ce-ac8beb1bb289','UPDATE_PASSWORD');
/*!40000 ALTER TABLE `USER_REQUIRED_ACTION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_ROLE_MAPPING`
--

DROP TABLE IF EXISTS `USER_ROLE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_ROLE_MAPPING` (
  `ROLE_ID` varchar(255) NOT NULL,
  `USER_ID` varchar(36) NOT NULL,
  PRIMARY KEY (`ROLE_ID`,`USER_ID`),
  KEY `IDX_USER_ROLE_MAPPING` (`USER_ID`),
  CONSTRAINT `FK_C4FQV34P1MBYLLOXANG7B1Q3L` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_ROLE_MAPPING`
--

LOCK TABLES `USER_ROLE_MAPPING` WRITE;
/*!40000 ALTER TABLE `USER_ROLE_MAPPING` DISABLE KEYS */;
INSERT INTO `USER_ROLE_MAPPING` (`ROLE_ID`, `USER_ID`) VALUES ('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','0d368ccc-e005-4bfd-b0ce-ac8beb1bb289'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','242d893e-ff24-41c2-abac-e50df93cb864'),('06ad34dc-2164-4067-896b-933e9dadb08c','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('07178d4e-bdd9-412f-924b-62977dc1d23b','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('07755138-7efe-46df-9ae5-4f7869d5c570','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('0edbe58b-31c3-4503-81e6-3426b88107d7','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('0f5abc1b-e9de-4f60-94fa-7a9f537eb90a','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('10114928-9ec5-4ce4-baf0-696dd614ea04','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('1b146798-0c31-4fe4-8d8d-8341e5b463ec','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('241a9760-03f0-41fa-911c-0a023532161e','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('24e2a886-cc15-44c9-a176-1b12c038f88d','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('2e299e39-7ae8-4394-82af-0488b8ba850c','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('39e1e0ae-01d2-40c5-82b7-66727d7542bb','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('43087049-6e35-422e-a082-97d32212cfcf','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('48eeadab-d1c8-403f-ae89-2270450ae4f6','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('57e92e5d-9b6b-4782-9045-d5ac3b1cd017','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('581d6c1d-f6d2-4d8c-bfcc-56d8e5f8de3d','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('623ac08d-4ff5-49c7-9392-c6f8a49ef26b','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('6552d7cc-3d28-424a-b9ab-c0719ad9d1f1','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('76ab719c-a862-441e-a66b-329547590e1d','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('84ec0238-dc71-4066-b99e-49e59cc6bc7a','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('888d64fb-b137-4f1e-b530-b4aef00ab967','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('8a2220b6-356e-4d7e-995f-83ca8da2e9fe','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('8b6e8ddd-250e-4bac-8bdd-de3a9d4ab9e6','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('94526ddc-9e84-4c0b-958b-07184a18ab32','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('a2fa72c7-f176-4d55-aab3-9c6e18a6c3e3','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('a6381b94-b388-4ac0-94a4-a9c24384f9ba','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('b717dbd5-35d8-4865-94c2-bf876c0146f4','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('c8ee3a69-7a2d-48f8-9373-6677e3724523','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('ce73d6bd-9432-4068-a9e6-02f000f8b4d3','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('d5f7c5c8-33c4-483c-802e-6ed2c1bc2f29','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('d69ea56d-5246-4e89-a45f-4effdf6d9bcf','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('db24697d-b20d-4a2a-99b6-f41d26e73b7e','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('e293b015-eada-4709-8d6c-da9c98a5ec04','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('e7b5b42b-09a2-46bf-b6b5-87f71892c9dd','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('ea082079-5ed5-4688-8d02-1a0a2dbed4fa','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('ee7eb7c2-62cf-4526-9560-a1e1c4bd27f2','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('f9e0e15e-14f9-4327-b781-04bc3078317e','491fe75e-534f-48ae-ada1-36ea0ee54a1f'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','5aeb1754-527c-4348-9c73-fe216819b4b0'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','6dbe094f-d8d0-4e8b-acc4-3b7aced25d2d'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','8a36ba4d-3a1a-4536-9d99-5f21304735c4'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','8f52ffc9-35e3-4000-9e76-d0ac17e56104'),('8a32586f-9583-40d7-9cd2-a548f8220c80','a80e2adf-05a7-4caf-9411-82efd36d23db'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','a80e2adf-05a7-4caf-9411-82efd36d23db'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','abc5c932-3904-46f8-8203-bf1943181ef0'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','b5b04553-1091-4272-9b93-91eb4e115d74'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','b6abd49e-7ea7-4e74-9081-8083ca72b15f'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','b6bd7af0-d9bc-442d-bef6-4224770b4b5a'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','dc00724e-6a75-42f8-ba4e-9cadc0fe516a'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','e38e0455-3afc-43df-945a-c1bd2d5274de'),('14fce4d6-46b4-4c14-bd14-4ac7227ea3b0','e90eb5ff-fd2e-404c-93bb-e5a79dd750bb'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','f9416272-2f08-49b3-a877-0c65127a2f03'),('c8496af8-af9d-49ea-9bb2-26ba4bcbef2f','ffce733f-1bac-4e35-be5a-00eba5a8483f');
/*!40000 ALTER TABLE `USER_ROLE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_SESSION`
--

DROP TABLE IF EXISTS `USER_SESSION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_SESSION` (
  `ID` varchar(36) NOT NULL,
  `AUTH_METHOD` varchar(255) DEFAULT NULL,
  `IP_ADDRESS` varchar(255) DEFAULT NULL,
  `LAST_SESSION_REFRESH` int DEFAULT NULL,
  `LOGIN_USERNAME` varchar(255) DEFAULT NULL,
  `REALM_ID` varchar(255) DEFAULT NULL,
  `REMEMBER_ME` tinyint NOT NULL DEFAULT '0',
  `STARTED` int DEFAULT NULL,
  `USER_ID` varchar(255) DEFAULT NULL,
  `USER_SESSION_STATE` int DEFAULT NULL,
  `BROKER_SESSION_ID` varchar(255) DEFAULT NULL,
  `BROKER_USER_ID` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_SESSION`
--

LOCK TABLES `USER_SESSION` WRITE;
/*!40000 ALTER TABLE `USER_SESSION` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_SESSION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `USER_SESSION_NOTE`
--

DROP TABLE IF EXISTS `USER_SESSION_NOTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_SESSION_NOTE` (
  `USER_SESSION` varchar(36) NOT NULL,
  `NAME` varchar(255) NOT NULL,
  `VALUE` text,
  PRIMARY KEY (`USER_SESSION`,`NAME`),
  CONSTRAINT `FK5EDFB00FF51D3472` FOREIGN KEY (`USER_SESSION`) REFERENCES `USER_SESSION` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `USER_SESSION_NOTE`
--

LOCK TABLES `USER_SESSION_NOTE` WRITE;
/*!40000 ALTER TABLE `USER_SESSION_NOTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_SESSION_NOTE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `WEB_ORIGINS`
--

DROP TABLE IF EXISTS `WEB_ORIGINS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `WEB_ORIGINS` (
  `CLIENT_ID` varchar(36) NOT NULL,
  `VALUE` varchar(255) NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`VALUE`),
  KEY `IDX_WEB_ORIG_CLIENT` (`CLIENT_ID`),
  CONSTRAINT `FK_LOJPHO213XCX4WNKOG82SSRFY` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `WEB_ORIGINS`
--

LOCK TABLES `WEB_ORIGINS` WRITE;
/*!40000 ALTER TABLE `WEB_ORIGINS` DISABLE KEYS */;
INSERT INTO `WEB_ORIGINS` (`CLIENT_ID`, `VALUE`) VALUES ('05021082-bfbc-4ec6-872a-e6c0916922c1','+'),('33ce64f8-6ffd-430d-a260-c5c8f6d92308','+'),('71cf0db7-6cfe-49a2-84e0-f58fd821b17f','*'),('95e70707-5b66-41b0-a127-590f786b2fba','*'),('c12ebc78-3392-401a-8328-5dbb4cddf222','+'),('ecabf8ab-e548-4b7a-a158-4d3e774afd77','*');
/*!40000 ALTER TABLE `WEB_ORIGINS` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-08-11 16:10:15
