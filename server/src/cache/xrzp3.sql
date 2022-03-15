/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50553
Source Host           : localhost:3306
Source Database       : xrzp3

Target Server Type    : MYSQL
Target Server Version : 50553
File Encoding         : 65001

Date: 2017-10-13 22:06:54
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `san_account`
-- ----------------------------
DROP TABLE IF EXISTS `san_account`;
CREATE TABLE `san_account` (
`uid`  bigint(20) NOT NULL AUTO_INCREMENT ,
`account`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`password`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`creator`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`time`  bigint(20) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=7967

;

-- ----------------------------
-- Table structure for `san_accountaward`
-- ----------------------------
DROP TABLE IF EXISTS `san_accountaward`;
CREATE TABLE `san_accountaward` (
`uid`  bigint(20) NOT NULL ,
`account`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`password`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`creator`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`time`  bigint(20) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_activity`
-- ----------------------------
DROP TABLE IF EXISTS `san_activity`;
CREATE TABLE `san_activity` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`jj`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`month`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_activitymask`
-- ----------------------------
DROP TABLE IF EXISTS `san_activitymask`;
CREATE TABLE `san_activitymask` (
`id`  int(11) UNSIGNED ZEROFILL NOT NULL AUTO_INCREMENT ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`items`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`topfight`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`toplevel`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=9005

;

-- ----------------------------
-- Table structure for `san_affairs`
-- ----------------------------
DROP TABLE IF EXISTS `san_affairs`;
CREATE TABLE `san_affairs` (
`uid`  int(64) NOT NULL ,
`data`  text CHARACTER SET utf8 COLLATE utf8_general_ci NULL ,
PRIMARY KEY (`uid`),
UNIQUE INDEX `uid_UNIQUE` (`uid`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_arena`
-- ----------------------------
DROP TABLE IF EXISTS `san_arena`;
CREATE TABLE `san_arena` (
`rankid`  int(11) NOT NULL ,
`uid`  bigint(20) NOT NULL ,
`name`  varchar(48) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL COMMENT '积分' ,
`format`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`best`  int(11) NOT NULL ,
`award`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`rankid`),
UNIQUE INDEX `uid_UNIQUE` (`uid`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_armsarena1`
-- ----------------------------
DROP TABLE IF EXISTS `san_armsarena1`;
CREATE TABLE `san_armsarena1` (
`uid`  bigint(20) NOT NULL ,
`rankid`  int(11) NOT NULL ,
`name`  varchar(48) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`format`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`best`  int(11) NOT NULL ,
`award`  bigint(20) NOT NULL ,
`worship`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`topaward`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_armsarena2`
-- ----------------------------
DROP TABLE IF EXISTS `san_armsarena2`;
CREATE TABLE `san_armsarena2` (
`uid`  bigint(20) NOT NULL ,
`rankid`  int(11) NOT NULL ,
`name`  varchar(48) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`format`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`best`  int(11) NOT NULL ,
`award`  bigint(20) NOT NULL ,
`worship`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`topaward`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_armsarena3`
-- ----------------------------
DROP TABLE IF EXISTS `san_armsarena3`;
CREATE TABLE `san_armsarena3` (
`uid`  bigint(20) NOT NULL ,
`rankid`  int(11) NOT NULL ,
`name`  varchar(48) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`format`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`best`  int(11) NOT NULL ,
`award`  bigint(20) NOT NULL ,
`worship`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`topaward`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_bag`
-- ----------------------------
DROP TABLE IF EXISTS `san_bag`;
CREATE TABLE `san_bag` (
`uid`  int(20) NOT NULL COMMENT 'key' ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`),
INDEX `uid` (`uid`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_battlearray`
-- ----------------------------
DROP TABLE IF EXISTS `san_battlearray`;
CREATE TABLE `san_battlearray` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`heroatt`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`solideratt`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_belog`
-- ----------------------------
DROP TABLE IF EXISTS `san_belog`;
CREATE TABLE `san_belog` (
`id`  bigint(20) UNSIGNED ZEROFILL NOT NULL AUTO_INCREMENT ,
`time`  bigint(20) NOT NULL ,
`type`  int(11) NOT NULL ,
`value`  int(11) NOT NULL ,
`param1`  int(11) NOT NULL ,
`param2`  int(11) NOT NULL ,
`uid`  bigint(20) NOT NULL ,
`dec`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`cur`  int(11) NOT NULL ,
`param3`  int(11) NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=MyISAM
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=5438

;

-- ----------------------------
-- Table structure for `san_book`
-- ----------------------------
DROP TABLE IF EXISTS `san_book`;
CREATE TABLE `san_book` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`bookcd`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_campdiplomacy`
-- ----------------------------
DROP TABLE IF EXISTS `san_campdiplomacy`;
CREATE TABLE `san_campdiplomacy` (
`id`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_casern`
-- ----------------------------
DROP TABLE IF EXISTS `san_casern`;
CREATE TABLE `san_casern` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`award`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_city`
-- ----------------------------
DROP TABLE IF EXISTS `san_city`;
CREATE TABLE `san_city` (
`id`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`fight`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`build`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_cityex`
-- ----------------------------
DROP TABLE IF EXISTS `san_cityex`;
CREATE TABLE `san_cityex` (
`id`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_country`
-- ----------------------------
DROP TABLE IF EXISTS `san_country`;
CREATE TABLE `san_country` (
`id`  int(11) NOT NULL ,
`level`  int(11) NOT NULL ,
`exp`  int(11) NOT NULL ,
`big`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`praise`  int(11) NOT NULL ,
`step`  int(11) NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_dartinfo`
-- ----------------------------
DROP TABLE IF EXISTS `san_dartinfo`;
CREATE TABLE `san_dartinfo` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`dartid`  varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'key' ,
`beginsection`  int(11) NOT NULL ,
`endsection`  int(11) NOT NULL ,
`begintime`  int(11) NOT NULL ,
`lasttime`  int(11) NOT NULL ,
`duan`  int(11) NOT NULL ,
`definfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`item`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`isfight`  int(11) NOT NULL DEFAULT 0 ,
`fightbegintime`  int(11) NOT NULL DEFAULT 0 ,
`fightuid`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`dartid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_diplomacy`
-- ----------------------------
DROP TABLE IF EXISTS `san_diplomacy`;
CREATE TABLE `san_diplomacy` (
`id`  int(11) NOT NULL ,
`progress`  int(11) NULL DEFAULT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_diplomacytask`
-- ----------------------------
DROP TABLE IF EXISTS `san_diplomacytask`;
CREATE TABLE `san_diplomacytask` (
`uid`  bigint(20) NOT NULL ,
`campver`  bigint(20) NOT NULL ,
`campaward`  int(11) NOT NULL ,
`camppersonal`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`campfightsignup`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_expedition`
-- ----------------------------
DROP TABLE IF EXISTS `san_expedition`;
CREATE TABLE `san_expedition` (
`uid`  int(11) NOT NULL ,
`acthero`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`defhero`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`deffight`  int(11) NOT NULL DEFAULT 0 ,
`score`  int(11) NOT NULL DEFAULT 0 ,
`allscore`  int(11) NOT NULL DEFAULT 0 ,
`bestdayscore`  int(11) NOT NULL DEFAULT 0 ,
`bestdayscoreaward`  int(11) NOT NULL DEFAULT 0 ,
`basefight`  int(11) NOT NULL DEFAULT 0 ,
`bestfight`  int(11) NOT NULL DEFAULT 0 ,
`curindex`  int(11) NOT NULL DEFAULT 0 ,
`star`  int(11) NOT NULL DEFAULT 0 ,
`myheroinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`mybuffinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`fightinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`boxinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`buffinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
`isrobot`  int(10) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_expedition1`
-- ----------------------------
DROP TABLE IF EXISTS `san_expedition1`;
CREATE TABLE `san_expedition1` (
`uid`  int(11) NOT NULL ,
`curindex`  int(11) NOT NULL DEFAULT 0 ,
`teaminfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`myheroinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`mybuffinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`boxinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`buffinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`fightinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`resetcount`  int(11) NOT NULL DEFAULT 1 ,
`isfinish`  int(11) NOT NULL DEFAULT 0 ,
`flowinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=MyISAM
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_fate`
-- ----------------------------
DROP TABLE IF EXISTS `san_fate`;
CREATE TABLE `san_fate` (
`id`  int(11) NOT NULL ,
`data`  text CHARACTER SET utf8 COLLATE utf8_general_ci NULL ,
PRIMARY KEY (`id`),
UNIQUE INDEX `id_UNIQUE` (`id`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_fief`
-- ----------------------------
DROP TABLE IF EXISTS `san_fief`;
CREATE TABLE `san_fief` (
`uid`  int(10) UNSIGNED NOT NULL ,
`data`  text CHARACTER SET utf8 COLLATE utf8_general_ci NULL ,
PRIMARY KEY (`uid`),
UNIQUE INDEX `uid_UNIQUE` (`uid`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_friend`
-- ----------------------------
DROP TABLE IF EXISTS `san_friend`;
CREATE TABLE `san_friend` (
`uid`  bigint(20) NOT NULL ,
`friend`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`apply`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`black`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`commend`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`hasapply`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`count`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_heroequ2`
-- ----------------------------
DROP TABLE IF EXISTS `san_heroequ2`;
CREATE TABLE `san_heroequ2` (
`id`  int(11) NOT NULL AUTO_INCREMENT ,
`uid`  int(11) NOT NULL COMMENT 'key' ,
`heroid`  int(11) NOT NULL COMMENT 'key' ,
`equid`  int(11) NOT NULL ,
`level`  int(11) NOT NULL DEFAULT 0 ,
`stars`  int(11) NOT NULL DEFAULT 0 ,
`color`  int(11) NOT NULL DEFAULT 1 ,
`open`  int(11) NOT NULL DEFAULT 0 ,
`type`  int(11) NOT NULL DEFAULT 0 COMMENT 'key' ,
`exp`  int(11) NOT NULL DEFAULT 0 ,
`spillexp`  int(11) NOT NULL DEFAULT 0 ,
`awakelevel`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`id`),
INDEX `uid` (`uid`) USING BTREE ,
INDEX `type` (`type`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=1

;

-- ----------------------------
-- Table structure for `san_herosupport`
-- ----------------------------
DROP TABLE IF EXISTS `san_herosupport`;
CREATE TABLE `san_herosupport` (
`heroid`  int(11) NOT NULL DEFAULT 0 COMMENT 'key' ,
`uid`  int(11) NOT NULL DEFAULT 1 COMMENT 'key' ,
`support`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`heroid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_industry`
-- ----------------------------
DROP TABLE IF EXISTS `san_industry`;
CREATE TABLE `san_industry` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_king1`
-- ----------------------------
DROP TABLE IF EXISTS `san_king1`;
CREATE TABLE `san_king1` (
`uid`  bigint(20) NOT NULL ,
`name`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`class`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`best`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`award`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`salary`  int(11) NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_king2`
-- ----------------------------
DROP TABLE IF EXISTS `san_king2`;
CREATE TABLE `san_king2` (
`uid`  bigint(20) NOT NULL ,
`name`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`class`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`best`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`award`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`salary`  int(11) NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_king3`
-- ----------------------------
DROP TABLE IF EXISTS `san_king3`;
CREATE TABLE `san_king3` (
`uid`  bigint(20) NOT NULL ,
`name`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`class`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`best`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`award`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`salary`  int(11) NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_king4`
-- ----------------------------
DROP TABLE IF EXISTS `san_king4`;
CREATE TABLE `san_king4` (
`id`  int(11) NOT NULL ,
`kinguid`  bigint(20) NOT NULL ,
`kingname`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`kingicon`  int(11) NOT NULL ,
`isfight`  int(11) NOT NULL ,
`openuid`  int(11) NOT NULL ,
`opentime`  bigint(20) NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_kingcity`
-- ----------------------------
DROP TABLE IF EXISTS `san_kingcity`;
CREATE TABLE `san_kingcity` (
`id`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`fight`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_linelog`
-- ----------------------------
DROP TABLE IF EXISTS `san_linelog`;
CREATE TABLE `san_linelog` (
`id`  bigint(20) NOT NULL AUTO_INCREMENT ,
`uid`  bigint(20) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`ip`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`line`  int(11) NOT NULL ,
`creator`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=MyISAM
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=1553

;

-- ----------------------------
-- Table structure for `san_log`
-- ----------------------------
DROP TABLE IF EXISTS `san_log`;
CREATE TABLE `san_log` (
`id`  bigint(20) NOT NULL AUTO_INCREMENT ,
`time`  bigint(20) NOT NULL ,
`type`  int(11) NOT NULL ,
`value`  int(11) NOT NULL ,
`param1`  int(11) NOT NULL ,
`param2`  int(11) NOT NULL ,
`param3`  int(11) NOT NULL DEFAULT 0 ,
`uid`  bigint(20) NOT NULL ,
`dec`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`cur`  bigint(11) NULL DEFAULT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=11711463

;

-- ----------------------------
-- Table structure for `san_mail`
-- ----------------------------
DROP TABLE IF EXISTS `san_mail`;
CREATE TABLE `san_mail` (
`id`  int(11) NOT NULL AUTO_INCREMENT ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=52

;

-- ----------------------------
-- Table structure for `san_money`
-- ----------------------------
DROP TABLE IF EXISTS `san_money`;
CREATE TABLE `san_money` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_monsteratfd`
-- ----------------------------
DROP TABLE IF EXISTS `san_monsteratfd`;
CREATE TABLE `san_monsteratfd` (
`uid`  int(11) NOT NULL ,
`energyuse`  int(11) NOT NULL DEFAULT 0 ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
`refreshcount`  int(11) NOT NULL DEFAULT 0 ,
`item`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`outmonster`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_notice`
-- ----------------------------
DROP TABLE IF EXISTS `san_notice`;
CREATE TABLE `san_notice` (
`id`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_office1`
-- ----------------------------
DROP TABLE IF EXISTS `san_office1`;
CREATE TABLE `san_office1` (
`uid`  bigint(20) NOT NULL ,
`name`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`class`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`best`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`award`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`salary`  int(11) NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_office2`
-- ----------------------------
DROP TABLE IF EXISTS `san_office2`;
CREATE TABLE `san_office2` (
`uid`  bigint(20) NOT NULL ,
`name`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`class`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`best`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`award`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`salary`  int(11) NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_office3`
-- ----------------------------
DROP TABLE IF EXISTS `san_office3`;
CREATE TABLE `san_office3` (
`uid`  bigint(20) NOT NULL ,
`name`  varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`class`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`num`  int(11) NOT NULL ,
`time`  bigint(20) NOT NULL ,
`best`  int(11) NOT NULL ,
`buynum`  int(11) NOT NULL ,
`award`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`point`  int(11) NOT NULL ,
`pointaward`  bigint(20) NOT NULL ,
`salary`  int(11) NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_pass`
-- ----------------------------
DROP TABLE IF EXISTS `san_pass`;
CREATE TABLE `san_pass` (
`uid`  bigint(20) NOT NULL ,
`warinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`missioninfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`boxinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`starboxinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`losehero`  int(11) NOT NULL ,
`passitem`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`passinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`jjinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`isfight`  int(11) NOT NULL ,
`totalstars`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_passguide`
-- ----------------------------
DROP TABLE IF EXISTS `san_passguide`;
CREATE TABLE `san_passguide` (
`id`  int(11) UNSIGNED ZEROFILL NOT NULL AUTO_INCREMENT ,
`firstteam`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`lastteam`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`fightteam`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=MyISAM
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=1229101

;

-- ----------------------------
-- Table structure for `san_science`
-- ----------------------------
DROP TABLE IF EXISTS `san_science`;
CREATE TABLE `san_science` (
`id`  int(11) NOT NULL ,
`num`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_sectiontax3`
-- ----------------------------
DROP TABLE IF EXISTS `san_sectiontax3`;
CREATE TABLE `san_sectiontax3` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`heroinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`buildinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`gocount`  int(10) NOT NULL DEFAULT 0 ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_shop1`
-- ----------------------------
DROP TABLE IF EXISTS `san_shop1`;
CREATE TABLE `san_shop1` (
`uid`  int(11) NOT NULL ,
`shoptype`  int(11) NOT NULL ,
`shopgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`refindex`  int(11) NOT NULL ,
`shopnextgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`todayrefcount`  int(11) NOT NULL ,
`sysreftime`  int(11) NOT NULL ,
`lastupdtime`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_shop2`
-- ----------------------------
DROP TABLE IF EXISTS `san_shop2`;
CREATE TABLE `san_shop2` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`shoptype`  int(11) NOT NULL COMMENT 'key' ,
`shopgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`refindex`  int(11) NOT NULL ,
`shopnextgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`todayrefcount`  int(11) NOT NULL ,
`sysreftime`  int(11) NOT NULL ,
`lastupdtime`  int(11) NOT NULL ,
PRIMARY KEY (`uid`),
INDEX `uid` (`uid`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_shop4`
-- ----------------------------
DROP TABLE IF EXISTS `san_shop4`;
CREATE TABLE `san_shop4` (
`uid`  int(11) NOT NULL ,
`shoptype`  int(11) NOT NULL ,
`shopgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`refindex`  int(11) NOT NULL ,
`shopnextgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`todayrefcount`  int(11) NOT NULL ,
`sysreftime`  int(11) NOT NULL ,
`lastupdtime`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_shop5`
-- ----------------------------
DROP TABLE IF EXISTS `san_shop5`;
CREATE TABLE `san_shop5` (
`uid`  int(11) NOT NULL ,
`shoptype`  int(11) NOT NULL ,
`shopgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`refindex`  int(11) NOT NULL ,
`shopnextgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`todayrefcount`  int(11) NOT NULL ,
`sysreftime`  int(11) NOT NULL ,
`lastupdtime`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_shop6`
-- ----------------------------
DROP TABLE IF EXISTS `san_shop6`;
CREATE TABLE `san_shop6` (
`uid`  int(11) NOT NULL ,
`shoptype`  int(11) NOT NULL ,
`shopgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`refindex`  int(11) NOT NULL ,
`shopnextgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`todayrefcount`  int(11) NOT NULL ,
`sysreftime`  int(11) NOT NULL ,
`lastupdtime`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_shop9`
-- ----------------------------
DROP TABLE IF EXISTS `san_shop9`;
CREATE TABLE `san_shop9` (
`uid`  int(11) NOT NULL ,
`shoptype`  int(11) NOT NULL ,
`shopgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`refindex`  int(11) NOT NULL ,
`shopnextgood`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`todayrefcount`  int(11) NOT NULL ,
`sysreftime`  int(11) NOT NULL ,
`lastupdtime`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_soldiergift`
-- ----------------------------
DROP TABLE IF EXISTS `san_soldiergift`;
CREATE TABLE `san_soldiergift` (
`uid`  int(11) NOT NULL ,
`soldiertype1`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`soldiertype2`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`soldiertype3`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_trainwar`
-- ----------------------------
DROP TABLE IF EXISTS `san_trainwar`;
CREATE TABLE `san_trainwar` (
`uid`  int(10) NOT NULL COMMENT 'key' ,
`item`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`lastupdtime`  int(10) NOT NULL DEFAULT 0 ,
`type1count`  int(5) NOT NULL DEFAULT 0 ,
`type1maxcount`  int(5) NOT NULL DEFAULT 0 ,
`beststar1`  int(10) NOT NULL DEFAULT 100000 ,
`acthero1`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type2count`  int(5) NOT NULL DEFAULT 0 ,
`type2maxcount`  int(5) NOT NULL DEFAULT 0 ,
`beststar2`  int(10) NOT NULL DEFAULT 100000 ,
`acthero2`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type3count`  int(5) NOT NULL DEFAULT 0 ,
`type3maxcount`  int(5) NOT NULL DEFAULT 0 ,
`beststar3`  int(10) NOT NULL DEFAULT 100000 ,
`acthero3`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type4count`  int(5) NOT NULL DEFAULT 0 ,
`type4maxcount`  int(5) NOT NULL DEFAULT 0 ,
`beststar4`  int(10) NOT NULL DEFAULT 100000 ,
`acthero4`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type5count`  int(5) NOT NULL DEFAULT 0 ,
`type5maxcount`  int(5) NOT NULL DEFAULT 0 ,
`beststar5`  int(10) NOT NULL DEFAULT 100000 ,
`acthero5`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type6count`  int(5) NOT NULL DEFAULT 0 ,
`type6maxcount`  int(5) NOT NULL DEFAULT 0 ,
`beststar6`  int(10) NOT NULL DEFAULT 100000 ,
`acthero6`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_unioninfo`
-- ----------------------------
DROP TABLE IF EXISTS `san_unioninfo`;
CREATE TABLE `san_unioninfo` (
`id`  int(11) NOT NULL AUTO_INCREMENT COMMENT 'key' ,
`icon`  int(11) NOT NULL ,
`unionname`  varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`masteruid`  int(11) NOT NULL DEFAULT 0 ,
`mastername`  varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`level`  int(11) NOT NULL DEFAULT 1 ,
`jointype`  int(11) NOT NULL DEFAULT 0 ,
`joinlevel`  int(11) NOT NULL DEFAULT 0 ,
`money`  int(11) NOT NULL DEFAULT 0 ,
`member`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`apply`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`state`  int(11) NOT NULL DEFAULT 0 ,
`cityinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`notice`  varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`createtime`  int(11) NOT NULL DEFAULT 0 ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
`fight`  int(11) NOT NULL ,
`camp`  int(11) NOT NULL ,
`exp`  int(11) NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`copy`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=85

;

-- ----------------------------
-- Table structure for `san_userbase`
-- ----------------------------
DROP TABLE IF EXISTS `san_userbase`;
CREATE TABLE `san_userbase` (
`uid`  bigint(20) NOT NULL COMMENT 'key' ,
`uname`  varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户名称' ,
`iconid`  int(5) NOT NULL DEFAULT 1002 COMMENT '图标' ,
`juqingid`  int(11) NOT NULL COMMENT '当前进度1' ,
`juqingid2`  int(11) NOT NULL DEFAULT 0 COMMENT '当前进度2' ,
`zhiyinid`  int(11) NOT NULL COMMENT '指引id' ,
`zhiyinid1`  int(11) NOT NULL DEFAULT 0 ,
`gem`  int(11) NOT NULL COMMENT '钻石' ,
`gold`  int(11) NOT NULL COMMENT '金币' ,
`exp`  int(11) NOT NULL COMMENT '经验' ,
`level`  int(11) NOT NULL COMMENT '等级' ,
`regtime`  timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '注册时间' ,
`lastlogintime`  timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后登录时间' ,
`lastupdtime`  timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后下线时间' ,
`lastlivetime`  timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '刷新次数' ,
`face`  int(10) NOT NULL DEFAULT 0 ,
`vip`  int(10) NOT NULL DEFAULT 0 COMMENT 'vip' ,
`vipexp`  int(10) NOT NULL DEFAULT 0 ,
`tili`  int(10) NOT NULL DEFAULT 0 COMMENT '体力' ,
`partyid`  int(10) NOT NULL DEFAULT 0 COMMENT '公会' ,
`skillpoint`  int(10) NOT NULL DEFAULT 10 ,
`tililastupdatatime`  int(10) NOT NULL DEFAULT 0 COMMENT '体力恢复时间' ,
`splastupdatatime`  int(10) NOT NULL DEFAULT 0 ,
`lastcheckintime`  timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后签到时间' ,
`checkinnum`  int(10) NOT NULL DEFAULT 0 COMMENT '签到次数' ,
`checkinaward`  int(11) NOT NULL DEFAULT 0 COMMENT '签到奖励' ,
`position`  int(11) NOT NULL DEFAULT 10010 ,
`isrename`  int(2) NOT NULL DEFAULT 0 ,
`logindays`  int(11) NOT NULL DEFAULT 1 COMMENT '登录天数' ,
`loginaward`  int(11) NOT NULL DEFAULT 0 COMMENT '登录奖励' ,
`levelaward`  int(11) NOT NULL DEFAULT 0 COMMENT '等级奖励' ,
`worldaward`  int(10) NOT NULL DEFAULT 0 ,
`morale`  int(11) NOT NULL DEFAULT 0 ,
`citylevel`  int(10) NOT NULL DEFAULT 0 COMMENT '城市等级' ,
`camp`  int(10) NOT NULL DEFAULT 0 ,
`city`  int(11) NOT NULL DEFAULT 0 COMMENT '城市位置' ,
`fight`  int(11) NOT NULL DEFAULT 0 COMMENT '战力' ,
`isgag`  int(11) NOT NULL DEFAULT 0 COMMENT '是否禁言' ,
`isblock`  int(11) NOT NULL DEFAULT 0 ,
`blockday`  int(11) NOT NULL DEFAULT 0 ,
`ip`  varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '上次登录ip' ,
`promotebox`  int(11) NOT NULL DEFAULT 0 ,
`linetime`  bigint(20) NOT NULL ,
PRIMARY KEY (`uid`),
UNIQUE INDEX `uname_2` (`uname`) USING BTREE ,
INDEX `uname` (`uname`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userbeauty3`
-- ----------------------------
DROP TABLE IF EXISTS `san_userbeauty3`;
CREATE TABLE `san_userbeauty3` (
`uid`  int(11) NOT NULL ,
`beautyinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
`count`  int(11) NOT NULL DEFAULT 0 ,
`fight`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userbuild`
-- ----------------------------
DROP TABLE IF EXISTS `san_userbuild`;
CREATE TABLE `san_userbuild` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_usercity`
-- ----------------------------
DROP TABLE IF EXISTS `san_usercity`;
CREATE TABLE `san_usercity` (
`uid`  bigint(20) NOT NULL ,
`campteam`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`cityevent`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`visitevent`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`feelings`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`citybox`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`treasure`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`power`  int(11) NOT NULL ,
`powertime`  bigint(20) NOT NULL ,
`visitouyu`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`ggzjhard`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`ggzjlevelnum`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`ggzjchosehard`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`curteam`  int(11) NOT NULL DEFAULT 0 ,
`state`  int(11) NOT NULL DEFAULT 0 ,
`cityid`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_usercitybuild`
-- ----------------------------
DROP TABLE IF EXISTS `san_usercitybuild`;
CREATE TABLE `san_usercitybuild` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userconsume`
-- ----------------------------
DROP TABLE IF EXISTS `san_userconsume`;
CREATE TABLE `san_userconsume` (
`uid`  int(11) NOT NULL ,
`gettilicount`  int(11) NOT NULL ,
`gettilimaxcount`  int(11) NOT NULL ,
`getgoldcount`  int(11) NOT NULL ,
`getgoldmaxcount`  int(11) NOT NULL ,
`getfreegold`  int(11) NOT NULL ,
`spopen`  int(2) NOT NULL ,
`getspcount`  int(5) NOT NULL ,
`npcaward`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_usercountrytask`
-- ----------------------------
DROP TABLE IF EXISTS `san_usercountrytask`;
CREATE TABLE `san_usercountrytask` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`choseStep`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userdart2`
-- ----------------------------
DROP TABLE IF EXISTS `san_userdart2`;
CREATE TABLE `san_userdart2` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`duan`  int(11) NOT NULL DEFAULT 0 ,
`exp`  int(11) NOT NULL DEFAULT 0 ,
`fightcount`  int(11) NOT NULL DEFAULT 0 ,
`fightmaxcount`  int(11) NOT NULL DEFAULT 0 ,
`addcounttime`  int(11) NOT NULL DEFAULT 0 ,
`acthero`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`dartcar`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`mydartcar`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`inmydartcarhero`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`diehero`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
`guidedartfinish`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userevent`
-- ----------------------------
DROP TABLE IF EXISTS `san_userevent`;
CREATE TABLE `san_userevent` (
`uid`  bigint(20) NOT NULL ,
`eventnum`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userfeats`
-- ----------------------------
DROP TABLE IF EXISTS `san_userfeats`;
CREATE TABLE `san_userfeats` (
`uid`  int(11) NOT NULL ,
`feats`  int(11) NOT NULL ,
`onebyone`  int(11) NOT NULL ,
`kills`  int(11) NOT NULL ,
`attack`  int(11) NOT NULL ,
`task`  int(11) NOT NULL ,
`boxnum`  int(11) NOT NULL ,
`boxnumall`  int(11) NOT NULL ,
`helpack`  int(11) NOT NULL ,
`helpackall`  int(11) NOT NULL ,
`award1`  int(11) NOT NULL ,
`award2`  int(11) NOT NULL ,
`award3`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userfind`
-- ----------------------------
DROP TABLE IF EXISTS `san_userfind`;
CREATE TABLE `san_userfind` (
`uid`  int(10) NOT NULL ,
`gemfreetime`  bigint(20) NOT NULL ,
`gold1`  int(11) NOT NULL ,
`gold9`  int(11) NOT NULL ,
`gold99`  int(11) NOT NULL ,
`gold101`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`gem1`  int(11) NOT NULL ,
`gem9`  int(11) NOT NULL ,
`gem99`  int(11) NOT NULL ,
`gem101`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`mj20`  int(11) NOT NULL ,
`sj20`  int(11) NOT NULL ,
`hgfreetime`  bigint(20) NOT NULL ,
`hg1`  int(11) NOT NULL ,
`hg9`  int(11) NOT NULL ,
`hg20`  int(11) NOT NULL ,
`yufu`  int(11) NOT NULL ,
`boxfreetime`  bigint(20) NOT NULL ,
`box1`  int(11) NOT NULL ,
`box9`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_usergamerecord`
-- ----------------------------
DROP TABLE IF EXISTS `san_usergamerecord`;
CREATE TABLE `san_usergamerecord` (
`uid`  int(11) NOT NULL ,
`getallgold`  int(11) NOT NULL DEFAULT 0 ,
`getallgem`  int(11) NOT NULL DEFAULT 0 ,
`consumegold`  int(11) NOT NULL DEFAULT 0 ,
`consumegoldrecord`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`consumegem`  int(11) NOT NULL DEFAULT 0 ,
`consumegemrecord`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`getalltax`  int(11) NOT NULL DEFAULT 0 ,
`allsection`  int(11) NOT NULL DEFAULT 0 ,
`jjcwin`  int(11) NOT NULL DEFAULT 0 ,
`jjcbest`  int(11) NOT NULL DEFAULT 0 ,
`orgheronum`  int(11) NOT NULL DEFAULT 0 ,
`orgartifatnum`  int(11) NOT NULL DEFAULT 0 ,
`allherostar`  int(11) NOT NULL DEFAULT 0 ,
`allpassstar`  int(11) NOT NULL DEFAULT 0 ,
`opensectiontime`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`addherotime`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`openmeinvtime`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`finishevent`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userhero2`
-- ----------------------------
DROP TABLE IF EXISTS `san_userhero2`;
CREATE TABLE `san_userhero2` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`totalstars`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`),
INDEX `uid` (`uid`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userjjc2`
-- ----------------------------
DROP TABLE IF EXISTS `san_userjjc2`;
CREATE TABLE `san_userjjc2` (
`uid`  int(10) NOT NULL COMMENT 'key' ,
`rankid`  int(11) NOT NULL AUTO_INCREMENT ,
`jjcpoint`  int(10) NOT NULL DEFAULT 0 ,
`featspoint`  int(10) NOT NULL DEFAULT 0 ,
`acthero`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`actfight`  int(11) NOT NULL DEFAULT 0 ,
`defhero`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`deffight`  int(11) NOT NULL DEFAULT 0 ,
`defver`  int(11) NOT NULL DEFAULT 0 ,
`fightcount`  int(5) NOT NULL DEFAULT 0 ,
`fightmaxcount`  int(5) NOT NULL DEFAULT 0 ,
`isfightting`  int(2) NOT NULL DEFAULT 0 ,
`isbefighting`  int(11) NOT NULL DEFAULT 0 ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
`worship`  int(10) NOT NULL DEFAULT 0 ,
`isrobot`  int(5) NOT NULL DEFAULT 0 ,
`bestrankid`  int(11) NOT NULL DEFAULT 0 ,
`getfightertime`  int(11) NOT NULL DEFAULT 0 ,
`beforerankid`  int(11) NOT NULL DEFAULT 0 ,
`lastfighttime`  int(11) NOT NULL DEFAULT 0 ,
`lastbefighttime`  int(11) NOT NULL DEFAULT 0 ,
`fighterinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`jjcrecord`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`bestaward`  int(11) NOT NULL DEFAULT 0 ,
`pointaward`  int(11) NOT NULL DEFAULT 0 ,
`isjoined`  int(2) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`rankid`),
UNIQUE INDEX `uid` (`uid`) USING BTREE 
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=1

;

-- ----------------------------
-- Table structure for `san_userking`
-- ----------------------------
DROP TABLE IF EXISTS `san_userking`;
CREATE TABLE `san_userking` (
`uid`  bigint(20) NOT NULL ,
`chose`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_usermail`
-- ----------------------------
DROP TABLE IF EXISTS `san_usermail`;
CREATE TABLE `san_usermail` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`recv`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userrecharge`
-- ----------------------------
DROP TABLE IF EXISTS `san_userrecharge`;
CREATE TABLE `san_userrecharge` (
`uid`  int(11) NOT NULL ,
`money`  int(11) NOT NULL DEFAULT 0 ,
`getallgem`  int(11) NOT NULL DEFAULT 0 ,
`type1`  int(11) NOT NULL DEFAULT 0 ,
`type2`  int(11) NOT NULL DEFAULT 0 ,
`type3`  int(11) NOT NULL DEFAULT 0 ,
`type4`  int(11) NOT NULL DEFAULT 0 ,
`type5`  int(11) NOT NULL DEFAULT 0 ,
`type6`  int(11) NOT NULL DEFAULT 0 ,
`record`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`firstaward`  int(5) NOT NULL DEFAULT 0 ,
`moneyday`  int(11) NOT NULL DEFAULT 0 ,
`moneyweek`  int(11) NOT NULL DEFAULT 0 ,
`monthcount1`  int(11) NOT NULL DEFAULT 0 ,
`monthcount2`  int(11) NOT NULL DEFAULT 0 ,
`monthcount3`  int(11) NOT NULL DEFAULT 0 ,
`vipbox`  int(20) NOT NULL DEFAULT 0 ,
`fundtype`  int(11) NOT NULL DEFAULT 0 ,
`fundget`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userscience`
-- ----------------------------
DROP TABLE IF EXISTS `san_userscience`;
CREATE TABLE `san_userscience` (
`uid`  bigint(20) NOT NULL ,
`sciencenum`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userspecial`
-- ----------------------------
DROP TABLE IF EXISTS `san_userspecial`;
CREATE TABLE `san_userspecial` (
`uid`  int(11) NOT NULL ,
`energyuse`  int(11) NOT NULL DEFAULT 0 ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`refreshcount`  int(11) NOT NULL DEFAULT 0 ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_usertask`
-- ----------------------------
DROP TABLE IF EXISTS `san_usertask`;
CREATE TABLE `san_usertask` (
`uid`  int(11) NOT NULL ,
`taskinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`liveness`  int(11) NOT NULL ,
`award`  int(11) NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_usertitle`
-- ----------------------------
DROP TABLE IF EXISTS `san_usertitle`;
CREATE TABLE `san_usertitle` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`level`  int(11) NOT NULL DEFAULT 1 ,
`task`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userunioninfo`
-- ----------------------------
DROP TABLE IF EXISTS `san_userunioninfo`;
CREATE TABLE `san_userunioninfo` (
`uid`  int(11) NOT NULL COMMENT 'key' ,
`position`  int(11) NOT NULL DEFAULT 0 ,
`donation`  int(11) NOT NULL DEFAULT 0 ,
`givecount`  int(11) NOT NULL DEFAULT 5 ,
`lastupdtime`  int(11) NOT NULL DEFAULT 0 ,
`unionid`  int(11) NOT NULL DEFAULT 0 ,
`applyinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`copynum`  int(11) NOT NULL ,
`copyver`  bigint(20) NOT NULL ,
`copyaward`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_userworldevent`
-- ----------------------------
DROP TABLE IF EXISTS `san_userworldevent`;
CREATE TABLE `san_userworldevent` (
`uid`  bigint(20) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_wartask`
-- ----------------------------
DROP TABLE IF EXISTS `san_wartask`;
CREATE TABLE `san_wartask` (
`uid`  bigint(20) NOT NULL ,
`taskinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_weekplan`
-- ----------------------------
DROP TABLE IF EXISTS `san_weekplan`;
CREATE TABLE `san_weekplan` (
`uid`  int(11) NOT NULL ,
`finishcount`  int(11) NOT NULL DEFAULT 0 ,
`type1`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type2`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type3`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type4`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type5`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type6`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type7`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`type8`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
`completeaward`  int(11) NOT NULL DEFAULT 0 ,
`lastupdtime`  int(11) NOT NULL ,
`taskinfo`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`uid`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Table structure for `san_worldevent`
-- ----------------------------
DROP TABLE IF EXISTS `san_worldevent`;
CREATE TABLE `san_worldevent` (
`id`  int(11) NOT NULL ,
`info`  text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL ,
PRIMARY KEY (`id`)
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci

;

-- ----------------------------
-- Auto increment value for `san_account`
-- ----------------------------
ALTER TABLE `san_account` AUTO_INCREMENT=7967;

-- ----------------------------
-- Auto increment value for `san_activitymask`
-- ----------------------------
ALTER TABLE `san_activitymask` AUTO_INCREMENT=9005;

-- ----------------------------
-- Auto increment value for `san_belog`
-- ----------------------------
ALTER TABLE `san_belog` AUTO_INCREMENT=5438;

-- ----------------------------
-- Auto increment value for `san_heroequ2`
-- ----------------------------
ALTER TABLE `san_heroequ2` AUTO_INCREMENT=1;

-- ----------------------------
-- Auto increment value for `san_linelog`
-- ----------------------------
ALTER TABLE `san_linelog` AUTO_INCREMENT=1553;

-- ----------------------------
-- Auto increment value for `san_log`
-- ----------------------------
ALTER TABLE `san_log` AUTO_INCREMENT=11711463;

-- ----------------------------
-- Auto increment value for `san_mail`
-- ----------------------------
ALTER TABLE `san_mail` AUTO_INCREMENT=52;

-- ----------------------------
-- Auto increment value for `san_passguide`
-- ----------------------------
ALTER TABLE `san_passguide` AUTO_INCREMENT=1229101;

-- ----------------------------
-- Auto increment value for `san_unioninfo`
-- ----------------------------
ALTER TABLE `san_unioninfo` AUTO_INCREMENT=85;

-- ----------------------------
-- Auto increment value for `san_userjjc2`
-- ----------------------------
ALTER TABLE `san_userjjc2` AUTO_INCREMENT=1;
