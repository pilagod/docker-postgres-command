ARG POSTGRES

FROM postgres:$POSTGRES-alpine
COPY ./bin /bin
