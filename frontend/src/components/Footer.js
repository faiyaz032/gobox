import React from 'react';
import { Github, Twitter, Mail, Heart } from 'lucide-react';
import './Footer.css';

function Footer() {
  return (
    <footer className="footer" id="about">
      <div className="footer-inner">
        <div className="footer-top">
          <div className="footer-brand-col">
            <span className="footer-brand">gobox</span>
            <p className="footer-tagline">
              Your isolated Linux playground in the cloud. 
              Learn, practice, and experiment — safely.
            </p>
            <div className="footer-socials">
              <a href="https://github.com/faiyaz032/gobox" aria-label="Github"><Github size={20} /></a>
              <a href="https://twitter.com" aria-label="Twitter"><Twitter size={20} /></a>
              <a href="mailto:hello@gobox.dev" aria-label="Email"><Mail size={20} /></a>
            </div>
          </div>
          
          <div className="footer-links-col">
            <h4>Product</h4>
            <a href="#features">Features</a>
            <a href="#how-it-works">How it Works</a>
            <a href="/terminal">Launch Box</a>
          </div>

          <div className="footer-links-col">
            <h4>Resources</h4>
            <a href="#">Documentation</a>
            <a href="#">Blog</a>
            <a href="#">Community</a>
          </div>

          <div className="footer-newsletter-col">
            <h4>Stay Updated</h4>
            <p>Get the latest updates and Linux tips.</p>
            <form className="footer-newsletter-form" onSubmit={(e) => e.preventDefault()}>
              <input type="email" placeholder="Enter your email" />
              <button type="submit">Subscribe</button>
            </form>
          </div>
        </div>

        <div className="footer-bottom">
          <span className="footer-copyright">
            © {new Date().getFullYear()} gobox. All rights reserved.
          </span>
          <span className="footer-made-with">
            Made with <Heart size={14} fill="#EF4444" color="#EF4444" style={{ margin: '0 4px' }} /> by developers
          </span>
        </div>
      </div>
    </footer>
  );
}

export default Footer;
