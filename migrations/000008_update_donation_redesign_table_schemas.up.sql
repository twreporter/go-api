ALTER TABLE `user` MODIFY `gender` enum(male, female, other, unreveal);
ALTER TABLE `pay_by_prime_donations` MODIFY `send_receipt` enum(no_receipt, digital_receipt_by_month, digital_receipt_by_year, paperback_receipt_by_month, paperback_receipt_by year);
ALTER TABLE `periodic_donations` MODIFY `send_receipt` enum(no_receipt, digital_receipt_by_month, digital_receipt_by_year, paperback_receipt_by_month, paperback_receipt_by year);
