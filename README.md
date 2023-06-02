# timestamp_mcast

Multicast timestamp every second.  
Compiled versions can be downloaded from releases.

##### Compile for Linux
```shell
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client/server
```

##### Compile for Windows
```shell
$ GOOS=windows GOARCH=amd64 go build -o client.exe/server.exe
```