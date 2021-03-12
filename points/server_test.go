package main

import (
	"encoding/json"
	"testing"
)

/*

1. call to "add transaction" route:

	{ "payer": "DANNON", "points": 1000, "timestamp": "2020-11-02T14:00:00Z" }
	{ "payer": "UNILEVER", "points": 200, "timestamp": "2020-10-31T11:00:00Z" }
	{ "payer": "DANNON", "points": -200, "timestamp": "2020-10-31T15:00:00Z" }
	{ "payer": "MILLER COORS", "points": 10000, "timestamp": "2020-11-01T14:00:00Z" }
	{ "payer": "DANNON", "points": 300, "timestamp": "2020-10-31T10:00:00Z" }

2.0. call to "spend points" route:

	{ "points": 5000 }

2.1. response from "spend points" call:

	[
		{ "payer": "DANNON", "points": -100 },
		{ "payer": "UNILEVER", "points": -200 },
		{ "payer": "MILLER COORS", "points": -4,700 }
	]

3.0. call to "points balance" route after the spend call
3.1. response from "points balance" call:

	{
		"DANNON": 1000,
		"UNILEVER": 0,
		"MILLER COORS": 5300
	}

*/

func TestServer(t *testing.T) {

	// init the (in-memory) db
	db := &PointsDB{
		PayerBalances: make(map[string]int),
		// UnspentTransactions: make([]*Transaction, 0),
		// SpentTransactions:   make([]*Transaction, 0),
	}

	// init the server over the db
	ps := &PointServer{
		DB: db,
	}

	// ---- 1 ---- //

	txListJSON := []byte(`
	{
		"txList" : 	[
			{ "payer": "DANNON", "points": 1000, "timestamp": "2020-11-02T14:00:00Z" },
			{ "payer": "UNILEVER", "points": 200, "timestamp": "2020-10-31T11:00:00Z" },
			{ "payer": "DANNON", "points": -200, "timestamp": "2020-10-31T15:00:00Z" },
			{ "payer": "MILLER COORS", "points": 10000, "timestamp": "2020-11-01T14:00:00Z" },
			{ "payer": "DANNON", "points": 300, "timestamp": "2020-10-31T10:00:00Z" }	
		]	
	}
	`)

	txList := &TXListJSON{}
	err := json.Unmarshal(txListJSON, txList)
	if err != nil {
		t.Fatalf("failed to unmarshal json: %s", err)
	}

	// fmt.Println("txList:")
	// printJSON(txList)

	ps.addTransactions(txList.TXList)

	// ---- 2 ---- //

	spendOrderJSON := []byte(`
		{ 
			"points": 5000
		}
	`)
	spendOrder := &SpendOrder{}
	err = json.Unmarshal(spendOrderJSON, spendOrder)
	if err != nil {
		t.Fatalf("failed to unmarshal json: %s", err)
	}
	payoutList, err := ps.spendPoints(spendOrder)
	if err != nil {
		t.Fatalf("spendPoints failed: %s", err)
	}

	desiredPayout := map[string]int{
		"DANNON":       -100,
		"UNILEVER":     -200,
		"MILLER COORS": -4700,
	}

	// todo: check that all keys appear in payout list (?)
	var got, want int
	for _, payout := range payoutList {
		got = payout.Points
		want = desiredPayout[payout.Payer]
		if got != want {
			t.Errorf("got: %v; want: %v", got, want)
		}
	}

	// ---- 3 ---- //

	desiredBalance := map[string]int{
		"DANNON":       1000,
		"UNILEVER":     0,
		"MILLER COORS": 5300,
	}
	balances := ps.fetchBalance()
	for payer, balance := range balances {
		want = desiredBalance[payer]
		if balance != want {
			t.Errorf("got: %v; want: %v", balance, want)
		}
	}

}
