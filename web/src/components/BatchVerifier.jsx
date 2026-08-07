import React, { useState } from 'react';
import { Layers, Play, CheckCircle2, XCircle, RefreshCw, Download, Send } from 'lucide-react';
import { verifyBatchReceipts } from '../utils/api';

export default function BatchVerifier() {
  const [receiptsText, setReceiptsText] = useState(
    'CJP9OSP9WZ\nFT25292FRPWD89873710\nFT25284X11PS79328\nAB1234567890'
  );
  const [amount, setAmount] = useState('100');
  const [callbackUrl, setCallbackUrl] = useState('');
  const [proxy, setProxy] = useState(false);

  const [loading, setLoading] = useState(false);
  const [batchResult, setBatchResult] = useState(null);
  const [error, setError] = useState(null);

  const handleBatchSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setBatchResult(null);

    const receiptList = receiptsText
      .split('\n')
      .map((r) => r.trim())
      .filter((r) => r.length > 0);

    if (receiptList.length === 0) {
      setError('Please enter at least one receipt ID or URL');
      setLoading(false);
      return;
    }

    const payload = {
      receipt: receiptList,
      expected: {
        amount: amount ? parseFloat(amount) : undefined,
      },
      defaultVerification: true,
      proxy,
      callbackUrl: callbackUrl.trim() || undefined,
    };

    try {
      const data = await verifyBatchReceipts(payload);
      setBatchResult(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleExportCSV = () => {
    if (!batchResult || !batchResult.result) return;
    let csvContent = 'data:text/csv;charset=utf-8,Receipt ID,Status,Error\n';
    
    batchResult.result.forEach((item) => {
      csvContent += `"${item.receiptId}","Valid",""\n`;
    });
    
    if (batchResult.failed) {
      batchResult.failed.forEach((f) => {
        csvContent += `"${f.receiptId}","Invalid","${f.error}"\n`;
      });
    }

    const encodedUri = encodeURI(csvContent);
    const link = document.createElement('a');
    link.setAttribute('href', encodedUri);
    link.setAttribute('download', 'payment_verification_report.csv');
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <div className="glass-panel p-6 mb-12">
      <div className="flex items-center justify-between mb-6 pb-4 border-b border-white/5">
        <div>
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <Layers className="text-indigo-400" size={20} />
            Parallel Batch Receipt Verifier
          </h2>
          <p className="text-xs text-gray-400">
            Verify up to 10+ Telebirr, CBE, BOA, and Amhara receipts simultaneously using Go worker pools.
          </p>
        </div>
        <span className="text-xs px-2.5 py-1 rounded-full bg-indigo-500/10 border border-indigo-500/30 text-indigo-300 font-mono">
          POST /api/verify/batch
        </span>
      </div>

      <form onSubmit={handleBatchSubmit} className="space-y-6">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
          {/* Input list */}
          <div className="lg:col-span-8">
            <label className="form-label">
              Receipt IDs or URLs (One per line)
            </label>
            <textarea
              rows={6}
              value={receiptsText}
              onChange={(e) => setReceiptsText(e.target.value)}
              placeholder="Paste receipt IDs or URLs here..."
              className="form-input font-mono text-xs leading-relaxed"
              required
            />
          </div>

          {/* Config column */}
          <div className="lg:col-span-4 space-y-4">
            <div>
              <label className="form-label">Expected Price / Amount (Birr)</label>
              <input
                type="number"
                step="any"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="e.g. 100"
                className="form-input"
              />
            </div>

            <div>
              <label className="form-label">Async Webhook URL (Optional)</label>
              <input
                type="url"
                value={callbackUrl}
                onChange={(e) => setCallbackUrl(e.target.value)}
                placeholder="https://your-api.com/webhooks"
                className="form-input text-xs"
              />
            </div>

            <div className="pt-2">
              <label className="flex items-center gap-2 text-xs font-semibold text-gray-300 cursor-pointer">
                <input
                  type="checkbox"
                  checked={proxy}
                  onChange={(e) => setProxy(e.target.checked)}
                  className="rounded border-gray-700 text-indigo-600 focus:ring-indigo-500 h-4 w-4"
                />
                Use Telebirr Proxy for Batch
              </label>
            </div>
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="btn-primary w-full justify-center py-3 text-base"
        >
          {loading ? (
            <>
              <RefreshCw size={20} className="animate-spin" />
              Processing Parallel Batch Receipts...
            </>
          ) : (
            <>
              <Play size={20} />
              Execute Concurrent Batch Verification
            </>
          )}
        </button>
      </form>

      {/* Error message */}
      {error && (
        <div className="mt-6 p-4 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-400 animate-fade-in">
          <p className="text-sm font-semibold">{error}</p>
        </div>
      )}

      {/* Batch Results View */}
      {batchResult && (
        <div className="mt-8 pt-6 border-t border-white/5 space-y-6 animate-fade-in">
          {/* Webhook notification banner if 202 Async */}
          {batchResult.callbackUrl && (
            <div className="p-4 rounded-xl bg-indigo-500/10 border border-indigo-500/30 text-indigo-300 flex items-center gap-3">
              <Send size={20} />
              <div>
                <h4 className="text-sm font-bold">Asynchronous Batch Dispatched</h4>
                <p className="text-xs text-gray-400">
                  Results will be posted to <span className="font-mono text-white">{batchResult.callbackUrl}</span> upon completion.
                </p>
              </div>
            </div>
          )}

          {/* Synchronous Summary Counters */}
          {batchResult.summary && (
            <div className="flex flex-wrap items-center justify-between gap-4">
              <div className="flex items-center gap-6">
                <div>
                  <span className="text-xs text-gray-400 block">Total Items:</span>
                  <span className="text-xl font-bold text-white">
                    {batchResult.summary.total}
                  </span>
                </div>

                <div>
                  <span className="text-xs text-emerald-400 block">Valid Receipts:</span>
                  <span className="text-xl font-bold text-emerald-400">
                    {batchResult.summary.valid}
                  </span>
                </div>

                <div>
                  <span className="text-xs text-rose-400 block">Failed / Invalid:</span>
                  <span className="text-xl font-bold text-rose-400">
                    {batchResult.summary.invalid}
                  </span>
                </div>
              </div>

              <button
                onClick={handleExportCSV}
                className="btn-secondary text-xs"
              >
                <Download size={14} />
                Export CSV Report
              </button>
            </div>
          )}

          {/* Results Table */}
          {batchResult.result && (
            <div className="overflow-x-auto bg-slate-900/60 rounded-xl border border-white/5">
              <table className="w-full text-left text-xs">
                <thead className="bg-slate-800/50 text-gray-400 uppercase font-mono text-[11px]">
                  <tr>
                    <th className="p-3">Receipt ID / URL</th>
                    <th className="p-3">Provider</th>
                    <th className="p-3">Status</th>
                    <th className="p-3">Details / Mismatch</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/5">
                  {/* Valid Results */}
                  {batchResult.result.map((res, idx) => (
                    <tr key={`valid-${idx}`} className="hover:bg-white/5">
                      <td className="p-3 font-mono text-gray-200">{res.receiptId}</td>
                      <td className="p-3 font-semibold text-indigo-400">{res.provider}</td>
                      <td className="p-3">
                        <span className="badge-valid">
                          <CheckCircle2 size={12} /> VALID
                        </span>
                      </td>
                      <td className="p-3 text-gray-400">Verified Amount: {res.parsed?.amount} ETB</td>
                    </tr>
                  ))}

                  {/* Failed Results */}
                  {batchResult.failed &&
                    batchResult.failed.map((f, idx) => (
                      <tr key={`failed-${idx}`} className="hover:bg-white/5 bg-rose-950/10">
                        <td className="p-3 font-mono text-gray-200">{f.receiptId}</td>
                        <td className="p-3 font-semibold text-gray-400">
                          {f.details?.provider || 'Unknown'}
                        </td>
                        <td className="p-3">
                          <span className="badge-mismatch">
                            <XCircle size={12} /> INVALID
                          </span>
                        </td>
                        <td className="p-3 text-rose-300">{f.error}</td>
                      </tr>
                    ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
