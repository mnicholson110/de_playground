FROM postgres:latest
ADD ./init.sql /docker-entrypoint-initdb.d/
EXPOSE 5432
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["postgres"]
