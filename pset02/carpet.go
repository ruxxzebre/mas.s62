package main

// PUT UNDER THE CARPET

type Nonce struct {
	nonceRef string
	nonceRefLen int
}

func modulo(a int, b int) int {
	return int(a/b)
}

func (self *Nonce) setRef(ref string) {
	self.nonceRef = ref
	self.nonceRefLen = len(ref)
}

func (self Nonce) generateNonce(i int) string {
	if i < self.nonceRefLen {
		return string(self.nonceRef[i])
	}

	mod := modulo(i, self.nonceRefLen)

	nonceValue := ""

	for i := 0; i < mod; i++ {
		nonceValue += string(self.nonceRef[self.nonceRefLen - 1])
	}

	nonceValue += string(self.nonceRef[i%self.nonceRefLen])

	return nonceValue
}