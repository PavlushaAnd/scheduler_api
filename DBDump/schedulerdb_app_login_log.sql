-- MySQL dump 10.13  Distrib 8.0.36, for Win64 (x86_64)
--
-- Host: localhost    Database: schedulerdb
-- ------------------------------------------------------
-- Server version	8.0.36

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `app_login_log`
--

DROP TABLE IF EXISTS `app_login_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `app_login_log` (
  `id` int NOT NULL AUTO_INCREMENT,
  `login_user` varchar(255) NOT NULL DEFAULT '',
  `login_ip` varchar(255) NOT NULL DEFAULT '',
  `login_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_login_log`
--

LOCK TABLES `app_login_log` WRITE;
/*!40000 ALTER TABLE `app_login_log` DISABLE KEYS */;
INSERT INTO `app_login_log` VALUES (1,'brunton','127.0.0.1','2024-04-11 08:34:56'),(2,'brunton','127.0.0.1','2024-04-11 08:40:29'),(3,'brunton','127.0.0.1','2024-04-11 08:45:53'),(4,'brunton','127.0.0.1','2024-04-11 08:46:18'),(5,'brunton','127.0.0.1','2024-04-11 08:47:02'),(6,'brunton','127.0.0.1','2024-04-11 08:53:10'),(7,'brunton','127.0.0.1','2024-04-11 08:55:34'),(8,'brunton','127.0.0.1','2024-04-11 09:00:51'),(9,'brunton','127.0.0.1','2024-04-11 09:02:41'),(10,'brunton','127.0.0.1','2024-04-11 09:05:28'),(11,'brunton','127.0.0.1','2024-04-11 09:38:23'),(12,'brunton','127.0.0.1','2024-04-11 09:56:52'),(13,'brunton','127.0.0.1','2024-04-11 13:19:29'),(14,'brunton','127.0.0.1','2024-04-11 13:19:56'),(15,'brunton','127.0.0.1','2024-04-11 14:48:58'),(16,'brunton','127.0.0.1','2024-04-11 16:13:26'),(17,'brunton','127.0.0.1','2024-04-11 17:05:01'),(18,'','127.0.0.1','2024-04-15 13:18:44'),(19,'','127.0.0.1','2024-04-15 13:19:02'),(20,'','127.0.0.1','2024-04-15 13:19:25'),(21,'','127.0.0.1','2024-04-15 13:21:44'),(22,'brunton','127.0.0.1','2024-04-15 13:22:30'),(23,'brunton','127.0.0.1','2024-04-15 13:22:57'),(24,'brunton','127.0.0.1','2024-04-15 15:24:03'),(25,'','127.0.0.1','2024-04-15 17:04:56'),(26,'brunton','127.0.0.1','2024-04-15 17:05:18'),(27,'brunton','127.0.0.1','2024-04-17 08:20:38'),(28,'brunton','127.0.0.1','2024-04-17 10:47:00'),(29,'dev','127.0.0.1','2024-04-17 12:53:08'),(30,'brunton','127.0.0.1','2024-04-17 12:56:49'),(31,'brunton','127.0.0.1','2024-04-17 12:59:50'),(32,'brunton','127.0.0.1','2024-04-17 13:06:15'),(33,'dev','127.0.0.1','2024-04-17 14:18:11'),(34,'brunton','127.0.0.1','2024-04-17 14:53:14'),(35,'dev','127.0.0.1','2024-04-17 15:33:35'),(36,'brunton','127.0.0.1','2024-04-17 17:02:10'),(37,'brunton','127.0.0.1','2024-04-18 08:41:22'),(38,'dev','127.0.0.1','2024-04-18 09:55:00'),(39,'dev','127.0.0.1','2024-04-18 11:46:54'),(40,'brunton','127.0.0.1','2024-04-18 12:56:45'),(41,'dev','127.0.0.1','2024-04-18 13:00:35'),(42,'dev','127.0.0.1','2024-04-18 15:00:03'),(43,'dev','127.0.0.1','2024-04-18 16:28:07'),(44,'brunton','127.0.0.1','2024-04-22 13:11:05'),(45,'brunton','127.0.0.1','2024-04-23 11:27:09');
/*!40000 ALTER TABLE `app_login_log` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-04-23 17:17:47
