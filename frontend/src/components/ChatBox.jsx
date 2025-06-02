import { useEffect, useState, useRef } from 'react';
import axiosInstance from '../utils/axiosInstance';
import { WS_URL } from '../config';

export default function ChatBox() {
  const [messages, setMessages] = useState([]);
  const [newMsg, setNewMsg] = useState('');
  const [receiverId, setReceiverId] = useState('');
  const ws = useRef(null);
  const [users, setUsers] = useState([]);

useEffect(() => {
  axiosInstance.get('/api/users').then((response) => {
    setUsers(response.data);
  });
}, []);

  useEffect(() => {
    const token = localStorage.getItem('accessToken');
    if (!token) return;

    const wsUrl = `${WS_URL}?token=${token}`;
    ws.current = new WebSocket(wsUrl);

    ws.current.onopen = () => {
      console.log('WebSocket connected');
    };

    ws.current.onmessage = (event) => {
      const receivedMessage = JSON.parse(event.data);
      console.log('WebSocket message received:', receivedMessage);
      setMessages((prev) => [...prev, receivedMessage]);
    };

    ws.current.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.current.onclose = () => {
      console.log('WebSocket closed');
    };

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, []);

  const fetchMessages = async () => {
    try {
      const res = await axiosInstance.get('/api/messages', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('accessToken')}`,
        },
      });
      setMessages(res.data);
    } catch (error) {
      console.error('Error fetching messages:', error);
    }
  };

  useEffect(() => {
    fetchMessages();
  }, []);

  const handleSend = async () => {
    if (newMsg.trim() === '' || receiverId.trim() === '') return;

    const msgPayload = {
      receiver_id: parseInt(receiverId),
      message: newMsg,
    };
    console.log('Sending message:', msgPayload);

    try {
      await axiosInstance.post('/api/send-message', msgPayload, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('accessToken')}`,
        },
      });

      const senderId = localStorage.getItem('userId') || 'Unknown';
      setMessages((prev) => [...prev, { ...msgPayload, sender_id: senderId }]);
      setNewMsg('');
    } catch (err) {
      console.error('Error sending message:', err);
    }
  };

  return (
    <div>
      <h3>Send Message</h3>
      <select
      value={receiverId}
      onChange={(e) => setReceiverId(e.target.value)}
    >
      <option value="">Select a user</option>
      {users.map((user) => (
        <option key={user.id} value={user.id}>{user.email}</option>
      ))}
    </select>
      <input
        placeholder="Receiver ID"
        value={receiverId}
        onChange={(e) => setReceiverId(e.target.value)}
        style={{ width: '200px', marginRight: '10px' }}
      />
      <input
        placeholder="Type a message"
        value={newMsg}
        onChange={(e) => setNewMsg(e.target.value)}
        style={{ width: '300px', marginRight: '10px' }}
      />
      <button onClick={handleSend}>Send</button>

      <hr />
      <h3>Messages</h3>
      <div>
        {messages.map((msg, i) => (
          <div key={i}>
            <strong>From {msg.sender_id} to {msg.receiver_id}:</strong> {msg.message}
          </div>
        ))}
      </div>
    </div>
  );
}