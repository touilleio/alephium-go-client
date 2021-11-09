package alephium

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestALFConstruct(t *testing.T) {
	a1, ok := ALPHFromALPHString("12")
	assert.True(t, ok)
	assert.Equal(t, 12.000000000, a1.FloatALPH())
	assert.Equal(t, fmt.Sprintf("12.000000000%s", N), a1.PrettyString())
	assert.Equal(t, "12000000000000000000", a1.String())

	a2, ok := ALPHFromCoinString("12")
	assert.True(t, ok)
	assert.Equal(t, "12", a2.PrettyString())
	assert.Equal(t, "12", a2.String())
}

func TestZeroALF(t *testing.T) {
	a1 := ALPH{}
	assert.Equal(t, "0", a1.String())
}

func TestALFConstructWithDot(t *testing.T) {
	a1, ok := ALPHFromALPHString("12.12")
	assert.True(t, ok)
	assert.Equal(t, 12.120000000, a1.FloatALPH())
}

func TestALFAdd(t *testing.T) {
	a1, ok := ALPHFromALPHString("10.1")
	assert.True(t, ok)
	a2, ok := ALPHFromALPHString("2.02")
	assert.True(t, ok)
	a3, ok := ALPHFromALPHString("12.12")
	assert.True(t, ok)

	res := a1.Add(a2)
	assert.Equal(t, 0, res.Cmp(a3))
	assert.Equal(t, res, a3)
}

func TestALFSub(t *testing.T) {
	a1, ok := ALPHFromALPHString("10")
	assert.True(t, ok)
	a2, ok := ALPHFromALPHString("2")
	assert.True(t, ok)
	a3, ok := ALPHFromALPHString("12")
	assert.True(t, ok)

	res := a3.Subtract(a2)
	assert.Equal(t, 0, res.Cmp(a1))
	assert.Equal(t, res, a1)
}

func TestALFMul(t *testing.T) {
	a1, ok := ALPHFromALPHString("10")
	assert.True(t, ok)
	a2, ok := ALPHFromALPHString("2")
	assert.True(t, ok)
	a3, ok := ALPHFromALPHString("12")
	assert.True(t, ok)

	res := a2.Multiply(5)
	assert.Equal(t, 0, res.Cmp(a1))
	assert.Equal(t, res, a1)
	res = a2.Multiply(6)
	assert.Equal(t, 0, res.Cmp(a3))
	assert.Equal(t, res, a3)
}

func TestALFDiv(t *testing.T) {
	a1, ok := ALPHFromALPHString("10")
	assert.True(t, ok)
	a2, ok := ALPHFromALPHString("2")
	assert.True(t, ok)
	a3, ok := ALPHFromALPHString("12")
	assert.True(t, ok)

	res := a1.Divide(5)
	assert.Equal(t, 0, res.Cmp(a2))
	assert.Equal(t, res, a2)
	res = a3.Divide(6)
	assert.Equal(t, 0, res.Cmp(a2))
	assert.Equal(t, res, a2)
}

func TestALFDecode(t *testing.T) {
	a1, ok := ALPHFromALPHString("10")
	assert.True(t, ok)
	a2, ok := ALPHFromALPHString("2")
	assert.True(t, ok)
	a3, ok := ALPHFromALPHString("12")
	assert.True(t, ok)

	res := a1.Divide(5)
	assert.Equal(t, 0, res.Cmp(a2))
	assert.Equal(t, res, a2)
	res = a3.Divide(6)
	assert.Equal(t, 0, res.Cmp(a2))
	assert.Equal(t, res, a2)
}

func TestRoundAmount(t *testing.T) {
	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))
	alf := RandomALPHAmount(100)
	fmt.Printf("%s", alf.PrettyString())

	TenALF, ok := ALPHFromALPHString("10")
	assert.True(t, ok)

	// just to ensure at least 10 ALPH
	alf = alf.Add(TenALF)
	assert.True(t, alf.Cmp(TenALF) > 0)
	var finalTxAmount ALPH
	fiveALF := TenALF.Divide(2)
	txAmount := alf.Subtract(fiveALF)
	if txAmount.Cmp(TenALF) > 0 {
		finalTxAmount = TenALF
	} else {
		finalTxAmount = txAmount
	}

	assert.True(t, finalTxAmount.Cmp(TenALF) <= 0)
}

func TestNanoAmount(t *testing.T) {
	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))
	alf := RandomNanoALPHAmount(int(OneBillionInt64))
	fmt.Printf("%s\n", alf.PrettyString())
	fmt.Printf("%.4f\n", alf.FloatALPH())

	alf = RandomNanoALPHAmount(10000000)
	fmt.Printf("%s\n", alf.PrettyString())
}

type TestJsonStruct struct {
	Amount ALPH `json:"Amount"`
}

func TestJSON(t *testing.T) {
	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))
	alf := RandomNanoALPHAmount(int(OneBillionInt64))

	j1 := TestJsonStruct{
		Amount: alf,
	}

	b, err := json.Marshal(j1)
	assert.Nil(t, err)

	fmt.Printf("-->%s\n", string(b))

	j2 := &TestJsonStruct{}
	err = json.Unmarshal(b, j2)
	assert.Nil(t, err)

	fmt.Printf("%s\n", j2.Amount.PrettyString())
}
