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
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_code` varchar(255) NOT NULL DEFAULT '',
  `user_name` varchar(255) NOT NULL DEFAULT '',
  `position_code` varchar(255) NOT NULL DEFAULT '',
  `inactive` tinyint(1) NOT NULL DEFAULT '0',
  `phone_no` varchar(255) NOT NULL DEFAULT '',
  `email_address` varchar(255) NOT NULL DEFAULT '',
  `has_uploaded_page` tinyint(1) NOT NULL DEFAULT '0',
  `has_recognised_page` tinyint(1) NOT NULL DEFAULT '0',
  `has_confirmed_page` tinyint(1) NOT NULL DEFAULT '0',
  `has_posted_page` tinyint(1) NOT NULL DEFAULT '0',
  `password` varchar(255) NOT NULL DEFAULT '',
  `role` varchar(255) NOT NULL DEFAULT '',
  `color_text` varchar(255) NOT NULL DEFAULT '',
  `color_background` varchar(255) NOT NULL DEFAULT '',
  `last_modified` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  `creator_code` varchar(255) NOT NULL DEFAULT '',
  `editor_code` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'dev','Developer',' ',0,'',' dev@gmail.com',0,0,0,0,'ff2dc05f8dddcebd288c80f75297454f','admin','#ffffff','#5a2626','2024-04-22 14:07:50','2024-04-10 15:34:55',' ',''),(2,'brunton','one','',0,'','b@gmail.com',0,0,0,0,'43960bae0e10d95717ba939201dbba61','admin','#ffffff','#62783b','2024-04-17 15:20:01','2024-04-11 11:39:49','admin',''),(3,'Mark','Mark','',0,'','string',0,0,0,0,'756472bf913d0fb6e8a57eb708b01e47','admin','#ffffff','#e100ff','2024-04-18 16:32:26','2024-04-11 11:44:44','admin','dev'),(4,'chris','Gang Lee','',0,'','c@gmail.com',0,0,0,0,'053c068fcd7ad669a03414bda21d3c34','Developer','#d6d6d6','#ff0000','2024-04-17 11:46:09','2024-04-15 13:24:39','admin',''),(5,'pa','Paul','',0,'','p@gmail.com',0,0,0,0,'0b6817c6aace8f5aab295a31196fc424','admin','#007bff','#ffa200','2024-04-15 13:40:09','2024-04-15 13:40:09','admin',''),(9,'asmie','Asmie Tmng','',0,'','asmie@gmail.com',0,0,0,0,'7bb015b06aef5685fb7d214d8c8a8915','admin','#ff0000','#75efff','2024-04-17 15:00:48','2024-04-15 16:56:35','admin',''),(10,'asmie','Asmie Tamang','',0,'','',0,0,0,0,'50509da28a05c97ffc07a072020fd53a','','#ffffff','#5900ff','2024-04-17 15:00:28','2024-04-15 16:57:22','admin',''),(12,'sb','Sue','',0,'','sue@gmail',0,0,0,0,'28192867d58d818a72d4e46c61ba39c2','admin','#ffffff','#004cff','2024-04-17 14:57:51','2024-04-15 17:07:00','admin',''),(13,'james','James','',1,'','j@gmail.com',0,0,0,0,'de0fedaeedc7cae7db12a757e1b835ee','admin','#ffffff','#00a34f','2024-04-22 14:09:28','2024-04-17 11:17:01','brunton',''),(14,'Mark22','MarkTmng','MNGR',1,'','mM@gmail.com',0,0,0,0,'cd16b08066e89987a3e4edccc9398337','admin','#ffffff','#000000','2024-04-22 14:08:46','2024-04-18 11:54:57','dev',''),(16,'brrr','Brun','',1,'','br@gmail.com',0,0,0,0,'8c3addd93f52b4e8e75e049121d8b125','developer','#ffffff','#ff0000','2024-04-18 16:32:52','2024-04-18 14:41:52','brunton',''),(17,'f','fu','',0,'','f@',0,0,0,0,'e1cafebf35b3c64e75a5f2d5f70b180c','user','#d9b4b4','#eb0000','2024-04-22 15:11:02','2024-04-22 15:07:32','brunton',''),(18,'test','test','',0,'','f@gmail.com',0,0,0,0,'f54fea2ad1eecc224e6aeca3dbadc3e0','user','#f2e3e3','#245011','2024-04-23 11:29:44','2024-04-22 15:09:36','brunton',''),(19,'test','fgh','',1,'','f@',0,0,0,0,'f54fea2ad1eecc224e6aeca3dbadc3e0','user','#f2e3e3','#368514','2024-04-22 15:11:24','2024-04-22 15:11:24','brunton',''),(20,'Tom','tom','',0,'','t@gmail.com',0,0,0,0,'ca9f4134b469948328e1644982484826','user','#e8a6a6','#420000','2024-04-23 16:46:07','2024-04-23 16:46:07','brunton','');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
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
