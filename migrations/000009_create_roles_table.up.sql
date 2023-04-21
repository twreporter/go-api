CREATE TABLE IF NOT EXISTS `roles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `name_en` varchar(50) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `roles` (`name`, `name_en`, `created_at`, `updated_at`) VALUES ('探索者', 'explorer', now(), now());
INSERT INTO `roles` (`name`, `name_en`, `created_at`, `updated_at`) VALUES ('行動者', 'action taker', now(), now());
INSERT INTO `roles` (`name`, `name_en`, `created_at`, `updated_at`) VALUES ('開創者', 'trailblazer', now(), now());
