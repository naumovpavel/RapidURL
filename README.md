# RapidURL

RapidURL is a robust REST API URL shortener crafted in Go. This versatile tool not only shortens URLs but also incorporates advanced features to enhance performance, security, and monitoring.

## Features:

- **Clean Architecture:** RapidURL is built on a clean and modular architecture, promoting maintainability, scalability, and testability.

- **JWT Authentication:** Secure your shortened URLs with JWT authentication, ensuring that only authorized users can create and manage links.

- **Caching with Memcached:** Employ a caching mechanism using Memcached to optimize the search of link, reducing latency and enhancing overall system performance.

- **Metrics collection with Prometheus:** RapidURL integrates Prometheus for collecting detailed metrics, providing insights into system behavior, performance, and usage patterns.

- **Visualization with Grafana:** Visualize the collected metrics using Grafana, allowing to gain meaningful insights through interactive dashboards.


## API Usage Examples

### 1. User Registration

**Description:** This endpoint allows users to register for a RapidURL account.

**Endpoint:** `/user/register`

**Method:** `POST`


**Request Body:**
```json
{
    "name": "bob",
    "email": "bob@gmail.com",
    "password": "pass"
}
```

**Response:**
200 OK
```json
{
  "message": "successfully registered"
}
```

### 2. User Authorization

**Description:** This endpoint allows users to login for a RapidURL account. It also set a cookie with JWT token

**Endpoint:** `/user/login`

**Method:** `POST`


**Request Body:**
```json
{
    "name": "bob",
    "password": "pass"
}
```

**Response:**
200 OK
```json
{
  "jwt": "eyJhbGciRgtyUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU5NDY1MjksInVzZXJpZCI6MX0.QCrA2bJ8ekMJyuKZFwPeWbqu8DxekKrPMvjIOc51gCU"
}
```

### 3. Creating alias for a link

**Description:** This endpoint allows users to create alias for link. If there is no alias in request, it will be a random. You need to be authorized to use this endpoint

**Endpoint:** `/link/add`

**Method:** `POST`


**Request Body:**
```json
{
    "alias" : "gh",
    "url" : "https://github.com/"
}
```

**Response:**
200 OK
```json
{
  "alias": "gh"
}
```

### 4. Redirect

**Description:** This endpoint allows users to be redirected using an alias

**Endpoint:** `/gh`

**Method:** `GET`

**Response:**
307 Temporary redirect

Location: github.com
