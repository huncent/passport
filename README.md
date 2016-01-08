# 用户中心(Passport)
把用户权限，会话管理提出一个单独的服务。

##接口
###添加用户
	POST /user/add
	Body:
	{
		"cellphone":"18510511015", 
		"email":"liuhengloveyou@gmail.com",
		"nickname":"恒恒",
		"password":"123456"
	}
	
	成功: 200 {"userid":"用户ID"}
	失败: ^200 {"message":"错误信息."}
	
###用户登录
	POST /user/login[?token=xxxxxx]
	Body:
	{
		"cellphone":"18510511015", 
		"email":"liuhengloveyou@gmail.com",
		"nickname":"恒恒",
		"password":"123456"
	}	
  
	成功: 200 {"userid":"用户ID", "token":"xxxxxx"}
	失败: ^200 {"message":"错误信息."}

##数据库

	CREATE DATABASE `passport` IF NOTEXISTS /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

	CREATE TABLE `passport`.`user` (
	  `id` bigint(64) NOT NULL,
	  `cellphone` varchar(11) COLLATE utf8_bin DEFAULT NULL,
	  `email` varchar(45) COLLATE utf8_bin DEFAULT NULL,
	  `nickname` varchar(45) CHARACTER SET utf8 DEFAULT NULL,
	  `password` varchar(45) COLLATE utf8_bin NOT NULL,
	  `add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  `update_time` datetime NOT NULL,
	  `stat` int(11) NOT NULL DEFAULT '0',
	  `version` int(11) DEFAULT NULL,
	  PRIMARY KEY (`id`),
	  UNIQUE KEY `phone_UNIQUE` (`cellphone`),
	  UNIQUE KEY `email_UNIQUE` (`email`),
	  UNIQUE KEY `nickname_UNIQUE` (`nickname`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
