package types

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
)

// Super nodes
type Supers struct {
	Candidates []*Member
	PreHash    arry.Hash
}

type Member struct {
	Signer arry.Address
	PeerId string
	Weight uint64
}

type Candidates struct {
	Members []*Member
}

func NewCandidates() *Candidates {
	return &Candidates{Members: make([]*Member, 0)}
}

func (c *Candidates) Set(newMem *Member) {
	for _, mem := range c.Members {
		if mem.Signer.IsEqual(newMem.Signer) {
			return
		}
	}
	c.Members = append(c.Members, newMem)
}

func (c *Candidates) Remove(reMem *Member) {
	for i, mem := range c.Members {
		if mem.Signer.IsEqual(reMem.Signer) {
			c.Members = append(c.Members[0:i], c.Members[i+1:]...)
			return
		}
	}
}

func (c *Candidates) Len() int {
	return len(c.Members)
}

type SortableCandidates []*Member

func (p SortableCandidates) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p SortableCandidates) Len() int      { return len(p) }
func (p SortableCandidates) Less(i, j int) bool {
	if p[i].Weight < p[j].Weight {
		return false
	} else if p[i].Weight > p[j].Weight {
		return true
	} else {
		return p[i].Weight < p[j].Weight
	}
}