**Task Test for Back-End Developer (Golang)**

**Task: Real-Time Chat Application with Microservices Architecture, Asynchronous Task Processing, WebSocket, and Real-Time Notifications**

**Objective:** Develop a real-time chat application using Golang with a microservices architecture. Implement RESTful API endpoints, AsyncQ for asynchronous task processing, WebSocket for real- time communication, Redis for caching and a notification system for real-time alerts. This will test the candidate’s proficiency in Golang, RESTful API development, AsyncQ, WebSocket, Redis, and Test- Driven Development (TDD).

**Requirements:**

1. **Set up Golang Project:**

- Initialize a new Golang project.
- Set up MongoDB as the primary database.

2. **Models:**

- Create a User model with the following fields:
  - ID (ObjectID)
  - Username (String, unique)
  - Email (String, unique)
  - PasswordHash (String)
  - CreatedAt (Timestamp, default to current time)
- Create a Message model with the following fields:
  - ID (ObjectID)
  - Content (String)
  - SenderID ( ObjectID of User)
  - RoomID ( ObjectID of Room)
  - Timestamp (Timestamp, default to current time)
- Create a Room model with the following fields:
  - ID (ObjectID)
  - Name (String, unique)
  - CreatedAt (Timestamp, default to current time)

3. **REST API:**

- Implement the following API endpoints:
  - POST /register/ - Register a new user.
  - POST /login/ - Authenticate a user and return a JWT token.
  - GET /messages/ - List all messages.
  - POST /messages/ - Create a new message.
  - GET /rooms/ - List all chat rooms.
  - POST /rooms/ - Create a new chat room.

4. **Asynchronous Task Processing with AsyncQ:**

- Set up AsyncQ for task scheduling.
- Create a task to periodically archive old messages (e.g., messages older than 30 days) into an archive collection in MongoDB.
- Create a task to send daily summaries of chat activities to users.
- Schedule the tasks to run daily.

5. **WebSocket for Real-Time Communication:**

- Implement WebSocket for real-time communication.
- Create an endpoint /ws to handle WebSocket connections.
- Implement chat rooms where users can join and communicate in real-time.
- Clients connected via WebSocket should receive new messages in real-time.
- Notify all connected clients in a room when a message is created through the WebSocket.

6. **Real-Time Notifications:**

- Implement a notification system to send real-time alerts to users.
- Create an endpoint /notifications to manage user notification preferences (e.g., enable/disable notifications).
- Send real-time notifications to users via WebSocket when they receive a new message.
- Ensure that notifications can be customized based on user preferences (e.g., mute notifications for specific rooms).

7. **Caching with Redis:**

- Set up Redis for caching frequently accessed data.
- Cache the list of messages and chat rooms to reduce database load.
- Implement cache invalidation strategies when data changes.

8. **Testing:**

- Write unit tests for the REST API endpoints.
- Write tests for the AsyncQ tasks to ensure they are scheduled and executed correctly.
- Implement WebSocket testing to ensure real-time communication is working as expected.
- Write tests to verify Redis caching and invalidation logic.
- Write tests to verify the real-time notification system.

9. **Docker:**

- Create a Dockerfile and docker-compose.yml file to containerize the application.
- The docker-compose.yml should include services for the Golang app, MongoDB,

  Redis, and AsyncQ worker.

10\.**Git Commits:**

- Ensure meaningful and descriptive commit messages.
- Commit frequently to demonstrate incremental progress and document your development process.

**Submission:**

- Provide a link to the GitHub repository with the project code.
- Ensure the repository includes a README.md file with instructions on how to set up and run

  the project using Docker.

- Include screenshots or logs of the test results, including AsyncQ task execution, WebSocket communication search results, Redis caching, and real-time notifications.
