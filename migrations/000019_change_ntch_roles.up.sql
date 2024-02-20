-- add new roles
INSERT INTO `roles` (`key`, `name`, `name_en`) VALUES ('trailblazer_ntch_3', '開創者', 'trailblazer');
INSERT INTO `roles` (`key`, `name`, `name_en`) VALUES ('trailblazer_ntch_12', '開創者', 'trailblazer');

-- update roles.weight value
UPDATE `roles` SET `weight` = 14 WHERE (`key` = 'trailblazer_ntch_3');
UPDATE `roles` SET `weight` = 15 WHERE (`key` = 'trailblazer_ntch_12');
