package implementations

import (
	"hash"
	"math"
	"sync"

	"github.com/cespare/xxhash/v2"
)

// BloomFilter provides a common interface for bloom filter implementations.
type BloomFilter interface {
	// Add a value to the bloom filter.
	Add(value []byte) error
	// Test if a value is in the bloom filter.
	Test(value []byte) (bool, error)
}

// Hasher provides a common interface for hash functions.
type Hasher interface {
	// Hashes returns a slice of n hash functions.
	Hashes(n uint64) []hash.Hash64
	// Name of the hash function.
	Name() string
}

// BasicBloomFilter is a basic bloom filter implementation.
type BasicBloomFilter struct {
	m  uint64 // Number of bits in the bloom filter.
	mu sync.RWMutex
	// Bit array representing set membership of elements.
	bitArr []uint64
	k      uint64 // Number of hash functions.
	// Hash functions used to map elements to bits in the bit array.
	hashFns []hash.Hash64
}

// Config needed to create a new bloom filter.
type Config struct {
	// Capacity is the number of elements expected to be added to the bloom filter.
	// Default is 100_000.
	Capacity uint64
	// HashesCount is the number of hash functions used to map elements to bits in the bit array.
	// Default is calculated based on the desired false positive rate.
	// Only override if you know what you are doing.
	HashesCount uint16
	// FalsePositiveRate is the desired false positive rate.
	// Default is 0.01.
	FalsePositiveRate float64
	// Hasher is the hash function used to map elements to bits in the bit array.
	// Default is xxhash.
	Hasher Hasher
}

// BloomFilterOption is used to configure a new bloom filter.
type BloomFilterOption func(*Config)

// WithCapacity sets the capacity of the bloom filter.
func WithCapacity(capacity uint64) BloomFilterOption {
	return func(c *Config) {
		c.Capacity = capacity
	}
}

// WithHashesCount sets the number of hash functions used to map elements to bits in the bit array.
func WithHashesCount(count uint16) BloomFilterOption {
	return func(c *Config) {
		c.HashesCount = count
	}
}

// WithFalsePositiveRate sets the desired false positive rate of the bloom filter.
func WithFalsePositiveRate(rate float64) BloomFilterOption {
	return func(c *Config) {
		c.FalsePositiveRate = rate
	}
}

// WithHasher sets the hash function used to map elements to bits in the bit array.
func WithHasher(hasher Hasher) BloomFilterOption {
	return func(c *Config) {
		c.Hasher = hasher
	}
}

// NewBasicBloomFilter creates a new basic bloom filter.
func NewBasicBloomFilter(opts ...BloomFilterOption) *BasicBloomFilter {
	const (
		defaultCapacity          = 100_000
		defaultFalsePositiveRate = 0.01
	)
	c := &Config{
		Capacity:          defaultCapacity,
		FalsePositiveRate: defaultFalsePositiveRate,
		Hasher:            xxHasher{},
	}
	for _, opt := range opts {
		opt(c)
	}

	m := calculateBitArraySize(c.Capacity, c.FalsePositiveRate)
	k := calculateHashesCount(c.Capacity, c.FalsePositiveRate)
	return &BasicBloomFilter{
		m:       m,
		bitArr:  make([]uint64, (m+63)/64),
		k:       uint64(k),
		hashFns: c.Hasher.Hashes(uint64(k)),
	}
}

// calculateHashesCount calculates the number of hash functions needed to achieve the desired false positive rate.
// p = (1 - e^(-kn/m))^k
func calculateHashesCount(cap uint64, fpRate float64) uint16 {
	return uint16(math.Ceil(-math.Log(fpRate) / math.Log(2)))
}

func calculateBitArraySize(cap uint64, fpRate float64) uint64 {
	return uint64(math.Ceil(-float64(cap) * math.Log(fpRate) / math.Pow(math.Log(2), 2)))
}

// Add a value to the bloom filter.
func (b *BasicBloomFilter) Add(value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := 0; i < int(b.k); i++ {
		h := xxhash.New()
		h.Reset()
		h.Write(value)
		b.setBit(h.Sum64())
	}
	return nil
}

func (b *BasicBloomFilter) setBit(sum uint64) {
	b.bitArr[(sum%b.m)/64] |= 1 << (sum % 64)
}

// Test if a value is in the bloom filter.
func (b *BasicBloomFilter) Test(value []byte) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for i := 0; i < int(b.k); i++ {
		h := xxhash.New()
		h.Reset()
		h.Write(value)
		if !b.isBitSet(h.Sum64()) {
			return false, nil
		}
	}
	return true, nil
}

func (b *BasicBloomFilter) isBitSet(sum uint64) bool {
	return b.bitArr[(sum%b.m)/64]&(1<<(sum%64)) != 0
}

const xxhashName = "xxhash"

type xxHasher struct{}

func (xxHasher) Hashes(n uint64) []hash.Hash64 {
	hashes := make([]hash.Hash64, n)
	for i := range hashes {
		hashes[i] = xxhash.New()
	}
	return hashes
}

func (xxHasher) Name() string {
	return xxhashName
}
