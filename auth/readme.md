Environment Variables
VALKEY_URL
POSTGRES_URL
AUTH_PORT ":10221"

OAuth Flow
frontend calls /oauth/login
backend redirects to github
user auths via github
github redirects to /oauth/callback?code=CODE
backend use code to get access token
backend use access token to get user id
backend use user id to link to user in db
backend issues session
