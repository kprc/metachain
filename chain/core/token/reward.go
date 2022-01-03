package token

import "math/big"

const(
	//3 seconds one block
	BlockHeightPerRound uint = 2*10512000  //365*24*60*(60/3)  two years
	FirstRoundReward uint = 100
)

func Reward(blkHeight uint64) *big.Int  {
	l:=blkHeight/uint64(BlockHeightPerRound)

	r:=&big.Int{}
	r.SetUint64(uint64(FirstRoundReward))

	u:=&big.Int{}
	u.SetInt64(1*M)

	r = r.Mul(r,u)

	z:=r

	for i:=0;i<int(l);i++{
		z = div2(z)
	}

	return z
}

func div2(r *big.Int) *big.Int  {
	z:=&big.Int{}

	d2:=&big.Int{}
	d2.SetInt64(2)

	z = r.Div(r,d2)

	return z
}