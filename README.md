# Microservices Architecture Documentation

## Table of Contents

1. [System Overview](#system-overview)
2. [Architecture Diagram](#architecture-diagram)
3. [Services Description](#services-description)
4. [Database Schema Design](#database-schema-design)
5. [Detailed Flow Processes](#detailed-flow-processes)
6. [API Specifications](#api-specifications)
7. [Event Schemas](#event-schemas)
8. [Security Implementation](#security-implementation)

---

## System Overview

A distributed microservices architecture built with GoLang implementing:

- JWT-based authentication
- gRPC for inter-service communication
- REST & GraphQL for client-facing APIs
- Kafka for event-driven architecture
- Multi-database strategy (PostgreSQL, MySQL, MongoDB, Redis)

### Technology Stack

- **Language**: GoLang
- **RPC Framework**: gRPC
- **API**: REST (Gin/Echo), GraphQL (gqlgen)
- **Message Broker**: Kafka
- **Databases**: PostgreSQL, MySQL, MongoDB, Redis
- **Authentication**: JWT (RS256)

---

## Architecture Diagram

```text
┌─────────────────────────────────────────────────────────────────┐
│                          API Gateway                             │
│                    (JWT Validation Layer)                        │
└───────────────────────┬─────────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┬──────────────┐
        │               │               │              │
        ▼               ▼               ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   Auth       │ │    User      │ │   Product    │ │   Order      │
│   Service    │ │   Service    │ │   Service    │ │   Service    │
│              │ │              │ │              │ │              │
│   REST       │ │ REST/GraphQL │ │ REST/GraphQL │ │   gRPC       │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │                │
       │                │                │                │
   ┌───▼────┐      ┌───▼────┐      ┌───▼────┐      ┌───▼────┐
   │ Postgre│      │ Postgre│      │ Postgre│      │  MySQL │
   │  SQL   │      │  SQL   │      │  SQL   │      │        │
   └────────┘      └───┬────┘      └────────┘      └────────┘
                       │
                   ┌───▼────┐
                   │ MongoDB│
                   └────────┘
       │                │                │                │
       └────────────────┴────────────────┴────────────────┘
                        │
                   ┌────▼─────┐
                   │  Redis   │
                   │ (Cache & │
                   │ Sessions)│
                   └──────────┘
                        │
        ┌───────────────┴───────────────┐
        ▼                               ▼
┌──────────────┐              ┌──────────────┐
│   Kafka      │              │   Kafka      │
│  Producer    │              │  Consumer    │
└──────────────┘              └──────┬───────┘
                                     │
                        ┌────────────┼────────────┐
                        ▼            ▼            ▼
                ┌──────────┐  ┌──────────┐  ┌──────────┐
                │Notification│ │Analytics │  │  Other   │
                │  Service   │ │ Service  │  │ Consumers│
                └─────┬──────┘ └────┬─────┘  └──────────┘
                      │             │
                  ┌───▼───┐     ┌───▼───┐
                  │MongoDB│     │MongoDB│
                  └───────┘     └───────┘
```

---

## Services Description

### 1. Authentication Service

**Port**: 8081  
**Protocol**: REST (Client), gRPC Server (Internal)  
**Database**: PostgreSQL (primary), Redis (sessions/tokens)

**Responsibilities**:

- User credential management
- JWT token generation and validation
- Refresh token management
- Token revocation/blacklisting
- Password hashing (bcrypt)

**Database Tables** (PostgreSQL):

- `users_auth` - authentication credentials
- `refresh_tokens` - refresh token tracking

**Redis Keys**:

- `session:{user_id}` - active sessions
- `blacklist:{token_id}` - revoked tokens
- `refresh:{token_id}` - refresh tokens (TTL: 7 days)

---

### 2. User Service

**Port**: 8082  
**Protocol**: REST + GraphQL (Client), gRPC Server (Internal)  
**Database**: PostgreSQL (primary), MongoDB (activity logs)

**Responsibilities**:

- User profile management (one-to-one: user → profile)
- Address management (one-to-many: user → addresses)
- User activity logging
- Profile photo management

**Database Tables** (PostgreSQL):

- `users` - basic user information
- `user_profiles` - detailed profile (one-to-one)
- `user_addresses` - multiple addresses per user (one-to-many)

**MongoDB Collections**:

- `user_activities` - activity logs with flexible schema
- `user_preferences` - JSON document storage

---

### 3. Product Service

**Port**: 8083  
**Protocol**: REST + GraphQL (Client), gRPC Server (Internal)  
**Database**: PostgreSQL (primary), Redis (cache)

**Responsibilities**:

- Product catalog management
- Category management
- Product-Category relationships (many-to-many)
- Inventory tracking
- Product search and filtering

**Database Tables** (PostgreSQL):

- `products` - product information
- `categories` - product categories
- `product_categories` - junction table (many-to-many)
- `inventory` - stock levels

**Redis Keys**:

- `product:{id}` - product cache (TTL: 1 hour)
- `category:{id}:products` - category products list
- `inventory:{product_id}` - real-time inventory

---

### 4. Order Service

**Port**: 8084  
**Protocol**: gRPC (All communications)  
**Database**: MySQL (primary), Kafka (events)

**Responsibilities**:

- Order creation and management
- Order-Items relationship (one-to-many)
- Payment processing coordination
- Order status tracking
- Publishing order events

**Database Tables** (MySQL):

- `orders` - order headers
- `order_items` - order line items (one-to-many)
- `payments` - payment records
- `order_status_history` - status audit trail

**Kafka Topics**:

- `order.created` - new order events
- `order.updated` - order status changes
- `order.cancelled` - cancellation events

---

### 5. Notification Service

**Port**: 8085  
**Protocol**: Kafka Consumer  
**Database**: MongoDB (primary), Redis (rate limiting)

**Responsibilities**:

- Consuming events from Kafka
- Sending notifications (email, SMS, push)
- Notification history tracking
- Rate limiting per user

**MongoDB Collections**:

- `notifications` - notification history
- `notification_preferences` - user preferences
- `notification_templates` - message templates

**Redis Keys**:

- `rate_limit:{user_id}:{type}` - rate limiting counters

---

### 6. Analytics Service

**Port**: 8086  
**Protocol**: Kafka Consumer, REST (Dashboard)  
**Database**: MongoDB (primary), Redis (cache)

**Responsibilities**:

- Event aggregation from Kafka
- Real-time metrics calculation
- Dashboard data generation
- Trend analysis

**MongoDB Collections**:

- `order_analytics` - aggregated order data
- `user_analytics` - user behavior metrics
- `product_analytics` - product performance

---

## Database Schema Design

### PostgreSQL - Authentication Service

```sql
-- users_auth table
CREATE TABLE users_auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_auth_email ON users_auth(email);

-- refresh_tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users_auth(id),
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users_auth(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
```

### PostgreSQL - User Service

```sql
-- users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    auth_id UUID UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_profiles table (One-to-One)
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    bio TEXT,
    avatar_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- user_addresses table (One-to-Many)
CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    address_type VARCHAR(20) CHECK (address_type IN ('shipping', 'billing', 'both')),
    street_address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_addresses_user_id ON user_addresses(user_id);
```

### PostgreSQL - Product Service

```sql
-- products table
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    parent_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- product_categories table (Many-to-Many)
CREATE TABLE product_categories (
    product_id UUID NOT NULL,
    category_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, category_id),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_categories_product_id ON product_categories(product_id);
CREATE INDEX idx_product_categories_category_id ON product_categories(category_id);

-- inventory table
CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID UNIQUE NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    reserved_quantity INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
```

### MySQL - Order Service

```sql
-- orders table
CREATE TABLE orders (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    status ENUM('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    shipping_address_id CHAR(36),
    billing_address_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_order_number (order_number),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB;

-- order_items table (One-to-Many)
CREATE TABLE order_items (
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    INDEX idx_order_id (order_id),
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB;

-- payments table
CREATE TABLE payments (
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status ENUM('pending', 'completed', 'failed', 'refunded') DEFAULT 'pending',
    transaction_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    INDEX idx_order_id (order_id),
    INDEX idx_transaction_id (transaction_id)
) ENGINE=InnoDB;

-- order_status_history table
CREATE TABLE order_status_history (
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    status VARCHAR(50) NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    INDEX idx_order_id (order_id)
) ENGINE=InnoDB;
```

### MongoDB - User Service

```javascript
// user_activities collection
{
  _id: ObjectId("..."),
  user_id: "uuid-string",
  activity_type: "login|logout|profile_update|password_change",
  metadata: {
    ip_address: "192.168.1.1",
    user_agent: "Mozilla/5.0...",
    location: "San Francisco, CA",
    // flexible fields based on activity type
  },
  timestamp: ISODate("2024-01-01T00:00:00Z"),
  created_at: ISODate("2024-01-01T00:00:00Z")
}

// user_preferences collection
{
  _id: ObjectId("..."),
  user_id: "uuid-string",
  preferences: {
    theme: "dark",
    language: "en",
    notifications: {
      email: true,
      sms: false,
      push: true
    },
    // other custom preferences
  },
  updated_at: ISODate("2024-01-01T00:00:00Z")
}
```

### MongoDB - Notification Service

```javascript
// notifications collection
{
  _id: ObjectId("..."),
  user_id: "uuid-string",
  type: "email|sms|push",
  category: "order|user|system",
  title: "Order Confirmation",
  message: "Your order #12345 has been confirmed",
  metadata: {
    order_id: "uuid-string",
    template_id: "order_confirmation_v1"
  },
  status: "pending|sent|failed",
  sent_at: ISODate("2024-01-01T00:00:00Z"),
  created_at: ISODate("2024-01-01T00:00:00Z")
}

// notification_templates collection
{
  _id: ObjectId("..."),
  template_id: "order_confirmation_v1",
  type: "email|sms|push",
  subject: "Order Confirmation - {{order_number}}",
  body: "Hi {{user_name}}, your order...",
  variables: ["user_name", "order_number", "total_amount"],
  active: true,
  created_at: ISODate("2024-01-01T00:00:00Z")
}
```

### MongoDB - Analytics Service

```javascript
// order_analytics collection
{
  _id: ObjectId("..."),
  date: ISODate("2024-01-01T00:00:00Z"),
  metrics: {
    total_orders: 150,
    total_revenue: 45000.00,
    average_order_value: 300.00,
    cancelled_orders: 5,
    completed_orders: 145
  },
  hourly_breakdown: [
    { hour: 0, orders: 5, revenue: 1500.00 },
    { hour: 1, orders: 3, revenue: 900.00 },
    // ... 24 hours
  ],
  updated_at: ISODate("2024-01-01T23:59:59Z")
}
```

---

## Detailed Flow Processes

### 1. User Registration Flow

```text
Client → User Service → Auth Service → PostgreSQL → Redis → User Service → MongoDB → Client
```

**Step-by-Step Process**:

1. **Client sends registration request**
   - **Endpoint**: `POST /api/v1/users/register`
   - **Protocol**: REST (HTTPS)
   - **Payload**:

     ```json
     {
       "email": "user@example.com",
       "password": "SecurePass123!",
       "first_name": "John",
       "last_name": "Doe"
     }
     ```

2. **User Service receives request**
   - **Service**: User Service (Port 8082)
   - **Action**: Validates input (email format, password strength)
   - **Validation checks**:
     - Email format validation
     - Password complexity (min 8 chars, uppercase, lowercase, number, special char)
     - Required fields present

3. **User Service calls Auth Service via gRPC**
   - **Protocol**: gRPC
   - **Method**: `authpb.CreateAuthUser`
   - **Request**:

     ```protobuf
     {
       email: "user@example.com",
       password: "SecurePass123!"
     }
     ```

4. **Auth Service processes authentication**
   - **Service**: Auth Service (Port 8081)
   - **Database**: PostgreSQL (auth_db)
   - **Action**:
     - Check if email already exists in `users_auth` table
     - Hash password using bcrypt (cost: 12)
     - Generate UUID for user
   - **Query**:

     ```sql
     SELECT id FROM users_auth WHERE email = 'user@example.com';
     ```

   - If exists: Return error "Email already registered"
   - If not exists: Continue

5. **Auth Service inserts user credentials**
   - **Database**: PostgreSQL (auth_db)
   - **Table**: `users_auth`
   - **Query**:

     ```sql
     INSERT INTO users_auth (id, email, password_hash, is_active, created_at, updated_at)
     VALUES ('550e8400-e29b-41d4-a716-446655440000', 'user@example.com',
             '$2a$12$hashed_password', true, NOW(), NOW())
     RETURNING id;
     ```

   - **Result**: Returns user auth_id

6. **Auth Service returns success to User Service**
   - **Protocol**: gRPC Response
   - **Response**:

     ```protobuf
     {
       auth_id: "550e8400-e29b-41d4-a716-446655440000",
       success: true
     }
     ```

7. **User Service creates user profile**
   - **Database**: PostgreSQL (user_db)
   - **Table**: `users`
   - **Transaction Start**
   - **Query 1** (Insert user):

     ```sql
     INSERT INTO users (id, auth_id, email, first_name, last_name, created_at, updated_at)
     VALUES (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440000',
             'user@example.com', 'John', 'Doe', NOW(), NOW())
     RETURNING id;
     ```

   - **Result**: user_id = '660e8400-e29b-41d4-a716-446655440000'

8. **User Service creates user profile record**
   - **Database**: PostgreSQL (user_db)
   - **Table**: `user_profiles`
   - **Query 2** (Insert profile - one-to-one relationship):

     ```sql
     INSERT INTO user_profiles (id, user_id, created_at, updated_at)
     VALUES (gen_random_uuid(), '660e8400-e29b-41d4-a716-446655440000',
             NOW(), NOW());
     ```

   - **Transaction Commit**

9. **User Service logs activity to MongoDB**
   - **Database**: MongoDB (user_db)
   - **Collection**: `user_activities`
   - **Operation**: Insert
   - **Document**:

     ```javascript
     {
       user_id: "660e8400-e29b-41d4-a716-446655440000",
       activity_type: "registration",
       metadata: {
         ip_address: "192.168.1.100",
         user_agent: "Mozilla/5.0...",
         registration_method: "email"
       },
       timestamp: new Date(),
       created_at: new Date()
     }
     ```

10. **User Service creates default preferences**
    - **Database**: MongoDB (user_db)
    - **Collection**: `user_preferences`
    - **Operation**: Insert
    - **Document**:

      ```javascript
      {
        user_id: "660e8400-e29b-41d4-a716-446655440000",
        preferences: {
          theme: "light",
          language: "en",
          notifications: {
            email: true,
            sms: false,
            push: true
          }
        },
        updated_at: new Date()
      }
      ```

11. **User Service returns success response**
    - **Response to Client**:

      ```json
      {
        "success": true,
        "message": "User registered successfully",
        "user": {
          "id": "660e8400-e29b-41d4-a716-446655440000",
          "email": "user@example.com",
          "first_name": "John",
          "last_name": "Doe"
        }
      }
      ```

**Error Handling**:

- If Auth Service fails (email exists): Rollback, return 409 Conflict
- If PostgreSQL transaction fails: Rollback all, call Auth Service to delete auth record
- If MongoDB fails: Log error but continue (non-critical data)

---

### 2. User Login Flow

```text
Client → Auth Service → PostgreSQL → Redis → Client
```

**Step-by-Step Process**:

1. **Client sends login request**
   - **Endpoint**: `POST /api/v1/auth/login`
   - **Protocol**: REST (HTTPS)
   - **Payload**:

     ```json
     {
       "email": "user@example.com",
       "password": "SecurePass123!"
     }
     ```

2. **Auth Service receives and validates request**
   - **Service**: Auth Service (Port 8081)
   - **Action**: Validate email format and password presence

3. **Auth Service queries user credentials**
   - **Database**: PostgreSQL (auth_db)
   - **Table**: `users_auth`
   - **Query**:

     ```sql
     SELECT id, email, password_hash, is_active
     FROM users_auth
     WHERE email = 'user@example.com';
     ```

   - **Result**:

     ```gRPC
     id: 550e8400-e29b-41d4-a716-446655440000
     password_hash: $2a$12$hashed_password
     is_active: true
     ```

4. **Auth Service validates password**
   - **Action**: Compare provided password with stored hash using bcrypt
   - **Function**: `bcrypt.CompareHashAndPassword(storedHash, providedPassword)`
   - If password invalid: Return 401 Unauthorized
   - If account inactive: Return 403 Forbidden
   - If password valid: Continue

5. **Auth Service generates JWT access token**
   - **Algorithm**: RS256 (RSA Signature with SHA-256)
   - **Claims**:

     ```json
     {
       "sub": "550e8400-e29b-41d4-a716-446655440000",
       "email": "user@example.com",
       "iat": 1704067200,
       "exp": 1704070800,
       "jti": "token-id-12345"
     }
     ```

   - **Expiry**: 1 hour (3600 seconds)

6. **Auth Service generates refresh token**
   - **Action**: Generate random secure token (32 bytes)
   - **Token**: `a7b8c9d0e1f2g3h4i5j6k7l8m9n0o1p2`
   - **Hash**: SHA-256 hash of token
   - **Expiry**: 7 days

7. **Auth Service stores refresh token in PostgreSQL**
   - **Database**: PostgreSQL (auth_db)
   - **Table**: `refresh_tokens`
   - **Query**:

     ```sql
     INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
     VALUES (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440000',
             'hashed_refresh_token', NOW() + INTERVAL '7 days', NOW())
     RETURNING id;
     ```

8. **Auth Service stores session in Redis**
   - **Database**: Redis
   - **Key**: `session:550e8400-e29b-41d4-a716-446655440000`
   - **Value**:

     ```json
     {
       "user_id": "550e8400-e29b-41d4-a716-446655440000",
       "email": "user@example.com",
       "token_id": "token-id-12345",
       "login_time": "2024-01-01T00:00:00Z",
       "ip_address": "192.168.1.100"
     }
     ```

   - **Command**: `SETEX session:550e8400-e29b-41d4-a716-446655440000 3600 "json_value"`
   - **TTL**: 3600 seconds (1 hour)

9. **Auth Service stores refresh token in Redis**
   - **Database**: Redis
   - **Key**: `refresh:a7b8c9d0e1f2g3h4i5j6k7l8m9n0o1p2`
   - **Value**: `550e8400-e29b-41d4-a716-446655440000`
   - **Command**: `SETEX refresh:a7b8c9d0e1f2g3h4i5j6k7l8m9n0o1p2 604800 "user_id"`
   - **TTL**: 604800 seconds (7 days)

10. **Auth Service calls User Service to log activity**
    - **Protocol**: gRPC (Async)
    - **Method**: `userpb.LogActivity`
    - **Request**:

      ```protobuf
      {
        user_id: "550e8400-e29b-41d4-a716-446655440000",
        activity_type: "login",
        metadata: {
          ip_address: "192.168.1.100",
          user_agent: "Mozilla/5.0..."
        }
      }
      ```

11. **User Service logs login activity**
    - **Database**: MongoDB (user_db)
    - **Collection**: `user_activities`
    - **Operation**: Insert
    - **Document**:

      ```javascript
      {
        user_id: "550e8400-e29b-41d4-a716-446655440000",
        activity_type: "login",
        metadata: {
          ip_address: "192.168.1.100",
          user_agent: "Mozilla/5.0...",
          device_type: "desktop",
          location: "San Francisco, CA"
        },
        timestamp: new Date(),
        created_at: new Date()
      }
      ```

12. **Auth Service returns tokens to client**
    - **Response**:

      ```json
      {
        "success": true,
        "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "a7b8c9d0e1f2g3h4i5j6k7l8m9n0o1p2",
        "token_type": "Bearer",
        "expires_in": 3600
      }
      ```

**Error Handling**:

- Invalid credentials: Return 401 with "Invalid email or password"
- Account inactive: Return 403 with "Account is disabled"
- Database connection error: Return 503 Service Unavailable
- Redis failure: Log error but continue (session can be validated from JWT)

---

### 3. Token Refresh Flow

```text
Client → Auth Service → Redis → PostgreSQL → Redis → Client
```

**Step-by-Step Process**:

1. **Client sends refresh token request**
   - **Endpoint**: `POST /api/v1/auth/refresh`
   - **Protocol**: REST (HTTPS)
   - **Payload**:

     ```json
     {
       "refresh_token": "a7b8c9d0e1f2g3h4i5j6k7l8m9n0o1p2"
     }
     ```

2. **Auth Service validates refresh token in Redis**
   - **Database**: Redis
   - **Key**: `refresh:a7b8c9d0e1f2g3h4i5j6k7l8m9n0o1p2`
   - **Command**: `GET refresh:a7b8c9d0e1f2g3h4i5j6k7l8m9n0o1p2`
   - **Result**: `550e8400-e29b-41d4-a716-446655440000` (user_id)
   - If not found or expired: Return 401 Unauthorized

3. **Auth Service verifies token in PostgreSQL**
   - **Database**: PostgreSQL (auth_db)
   - **Table**: `refresh_tokens`
   - **Query**:

     ```sql
     SELECT id, user_id, expires_at
     FROM refresh_tokens
     WHERE token_hash = 'hashed_token'
     AND user_id = '550e8400-e29b-41d4-a716-446655440000'
     AND expires_at > NOW();
     ```

   - If not found: Return 401 Unauthorized

4. **Auth Service retrieves user information**
   - **Database**: PostgreSQL (auth_db)
   - **Table**: `users_auth`
   - **Query**:

     ```sql
     SELECT id, email, is_active
     FROM users_auth
     WHERE id = '550e8400-e29b-41d4-a716-446655440000';
     ```

5. **Auth Service generates new JWT access token**
   - **Algorithm**: RS256
   - **Claims** (same structure as login)
   - **Expiry**: 1 hour

6. **Auth Service updates session in Redis**
   - **Database**: Redis
   - **Key**: `session:550e8400-e29b-41d4-a716-446655440000`
   - **Command**: `SETEX` with new session data
   - **TTL**: 3600 seconds

7. **Auth Service returns new access token**
   - **Response**:

     ```json
     {
       "success": true,
       "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...(new)",
       "token_type": "Bearer",
       "expires_in": 3600
     }
     ```

---

### 4. User Logout Flow

```text
Client → Auth Service → Redis → PostgreSQL → Client
```

**Step-by-Step Process**:

1. **Client sends logout request with JWT**
   - **Endpoint**: `POST /api/v1/auth/logout`
   - **Protocol**: REST (HTTPS)
   - **Headers**: `Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...`

2. **Auth Service validates and decodes JWT**
   - **Action**: Verify JWT signature and expiry
   - **Extract claims**: user_id, token_id (jti)

3. **Auth Service blacklists token in Redis**
   - **Database**: Redis
   - **Key**: `blacklist:token-id-12345`
   - **Value**: `revoked`
   - **Command**: `SETEX blacklist:token-id-12345 3600 "revoked"`
   - **TTL**: Remaining token lifetime (3600 seconds max)

4. **Auth Service deletes session from Redis**
   - **Database**: Redis
   - **Key**: `session:550e8400-e29b-41d4-a716-446655440000`
   - **Command**: `DEL session:550e8400-e29b-41d4-a716-446655440000`

5. **Auth Service deletes refresh tokens from PostgreSQL**
   - **Database**: PostgreSQL (auth_db)
   - **Table**: `refresh_tokens`
   - **Query**:

     ```sql
     DELETE FROM refresh_tokens
     WHERE user_id = '550e8400-e29b-41d4-a716-446655440000';
     ```

6. **Auth Service deletes refresh tokens from Redis**
   - **Database**: Redis
   - **Action**: Delete all refresh tokens for user
   - **Pattern**: `refresh:*`
   - **Note**: This requires tracking user's refresh tokens separately

7. **Auth Service returns success**
   - **Response**:

     ```json
     {
       "success": true,
       "message": "Logged out successfully"
     }
     ```

---

### 5. Create Order Flow (Complex Multi-Service)

```text
Client → Order Service → Auth Service → User Service → Product Service →
MySQL → Kafka → Notification Service → Analytics Service
```

**Step-by-Step Process**:

1. **Client sends create order request**
   - **Endpoint**: `POST /api/v1/orders`
   - **Protocol**: REST (HTTPS)
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Payload**:

     ```json
     {
       "items": [
         {
           "product_id": "770e8400-e29b-41d4-a716-446655440000",
           "quantity": 2
         },
         {
           "product_id": "880e8400-e29b-41d4-a716-446655440000",
           "quantity": 1
         }
       ],
       "shipping_address_id": "990e8400-e29b-41d4-a716-446655440000",
       "billing_address_id": "990e8400-e29b-41d4-a716-446655440000",
       "payment_method": "credit_card"
     }
     ```

2. **API Gateway validates JWT token**
   - **Action**: Extract JWT from Authorization header
   - **Validation**: Verify signature, check expiry, extract user_id
   - **Check Redis blacklist**:
     - **Database**: Redis
     - **Key**: `blacklist:token-id-12345`
     - **Command**: `EXISTS blacklist:token-id-12345`
     - If exists: Return 401 Unauthorized

3. **API Gateway calls Auth Service for additional validation**
   - **Protocol**: gRPC
   - **Method**: `authpb.ValidateToken`
   - **Request**: `{token: "jwt_string"}`
   - **Response**: `{valid: true, user_id: "550e8400-..."}`

4. **Request forwarded to Order Service**
   - **Service**: Order Service (Port 8084)
   - **Protocol**: gRPC
   - **User Context**: Attached from JWT (user_id)

5. **Order Service validates user via User Service**
   - **Protocol**: gRPC
   - **Method**: `userpb.GetUser`
   - **Request**: `{user_id: "550e8400-e29b-41d4-a716-446655440000"}`
   - **User Service queries PostgreSQL**:
     - **Database**: PostgreSQL (user_db)
     - **Table**: `users`
     - **Query**:

       ```sql
       SELECT id, email, first_name, last_name
       FROM users
       WHERE id = '550e8400-e29b-41d4-a716-446655440000'
       AND deleted_at IS NULL;
       ```

   - **Response**: User details or error if not found

6. **Order Service validates shipping address**
   - **Protocol**: gRPC
   - **Method**: `userpb.GetAddress`
   - **Request**: `{address_id: "990e8400-...", user_id: "550e8400-..."}`
   - **User Service queries PostgreSQL**:
     - **Database**: PostgreSQL (user_db)
     - **Table**: `user_addresses`
     - **Query**:

       ```sql
       SELECT id, user_id, street_address, city, state, postal_code, country
       FROM user_addresses
       WHERE id = '990e8400-e29b-41d4-a716-446655440000'
       AND user_id = '550e8400-e29b-41d4-a716-446655440000';
       ```

   - **Response**: Address details

7. **Order Service validates products and inventory (Loop for each item)**

   **For Product 1**:
   - **Protocol**: gRPC
   - **Method**: `productpb.GetProduct`
   - **Request**: `{product_id: "770e8400-e29b-41d4-a716-446655440000"}`

   - **Product Service checks Redis cache first**:
     - **Database**: Redis
     - **Key**: `product:770e8400-e29b-41d4-a716-446655440000`
     - **Command**: `GET product:770e8400-e29b-41d4-a716-446655440000`
     - **Cache Hit**: Return cached product
     - **Cache Miss**: Query PostgreSQL

   - **Product Service queries PostgreSQL** (if cache miss):
     - **Database**: PostgreSQL (product_db)
     - **Table**: `products`
     - **Query**:

       ```sql
       SELECT id, name, description, price, sku, is_active
       FROM products
       WHERE id = '770e8400-e29b-41d4-a716-446655440000'
       AND is_active = true;
       ```

     - **Result**:

       ```gRPC
       name: "Wireless Mouse"
       price: 29.99
       is_active: true
       ```

   - **Product Service caches result in Redis**:
     - **Database**: Redis
     - **Key**: `product:770e8400-e29b-41d4-a716-446655440000`
     - **Command**: `SETEX product:770e8400... 3600 "json_product_data"`
     - **TTL**: 3600 seconds

   **Check Inventory**:
   - **Protocol**: gRPC
   - **Method**: `productpb.CheckInventory`
   - **Request**: `{product_id: "770e8400-...", quantity: 2}`

   - **Product Service checks Redis inventory first**:
     - **Database**: Redis
     - **Key**: `inventory:770e8400-e29b-41d4-a716-446655440000`
     - **Command**: `GET inventory:770e8400-e29b-41d4-a716-446655440000`
     - **Result**: `{"available": 50, "reserved": 10}`

   - **Product Service queries PostgreSQL inventory**:
     - **Database**: PostgreSQL (product_db)
     - **Table**: `inventory`
     - **Query**:

       ```sql
       SELECT quantity, reserved_quantity
       FROM inventory
       WHERE product_id = '770e8400-e29b-41d4-a716-446655440000'
       FOR UPDATE; -- Lock row for update
       ```

     - **Result**: `quantity: 50, reserved_quantity: 10`
     - **Available**: 50 - 10 = 40 units
     - **Requested**: 2 units
     - **Check**: 40 >= 2 ✓ Pass

   **Repeat for Product 2** (same process)

8. **Order Service reserves inventory**
   - **Protocol**: gRPC
   - **Method**: `productpb.ReserveInventory`
   - **Transaction Start** in Product Service

   - **Product Service updates PostgreSQL**:
     - **Database**: PostgreSQL (product_db)
     - **Table**: `inventory`
     - **Query for Product 1**:

       ```sql
       UPDATE inventory
       SET reserved_quantity = reserved_quantity + 2,
           updated_at = NOW()
       WHERE product_id = '770e8400-e29b-41d4-a716-446655440000';
       ```

     - **Query for Product 2**:

       ```sql
       UPDATE inventory
       SET reserved_quantity = reserved_quantity + 1,
           updated_at = NOW()
       WHERE product_id = '880e8400-e29b-41d4-a716-446655440000';
       ```

   - **Product Service updates Redis inventory**:
     - **Database**: Redis
     - **Key**: `inventory:770e8400-...`
     - **Command**: `HINCRBY inventory:770e8400-... reserved 2`
     - **Repeat for Product 2**

   - **Transaction Commit**

9. **Order Service calculates order total**
   - **Calculation**:

     ```text
     Product 1: 2 × $29.99 = $59.98
     Product 2: 1 × $49.99 = $49.99
     Subtotal: $109.97
     Tax (8%): $8.80
     Shipping: $10.00
     Total: $128.77
     ```

10. **Order Service creates order in MySQL**
    - **Database**: MySQL (order_db)
    - **Transaction Start**

    - **Generate order number**:

      ```text
      ORD-20240101-001234
      ```

    - **Table**: `orders`
    - **Query**:

      ```sql
      INSERT INTO orders (
        id, user_id, order_number, status, total_amount,
        shipping_address_id, billing_address_id, created_at, updated_at
      ) VALUES (
        '111e8400-e29b-41d4-a716-446655440000',
        '550e8400-e29b-41d4-a716-446655440000',
        'ORD-20240101-001234',
        'pending',
        128.77,
        '990e8400-e29b-41d4-a716-446655440000',
        '990e8400-e29b-41d4-a716-446655440000',
        NOW(),
        NOW()
      );
      ```

11. **Order Service creates order items (one-to-many relationship)**
    - **Database**: MySQL (order_db)
    - **Table**: `order_items`
    - **Query for Item 1**:

      ```sql
      INSERT INTO order_items (
        id, order_id, product_id, product_name, quantity, unit_price, subtotal, created_at
      ) VALUES (
        '222e8400-e29b-41d4-a716-446655440000',
        '111e8400-e29b-41d4-a716-446655440000',
        '770e8400-e29b-41d4-a716-446655440000',
        'Wireless Mouse',
        2,
        29.99,
        59.98,
        NOW()
      );
      ```

    - **Query for Item 2**:

      ```sql
      INSERT INTO order_items (
        id, order_id, product_id, product_name, quantity, unit_price, subtotal, created_at
      ) VALUES (
        '333e8400-e29b-41d4-a716-446655440000',
        '111e8400-e29b-41d4-a716-446655440000',
        '880e8400-e29b-41d4-a716-446655440000',
        'Mechanical Keyboard',
        1,
        49.99,
        49.99,
        NOW()
      );
      ```

12. **Order Service creates payment record**
    - **Database**: MySQL (order_db)
    - **Table**: `payments`
    - **Query**:

      ```sql
      INSERT INTO payments (
        id, order_id, payment_method, amount, status, created_at, updated_at
      ) VALUES (
        '444e8400-e29b-41d4-a716-446655440000',
        '111e8400-e29b-41d4-a716-446655440000',
        'credit_card',
        128.77,
        'pending',
        NOW(),
        NOW()
      );
      ```

13. **Order Service creates status history**
    - **Database**: MySQL (order_db)
    - **Table**: `order_status_history`
    - **Query**:

      ```sql
      INSERT INTO order_status_history (
        id, order_id, status, note, created_at
      ) VALUES (
        '555e8400-e29b-41d4-a716-446655440000',
        '111e8400-e29b-41d4-a716-446655440000',
        'pending',
        'Order created',
        NOW()
      );
      ```

    - **Transaction Commit** (All MySQL operations)

14. **Order Service publishes event to Kafka**
    - **Message Broker**: Kafka
    - **Topic**: `order.created`
    - **Partition Key**: user_id (for ordering)
    - **Message**:

      ```json
      {
        "event_id": "evt-123456",
        "event_type": "order.created",
        "timestamp": "2024-01-01T12:00:00Z",
        "data": {
          "order_id": "111e8400-e29b-41d4-a716-446655440000",
          "order_number": "ORD-20240101-001234",
          "user_id": "550e8400-e29b-41d4-a716-446655440000",
          "user_email": "user@example.com",
          "total_amount": 128.77,
          "status": "pending",
          "items": [
            {
              "product_id": "770e8400-e29b-41d4-a716-446655440000",
              "product_name": "Wireless Mouse",
              "quantity": 2,
              "unit_price": 29.99
            },
            {
              "product_id": "880e8400-e29b-41d4-a716-446655440000",
              "product_name": "Mechanical Keyboard",
              "quantity": 1,
              "unit_price": 49.99
            }
          ]
        }
      }
      ```

15. **Notification Service consumes Kafka event**
    - **Message Broker**: Kafka
    - **Consumer Group**: `notification-service-group`
    - **Topic**: `order.created`
    - **Action**: Receives event message

16. **Notification Service checks user preferences**
    - **Database**: MongoDB (notification_db)
    - **Collection**: `notification_preferences`
    - **Query**:

      ```javascript
      db.notification_preferences.findOne({
        user_id: "550e8400-e29b-41d4-a716-446655440000",
      });
      ```

    - **Result**: User wants email notifications for orders

17. **Notification Service checks rate limit**
    - **Database**: Redis
    - **Key**: `rate_limit:550e8400-...:order_notifications`
    - **Command**: `INCR rate_limit:550e8400-...:order_notifications`
    - **Result**: Current count
    - **Check**: If count <= 10 per hour, proceed

18. **Notification Service retrieves template**
    - **Database**: MongoDB (notification_db)
    - **Collection**: `notification_templates`
    - **Query**:

      ```javascript
      db.notification_templates.findOne({
        template_id: "order_confirmation_v1",
        type: "email",
        active: true,
      });
      ```

19. **Notification Service sends email**
    - **Action**: Send email via SMTP/SendGrid
    - **To**: <user@example.com>
    - **Subject**: "Order Confirmation - ORD-20240101-001234"
    - **Body**: Rendered template with order details

20. **Notification Service stores notification history**
    - **Database**: MongoDB (notification_db)
    - **Collection**: `notifications`
    - **Operation**: Insert
    - **Document**:

      ```javascript
      {
        user_id: "550e8400-e29b-41d4-a716-446655440000",
        type: "email",
        category: "order",
        title: "Order Confirmation",
        message: "Your order #ORD-20240101-001234 has been confirmed",
        metadata: {
          order_id: "111e8400-e29b-41d4-a716-446655440000",
          template_id: "order_confirmation_v1",
          sent_to: "user@example.com"
        },
        status: "sent",
        sent_at: new Date(),
        created_at: new Date()
      }
      ```

21. **Analytics Service consumes Kafka event**
    - **Message Broker**: Kafka
    - **Consumer Group**: `analytics-service-group`
    - **Topic**: `order.created`
    - **Action**: Receives event message

22. **Analytics Service updates daily metrics**
    - **Database**: MongoDB (analytics_db)
    - **Collection**: `order_analytics`
    - **Operation**: Update (Upsert)
    - **Query**:

      ```javascript
      db.order_analytics.updateOne(
        { date: new Date("2024-01-01") },
        {
          $inc: {
            "metrics.total_orders": 1,
            "metrics.total_revenue": 128.77,
            "hourly_breakdown.12.orders": 1,
            "hourly_breakdown.12.revenue": 128.77,
          },
          $set: {
            updated_at: new Date(),
          },
        },
        { upsert: true },
      );
      ```

23. **Analytics Service updates product analytics**
    - **Database**: MongoDB (analytics_db)
    - **Collection**: `product_analytics`
    - **Operation**: Update (for each product)
    - **Query for Product 1**:

      ```javascript
      db.product_analytics.updateOne(
        {
          product_id: "770e8400-e29b-41d4-a716-446655440000",
          date: new Date("2024-01-01"),
        },
        {
          $inc: {
            sales_count: 2,
            revenue: 59.98,
          },
        },
        { upsert: true },
      );
      ```

24. **Analytics Service caches popular metrics in Redis**
    - **Database**: Redis
    - **Key**: `analytics:daily:2024-01-01`
    - **Command**:

      ```text
      HMSET analytics:daily:2024-01-01
        total_orders 150
        total_revenue 45000.00
      ```

    - **TTL**: 86400 seconds (24 hours)

25. **Order Service returns success response to client**
    - **Response**:

      ```json
      {
        "success": true,
        "message": "Order created successfully",
        "order": {
          "id": "111e8400-e29b-41d4-a716-446655440000",
          "order_number": "ORD-20240101-001234",
          "status": "pending",
          "total_amount": 128.77,
          "items_count": 2,
          "created_at": "2024-01-01T12:00:00Z"
        }
      }
      ```

**Error Handling & Rollback**:

- **User not found**: Return 404, no database changes
- **Invalid address**: Return 400, no database changes
- **Product not found**: Return 404, no database changes
- **Insufficient inventory**: Return 409, release any partial reservations
- **MySQL transaction fails**:
  - Rollback MySQL transaction
  - Call Product Service to release reserved inventory
  - Return 500 Internal Server Error
- **Kafka publish fails**:
  - Order still created (eventual consistency)
  - Retry Kafka publish with exponential backoff
  - Log error for manual intervention if retries fail

---

### 6. Get User Activity Log Flow

```text
Client → User Service → MongoDB → Client
```

**Step-by-Step Process**:

1. **Client sends request for activity logs**
   - **Endpoint**: `GET /api/v1/users/me/activities`
   - **Protocol**: REST (HTTPS)
   - **Headers**: `Authorization: Bearer <jwt_token>`
   - **Query Parameters**:

     ```text
     ?page=1&limit=20&activity_type=login&start_date=2024-01-01
     ```

2. **API Gateway validates JWT**
   - **Action**: Extract user_id from JWT
   - **User ID**: `550e8400-e29b-41d4-a716-446655440000`

3. **Request forwarded to User Service**
   - **Service**: User Service (Port 8082)
   - **Protocol**: REST

4. **User Service queries MongoDB**
   - **Database**: MongoDB (user_db)
   - **Collection**: `user_activities`
   - **Query**:

     ```javascript
     db.user_activities
       .find({
         user_id: "550e8400-e29b-41d4-a716-446655440000",
         activity_type: "login",
         timestamp: { $gte: new Date("2024-01-01T00:00:00Z") },
       })
       .sort({ timestamp: -1 })
       .skip(0) // (page - 1) * limit
       .limit(20);
     ```

5. **User Service retrieves activity count**
   - **Database**: MongoDB (user_db)
   - **Collection**: `user_activities`
   - **Query**:

     ```javascript
     db.user_activities.countDocuments({
       user_id: "550e8400-e29b-41d4-a716-446655440000",
       activity_type: "login",
       timestamp: { $gte: new Date("2024-01-01T00:00:00Z") },
     });
     ```

   - **Result**: Total count = 45

6. **User Service returns paginated response**
   - **Response**:

     ```json
     {
       "success": true,
       "data": [
         {
           "id": "60a7b8c9d0e1f2g3h4i5j6k7",
           "activity_type": "login",
           "metadata": {
             "ip_address": "192.168.1.100",
             "user_agent": "Mozilla/5.0...",
             "device_type": "desktop",
             "location": "San Francisco, CA"
           },
           "timestamp": "2024-01-15T08:30:00Z"
         }
         // ... 19 more items
       ],
       "pagination": {
         "page": 1,
         "limit": 20,
         "total": 45,
         "total_pages": 3
       }
     }
     ```

---

### 7. GraphQL Complex Product Query Flow

```text
Client → Product Service → PostgreSQL → Redis → Product Service → Client
```

**Step-by-Step Process**:

1. **Client sends GraphQL query**
   - **Endpoint**: `POST /graphql`
   - **Protocol**: HTTPS
   - **Query**:

     ```graphql
     query {
       product(id: "770e8400-e29b-41d4-a716-446655440000") {
         id
         name
         description
         price
         categories {
           id
           name
           slug
         }
         inventory {
           quantity
           available
         }
         relatedProducts {
           id
           name
           price
         }
       }
     }
     ```

2. **Product Service receives GraphQL request**
   - **Service**: Product Service (Port 8083)
   - **Action**: Parse and validate GraphQL query

3. **Product Service resolves product field**
   - **Check Redis cache**:
     - **Database**: Redis
     - **Key**: `product:770e8400-e29b-41d4-a716-446655440000`
     - **Command**: `GET product:770e8400-...`
     - **Cache Miss**: Continue to database

   - **Query PostgreSQL**:
     - **Database**: PostgreSQL (product_db)
     - **Table**: `products`
     - **Query**:

       ```sql
       SELECT id, name, description, price, sku, is_active
       FROM products
       WHERE id = '770e8400-e29b-41d4-a716-446655440000'
       AND is_active = true;
       ```

     - **Result**: Base product data

4. **Product Service resolves categories field (many-to-many)**
   - **Database**: PostgreSQL (product_db)
   - **Query** (JOIN junction table):

     ```sql
     SELECT c.id, c.name, c.slug
     FROM categories c
     INNER JOIN product_categories pc ON c.id = pc.category_id
     WHERE pc.product_id = '770e8400-e29b-41d4-a716-446655440000';
     ```

   - **Result**:

     ```json
     [
       { "id": "cat-001", "name": "Electronics", "slug": "electronics" },
       {
         "id": "cat-002",
         "name": "Computer Accessories",
         "slug": "computer-accessories"
       }
     ]
     ```

5. **Product Service resolves inventory field**
   - **Check Redis cache**:
     - **Database**: Redis
     - **Key**: `inventory:770e8400-...`
     - **Command**: `HGETALL inventory:770e8400-...`
     - **Result**: `{ quantity: "50", reserved: "12" }`

   - **If cache miss, query PostgreSQL**:
     - **Database**: PostgreSQL (product_db)
     - **Table**: `inventory`
     - **Query**:

       ```sql
       SELECT quantity, reserved_quantity
       FROM inventory
       WHERE product_id = '770e8400-e29b-41d4-a716-446655440000';
       ```

   - **Calculate available**: 50 - 12 = 38

6. **Product Service resolves relatedProducts field**
   - **Logic**: Products in same categories
   - **Database**: PostgreSQL (product_db)
   - **Query**:

     ```sql
     SELECT DISTINCT p.id, p.name, p.price
     FROM products p
     INNER JOIN product_categories pc ON p.id = pc.product_id
     WHERE pc.category_id IN (
       SELECT category_id
       FROM product_categories
       WHERE product_id = '770e8400-e29b-41d4-a716-446655440000'
     )
     AND p.id != '770e8400-e29b-41d4-a716-446655440000'
     AND p.is_active = true
     LIMIT 5;
     ```

7. **Product Service caches complete result in Redis**
   - **Database**: Redis
   - **Key**: `product:770e8400-...`
   - **Value**: Complete product JSON with all resolved fields
   - **Command**: `SETEX product:770e8400-... 3600 "json_data"`
   - **TTL**: 3600 seconds

8. **Product Service returns GraphQL response**
   - **Response**:

     ```json
     {
       "data": {
         "product": {
           "id": "770e8400-e29b-41d4-a716-446655440000",
           "name": "Wireless Mouse",
           "description": "Ergonomic wireless mouse with 6 buttons",
           "price": 29.99,
           "categories": [
             {
               "id": "cat-001",
               "name": "Electronics",
               "slug": "electronics"
             },
             {
               "id": "cat-002",
               "name": "Computer Accessories",
               "slug": "computer-accessories"
             }
           ],
           "inventory": {
             "quantity": 50,
             "available": 38
           },
           "relatedProducts": [
             {
               "id": "880e8400-...",
               "name": "Mechanical Keyboard",
               "price": 49.99
             }
             // ... more products
           ]
         }
       }
     }
     ```

---

## API Specifications

### REST API Endpoints

#### Authentication Service

```list
POST   /api/v1/auth/register        - Register new user
POST   /api/v1/auth/login           - User login
POST   /api/v1/auth/refresh         - Refresh access token
POST   /api/v1/auth/logout          - User logout
POST   /api/v1/auth/forgot-password - Request password reset
POST   /api/v1/auth/reset-password  - Reset password
```

#### User Service

```list
GET    /api/v1/users/me             - Get current user profile
PUT    /api/v1/users/me             - Update user profile
GET    /api/v1/users/me/activities  - Get activity logs
GET    /api/v1/users/me/addresses   - List user addresses
POST   /api/v1/users/me/addresses   - Add new address
PUT    /api/v1/users/me/addresses/:id - Update address
DELETE /api/v1/users/me/addresses/:id - Delete address
```

#### Product Service

```list
GET    /api/v1/products             - List products (with filters)
GET    /api/v1/products/:id         - Get product details
GET    /api/v1/categories           - List categories
GET    /api/v1/categories/:id       - Get category details
GET    /api/v1/categories/:id/products - Get products in category
```

#### Order Service (via API Gateway)

```list
POST   /api/v1/orders               - Create new order
GET    /api/v1/orders               - List user orders
GET    /api/v1/orders/:id           - Get order details
PUT    /api/v1/orders/:id/cancel    - Cancel order
```

### GraphQL Schema

```graphql
type User {
  id: ID!
  email: String!
  firstName: String
  lastName: String
  profile: UserProfile
  addresses: [Address!]!
  orders: [Order!]!
}

type UserProfile {
  id: ID!
  phone: String
  dateOfBirth: String
  bio: String
  avatarUrl: String
}

type Address {
  id: ID!
  addressType: AddressType!
  streetAddress: String!
  city: String!
  state: String
  postalCode: String!
  country: String!
  isDefault: Boolean!
}

type Product {
  id: ID!
  name: String!
  description: String
  price: Float!
  sku: String!
  categories: [Category!]!
  inventory: Inventory!
  relatedProducts: [Product!]!
}

type Category {
  id: ID!
  name: String!
  slug: String!
  products: [Product!]!
  parent: Category
  children: [Category!]!
}

type Inventory {
  quantity: Int!
  reserved: Int!
  available: Int!
}

type Order {
  id: ID!
  orderNumber: String!
  status: OrderStatus!
  totalAmount: Float!
  items: [OrderItem!]!
  shippingAddress: Address
  billingAddress: Address
  payment: Payment
  createdAt: String!
}

type OrderItem {
  id: ID!
  product: Product!
  quantity: Int!
  unitPrice: Float!
  subtotal: Float!
}

type Payment {
  id: ID!
  paymentMethod: String!
  amount: Float!
  status: PaymentStatus!
  transactionId: String
}

enum AddressType {
  SHIPPING
  BILLING
  BOTH
}

enum OrderStatus {
  PENDING
  CONFIRMED
  PROCESSING
  SHIPPED
  DELIVERED
  CANCELLED
}

enum PaymentStatus {
  PENDING
  COMPLETED
  FAILED
  REFUNDED
}

type Query {
  me: User!
  product(id: ID!): Product
  products(
    page: Int
    limit: Int
    category: String
    minPrice: Float
    maxPrice: Float
  ): ProductConnection!
  category(id: ID!): Category
  order(id: ID!): Order
  myOrders(page: Int, limit: Int): OrderConnection!
}

type Mutation {
  updateProfile(input: UpdateProfileInput!): User!
  addAddress(input: AddressInput!): Address!
  createOrder(input: CreateOrderInput!): Order!
  cancelOrder(orderId: ID!): Order!
}
```

---

## Event Schemas

### Kafka Event: order.created

```json
{
  "event_id": "string (uuid)",
  "event_type": "order.created",
  "version": "1.0",
  "timestamp": "ISO 8601 datetime",
  "source": "order-service",
  "data": {
    "order_id": "string (uuid)",
    "order_number": "string",
    "user_id": "string (uuid)",
    "user_email": "string",
    "total_amount": "decimal",
    "status": "pending",
    "items": [
      {
        "product_id": "string (uuid)",
        "product_name": "string",
        "quantity": "integer",
        "unit_price": "decimal",
        "subtotal": "decimal"
      }
    ],
    "shipping_address": {
      "street": "string",
      "city": "string",
      "state": "string",
      "postal_code": "string",
      "country": "string"
    }
  }
}
```

### Kafka Event: order.updated

```json
{
  "event_id": "string (uuid)",
  "event_type": "order.updated",
  "version": "1.0",
  "timestamp": "ISO 8601 datetime",
  "source": "order-service",
  "data": {
    "order_id": "string (uuid)",
    "order_number": "string",
    "previous_status": "string",
    "new_status": "string",
    "updated_by": "string (user_id or system)",
    "note": "string (optional)"
  }
}
```

### Kafka Event: order.cancelled

```json
{
  "event_id": "string (uuid)",
  "event_type": "order.cancelled",
  "version": "1.0",
  "timestamp": "ISO 8601 datetime",
  "source": "order-service",
  "data": {
    "order_id": "string (uuid)",
    "order_number": "string",
    "user_id": "string (uuid)",
    "cancellation_reason": "string",
    "refund_amount": "decimal",
    "items_to_restock": [
      {
        "product_id": "string (uuid)",
        "quantity": "integer"
      }
    ]
  }
}
```

---

## Security Implementation

### JWT Token Structure

**Access Token**:

- Algorithm: RS256
- Expiry: 1 hour
- Claims: sub (user_id), email, iat, exp, jti (token_id)

**Refresh Token**:

- Type: Opaque token (random 32 bytes)
- Storage: PostgreSQL + Redis
- Expiry: 7 days

### Password Hashing

- Algorithm: bcrypt
- Cost factor: 12
- Salt: Automatically generated per password

### API Security

- All endpoints require HTTPS
- JWT validation on protected routes
- Token blacklisting for logout
- Rate limiting per user/IP
- CORS configuration
- Request validation and sanitization

---

## Performance Optimizations

### Caching Strategy

**Redis Cache TTLs**:

- User sessions: 1 hour
- Product data: 1 hour
- Inventory: 5 minutes
- Category listings: 24 hours
- Analytics dashboard: 1 hour

### Database Indexing

**PostgreSQL Indexes**:

- All foreign keys
- Email fields (unique)
- Composite indexes for common queries
- Partial indexes for filtered queries

**MySQL Indexes**:

- Order user_id, order_number, status
- Order items order_id
- Created_at timestamps for time-based queries

### Connection Pooling

- PostgreSQL: Max 100 connections per service
- MySQL: Max 100 connections
- Redis: Max 50 connections
- MongoDB: Max 100 connections

---

## License

This documentation is provided as-is for reference and implementation purposes.

## Contributing

Contributions are welcome! Please follow the standard pull request process.

## Support

For questions or support, please open an issue in the repository.
