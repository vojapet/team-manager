#!/bin/bash

HOST=${1:-"localhost:8000"}
echo "Using host '${HOST}'"

#create users
for i in {1..500}
do
curl -X PUT \
  http://${HOST}/api/user \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
   "email" : "jdoe_'$i'@email.com",
   "firstname" : "John",
   "lastname" : "Doe",
   "password" : "a"
}'
done

#get user info
curl -X GET \
  http://${HOST}/api/user \
  -H 'Authorization: Basic amRvZV8xQGVtYWlsLmNvbTph' \
  -H 'cache-control: no-cache'

#create teams
for i in {1..500}
do
curl -X PUT \
  "http://${HOST}/api/team" \
  -H 'Authorization: Basic amRvZV8xQGVtYWlsLmNvbTph' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"name" : "new_team_'${i}'",
	"Description" : "brand new team"
}'
done

#subscribe to team (jdoe_1)
for i in {1..200}
do
curl -X POST \
  http://${HOST}/api/team/new_team_${i}/subscribe \
  -H 'Authorization: Basic amRvZV8xQGVtYWlsLmNvbTph' \
  -H 'cache-control: no-cache'
done

#subscribe to team (jdoe_2)
for i in {100..300}
do
curl -X POST \
  http://${HOST}/api/team/new_team_${i}/subscribe \
  -H 'Authorization: Basic amRvZV8yQGVtYWlsLmNvbTph' \
  -H 'cache-control: no-cache'
done

#subscribe to team (jdoe_3)
for i in {1..400}
do
curl -X POST \
  http://${HOST}/api/team/new_team_${i}/subscribe \
  -H 'Authorization: Basic amRvZV8zQGVtYWlsLmNvbTph' \
  -H 'cache-control: no-cache'
done

// Comment
