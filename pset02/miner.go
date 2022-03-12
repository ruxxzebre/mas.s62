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
	}
}

func (miner *Miner) poolTip() (Block, error) {
	prevBlock, err := GetTipFromServer()

	if err != nil {
		fmt.Println("Block error.")
		return prevBlock, err
	}
	return prevBlock, nil
}

func (miner *Miner) Run(channel chan Block, chanAmount int) {
	prevBlock, err := GetTipFromServer()

	if err != nil {
		fmt.Println("Block error.")
	}

	uchanAmount := uint32(chanAmount)

	incval := uint32(miner.nonceSize / uchanAmount - 1)

	for i := uint32(0); i < uchanAmount; i++ {
		go miner.mine(channel, prevBlock, i, incval, chanAmount)
	}
}

func (miner *Miner) mine(channel chan Block, prevBlock Block, jobIdx uint32, incval uint32, chanam int) {

	start := time.Now()

	var newBlock Block
	newBlock.PrevHash = prevBlock.Hash()
	newBlock.Name = "ruxxzebre"

	i := jobIdx * incval
	top := (jobIdx + 1) * incval
	
	for i < top {
		newBlock.Nonce = fmt.Sprint(i)
		if OldCheckWork(newBlock, miner.targetBits) {
			channel <- newBlock
			log.Println("DONE")
			miner.logWork(start, int(jobIdx), chanam)
		}
		i++
	}
}

func check (err error) {
	if err != nil {
		panic(err)
	}
}

func (miner *Miner) logWork(start time.Time, cidx int, chanam int) {
	elapsed := time.Since(start)

	f, err := os.OpenFile("./measure.csv", os.O_APPEND|os.O_WRONLY, 0666)
	check(err)

	stat, err := f.Stat()
	check(err)

	if stat.Size() == 0 {
		f.WriteString("DIFFICULTY,CHANNELS,CHANNEL_INDEX,ELAPSED_TIME\n")
	}

	template := "%d,%d,%d,%v\n"
	f.WriteString(fmt.Sprintf(template, miner.targetBits, chanam, cidx, 
		elapsed.String()),
	)
	f.Close()
}

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
