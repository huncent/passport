# Passport


数据库

	CREATE DATABASE `passport` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

	CREATE TABLE `passport`.`users` (
  		`id` int(11) NOT NULL AUTO_INCREMENT,
  		`phone` varchar(11) COLLATE utf8_bin DEFAULT NULL,
  		`email` varchar(45) COLLATE utf8_bin DEFAULT NULL,
  		`passworld` varchar(45) COLLATE utf8_bin NOT NULL,
  		`add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  		`update_time` datetime NOT NULL,
  		`stat` int(11) NOT NULL DEFAULT '1',
  		PRIMARY KEY (`id`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
