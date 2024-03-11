### Run project
```
go run ./cmd/main.go

```

### Gen swagger

```
swag init -g ./api/http.go

```

### Install packages (LINUX)

1. Install C compiler

```
export CGO_ENABLED=1
apt-get install build-essential
```

2. Install IBM MQ
```
export genmqpkg_incnls=1
export genmqpkg_incsdk=1
export genmqpkg_inctls=1
sudo mkdir -p /opt/mqm
cd /opt/mqm
sudo curl -LO "https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist/9.3.4.1-IBM-MQC-Redist-LinuxX64.tar.gz"
sudo tar -zxf ./*.tar.gz
sudo rm -f ./*.tar.gz
```

3. Install XSLT package
```
sudo apt install libxml2-dev libxslt1-dev liblzma-dev zlib1g-dev -y
```