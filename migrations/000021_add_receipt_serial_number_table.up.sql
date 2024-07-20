CREATE TABLE IF NOT EXISTS `receipt_serial_numbers` (
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `serial_number` int(10) unsigned NOT NULL,
  `month` varchar(6) NOT NULL,
  PRIMARY KEY (`month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;