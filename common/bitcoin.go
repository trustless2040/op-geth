package common

func GetFeeRateByBlockHeight(blockHeight uint64, size uint64, feeRate float64) uint64 {
	if feeRate <= 0 {
		feeRate = 50
	}
	if blockHeight > 0 {
		fee := uint64(float64(size/4) * feeRate)
		return fee
	}
	return 0
}
