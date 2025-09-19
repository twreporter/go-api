ALTER TABLE `users`
ADD COLUMN `is_showofflinedonation` BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE `users`
DROP COLUMN `should_merge_offline_donation_by_identity`;