package dictpool

import (
	"sort"

	"github.com/savsgio/gotils/strconv"
)

func New() *Dict {
	return new(Dict)
}

func (d *Dict) allocKV() *KV {
	n := d.len()

	if cap(d.D) > n {
		d.D = d.D[:n+1]
	} else {
		d.D = append(d.D, KV{}) // nolint:exhaustruct
	}

	return &d.D[n]
}

func (d *Dict) append(key string, value interface{}) {
	kv := d.allocKV()
	kv.Key = key
	kv.Value = value
}

func (d *Dict) indexOf(key string) int {
	n := d.len()

	if d.BinarySearch {
		idx := sort.Search(n, func(i int) bool {
			return key <= d.D[i].Key
		})

		if idx < n && d.D[idx].Key == key {
			return idx
		}
	} else {
		for i := 0; i < n; i++ {
			if d.D[i].Key == key {
				return i
			}
		}
	}

	return -1
}

func (d *Dict) len() int {
	return len(d.D)
}

func (d *Dict) swap(i, j int) {
	d.D[i], d.D[j] = d.D[j], d.D[i]
}

func (d *Dict) less(i, j int) bool {
	return d.D[i].Key < d.D[j].Key
}

func (d *Dict) get(key string) interface{} {
	if idx := d.indexOf(key); idx > -1 {
		return d.D[idx].Value
	}

	return nil
}

func (d *Dict) set(key string, value interface{}) {
	if idx := d.indexOf(key); idx > -1 {
		d.D[idx].Value = value
	} else {
		d.append(key, value)

		if d.BinarySearch {
			sort.Sort(d)
		}
	}
}

func (d *Dict) del(key string) {
	if idx := d.indexOf(key); idx > -1 {
		d.D = append(d.D[:idx], d.D[idx+1:]...)
	}
}

func (d *Dict) has(key string) bool {
	return d.indexOf(key) > -1
}

func (d *Dict) reset() {
	d.D = d.D[:0]
}

// Len is the number of elements in the Dict.
func (d *Dict) Len() int {
	return d.len()
}

// Swap swaps the elements with indexes i and j.
func (d *Dict) Swap(i, j int) {
	d.swap(i, j)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (d *Dict) Less(i, j int) bool {
	return d.less(i, j)
}

// Get get data from key.
func (d *Dict) Get(key string) interface{} {
	return d.get(key)
}

// GetBytes get data from key.
func (d *Dict) GetBytes(key []byte) interface{} {
	return d.Get(strconv.B2S(key))
}

// Set set new key.
func (d *Dict) Set(key string, value interface{}) {
	d.set(key, value)
}

// SetBytes set new key.
func (d *Dict) SetBytes(key []byte, value interface{}) {
	d.Set(strconv.B2S(key), value)
}

// Del delete key.
func (d *Dict) Del(key string) {
	d.del(key)
}

// DelBytes delete key.
func (d *Dict) DelBytes(key []byte) {
	d.Del(strconv.B2S(key))
}

// Has check if key exists.
func (d *Dict) Has(key string) bool {
	return d.has(key)
}

// HasBytes check if key exists.
func (d *Dict) HasBytes(key []byte) bool {
	return d.Has(strconv.B2S(key))
}

// Reset reset dict.
func (d *Dict) Reset() {
	d.reset()
}

// Map convert to map.
func (d *Dict) Map(dst DictMap) {
	for i := range d.D {
		kv := &d.D[i]

		sd, ok := kv.Value.(*Dict)
		if ok {
			subDst := make(DictMap)
			sd.Map(subDst)
			dst[kv.Key] = subDst
		} else {
			dst[kv.Key] = kv.Value
		}
	}
}

// Parse convert map to Dict.
func (d *Dict) Parse(src DictMap) {
	d.reset()

	for k, v := range src {
		sv, ok := v.(map[string]interface{})
		if ok {
			subDict := new(Dict)
			subDict.Parse(sv)
			d.append(k, subDict)
		} else {
			d.append(k, v)
		}
	}
}
