#!/bin/bash

docker exec -it auth-service-postgres-1\
  bash -c "psql -U rat -d auth_service -c \"delete from auth_whitelist where created_at <= '$CUT_TOKENS_FROM'\";"
