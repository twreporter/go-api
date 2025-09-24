--
-- Table structure for table `jobs_mailchimp`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;

CREATE TABLE IF NOT EXISTS `jobs_mailchimp` (
  `id` int (10) unsigned NOT NULL AUTO_INCREMENT,
  `receiver` varchar(255) NOT NULL COMMENT 'email address of the subscriber',
  `interests` varchar(255) NOT NULL COMMENT 'JSON array for interest IDs',
  `state` varchar(10) NOT NULL DEFAULT 'new' COMMENT 'new / processing / fail',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users_mailgroup`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;

CREATE TABLE `users_mailgroups` (
  `user_id` int(10) unsigned NOT NULL,
  `mailgroup_id` varchar(255) NOT NULL COMMENT 'interest ID from MailChimp',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`,`mailgroup_id`),
  CONSTRAINT `fk_users_mailgroup_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

/*!40101 SET character_set_client = @saved_cs_client */;