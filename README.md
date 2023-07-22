# NexusNet

Rest Api

## Description

Social media rest api with user registration/authorisation.
Authenticated user can create posts, delete own posts, get specific posts, and the seim manipulations with stories.
Signed users can comment posts and stories.
Each user can see all posts, user can sort them or search with filters.
Also, users can message each other.

## Technologies 

 - Language: Go
 - Frontend: not present
 - Database: PostrgeSql
 - Other: JSON, SMTP, HTTP

## API Reference

```http
Get, "/healthcheck"

Post, "/signup"
Post, "/users/activated"
Post, "/login"
Get, "/users/:id"

Get, "/"
Post, "/posts/create"
Patch, "/posts/:id"
Delete, "/posts/:id"
Get, "/posts/:id"

Get, "/stories/:id"
Get, "/users/:id/stories"
Delete, "/stories/:id"
Post, "/stories"

Post, "/post/:id/comments/create"
Patch, "/post/:post_id/comments/:id"
Get, "/post/:id/comments"

Get, "/direct/:id"
Get, "/direct"
Post, "/direct/:id"
Post, "/directs/:id"
```
