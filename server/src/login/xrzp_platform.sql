/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50721
Source Host           : localhost:3306
Source Database       : xrzp_platform

Target Server Type    : MYSQL
Target Server Version : 50721
File Encoding         : 65001

Date: 2018-03-30 10:52:24
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `logs_accountlogin`
-- ----------------------------
DROP TABLE IF EXISTS `logs_accountlogin`;
CREATE TABLE `logs_accountlogin` (
  `uid` int(11) DEFAULT NULL,
  `logindate` datetime DEFAULT NULL,
  `ip` varchar(64) DEFAULT NULL,
  `device` varchar(256) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of logs_accountlogin
-- ----------------------------

-- ----------------------------
-- Table structure for `logs_payment`
-- ----------------------------
DROP TABLE IF EXISTS `logs_payment`;
CREATE TABLE `logs_payment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `channelid` int(11) NOT NULL,
  `serverid` int(11) NOT NULL,
  `username` varchar(100) NOT NULL DEFAULT '',
  `orderid` varchar(100) NOT NULL,
  `amount` int(18) NOT NULL,
  `uid` int(11) DEFAULT NULL,
  `roleid` int(11) NOT NULL,
  `paydate` int(11) DEFAULT NULL,
  `status` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_order` (`channelid`,`orderid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- ----------------------------
-- Records of logs_payment
-- ----------------------------

-- ----------------------------
-- Table structure for `logs_rolelogin`
-- ----------------------------
DROP TABLE IF EXISTS `logs_rolelogin`;
CREATE TABLE `logs_rolelogin` (
  `uid` int(11) DEFAULT NULL,
  `roleid` int(11) DEFAULT NULL,
  `serverid` int(11) DEFAULT NULL,
  `ip` varchar(64) DEFAULT NULL,
  `logindate` datetime DEFAULT NULL,
  `device` varchar(256) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of logs_rolelogin
-- ----------------------------

-- ----------------------------
-- Table structure for `san_account`
-- ----------------------------
DROP TABLE IF EXISTS `san_account`;
CREATE TABLE `san_account` (
  `uid` bigint(20) NOT NULL AUTO_INCREMENT,
  `account` varchar(45) NOT NULL,
  `password` varchar(45) NOT NULL,
  `creator` varchar(45) NOT NULL,
  `channelid` varchar(31) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`uid`),
  KEY `INDEX_PASSWORD` (`account`,`password`) USING BTREE,
  KEY `INDEX_ACCOUNT` (`account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of san_account
-- ----------------------------

-- ----------------------------
-- Table structure for `san_announce`
-- ----------------------------
DROP TABLE IF EXISTS `san_announce`;
CREATE TABLE `san_announce` (
  `id` int(11) NOT NULL,
  `channelid` int(11) DEFAULT NULL,
  `title` varchar(128) DEFAULT NULL,
  `content` varchar(2048) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of san_announce
-- ----------------------------

-- ----------------------------
-- Table structure for `san_channel`
-- ----------------------------
DROP TABLE IF EXISTS `san_channel`;
CREATE TABLE `san_channel` (
  `id` int(11) NOT NULL,
  `channelid` int(11) NOT NULL DEFAULT '0',
  `name` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_channel` (`channelid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- ----------------------------
-- Records of san_channel
-- ----------------------------

-- ----------------------------
-- Table structure for `san_server`
-- ----------------------------
DROP TABLE IF EXISTS `san_server`;
CREATE TABLE `san_server` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `serverid` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `host` varchar(50) NOT NULL,
  `state` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unque_serverId` (`serverid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- ----------------------------
-- Records of san_server
-- ----------------------------

-- ----------------------------
-- Table structure for `san_serverchannel`
-- ----------------------------
DROP TABLE IF EXISTS `san_serverchannel`;
CREATE TABLE `san_serverchannel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `serverid` int(11) NOT NULL DEFAULT '0',
  `servername` varchar(250) NOT NULL,
  `channelid` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_map` (`serverid`,`channelid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

-- ----------------------------
-- Records of san_serverchannel
-- ----------------------------

-- ----------------------------
-- View structure for `uccardroles`
-- ----------------------------
DROP VIEW IF EXISTS `uccardroles`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `uccardroles` AS select `faccount`.`FId` AS `Id`,`faccount`.`FChannelId` AS `ChannelId`,`faccount`.`FUserName` AS `AccountId`,`user_role`.`roleid` AS `RoleId`,`user_role`.`rolename` AS `RoleName`,`user_role`.`rolelvl` AS `RoleLevel`,`user_role`.`serverid` AS `ServerId`,`fserver`.`FName` AS `ServerName`,`user_role`.`lastlogin` AS `LastLogin` from ((`faccount` join `user_role` on((`faccount`.`FId` = `user_role`.`pt_uid`))) join `fserver` on((`user_role`.`serverid` = `fserver`.`FServerId`))) order by `user_role`.`lastlogin` desc ;

-- ----------------------------
-- Procedure structure for `creatUser`
-- ----------------------------
DROP PROCEDURE IF EXISTS `creatUser`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `creatUser`()
BEGIN
DECLARE ID INT DEFAULT 0;
DECLARE I INT DEFAULT 0;
DECLARE TEMP VARCHAR(20);
SET TEMP='';
WHILE I<500 DO
SET I=I+1;
SET TEMP=CONCAT('cs',I);
INSERT INTO faccount(FChannelId,FUserName,FLoginDate,FCreateDate,FCreateIp,FPhone) VALUES (1,TEMP,NOW(),NOW(),'127.0.0.1',0);
END WHILE;
END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for `sp_use_giftcard`
-- ----------------------------
DROP PROCEDURE IF EXISTS `sp_use_giftcard`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_use_giftcard`(_code VARCHAR(50),
   _pid int,
   _regid int)
begin

   DECLARE mtype int DEFAULT 0;
   DECLARE mfid int DEFAULT 0;
   DECLARE menddate DATETIME DEFAULT NOW();
	 SELECT a.fid,b.cardType,b.enddate 
		INTO mfid,mtype,menddate
		FROM giftcard_use a join giftcard b on a.cardId=b.fid 
    where a.`code`=_code and a.fisuse=FALSE;
   SELECT _code,_pid,_regid;
end
;;
DELIMITER ;

-- ----------------------------
-- Function structure for `DiffDays`
-- ----------------------------
DROP FUNCTION IF EXISTS `DiffDays`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` FUNCTION `DiffDays`(`d1` datetime,`d2` datetime) RETURNS int(11)
BEGIN
	RETURN DATEDIFF(d1, d2) ;
END
;;
DELIMITER ;

-- ----------------------------
-- Function structure for `DiffMonths`
-- ----------------------------
DROP FUNCTION IF EXISTS `DiffMonths`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` FUNCTION `DiffMonths`(`d1` datetime,`d2` datetime) RETURNS int(2)
BEGIN
 RETURN TIMESTAMPDIFF(MONTH, d2, d1);
END
;;
DELIMITER ;

-- ----------------------------
-- Function structure for `DiffSeconds`
-- ----------------------------
DROP FUNCTION IF EXISTS `DiffSeconds`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` FUNCTION `DiffSeconds`(`d1` datetime,`d2` datetime) RETURNS int(11)
BEGIN
	RETURN TIMESTAMPDIFF(SECOND, d2, d1);
END
;;
DELIMITER ;
