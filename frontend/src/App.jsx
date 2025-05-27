import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useState } from 'react';
import Login from './pages/Login';
import SignUp from './pages/SignUp';
import Chat from './pages/Chat';
import './App.css';

export default function App() {
  const [loggedIn, setLoggedIn] = useState(!!localStorage.getItem('accessToken'));

  return (
    <Router>
      <div className="container">
        <Routes>
          <Route
            path="/login"
            element={
              !loggedIn ? (
                <Login onLogin={() => setLoggedIn(true)} />
              ) : (
                <Navigate to="/chat" />
              )
            }
          />
          <Route
            path="/signup"
            element={
              !loggedIn ? (
                <SignUp onLogin={() => setLoggedIn(true)} />
              ) : (
                <Navigate to="/chat" />
              )
            }
          />
          <Route
            path="/chat"
            element={
              loggedIn ? (
                <Chat onLogout={() => setLoggedIn(false)} />
              ) : (
                <Navigate to="/login" />
              )
            }
          />
          <Route path="/" element={<Navigate to="/login" />} />
        </Routes>
      </div>
    </Router>
  );
}