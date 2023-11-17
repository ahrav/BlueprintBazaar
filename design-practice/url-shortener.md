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
- **Total Storage Requirement**: Approximately 6,000,000,000,000 bytes (or 6 TB).

### 2. Queries per Second (QPS)

### Daily Active Users (DAU) and Read:Create Ratio
- **DAU**: 100 million
- **Read:Create Ratio**: 100:1

### Daily Queries
- **Total Daily Queries**: DAU × (Reads + Creates)
- **Total Daily Queries**: 100 million DAU × (100 reads + 1 create)
- **Total Daily Queries**: 100 million DAU × 101

### QPS Calculation
- **Seconds in a Day**: 86,400 seconds
- QPS = Total Daily Queries / Seconds in a Day
- **Queries Per Second (QPS)**: Approximately 116,898 QPS.

### 3. Bandwidth

### Bandwidth for Read and Write Operations
- **Read Bandwidth**: 100 reads per user per day × 500 bytes × 100 million DAU
- **Write Bandwidth**: 1 create per user per day × 500 bytes × 100 million DAU

### Total Daily Bandwidth
- Total Daily Bandwidth = Read Bandwidth + Write Bandwidth
- **Average Bandwidth per Second**: Total Daily Bandwidth / Seconds in a Day
- **Average Bandwidth per Second**: Approximately 58,449,074 bytes/sec (or about 55.7 MB/sec).

### 4. Cache Size Calculation

### Caching the Top 20% URLs
- **Top 20% URLs**: 20% of 12 billion URLs (calculated for 5 years)
- **Cache Refresh Rate**: Daily

### Cache Size
- Cache Size = Top 20% URLs × Size per URL
- **Cache Size for Top 20% URLs**: Approximately 1,200,000,000,000 bytes (or 1,200 TB).

