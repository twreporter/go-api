ALTER TABLE `user` MODIFY `gender` varchar(2);
ALTER TABLE `pay_by_prime_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no');
ALTER TABLE `periodic_donations` MODIFY `send_receipt` enum('yearly', 'monthly', 'no');
