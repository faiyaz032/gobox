import React, { useState, useEffect } from 'react';
import { ArrowRight, Github, Play } from 'lucide-react';
import './HeroSection.css';

function HeroSection() {
  const [text, setText] = useState('');
  const fullText = 'apt update && apt install -y htop';
  
  useEffect(() => {
    let index = 0;
    const timer = setInterval(() => {
      setText((prev) => {
        if (index < fullText.length) {
          index++;
          return fullText.slice(0, index);
        }
        return prev; // done
      });
      
      if (index >= fullText.length) {
        clearInterval(timer);
        // Optional: loop?
        // setTimeout(() => { index = 0; setText(''); }, 3000);
      }
    }, 50); // Typing speed
    
    return () => clearInterval(timer);
  }, []);

  const handleLaunchBox = () => {
    window.open('/terminal', '_blank');
  };

  return (
    <section className="hero">
      {/* Background Effects */}
      <div className="hero-bg">
        <div className="hero-orb hero-orb-1"></div>
        <div className="hero-orb hero-orb-2"></div>
        <div className="hero-orb hero-orb-3"></div>
        <div className="hero-dots"></div>
      </div>

      <div className="hero-content">
        {/* Left Column */}
        <div className="hero-text">
          <div className="hero-badge">
            <span className="hero-badge-dot"></span>
            Open Source Linux Playground
          </div>

          <h1 className="hero-title">
            Spin Up a{' '}
            <span className="hero-title-accent">Linux Box</span>{' '}
            in Seconds
          </h1>

          <p className="hero-description">
            Practice Linux commands, experiment with tools, and learn system administration — 
            all inside isolated, disposable containers with <strong>full root access</strong> and network connectivity. 
            No setup. No risk. Just launch and go.
          </p>

          <div className="hero-actions">
            <button className="btn-launch" onClick={handleLaunchBox} id="launch-box-btn">
              <Play size={20} fill="currentColor" />
              Launch a Box
              <ArrowRight size={18} className="btn-launch-arrow" />
            </button>
            <a href="https://github.com/faiyaz032/gobox" target="_blank" rel="noopener noreferrer" className="btn-secondary">
              <Github size={20} />
              View Source
            </a>
          </div>

          <div className="hero-stats">
            <div className="hero-stat">
              <div className="hero-stat-value">100%</div>
              <div className="hero-stat-label">Isolated Containers</div>
            </div>
            <div className="hero-stat">
              <div className="hero-stat-value">Full</div>
              <div className="hero-stat-label">Root Access</div>
            </div>
            <div className="hero-stat">
              <div className="hero-stat-value">&lt;3s</div>
              <div className="hero-stat-label">Boot Time</div>
            </div>
          </div>
        </div>

        {/* Right Column - Terminal Mockup */}
        <div className="hero-visual">
          <div className="terminal-mockup">
            <div className="terminal-header">
              <div className="terminal-dot red"></div>
              <div className="terminal-dot yellow"></div>
              <div className="terminal-dot green"></div>
              <div className="terminal-title">gobox — bash</div>
            </div>
            <div className="terminal-body">
              <div className="terminal-line">
                <span className="terminal-prompt">root@gobox:~$</span>
                <span className="terminal-command"> curl -s https://api.ipify.org</span>
              </div>
              <div className="terminal-output">203.0.113.42</div>
              <div className="terminal-line">
                <span className="terminal-prompt">root@gobox:~$</span>
                <span className="terminal-command"> {text}</span>
                <span className="terminal-cursor"></span>
              </div>
              
              {/* Show output only when typing finishes? Or mock it appearing later. 
                  For now, let's keep it simple and just show typing prompt. 
                  Alternatively, show output if text is full length. */}
              {text === fullText && (
                <>
                  <div className="terminal-output fade-in">Reading package lists... Done</div>
                  <div className="terminal-output fade-in delay-1">Setting up htop (3.2.1) ...</div>
                  <div className="terminal-line fade-in delay-2">
                    <span className="terminal-prompt">root@gobox:~$</span>
                    <span className="terminal-command"> neofetch</span>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default HeroSection;
