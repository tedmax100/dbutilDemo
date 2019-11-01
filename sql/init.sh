#!/bin/bash

set -e

mysql -u root -p${MYSQL_ROOT_PASSWORD} -h localhost \
-e "CREATE USER 'db_user'@'%' IDENTIFIED BY 'pwd';"

mysql -u root -p${MYSQL_ROOT_PASSWORD} -h localhost \
-e "GRANT ALL PRIVILEGES ON lottery.* TO 'db_user'@'%';"

mysql -u root -p${MYSQL_ROOT_PASSWORD} -h localhost \
-e "FLUSH PRIVILEGES;"
