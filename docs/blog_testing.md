# Blog testing summary

This document summarizes what was tested for the blog features, why these tests matter, and how they were written and run.

## scope covered

- Blog use case
  - CreateBlog: sets ID/userID/timestamps and calls repo
  - UpdateBlog: updates UpdatedAt and calls repo
  - FilterBlogs: date range validation, popularity_sort and sort_order validation, paging math and response
  - SearchBlogs: parameter validation (title/author required), defaulting and non-negative limit/skip
  - GetPopularBlogs: computes score, sorts by score, applies limit

- Blog HTTP handlers (Gin)
  - Create/Get/Update/Delete blog
  - Popular blogs: limit parsing (default/custom)
  - Filter: query parsing including page→skip conversion and date formats
  - Search: query parsing, limit/skip/page

- Blog interactions
  - Like: toggle off when already liked; switch from dislike to like
  - Dislike: new dislike increments counters
  - View: anonymous debounce (no duplicate recent view) and new view increments counter

- Comments (light coverage to support blog flows)
  - Use case: CreateComment invalid blogID error; success path
  - Handlers: create/list/update/delete with auth/ownership checks

## testing strategy and logic

Layer isolation using mocks
- Use case tests mock repositories (BlogRepositoryInterface, CommentRepositoryInterface, BlogInteractionRepositoryInterface) to keep tests pure unit and fast. This isolates business logic from DB and HTTP layers.
- Handler tests mock the corresponding use cases to test only HTTP parsing/validation/response wiring.

What we assert and why
- Validation errors: Ensure early validation guards work and return clear messages (e.g., date_from > date_to, missing search params, negative limit/skip). This prevents bad queries hitting the DB.
- Parameter parsing: Verify handlers parse query/body correctly, including page→skip logic and limit caps; this protects public API behavior.
- Happy paths: Ensure correct calls are made with expected arguments and responses are shaped and status-coded properly (201, 200, 204, 400/401/403/404).
- Popularity logic: Confirm score ordering and limit application so the endpoint is predictable and stable.
- Interaction toggles: Exercise like/dislike toggle/switch behavior and view debounce to ensure counters remain consistent.

Implementation details
- testify and mockery generated mocks under `mocks/` are used to set expectations and stub returns.
- gin handler tests use `httptest` and minimal routers; `userID` is injected via a simple middleware when needed to simulate authenticated requests.
- Mongo ObjectID parsing is triggered in some use cases (views/comments). Tests use a valid 24-hex string (`507f1f77bcf86cd799439011`) to avoid conversion errors when needed, or assert error paths for invalid IDs.
- Matchers are kept specific where it adds value (e.g., checking page→skip) and relaxed where coupling to defaults would be brittle.

## files added/updated (tests)

- `delivery/controllers/blog_handler_test.go` — blog handler tests (create/get/update/delete, popular, filter, search)
- `delivery/controllers/blog_interaction_handler_test.go` — like/dislike/view handler tests
- `delivery/controllers/comment_handler_test.go` — comment handler auth/ownership tests
- `usecase/blog_usecase_test.go` — blog use case validations, popularity, create/update
- `usecase/blog_interaction_usecase_test.go` — like/dislike/view logic
- `usecase/comment_usecase_test.go` — comment create invalid/valid blog ID

## how to run

Windows PowerShell examples:

```
go test ./... -count=1 -v
```

Run specific package:

```
go test ./delivery/controllers -v
go test ./usecase -v
```

## results snapshot

- Current status: All blog-related tests PASS across controllers and use cases.
- Fast feedback: Tests are unit-level with mocks; no DB or network required.

## next steps (optional)

- Add a few negative-path handler tests for 500 propagation when use cases/repositories return errors.
- Add middleware-only tests (auth/role/ownership) for complete HTTP layer coverage.
- Consider a small integration test suite for the repository layer using a test MongoDB if/when needed.
