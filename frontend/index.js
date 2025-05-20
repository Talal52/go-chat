const { useState, useEffect } = React;

const ChatApp = () => {
    const [token, setToken] = useState(localStorage.getItem('token') || '');
    const [username, setUsername] = useState(localStorage.getItem('username') || '');
    const [messages, setMessages] = useState([]);
    const [newMessage, setNewMessage] = useState('');
    const [groups, setGroups] = useState([]);
    const [selectedGroup, setSelectedGroup] = useState(null);
    const [ws, setWs] = useState(null);
    const [signupData, setSignupData] = useState({ username: '', password: '' });
    const [loginData, setLoginData] = useState({ username: '', password: '', token: '' }); // Added token field
    const [error, setError] = useState('');

    const BASE_URL = 'http://localhost:8080';

    // Initialize WebSocket
    useEffect(() => {
        if (token) {
            const websocket = new WebSocket('ws://localhost:8081/ws');
            websocket.onopen = () => console.log('WebSocket connected');
            websocket.onmessage = (event) => {
                const msg = JSON.parse(event.data);
                setMessages((prev) => [...prev, msg]);
            };
            websocket.onclose = () => console.log('WebSocket disconnected');
            websocket.onerror = (error) => console.error('WebSocket error:', error);
            setWs(websocket);
            return () => websocket.close();
        }
    }, [token]);

    useEffect(() => {
        if (!token) return;
        const url = selectedGroup 
            ? `${BASE_URL}/api/messages?group_id=${selectedGroup}` 
            : `${BASE_URL}/api/messages`;
        console.log('Fetching messages from:', url);
        axios.get(url, {
            headers: { Authorization: `Bearer ${token}` }
        })
        .then(response => {
            console.log('Messages fetched:', response.data);
            setMessages(response.data);
            setError('');
        })
        .catch(error => {
            console.error('Error fetching messages:', error);
            setError('Failed to fetch messages: ' + (error.response?.data?.error || error.message));
        });
    }, [token, selectedGroup]);

    useEffect(() => {
        setGroups([
            { id: 'group1', name: 'General' },
            { id: 'group2', name: 'Friends' },
        ]);
    }, []);

    const handleSignup = () => {
        console.log('Signup button clicked:', signupData);
        axios.post(`${BASE_URL}/api/signup`, signupData)
            .then(response => {
                console.log('Signup successful:', response.data);
                alert('Signup successful! Please log in.');
                setSignupData({ username: '', password: '' });
                setError('');
            })
            .catch(error => {
                console.error('Signup error:', error);
                setError('Signup failed: ' + (error.response?.data?.error || error.message));
            });
    };

    const handleLogin = () => {
        console.log('Login button clicked:', loginData);
        if (!loginData.token) {
            setError('A valid JWT token is required for login.');
            return;
        }
        axios.post(`${BASE_URL}/api/login`, loginData, {
            headers: { Authorization: `Bearer ${loginData.token}` }
        })
        .then(response => {
            console.log('Login successful:', response.data);
            setToken(response.data.token);
            setUsername(loginData.username);
            localStorage.setItem('token', response.data.token);
            localStorage.setItem('username', loginData.username);
            setLoginData({ username: '', password: '', token: '' });
            setError('');
        })
        .catch(error => {
            console.error('Login error:', error);
            setError('Login failed: ' + (error.response?.data?.error || error.message));
        });
    };

    const handleLogout = () => {
        console.log('Logout button clicked');
        setToken('');
        setUsername('');
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        setMessages([]);
        setSelectedGroup(null);
        if (ws) ws.close();
        setError('');
    };

    const handleSendMessage = () => {
        console.log('Send button clicked:', newMessage);
        if (!newMessage.trim()) return;
        const message = {
            sender: username,
            content: newMessage,
            created_at: new Date().toISOString(),
            group_id: selectedGroup || undefined
        };
        axios.post(`${BASE_URL}/api/send-message`, message, {
            headers: { Authorization: `Bearer ${token}` }
        })
        .then(() => {
            console.log('Message sent successfully');
            setNewMessage('');
            setError('');
        })
        .catch(error => {
            console.error('Error sending message:', error);
            setError('Failed to send message: ' + (error.response?.data?.error || error.message));
        });
    };

    const handleGroupSelect = (groupId) => {
        console.log('Group selected:', groupId);
        setSelectedGroup(groupId);
    };

    // Input change handlers
    const handleSignupChange = (e) => {
        setSignupData({ ...signupData, [e.target.name]: e.target.value });
    };

    const handleLoginChange = (e) => {
        setLoginData({ ...loginData, [e.target.name]: e.target.value });
    };

    return (
        <div className="w-full max-w-2xl bg-white rounded-lg shadow-lg p-6">
            {error && (
                <div className="bg-red-100 text-red-700 p-2 mb-4 rounded">
                    {error}
                </div>
            )}
            {!token ? (
                <div>
                    <h2 className="text-2xl font-bold mb-4">Chat App</h2>
                    <div className="mb-6">
                        <h3 className="text-lg font-semibold mb-2">Sign Up</h3>
                        <input
                            type="text"
                            name="username"
                            placeholder="Username"
                            value={signupData.username}
                            onChange={handleSignupChange}
                            className="w-full p-2 mb-2 border rounded"
                        />
                        <input
                            type="password"
                            name="password"
                            placeholder="Password"
                            value={signupData.password}
                            onChange={handleSignupChange}
                            className="w-full p-2 mb-2 border rounded"
                        />
                        <button
                            onClick={handleSignup}
                            className="w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600"
                        >
                            Sign Up
                        </button>
                    </div>
                    <div>
                        <h3 className="text-lg font-semibold mb-2">Log In</h3>
                        <input
                            type="text"
                            name="username"
                            placeholder="Username"
                            value={loginData.username}
                            onChange={handleLoginChange}
                            className="w-full p-2 mb-2 border rounded"
                        />
                        <input
                            type="password"
                            name="password"
                            placeholder="Password"
                            value={loginData.password}
                            onChange={handleLoginChange}
                            className="w-full p-2 mb-2 border rounded"
                        />
                        <input
                            type="text"
                            name="token"
                            placeholder="JWT Token"
                            value={loginData.token}
                            onChange={handleLoginChange}
                            className="w-full p-2 mb-2 border rounded"
                        />
                        <button
                            onClick={handleLogin}
                            className="w-full bg-green-500 text-white p-2 rounded hover:bg-green-600"
                        >
                            Log In
                        </button>
                    </div>
                </div>
            ) : (
                <div>
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-2xl font-bold">Chat App - Welcome, {username}</h2>
                        <button
                            onClick={handleLogout}
                            className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
                        >
                            Log Out
                        </button>
                    </div>
                    <div className="mb-4">
                        <h3 className="text-lg font-semibold mb-2">Groups</h3>
                        <div className="flex space-x-2">
                            <button
                                onClick={() => handleGroupSelect(null)}
                                className={`px-4 py-2 rounded ${!selectedGroup ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}
                            >
                                All Messages
                            </button>
                            {groups.map(group => (
                                <button
                                    key={group.id}
                                    onClick={() => handleGroupSelect(group.id)}
                                    className={`px-4 py-2 rounded ${selectedGroup === group.id ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}
                                >
                                    {group.name}
                                </button>
                            ))}
                        </div>
                    </div>
                    <div className="h-96 overflow-y-auto p-4 bg-gray-50 rounded-lg mb-4">
                        {messages.map((message, index) => (
                            <div
                                key={message.created_at + index}
                                className={`mb-2 p-2 rounded ${message.sender === username ? 'bg-blue-100 ml-auto' : 'bg-gray-200 mr-auto'} max-w-xs`}
                            >
                                <p className="text-sm font-semibold">{message.sender}</p>
                                <p>{message.content}</p>
                                <p className="text-xs text-gray-500">{new Date(message.created_at).toLocaleTimeString()}</p>
                            </div>
                        ))}
                    </div>
                    <div className="flex space-x-2">
                        <input
                            type="text"
                            value={newMessage}
                            onChange={(e) => setNewMessage(e.target.value)}
                            placeholder="Type a message..."
                            className="flex-1 p-2 border rounded"
                        />
                        <button
                            onClick={handleSendMessage}
                            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
                        >
                            Send
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

try {
    ReactDOM.render(<ChatApp />, document.getElementById('root'));
} catch (error) {
    console.error('Error rendering React app:', error);
}