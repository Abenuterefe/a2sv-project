# Blog API (Postman-ready)

Base URL: http://localhost:8080/api/v1

Note on casing: some responses use PascalCase field names (e.g., `Title`, `UserID`) because the base `Blog` struct lacks JSON tags. Others (like Popular and Search) include additional fields with `snake_case` (e.g., `author_name`). Examples below reflect actual responses.

Postman tip:
- Set an environment variable: `baseUrl = http://localhost:8080/api/v1`
- For protected routes, add header: `Authorization: Bearer {{token}}`

---

## 1) Create Blog

- Method: POST
- URL: {{baseUrl}}/blogs
- Auth: Required (Bearer token)

Headers:
- Content-Type: application/json
- Authorization: Bearer {{token}}

Body (raw JSON):
{
  "title": "The Road map to Google",
  "content": "This is a Tech Blog",
  "tags": ["golang", "programming"]
}

Success 201:
```
{
  "ID": "6895aec7f39172726e146c27",
  "UserID": "68935e8b56ce1bbf14b7a95f",
  "Title": "The Road map to Google",
  "Content": "This is a Tech Blog",
  "Tags": ["golang", "programming"],
  "CreatedAt": "2025-08-08T08:01:11.009Z",
  "UpdatedAt": "2025-08-08T08:01:11.009Z",
  "ViewCount": 0,
  "LikeCount": 0,
  "DislikeCount": 0
}
```
Errors:
- 400 Invalid request payload
- 401 User not authenticated
- 500 Server error

---

## 2) Get My Blogs (or a specific user's blogs)

- Method: GET
- URL: {{baseUrl}}/blogs
- Auth: Required

Query Params:
- user_id (optional): hex string of another user
- page (optional): default 1
- limit (optional): default 5, maximum 5

Success 200:
```
[
  {
    "ID": "6893680169083d319f036463",
    "UserID": "68935e8b56ce1bbf14b7a95f",
    "Title": "CARs",
    "Content": "They are too high",
    "Tags": ["golang","programming","clean-architecture","web-development"],
    "CreatedAt": "2025-08-06T14:34:41.669Z",
    "UpdatedAt": "2025-08-06T15:46:46.209Z",
    "ViewCount": 0,
    "LikeCount": 0,
    "DislikeCount": 0
  }
]
```

Errors:
- 401 User not authenticated
- 500 Server error

---

## 3) Get Blog by ID

- Method: GET
- URL: {{baseUrl}}/blogs/:id
- Auth: Public

Success 200:
```
{
  "ID": "6895aec7f39172726e146c27",
  "UserID": "68935e8b56ce1bbf14b7a95f",
  "Title": "The Road map to Google",
  "Content": "This is a Tech Blog",
  "Tags": null,
  "CreatedAt": "2025-08-08T08:01:11.009Z",
  "UpdatedAt": "2025-08-08T08:01:11.009Z",
  "ViewCount": 1,
  "LikeCount": 0,
  "DislikeCount": 0
}
```

Errors:
- 404 Blog not found

---

## 4) Update Blog

- Method: PUT
- URL: {{baseUrl}}/blogs/:id
- Auth: Required + Owner only

Headers:
- Content-Type: application/json
- Authorization: Bearer {{token}}

Body (raw JSON, any updatable fields):
```
{
  "title": "Updated title",
  "content": "Updated content",
  "tags": ["golang","web"]
}
```
Success 200: same shape as "Get Blog by ID"

Errors:
- 400 Invalid blog ID | Invalid request payload
- 404 Blog not found
- 500 Server error

---

## 5) Delete Blog

- Method: DELETE
- URL: {{baseUrl}}/blogs/:id
- Auth: Required + Owner only

Success:
- 204 No Content

Errors:
- 500 Server error

---

## 6) Popular Blogs

- Method: GET
- URL: {{baseUrl}}/blogs/popular
- Auth: Public

Query Params:
- limit (optional): default 10

Success 200:
```
{
  "message": "Popular blogs retrieved successfully",
  "data": [
    {
      "id": "6895aec7f39172726e146c27",
      "title": "The Road map to Google",
      "content": "This is a Tech Blog",
      "user_id": "68935e8b56ce1bbf14b7a95f",
      "like_count": 12,
      "dislike_count": 1,
      "view_count": 250,
      "comment_count": 3,
      "popularity_score": 57.5,
      "created_at": "2025-08-08T08:01:11.009Z",
      "updated_at": "2025-08-08T08:01:11.009Z"
    }
  ],
  "count": 1
}
```
Errors:
- 500 Server error

---

## 7) Filter Blogs

- Method: GET
- URL: {{baseUrl}}/blogs/filter
- Auth: Public

Query Params:
- tags (repeatable): tags=tech&tags=golang
- date_from (YYYY-MM-DD)
- date_to (YYYY-MM-DD)
- popularity_sort: views | likes | dislikes | engagement
- sort_order: asc | desc
- limit (default 20)
- skip (default 0)
- page (alternative to skip)

Success 200:
```
{
  "message": "Blogs filtered successfully",
  "data": {
    "blogs": [
      {
        "ID": "6893680169083d319f036463",
        "UserID": "68935e8b56ce1bbf14b7a95f",
        "Title": "CARs",
        "Content": "They are too high",
        "Tags": ["golang","programming"],
        "CreatedAt": "2025-08-06T14:34:41.669Z",
        "UpdatedAt": "2025-08-06T15:46:46.209Z",
        "ViewCount": 0,
        "LikeCount": 0,
        "DislikeCount": 0
      }
    ],
    "count": 1,
    "total_count": 1,
    "page": 1,
    "limit": 20
  }
}
```
Validation errors (400):
- Invalid date_from/date_to format (use YYYY-MM-DD)
- Invalid popularity_sort or sort_order
- Invalid limit/skip

