import React, { useEffect, useState } from 'react';
import { BarChart2, Activity, CheckCircle2, XCircle, Clock, RefreshCw } from 'lucide-react';
import { getMetrics } from '../utils/api';

export default function SystemMetrics() {
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);

  const fetchMetrics = async () => {
    setLoading(true);
    const data = await getMetrics();
    setMetrics(data);
    setLoading(false);
  };

  useEffect(() => {
    fetchMetrics();
    const interval = setInterval(fetchMetrics, 5000);
    return () => clearInterval(interval);
  }, []);

  const formatUptime = (seconds) => {
    if (!seconds) return '0s';
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);
    if (hrs > 0) return `${hrs}h ${mins}m ${secs}s`;
    if (mins > 0) return `${mins}m ${secs}s`;
    return `${secs}s`;
  };

  return (
    <div className="glass-panel p-6 mb-12 animate-fade-in">
      <div className="flex items-center justify-between mb-6 pb-4 border-b border-white/5">
        <div>
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <BarChart2 className="text-indigo-400" size={20} />
            Live System Observability & Performance Metrics
          </h2>
          <p className="text-xs text-gray-400">
            Real-time metrics exposed via <span className="font-mono text-indigo-300">GET /metrics</span>.
          </p>
        </div>

        <button
          onClick={fetchMetrics}
          className="btn-secondary text-xs"
        >
          <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
          Refresh Metrics
        </button>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Metric Card 1 */}
        <div className="bg-slate-900/60 p-4 rounded-xl border border-white/5 flex items-center gap-4">
          <div className="p-3 rounded-xl bg-indigo-500/10 text-indigo-400">
            <Clock size={24} />
          </div>
          <div>
            <span className="text-xs text-gray-400 block">Server Uptime</span>
            <span className="text-lg font-bold text-white font-mono">
              {formatUptime(metrics?.uptimeSeconds)}
            </span>
          </div>
        </div>

        {/* Metric Card 2 */}
        <div className="bg-slate-900/60 p-4 rounded-xl border border-white/5 flex items-center gap-4">
          <div className="p-3 rounded-xl bg-cyan-500/10 text-cyan-400">
            <Activity size={24} />
          </div>
          <div>
            <span className="text-xs text-gray-400 block">Total Processed</span>
            <span className="text-lg font-bold text-white font-mono">
              {metrics?.totalRequests ?? 0}
            </span>
          </div>
        </div>

        {/* Metric Card 3 */}
        <div className="bg-slate-900/60 p-4 rounded-xl border border-white/5 flex items-center gap-4">
          <div className="p-3 rounded-xl bg-emerald-500/10 text-emerald-400">
            <CheckCircle2 size={24} />
          </div>
          <div>
            <span className="text-xs text-emerald-400 block">Valid Receipts</span>
            <span className="text-lg font-bold text-emerald-400 font-mono">
              {metrics?.validReceipts ?? 0}
            </span>
          </div>
        </div>

        {/* Metric Card 4 */}
        <div className="bg-slate-900/60 p-4 rounded-xl border border-white/5 flex items-center gap-4">
          <div className="p-3 rounded-xl bg-rose-500/10 text-rose-400">
            <XCircle size={24} />
          </div>
          <div>
            <span className="text-xs text-rose-400 block">Failed / Invalid</span>
            <span className="text-lg font-bold text-rose-400 font-mono">
              {metrics?.failedReceipts ?? 0}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
