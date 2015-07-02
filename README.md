# Passport


数据库

	CREATE DATABASE `passport` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

	CREATE TABLE `user` (
  	       `id` bigint(64) NOT NULL,
 	       `phone` varchar(11) COLLATE utf8_bin DEFAULT NULL,
  	       `email` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  	       `password` varchar(45) COLLATE utf8_bin NOT NULL,
  	       `add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  	       `update_time` datetime NOT NULL,
  	       `stat` int(11) NOT NULL DEFAULT '0',
  	       `version` int(11) DEFAULT NULL,
  	       PRIMARY KEY (`id`),
  	       UNIQUE KEY `phone_unique` (`phone`),
  	       UNIQUE KEY `email_unique` (`email`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

