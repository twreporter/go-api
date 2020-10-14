ALTER TABLE `pay_by_prime_donations` ADD `customer_service_notes` varchar(256) DEFAULT NULL;
ALTER TABLE `periodic_donations` ADD `customer_service_notes` varchar(256) DEFAULT NULL;
ALTER TABLE `pay_by_other_method_donations` ADD `customer_service_notes` varchar(256) DEFAULT NULL;
