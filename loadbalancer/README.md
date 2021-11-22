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
