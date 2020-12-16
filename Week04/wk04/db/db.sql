-- 创建db
Create Database gocamp;

use gocamp;

-- 创建表
CREATE TABLE IF NOT EXISTS wk04 (
  `Id` bigint(20) NOT NULL COMMENT 'Id',
  `Password` varchar(50) NOT NULL COMMENT 'Password',
  `Name` varchar(50) NOT NULL COMMENT 'Name',
  PRIMARY KEY (`Id`,`Password`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 假设这里是另外一种数据库方式，通过grom来简单模拟
Create Database gocamp_bak;

use gocamp_bak;

-- 创建表
CREATE TABLE IF NOT EXISTS wk04_bak (
  `Id` bigint(20) NOT NULL COMMENT 'Id',
  `Password` varchar(50) NOT NULL COMMENT 'Password',
  `Name` varchar(50) NOT NULL COMMENT 'Name',
  PRIMARY KEY (`Id`,`Password`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
