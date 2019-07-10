SET FOREIGN_KEY_CHECKS=0;

DROP TABLE IF EXISTS `bookmarks`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `o_auth_accounts`;
DROP TABLE IF EXISTS `services`;
DROP TABLE IF EXISTS `registrations`;
DROP TABLE IF EXISTS `reporter_accounts`;
DROP TABLE IF EXISTS `users_bookmarks`;
DROP TABLE IF EXISTS `web_push_subs`;
DROP TABLE IF EXISTS `pay_by_prime_donations`;
DROP TABLE IF EXISTS `pay_by_other_method_donations`;
DROP TABLE IF EXISTS `periodic_donations`;
DROP TABLE IF EXISTS `pay_by_card_token_donations`;

SET FOREIGN_KEY_CHECKS=1;
