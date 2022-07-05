#!/bin/sh

USERNAME=csamuro
PASSWORD=1234
HOST=localhost
PORT=5432
DB=csamuro

echo "USERNAME=$USERNAME"
export TEC_DOC_POSTGRES_USERNAME=$USERNAME
echo "PASSWORD=$PASSWORD"
export TEC_DOC_POSTGRES_PASSWORD=$PASSWORD
echo "HOST=$HOST"
export TEC_DOC_POSTGRES_HOST=$HOST
echo "PORT=$PORT"
export TEC_DOC_POSTGRES_PORT=$PORT
echo "DB=$DB"
export TEC_DOC_POSTGRES_DB=$DB

export TEC_DOC_TECDOC_URL="https://webservice.tecalliance.services/pegasus-3-0/services/TecdocToCatDLB.jsonEndpoint"
export  TECDOC_API_KEY=
