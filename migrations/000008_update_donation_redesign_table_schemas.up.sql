ALTER TABLE `pay_by_prime_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no', 'no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');
UPDATE `pay_by_prime_donations`
SET `send_receipt` = CASE `send_receipt`
                     WHEN 'no' THEN 'no_receipt'
                     WHEN 'yearly' THEN 'paperback_receipt_by_year'
                     WHEN 'monthly' THEN 'paperback_receipt_by_month'
                   END;
ALTER TABLE `pay_by_prime_donations` MODIFY `send_receipt` enum('no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');

ALTER TABLE `periodic_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no', 'no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');
UPDATE `periodic_donations`
SET `send_receipt` = CASE `send_receipt`
                     WHEN 'no' THEN 'no_receipt'
                     WHEN 'yearly' THEN 'paperback_receipt_by_year'
                     WHEN 'monthly' THEN 'paperback_receipt_by_month'
                   END;
ALTER TABLE `periodic_donations` MODIFY `send_receipt` enum('no_receipt', 'digital_receipt_by_month', 'digital_receipt_by_year', 'paperback_receipt_by_month', 'paperback_receipt_by_year');
