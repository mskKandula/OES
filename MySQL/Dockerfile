# start with base image
FROM mysql

# import data into container

# All scripts in docker-entrypoint-initdb.d/ are automatically executed during container startup
COPY ./MySQL/*.sql /docker-entrypoint-initdb.d/

#// mysql -h 127.0.0.1 -P 3306 -u root -p OES