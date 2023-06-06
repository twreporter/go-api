ALTER TABLE `users`
ADD COLUMN `activated` TIMESTAMP NULL DEFAULT NULL COMMENT 'null - deactivated; timestamp - activated;';
