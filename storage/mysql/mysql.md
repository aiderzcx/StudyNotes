# Mysql限制

## 每一行的最大长度
	InnoDB：The maximum row size for the used table type, not counting BLOBs, is 65535

## 每一行的最大字段数量
	innodb引擎支持最大字段上线为1017
	myisam引擎最大字段上限为2410 

## vachar的最大长度
	latin1字符集下的表varchar上限为65532，即一个字符一个字节
	utf8字符集下的表varchar上限为21844，即一个字符三个字节 65535-1-2 结果除以3 ==21844
		 -1表示第一个字节不存数据，-2表示两个字节存放varchar的长度，除以3是utf8字符特性，一个字符三个字节

## MySQL数据类型（留作备忘）
		类 型			大 小				描 述
	CAHR(n)				n字节			定长字段，长度为0~255个字符
	VARCHAR(n)			n+1或+2字节		变长字段，长度为0~65 535个字符
	TINYTEXT			n+1字节			字符串，最大长度为255个字符
	TEXT				n+2字节			字符串，最大长度为65 535个字符
	MEDIUMINT			n+3字节			字符串，最大长度为16 777 215个字符
	LONGTEXT			n+4字节			字符串，最大长度为4 294 967 295个字符
	TINYINT(n)			1字节			范围：-128~127，或者0~255（无符号）
	SMALLINT(n)			2字节			范围：-32 768~32 767，或者0~65 535（无符号）
	MEDIUMINT(n)		3字节			范围：-838608~8388607，0~16777215（无符号）
	INT(n)				4字节			范围：-2147483648~2147483647, 0~4294967295（无符号）
	BIGINT(n)			8字节			范围：-9223372036854775808~9223372036854775807，
										 		 或者0~18446744073709551615（无符号）
	FLOAT(n,Decimals)	4字节			具有浮动小数点的较小的数
	DOUBLE(n,Decimals)	8字节			具有浮动小数点的较大的数
	DECIMAL(n,Decimals)	n+1或+2字节		存储为字符串的DOUBLE，允许固定的小数点
	DATE				3字节			采用YYYY-MM-DD格式
	DATETIME			8字节			采用YYYY-MM-DD HH:MM:SS格式
	TIMESTAMP			4字节			采用YYYYMMDDHHMMSS格式；可接受的范围终止于2037年
	TIME				3字节			采用HH:MM:SS格式
	ENUM				1或2字节			Enumeration(枚举)的简写，这意味着每一列都可以具有多个可能的值之一
	SET				1、2、3、4或8字节		与ENUM一样，只不过每一列都可以具有多个可能的值




