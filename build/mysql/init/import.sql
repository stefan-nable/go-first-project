CREATE DATABASE IF NOT EXISTS `log`;
USE `log`;

CREATE TABLE IF NOT EXISTS `log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `timestamp` datetime NOT NULL,
  `message` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
);