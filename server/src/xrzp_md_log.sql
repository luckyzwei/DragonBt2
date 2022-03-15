/*
Navicat MySQL Data Transfer

Source Server         : 192.168.10.155[本机]
Source Server Version : 50717
Source Host           : 192.168.10.155:3306
Source Database       : xrzp_md_log

Target Server Type    : MYSQL
Target Server Version : 50717
File Encoding         : 65001

Date: 2018-07-23 14:07:15
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for san_belog
-- ----------------------------
DROP TABLE IF EXISTS `san_belog`;
CREATE TABLE `san_belog` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `time` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `value` int(11) NOT NULL,
  `param1` int(11) NOT NULL,
  `param2` int(11) NOT NULL,
  `uid` bigint(20) NOT NULL,
  `dec` varchar(45) NOT NULL,
  `cur` int(11) NOT NULL,
  `param3` int(11) NOT NULL,
  `level` int(4) NOT NULL,
  `vip` int(4) NOT NULL,
  `fight` int(4) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `INDEX_UID` (`uid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for san_linelog
-- ----------------------------
DROP TABLE IF EXISTS `san_linelog`;
CREATE TABLE `san_linelog` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `time` bigint(20) NOT NULL,
  `ip` varchar(45) NOT NULL,
  `line` int(11) NOT NULL,
  `creator` varchar(45) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `INDEX_UID` (`uid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=78 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for san_log
-- ----------------------------
DROP TABLE IF EXISTS `san_log`;
CREATE TABLE `san_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `time` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `value` int(11) NOT NULL,
  `param1` int(11) NOT NULL,
  `param2` int(11) NOT NULL,
  `param3` int(11) NOT NULL DEFAULT '0',
  `uid` bigint(20) NOT NULL,
  `dec` varchar(45) NOT NULL,
  `cur` bigint(11) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `INDEX_UID` (`uid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1137 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for san_mail
-- ----------------------------
DROP TABLE IF EXISTS `san_mail`;
CREATE TABLE `san_mail` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '唯一Id',
  `uid` bigint(20) NOT NULL DEFAULT '0' COMMENT '玩家Id',
  `msg` text COLLATE utf8_unicode_ci NOT NULL COMMENT '消息',
  `item` text COLLATE utf8_unicode_ci NOT NULL COMMENT '物品信息',
  `time` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for san_userbase2
-- ----------------------------
DROP TABLE IF EXISTS `san_userbase2`;
CREATE TABLE `san_userbase2` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL COMMENT 'key',
  `timestamp` int(11) NOT NULL,
  `uname` varchar(64) NOT NULL COMMENT '用户名称',
  `iconid` int(5) NOT NULL DEFAULT '1002' COMMENT '图标',
  `juqingid` int(11) NOT NULL COMMENT '当前进度1',
  `juqingid2` int(11) NOT NULL DEFAULT '0' COMMENT '当前进度2',
  `zhiyinid` int(11) NOT NULL COMMENT '指引id',
  `zhiyinid1` int(11) NOT NULL DEFAULT '0',
  `gem` int(11) NOT NULL COMMENT '钻石',
  `gold` int(11) NOT NULL COMMENT '金币',
  `exp` int(11) NOT NULL COMMENT '经验',
  `level` int(11) NOT NULL COMMENT '等级',
  `regtime` timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '注册时间',
  `lastlogintime` timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后登录时间',
  `lastupdtime` timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后下线时间',
  `lastlivetime` timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '刷新次数',
  `face` int(10) NOT NULL DEFAULT '0',
  `vip` int(10) NOT NULL DEFAULT '0' COMMENT 'vip',
  `vipexp` int(10) NOT NULL DEFAULT '0',
  `tili` int(10) NOT NULL DEFAULT '0' COMMENT '体力',
  `partyid` int(10) NOT NULL DEFAULT '0' COMMENT '军团',
  `skillpoint` int(10) NOT NULL DEFAULT '10',
  `tililastupdatatime` int(10) NOT NULL DEFAULT '0' COMMENT '体力恢复时间',
  `splastupdatatime` int(10) NOT NULL DEFAULT '0',
  `lastcheckintime` timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后签到时间',
  `checkinnum` int(10) NOT NULL DEFAULT '0' COMMENT '签到次数',
  `checkinaward` int(11) NOT NULL DEFAULT '0' COMMENT '签到奖励',
  `position` int(11) NOT NULL DEFAULT '10010',
  `isrename` int(2) NOT NULL DEFAULT '0',
  `logindays` int(11) NOT NULL DEFAULT '1' COMMENT '登录天数',
  `loginaward` int(11) NOT NULL DEFAULT '0' COMMENT '登录奖励',
  `levelaward` int(11) NOT NULL DEFAULT '0' COMMENT '等级奖励',
  `worldaward` int(10) NOT NULL DEFAULT '0',
  `morale` int(11) NOT NULL DEFAULT '0',
  `citylevel` int(10) NOT NULL DEFAULT '0' COMMENT '城市等级',
  `camp` int(10) NOT NULL DEFAULT '0',
  `city` int(11) NOT NULL DEFAULT '0' COMMENT '城市位置',
  `fight` int(11) NOT NULL DEFAULT '0' COMMENT '战力',
  `isgag` int(11) NOT NULL DEFAULT '0' COMMENT '是否禁言',
  `isblock` int(11) NOT NULL DEFAULT '0',
  `blockday` int(11) NOT NULL,
  `ip` varchar(64) NOT NULL DEFAULT '0' COMMENT '上次登录ip',
  `promotebox` int(11) NOT NULL DEFAULT '0',
  `linetime` bigint(20) NOT NULL DEFAULT '0',
  `passmax` int(11) NOT NULL DEFAULT '0',
  `fitserver` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `INDEX_UID` (`uid`,`timestamp`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
