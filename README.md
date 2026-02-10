# chirpy
### chirpy is an app to practice using Go's http.server methods.

## What does it do? 
### Chirpy emulates the backend of a social media app like "X", including some effort to create authentication and authorization.

NB: to run this on a local machine, you would need to have GO installed, run `go build`
Then, you would need to start the server with `go run .` from the app's home directory: /chirpy

## Key Endpoints
### Once the server is running, you can use curl to query the endpoints:
	- "POST /admin/reset" (resets all tables)
	- "GET /api/healthz" (confirms that the app is running with "OK")
	- "POST /api/users" (returns all users)
	- "POST /api/chirps" (returns all chirps, but one can add `author_id=` to search by author and `sort={asc or desc} to sort)
	- "GET /api/chirps" (returns all chirps)
	- "GET /api/chirps/{chirpID}" (returns chirps by chirpID)
	- "POST /api/login" (logs in user)
	- "PUT /api/users" (lists all users)
