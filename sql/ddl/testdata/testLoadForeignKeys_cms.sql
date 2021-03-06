SET foreign_key_checks = 0;
SET NAMES utf8mb4;

DROP TABLE IF EXISTS `x910store`;
CREATE TABLE `x910store` (
  `store_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(32) DEFAULT NULL,
  `website_id` smallint(5) unsigned NOT NULL DEFAULT 0,
  `group_id` smallint(5) unsigned NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `sort_order` smallint(5) unsigned NOT NULL DEFAULT 0,
  `is_active` smallint(5) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`store_id`),
  UNIQUE KEY `STORE_CODE` (`code`),
  KEY `STORE_WEBSITE_ID` (`website_id`),
  KEY `STORE_IS_ACTIVE_SORT_ORDER` (`is_active`,`sort_order`),
  KEY `STORE_GROUP_ID` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `x910cms_block`;
CREATE TABLE `x910cms_block` (
  `block_id` smallint(6) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`block_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `x910cms_block_store`;
CREATE TABLE `x910cms_block_store` (
  `block_id` smallint(6) NOT NULL,
  `store_id` smallint(5) unsigned NOT NULL,
  PRIMARY KEY (`block_id`,`store_id`),
  KEY `CMS_BLOCK_STORE_STORE_ID` (`store_id`),
  CONSTRAINT `CMS_BLOCK_STORE_BLOCK_ID_CMS_BLOCK_BLOCK_ID` FOREIGN KEY (`block_id`) REFERENCES `x910cms_block` (`block_id`) ON DELETE CASCADE,
  CONSTRAINT `CMS_BLOCK_STORE_STORE_ID_STORE_STORE_ID` FOREIGN KEY (`store_id`) REFERENCES `x910store` (`store_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


DROP TABLE IF EXISTS `x910cms_page`;
CREATE TABLE `x910cms_page` (
  `page_id` smallint(6) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`page_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `x910cms_page_store`;
CREATE TABLE `x910cms_page_store` (
  `page_id` smallint(6) NOT NULL,
  `store_id` smallint(5) unsigned NOT NULL,
  PRIMARY KEY (`page_id`,`store_id`),
  KEY `CMS_PAGE_STORE_STORE_ID` (`store_id`),
  CONSTRAINT `CMS_PAGE_STORE_PAGE_ID_CMS_PAGE_PAGE_ID` FOREIGN KEY (`page_id`) REFERENCES `x910cms_page` (`page_id`) ON DELETE CASCADE,
  CONSTRAINT `CMS_PAGE_STORE_STORE_ID_STORE_STORE_ID` FOREIGN KEY (`store_id`) REFERENCES `x910store` (`store_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `x910catalog_eav_attribute`;
CREATE TABLE `x910catalog_eav_attribute` (
  `attribute_id` smallint(5) unsigned NOT NULL COMMENT 'Attribute ID',
  PRIMARY KEY (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET foreign_key_checks = 1;
