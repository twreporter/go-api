ALTER TABLE `pay_by_prime_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no', 'no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');
UPDATE `pay_by_prime_donations`
SET `send_receipt` = CASE `send_receipt`
                     WHEN 'no_receipt' THEN 'no'
                     WHEN 'paperback_receipt_by_year' THEN 'yearly'
                     WHEN 'paperback_receipt_by_month' THEN 'monthly'
                   END
ALTER TABLE `pay_by_prime_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no');

ALTER TABLE `periodic_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no', 'no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');
UPDATE `periodic_donations`
SET `send_receipt` = CASE `send_receipt`
                     WHEN 'no_receipt' THEN 'no'
                     WHEN 'paperback_receipt_by_year' THEN 'yearly'
                     WHEN 'paperback_receipt_by_month' THEN 'monthly'
                   END
ALTER TABLE `periodic_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no');
