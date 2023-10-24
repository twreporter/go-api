-- drop column
ALTER TABLE `users` DROP `source`;
ALTER TABLE `roles` DROP `weight`;
ALTER TABLE `users_roles` DROP `expired_at`;

-- remove added roles
DELETE FROM `roles` WHERE `key` = 'action_taker_ntch';
DELETE FROM `roles` WHERE `key` = 'trailblazer_ntch';
