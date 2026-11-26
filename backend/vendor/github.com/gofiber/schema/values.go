package schema

import (
	"errors"
	"hash/maphash"
	"slices"
)

// errValuesLength is returned by DecodeValues for keys and values of
// different lengths.
var errValuesLength = errors.New("schema: keys and values differ in length")

const (
	// maxInlinePairs is how many pairs DecodeValues groups in arrays on its
	// own stack; more are grouped in slices allocated for the call.
	maxInlinePairs = 32
	// maxScanPairs is how many pairs are grouped by comparing each key with
	// the keys before it rather than through a hash table, which for so few
	// costs more to hash than to compare.
	maxScanPairs = 8
)

// source is what a decode reads the values of a key from: the map Decode was
// given, or the pairs DecodeValues was, grouped by key. It is a struct rather
// than an interface so that a DecodeValues call's index, held on its stack,
// does not escape through a method call the compiler cannot see into.
type source struct {
	m     map[string][]string
	pairs *pairIndex
}

// lookup returns the values of key as far as the required-key and default
// checks read them, which is whether there are any and what the last one is:
// all of them from a map, and the last alone from pairs, which saves
// gathering them. It reports whether src has the key.
func (src *source) lookup(key string) ([]string, bool) {
	if src.pairs != nil {
		if i := src.pairs.find(key); i >= 0 {
			l := src.pairs.last[i]
			return src.pairs.values[l : l+1 : l+1], true
		}
		return nil, false
	}
	v, ok := src.m[key]
	return v, ok
}

// hasBytes reports whether src has the key held in key.
func (src *source) hasBytes(key []byte) bool {
	if src.pairs != nil {
		return src.pairs.findBytes(key) >= 0
	}
	_, ok := src.m[string(key)]
	return ok
}

// pairIndex groups key/value pairs by key without building a map: keys are
// hashed into an open-addressing table that holds the first pair of each
// key, or for a few pairs compared with the keys before them, and each pair
// links to the next with the same key. A key's values are then the values of
// its pairs, which a single pair or a run of adjacent ones hands out as a
// subslice of values, and any other pattern gathers into a buffer of the
// caller's.
//
// The keys are the request's, so the table hashes them as a Go map would:
// exactly, as they are grouped, and with a seed drawn per process (see
// pairSeed), so no choice of keys can pile them onto one slot.
//
// DecodeValues sets its slices on it where it lies, and its methods store
// integers only, into those slices and mask: a slice stored through a pointer
// to it would count, for the compiler, as a store to the heap, and take the
// arrays DecodeValues keeps the index in on its stack with it.
type pairIndex struct {
	keys, values []string
	// table holds, for each key, its first pair plus one, so that 0 marks a
	// free slot; it is a power of two long and at most half full, and nil
	// for maxScanPairs pairs or fewer
	table []int32
	// next holds, for each pair, the next pair with the same key plus one,
	// or 0 for its key's last pair
	next []int32
	// last holds, for a key's first pair, the key's last pair, and -1 for
	// every other pair
	last []int32
	mask uint64
}

// pairSeed seeds the hash of pairIndex. A hash an attacker could predict,
// or one under which distinct keys collide by construction, as keys that
// differ only in case do under a case-folding hash, would let a request
// with enough such keys make every insert probe past all the keys before
// it, quadratic in their number.
var pairSeed = maphash.MakeSeed()

// pairHash is the hash pairIndex keeps key in, and pairHashBytes the same
// hash of a key held in a byte slice.
func pairHash(key string) uint64 { return maphash.String(pairSeed, key) }

func pairHashBytes(key []byte) uint64 { return maphash.Bytes(pairSeed, key) }

// indexSize returns the table length that indexes n pairs.
func indexSize(n int) int {
	size := 2
	for size < 2*n {
		size *= 2
	}
	return size
}

// index fills in the index of the pairs p.keys[i]=p.values[i], kept in
// p.table, p.next and p.last, which must be zeroed and sized
// indexSize(len(p.keys)) (or nil, for maxScanPairs pairs or fewer),
// len(p.keys) and len(p.keys).
func (p *pairIndex) index() {
	if p.table == nil {
		p.scan()
		return
	}
	p.mask = uint64(len(p.table) - 1)
	p.build()
}

// scan fills in the index of p's pairs, few enough to compare each key with
// the first pairs before it.
func (p *pairIndex) scan() {
	keys := p.keys
	for i, key := range keys {
		p.last[i] = int32(i) //nolint:gosec // G115 - bounded by the number of pairs
		for first := range i {
			if p.last[first] >= 0 && keys[first] == key {
				p.next[p.last[first]] = int32(i + 1) //nolint:gosec // G115 - bounded by the number of pairs
				p.last[first] = int32(i)             //nolint:gosec // G115 - bounded by the number of pairs
				p.last[i] = -1
				break
			}
		}
	}
}

// build fills in the index of p's pairs through the hash table.
func (p *pairIndex) build() {
	keys := p.keys
	for i, key := range keys {
		for j := pairHash(key) & p.mask; ; j = (j + 1) & p.mask {
			e := p.table[j]
			if e == 0 {
				p.table[j] = int32(i + 1) //nolint:gosec // G115 - bounded by the number of pairs
				p.last[i] = int32(i)      //nolint:gosec // G115 - bounded by the number of pairs
				break
			}
			if first := e - 1; keys[first] == key {
				p.next[p.last[first]] = int32(i + 1) //nolint:gosec // G115 - bounded by the number of pairs
				p.last[first] = int32(i)             //nolint:gosec // G115 - bounded by the number of pairs
				p.last[i] = -1
				break
			}
		}
	}
}

// first reports whether pair i is the first pair of its key.
func (p *pairIndex) first(i int) bool {
	return p.last[i] >= 0
}

// find returns the first pair of key, or -1 when no pair has it.
func (p *pairIndex) find(key string) int {
	if p.table == nil {
		for first, k := range p.keys {
			if p.last[first] >= 0 && k == key {
				return first
			}
		}
		return -1
	}
	for j := pairHash(key) & p.mask; ; j = (j + 1) & p.mask {
		e := p.table[j]
		if e == 0 {
			return -1
		}
		if first := int(e - 1); p.keys[first] == key {
			return first
		}
	}
}

// findBytes is find for a key held in a byte slice.
func (p *pairIndex) findBytes(key []byte) int {
	if p.table == nil {
		for first, k := range p.keys {
			if p.last[first] >= 0 && k == string(key) {
				return first
			}
		}
		return -1
	}
	for j := pairHashBytes(key) & p.mask; ; j = (j + 1) & p.mask {
		e := p.table[j]
		if e == 0 {
			return -1
		}
		if first := int(e - 1); p.keys[first] == string(key) {
			return first
		}
	}
}

// group returns the values of the key whose first pair is i, capped at their
// length: a subslice of values when its pairs are adjacent, else the values
// appended to scratch.
func (p *pairIndex) group(i int, scratch []string) []string {
	end := i
	for n := p.next[end]; n != 0; n = p.next[end] {
		if int(n-1) != end+1 {
			for {
				scratch = append(scratch, p.values[i])
				n := p.next[i]
				if n == 0 {
					return slices.Clip(scratch)
				}
				i = int(n - 1)
			}
		}
		end++
	}
	return p.values[i : end+1 : end+1]
}
