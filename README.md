# Go NetHunter

```text
This was my "Hello World" project on GOLANG.
This program is created for educational purposes only.
The developer is not responsible for its usage.
Scanning any network without permission may be considered illegal in your country. 
Users are solely responsible for their actions.
```

[You can watch the video on Twitter](https://twitter.com/zoBkaruB/status/1708356471983858082/)

### Build
```bash
go build gonethunter.go
```

### Build macOS
```bash
export GOOS="darwin"; export GOARCH="amd64"; go build gonethunter.go
```

### Sample usage
```bash
./gonethunter -stopOnFound -input "iplist.txt" -output "found.txt" -hostname "google.com" -search "www.google.com" -threads 300
```

### Run from source
```bash
go run gonethunter.go -h
```

### Cross Compile
```bash
rm -f bin/gonethunter*
export GOOS="darwin"; export GOARCH="amd64"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH"" gonethunter.go
export GOOS="darwin"; export GOARCH="arm64"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH"" gonethunter.go
export GOOS="windows"; export GOARCH="amd64"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH".exe" gonethunter.go
export GOOS="windows"; export GOARCH="386"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH".exe" gonethunter.go
export GOOS="linux"; export GOARCH="386"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH"" gonethunter.go
export GOOS="linux"; export GOARCH="amd64"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH"" gonethunter.go
export GOOS="linux"; export GOARCH="arm64"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH"" gonethunter.go
export GOOS="linux"; export GOARCH="arm"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH"" gonethunter.go
#export GOOS="android"; export GOARCH="arm"; go build -o "bin/gonethunter_"$GOOS"_"$GOARCH"" gonethunter.go
```

