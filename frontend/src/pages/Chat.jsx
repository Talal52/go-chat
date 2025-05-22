import ChatBox from '../components/ChatBox';

export default function Chat({ onLogout }) {
  return (
    <div>
      <h2>Chat-App</h2>
      <ChatBox />
      <button
        onClick={() => {
          localStorage.removeItem('accessToken');
          onLogout();
        }}
      >
        Logout
      </button>
    </div>
  );
}