# Consistent Hashing

## What is Consistent Hashing?

Consistent hashing is a technique used to distribute requests or data across a cluster of servers in a way that minimizes
reorganization when servers are added or removed. It involves mapping each object to a point on a circle (or hash ring)
and each server to a point on the same circle.

# Key Concepts of Consistent Hashing

Before exploring the process of consistent hashing, it's essential to understand some key concepts that form the basis of this technique.
These concepts are crucial for appreciating how consistent hashing achieves its goals in distributed systems.

## Key Concepts

### 1. **Hash Ring or Hash Space**
- A consistent hashing algorithm maps both servers (nodes) and data to a continuous ring-like space, often referred to as the hash ring.
- This ring represents a range of hash values from the hash function used.

### 2. **Hash Function**
- A hash function is used to assign each server and each data item a position on the hash ring.
- The choice of hash function can impact the uniformity of data and load distribution.

### 3. **Virtual Nodes**
- To improve the distribution uniformity, each server can be represented by multiple points (virtual nodes) on the hash ring.
- This approach helps in balancing the load more evenly among servers.

### 4. **Data Distribution**
- Data items are distributed across servers based on their position in the hash space.
- Each piece of data is managed by the nearest server on the ring in the clockwise direction.

### 5. **Scalability and Elasticity**
- Consistent hashing is inherently designed to handle changes in the cluster size (adding or removing servers) with minimal impact on the overall system.

### 6. **Fault Tolerance and Redundancy**
- The technique supports replicating data across multiple nodes to ensure high availability and fault tolerance.

Understanding these key concepts is vital for grasping how consistent hashing provides efficient data distribution and load
balancing in a dynamic distributed system environment. Next, we'll look at the process overview to see consistent hashing in action.

### Process Overview

1. **Hash Function**: Both servers (or nodes) and data items are hashed using the same hash function.
2. **Mapping to Hash Ring**: The results are placed on the edge of a circle (hash ring).
3. **Data Assignment**: Each data item is assigned to the nearest server on the ring in the clockwise direction.

## Why Use Consistent Hashing?

### 1. **Minimizing Rehashing**
- When a server is added or removed, only a small portion of data needs to be reassigned, reducing the amount of data transfer and rehashing required.

### 2. **Load Balancing**
- It spreads data and load uniformly across servers, ensuring no single server becomes a bottleneck.

### 3. **Scalability**
- Makes it easier to scale a distributed system by adding or removing servers without a significant reorganization of data.

### 4. **High Availability**
- Facilitates replication and redundancy. If a server fails, its portion of the hash ring can be covered by other servers, ensuring continuous availability of data.

### 5. **Flexibility**
- Suitable for environments where the number of servers can change dynamically, such as cloud computing and peer-to-peer networks.

### 6. **Distributed Caching**
- Widely used in caching systems to distribute cache data across multiple cache servers and efficiently locate where data is stored.

# Consistent Hashing: Problems Solved and Practical Examples

Consistent hashing addresses specific challenges in distributed systems that simpler approaches like modulo hashing can't effectively handle.

## Problems Solved by Consistent Hashing

### 1. **Handling Server Changes**
- **Issue with Modulo Hashing**: In modulo hashing, adding or removing a server requires remapping a significant portion of the keys, leading to massive data reshuffling.
- **Solution in Consistent Hashing**: Minimizes the amount of remapping when servers are added or removed, reducing the data reshuffling burden.

### 2. **Uneven Load Distribution**
- **Issue with Modulo Hashing**: Can lead to uneven load distribution, especially with a small number of servers.
- **Solution in Consistent Hashing**: Offers more uniform load distribution, especially when combined with the concept of virtual nodes.

### 3. **Scalability and Flexibility**
- **Issue with Modulo Hashing**: Less flexible to scale in dynamic environments.
- **Solution in Consistent Hashing**: Better suited for environments where the number of servers can change frequently and dynamically.

## Examples of Consistent Hashing in Practice

### 1. **Distributed Caching Systems**
- **Example**: Memcached, Redis.
- **Use Case**: Distributing cache data across multiple cache servers; when a cache server is added or removed, only a small portion of the data needs to be moved.

### 2. **Distributed Databases**
- **Example**: Amazon DynamoDB, Cassandra.
- **Use Case**: In NoSQL databases, consistent hashing helps in distributing data across various nodes, facilitating scalability and high availability.

### 3. **Load Balancing**
- **Example**: Load balancers in cloud computing environments.
- **Use Case**: Distributing incoming requests or network connections across multiple servers to ensure balanced load handling.

### 4. **Content Delivery Networks (CDNs)**
- **Example**: Akamai, Cloudflare.
- **Use Case**: Efficiently caching and retrieving content from edge servers closest to the end-users.

### 5. **Peer-to-Peer (P2P) Networks**
- **Example**: BitTorrent's Distributed Hash Table (DHT).
- **Use Case**: Locating nodes in a P2P network and efficiently distributing data amongst peers.

By addressing the shortcomings of simpler hashing mechanisms and providing a more robust solution for dynamic and large-scale
distributed systems, consistent hashing has become a fundamental component in various modern applications and services.
