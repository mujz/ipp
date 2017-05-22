# Demo

Go to [ipp.mujz.ca](https://ipp.mujz.ca). Sign up with email or Facebook. Break it if you can ;)

# Run it yourself

All you need is docker and docker-compose

```Bash
docker-compose up
```

# Planning

The main three features of this service is authentication, number generation and incrementation, and exposing a RESTful API. Since this application is very simple and small, we should aim to use a minimalistic language/framework with as little overhead as possible. For example, it's probably not worth it to use ruby on rails or Express.js for such a simple service. I will choose to do this in Go because it's excellent for writing lightweight API services and it's very fast.

## Architecture

We need a way to persist the users' information and numbers. The main decision to make is whether to choose a relational or a non-relational data store. Well, let's look at the data that we need to persist. We have user's email, password, and number. That's it! We only need one table for this. Let's call it `users` table. Therefore, it doesn't really matter what database engine we use. I'll go with PostgresSQL since it's open-source and I've worked a lot with it before. Here's what our schema looks like:

```SQL
CREATE TABLE users (
  id          SERIAL PRIMARY KEY,
  email       VARCHAR(254) UNIQUE,
  password    VARCHAR(128),
  facebook_id VARCHAR(128) UNIQUE,
  num         INTEGER NOT NULL DEFAULT 1,
  CONSTRAINT email_or_facebook_id
    CHECK(facebook_id IS NOT NULL OR email IS NOT NULL),
  CONSTRAINT email_and_password CHECK(
    (email IS NOT NULL AND password IS NOT NULL) OR
    (email IS NULL AND password IS NULL)
  )
)
```

We can use `github.com/lib/pq` Postgres driver for Go and the standard `database/sql` for queries.

We also should cache users' numbers so that we don't have to keep fetching them on every call. We can use Redis, but since the requirements aren't very clear on this, we can just simply create an LRU cache in memory so that:

- We look for the number in the cache every time the current number is requested. If the number is not in cache, we'd fetch it from the DB and add it to the cache.
- Update the cache every time a number is updated. If the number is not in cache, we'd add it.

I didn't get enough time to do that, but it should be simple to implement.

## Exposing RESTful API

All API conform to the [jsonapi](http://jsonapi.org) specifications.

The builtin `net/http` package should be sufficient for our needs. We will have 4 main endpoints:

- `GET /next`
  - Increments the user's number in the database and cache
    - return 500 on overflow (There's probably a better error code for this though)
  - Returns the new user's number
- `GET /current`
  - Looks for the user's number in the cache,
    - hit: returns it from cache
    - miss: fetches it from the user's table in the DB and adds it to the cache
  - Returns the user's number
- `PUT /current`
  - Verifies body (for example not integer, overflow, underflow, null, etc.)
    - returns 400 if verification fails
  - Updates the user's number in the DB and the cache
  - Returns the updated number.
- `POST /login`
  - Verifies the email and password
    - returns 401 if it fails
  - Returns the user's API key
- `POST /signup`
  - Validates the email and password
    - If email already used or invalid or password violates the rules, return 400 with the reason in the response body
  - Creates a user with the passed email and password
  - Initializes the user's number at 1
  - Returns the API key

If any of these endpoints is called without an API key or with an invalid API key, a 401 is returned.

For any other endpoint, a 404 is returned.

Database errors or any other internal errors return 500.

## Authentication

We can use JSON Web Tokens (JWT) for email and password auth and Oauth for third party login services (Google, Facebook, etc.). We can use [`jwt-go`](github.com/dgrijalva/jwt-go) for the former and the official [`oauth2`](https://github.com/golang/oauth2) for the latter.

## Incrementing

It's pretty straight forward. The only thing that I want to say about this is that we should have the database do the incrementation and not the Go service. This is to make sure that the increment is done safely and we don't get concurrency issues or race conditions. For example, if the current number is 11:

```SQL
-- Don't do this
UPDATE user SET number = 12 WHERE email='mujtaba@example.com';

-- Instead do this
UPDATE user SET number = number + 1 WHERE email='mujtaba@example.com`;
```

## Frontend

Again, a minimalistic SPA framework would be appropriate here. `riot.js`, being very small in size and having all the features we need, makes it a great choice. The rest can be done with vanilla javascript.

# Known Issues

- We're using a 4 byte integer. This means that if our counter exceeds 2147483647 we get an overflow. This is easy to fix. Just use a larger integer. However, since the requirements don't specify this, we should be fine with this current limitation.
- We're not enforcing the use of strong passwords. This is also easy to implement.
- Although I explained how to implement the cache, I didn't get to do it. It would be nice to have that in place.

# Specification

## User Story

As a developer I need a way to get integers that automatically increment so that I can generate identifiers from within my code. My code is javascript running as an AWS Lambda function and I don’t want the extra complexity of creating and managing a database or key/value store. I also don’t want the complexity of dealing with other AWS services such as S3. I’m in a hurry and would rather just call a REST endpoint that returns the next available integer so that I can get on with my job. Additionally, my code needs to be stateless, meaning that I can’t store any data between calls. Why I need to generate identifiers using sequential integers is not important ;) Suffice it to say that this challenge is based on a real-world scenario.

## Tasks

Develop a rest service that:

1. Allows me to register as a user. At a minimum, this should be a REST endpoint that accepts an email address and a password and returns an API key.
1. Returns the next integer in my sequence when called. For example, if my current integer is 12, the service should return 13 when it is called. The endpoint should be secured by API key. I should not have to provide the previous value of the integer for this to work. Fetching the next integer should cause my current integer to increment by 1 on the server so that if I call the endpoint again, I get the next integer in the sequence.
1. Allows me to fetch my current integer. For example, if my current integer is 12, the service should return 12. The endpoint should be secured by API key.
1. Allows me to reset my integer to an arbitrary, non-negative value. For example, my integer may be currently 1005. I would like to reset it to 1000. The endpoint should be secured by API key.
1. Allow sign up using OAuth
  - Github, Facebook, Google, anything that supports it!
1. Build a UI for the service, especially the account creation functionality, as a single page app that consumes your API.

Deploy your API somewhere and include the link in your README so we can try it out without having to run it.

## Examples

```bash
$ # Get the next integer in the sequence
$ curl https://myapi.com/v1/next -H “Authorization: Bearer XXXXXX”
$ # Get the current integer
$ curl https://myapi.com/v1/current -H “Authorization: Bearer XXXXX”
$ # Reset the current integer
$ curl -X “PUT” https://myapi.com/v1/current -H “Authorization: Bearer XXXXX” --data “current=1000”
```
