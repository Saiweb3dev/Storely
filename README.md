# Storely

This project, **Storely**, is a file-sharing application inspired by Google Drive. It allows users to upload and download files, manage their storage, and includes features like user authentication and Prometheus metrics for monitoring. The project is built with a **Next.js** frontend and a **Go** backend for a seamless user experience.

---

## Features

- **Chunked File Upload**: Support for uploading large files in chunks.
- **File Storage in MinIO**: Files are stored securely in a MinIO bucket.
- **User Authentication**: Register and login endpoints for managing user access.
- **Storage Health**: Monitor user's storage usage and health.
- **Prometheus Metrics**: Metrics endpoint for application monitoring.
- **Modern Frontend**: Built with Next.js, TypeScript, and TailwindCSS for a responsive and dynamic user interface.

---

## Project Structure

- **Frontend**: Located in the `frontend` folder.
  - Built with Next.js, TypeScript, and TailwindCSS.
  - Provides a user-friendly interface for file management, user registration, and login.
  - Run using `npm run dev`.

- **Backend**: Located in the `backend` folder.
  - Powered by Go with Gorilla Mux for routing.
  - Handles file storage, authentication, and metrics.
  - Run using `go run main.go`.

---

## Endpoints

### File Upload and Management

- **`POST /upload-chunk`**
  - Upload a file chunk.

- **`POST /api/minio/files/init`**
  - Initialize file upload in MinIO.

- **`GET /files/minio/{fileId}`**
  - Retrieve a file from MinIO using its ID.

- **`DELETE /api/minio/files/delete`**
  - Delete a file from MinIO.

- **`POST /api/minio/files/{fileId}/complete`**
  - Mark a file upload as complete in MinIO.

### User Management

- **`POST /api/auth/register`**
  - Register a new user.

- **`POST /api/auth/login`**
  - Login an existing user.

### Storage Monitoring

- **`GET /get/user/storageHealth`**
  - Get the storage usage and health details of the logged-in user.

### Metrics

- **`GET /metrics`**
  - Access Prometheus metrics for monitoring application performance.

### Testing Endpoints

- **`POST /test/post`**
  - Insert test data into the database.

- **`GET /test/last`**
  - Retrieve the last test data from the database.

---

## Technologies Used

- **Frontend**: Next.js, TypeScript, TailwindCSS
- **Backend Framework**: Gorilla Mux
- **Database**: MongoDB
- **Object Storage**: MinIO
- **Monitoring**: Prometheus
- **Programming Language**: Go

---

## Setup Instructions

1. Clone the repository.
2. Navigate to the `frontend` folder and install dependencies with `npm install`.
3. Start the frontend server with `npm run dev`.
4. Navigate to the `backend` folder and set up your environment variables for MongoDB, MinIO, and Prometheus.
5. Start the backend server using `go run main.go`.
6. Use tools like Postman or a browser to interact with the API endpoints via the frontend interface.

---

## Future Enhancements

- Implement user-specific file quotas.
- Add payment gateways for expanding user storage limits.
- Enhance the frontend with additional file management features.

---

## Contributing

Feel free to contribute to this project by submitting issues or pull requests. For major changes, please discuss them via issue first.

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.

