#!/bin/bash

#create user
curl -X PUT \
  http://localhost:8000/api/user \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
   "email" : "jdoe@email.com",
   "firstname" : "John",
   "lastname" : "Doe",
   "password" : "a"
}'

#get user info
curl -X GET \
  http://localhost:8000/api/user \
  -H 'Authorization: Basic amRvZUBlbWFpbC5jb206YQ==' \
  -H 'cache-control: no-cache'

#create team
curl -X PUT \
  'http://localhost:8000/api/team?email=pvojacek@foxconndrc.com' \
  -H 'Authorization: Basic amRvZUBlbWFpbC5jb206YQ==' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"name" : "new_team",
	"Description" : "brand new team"
}'

#create team
curl -X PUT \
  'http://localhost:8000/api/team?email=pvojacek@foxconndrc.com' \
  -H 'Authorization: Basic amRvZUBlbWFpbC5jb206YQ==' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"name" : "new_team_1",
	"Description" : "brand new team"
}'

#create team
curl -X PUT \
  'http://localhost:8000/api/team?email=pvojacek@foxconndrc.com' \
  -H 'Authorization: Basic amRvZUBlbWFpbC5jb206YQ==' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"name" : "new_team_2",
	"Description" : "brand new team"
}'

#create team
curl -X PUT \
  'http://localhost:8000/api/team?email=pvojacek@foxconndrc.com' \
  -H 'Authorization: Basic amRvZUBlbWFpbC5jb206YQ==' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"name" : "new_team_3",
	"Description" : "brand new team"
}'


#subscribe to team
curl -X POST \
  http://localhost:8000/api/team/new_team/subscribe \
  -H 'Authorization: Basic amRvZUBlbWFpbC5jb206YQ==' \
  -H 'cache-control: no-cache'
