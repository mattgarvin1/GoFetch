# Fetch Points Server

-> describe API endpoints

## Note To Reviewer

Thanks for taking the time to review my work!
[Here](./problemStatement.pdf) is the original problem statement.
If you would like to see a bit of what my thought process was,
I jotted down some [notes](./thoughtProcess.txt) while I was working.

This application functions at a basic level in that
you can walk through the example flow provided in the problem statement
and the server behaves as desired and gives the desired responses.

If I were to spend some more time on this application, I would 
- create more unit tests
- hook up [Travis](https://travis-ci.com/) to run the tests against commits and PRs
- tighten up HTTP response codes / error handling
- add other useful routes to the API - for example, a route to
fetch the balances corresponding to a list of payers so that
you don't have to receive all payer balances if you're 
only interested in querying a subset of payers

## Running The Server

The points server can be run locally in a [Docker](https://www.docker.com/) container.

### Prereq's

This server is written in [Go](https://golang.org/) and uses [mux](https://github.com/gorilla/mux).

### How To Run With Docker

At the root directory of this repo, do the following:

1. Build the points server image:
```docker build -t points .```

2. Run the image you just built as a container:
```docker run -d -p 8080:8080 points```

3. Ping the points server:
```curl http://localhost:8080/points/_status```

### (chk) Example Flow Walkthrough

Here is a walkthrough of me testing the server locally, setup exactly as detailed in the above section.

```
MacBook-Air:GoFetch matthewgarvin$ curl http://localhost:8080/points/_status
{
    "status": "healthy"
}
MacBook-Air:GoFetch matthewgarvin$ curl http://localhost:8080/points/payerBalance
{}
MacBook-Air:GoFetch matthewgarvin$ curl -X POST -d "@testData/txList.json" http://localhost:8080/points/addTransactions
{
    "nTX": "5",
    "result": "success"
}
MacBook-Air:GoFetch matthewgarvin$ curl -X POST -d "@testData/spendOrder.json" http://localhost:8080/points/spend
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
MacBook-Air:GoFetch matthewgarvin$ curl http://localhost:8080/points/payerBalance
{
    "DANNON": 1000,
    "MILLER COORS": 5300,
    "UNILEVER": 0
}
```

