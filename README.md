# sample-msg-broker
A toy example for message broker.

Available commands:

* `PUB <topic> <msg>`
* `SUB <topic>`


# Run the app

- Run server app


    go run sample/server/main.go
 

- Run client app (you can run multiple instances)


    go run sample/client/main.go
    
You will see three predefined outputs.

Then from client window you can publish new message: 

    PUB qwerty test msg
    
Also client can subscribe to new topic with command: 

    SUB test 

    
    
