# Chat Application

## File Upload System

### Secure File Uploads with Presigned URLs
Our chat application enhances security and efficiency by utilizing presigned URLs for file uploads. This method allows users to securely upload files directly to our cloud storage without routing them through our chat server, thereby reducing latency and server load.

#### How It Works
When a user selects a file for upload, the chat server generates a presigned URL that grants temporary access to the storage bucket. The file is then uploaded directly to the storage bucket(S3) using this URL, ensuring that the data transfer is secure and direct.

![Media Meta Data (3)](https://github.com/user-attachments/assets/daeadf78-3862-434e-8f57-ddf597438e2f)


#### Future Implementations
**Content Delivery Network (CDN) Integration:** To implement a CDN to further enhance the delivery speeds of static files like images and videos across global locations. This will ensure users experience faster loading times and improved performance as shown in design.

## Heartbeat Service

### Maintaining User Connections with Heartbeat Messages
The Heartbeat Service in our chat application plays a crucial role in managing connection states. This service ensures that the application can handle user presence and disconnections gracefully.

#### How It Works
- **Ping/Pong Mechanism**: The client (frontend) sends a 'ping' message at regular intervals to the chat server via a WebSocket connection. In response, the server sends back a 'pong' message, confirming the connection's activeness.
- **Publish to Pub/Sub**: Upon receiving a ping, the chat server publishes a heartbeat message to a dedicated Pub/Sub channel named "HEARTBEAT".
- **Subscriber Actions**: Services subscribed to the HEARTBEAT channel listen for these messages and update the Time-To-Live (TTL) of each user's session in Redis. Maintaining a connection pool with redis.
- **Connection Monitoring**: The chat server is also subscribed to a Redis channel that notifies about TTL expirations (`__keyevent@0__:expired`). If a user's TTL expires (indicating inactivity), the chat server receives this expiration event and terminates the inactive connection, ensuring system resources are efficiently managed.

  ![Chat Server (3)](https://github.com/user-attachments/assets/75570a91-31bb-47c2-a462-7b95ba0d5389)
