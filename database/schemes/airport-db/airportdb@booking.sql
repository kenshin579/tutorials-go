-- MySQLShell dump 1.0.2  Distrib Ver 8.0.25 for Linux on x86_64 - for MySQL 8.0.25 (MySQL Community Server (GPL)), for Linux (x86_64)
--
-- Host: 10.0.1.81    Database: airportdb    Table: booking
-- ------------------------------------------------------
-- Server version	8.0.26

--
-- Table structure for table `booking`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `booking` (
  `booking_id` int NOT NULL AUTO_INCREMENT,
  `flight_id` int NOT NULL,
  `seat` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `passenger_id` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  PRIMARY KEY (`booking_id`),
  UNIQUE KEY `seatplan_unq` (`flight_id`,`seat`),
  KEY `flight_idx` (`flight_id`),
  KEY `passenger_idx` (`passenger_id`),
  CONSTRAINT `booking_ibfk_1` FOREIGN KEY (`flight_id`) REFERENCES `flight` (`flight_id`),
  CONSTRAINT `booking_ibfk_2` FOREIGN KEY (`passenger_id`) REFERENCES `passenger` (`passenger_id`)
) ENGINE=InnoDB AUTO_INCREMENT=55099799 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Flughafen DB by Stefan Pröll, Eva Zangerle, Wolfgang Gassler is licensed under CC BY 4.0. To view a copy of this license, visit https://creativecommons.org/licenses/by/4.0';
/*!40101 SET character_set_client = @saved_cs_client */;
