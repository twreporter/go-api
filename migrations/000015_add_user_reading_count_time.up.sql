-- add new columns
ALTER TABLE `users` ADD `agree_data_collection` tinyint(1) DEFAULT 1;

-- add new tables
CREATE TABLE IF NOT EXISTS `users_posts_reading_count` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `post_id` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_users_read_posts_users1_idx` (`user_id`),
  CONSTRAINT `fk_users_read_posts_users1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `users_posts_reading_time` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `seconds` int(10) NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `post_id` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_users_posts_seconds_users1_idx` (`user_id`),
  CONSTRAINT `fk_users_posts_seconds_users1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
