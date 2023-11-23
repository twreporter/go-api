-- drop column
ALTER TABLE `users` DROP agree_data_collection;

-- drop tables
DROP TABLE IF EXISTS `users_posts_reading_count`;
DROP TABLE IF EXISTS `users_posts_reading_time`;
