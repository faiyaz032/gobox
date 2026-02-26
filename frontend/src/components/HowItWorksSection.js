import React from 'react';
import { Play, Terminal, RefreshCcw } from 'lucide-react';
import './HowItWorksSection.css';

const steps = [
  {
    icon: <Play size={32} color="#E95420" />,
    title: "1. Launch",
    description: "Click the launch button and a fresh container is provisioned just for you."
  },
  {
    icon: <Terminal size={32} color="#E95420" />,
    title: "2. Practice",
    description: "Use the terminal to practice Bash scripting, file manipulation, and more."
  },
  {
    icon: <RefreshCcw size={32} color="#E95420" />,
    title: 'Control',
    description: 'Typing feels native. Full root access to install, configure, and hack.',
  },
];

function HowItWorksSection() {
  return (
    <section className="how-it-works">
      <div className="hiw-header">
        <span className="hiw-label">Workflow</span>
        <h2 className="hiw-title">From Click to Command Line</h2>
      </div>

      <div className="hiw-steps">
        <div className="hiw-line"></div>
        {steps.map((step, index) => (
          <div className="hiw-step" key={index} style={{ animationDelay: `${index * 0.2}s` }}>
            <div className="hiw-icon-wrapper">
              <div className="hiw-icon">{step.icon}</div>
            </div>
            <h3 className="hiw-step-title">{step.title}</h3>
            <p className="hiw-step-desc">{step.description}</p>
          </div>
        ))}
      </div>
    </section>
  );
}

export default HowItWorksSection;
