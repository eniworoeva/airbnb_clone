# Airbnb Clone - System Design Document

## Table of Contents
1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Database Design](#database-design)
4. [API Design](#api-design)
5. [Important Questions](#important-questions)
6. [Conclusion](#conclusion)


## Overview

This document outlines the system design for a scalable Airbnb-like backend system designed to support up to 10M users initially, with a roadmap to scale to 100M users.

### Requirements
- **Functional**: User management, property listings, booking system, reviews
- **Non-functional**: 10M users, high availability, low latency, data consistency

### Key Metrics
- **Target TPS**: 5,000-10,000 transactions per second
- **Response Time**: < 200ms for 95th percentile
- **Availability**: 99.9% uptime
- **Data Consistency**: Strong consistency for critical operations

## Architecture

### High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐  
│   Load Balancer │    │   Client        │ 
│                 │◄───┤                 │
└─────────────────┘    └─────────────────┘
          │
          ▼
┌────────────────────────────────────────────────────────────────────────────────┐
│                    Application Layer                                           │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐   ┌───────────────┐   │
│  │   User API    │  │ Property API  │  │  Booking API  │   |  Review API   |   │
│  │   Service     │  │   Service     │  │   Service     │   |    Service    │   |
│  └───────────────┘  └───────────────┘  └───────────────┘   └───────────────┘   │
└────────────────────────────────────────────────────────────────────────────────┘
          │
          ▼
┌────────────────────────────────────────┐
│             Database Layer             │
│  ┌───────────────┐  ┌───────────────┐  │
│  │   PostgreSQL  │  │     Redis     │  │
│  │   (Primary)   │  │   (Cache)     │  │
│  └───────────────┘  └───────────────┘  │
└────────────────────────────────────────┘
```

### Clean Architecture Layers

1. **API Layer**: HTTP handlers, routing, middleware
2. **Service Layer**: Business logic, validation, orchestration
3. **Repository Layer**: Data access abstraction
4. **Infrastructure Layer**: Database, external services

## Database Design

### Entity Relationship Diagram

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│    Users    │     │  Properties  │     │  Bookings   │
│─────────────│     │──────────────│     │─────────────│
│ id (PK)     │     │ id (PK)      │     │ id (PK)     │
│ email       │────┐│ host_id (FK) │────┐│ property_id │
│ password    │    ││ title        │    ││ guest_id    │
│ first_name  │    ││ description  │    ││ check_in    │
│ last_name   │    ││ type         │    ││ check_out   │
│ role        │    ││ price        │    ││ status      │
│ created_at  │    ││ location     │    ││ total_price │
└─────────────┘    │└──────────────┘    │└─────────────┘
                   │                    │
                   │ ┌─────────────┐    │
                   └─┤   Reviews   │────┘
                     │─────────────│
                     │ id (PK)     │
                     │ property_id │
                     │ booking_id  │
                     │ reviewer_id │
                     │ rating      │
                     │ comment     │
                     └─────────────┘
```

### Key Database Design Decisions

1. **PostgreSQL**: Chosen for ACID compliance and complex queries
2. **UUID Primary Keys**: Better for distributed systems
3. **Indexes**: Strategic indexing for common query patterns
4. **Soft Deletes**: Maintain data integrity and audit trails


## API Design

### RESTful Principles
- Resource-based URLs
- HTTP methods for operations
- Stateless requests
- Standard HTTP status codes

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (Guest, Host, Admin)
- Token refresh mechanism
- Secure password storage (bcrypt)

### Key Endpoints

#### Property Search (Critical Path)
```
GET /api/v1/properties/search?city=NYC&check_in=2024-01-01&check_out=2024-01-05&guests=2
```

**Performance Optimizations**:
- Database query optimization
- Result caching
- Pagination
- Relevant indexes


## Important Questions


### How Do You Increase Capacity/Scale To 100m Users?

#### 1. Horizontal Scaling

- **Application Layer**
   ```
   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
   │   App       │    │   App       │    │   App       │
   │ Instance 1  │    │ Instance 2  │    │ Instance N  │
   └─────────────┘    └─────────────┘    └─────────────┘
            │                 │                 │
            └─────────────────┼─────────────────┘
                              │
                    ┌─────────────┐
                    │ Load        │
                    │ Balancer    │
                    └─────────────┘
   ```

- **Database Scaling**
   ```
   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
   │  Primary    │    │ Read        │    │ Read        │
   │  Database   │───▶│ Replica 1   │    │ Replica N   │
   │ (Write)     │    │ (Read)      │    │ (Read)      │
   └─────────────┘    └─────────────┘    └─────────────┘
   ```

#### 2. Microservices Architecture

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   User      │  │  Property   │  │  Booking    │  │   Review    │
│  Service    │  │  Service    │  │  Service    │  │  Service    │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
       │                │                │                │
       └────────────────┼────────────────┼────────────────┘
                        │                │
              ┌─────────────┐    ┌─────────────┐
              │  Message    │    │   Event     │
              │  Queue      │    │   Store     │
              └─────────────┘    └─────────────┘
```

#### 3. Caching Strategy

- **Redis Layers**:
   - Session cache (L1)
   - Query result cache (L2)
   - Static data cache (L3)

- **CDN**: Static assets and API responses

- **Application Cache**: In-memory caching for frequently accessed data

#### 4. Database Sharding

- **Geographic Sharding**: By user location
- **Functional Sharding**: By feature (users, properties, bookings)
- **Hash-based Sharding**: By user ID


### What's The TPS The System Will support? And How Can This Be Measured
- **Assumptions**: 10M total users
- **Assumptions**: 10% daily active users = 1M DAU
- **Daily Requests**: 20 requests per user per day = 20M requests/day
- **Traffic Factor**Peak traffic factor = 3x
- **Base TPS**: 20M ÷ (24 × 3600) = 231 TPS
- **Peak TPS**: 231 × 3 = 694 TPS
- **With optimization**: 5,000-10,000 TPS achievable




### How Do You Increase Performance Without Increasing Cost?

1. **Efficient Queries**: Optimize database queries to reduce compute time
2. **Smart Caching**: Reduce database load through strategic caching
3. **Resource Right-sizing**: Match resources to actual usage patterns
4. **Auto-scaling**: Scale down during low traffic periods
5. **Reserved Instances**: Long-term commitments for predictable workloads


## Conclusion

This system design provides a solid foundation for scaling from 10M to 100M users through:

1. **Clean Architecture**: Maintainable and testable code
2. **Strategic Scaling**: Horizontal scaling with microservices
3. **Performance Optimization**: Database indexing and caching
4. **Monitoring**: Comprehensive observability stack