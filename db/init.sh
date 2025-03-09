#!/bin/sh

echo "Creating tables if not exists..."

# Register
echo "Registering users..."
sleep 1
REGISTER_RESPONSE_1=$(curl -s -X POST http://$HOST:$EXPOSE_PORT/register -H "Content-Type: application/json" -d '{"Username": "user11", "Password": "passWord1234"}')
USER_1=$(echo $REGISTER_RESPONSE_1 | jq -r '.Id')
LOGIN_RESPONSE_1=$(curl -s -X POST http://$HOST:$EXPOSE_PORT/login -H "Content-Type: application/json" -d "{\"Username\": \"user11\", \"Password\": \"passWord1234\"}")
TOKEN_1=$(echo $LOGIN_RESPONSE_1 | jq -r '.token')

REGISTER_RESPONSE_2=$(curl -s -X POST http://$HOST:$EXPOSE_PORT/register -H "Content-Type: application/json" -d '{"Username": "user22", "Password": "passWord12345"}')
USER_2=$(echo $REGISTER_RESPONSE_2 | jq -r '.Id')
LOGIN_RESPONSE_2=$(curl -s -X POST http://$HOST:$EXPOSE_PORT/login -H "Content-Type: application/json" -d "{\"Username\": \"user22\", \"Password\": \"passWord12345\"}")
TOKEN_2=$(echo $LOGIN_RESPONSE_2 | jq -r '.token')

echo "Registered users ID: {$USER_1, $USER_2}"

# Creating some films...
echo "Posting films..."
curl -X POST http://$HOST:$EXPOSE_PORT/film -H "Content-Type: application/json" -H "Authorization: $TOKEN_1" -d '{"Title": "Movie 1111", "Director": "Director 1", "Release": "2020-01-01T00:00:00Z"}'
curl -X POST http://$HOST:$EXPOSE_PORT/film -H "Content-Type: application/json" -H "Authorization: $TOKEN_1" -d '{"Title": "Movie 2222", "Director": "Director 2", "Release": "2020-01-01T00:00:00Z"}'
curl -X POST http://$HOST:$EXPOSE_PORT/film -H "Content-Type: application/json" -H "Authorization: $TOKEN_1" -d '{"Title": "Movie 8888", "Director": "Director 5", "Release": "2020-01-01T00:00:00Z"}'
curl -X POST http://$HOST:$EXPOSE_PORT/film -H "Content-Type: application/json" -H "Authorization: $TOKEN_1" -d '{"Title": "Movie 7777", "Director": "Director 3", "Release": "2020-01-01T00:00:00Z"}'

curl -X POST http://$HOST:$EXPOSE_PORT/film -H "Content-Type: application/json" -H "Authorization: $TOKEN_2" -d '{"Title": "Movie 3333", "Director": "Director 3", "Release": "2020-01-01T00:00:00Z"}'
curl -X POST http://$HOST:$EXPOSE_PORT/film -H "Content-Type: application/json" -H "Authorization: $TOKEN_2" -d '{"Title": "Movie 4444", "Director": "Director 2", "Release": "2020-01-01T00:00:00Z"}'
curl -X POST http://$HOST:$EXPOSE_PORT/film -H "Content-Type: application/json" -H "Authorization: $TOKEN_2" -d '{"Title": "Movie 6666", "Director": "Director 3", "Release": "2020-01-01T00:00:00Z"}'

echo "Setup completed!"