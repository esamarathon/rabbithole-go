# Rabbithole
RabbitMQ to Postgres sink. Stores messages in a queriable jsonb column.

<img src="https://raw.githubusercontent.com/esamarathon/rabbithole-go/master/Rabbithole.png" alt="example visualization used by ESA" />

## Usage
For regular use, please use the prebuilt Docker container, found on [Docker Hub](https://hub.docker.com/r/esamarathon/rabbithole:2)

````sh
$ docker pull esamarathon/rabbithole:2
````

But please make not of the configuration necessary to get it started. 
I recommend using something like Docker-Compose to avoid having to type it all in manually on the commandline.
For more information on this, see below.


## Background ##
We needed a way to store a log of the RabbitMQ events sent during our events for graphing and evaluation purposes.
So I wrote this program.
It connects to a RabbitMQ server and a PostgreSQL server and stores the messages in there.
The messages stored are configurable using `/app/appsettings.json` or environment variables.

## Configuration ##
Application comes with the following default configuration (found in `/app/appsettings.default.json`)
````json
{
    "Logging": {
        "LogLevel": {
            "Default": "Information",
            "Microsoft": "Error"
        },
        "Console": {
            "IncludeScopes": true
        }
    },
    "ConnectionStrings": {
        "Events": "User ID=postgres;Password=password;Host=localhost;Port=5432;Database=rabbithole;"
    },
    "RabbitMQ": {
        "HostName": "localhost",
        "UserName": "guest",
        "Password": "guest",
        "VirtualHost": "esa_dev",
        "Bindings": [
            {
                "Exchange": "demo1",
                "Topic": "#"
            }
        ]
    }
}
````
This showcases all possible settings and works fine for most development work.
To override any settings, simply add a appsettings.json file with the overriding settings and it will load during startup.
Additionally, using the rules provided at https://docs.microsoft.com/en-us/aspnet/core/fundamentals/configuration/?view=aspnetcore-2.2#environment-variables-configuration-provider, 
you can also override settings using environment variables prefixed with RABBITHOLE_.
Example docker-compose.yml file:

````yaml
version: "3.7"
services:
  app:
    environment:
      - RABBITHOLE_Logging__LogLevel__Default=Debug
      - RABBITHOLE_RabbitMQ__HostName=rabbit
      - RABBITHOLE_ConnectionStrings__Events=User ID=postgres;Password=password;Host=db;Port=5432;Database=rabbithole
````

Please make note of the double underscores between keys.
