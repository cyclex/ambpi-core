# AMBPI Day Core Services

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Configuration](#configuration)

## Features

There is 2 command in this core services
1. server
    supporting rest endpoint for CMS
2. webhook
    supporting webhook for WABA
3. cms
    provide content management system for campaign AMBPI

## Installation


```bash
# Clone the repository
git clone https://github.com/cyclex/ambpi-core.git

# Navigate to the project directory
cd ambpi-core

# Install dependencies
make build

# Execute sql script
source schema/master.sql
```

## Usage

```bash
# Run the application
# Server apps
./engine server -p :8081 -c config.json -d true

NAME:
   server server - start cms and chatbot service

USAGE:
   server server [command options] [arguments...]

OPTIONS:
   --port value, -p value    Listen to port
   --config value, -c value  Load configuration file
   --debug, -d               Debug mode (default: false)
   --help, -h                show help

# webhook apps
./engine webhook -p :8082 -c config.json -d true

NAME:
   server webhook - start webhook

USAGE:
   server webhook [command options] [arguments...]

OPTIONS:
   --port value, -p value    Listen to port
   --config value, -c value  Load configuration file
   --debug, -d               Debug mode (default: false)
   --help, -h                show help
```

## Configuration

```json
// config.json
{
  "log": {
    "maxsize": 10,
    "maxbackups": 10
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "user",
    "password": "password",
    "name": "ambpi"
  },
  "chatbot": {
    "host": "https://869de8ac-2cac-44eb-ab04-fb7ff83cb000.mock.pstmn.io/bot/webhook",
    "host_push": "https://869de8ac-2cac-44eb-ab04-fb7ff83cb000.mock.pstmn.io",
    "phone_id":"6282122912548",
    "account_id": "xxx",
    "division_id": "xxx",
    "waba_account_number":"6282311333723",
    "access_token":"xxx",
    "access_token_push":"12345"
  },
  "queue":{
    "host" : "localhost",
    "port" : 27017,
    "name" : "ambpi",
    "expired" : 168 // hours in a week
  },
  "download_folder":"/Users/ansharharyadi/Desktop"
}

```
