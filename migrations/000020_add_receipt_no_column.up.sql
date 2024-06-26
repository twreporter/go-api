ALTER TABLE `pay_by_prime_donations`
ADD COLUMN `receipt_number` varchar(13) DEFAULT NULL;

ALTER TABLE `pay_by_card_token_donations`
ADD COLUMN `receipt_number` varchar(13) DEFAULT NULL;