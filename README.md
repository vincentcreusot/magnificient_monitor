# magnificient_monitor
Monitors Magnificient service and gives some status about its health
## Dependencies
The program depends on 
- Gorrilla mux http://github.com/gorilla/mux
## Running the program
run
```bash
go run main.go
```
## REST API
The program exposes 2 get methods :
- / which returns a json version of the status of th magnificent service
- /callit which calls the service and updates the statuses
- /muststop stops the call of the magnificent service in background