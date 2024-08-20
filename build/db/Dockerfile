FROM postgres:13.14-alpine

COPY ./build/db/init.sh /docker-entrypoint-initdb.d/

RUN chmod +x /docker-entrypoint-initdb.d/init.sh