# Start from official MySQL image
From mysql:5.7.21

# Add go-api test user
ENV MYSQL_ALLOW_EMPTY_PASSWORD=true \
    MYSQL_DATABASE=gorm \
    MYSQL_USER=gorm \
    MYSQL_PASSWORD=gorm

# Add scripts for database scheme and developer account initialization
RUN mkdir -p /init

ADD initdb.sql /init/
ADD initdb.sh /docker-entrypoint-initdb.d/
ADD mysql.cnf /etc/mysql/conf.d/

EXPOSE 3306
