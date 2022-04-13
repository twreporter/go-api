ALTER TABLE `users` MODIFY `age_range` enum('less_than_18', '18_to_24', '25_to_34', '35_to_44', '45_to_54', '55_to_64', 'above_65') DEFAULT NULL;
ALTER TABLE `pay_by_prime_donations` MODIFY `cardholder_age_range` enum('less_than_18', '18_to_24', '25_to_34', '35_to_44', '45_to_54', '55_to_64', 'above_65') DEFAULT NULL;
ALTER TABLE `periodic_donations` MODIFY `cardholder_age_range` enum('less_than_18', '18_to_24', '25_to_34', '35_to_44', '45_to_54', '55_to_64', 'above_65') DEFAULT NULL;
