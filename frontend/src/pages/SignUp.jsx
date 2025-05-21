import { useState } from 'react';
import { signup } from '../auth';
import { useNavigate } from 'react-router-dom';

export default function SignUp({ onLogin }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSignup = async (e) => {
    e.preventDefault();
    if (!email || !password) {
      alert('Please fill in all fields');
      return;
    }
    try {
      await signup(email, password);
      alert('Signup successful. Please login.');
      navigate('/login');
    } catch (error) {
      alert('Signup failed: ' + (error.response?.data?.error || error.message));
    }
  };

  return (
    <form onSubmit={handleSignup}>
      <h2>Sign Up</h2>
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button type="submit">Sign Up</button>
      <button type="button" onClick={() => navigate('/login')}>
        Already account? Login
      </button>
    </form>
  );
}