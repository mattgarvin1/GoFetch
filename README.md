# Fetch Points Server

## (chk) Note To Reviewer

[Here](./problemStatement.pdf) is the original problem statement.

If you would like to see a bit of my thought process as I was working, 
I jotted down some [notes](./thoughtProcess.txt) while working on my solution.

This application functions at a basic level in that
you can walk through the example flow provided in the problem statement
and the server behaves as desired and gives the desired responses.

If I were to spend some more time on this application, I would 
- create more unit tests
- hook up Travis to run the tests against commits and PRs
- tighten up HTTP response codes / error handling
- add other useful routes to the API - for example, a route to
fetch the balances corresponding to a list of payers so that
you don't have to receive all payer balances if you're 
only interested in querying a subset of payers

Thanks for taking the time to review my work!

## (todo) Running the Server

The application can easily be run and tested locally.
We run the application as a Docker container.

### (todo) Prereq's

- Go
- Mux
- Docker
- .. anything else?

### (chk) Run with Docker

At the root directory of this repo, do the following:

1. Build the points server image:
```docker build -t points .```

2. Run the image you just built as a container:
```docker run -d -p 8080:8080 points```

3. Ping the points server:
```curl http://localhost:8080/points/_status```

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

#### (chk) Server Logs From Example Flow

Here are the server logs corresponding to the above example:

```
2021/03/12 15:49:28 Fetch Points Server serving at :8080
172.17.0.1 - - [12/Mar/2021:15:49:52 +0000] "GET /points/_status HTTP/1.1" 200 28 "" "curl/7.64.1"
172.17.0.1 - - [12/Mar/2021:15:50:06 +0000] "POST /points/addTransactions HTTP/1.1" 200 44 "" "curl/7.64.1"
172.17.0.1 - - [12/Mar/2021:15:50:27 +0000] "POST /points/spend HTTP/1.1" 200 201 "" "curl/7.64.1"
172.17.0.1 - - [12/Mar/2021:15:50:39 +0000] "GET /points/payerBalance HTTP/1.1" 200 68 "" "curl/7.64.1"
```


