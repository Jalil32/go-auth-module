### Running The Backend
cd backend
go run cmd/main.go

### Setup Environment
- Create a file '.env' in /backend/config/
- add environment variables

### Start Dev Server
```bash
bun run dev
```
### Start Dev Database
1. Ensure you have docker, docker compose and postgres installed
2. Start the docker daemon using
- sudo systemctl start docker
- sudo systemctl enable docker
2. cd into directory where compose.yml is located (deployments)
3. Run docker-compose up -d to spin up the pgAdmin and postgres containers
- Note that -d stands for "detatched mode", which will start the containers in the background and will not be attached to the current terminal.
4. Navigate to http://localhost:5050 to access the pgAdmin web interface
5. Log in using the email and password that you set in the PGADMIN_DEFAULT_EMAIL and PGADMIN_DEFAULT_PASSWORD environment variables in the docker compose file.
In the left-hand pane, expand the Servers node.
6. Right-click on the Servers node, and select Create-> Server… from the context menu.
In the Create — Server dialog that appears, enter a name for the server in the Name field.
In the Connection tab, enter the following information:
Host name/address: the hostname or IP address of the machine where the PostgreSQL database is running. If you are running the PostgreSQL container on the same machine as the pgAdmin container, you can use postgres as the hostname.
Port: the port number where the PostgreSQL database is listening for connections. In the Docker Compose file I provided earlier, the PostgreSQL container is exposing port 5432, so you can use 5432 as the port number.
Maintenance database: the name of the database that you want to use for maintenance tasks. You can use the postgres database for this purpose.
Username: the username that you want to use to connect to the database. You can use the POSTGRES_USER environment variable that you set in the Docker Compose file.
Password: the password for the user that you want to use to connect to the database. You can use the POSTGRES_PASSWORD environment variable that you set in the Docker Compose file.
