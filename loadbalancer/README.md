# Setup

`go build server.go`

`go build loadbalancer.go`

# Run

`./loadbalancer &`

`./server 3000 &`

`./server 3001 &`

`./server 3002 &`

`curl -i localhost:8080` multiple times

`curl -w "\n" localhost:8080 & curl -w "\n" localhost:8080 & curl -w "\n" localhost:8080` to run in parallel

`curl --header "UserID: 1" localhost:8080`

`curl -X POST -d '{"name":"node4","ref":"http://localhost:3003"}' localhost:8080/add`

`curl localhost:8080/report`

# Round Robin

```golang
lb := newLoadBalancer(&roundRobin{
	handlers: []*handler{
		newHandler("http://localhost:3000"),
		newHandler("http://localhost:3001"),
		newHandler("http://localhost:3002"),
	},
})

http.HandleFunc("/", lb.handleFunc)
log.Fatal(http.ListenAndServe(":8080", nil))
```

# Random

```golang
lb := newLoadBalancer(random{
	handlers: []*handler{
		newHandler("http://localhost:3000"),
		newHandler("http://localhost:3001"),
		newHandler("http://localhost:3002"),
	},
})

http.HandleFunc("/", lb.handleFunc)
log.Fatal(http.ListenAndServe(":8080", nil))
```

# Least Loaded

```golang
lb := newLoadBalancer(
	newLeastConn([]*handler{
		newHandler("http://localhost:3000"),
		newHandler("http://localhost:3001"),
		newHandler("http://localhost:3002"),
	}),
)

http.HandleFunc("/", lb.handleFunc)
log.Fatal(http.ListenAndServe(":8080", nil))
```

# Weighted Round Robin

```golang
lb := newLoadBalancer(
	newWeightedRoundRobin([]*handler{
		newHandler(&handlerOpts{
			ref:    "http://localhost:3000",
			weight: 3,
		}),
		newHandler(&handlerOpts{
			ref:    "http://localhost:3001",
			weight: 1,
		}),
		newHandler(&handlerOpts{
			ref:    "http://localhost:3002",
			weight: 1,
		}),
	}),
)

http.HandleFunc("/", lb.handleFunc)
log.Fatal(http.ListenAndServe(":8080", nil))
```

# Simple Hashing

```golang
lb := newLoadBalancer(simpleHashing{
  hs: []*handler{
    newHandler(&handlerOpts{ref: "http://localhost:3000"}),
    newHandler(&handlerOpts{ref: "http://localhost:3001"}),
    newHandler(&handlerOpts{ref: "http://localhost:3002"}),
  },
})

http.HandleFunc("/", lb.handleFunc)
log.Fatal(http.ListenAndServe(":8080", nil))
```

# Consistent Hashing

```golang
h1, err := newHandler(&handlerOpts{
  Name: "node1",
  Ref:  "http://localhost:3000",
})
if err != nil {
  log.Fatal(err)
}

h2, err := newHandler(&handlerOpts{
  Name: "node2",
  Ref:  "http://localhost:3001",
})
if err != nil {
  log.Fatal(err)
}

h3, err := newHandler(&handlerOpts{
  Name: "node3",
  Ref:  "http://localhost:3002",
})
if err != nil {
  log.Fatal(err)
}

lb := newLoadBalancer(newConsistentHashing(
  &consistentHashingOpts{
    handlers:    []*handler{h1, h2, h3},
    numReplicas: 2,
  },
))
```
