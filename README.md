# team-manager
The RestAPI that allows you to manage the teams and their members

## Server

### API endpoints

Create new user  
*/api/user - PUT (input: User)*


retrieve user info (http auth)  
*/api/user - GET (output: User)*


Updates user info (http auth)  
*/api/user - POST (Input: User)*


Retrieve list of all teams (http auth)  
*/api/team - GET (output: List of teams)*


Creates new team (http auth)  
*/api/team - PUT (input: Team)*


Retrieve team info (http auth)  
*/api/team/{team_name} - GET (output: Team)*


User will be added to members of team (http auth)  
*/api/team/{team_name}/subscribe - POST*


User will be removed from members of team (http auth)  
*/api/team/{team_name}/unsubscribe - POST*

### Data description

#### User data structure

```
{
    "email": "",
    "firstname": "",
    "secondname": "",
    "password": ""
}
```

#### Team data structure

```
{
    "name": "",
    "description": ""
}
```

## Docker

### Server

To build server image:

```
# docker build -t team-server .
```

To run sever image (using default port 8000):

```
# docker run -p 8000:8000 team-server
```

### Client

To build 'client' image:

```
docker build -f Dockerfile.client -t client .
```

To run  'client' image (the login and password is taken from data in *server_test.sh* ):

```
docker run -it -e TEAM_URL='<add-server-url>' -e TEAM_LOGIN='jdoe_1@email.com' -e TEAM_PASSWORD='a' client
```