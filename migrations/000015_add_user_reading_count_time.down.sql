-- drop column
ALTER TABLE `users` DROP `agree_data_collection`;
ALTER TABLE `users` DROP `read_posts_count`;
ALTER TABLE `users` DROP `read_posts_sec`;

-- drop tables
DROP TABLE IF EXISTS `users_posts_reading_counts`;
DROP TABLE IF EXISTS `users_posts_reading_times`;
