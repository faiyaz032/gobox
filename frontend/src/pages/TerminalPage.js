import React, { useEffect, useRef, useState } from 'react';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { getFingerprint } from '../utils/fingerprint';
import '@xterm/xterm/css/xterm.css';
import './TerminalPage.css';

const WS_BASE_URL = 'ws://localhost:8010/api/v1/box/connect';

// GoBox brand-themed terminal colors
// Ubuntu-themed terminal colors (Hyper.js style)
const GOBOX_THEME = {
  background: '#1A0A14', // Deep, dark aubergine/black
  foreground: '#FFFFFF',
  cursor: '#E95420', // Ubuntu Orange
  cursorAccent: '#1A0A14',
  selectionBackground: 'rgba(233, 84, 32, 0.4)',
  selectionForeground: '#FFFFFF',
  black: '#000000',
  red: '#E0211D',
  green: '#4E9A06',
  yellow: '#C4A000',
  blue: '#3465A4',
  magenta: '#75507B',
  cyan: '#06989A',
  white: '#D3D7CF',
  brightBlack: '#555753',
  brightRed: '#EF2929',
  brightGreen: '#8AE234',
  brightYellow: '#FCE94F',
  brightBlue: '#729FCF',
  brightMagenta: '#AD7FA8',
  brightCyan: '#34E2E2',
  brightWhite: '#EEEEEC',
};

const TerminalPage = () => {
  const terminalRef = useRef(null);
  const termRef = useRef(null);
  const wsRef = useRef(null);
  const fitAddonRef = useRef(null);
  
  // Connection State Tracking
  const [connectionStatus, setConnectionStatus] = useState('connecting');
  
  // Refs for strict lifecycle management
  const initialized = useRef(false);

  useEffect(() => {
    let isMounted = true;

    // Prevent double initialization in StrictMode
    if (initialized.current) return;
    initialized.current = true;

    console.log('[GoBox] Component Mounted');

    // 1. Initialize xterm.js
    const term = new Terminal({
      theme: GOBOX_THEME,
      fontFamily: '"JetBrains Mono", "Fira Code", monospace',
      fontSize: 14,
      lineHeight: 1.4,
      cursorBlink: true,
      cursorStyle: 'bar',
      allowTransparency: true,
      scrollback: 5000,
      convertEol: true,
    });

    const fitAddon = new FitAddon();
    const webLinksAddon = new WebLinksAddon();
    fitAddonRef.current = fitAddon;

    term.loadAddon(fitAddon);
    term.loadAddon(webLinksAddon);

    if (terminalRef.current) {
      term.open(terminalRef.current);
      requestAnimationFrame(() => {
        try {
          fitAddon.fit();
        } catch (e) {
          console.warn('Initial fit failed:', e);
        }
      });
    }
    
    termRef.current = term;
    term.writeln('  \x1b[33m⏳ Connecting to GoBox Container...\x1b[0m');

    // 2. Connect to WebSocket
    const connect = async () => {
      try {
        const fingerprint = await getFingerprint();
        
        // Safety: If unmounted during await, stop
        if (!isMounted) return;

        console.log('[GoBox] Using Fingerprint:', fingerprint);
        
        const wsUrl = `${WS_BASE_URL}?fingerprint=${encodeURIComponent(fingerprint)}`;
        console.log('[GoBox] Opening WebSocket:', wsUrl);
        
        const ws = new WebSocket(wsUrl);
        ws.binaryType = 'arraybuffer';
        wsRef.current = ws;

        ws.onopen = () => {
          if (!isMounted) {
            ws.close();
            return;
          }
          console.log('[GoBox] WebSocket Connected');
          setConnectionStatus('connected');
          
          term.reset(); 
          term.writeln('  \x1b[1;32m✅ Connected to GoBox!\x1b[0m\r\n');
          term.focus();
          // Trigger initial prompt
          ws.send('\n');
        };

        ws.onmessage = (event) => {
           if (!isMounted) return;
           if (event.data instanceof ArrayBuffer) {
             term.write(new Uint8Array(event.data));
           } else {
             term.write(event.data);
           }
        };

        ws.onclose = (event) => {
          if (!isMounted) return;
          console.log(`[GoBox] WebSocket Closed (Code: ${event.code})`);
          setConnectionStatus('disconnected');
          term.writeln('\r\n\x1b[1;33m⚠ Disconnected from server.\x1b[0m');
          wsRef.current = null;
        };

        ws.onerror = (error) => {
          console.error('[GoBox] WebSocket Error:', error);
          if (isMounted) setConnectionStatus('error');
        };

      } catch (err) {
        console.error('Failed to connect:', err);
        if (isMounted) setConnectionStatus('error');
      }
    };

    connect();

    // 3. Handle Terminal Input
    term.onData((data) => {
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(data);
      }
    });

    // 4. Handle Resize
    const handleResize = () => {
      if (fitAddonRef.current) {
        fitAddonRef.current.fit();
      }
    };
    window.addEventListener('resize', handleResize);

    // Cleanup Function (Runs on unmount)
    return () => {
      console.log('[GoBox] Cleaning up TerminalPage...');
      isMounted = false; // Prevent async callbacks
      initialized.current = false; 
      
      window.removeEventListener('resize', handleResize);
      
      if (wsRef.current) {
        console.log('[GoBox] Closing WebSocket...');
        wsRef.current.close();
        wsRef.current = null;
      }
      
      term.dispose();
      termRef.current = null;
    };
  }, []); // Empty dependency array = run once on mount

  const handleClose = () => {
    if (wsRef.current) wsRef.current.close();
    window.close();
  };

  const statusConfig = {
    connecting: { label: 'connecting...', dotClass: 'terminal-page-status-dot connecting' },
    connected: { label: 'connected', dotClass: 'terminal-page-status-dot connected' },
    disconnected: { label: 'disconnected', dotClass: 'terminal-page-status-dot disconnected' },
    reconnecting: { label: 'reconnecting...', dotClass: 'terminal-page-status-dot reconnecting' },
    error: { label: 'error', dotClass: 'terminal-page-status-dot error' },
  };
  const status = statusConfig[connectionStatus] || statusConfig.connecting;

  return (
    <div className="terminal-page">
      <div className="terminal-page-bar">
        <a href="/" className="terminal-page-brand">gobox</a>
        <div className="terminal-page-info">
          <div className="terminal-page-status">
            <span className={status.dotClass}></span>
            {status.label}
          </div>
          <button className="terminal-page-close" onClick={handleClose}>
            ✕ Destroy Box
          </button>
        </div>
      </div>
      <div className="terminal-page-main">
        <div className="terminal-window">
          <div className="terminal-window-header">
            <div className="terminal-window-dots">
              <span className="dot red"></span>
              <span className="dot yellow"></span>
              <span className="dot green"></span>
            </div>
            <div className="terminal-window-title">root@gobox: ~</div>
          </div>
          <div className="terminal-page-container" ref={terminalRef}></div>
        </div>
      </div>
    </div>
  );
};

export default TerminalPage;
