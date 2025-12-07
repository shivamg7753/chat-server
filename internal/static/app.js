// Global state
let ws = null;
let currentUser = null;
let currentRoom = 'general';
let authToken = null;
let roomStatsInterval = null;

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
  checkAuthStatus();
});

// ============================================
// AUTH FUNCTIONS
// ============================================

function checkAuthStatus() {
  const token = sessionStorage.getItem('authToken');
  const username = sessionStorage.getItem('username');
  const userId = sessionStorage.getItem('userId');

  if (token && username && userId) {
    authToken = token;
    currentUser = { username, userId: parseInt(userId) };
    showChat();
    connectWebSocket();
  } else {
    showAuth();
  }
}

function showAuth() {
  document.getElementById('auth-container').classList.remove('hidden');
  document.getElementById('chat-container').classList.add('hidden');
}

function showChat() {
  document.getElementById('auth-container').classList.add('hidden');
  document.getElementById('chat-container').classList.remove('hidden');

  if (currentUser) {
    document.getElementById('current-username').textContent = currentUser.username;
    document.getElementById('user-avatar').textContent = currentUser.username.charAt(0).toUpperCase();
  }

  loadMessages();
  updateRoomCounts();

  // Clear existing interval to prevent duplicates
  if (roomStatsInterval) {
    clearInterval(roomStatsInterval);
  }

  // Update room counts every 5 seconds
  roomStatsInterval = setInterval(updateRoomCounts, 5000);
}

function switchToRegister(e) {
  e.preventDefault();
  document.getElementById('login-form').classList.remove('active');
  document.getElementById('register-form').classList.add('active');
  clearErrors();
}

function switchToLogin(e) {
  e.preventDefault();
  document.getElementById('register-form').classList.remove('active');
  document.getElementById('login-form').classList.add('active');
  clearErrors();
}

function clearErrors() {
  document.getElementById('login-error').textContent = '';
  document.getElementById('register-error').textContent = '';
}

async function handleRegister(e) {
  e.preventDefault();

  const username = document.getElementById('register-username').value.trim();
  const password = document.getElementById('register-password').value;
  const errorEl = document.getElementById('register-error');

  errorEl.textContent = '';

  try {
    const response = await fetch('/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      errorEl.textContent = data.error || 'Registration failed';
      return;
    }

    // Save auth data
    sessionStorage.setItem('authToken', data.token);
    sessionStorage.setItem('username', data.username);
    sessionStorage.setItem('userId', data.user_id);

    authToken = data.token;
    currentUser = { username: data.username, userId: data.user_id };

    // Show chat
    showChat();
    connectWebSocket();

  } catch (error) {
    console.error('Registration error:', error);
    errorEl.textContent = 'Network error. Please try again.';
  }
}

async function handleLogin(e) {
  e.preventDefault();

  const username = document.getElementById('login-username').value.trim();
  const password = document.getElementById('login-password').value;
  const errorEl = document.getElementById('login-error');

  errorEl.textContent = '';

  try {
    const response = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      errorEl.textContent = data.error || 'Login failed';
      return;
    }

    // Save auth data
    sessionStorage.setItem('authToken', data.token);
    sessionStorage.setItem('username', data.username);
    sessionStorage.setItem('userId', data.user_id);

    authToken = data.token;
    currentUser = { username: data.username, userId: data.user_id };

    // Show chat
    showChat();
    connectWebSocket();

  } catch (error) {
    console.error('Login error:', error);
    errorEl.textContent = 'Network error. Please try again.';
  }
}

function handleLogout() {
  // Disable reconnection
  shouldReconnect = false;

  // Close WebSocket
  if (ws) {
    ws.close(1000); // Normal closure
    ws = null;
  }

  // Clear reconnect timeout
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout);
    reconnectTimeout = null;
  }

  // Clear room stats interval
  if (roomStatsInterval) {
    clearInterval(roomStatsInterval);
    roomStatsInterval = null;
  }

  // Clear auth data
  sessionStorage.removeItem('authToken');
  sessionStorage.removeItem('username');
  sessionStorage.removeItem('userId');

  authToken = null;
  currentUser = null;

  // Show auth screen
  showAuth();

  // Clear forms
  document.getElementById('login-username').value = '';
  document.getElementById('login-password').value = '';
  document.getElementById('register-username').value = '';
  document.getElementById('register-password').value = '';
  clearErrors();
}

// ============================================
// WEBSOCKET FUNCTIONS
// ============================================

let shouldReconnect = true;
let reconnectTimeout = null;

function connectWebSocket() {
  if (!authToken) {
    console.error('No auth token available');
    return;
  }

  // Close existing connection if any
  if (ws) {
    if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
      shouldReconnect = false;
      ws.close();
      // Wait a bit for the connection to fully close
      setTimeout(() => {
        shouldReconnect = true;
        createWebSocketConnection();
      }, 100);
      return;
    }
  }

  createWebSocketConnection();
}

