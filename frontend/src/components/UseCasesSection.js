import React from 'react';
import { BookOpen, Server, Code } from 'lucide-react';
import './UseCasesSection.css';

const cases = [
  {
    icon: <BookOpen size={24} />,
    title: 'Students & Learners',
    items: [
      'Practice shell commands risk-free',
      'Learn file permissions & user management',
      'Experiment with breaking changes',
    ],
    color: '#E95420', // orange
  },
  {
    icon: <Server size={24} />,
    title: 'DevOps Engineers',
    items: [
      'Test installation scripts',
      'Debug network configurations',
      'Verify package dependencies',
    ],
    color: '#772953', // aubergine
  },
  {
    icon: <Code size={24} />,
    title: 'Developers',
    items: [
      'Quick environment for CLI tools',
      'Test code in isolated Linux',
      'No local Docker setup required',
    ],
    color: '#C7461B', // auburn
  },
];

function UseCasesSection() {
  return (
    <section className="use-cases">
      <div className="cases-content">
        <div className="cases-text">
          <span className="cases-label">Use Cases</span>
          <h2 className="cases-title">Who is GoBox For?</h2>
          <p className="cases-desc">
             Whether you're writing your first `ls` command or debugging complex pipelines, 
             GoBox provides the instant, clean slate you need.
          </p>
        </div>
        
        <div className="cases-grid">
          {cases.map((useCase, index) => (
            <div className="case-card" key={index} style={{ '--accent': useCase.color }}>
              <div className="case-header">
                <div className="case-icon">{useCase.icon}</div>
                <h3 className="case-title">{useCase.title}</h3>
              </div>
              <ul className="case-list">
                {useCase.items.map((item, i) => (
                  <li key={i}>{item}</li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

export default UseCasesSection;
