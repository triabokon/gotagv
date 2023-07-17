# Gotagv

## Overview

Gotagv is a simple video and annotations management service written in Go.

The service has Restful API to manage videos and annotations with basic API security using JWT Token. 

The project is structured into three main modules:

1. **Auth**: creating and validating JWT token to authorize user into system.
2. **Server**: serves HTTP requests and converts data from request body and query params to Go structs.
3. **Controller**: implementation of main business logic like data validation, conversion etc.
4. **Storage**: implements storage methods and queries to  database.

## Getting Started

The project uses Go Modules for dependency management.
All the dependencies for the application are specified in the `go.mod` file.
Dependencies that are needed for linting are specified in the `tools/go.mod`.

Also, service `Docker` image and can be launched with docker-compose.

### Prerequisites

Make sure you have installed `Docker` and `docker-compose` on your system. 
Refer to [official Docker installation guide](https://docs.docker.com/engine/install/) and
[official docker-compose installation guide](https://docs.docker.com/compose/install/).

### Installation

Clone the repository:

```bash
git clone https://github.com/triabokon/gotagv.git
```

Then, navigate into the cloned directory:

```bash
cd gotagv
```

## Building the Project

You can build the project using the following steps:

Build and run the project:
```bash
docker-compose up --build
```
This will launch two containers - `postgres` database and HTTP web server.

## Migrating database

You can migrate database using commands:

1. to apply migrations:
```bash
make migrate-up
```
2. to rollback migrations:
```bash
make migrate-down
```

## Usage example

1. Launch service:
```bash
docker-compose up
```

2. Create user to login into system:

```bash
curl -X POST 'localhost:8080/signup'
```

Example response:
```
{
  "user_id": "6ce179e9-53a6-430f-9833-3de929d9696b",
  "token": <jwt_token>
}
```
3. Create video

```bash
curl -X POST 'localhost:8080/videos/add' --header 'Authorization: Bearer <jwt_token>' -d '{"url": "https://youtube.com/test", "duration": "2m37s"}'
```

Example response:
```
{
  "video_id": "0bb49819-a5be-437e-8fc2-d4f3cebef283"
}
```

4. Get all videos

```bash
curl -X POST 'localhost:8080/videos' --header 'Authorization: Bearer <jwt_token>'
```

Example response:
```
{
  "videos": [
    {
      "id": "0bb49819-a5be-437e-8fc2-d4f3cebef283",
      "user_id": "6ce179e9-53a6-430f-9833-3de929d9696b",
      "url": "https://youtube.com/test",
      "duration": 157000000000,
      "created_at": "2023-07-17T06:59:17.463301Z",
      "updated_at": "2023-07-17T06:59:17.463324Z"
    }
  ]
}
```

5. Create annotation

```bash
curl -X POST 'localhost:8080/annotations/add' --header 'Authorization: Bearer <jwt_token>' -d '{
    "video_id": "0bb49819-a5be-437e-8fc2-d4f3cebef283",
    "start_time": "2m",
    "end_time": "2m25s",
    "type": "title",
    "title": "First annotation!"
}'
```

Example response:
```
{
  "annotation_id": "fdf2d1ef-9f91-4adf-9723-75f3e777e56b"
}
```

6. Get annotations for specific video

```bash
curl -X POST 'localhost:8080/annotations' --header 'Authorization: Bearer <jwt_token>' -d '{"video_id": "0bb49819-a5be-437e-8fc2-d4f3cebef283"}'
```

Example response:
```
{
  "annotations": [
    {
      "id": "fdf2d1ef-9f91-4adf-9723-75f3e777e56b",
      "video_id": "0bb49819-a5be-437e-8fc2-d4f3cebef283",
      "user_id": "6ce179e9-53a6-430f-9833-3de929d9696b",
      "start_time": 120000000000,
      "end_time": 145000000000,
      "type": "title",
      "title": "First annotation!",
      "created_at": "2023-07-17T07:01:55.120516Z",
      "updated_at": "2023-07-17T07:01:55.120541Z"
    }
  ]
}
```

7. Update annotation

```bash
curl -X POST 'localhost:8080/annotations/update/fdf2d1ef-9f91-4adf-9723-75f3e777e56b' --header 'Authorization: Bearer <jwt_token>' -d '{
    "start_time": "1m10s",
    "end_time": "1m25s",
    "type": "commentary",
    "message": "Hello!"
}'
```
8. Delete annotation
```bash
curl -X POST 'localhost:8080/annotations/delete/fdf2d1ef-9f91-4adf-9723-75f3e777e56b' --header 'Authorization: Bearer <jwt_token>'
```
9. Delete video
```bash
curl -X POST 'localhost:8080/videos/delete/0bb49819-a5be-437e-8fc2-d4f3cebef283' --header 'Authorization: Bearer <jwt_token>'
```

## Linting

This project uses `golangci-lint` for linting, it's configuration is specified in `.golangci.yml`.

To lint the code, you need to run command:
```bash
make lint
```

## Makefile

Makefile is used for simplifying some utils commands of the project.

Here is output of `make help` command:

```
Usage: make <TARGETS> ... <OPTIONS>

Available targets are:

    help               Show this help
    clean              Remove binaries
    download-deps      Download and install dependencies
    tidy               Perform go tidy steps
    lint               Run all linters
    build              Compile packages and dependencies
    migrate-up         Applies migrations on database
    migrate-down       Rollbacks migrations on database
```

## Assumptions

While developing this task some assumptions were made:

- video files already uploaded to some resource, so this service stores only metadata about the video,
- any authenticated user can perform all operations, no role-based access control.

## Further improvements

Some other things could be done to improve the project:

1. Unit and integration tests: it would be great to write unit tests and integration tests.
2. Paging and Sorting: for APIs returning multiple items (e.g., listing all annotations), introduce paging to limit the response size and sorting to customize the order of the results.
3. Caching: caching could be used, to improve performance, especially for read-heavy APIs.
4. Video Upload and Processing: while the current API assumes videos are stored elsewhere, a future feature could allow users to upload videos directly, possibly with additional video processing functionalities (e.g., video transcoding, thumbnail generation).
5. Role-Based Access Control: now, it's assumed that any authenticated user can perform all operations, but in the future, roles and permissions could be added so that certain operations can be restricted (e.g., only video owner can delete video).