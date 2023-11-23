## Web Crawler: Scope, Requirements, and Purpose

### Purpose of a Web Crawler
1. **Data Collection**: To gather information from various websites, which can be used for indexing content in search engines, data mining, and web archiving.
2. **Indexing for Search Engines**: To enable search engines to provide relevant results. Crawlers collect and process data from web pages,
which is then indexed by the search engine.
3. **Website Health Checks**: To monitor the health of websites by checking links, validating HTML/CSS, and ensuring content adheres to SEO guidelines.
4. **Web Archiving**: To capture and store historical versions of web pages.

### Problem Scope
1. **Scalability**: Should be capable of handling large volumes of web pages and data.
2. **Efficiency**: Must efficiently navigate and download content from websites without causing undue load on web servers.
3. **Robustness**: Should be able to handle different website layouts, content types, and recover from errors.
4. **Respectful Crawling**: Adherence to `robots.txt` guidelines and ethical crawling practices to avoid overloading servers.

### Functional Requirements
1. **URL Discovery**: Identifying URLs to visit through links found in previously crawled content.
2. **Content Downloading**: Fetching the content of the web pages.
3. **Content Processing**: Extracting useful information from the downloaded content, such as text, links, and metadata.
4. **Link Extraction**: Identifying links to other pages for further crawling.
5. **Storage**: Storing the downloaded content and extracted information.
6. **Duplication Avoidance**: Avoiding the re-crawling of identical content.
7. **Adherence to Protocols**: Respecting `robots.txt` and meta `robots` tags, following sitemap directives, and managing crawl rates.
8. **Update Frequency**: Deciding how often to revisit sites to update the indexed content.
9. **Crawl Scheduling**: Scheduling the crawling process to ensure efficient use of resources.

### Non-Functional Requirements
1. **Performance**: High-speed crawling and data processing capabilities.
2. **Scalability**: Ability to scale up with the increasing number of web pages and sites.
3. **Politeness**: Crawling websites without impacting their normal operation.
4. **Distributed Architecture**: For improved performance and fault tolerance.
5. **Data Management**: Efficient storage and retrieval of crawled data.
6. **Customizability**: Flexibility to adapt to different crawling objectives.


## Back-of-the-Envelope Calculations

*Total Webpages:*
* 20 billion webpages total
* So if each webpage takes up 2MB storage when saved in the database, total storage needed is:
* 20 billion * 2MB per page = 40,000,000 GB = 40,000 TB = 40PB storage

*Crawling Rate:*
* 500 million webpages crawled per day
* So write throughput needed is: 500 million writes per day
* Converting this to per second rate:
* 500 million / (100,000 sec/day) = 5000 writes per second (WPS)

*User Traffic:*
* 100,000 queries per second (QPS)

*Read Bandwidth:*
* 100,000 QPS
* 500 bytes per query
* 500 bytes = 4000 bits per query
* 100,000 queries * 4000 bits per query = 400,000,000 bits per second
* 400,000,000 bits per second = 400 Megabits per second = 400 Mbps

*Write Bandwidth:*
* 5000 writes per second (WPS)
* 2 MB per write
* 2 MB = 16,000,000 bits
* 16,000,000 * 5000 writes per second = 80,000,000,000 bits per second = 80 Gbps

In summary the key metrics are:
* 40PB storage
* 5K WPS write throughput
* 100K QPS read throughput
* 80 Gbps write bandwidth
* 400 Mbps read bandwidth

## High-Level Components of a Web Crawler System

### URL Frontier
- **Function**: Acts as the queue of URLs to be visited by the crawler.
- **Purpose**: Prioritizes and schedules URLs for crawling based on predefined rules or algorithms.

### Scrapers
- **Function**: Responsible for downloading web pages.
- **Purpose**: Extracts content and metadata from web pages, then passes this data to other components for processing.

### DNS Resolver
- **Function**: Resolves domain names to IP addresses.
- **Purpose**: Essential for the scraper to connect to web servers and download content.

### Dedupe Service (De-duplication Service)
- **Function**: Detects and removes duplicate content.
- **Purpose**: Ensures the crawler does not process the same content multiple times, improving efficiency.

### Parser
- **Function**: Extracts useful information from the downloaded content.
- **Purpose**: Enables the crawler to understand the content and extract relevant data.

### Link Extractor
- **Function**: Identifies links to other pages for further crawling.
- **Purpose**: Enables the crawler to navigate the web and discover new URLs.

### Message Broker
- **Function**: Manages communication between different components of the system.
- **Purpose**: Facilitates asynchronous processing and decouples the system components for scalability and reliability.

### Object Stores
- **Function**: Stores the crawled content and metadata.
- **Purpose**: Provides scalable and reliable storage for large amounts of data.

### Graph Database for Link Associations
- **Function**: Stores and manages the relationships between web pages.
- **Purpose**: Useful for understanding the link structure of the web and for algorithms like PageRank.

### PageRank Workers
- **Function**: Executes the PageRank algorithm or other relevance-scoring algorithms.
- **Purpose**: Ranks pages based on their importance and relevance, which is crucial for search engine results.

### Recrawl Service
- **Function**: Manages the re-visiting of web pages.
- **Purpose**: Ensures that the crawler revisits pages to update the information and capture changes over time.

### Cache
- **Function**: Temporarily stores frequently accessed data.
- **Purpose**: Reduces the load on primary storage and improves response times for repeated requests.

### Bloom Filters
- **Function**: Determines if a URL has already been crawled.
- **Purpose**: Improves the efficiency of the dedupe service by filtering out URLs that have already been crawled.
