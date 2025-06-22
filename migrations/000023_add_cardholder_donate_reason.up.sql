ALTER TABLE `pay_by_prime_donations`
ADD COLUMN `cardholder_donate_reason` varchar(191) DEFAULT NULL;

ALTER TABLE `pay_by_card_token_donations`
ADD COLUMN `cardholder_donate_reason` varchar(191) DEFAULT NULL;