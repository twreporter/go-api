-- add should_merge_offline_donation_by_identity
ALTER TABLE `users`
ADD COLUMN `should_merge_offline_donation_by_identity` BOOLEAN NOT NULL DEFAULT false;

-- drop `is_showofflinedonation`
ALTER TABLE `users`
DROP COLUMN `is_showofflinedonation`;