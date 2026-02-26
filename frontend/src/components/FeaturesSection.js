import React from 'react';
import { Terminal, Shield, Globe, Zap, Lock, Laptop } from 'lucide-react';
import './FeaturesSection.css';

const features = [
  {
    icon: <Terminal size={28} color="#E95420" />,
    title: 'Learn Linux Hands-On',
    description:
      'The best way to learn Linux is by doing. Practice real commands, manage files, install packages, and build confidence — all in a safe environment.',
  },
  {
    icon: <Shield size={28} color="#772953" />,
    title: 'Fully Isolated Boxes',
    description:
      'Every box runs in its own isolated container. Break things, experiment freely — nothing you do affects other users or the host system.',
  },
  {
    icon: <Globe size={28} color="#E95420" />,
    title: 'Full Network Access',
    description:
      'Practice real-world scenarios with full network connectivity. Use curl, wget, apt, ping, SSH, and any network tool just like a real server.',
  },
  {
    icon: <Zap size={28} color="#C7461B" />,
    title: 'Instant Boot',
    description:
      'No waiting around. Spin up a fresh Linux box in under 3 seconds. When you\'re done, just close the tab — it\'s fully disposable.',
  },
  {
    icon: <Lock size={28} color="#FF6347" />,
    title: 'Full Root Access',
    description:
      'Every box gives you full root privileges. Install system packages, configure services, manage users — nothing is off-limits.',
  },
  {
    icon: <Laptop size={28} color="#E95420" />,
    title: 'Browser-Based Terminal',
    description:
      'No SSH client needed. Access your Linux box directly from the browser with a full-featured terminal emulator. Works on any device.',
  },
];

function FeaturesSection() {
  return (
    <section className="features" id="features">
      <div className="features-header">
        <span className="features-label">Features</span>
        <h2 className="features-title">Everything You Need to Master Linux</h2>
        <p className="features-subtitle">
          GoBox gives you real Linux environments to practice, learn, and experiment — 
          without the hassle of setting up virtual machines or cloud servers.
        </p>
      </div>

      <div className="features-grid">
        {features.map((feature, index) => (
          <div className="feature-card" key={index} style={{ animationDelay: `${index * 0.1}s` }}>
            <div className="feature-icon">{feature.icon}</div>
            <h3 className="feature-card-title">{feature.title}</h3>
            <p className="feature-card-desc">{feature.description}</p>
          </div>
        ))}
      </div>
    </section>
  );
}

export default FeaturesSection;
