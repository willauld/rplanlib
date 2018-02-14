#go build -ldflags "-X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash `git rev-parse HEAD`" myapp.go
Set-PSDebug -Trace 1
#Set-PSDebug -Trace 2
go build -ldflags "-X rplanlib.version.BuildTime `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X rplanlib.version.Version 1.5"

Set-PSDebug -Trace 0 