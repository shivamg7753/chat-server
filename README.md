# ChatFlow - Real-time Chat Application

A modern, real-time chat application built with Go and WebSockets featuring a beautiful dark-themed UI.

## Features

- 🔐 **User Authentication** - JWT-based login and registration
- 💬 **Real-time Messaging** - WebSocket-powered instant messaging
- 🏠 **Multiple Rooms** - General, Random, and Tech Talk chat rooms
- 👥 **Live User Counts** - See active users in each room
- 🎨 **Modern UI** - Beautiful dark theme with gradient accents
- 💾 **Message Persistence** - All messages saved to SQLite database
- 🔄 **Auto-reconnect** - Automatic WebSocket reconnection on disconnect

## Tech Stack

### Backend
- **Go** - Server-side language
- **Fiber** - Fast HTTP web framework
- **SQLite** - Lightweight database
- **WebSocket** - Real-time communication
- **JWT** - Secure authentication
- **bcrypt** - Password hashing

### Frontend
- **HTML/CSS/JavaScript** - Vanilla web technologies
- **WebSocket API** - Real-time messaging
- **Inter Font** - Modern typography

## Prerequisites

- Go 1.16 or higher
- A modern web browser

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd chat-server
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment** (optional)
   
   Copy `.env.example` to `.env` and modify if needed:
   ```bash
   cp .env.example .env
   ```

   Default configuration:
   - Port: `3000`
   - Database: `./chat.db`
   - JWT Secret: Change in production!

## Running the Application

1. **Start the server**
   ```bash
   go run cmd/main.go
   ```

2. **Open your browser**
   
   Navigate to: `http://localhost:3000`

3. **Create an account**
   - Click "Sign up"
   - Enter username and password
   - Start chatting!

## Usage

### Creating an Account
1. Click "Sign up" on the login screen
2. Enter a unique username
3. Create a password
4. You'll be automatically logged in

### Sending Messages
1. Type your message in the input field
2. Press Enter or click the send button
3. Messages appear instantly for all users in the room

### Switching Rooms
1. Click on any room in the sidebar (General, Random, Tech Talk)
2. Messages are saved per room
3. User counts update in real-time

### Logging Out
- Click the "Logout" button in the sidebar
- Your session will be cleared

## Project Structure

```
chat-server/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── auth/                # JWT authentication
│   ├── config/              # Configuration management
│   ├── database/            # SQLite database operations
│   ├── handlers/            # HTTP request handlers
│   ├── models/              # Data models
│   ├── static/              # Frontend assets (CSS, JS)
│   ├── views/               # HTML templates
│   └── websocket/           # WebSocket server
├── .env.example             # Example environment variables
├── go.mod                   # Go module dependencies
└── README.md                # This file
```

## API Endpoints

### Authentication
- `POST /api/register` - Create new user account
- `POST /api/login` - Login with credentials

### Messages
- `GET /api/messages?room={room}` - Get messages for a room
- `GET /api/room-stats` - Get active user counts per room

### WebSocket
- `WS /ws?token={jwt}&room={room}` - WebSocket connection

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3000` |
| `DB_PATH` | SQLite database path | `./chat.db` |
| `JWT_SECRET` | Secret key for JWT | Change in production! |
| `JWT_EXPIRY_HOURS` | JWT token expiry time | `24` |

## Features in Detail

### Authentication
- Secure password hashing with bcrypt
- JWT tokens for session management
- Token-based WebSocket authentication

### Real-time Messaging
- WebSocket connections for instant messaging
- Room-based message broadcasting
- Message persistence in SQLite

### User Interface
- Dark theme with purple-pink gradients
- Smooth animations and transitions
- Responsive design
- Real-time user count updates

## Development

### Running in Development
```bash
go run cmd/main.go
```

### Building for Production
```bash
go build -o chat-server cmd/main.go
./chat-server
```

## Security Notes

⚠️ **Important for Production:**
- Change `JWT_SECRET` in `.env` to a strong random string
- Use HTTPS in production
- Configure proper CORS settings
- Use a production-grade database (PostgreSQL, MySQL)
- Implement rate limiting
- Add input validation and sanitization

## Troubleshooting

### Port already in use
If port 3000 is busy, change `PORT` in `.env`:
```
PORT=8080
```

### Database locked
If you see database locked errors, ensure only one instance is running.

### WebSocket connection fails
- Check that the server is running
- Verify you're logged in
- Check browser console for errors

## License

MIT License - feel free to use this project for learning and development.

## Author

Built with ❤️ using Go and modern web technologies.
