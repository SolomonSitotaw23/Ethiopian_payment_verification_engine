import React, { useEffect, useState } from 'react';
import { ShieldCheck, Activity, Terminal, Layers, BarChart2 } from 'lucide-react';
import { getHealth } from '../utils/api';

export default function Navbar({ activeTab, setActiveTab }) {
  const [systemHealth, setSystemHealth] = useState({ status: 'CHECKING' });

  useEffect(() => {
    const checkHealth = async () => {
      const data = await getHealth();
      setSystemHealth(data);
    };
    checkHealth();
    const interval = setInterval(checkHealth, 10000);
    return () => clearInterval(interval);
  }, []);

  const isOnline = systemHealth.status === 'UP';

  return (
    <header className="glass-panel sticky top-0 z-50 px-6 py-4 mb-8 border-b border-[rgba(255,255,255,0.08)]">
      <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4">
        {/* Brand */}
        <div className="flex items-center gap-3">
          <div className="p-2.5 rounded-xl bg-indigo-600/20 border border-indigo-500/30 text-indigo-400">
            <ShieldCheck size={26} />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white flex items-center gap-2">
              Payment Verifier <span className="text-xs px-2 py-0.5 rounded-full bg-indigo-500/20 text-indigo-300 font-mono">Go API v1.0</span>
            </h1>
            <p className="text-xs text-gray-400">Ethiopian Digital Payment Receipt Verification Engine</p>
          </div>
        </div>

        {/* Navigation Tabs */}
        <nav className="flex items-center gap-2 bg-slate-900/60 p-1.5 rounded-xl border border-white/5">
          <button
            onClick={() => setActiveTab('single')}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              activeTab === 'single'
                ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            }`}
          >
            <ShieldCheck size={16} />
            Single Verifier
          </button>

          <button
            onClick={() => setActiveTab('batch')}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              activeTab === 'batch'
                ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            }`}
          >
            <Layers size={16} />
            Batch Verifier
          </button>

          <button
            onClick={() => setActiveTab('code')}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              activeTab === 'code'
                ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            }`}
          >
            <Terminal size={16} />
            Code Generator
          </button>

          <button
            onClick={() => setActiveTab('metrics')}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
              activeTab === 'metrics'
                ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            }`}
          >
            <BarChart2 size={16} />
            Observability
          </button>
        </nav>

        {/* System Health Badge */}
        <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-slate-900/80 border border-white/10 text-xs">
          <span className={`w-2.5 h-2.5 rounded-full ${isOnline ? 'bg-emerald-500 animate-pulse' : 'bg-rose-500'}`} />
          <span className="text-gray-300 font-mono">
            {isOnline ? 'API Operational' : 'API Offline'}
          </span>
          <Activity size={14} className="text-gray-400 ml-1" />
        </div>
      </div>
    </header>
  );
}
