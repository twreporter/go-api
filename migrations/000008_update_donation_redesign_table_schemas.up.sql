ALTER TABLE `pay_by_prime_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no', 'no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');
UPDATE `pay_by_prime_donations`
SET `send_receipt` = CASE `send_receipt`
                       WHEN 'no' THEN 'no_receipt'
                       WHEN 'yearly' THEN 'paperback_receipt_by_year'
                       WHEN 'monthly' THEN 'paperback_receipt_by_month'
                     END,
`cardholder_words_for_twreporter` = `notes`,
`cardholder_security_id` = `cardholder_national_id`,
`cardholder_address_zip_code` = `cardholder_zip_code`,
`cardholder_address_detail` = `cardholder_address`;

ALTER TABLE `periodic_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no', 'no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');
UPDATE `periodic_donations`
SET `send_receipt` = CASE `send_receipt`
                       WHEN 'no' THEN 'no_receipt'
                       WHEN 'yearly' THEN 'paperback_receipt_by_year'
                       WHEN 'monthly' THEN 'paperback_receipt_by_month'
                     END,
`cardholder_words_for_twreporter` = `notes`,
`cardholder_security_id` = `cardholder_national_id`,
`cardholder_address_zip_code` = `cardholder_zip_code`,
`cardholder_address_detail` = `cardholder_address`;
