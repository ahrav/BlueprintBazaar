# Notification Service

## Purpose

The purpose of this system is to provide a scalable, reliable, and efficient way to handle notifications for clients,
prioritizing them accordingly, and delivering them through various channels such as SMS, email, and push notifications.

## Problem Scope

The system is designed to manage a high volume of incoming notification requests, prioritize them,
and ensure timely delivery while maintaining user preferences and rate limiting to prevent spam.

## Functional Requirements

- **Notification Prioritization**: Classify notifications into high and low priority queues.
- **Rate Limiting**: Prevent notification spam by enforcing rate limits.
- **User Preference Management**: Allow users to set notification preferences.
- **Delivery Channels**: Support multiple delivery mechanisms including SMS, email, and push notifications.
- **Tracking**: Keep track of notification delivery status.

## Non-Functional Requirements

- **Scalability**: System must handle an increasing number of notifications.
- **Reliability**: Notifications should be delivered reliably to the intended recipients.
- **Performance**: The system should have low latency and high throughput.
- **Availability**: The system should be highly available and resilient to failures.

## High-Level Components

### FE Service
- **Function**: Serves as the entry point for incoming notification requests from various clients.
- **Purpose**: Validates and forwards requests to the appropriate internal services, ensuring a decoupled architecture for scalability.

### Message Broker
- **Function**: Handles message queuing, distribution, and decoupling of services.
- **Purpose**: Provides reliable message delivery, facilitates asynchronous processing, and helps in managing load across services.

### Prioritizers
- **Function**: Sorts incoming notifications into priority queues based on defined criteria.
- **Purpose**: Ensures that high-priority notifications are processed and delivered faster than lower-priority ones.

### Rate Limiter
- **Function**: Enforces limits on how many notifications a user can send or receive within a certain timeframe.
- **Purpose**: Prevents spam and abuse of the notification system, maintaining quality of service for all users.

### Notification Assemblers
- **Function**: Constructs notification messages from templates or parameters, readying them for delivery.
- **Purpose**: Personalizes notifications according to user preferences and prepares them in the correct format for each delivery channel.

### SMS/Email/Push Brokers
- **Function**: Specialized services that handle the delivery of notifications through their respective channels.
- **Purpose**: Isolates the complexity of delivering messages across different platforms and providers, offering flexibility and reliability.

### Databases
- **Function**: Stores data such as user preferences, notification logs, and system metrics.
- **Purpose**: Enables persistence, retrieval, and analysis of data crucial for the operation and monitoring of the service.

#### Document DB
- **Specific Function**: Stores unstructured data, possibly the content of the notifications themselves.
- **Specific Purpose**: Allows for flexible schema design and quick retrieval of notification content.

#### Metrics DB
- **Specific Function**: Collects and stores metrics about system performance and usage.
- **Specific Purpose**: Facilitates monitoring, alerting, and performance tuning.

#### User Pref DB
- **Specific Function**: Maintains user profiles and notification preferences.
- **Specific Purpose**: Ensures that notifications are delivered according to user-specified settings.

#### Cassandra
- **Specific Function**: Provides a highly available and scalable database for storing large volumes of notification logs.
- **Specific Purpose**: Allows for efficient writes and reads, especially suited for time-series data like notification events.

### Redis
- **Function**: In-memory data structure store used as a database, cache, and message broker.
- **Purpose**: Provides fast access to data for high-performance operations, often used for transient data such as session states or intermediate message queuing.

### Notification Tracking Service
- **Function**: Tracks the delivery status and engagement of notifications sent to users.
- **Purpose**: Provides insights into notification effectiveness and ensures that the delivery of notifications can be audited and verified.
