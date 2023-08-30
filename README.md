# RapidURL

RapidURL is a link shortener service implemented in the Go programming language. The service provides an API with authentication capabilities using JWT tokens. Users can utilize RapidURL to create concise, easy-to-remember links and share them with others.

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
  "status": "Ok"
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
  "status": "Ok",
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
  "status": "Ok",
  "alias": "gh"
}
```