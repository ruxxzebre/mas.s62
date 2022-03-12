package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Note that "targetBits" for this assignment, at least initially, is 33.
// This could change during the assignment duration!  I will post if it does.

// Mine mines a block by varying the nonce until the hash has targetBits 0s in
// the beginning.  Could take forever if targetBits is too high.
// Modifies a block in place by using a pointer receiver.
type Miner struct {
	targetBits uint8
	nonceSize uint32
}

func makeMiner(targetBits uint8, nonceSize uint32) Miner {
	return Miner{
		targetBits,
		nonceSize,
		// channelIndex*incrementval,
		// (channelIndex+1)*incrementval,
	}
}

func (self *Miner) Run(channel chan Block, chanAmount int) {
	prevBlock, err := GetTipFromServer()

	if err != nil {
		fmt.Println("Block error.")
	}

	uchanAmount := uint32(chanAmount)

	incval := uint32(self.nonceSize / uchanAmount - 1)

	for i := uint32(0); i < uchanAmount; i++ {
		go self.mine(channel, prevBlock, i, incval)
	}
}

func (self *Miner) mine(channel chan Block, prevBlock Block, jobIdx uint32, incval uint32) {

	start := time.Now()
	// var nonce Nonce
	// nonce.setRef("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	var newBlock Block
	newBlock.PrevHash = prevBlock.Hash()
	newBlock.Name = "ruxxzebre"

	i := jobIdx * incval
	top := (jobIdx + 1) * incval
	
	for i < top {
		newBlock.Nonce = fmt.Sprint(i)
		if OldCheckWork(newBlock, self.targetBits) {
			channel <- newBlock
			// close(channel)

			log.Println("DONE")
			logWork(start, int(self.targetBits), int(jobIdx))
		}
		i++
	}
	// your mining code here
	// also feel free to get rid of this method entirely if you want to
	// organize things a different way; this is just a suggestion
}

func logWork(start time.Time, targetBits int, cidx int) {
	elapsed := time.Since(start)

	f, err := os.OpenFile("./measure.txt", os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		template := "DIFFICULTY: %d; CHANNEL INDEX: %d; ELAPSED TIME: {ms: %d, s: %g}\n"
		f.WriteString(fmt.Sprintf(template, targetBits, cidx, elapsed.Milliseconds(), elapsed.Seconds()))
		f.Close()
}

// func NewCheckWork(bl Block, targetBits uint8) bool {
// 	h := bl.Hash()

// 	i := uint8(0)
// 	for {

// 	}
// }

func OldCheckWork(bl Block, targetBits uint8) bool {
	h := bl.Hash()

	for i := uint8(0); i < targetBits; i++ {
		// for every bit from the MSB down, check if it's a 1.
		// If it is, stop and fail.
		// Could definitely speed this up by checking bytes at a time.
		// Left as excercise for the reader...?
		if (h[i/8]>>(7-(i%8)))&0x01 == 1 {
			return false
		}
	}
	return true
}

// CheckWork checks if there's enough work
func CheckWork(bl Block, targetBits uint8) bool {
	h := bl.Hash()
	tbrem := targetBits%8
	tbmod := int(targetBits/8)

	// fmt.Printf("tb: %d, rem: %d, mod: %d\n", targetBits, tbrem, tbmod)

	for i := 0; i < tbmod; i++ {
		if h[i] > 0 {
			return false
		}
	}

	for i := uint8(targetBits - tbrem); i < tbrem; i++ {
		// for every bit from the MSB down, check if it's a 1.
		// If it is, stop and fail.
		// Could definitely speed this up by checking bytes at a time.
		// Left as excercise for the reader...?
		if (h[i/8]>>(7-(i%8)))&0x01 == 1 {
			return false
		}
	}
	return true
}
