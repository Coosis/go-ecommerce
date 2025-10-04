1. User visits frontend application.
2. Frontend application sends a request to this 
backend api gateway.
3. API gateway finds that there's no user session.
4. API gateway returns a HTTP Unauthorized response to the frontend application.
5. Frontend application sees this, redirects the user to the login page.
6. User hits login via github on the frontend.
7. Frontend application redirects the user to GitHub for authentication.
8. User logs in to GitHub and authorizes the application.
9. GitHub redirects the user back to the frontend application with an authorization code.
10. Frontend application sends the authorization code to the API gateway.
11. API gateway sends the authorization code to the GitHub OAuth service to exchange it for an access token.
12. API gateway uses grpc to create a user session from auth service.
13. API gateway returns the user session to the frontend application.
