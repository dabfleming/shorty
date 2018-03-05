SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for url
-- ----------------------------
DROP TABLE IF EXISTS `url`;
CREATE TABLE `url` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `slug` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `url` varchar(4096) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug_idx` (`slug`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of url
-- ----------------------------
BEGIN;
INSERT INTO `url` VALUES (1, 'goog', 'https://www.google.ca/');
INSERT INTO `url` VALUES (2, 'twitter', 'https://twitter.com/');
INSERT INTO `url` VALUES (3, 'fb', 'https://www.facebook.com/');
INSERT INTO `url` VALUES (4, 'yt', 'https://www.youtube.com/');
COMMIT;

-- ----------------------------
-- Table structure for visit
-- ----------------------------
DROP TABLE IF EXISTS `visit`;
CREATE TABLE `visit` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url_id` int(11) NOT NULL,
  `device` varchar(100) NOT NULL,
  `os` varchar(100) NOT NULL,
  `browser` varchar(100) NOT NULL,
  `ip` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `url_id` (`url_id`),
  CONSTRAINT `visit_ibfk_1` FOREIGN KEY (`url_id`) REFERENCES `url` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of visit
-- ----------------------------
BEGIN;
INSERT INTO `visit` VALUES (1, 1, 'Other', 'Mac OS X', 'Chrome', '127.0.0.1', '2018-03-05 04:16:06');
INSERT INTO `visit` VALUES (2, 1, 'Nexus 5', 'Android', 'Chrome Mobile', '127.0.0.1', '2018-03-05 05:08:38');
INSERT INTO `visit` VALUES (3, 1, 'iPhone', 'iOS', 'Chrome Mobile iOS', '127.0.0.1', '2018-03-05 05:08:46');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
