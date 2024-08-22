#!/bin/bash

docker run --rm -v $(pwd)/db/migrations:/migrations\
  --network $APP_NETWORK\
      migrate/migrate\
      -path migrations\
      -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@postgres:5432/$POSTGRES_DB?sslmode=disable"\
      up
