This is a toy project to experiment with k8s services and gRPC.

# Setup

## Kind
```
# On macOS
brew install kind

# On Ubuntu
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.18.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

create cluster
```
kind create cluster --name my-cluster
```

## Docker
On macOS: Download Docker Desktop from the Docker website.

On Ubuntu:
```
sudo apt update
sudo apt install docker.io
sudo systemctl start docker
sudo systemctl enable docker
```

## Ko

```
go install github.com/google/ko@latest
```

## Proto

```
# On macOS
brew install protoc

# On Ubuntu
sudo apt update
sudo apt install -y protobuf-compiler
```

install Go proto plugin
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Add go to path
```
export PATH="$PATH:$(go env GOPATH)/bin"
```

# Experiment
## Deploy Headless Service Example
This will deploy servers, the client, and the service with the headless model.

```
make deploy-headless
```

## Deploy ClusterIP Service Example
This will deploy servers, the client, and the service with the clusterIp model.
```
make deploy-cluster-ip
```

## Undeploy Everything
This deletes servers, the client, and services.
```
make undeploy-all
```




