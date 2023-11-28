package implementations

import (
	"hash"
	"hash/fnv"
	"math"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/dchest/siphash"
	"github.com/spaolacci/murmur3"
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
	// HashFn returns a hash function.
	HashFn() hash.Hash64
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
	// hasher is the hash function used to map elements to bits in the bit array.
	// Default is xxhash.
	hasher Hasher
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
		Hasher:            sipHasher{},
	}
	for _, opt := range opts {
		opt(c)
	}

	m := calculateBitArraySize(c.Capacity, c.FalsePositiveRate)
	k := calculateHashesCount(c.Capacity, c.FalsePositiveRate)
	return &BasicBloomFilter{
		m:      m,
		bitArr: make([]uint64, (m+63)/64),
		k:      uint64(k),
		hasher: c.Hasher,
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
		h := b.hasher.HashFn()
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
		h := b.hasher.HashFn()
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

const (
	xxhashName = "xxhash"
	fnvaName   = "fnva"
	murmurName = "murmur"
)

type sipHasher struct{}

func (sipHasher) HashFn() hash.Hash64 {
	zeroKey := make([]byte, 16)
	return siphash.New(zeroKey)
}

func (sipHasher) Name() string {
	return "siphash"
}

type murmurHasher struct{}

func (murmurHasher) HashFn() hash.Hash64 {
	return murmur3.New64()
}

func (murmurHasher) Name() string {
	return murmurName
}

type fnvAHasher struct{}

func (fnvAHasher) HashFn() hash.Hash64 {
	return fnv.New64a()
}

func (fnvAHasher) Name() string {
	return fnvaName
}

type xxHasher struct{}

func (xxHasher) HashFn() hash.Hash64 {
	return xxhash.New()
}

func (xxHasher) Name() string {
	return xxhashName
}
