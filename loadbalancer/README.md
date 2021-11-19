# Setup

`go build server.go`

`go build loadbalancer.go`

# Run

`./loadbalancer &`

`./server 3000 &`

`./server 3001 &`

`./server 3002 &`

`curl -i localhost:8080` multiple times
