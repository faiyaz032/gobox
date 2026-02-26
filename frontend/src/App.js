import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from './context/ThemeContext';
import LandingPage from './pages/LandingPage';
import TerminalPage from './pages/TerminalPage';
import './App.css';

function App() {
  return (
    <ThemeProvider>
      <Router>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/terminal" element={<TerminalPage />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;
