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
-- Table structure for table `task`
--

DROP TABLE IF EXISTS `task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_code` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL DEFAULT '',
  `user_code` varchar(255) NOT NULL DEFAULT '',
  `room_name` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) DEFAULT NULL,
  `location` varchar(255) NOT NULL DEFAULT '',
  `repeatable` varchar(255) NOT NULL DEFAULT '',
  `start_date` datetime NOT NULL,
  `end_date` datetime NOT NULL,
  `rec_end_date` datetime DEFAULT NULL,
  `rec_start_date` datetime DEFAULT NULL,
  `version` int NOT NULL DEFAULT '0',
  `last_modified` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  `creator_code` varchar(255) NOT NULL DEFAULT '',
  `editor_code` varchar(255) NOT NULL DEFAULT '',
  `project_name` varchar(255) NOT NULL DEFAULT '',
  `client_code` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task`
--

LOCK TABLES `task` WRITE;
/*!40000 ALTER TABLE `task` DISABLE KEYS */;
INSERT INTO `task` VALUES (51,'task_1713412736042508500','Weekly meeting','brunton','office','All eployees should prepare presentation with their week report','','FREQ=MONTHLY','2022-03-22 09:00:00','2022-03-22 13:00:00','2022-05-22 08:30:00','2022-03-22 09:00:00',0,'2024-04-18 15:58:56','2024-04-18 15:58:56','dev','','BNZL_Scheduler',''),(52,'task_1713412736042508501','Weekly meeting','brunton','office','All eployees should prepare presentation with their week report','','FREQ=MONTHLY','2022-04-22 09:00:00','2022-04-22 13:00:00','2022-05-22 08:30:00','2022-03-22 09:00:00',0,'2024-04-18 15:58:56','2024-04-18 15:58:56','dev','dev','BNZL_Scheduler',''),(53,'task_1713755771077406300','Weekly meeting','','','Project Planning','','FREQ=WEEKLY','2024-04-22 15:00:00','2024-04-22 15:30:00','2024-05-31 15:15:11','2024-04-22 15:14:11',1,'2024-04-23 16:37:12','2024-04-22 15:16:11','brunton','brunton','',''),(54,'task_1713755771077406301','Weekly meeting','','','Project Planning','','FREQ=WEEKLY','2024-04-29 15:14:11','2024-04-29 15:15:11','2024-05-31 15:15:11','2024-04-22 15:14:11',0,'2024-04-22 15:16:11','2024-04-22 15:16:11','brunton','brunton','',''),(55,'task_1713755771077406302','Weekly meeting','','','Project Planning','','FREQ=WEEKLY','2024-05-06 15:14:11','2024-05-06 15:15:11','2024-05-31 15:15:11','2024-04-22 15:14:11',0,'2024-04-22 15:16:11','2024-04-22 15:16:11','brunton','brunton','',''),(56,'task_1713755771077406303','Weekly meeting','','','Project Planning','','FREQ=WEEKLY','2024-05-13 15:14:11','2024-05-13 15:15:11','2024-05-31 15:15:11','2024-04-22 15:14:11',0,'2024-04-22 15:16:11','2024-04-22 15:16:11','brunton','brunton','',''),(57,'task_1713755771077406304','Weekly meeting','','','Project Planning','','FREQ=WEEKLY','2024-05-20 15:14:11','2024-05-20 15:15:11','2024-05-31 15:15:11','2024-04-22 15:14:11',0,'2024-04-22 15:16:11','2024-04-22 15:16:11','brunton','brunton','',''),(58,'task_1713755771077406305','Weekly meeting','','','Project Planning','','FREQ=WEEKLY','2024-05-27 15:14:11','2024-05-27 15:15:11','2024-05-31 15:15:11','2024-04-22 15:14:11',0,'2024-04-22 15:16:11','2024-04-22 15:16:11','brunton','brunton','',''),(59,'task_1713755895786180800','Scheduler App','','','','','','2024-04-20 15:16:00','2024-04-20 15:17:00','2024-04-22 15:17:05','2024-04-19 15:16:05',4,'2024-04-22 15:19:12','2024-04-22 15:18:16','brunton','brunton','',''),(60,'task_1713847137066851000','test','','','test','','FREQ=DAILY','2024-04-22 15:00:00','2024-04-22 16:00:00','2024-04-23 16:37:41','2024-04-22 15:00:00',0,'2024-04-23 16:38:57','2024-04-23 16:38:57','brunton','','',''),(62,'task_1713847256622411300','test2','','','','','','2024-04-22 14:00:00','2024-04-22 16:00:00','2024-04-23 16:40:32','2024-04-22 16:39:31',1,'2024-04-23 16:41:27','2024-04-23 16:40:57','brunton','brunton','','');
/*!40000 ALTER TABLE `task` ENABLE KEYS */;
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
