CREATE TABLE `quant`.`kline` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `open_time` DATETIME NULL,
  `close_time` DATETIME NULL,
  `open` DOUBLE NULL,
  `close` DOUBLE NULL,
  `high` DOUBLE NULL,
  `low` DOUBLE NULL,
  `volume` DOUBLE NULL COMMENT '总成交量',
  `buy_volume` DOUBLE NULL COMMENT '买成交量',
  `sell_volume` DOUBLE NULL COMMENT '卖成交量',
  `trade_number` INT NULL COMMENT '成交比数',
  `quote` DOUBLE NULL COMMENT '总成交额',
  `buy_quote` DOUBLE NULL COMMENT '买成交额'
  `sell_quote` DOUBLE NULL COMMENT '卖成交额'
  PRIMARY KEY (`id`))
COMMENT = 'k线数据';
