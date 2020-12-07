-- 创建db
Create Database gocamp;

use gocamp;

-- 创建表
CREATE TABLE IF NOT EXISTS WK02 (
  `PId` bigint(20) NOT NULL COMMENT 'PId',
  `UserId` varchar(32) NOT NULL COMMENT 'UserId',
  PRIMARY KEY (`PId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
