# app  
A simple web application which returns a random uid each time it recieves a get request on /uid endpoint.  

## endpoints  
- `/uid`: Whenever you send a get request to this endpoint it returns a random uid and store it in a mysql database.  
- `health`: An endpoint to check the health of application.  

## Vars  
Set these variables for the application:  

```yaml
Port: # The port that the application listens on
```
## config  

This application loads a config from `/config/config.json` and reads the database information:

```json
    {
        "database": {
            "host": "database host",
            "port": "database port",
            "user": "database user",
            "password": "database user",
            "dbname": "db name"
        }
    }
```
