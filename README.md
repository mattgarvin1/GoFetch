# (todo) Fetch Points Server

What is this and what does it do?

## (todo) Note To Interviewer

- copy of the problem statement
- Please see my raw scratch note for thoughts I had along the way of solving the problem
- Basic explanation of my solution
- How complete is it
- If I wanted to spend more time on it, I would implement X Y Z
- thank you for taking the time to review my work!

## (todo) Running the Server

The application can easily be run and tested locally.
We run the application as a Docker container.

### (todo) Prereq's

- Go
- Mux
- Docker
- .. anything else?

### (rf) Run with Docker

comment
```
MacBook-Air:GoFetch matthewgarvin$ docker build -t points .
```

comment
```
MacBook-Air:GoFetch matthewgarvin$ docker run -d -p 8080:8080 points
```

comment
```
MacBook-Air:testData matthewgarvin$ curl http://localhost:8080/points/_status
{
    "status": "healthy"
}
```

### (rf) Example Flow

comment
```
MacBook-Air:testData matthewgarvin$ curl -X POST -d "@txList.json" http://localhost:8080/points/addTransactions
{
    "nTX": "5",
    "result": "success"
}
```

comment
```
MacBook-Air:testData matthewgarvin$ curl -X POST -d "@spendOrder.json" http://localhost:8080/points/spend
[
    {
        "payer": "DANNON",
        "points": -100
    },
    {
        "payer": "UNILEVER",
        "points": -200
    },
    {
        "payer": "MILLER COORS",
        "points": -4700
    }
]
```

comment
```
MacBook-Air:testData matthewgarvin$ curl http://localhost:8080/points/payerBalance
{
    "DANNON": 1000,
    "MILLER COORS": 5300,
    "UNILEVER": 0
}
```

#### (rf) Server Logs From Example Flow

Here are the server logs corresponding to the above example:

```
2021/03/12 15:49:28 Fetch Points Server serving at :8080
172.17.0.1 - - [12/Mar/2021:15:49:52 +0000] "GET /points/_status HTTP/1.1" 200 28 "" "curl/7.64.1"
172.17.0.1 - - [12/Mar/2021:15:50:06 +0000] "POST /points/addTransactions HTTP/1.1" 200 44 "" "curl/7.64.1"
172.17.0.1 - - [12/Mar/2021:15:50:27 +0000] "POST /points/spend HTTP/1.1" 200 201 "" "curl/7.64.1"
172.17.0.1 - - [12/Mar/2021:15:50:39 +0000] "GET /points/payerBalance HTTP/1.1" 200 68 "" "curl/7.64.1"
```


