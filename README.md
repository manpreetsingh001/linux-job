### Linux job worker

#Design doc

https://docs.google.com/document/d/1wb2-LtYrLH8adqrIHsLu1IruhHdA-E27gsJXjw8CSxY/edit 

### Run Tests

Run this in root directory: 
```
go test ./...
```

### Build Binary

```
go mod vendor   in root directory 
api:
go build -o worker-api cmd/api/main.go

client:
go build -o worker-client cmd/client/main.go
	
```

### RUN  server
```
./worker-api
```

### Client commands
```
MANPSIN4-M-C3HM:linux-job manpsin4$ ./worker-client start bash a.sh
Job 8404f495-474c-4508-9704-bbb68372887b is started

MANPSIN4-M-C3HM:linux-job manpsin4$ ./worker-client status  8404f495-474c-4508-9704-bbb68372887b 
Job ID: 8404f495-474c-4508-9704-bbb68372887b 
Job Status:         running
Start time: 2021-08-30 12:38:09.759216 -0700 PDT


./worker-client stream 8404f495-474c-4508-9704-bbb68372887b   
Mon Aug 30 13:41:21 PDT 2021
Mon Aug 30 13:41:22 PDT 2021
Mon Aug 30 13:41:23 PDT 2021
Mon Aug 30 13:41:24 PDT 2021
Mon Aug 30 13:41:25 PDT 2021
Mon Aug 30 13:41:26 PDT 2021
Mon Aug 30 13:41:27 PDT 2021
Mon Aug 30 13:41:28 PDT 2021
Mon Aug 30 13:41:29 PDT 2021
Mon Aug 30 13:41:30 PDT 2021

MANPSIN4-M-C3HM:linux-job manpsin4$ ./worker-client  stop  8404f495-474c-4508-9704-bbb68372887b 
Job 8404f495-474c-4508-9704-bbb68372887b  has been stopped



untrusted client

MANPSIN4-M-C3HM:linux-job manpsin4$ ./worker-client status f3b30538-61b6-45d9-b4f6-07e61a6696ae
rpc error: code = DeadlineExceeded desc = latest balancer error: connection error: desc = "transport: authentication handshake failed: x509: certificate signed by unknown authority"MANPSIN4-M-C3HM:linux-job manpsin4$ 

```

### What is left:

> Makefile
> Fix module in go mod to use github.com

