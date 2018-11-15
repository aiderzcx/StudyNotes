# 系统启动前的数据库相关的初始化
## 创建数据库相关的表
	CREATE TABLE `pay_order` (
	  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
	  `product_id` varchar(64) NOT NULL COMMENT '医疗系统的ID',
	  `product_desc` varchar(64) NOT NULL COMMENT '产品的描述信息',
	  `total_fee` int(11) NOT NULL COMMENT '需要支付的钱，单位：分',
	  `pay_id` varchar(64) COMMENT '系统的支付ID',
	  `pay_type` varchar(32)  COMMENT '支付方式',
	  `third_pre_id` varchar(64) COMMENT '第三方的预支付ID',
	  `third_id` varchar(64) COMMENT '第三方的支付ID',
	  `state` int(11) default 0 COMMENT '订单状态',
	  `create_at` varchar(32) NOT NULL COMMENT '创建时间',
	  `pay_at`  varchar(32) default '' COMMENT '支付时间',
	  `refund_at`  varchar(32) default '' COMMENT '退款时间',
	  `query_count` int(11) default 0 COMMENT '订单查询的次数',
	  `remarks` varchar(128) default '' COMMENT '备注信息',
	  PRIMARY KEY (`id`),
	  KEY `i_product_id` (`product_id`),
	  KEY `i_pay_id` (`pay_id`),
	  KEY `i_third_id` (`third_id`),
	  KEY `i_create` (`create_at`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='产品支付的表'