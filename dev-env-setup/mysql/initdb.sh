#!/bin/bash

mysql -u root -e "CREATE DATABASE IF NOT EXISTS test_membership"

mysql -u root test_membership < /init/initdb.sql

mysql -u root -e "CREATE USER 'test_membership'@'%' IDENTIFIED BY 'test_membership'"

mysql -u root -e "GRANT ALL PRIVILEGES ON test_membership.* TO 'test_membership'@'%'"
