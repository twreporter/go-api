CREATE TABLE IF NOT EXISTS `receipt_serial_numbers` (
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `serial_number` int(10) unsigned NOT NULL,
  `YYYYMM` varchar(6) NOT NULL,
  PRIMARY KEY (`YYYYMM`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
