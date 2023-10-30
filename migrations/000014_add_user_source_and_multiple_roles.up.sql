-- add new columns
ALTER TABLE `users` ADD `source` set('ntch') DEFAULT NULL;
ALTER TABLE `roles` ADD `weight` int(10) DEFAULT NULL;
ALTER TABLE `users_roles` ADD `expired_at` timestamp NULL DEFAULT NULL;

-- add new roles
INSERT INTO `roles` (`key`, `name`, `name_en`) VALUES ('action_taker_ntch', '行動者', 'action taker');
INSERT INTO `roles` (`key`, `name`, `name_en`) VALUES ('trailblazer_ntch', '開創者', 'trailblazer');

-- update roles.weight value
UPDATE `roles` SET `weight` = 1 WHERE (`key` = 'explorer');
UPDATE `roles` SET `weight` = 5 WHERE (`key` = 'action_taker_ntch');
UPDATE `roles` SET `weight` = 9 WHERE (`key` = 'action_taker');
UPDATE `roles` SET `weight` = 13 WHERE (`key` = 'trailblazer_ntch');
UPDATE `roles` SET `weight` = 17 WHERE (`key` = 'trailblazer');
