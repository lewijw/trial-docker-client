# trial-docker-client

Uses the Go docker client to start a postgres database with a mapped port.

This works as expected on RHEL 7.7/Docker Community 20.10.21/Go 1.18.7 
but not on Windows 11/Dockers Desktop 20.10.23/Go 1.20.2.  It also does not work in WSL version 2.

The port mapping does not work.  For example on Linux you can see :::35093->5432/tcp:

```
go run main.go
Go version: go1.19.2
{"status":"Pulling from library/postgres","id":"13.10"}
{"status":"Digest: sha256:a06b381f1ed083cf85900b4814936b6c85a93e820c9924a0a9d622054ce353b9"}
{"status":"Status: Image is up to date for postgres:13.10"}
2023/03/16 11:37:07 free port: [::]:35093
2023/03/16 11:37:07 Port: 5432/tcp; Bindings: [{:: 35093}]
fsnlax05.unx.sas.com> docker ps | grep database
e2d11d974259   postgres:13.10                               "docker-entrypoint.sâ¦"   11 seconds ago   Up 10 seconds   :::35093->5432/tcp                                                                            database
```

However, on Windows (Git Bash terminal), the port is not mapped to 64215 as desired:

```
$ go run main.go
Go version: go1.20.2
{"status":"Pulling from library/postgres","id":"13.10"}
{"status":"Digest: sha256:a06b381f1ed083cf85900b4814936b6c85a93e820c9924a0a9d622054ce353b9"}
{"status":"Status: Image is up to date for postgres:13.10"}
2023/03/16 11:39:31 free port: [::]:64215
2023/03/16 11:39:31 Port: 5432/tcp; Bindings: [{:: 64215}]
$ docker ps | grep database
14c5e7719550   postgres:13.10   "docker-entrypoint.s…"   52 seconds ago   Up 51 seconds   5432/tcp   database
```

You can however do it at the command line.  Note the 0.0.0.0:64215->5432/tcp:

```
$ docker run -d -p 64215:5432 --name database -e POSTGRES_USER='myuser' -e POSTGRES_PASSWORD='mypassword' -e POSTGRES_DB='mydb' postgres:13.10

6a942b881de1cb7a66652926bbdfa199376d9c239983d37b2710da108dc98859

$ docker ps | grep database
6a942b881de1   postgres:13.10   "docker-entrypoint.s…"   13 seconds ago   Up 12 seconds   0.0.0.0:64215->5432/tcp   database
```
