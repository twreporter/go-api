UPDATE `pay_by_prime_donations`
SET `send_receipt` = CASE `send_receipt`
                       WHEN 'no_receipt' THEN 'no'
                       WHEN 'paperback_receipt_by_year' THEN 'yearly'
                       WHEN 'paperback_receipt_by_month' THEN 'monthly'
                     END,
`notes` = `cardholder_words_for_twreporter`,
`cardholder_national_id` = `cardholder_security_id`,
`cardholder_zip_code` = `cardholder_address_zip_code`,
`cardholder_address` = CONCAT(`cardholder_address_country`, `cardholder_address_state`, `cardholder_address_city`, `cardholder_address_detail`);
ALTER TABLE `pay_by_prime_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no') DEFAULT 'yearly';

UPDATE `periodic_donations`
SET `send_receipt` = CASE `send_receipt`
                       WHEN 'no_receipt' THEN 'no'
                       WHEN 'paperback_receipt_by_year' THEN 'yearly'
                       WHEN 'paperback_receipt_by_month' THEN 'monthly'
                     END,
`notes` = `cardholder_words_for_twreporter`,
`cardholder_national_id` = `cardholder_security_id`,
`cardholder_zip_code` = `cardholder_address_zip_code`,
`cardholder_address` = CONCAT(`cardholder_address_country`, `cardholder_address_state`, `cardholder_address_city`, `cardholder_address_detail`);
ALTER TABLE `periodic_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no') DEFAULT 'yearly';
