#!/usr/bin/env bash

docker container prune && docker run -v ./DbSchema:/docker-entrypoint-initdb.d -e POSTGRES_MAX_CONNECTIONS=20000 -e POSTGRES_PASSWORD=mystrongpassword -e POSTGRES_USER=admin -e POSTGRES_DB=rinha -p 5432:5432 -d --name=postgres-dev postgres:16.2-bullseye
