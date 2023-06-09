package common

func EstimateFeeUsage(txSize uint64) uint64 {
	return txSize / 4 * 50
}
