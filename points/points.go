package main

import (
	"fmt"
	"sort"
	"time"
)

type PointsDB struct {
	PayerBalances       map[string]int
	UnspentTransactions []*Transaction // all txs which have positive unspent points remaining
	SpentTransactions   []*Transaction // all txs which have all their points spent
}

// an order to spend points
type SpendOrder struct {
	Points int `json:"points"`
}

// "spend points" response list item
type Payout struct {
	Payer  string `json:"payer"`
	Points int    `json:"points"`
}

// basic points transaction
type Transaction struct {
	Payer         string    `json:"payer"`
	Points        int       `json:"points"`
	UnspentPoints int       `json:"unspentPoints"` // todo: "-" for json after testing
	RawTimestamp  time.Time `json:"-"`
	Timestamp     string    `json:"timestamp"`
}

func (ps *PointServer) spendPoints(order *SpendOrder) ([]*Payout, error) {
	// Q. desired behavior if spendOrder is bigger than total points balance for that user? error or warning or ?

	// 1. sort unspent tx by timestamp (oldest to newest)
	sort.Sort(ByTime(ps.DB.UnspentTransactions))

	// incrementally subtract points per tx - "spend" each tx in order
	// keep track of points paid per payer ("payout")
	// once a tx is spent, remove it from the unspent db list and put it into the spent db list

	// map payer to points paid
	payments := make(map[string]int)

	remainder := order.Points
	// fmt.Printf("--- payments order: %v ---\n", order.Points)
	for remainder > 0 {

		// DEBUG
		// fmt.Println("remainder:", remainder)
		// fmt.Println("unspentTX:")
		// printJSON(ps.DB.UnspentTransactions)

		// take the front of the queue
		tx := ps.DB.UnspentTransactions[0]

		// DEBUG
		// fmt.Println("handling tx:")
		// printJSON(tx)

		if tx.UnspentPoints <= remainder {
			remainder -= tx.UnspentPoints
			payments[tx.Payer] -= tx.UnspentPoints
			tx.UnspentPoints = 0

			// pop front of queue
			ps.DB.UnspentTransactions = ps.DB.UnspentTransactions[1:]

			// move spent tx to appropriate db list
			ps.DB.SpentTransactions = append(ps.DB.SpentTransactions, tx)

		} else {
			tx.UnspentPoints -= remainder
			payments[tx.Payer] -= remainder
			remainder = 0
		}
	}

	// reconcile PayerBalances records with Payout records
	// todo: put this in a fn
	out := make([]*Payout, 0)
	var p *Payout
	var err error
	for payer, points := range payments {
		p = &Payout{
			Payer:  payer,
			Points: points,
		}

		err = ps.updateBalance(p)
		if err != nil {
			// todo: handle err; log;
			return nil, fmt.Errorf("failed to update balance: %v", err)
		}

		out = append(out, p)
	}

	return out, nil
}

// given a payout, update the corresponding payer balance
func (ps *PointServer) updateBalance(payout *Payout) error {
	/*
		if payout == nil {
			// handle; log;
		}
	*/
	ps.DB.PayerBalances[payout.Payer] += payout.Points
	return nil
}

func (ps *PointServer) fetchBalance() map[string]int {
	return ps.DB.PayerBalances
}

type ByTime []*Transaction

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].RawTimestamp.Before(a[j].RawTimestamp) }

func (ps *PointServer) addTransactions(txs []*Transaction) error {
	// handle nil input ?

	// could make this a constant, or put it in the server, or otherwise keep it in a better place
	timeFormat := "2006-01-02T15:04:05Z"

	var t time.Time
	for _, tx := range txs {
		tx.UnspentPoints = tx.Points

		// compute raw timestamp for sorting later
		// todo: error handling, date-time formatting, etc.
		// state assumptions in my solution
		t, _ = time.Parse(timeFormat, tx.Timestamp)
		tx.RawTimestamp = t

		// update balance for payer of this transaction
		ps.DB.PayerBalances[tx.Payer] += tx.Points

		ps.DB.UnspentTransactions = append(ps.DB.UnspentTransactions, tx)
	}

	return nil
}
