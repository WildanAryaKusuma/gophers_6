package atm

import (
	"fmt"
	"sync"
)

type Rekening struct {
	Pemilik string
	Saldo   int
	mu sync.Mutex
}

// TIPS

// method receiver
// read
func (r *Rekening) Info() {
	fmt.Printf("Rekening %s, saldo %d\n", r.Pemilik, r.Saldo)
}

// method pointer receiver
// create, update, delete
func (r *Rekening) Setor(jumlah int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Saldo = r.Saldo + jumlah
}