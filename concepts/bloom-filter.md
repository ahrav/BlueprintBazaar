## Bloom Filters

### What is a Bloom Filter?
- **Definition**: A Bloom filter is a data structure designed to quickly and efficiently check whether an element is present in a set.
It's probabilistic, which means it can tell with certainty if an element is not in the set but only probabilistically if it is.
- **Space Efficiency**: It is highly space-efficient compared to other data structures but at the cost of a certain probability of false positives.

### Why are Bloom Filters Used?
- **Speed and Memory Efficiency**: They are incredibly fast and use significantly less memory than other structures, making them ideal for large datasets.
- **Network Efficiency**: In distributed systems, they reduce the need to query the main dataset by filtering out unnecessary queries.
- **False Positives Acceptable**: Useful in situations where false positives are acceptable but false negatives are not.

### How Does a Bloom Filter Work?
1. **Bit Array**: At its core, a Bloom filter is an array of bits, all set to 0 initially.
2. **Hash Functions**: Several hash functions are used to map elements of the set to positions in the bit array.
3. **Insertions**: When an element is added, the bits at positions indicated by the hash functions are set to 1.
4. **Lookups**: To check if an element is in the set, it is hashed with the same hash functions. If all the bits at these positions are 1, the element is *probably* in the set; if any bit is 0, it's definitely not in the set.

### Example: Using a Bloom Filter

#### Scenario
- We have a Bloom filter with a bit array of size 10 (indices 0-9).
- Two hash functions are used:
    - `hash1(x) = (x * 2) % 10`
    - `hash2(x) = (x * 3) % 10`
- Initially, all bits in the array are 0.

#### Adding Elements
- **Adding Element 1**:
    - `hash1(1) = 2`, `hash2(1) = 3`
    - Set bits at indices 2 and 3 to 1.
    - Bit array becomes `[0, 0, 1, 1, 0, 0, 0, 0, 0, 0]`.
- **Adding Element 2**:
    - `hash1(2) = 4`, `hash2(2) = 6`
    - Set bits at indices 4 and 6 to 1.
    - Bit array becomes `[0, 0, 1, 1, 1, 0, 1, 0, 0, 0]`.

#### Lookup
- **Checking for Element 5**:
    - `hash1(5) = 0`, `hash2(5) = 5`
    - Check bits at indices 0 and 5. Index 0 is 0, so element 5 is not in the set.
- **Checking for Element 1** (already added):
    - `hash1(1) = 2`, `hash2(1) = 3`
    - Both bits at indices 2 and 3 are 1, so element 1 is probably in the set.

### Explanation
- The Bloom filter uses hash functions to map elements to the bit array. When adding, it sets the bits at these positions to 1. During a lookup, if all bits are 1, the element is probably in the set; if any bit is 0, the element is not in the set.
- **False Positives**: There can be false positives (being told an element is in the set when it's not), but no false negatives.

### Real-World Systems Using Bloom Filters
- **Network Routers**: For quickly checking if a URL is in a list of malicious websites.
- **Databases**: Cassandra and Redis use Bloom filters to check if a key is in a dataset before querying the disk.
- **Web Browsers**: For checking if a URL is in a list of known malicious web addresses.
- **Big Data Applications**: Apache HBase and Apache Spark use Bloom filters for data lookups to reduce disk I/O.
