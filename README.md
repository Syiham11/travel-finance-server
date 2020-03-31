[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/matthewandrewpalmer/travel-finance-server) 

# Travel Finance Server

This is a simple Go application that connects to a MySQL Database, to handle requests for the  [React frontend](https://github.com/matthewandrewpalmer/travel-finance-website)

## Setup

To setup this project you'll need a MySQL database setup with the schema setup like in `db_setup.sql`.

Then you need to create a '.env' file with the DB credentials and the port you want the server setup. Setup like below
```.env
PORT=5000

DB_USERNAME=root
DB_PASSWORD=password
DB_NAME=travel
DB_URL=localhost
DB_PORT=3306
```

### API Guide

**Return all rail journeys**
```
/rail-journeys
```