ALTER TABLE `pay_by_prime_donations`
ADD `linepay_method` enum('CREDIT_CARD', 'BALANCE', 'POINT') DEFAULT NULL,
ADD `linepay_point` int DEFAULT NULL AFTER `is_anonymous`;
