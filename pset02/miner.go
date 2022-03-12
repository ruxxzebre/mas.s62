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
	lastBlock Block
	data string
	MiningDone chan bool
	chanAmount int
	staleSync map[int] bool
	ResultChan chan Block
}

func makeMiner(targetBits uint8, nonceSize uint32, data string) (Miner, error) {
	block, err := GetTipFromServer()

	if err != nil {
		fmt.Println("Block error.")
		return *new(Miner), err
	}	

	return Miner{
		targetBits,
		nonceSize,
		block,
		data,
		make(chan bool),
		0,
		make(map[int] bool),
		make(chan Block),
	}, nil
}

func (m *Miner) initStaleMap() {
	for i := 0; i < m.chanAmount; i++ {
		m.staleSync[i] = false
	}
}

func (m *Miner) poolNewBlocks() error {
	for {
		time.Sleep(time.Second)
		block, err := GetTipFromServer()

		if err != nil {
			fmt.Println("Block error.")
			return err
		}

		m.lastBlock = block
		
		if block.Hash().ToString() != m.lastBlock.Hash().ToString() {
			fmt.Println("Block mismatch. Updating...")
			for i := 0; i < m.chanAmount; i++ {
				m.staleSync[i] = true
			}
		}
	}
}

func (m *Miner) Run(chanAmount int) {
	uchanAmount := uint32(chanAmount)
	m.chanAmount = chanAmount
	m.initStaleMap()
	go m.poolNewBlocks()

	incval := uint32(m.nonceSize / uchanAmount - 1)

	m.logStart()
	for i := uint32(0); i < uchanAmount; i++ {
		go m.mine(i, incval, chanAmount)
	}
}

func (m *Miner) initNewBlock() Block {
	return Block{
		m.lastBlock.Hash(),
		m.data,
		"",
	}
}

func (m *Miner) initLoopingVars(jobIdx uint32, incval uint32) (uint32, uint32) {
	i := jobIdx * incval
	top := (jobIdx + 1) * incval
	return i, top
}

func (m *Miner) mine(jobIdx uint32, incval uint32, chanam int) {
	start := time.Now()

	newBlock := m.initNewBlock()

	i, top := m.initLoopingVars(jobIdx, incval)
	
	for i < top {
		if m.staleSync[chanam] {
			fmt.Println("Revoking stale blocks")
			newBlock = m.initNewBlock()
			i, top = m.initLoopingVars(jobIdx, incval)
			m.staleSync[chanam] = false
		}

		newBlock.Nonce = fmt.Sprint(i)
		if m.checkWork(newBlock) {
			m.MiningDone <- true	
			log.Println("DONE")
			m.logWork(start, int(jobIdx), chanam)
			m.ResultChan <- newBlock
			return
		}
		i++
	}
}

func check (err error) {
	if err != nil {
		panic(err)
	}
}

func (m *Miner) logStart() {
	fmt.Printf("Mining started. Difficulty: %d\n", m.targetBits)
}

func (m *Miner) logDone(t string) {
	fmt.Printf("Mining done. Difficulty: %d. Time: %v\n", m.targetBits, t)
}

func (m *Miner) logWork(start time.Time, cidx int, chanam int) {
	elapsed := time.Since(start)

	f, err := os.OpenFile("./measure.csv", os.O_APPEND|os.O_WRONLY, 0666)
	check(err)

	stat, err := f.Stat()
	check(err)

	if stat.Size() == 0 {
		f.WriteString("DIFFICULTY,CHANNELS,CHANNEL_INDEX,ELAPSED_TIME\n")
	}

	template := "%d,%d,%d,%v\n"
	f.WriteString(fmt.Sprintf(template, m.targetBits, chanam, cidx, 
		elapsed.String()),
	)
	f.Close()
	m.logDone(elapsed.String())
}

func (m *Miner) checkWork(bl Block) bool {
	return OldCheckWork(bl, m.targetBits)
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
