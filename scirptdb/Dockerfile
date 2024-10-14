FROM mysql:8.0.33
ENV MYSQL_ROOT_PASSWORD=mysql@root
ENV MYSQL_DATABASE=edugo_v3
ENV LANG en_US.UTF-8
COPY ./db.sql /docker-entrypoint-initdb.d/
COPY ./insert.sql /docker-entrypoint-initdb.d/
COPY ./user.sql /docker-entrypoint-initdb.d/
EXPOSE 3306
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["mysqld"]
