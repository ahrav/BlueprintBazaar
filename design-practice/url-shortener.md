# Designing a URL Shortener

i# Understanding the Scope and Requirements of a URL Shortener

## Purpose of a URL Shortener

A URL shortener is a service that converts long URLs into a shortened version, making them easier to share and manage.
The primary purposes are:

- **Ease of Sharing**: Short URLs are more user-friendly for sharing on platforms like social media, where character space is limited.
- **Improved Aesthetics**: Shorter URLs look cleaner and are more appealing to users.
- **Tracking and Analytics**: They can be used to track clicks, geographic locations of users, and other analytics data.
- **Redirection Control**: They allow the management of the redirection process, which can be useful for marketing and control of web traffic.

## Scope and Functional Requirements

### Core Features

1. **URL Shortening**: Convert a long URL to a short, unique identifier.
2. **Redirection**: When a user accesses the shortened URL, they should be redirected to the original URL.
3. **Deletion of URLs**: Allow users to delete shortened URLs, removing them from the system.
4. **Scalability**: The system should handle a high volume of requests.
5. **Performance**: Fast response times for both creating the short URL and redirecting to the original URL.
6. **Expiration**: Allow users to set expiration dates for shortened URLs.

### Additional Features

1. **Custom Short URLs**: Allow users to create custom short links.
2. **Analytics**: Provide analytics for each shortened URL (e.g., number of clicks, geographic location of the users, etc.).
3. **Expiration**: Option to set expiration dates for URLs.
4. **Security**: Measures to prevent abuse, such as rate limiting and filtering of malicious URLs.

### Non-functional Requirements

1. **High Availability**: The system should be highly available and reliable.
2. **Low Latency**: Both the URL shortening and redirection processes should have low latency.
3. **Data Integrity**: Ensure that the mapping between the original URL and the shortened URL remains consistent and accurate.
4. **Scalability**: The system should scale horizontally to accommodate growth.

### Constraints

1. **URL Character Limit**: The shortened URL should have a defined maximum character limit.
2. **Storage**: Efficient storage of URL mappings.
3. **Rate Limiting**: To prevent abuse, rate limiting may be necessary for URL creation requests.

## Back-of-the-Envelope Calculations

### 1. Total Storage Needed for 5 Years

### URLs Creation Rate
- **New URLs per Month**: 200 million
- **Time Period**: 5 years = 60 months

### Total URLs in 5 Years
- Total URLs = 200 million URLs/month × 60 months = 12 billion URLs

### Size per URL
- **Size per URL**: 500 bytes

### Total Storage Requirement
- Total Storage = Total URLs × Size per URL
- Total Storage = 12 billion URLs × 500 bytes/URL
- **Total Storage Requirement**: Approximately 6 TB.

### 2. Queries per Second (QPS) Calculation

### URLs Creation and Read Rate
- **New URLs per Month**: 200 million
- **Read:Create Ratio**: 100:1

### Monthly Reads
- **Monthly Reads**: 20 billion reads/month (200 million URLs × 100 reads/URL)

### Monthly Writes
- **Monthly Writes**: 200 million writes/month

### Total Monthly Queries
- Total Monthly Reads ≈ 20 billion
- Total Monthly Writes ≈ 200 million

### QPS Calculation for Reads
- QPS for Reads ≈ 20 billion / 2.5 million ≈ ~8K QPS

### QPS Calculation for Writes
- QPS for Writes ≈ 200 million / 2.5 million ≈ ~80 QPS

### 3. Bandwidth

### Bandwidth for Read Operations
- **Read Bandwidth**: 20 billion reads/month × 500 bytes

### Bandwidth for Write Operations
- **Write Bandwidth**: 200 million writes/month × 500 bytes

### Total Monthly Bandwidth
- Total Monthly Bandwidth ≈ Read Bandwidth + Write Bandwidth

### Average Bandwidth per Second for Reads
- **Average Bandwidth per Second for Reads**: ≈ 3.9 MB/sec

### Average Bandwidth per Second for Writes
- **Average Bandwidth per Second for Writes**: Smaller fraction of the total, due to fewer writes

### 4. Cache Size Calculation
#### Daily Reads
- **Daily Reads**: 7,700 QPS × 86,400 seconds/day ≈ 650 million reads/day

#### Top 20% of Daily Reads for Caching
- **Top 20% of Daily Reads**: 20% of 650 million ≈ 130 million reads/day

