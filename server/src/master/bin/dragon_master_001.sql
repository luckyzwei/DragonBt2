/*
 Navicat MySQL Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 50724
 Source Host           : localhost:3306
 Source Schema         : dragon_master_001

 Target Server Type    : MySQL
 Target Server Version : 50724
 File Encoding         : 65001

 Date: 25/05/2020 11:51:31
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for tbl_account
-- ----------------------------
DROP TABLE IF EXISTS `tbl_account`;
CREATE TABLE `tbl_account`  (
  `id` int(11) NOT NULL COMMENT '帐号ID',
  `account` varchar(63) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '帐号',
  `password` varchar(63) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '密码',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tbl_server
-- ----------------------------
DROP TABLE IF EXISTS `tbl_server`;
CREATE TABLE `tbl_server`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '表主键',
  `serverid` int(11) NOT NULL COMMENT '服务器ID',
  `name` varchar(63) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '服务器名字',
  `online` int(11) NOT NULL DEFAULT 0 COMMENT '实时在线人数',
  `lastupdate` int(11) NOT NULL DEFAULT 0 COMMENT '上次更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tbl_union
-- ----------------------------
DROP TABLE IF EXISTS `tbl_union`;
CREATE TABLE `tbl_union`  (
    `id` int(11) NOT NULL COMMENT '公会ID',
    icon int(11) NOT NULL COMMENT 'icon',
    unionname varchar(63) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '公会名称',
    masteruid bigint(20) NOT NULL COMMENT '所有者ID',
    mastername varchar(63) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
    level int(11) NOT NULL COMMENT '公会等级',
    jointype int(11) NOT NULL COMMENT '加入类型',
    joinlevel int(11) NOT NULL COMMENT '加入等级',
    serverid int(11) NOT NULL COMMENT '服务器id',
    notice text NOT NULL COMMENT '公告',
    board text NOT NULL COMMENT '对外展示',
    createtime bigint(20) NOT NULL COMMENT '创建时间',
    lastupdtime bigint(20) NOT NULL COMMENT '更新时间',
    fight bigint(20) NOT NULL COMMENT '总战力',
    exp int(11) NOT NULL COMMENT '经验',
    dayexp int(11) NOT NULL COMMENT '每日经验',
    activitypoint int(11) NOT NULL COMMENT '活跃点数',
    acitvitylimit int(11) NOT NULL COMMENT '活跃度限额',
    mailcd bigint(20) NOT NULL COMMENT '邮件cd',
    member text NOT NULL COMMENT '成员列表',
    applys text NOT NULL COMMENT '申请列表',
    record text NOT NULL COMMENT '操作记录',
    huntinfo text NOT NULL COMMENT '军团狩猎记录',
    bravehand text NOT NULL COMMENT '无畏之手',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tbl_user
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user`;
CREATE TABLE `tbl_user`  (
  `uid` bigint(11) NOT NULL COMMENT '角色id',
  `uname` varchar(63) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '昵称',
  `level` int(11) NULL DEFAULT NULL COMMENT '等级',
  `fight` int(11) NULL DEFAULT NULL COMMENT '战力',
  `passid` int(11) NULL DEFAULT NULL COMMENT '关卡id',
  `serverid` int(11) NULL DEFAULT NULL COMMENT '服务器id',
  `regtime` int(11) NULL DEFAULT NULL COMMENT '注册时间',
  `logintime` int(11) NULL DEFAULT NULL COMMENT '登录时间',
  `lastupdate` int(11) NULL DEFAULT NULL COMMENT '上次更新时间',
  `data` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '基础数据',
  `heros` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '阵容列表',
  `equips` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '阵容装备列表',
  PRIMARY KEY (`uid`) USING BTREE
) ENGINE = MyISAM CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tbl_userfriend
-- ----------------------------
DROP TABLE IF EXISTS `tbl_userfriend`;
CREATE TABLE `tbl_userfriend`  (
  `uid` bigint(11) NOT NULL COMMENT '注册ID',
  `friends` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '好友列表',
  `applys` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '申请列表',
  `applieds` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '被申请列表',
  PRIMARY KEY (`uid`) USING BTREE
) ENGINE = MyISAM CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
