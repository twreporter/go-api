--
-- Table structure for table `bookmarks`
--

DROP TABLE IF EXISTS `bookmarks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bookmarks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `slug` varchar(100) NOT NULL,
  `host` varchar(100) NOT NULL,
  `is_external` tinyint(4) DEFAULT '0',
  `title` varchar(100) NOT NULL,
  `category` varchar(20) DEFAULT NULL,
  `authors` varchar(10000) DEFAULT NULL,
  `pub_date` int(10) unsigned DEFAULT NULL,
  `desc` varchar(250) DEFAULT NULL,
  `thumbnail` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_bookmarks_slug_host` (`slug`,`host`),
  KEY `idx_bookmarks_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=119 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `first_name` varchar(50) DEFAULT NULL,
  `last_name` varchar(50) DEFAULT NULL,
  `security_id` varchar(20) DEFAULT NULL,
  `passport_id` varchar(30) DEFAULT NULL,
  `city` varchar(45) DEFAULT NULL,
  `state` varchar(45) DEFAULT NULL,
  `country` varchar(45) DEFAULT NULL,
  `zip` varchar(20) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `privilege` int(5) NOT NULL,
  `registration_date` timestamp NULL DEFAULT NULL,
  `birthday` timestamp NULL DEFAULT NULL,
  `gender` varchar(2) DEFAULT NULL,
  `education` varchar(20) DEFAULT NULL,
  `enable_email` int(5) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=790 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `o_auth_accounts`
--

DROP TABLE IF EXISTS `o_auth_accounts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `o_auth_accounts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `type` varchar(10) DEFAULT NULL,
  `a_id` varchar(255) NOT NULL,
  `email` varchar(100) DEFAULT NULL,
  `name` varchar(80) DEFAULT NULL,
  `first_name` varchar(50) DEFAULT NULL,
  `last_name` varchar(50) DEFAULT NULL,
  `gender` varchar(20) DEFAULT NULL,
  `picture` varchar(255) DEFAULT NULL,
  `birthday` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_o_auth_accounts_deleted_at` (`deleted_at`),
  KEY `fk_o_auth_accounts_users1_idx` (`user_id`),
  CONSTRAINT `fk_o_auth_accounts_users1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=758 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `services`
--

DROP TABLE IF EXISTS `services`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `services` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_services_name` (`name`),
  KEY `idx_services_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `registrations`
--

DROP TABLE IF EXISTS `registrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `registrations` (
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `email` varchar(100) NOT NULL,
  `active` tinyint(1) DEFAULT '0',
  `activate_token` varchar(20) DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `services_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`email`),
  KEY `fk_registrations_users1_idx` (`user_id`),
  KEY `fk_registrations_services1_idx` (`services_id`),
  CONSTRAINT `fk_registrations_services1` FOREIGN KEY (`services_id`) REFERENCES `services` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_registrations_users1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `reporter_accounts`
--

DROP TABLE IF EXISTS `reporter_accounts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `reporter_accounts` (
  `user_id` int(10) unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(100) NOT NULL,
  `activate_token` varchar(50) DEFAULT NULL,
  `act_exp_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_reporter_accounts_email` (`email`),
  KEY `fk_reporter_accounts_users1_idx` (`user_id`),
  CONSTRAINT `fk_reporter_accounts_users1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users_bookmarks`
--

DROP TABLE IF EXISTS `users_bookmarks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users_bookmarks` (
  `user_id` int(10) unsigned NOT NULL,
  `bookmark_id` int(10) unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`,`bookmark_id`),
  KEY `fk_users_has_bookmarks_bookmarks1_idx` (`bookmark_id`),
  KEY `fk_users_has_bookmarks_users1_idx` (`user_id`),
  CONSTRAINT `fk_users_has_bookmarks_bookmarks1` FOREIGN KEY (`bookmark_id`) REFERENCES `bookmarks` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_users_has_bookmarks_users1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `web_push_subs
--

DROP TABLE IF EXISTS `web_push_subs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `web_push_subs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `endpoint` varchar(500) NOT NULL,
  `crc32_endpoint` int(10) unsigned NOT NULL,
  `keys` varchar(200) NOT NULL,
  `expiration_time` timestamp NULL DEFAULT NULL,
  `user_id` int(10) unsigned NULL DEFAULT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_web_push_subs_endpoint` (`endpoint`),
  KEY `idx_web_push_subs_crc32_endpoint` (`crc32_endpoint`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `pay_by_prime_donations`
--

DROP TABLE IF EXISTS `pay_by_prime_donations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pay_by_prime_donations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `details` varchar(50) NOT NULL,
  `merchant_id` varchar(30) NOT NULL,
  `amount` int(10) unsigned NOT NULL,
  `order_number` varchar(50) NOT NULL,
  `currency` char(3) DEFAULT 'TWD' NOT NULL,
  `pay_method` enum('credit_card', 'line', 'apple', 'google', 'samsung') NOT NULL,
  `status` enum('paying', 'paid', 'fail') NOT NULL,
  `tappay_api_status` int NULL DEFAULT NULL,
  `msg` varchar(100) NULL DEFAULT NULL,
  `tappay_record_status` int NULL DEFAULT NULL,
  `rec_trade_id` varchar(20) NULL DEFAULT NULL,
  `bank_transaction_id` varchar(50) NULL DEFAULT NULL,
  `auth_code` varchar(6) NULL DEFAULT NULL,
  `acquirer` varchar(50) NULL DEFAULT NULL,
  `transaction_time` timestamp NULL DEFAULT NULL,
  `bank_transaction_start_time` timestamp NULL DEFAULT NULL,
  `bank_transaction_end_time` timestamp NULL DEFAULT NULL,
  `bank_result_code` varchar(50) NULL DEFAULT NULL,
  `bank_result_msg` varchar(50) NULL DEFAULT NULL,
  `cardholder_email` varchar(100) NOT NULL,
  `cardholder_phone_number` varchar(20) DEFAULT NULL,
  `cardholder_name` varchar(30) DEFAULT NULL,
  `cardholder_zip_code` varchar(10) DEFAULT NULL,
  `cardholder_address` varchar(100) DEFAULT NULL,
  `cardholder_national_id` varchar(20) DEFAULT NULL,
  `card_info_bin_code` varchar(6) DEFAULT NULL,  
  `card_info_last_four` varchar(4) DEFAULT NULL,
  `card_info_issuer` varchar(50) DEFAULT NULL,
  `card_info_funding` tinyint DEFAULT NULL,  
  `card_info_type` tinyint DEFAULT NULL, 
  `card_info_level` varchar(255) DEFAULT NULL, 
  `card_info_country` varchar(30) DEFAULT NULL, 
  `card_info_country_code` varchar(10) DEFAULT NULL, 
  `card_info_expiry_date` varchar(6) DEFAULT NULL, 
  `send_receipt` enum('yearly', 'monthly', 'no') DEFAULT 'yearly',
  `notes` varchar(100) DEFAULT NULL,
  `is_anonymous` tinyint(1) DEFAULT 0,
  `linepay_method` enum('CREDIT_CARD', 'BALANCE', 'POINT') DEFAULT NULL,
  `linepay_point` int DEFAULT NULL, 

  PRIMARY KEY (`id`),
  KEY `idx_pay_by_prime_donations_status` (`status`),
  KEY `idx_pay_by_prime_donations_pay_method` (`pay_method`),
  KEY `idx_pay_by_prime_donations_order_number` (`order_number`),
  CONSTRAINT `fk_pay_by_prime_donations_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `pay_by_other_method_donations`
--

DROP TABLE IF EXISTS `pay_by_other_method_donations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pay_by_other_method_donations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `email` varchar(100) NOT NULL,
  `pay_method` varchar(50) NOT NULL,
  `amount` int(10) unsigned NOT NULL,
  `currency` char(3) DEFAULT 'TWD' NOT NULL,
  `order_number` varchar(50) NOT NULL,
  `details` varchar(50) NOT NULL,
  `merchant_id` varchar(30) NOT NULL,
  `send_receipt` enum('yearly', 'monthly', 'no') DEFAULT 'yearly',
  `phone_number` varchar(20) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  `zip_code` varchar(10) DEFAULT NULL,
  `address` varchar(100) DEFAULT NULL,
  `national_id` varchar(20) DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pay_by_other_method_donations_pay_method` (`pay_method`),
  KEY `idx_pay_by_other_method_donations_amount` (`amount`),
  KEY `idx_pay_by_other_method_order_number` (`order_number`),
  CONSTRAINT `fk_pay_by_other_method_donations_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `periodic_donations`
--

DROP TABLE IF EXISTS `periodic_donations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `periodic_donations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `status` enum('to_pay', 'paying', 'paid', 'fail', 'stopped', 'invalid') NOT NULL,
  `card_token` tinyblob NULL DEFAULT NULL,
  `card_key` tinyblob NULL DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `currency` char(3) DEFAULT 'TWD' NOT NULL,
  `amount` int(10) unsigned NOT NULL,
  `last_success_at` timestamp NULL DEFAULT NULL,
  `cardholder_email` varchar(100) NOT NULL,
  `cardholder_phone_number` varchar(20) DEFAULT NULL,
  `cardholder_name` varchar(30) DEFAULT NULL,
  `cardholder_zip_code` varchar(10) DEFAULT NULL,
  `cardholder_address` varchar(100) DEFAULT NULL,
  `cardholder_national_id` varchar(20) DEFAULT NULL,
  `card_info_bin_code` varchar(6) DEFAULT NULL,  
  `card_info_last_four` varchar(4) DEFAULT NULL,
  `card_info_issuer` varchar(50) DEFAULT NULL,
  `card_info_funding` tinyint DEFAULT NULL,  
  `card_info_type` tinyint DEFAULT NULL, 
  `card_info_level` varchar(255) DEFAULT NULL, 
  `card_info_country` varchar(30) DEFAULT NULL, 
  `card_info_country_code` varchar(10) DEFAULT NULL, 
  `card_info_expiry_date` varchar(6) DEFAULT NULL, 
  `send_receipt` enum('yearly', 'no') DEFAULT 'yearly',
  `to_feedback` tinyint(1) DEFAULT 1,
  `order_number` varchar(50) NOT NULL,
  `details` varchar(50) NOT NULL,
  `frequency` enum('monthly', 'yearly') DEFAULT 'monthly',
  `notes` varchar(100) DEFAULT NULL,
  `max_paid_times` int NOT NULL DEFAULT 2147483647,
  `is_anonymous` tinyint(1) DEFAULT 0,

  PRIMARY KEY (`id`),
  KEY `idx_periodic_donations_status` (`status`),
  KEY `idx_periodic_donations_amount` (`amount`),
  KEY `idx_periodic_donations_order_number` (`order_number`),
  KEY `idx_periodic_donations_last_success_at` (`last_success_at`),
  CONSTRAINT `fk_periodic_donations_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `pay_by_card_token_donations`
--

DROP TABLE IF EXISTS `pay_by_card_token_donations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pay_by_card_token_donations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `periodic_id` int(10) unsigned NOT NULL,
  `status` enum('paying', 'paid', 'fail') NOT NULL,
  `tappay_api_status` int NULL DEFAULT NULL, 
  `msg` varchar(100) NULL DEFAULT NULL,
  `tappay_record_status` int NULL DEFAULT NULL,
  `currency` char(3) DEFAULT 'TWD' NOT NULL,
  `amount` int(10) unsigned NOT NULL,
  `merchant_id` varchar(30) NOT NULL,
  `order_number` varchar(50) NOT NULL,
  `details` varchar(50) NULL DEFAULT NULL,
  `rec_trade_id` varchar(20) NULL DEFAULT NULL,
  `bank_transaction_id` varchar(50) NULL DEFAULT NULL,
  `auth_code` varchar(6) NULL DEFAULT NULL,
  `acquirer` varchar(50) NULL DEFAULT NULL,
  `transaction_time` timestamp NULL DEFAULT NULL,
  `bank_transaction_start_time` timestamp NULL DEFAULT NULL,
  `bank_transaction_end_time` timestamp NULL DEFAULT NULL,
  `bank_result_code` varchar(50) NULL DEFAULT NULL,
  `bank_result_msg` varchar(50) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pay_by_card_token_donations_status` (`status`),
  KEY `idx_pay_by_card_token_donations_amount` (`amount`),
  KEY `idx_pay_by_card_token_donations_order_number` (`order_number`),
  CONSTRAINT `fk_pay_by_card_token_donations_periodic_id` FOREIGN KEY (`periodic_id`) REFERENCES `periodic_donations` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

