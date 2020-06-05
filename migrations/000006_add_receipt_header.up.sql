ALTER TABLE `pay_by_prime_donations` ADD `receipt_header` varchar(128) DEFAULT NULL;
ALTER TABLE `periodic_donations` ADD `receipt_header` varchar(128) DEFAULT NULL;
ALTER TABLE `pay_by_card_token_donations` ADD `receipt_header` varchar(128) DEFAULT NULL;
ALTER TABLE `pay_by_other_method_donations` ADD `receipt_header` varchar(128) DEFAULT NULL;