---

## 8) Search Blogs (Title and/or Author)

- Method: GET
- URL: {{baseUrl}}/blogs/search
- Auth: Public

At least one of `title` or `author` is required.

Query Params:
- title (optional): case-insensitive partial match
- author (optional): username, case-insensitive (resolved via user lookup)
- limit (default 20)
- skip (default 0)
- page (alternative to skip)

Success 200:
```
{
  "message": "Blog search completed successfully",
  "data": {
    "blogs": [
      {
        "ID": "6893680169083d319f036463",
        "UserID": "68935e8b56ce1bbf14b7a95f",
        "Title": "CARs",
        "Content": "They are too high",
        "Tags": ["golang","programming","clean-architecture","web-development"],
        "CreatedAt": "2025-08-06T14:34:41.669Z",
        "UpdatedAt": "2025-08-06T15:46:46.209Z",
        "ViewCount": 0,
        "LikeCount": 0,
        "DislikeCount": 0,
        "author_name": "Mo"
      }
    ],
    "count": 1,
    "total_count": 1,
    "query": {
      "title": "CARs",
      "author": "mo",
      "limit": 20
    }
  }
}
```
Errors:
- 400 At least one search parameter (title or author) must be provided

---

## 9) Blog Interactions: Like a Blog

- Method: POST
- URL: {{baseUrl}}/blogs/:id/like
- Auth: Required (Bearer token)

Path Params:
- id: blog id (hex string)

Success 200:
```
{
  "message": "Blog liked successfully"
}
```
Errors:
- 401 User not authenticated
- 500 Server error

---

## 10) Blog Interactions: Dislike a Blog

- Method: POST
- URL: {{baseUrl}}/blogs/:id/dislike
- Auth: Required (Bearer token)

Path Params:
- id: blog id (hex string)

Success 200:
```
{
  "message": "Blog disliked successfully"
}
```
Errors:
- 401 User not authenticated
- 500 Server error

---

## 11) Blog Interactions: Record a View

- Method: POST
- URL: {{baseUrl}}/blogs/:id/view
- Auth: Public (works for anonymous and authenticated users)

Path Params:
- id: blog id (hex string)

Behavior:
- Prevents rapid duplicate views per user/ip/agent in a short window.
- For anonymous users, dedup is based on IP + User-Agent.

Success 200:
```
{
  "message": "Blog view recorded"
}
```
Errors:
- 500 Server error

---

## 12) Comments: List Comments for a Blog

- Method: GET
- URL: {{baseUrl}}/blogs/:id/comments
- Auth: Public

Path Params:
- id: blog id (hex string)

Success 200:
```
[
  {
    "ID": "<commentId>",
    "UserID": "<userId>",
    "BlogID": "<blogId>",
    "Content": "Nice post!",
    "CreatedAt": "ISO",
    "UpdatedAt": "ISO"
  }
]
```
Errors:
- 400 Blog ID is required
- 500 Server error

---

## 13) Comments: Get Comment by ID

- Method: GET
- URL: {{baseUrl}}/comments/:id
- Auth: Public

Success 200:
```
{
  "ID": "<commentId>",
  "UserID": "<userId>",
  "BlogID": "<blogId>",
  "Content": "Nice post!",
  "CreatedAt": "ISO",
  "UpdatedAt": "ISO"
}
```
Errors:
- 404 Comment not found

---

## 14) Comments: Create Comment on a Blog

- Method: POST
- URL: {{baseUrl}}/blogs/:id/comments
- Auth: Required (Bearer token)

Headers:
- Content-Type: application/json
- Authorization: Bearer {{token}}

Path Params:
- id: blog id (hex string)

Body (raw JSON):
```
{
  "content": "Nice post!"
}
```
Success 201:
```
{
  "ID": "<commentId>",
  "UserID": "<userId>",
  "BlogID": "<blogId>",
  "Content": "Nice post!",
  "CreatedAt": "ISO",
  "UpdatedAt": "ISO"
}
```
Errors:
- 400 Invalid request payload | Blog ID is required
- 401 User not authenticated
- 500 Server error

---

## 15) Comments: Update Comment

- Method: PUT
- URL: {{baseUrl}}/comments/:id
- Auth: Required (Bearer token) and owner only

Headers:
- Content-Type: application/json
- Authorization: Bearer {{token}}

Body (raw JSON):
```
{
  "content": "Updated comment"
}
```
Success 200: same shape as Get Comment by ID

Errors:
- 400 Invalid request payload | Invalid comment ID
- 401 User not authenticated
- 403 You can only modify your own comments
- 404 Comment not found
- 500 Server error

---

## 16) Comments: Delete Comment

- Method: DELETE
- URL: {{baseUrl}}/comments/:id
- Auth: Required (Bearer token) and owner only

Success:
- 204 No Content

Errors:
- 401 User not authenticated
- 403 You can only delete your own comments
- 404 Comment not found
- 500 Server error

---

## Quick Postman Examples

- Create Blog
  - Method: POST
  - URL: {{baseUrl}}/blogs
  - Headers: Authorization: Bearer {{token}}, Content-Type: application/json
  - Body: { "title": "My Post", "content": "...", "tags": ["go","web"] }

- Popular Blogs
  - GET {{baseUrl}}/blogs/popular?limit=5

- Filter Blogs (tags + likes desc)
  - GET {{baseUrl}}/blogs/filter?tags=go&tags=web&popularity_sort=likes&sort_order=desc&limit=10

- Search Blogs (title + author)
  - GET {{baseUrl}}/blogs/search?title=golang&author=mo&limit=5
