#!/bin/sh

curl -i -X PUT "http://localhost:9090/employees/123/deactivate"

exit 0

curl -i -X POST "http://localhost:9090/employees" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com"
  }'
