# Chat Application

## Heartbeat Service

### Maintaining User Connections with Heartbeat Messages
The Heartbeat Service in our chat application plays a crucial role in monitoring user activity and managing connection states. This service ensures that the application can handle user presence and disconnections gracefully.

#### How It Works
- **Ping/Pong Mechanism**: The client (frontend) sends a 'ping' message at regular intervals to the chat server via a WebSocket connection. In response, the server sends back a 'pong' message, confirming the connection's activeness.
- **Publish to Pub/Sub**: Upon receiving a ping, the chat server publishes a heartbeat message to a dedicated Pub/Sub channel named "HEARTBEAT".
- **Subscriber Actions**: Services subscribed to the HEARTBEAT channel listen for these messages and update the Time-To-Live (TTL) of each user's session in Redis. Maintaining a connection pool with redis.
- **Connection Monitoring**: The chat server is also subscribed to a Redis channel that notifies about TTL expirations (`__keyevent@0__:expired`). If a user's TTL expires (indicating inactivity), the chat server receives this expiration event and terminates the inactive connection, ensuring system resources are efficiently managed.

![Heartbeat Service Diagram]![![ChatServer1-ezgif com-resize](https://github.com/user-attachments/assets/61221a58-c05a-452e-81d8-1a6de0ff68e1)