function createWebSocketConnection() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const wsUrl = `${protocol}//${window.location.host}/ws?token=${authToken}&room=${currentRoom}`;

  ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    console.log('WebSocket connected to room:', currentRoom);
    // Clear any pending reconnect attempts
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }
  };

  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data);
      displayMessage(message);
    } catch (error) {
      console.error('Error parsing message:', error);
    }
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  ws.onclose = (event) => {
    console.log('WebSocket disconnected, code:', event.code);

    // Clear any existing reconnect timeout
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }

    // Only attempt to reconnect if:
    // 1. User is still authenticated
    // 2. shouldReconnect is true (not intentionally closed)
    // 3. Close code indicates an abnormal closure
    if (authToken && shouldReconnect && event.code !== 1000) {
      console.log('Will attempt to reconnect in 3 seconds...');
      reconnectTimeout = setTimeout(() => {
        console.log('Attempting to reconnect...');
        createWebSocketConnection();
      }, 3000);
    }
  };
}

function sendMessage(e) {
  e.preventDefault();

  const input = document.getElementById('message-input');
  const text = input.value.trim();

  if (!text || !ws || ws.readyState !== WebSocket.OPEN) {
    return;
  }

  const message = {
    user: currentUser.username,
    text: text,
    room: currentRoom,
    timestamp: new Date().toISOString(),
  };

  ws.send(JSON.stringify(message));
  input.value = '';
}

// ============================================
// MESSAGE FUNCTIONS
// ============================================

async function loadMessages() {
  if (!authToken) return;

  try {
    console.log('Loading messages for room:', currentRoom);
    const response = await fetch(`/api/messages?room=${currentRoom}&limit=50`, {
      headers: {
        'Authorization': `Bearer ${authToken}`,
      },
    });

    if (!response.ok) {
      console.error('Failed to load messages');
      return;
    }

    const messages = await response.json();
    console.log('Loaded messages:', messages.length, messages);

    const container = document.getElementById('messages');
    container.innerHTML = '';

    if (messages && messages.length > 0) {
      messages.forEach(msg => {
        displayMessage(msg, false);
      });
    } else {
      console.log('No messages found for room:', currentRoom);
    }

    // Scroll to bottom
    scrollToBottom();

  } catch (error) {
    console.error('Error loading messages:', error);
  }
}

function displayMessage(message, animate = true) {
  const container = document.getElementById('messages');
  const messageEl = document.createElement('div');

  const isOwnMessage = currentUser && message.user === currentUser.username;
  messageEl.className = `message ${isOwnMessage ? 'own-message' : ''}`;

  if (!animate) {
    messageEl.style.animation = 'none';
  }

  const avatar = message.user.charAt(0).toUpperCase();
  const timestamp = formatTimestamp(message.timestamp || message.created_at);

  messageEl.innerHTML = `
    <div class="message-avatar">${avatar}</div>
    <div class="message-content">
      <div class="message-header">
        <span class="message-username">${escapeHtml(message.user)}</span>
        <span class="message-time">${timestamp}</span>
      </div>
      <div class="message-text">${escapeHtml(message.text)}</div>
    </div>
  `;

  container.appendChild(messageEl);
  scrollToBottom();
}

function scrollToBottom() {
  const container = document.getElementById('messages');
  container.scrollTop = container.scrollHeight;
}

function formatTimestamp(timestamp) {
  if (!timestamp) return '';

  const date = new Date(timestamp);
  const now = new Date();
  const diff = now - date;

  // Less than 1 minute
  if (diff < 60000) {
    return 'Just now';
  }

  // Less than 1 hour
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000);
    return `${minutes}m ago`;
  }

  // Less than 24 hours
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000);
    return `${hours}h ago`;
  }

  // Format as time
  return date.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  });
}

function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

// ============================================
// ROOM FUNCTIONS
// ============================================

function switchRoom(room) {
  if (room === currentRoom) return;

  currentRoom = room;

  // Update UI
  document.querySelectorAll('.room-item').forEach(item => {
    item.classList.remove('active');
  });
  document.querySelector(`[data-room="${room}"]`).classList.add('active');

  // Update header
  const roomNames = {
    'general': { name: '# general', desc: 'Welcome to the general chat room' },
    'random': { name: '# random', desc: 'Random discussions and off-topic chat' },
    'tech': { name: '# tech', desc: 'Technology and programming discussions' },
  };

  const roomInfo = roomNames[room] || { name: `# ${room}`, desc: '' };
  document.getElementById('current-room-name').textContent = roomInfo.name;
  document.getElementById('room-description').textContent = roomInfo.desc;

  // Reconnect WebSocket with new room (connectWebSocket handles closing old connection)
  connectWebSocket();
  loadMessages();
  updateRoomCounts();
}

// ============================================
// UTILITY FUNCTIONS
// ============================================

// Handle Enter key in message input
document.addEventListener('DOMContentLoaded', () => {
  const messageInput = document.getElementById('message-input');
  if (messageInput) {
    messageInput.addEventListener('keypress', (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        sendMessage(e);
      }
    });
  }
});

// ============================================
// ROOM STATS FUNCTIONS
// ============================================

async function updateRoomCounts() {
  try {
    const response = await fetch('/api/room-stats');
    if (!response.ok) {
      console.error('Failed to fetch room stats');
      return;
    }

    const counts = await response.json();

    // Update room member counts in the UI
    const rooms = ['general', 'random', 'tech'];
    rooms.forEach(room => {
      const roomElement = document.querySelector(`[data-room="${room}"] .room-members`);
      if (roomElement) {
        const count = counts[room] || 0;
        roomElement.textContent = `${count} online`;
      }
    });

  } catch (error) {
    console.error('Error updating room counts:', error);
  }
}