#### Cache Size
- **Cache Size for Top 20% Daily Reads**: 130 million reads/day × 500 bytes/read
- Cache Size ≈ 65 GB/day

## High-Level API Design

### Creating a Short URL from a Long URL

#### Endpoint
- `POST /create`

#### Request Parameters
- `longUrl`: The original long URL that needs to be shortened.

#### Response
- `shortUrl`: The shortened URL.
- `status`: Success or error message.

#### Description
This endpoint takes a long URL and returns a shortened version of it. It generates a unique identifier for the long URL,
stores this mapping in the database, and returns the short URL.

### Redirecting from a Short URL to a Long URL

#### Endpoint
- `GET /{shortUrl}`

#### Path Parameters
- `shortUrl`: The shortened URL.

#### Response
- Redirects to the original long URL.
- In case of an error or if the URL is not found, returns an appropriate error message.

#### Description
When a user accesses a short URL, this endpoint redirects them to the original long URL.
It looks up the short URL in the database to find the corresponding long URL and then performs the redirection.

### Deleting a Short URL

#### Endpoint
- `DELETE /delete`

#### Request Parameters
- `shortUrl`: The shortened URL to be deleted.

#### Response
- `status`: Success or error message.

#### Description
This endpoint allows for the deletion of a short URL. It removes the mapping of the short URL to its original URL from the database.
This is useful for managing the lifecycle of URLs and removing URLs that are no longer needed or violate certain policies.

## Key Components Overview

### Load Balancer
- **Purpose**: Distributes incoming network traffic across multiple servers to ensure no single server bears too much demand.
- **Benefits**: Enhances application responsiveness and availability.

### Rate Limiter
- **Purpose**: Limits the number of requests a user can make within a given time frame. Essential for preventing abuse and overloading the system.
- **Benefits**: Protects against DDoS attacks and ensures fair resource usage.

### Web Server
- **Purpose**: Handles HTTP requests from clients and provides responses. Acts as the entry point for the API.
- **Benefits**: Processes client requests efficiently, ensuring quick responses.

### Application Server
- **Purpose**: Hosts the application logic for URL shortening, redirecting, and deletion.
- **Benefits**: Separates business logic from client-server communication, enhancing maintainability and scalability.

### Cache
- **Purpose**: Temporarily stores the most frequently accessed data (e.g., top 20% of URLs).
- **Benefits**: Reduces database load and improves response time for popular URLs.

### Database for High Read Throughput
- **Purpose**: Stores mappings of short URLs to long URLs. Optimized for high read throughput.
- **Benefits**: Ensures quick retrieval of URL mappings, crucial for redirection.

### Sequencer
- **Purpose**: Generates unique identifiers for shortened URLs.
- **Benefits**: Ensures each short URL is unique and can be reliably traced back to its original URL.

### Encoder
- **Purpose**: Converts the unique identifier generated by the sequencer into a short URL.
- **Benefits**: Provides the final shortened URL in a user-friendly format.

## Using a Base58 Encoder in the URL Shortener System

### Purpose of Base58 Encoder
- **Main Function**: The encoder converts a unique identifier (generated by the sequencer) into a short URL.

### Why Base58?
- **Character Set**: Base58 uses a set of 58 characters, which includes the numbers 1-9, the letters A-Z and a-z,
excluding characters that can be easily mistaken for each other such as '0' (zero), 'O' (uppercase 'o'), 'I' (uppercase 'i'), and 'l' (lowercase 'L').
- **Advantages**:
    - **Readability and User-Friendliness**: By avoiding characters that look similar, base58 reduces the risk of misreading or mistyping URLs.
This is crucial for user experience, especially when URLs are manually entered or read aloud.
    - **Compactness**: Base58 encoding is more compact than hexadecimal (base16) or base64. It allows for a shorter URL,
which is a key objective of a URL shortener.
    - **No Special Characters**: Unlike base64, base58 doesn't include special characters like '+' or '/'.
This makes it suitable for URLs, where special characters can sometimes cause issues in browsers or when copied and pasted across different platforms.

### Usage in URL Shortener
- **URL Generation**: The sequencer generates a unique numeric identifier for each new URL using a snowflake ID generator. 
Using a snowflake ID generator ensures that the IDs are unique across different servers and can be generated quickly.
The base58 encoder then converts this number into a short string, forming the path of the shortened URL.
- **Scalability and Uniqueness**: With a pool of 58 characters, base58 provides a vast number of combinations,
ensuring that unique URLs can be generated for a large volume of entries, which is essential for scalability and uniqueness in the system.
