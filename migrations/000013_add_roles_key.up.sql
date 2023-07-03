ALTER TABLE `roles`
ADD COLUMN `key` varchar(50) NOT NULL AFTER `name_en`;

UPDATE `roles` SET `key` = 'explorer' WHERE (`name_en` = 'explorer');
UPDATE `roles` SET `key` = 'action_taker' WHERE (`name_en` = 'action taker');
UPDATE `roles` SET `key` = 'trailblazer' WHERE (`name_en` = 'trailblazer');