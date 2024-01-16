CREATE TABLE IF NOT EXISTS `users_posts_reading_footprints` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `post_id` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_users_read_posts_users3_idx` (`user_id`),
  CONSTRAINT `fk_users_read_posts_users3` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
