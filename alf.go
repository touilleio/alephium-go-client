package alephium

import (
	"bytes"
	"fmt"
	"github.com/willf/pad"
	"math/big"
	"math/rand"
	"strings"
)

type ALF struct {
	Amount *big.Int
}

var (
	OneQuintillionString = "1000000000000000000"
	OneQuintillionInt64  = int64(1000000000000000000)
	OneBillionInt64      = int64(1000000000)
	CoinInOneALF         = new(big.Int).SetInt64(OneQuintillionInt64)
	CoinInNanoALF        = new(big.Int).SetInt64(OneBillionInt64)
	//N = "א"
	N = "ALPH"
)

// TODO: should the ALF parsing function return a typed error instead of a bool?
func AFLFromALFString(amount string) (ALF, bool) {
	split := strings.Split(amount, ".")
	if len(split) == 1 {
		alfAmount, ok := new(big.Int).SetString(amount, 10)
		if ok {
			coinAmount := new(big.Int).Mul(alfAmount, CoinInOneALF)
			return ALF{Amount: coinAmount}, true
		}
	} else if len(split) == 2 {
		decimals := pad.Right(split[1], 18, "0")
		coinAmount, ok := new(big.Int).SetString(strings.Join([]string{split[0], decimals}, ""), 10)
		if ok {
			return ALF{Amount: coinAmount}, true
		}
	}
	return ALF{}, false
}

func ALFromCoinString(amount string) (ALF, bool) {
	alfAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return ALF{}, false
	}
	return ALF{Amount: alfAmount}, true
}

func (alf ALF) Add(other ALF) ALF {
	c := new(big.Int)
	c.Add(alf.Amount, other.Amount)
	return ALF{Amount: c}
}

func (alf ALF) Subtract(other ALF) ALF {
	c := new(big.Int)
	c.Sub(alf.Amount, other.Amount)
	return ALF{Amount: c}
}

func (alf ALF) Multiply(multiplier int64) ALF {
	c := new(big.Int)
	m := new(big.Int).SetInt64(multiplier)
	c.Mul(alf.Amount, m)
	return ALF{Amount: c}
}

func (alf ALF) Divide(divider int64) ALF {
	c := new(big.Int)
	m := new(big.Int).SetInt64(divider)
	c.Div(alf.Amount, m)
	return ALF{Amount: c}
}

func (alf ALF) Cmp(other ALF) int {
	return alf.Amount.Cmp(other.Amount)
}

func (alf ALF) String() string {
	if alf.Amount == nil {
		return "0"
	}
	return alf.Amount.String()
}

func (alf ALF) PrettyString() string {
	if alf.Amount == nil {
		return "0"
	}
	if alf.Amount.Cmp(CoinInNanoALF) > 0 {
		return fmt.Sprintf("%s%s", strings.TrimRight(fmt.Sprintf("%.9f", alf.FloatALF()), "0"), N)
	}
	return alf.Amount.String()
}

func (alf ALF) FloatALF() float64 {
	c := new(big.Int)
	nanoAFL := c.Div(alf.Amount, CoinInNanoALF).Int64()
	return float64(nanoAFL) / float64(OneBillionInt64)
}

func RandomALFAmount(upperLimit int) ALF {
	unit := rand.Intn(upperLimit)
	decimals := rand.Intn(int(OneBillionInt64))
	rAmountStr := fmt.Sprintf("%d.%d", unit, decimals)
	alf, _ := AFLFromALFString(rAmountStr)
	return alf
}

func RandomNanoALFAmount(upperLimit int) ALF {
	nanoALF := rand.Intn(upperLimit)
	c := new(big.Int).SetInt64(int64(nanoALF))
	m := new(big.Int).Mul(c, CoinInNanoALF)
	return ALF{Amount: m}
}

func (alf ALF) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	str := alf.Amount.String()
	buffer.WriteString(str)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (alf *ALF) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("NaN")
	}
	alf.Amount = new(big.Int)
	err := alf.Amount.UnmarshalJSON(b[1 : len(b)-1])
	return err
}

func ToNanoALF(alf ALF) int {
	m := new(big.Int).Div(alf.Amount, CoinInNanoALF)
	return int(m.Int64())
}
