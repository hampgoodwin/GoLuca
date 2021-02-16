package transaction

// Balanced checks that a transaction is balanced, that is to say that debits equal credits
func (t *Transaction) Balanced() bool {
	sum := float64(0)
	for _, e := range t.Entries {
		sum += e.Amount
	}
	return sum == float64(0)
}
