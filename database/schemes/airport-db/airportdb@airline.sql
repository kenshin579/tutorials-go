-- MySQLShell dump 1.0.2  Distrib Ver 8.0.25 for Linux on x86_64 - for MySQL 8.0.25 (MySQL Community Server (GPL)), for Linux (x86_64)
--
-- Host: 10.0.1.81    Database: airportdb    Table: airline
-- ------------------------------------------------------
-- Server version	8.0.26

--
-- Table structure for table `airline`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `airline` (
  `airline_id` smallint NOT NULL AUTO_INCREMENT,
  `iata` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `airlinename` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `base_airport` smallint NOT NULL,
  PRIMARY KEY (`airline_id`),
  UNIQUE KEY `iata_unq` (`iata`),
  KEY `base_airport_idx` (`base_airport`),
  CONSTRAINT `airline_ibfk_1` FOREIGN KEY (`base_airport`) REFERENCES `airport` (`airport_id`)
) ENGINE=InnoDB AUTO_INCREMENT=114 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Flughafen DB by Stefan Pröll, Eva Zangerle, Wolfgang Gassler is licensed under CC BY 4.0. To view a copy of this license, visit https://creativecommons.org/licenses/by/4.0';
/*!40101 SET character_set_client = @saved_cs_client */;
