package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/karlseguin/ccache"
	"github.com/sirupsen/logrus"
)

var (
	//modN    = uint32(1 << 26)
	//modNBig = big.NewInt(int64(modN))

	NCache                 = ccache.New(ccache.Configure().MaxSize(16))
	NBase                  = uint32(1 << 26)
	NBaseBig               = big.NewInt(int64(NBase))
	IncreaseStart          = int64(600 * 1024)
	IncreasePeriodForN     = int64(50 * 1024)
	NIncreasementHeightMax = int64(9216000)

	modM = func() []byte {
		buf := bytes.NewBuffer(nil)
		for i := 0; i < 1024; i++ {
			buf.Write(UInt64BEToBytes(uint64(i)))
		}
		return buf.Bytes()
	}()
)

func modN(height int64) uint32 {
	height = MinInt64(height, NIncreasementHeightMax)
	if height < IncreaseStart {
		return NBase
	} else if height >= NIncreasementHeightMax {
		return 2147387550
	} else {
		iterationsNumber := int((height-IncreaseStart)/IncreasePeriodForN + 1)
		k := fmt.Sprintf("iter%d", iterationsNumber)
		item := NCache.Get(k)
		if item == nil {
			res := NBase
			for i := 0; i < iterationsNumber; i++ {
				res = res / 100 * 105
			}
			NCache.Set(k, res, 30*time.Minute)
			return res
		}
		return item.Value().(uint32)
	}
}

func modNBig(height int64) *big.Int {
	return big.NewInt(int64(modN(height)))
}

func genIndexes(seed []byte, height int64) []uint32 {
	seedHash := Blake2b(seed)
	hash := append(seedHash, seedHash...)

	indexes := make([]uint32, 32)
	for i := range indexes {
		indexes[i] = binary.BigEndian.Uint32(hash[i:i+4]) % modN(height)
	}

	return indexes
}

func genIndexesByte(seed []byte, height int64) [][]byte {
	indexesByte := make([][]byte, 32)
	indexes := genIndexes(seed, height)
	for i := range indexes {
		indexesByte[i] = UInt32BEToBytes(indexes[i])
	}

	return indexesByte
}

// port from https://github.com/mhssamadani/ErgoStratumServer/tree/master
func Autolykos2Hash(msg, nonce1, nonce2 string, height int64, target *big.Int) []byte {
	msgWithNonce := MustStringToHexBytes(msg + nonce1 + nonce2)
	msgHash := Blake2b(msgWithNonce)
	eBuffer := bytes.NewBuffer(nil)
	i := UInt32BEToBytes(uint32(new(big.Int).Mod(new(big.Int).SetBytes(msgHash[24:32]), modNBig(height)).Int64()))
	h := UInt32BEToBytes(uint32(height))
	eBuffer.Write(i)
	eBuffer.Write(h)
	eBuffer.Write(modM)
	e := Blake2b(eBuffer.Bytes())[1:32]
	J := genIndexesByte(append(e, msgWithNonce...), height)

	f := func(indexes [][]byte) []byte {
		bigHash := big.NewInt(0)
		for i := range indexes {
			temp := Blake2b(append(append(indexes[i], h...), modM...))
			temp[0] = 0x00
			bigHash = new(big.Int).Add(bigHash, new(big.Int).SetBytes(temp))
		}
		bigHashBytes := bigHash.Bytes()
		return bigHashBytes
	}(J)
	fh := Blake2b(f)

	bT := Hash2BigTarget(fh)
	if bT.Cmp(target) <= 0 {
		logrus.Debugf("solving nonce2: %s", nonce2)
		logrus.Debugf("msgWithNonce: %x", msgWithNonce)
		logrus.Debugf("msgHash: %064x", msgHash)
		logrus.Debugf("i: %x", i)
		logrus.Debugf("h: %x", h)
		//logrus.Debugf("modM: %x", modM)
		logrus.Debugf("e: %x", e)
		logrus.Debugf("J: %x", J)
		logrus.Debugf("f: %064x", f)
		logrus.Debugf("fh: %064x", fh)
	}

	return fh
}
