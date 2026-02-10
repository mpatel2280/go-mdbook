# go-mdbook

Full-stack mdBook portal with Go + MongoDB backend and React frontend. Includes role-based auth (admin, reader) and Dockerized services.

## Quickstart

1. `docker compose up --build`
2. Open `http://localhost:3000`
3. Login with the default admin:
   - Email: `admin@example.com`
   - Password: `admin123`
4. Create users from the Admin panel (self-registration is disabled).
5. Create a book with slug `sample` and build it to view the bundled sample content.

The backend expects mdBook-compatible source folders under `backend/books/<slug>`.
A sample book exists at `backend/books/sample`.

## Services

- Backend API: `http://localhost:8080`
- Frontend: `http://localhost:3000`
- MongoDB: `mongodb://localhost:27017`

## Notes

- Set `JWT_SECRET`, `ADMIN_EMAIL`, `ADMIN_PASSWORD` in `docker-compose.yml` for production.
- mdBook is installed in the backend container. Use the admin endpoint to build.
