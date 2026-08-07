import React, { useState } from 'react';
import { Play, CheckCircle2, XCircle, AlertTriangle, ShieldAlert, Sparkles, RefreshCw, ChevronRight } from 'lucide-react';
import { verifySingleReceipt } from '../utils/api';

export default function SingleVerifier({ sampleData }) {
  const [receipt, setReceipt] = useState(sampleData?.sampleId || 'CJP9OSP9WZ');
  const [provider, setProvider] = useState(sampleData?.id || 'auto');
  const [amount, setAmount] = useState('100');
  const [minAmount, setMinAmount] = useState('');
  const [maxAmount, setMaxAmount] = useState('');
  const [recipientAccount, setRecipientAccount] = useState(sampleData?.samplePhone || '0912345678');
  const [recipientName, setRecipientName] = useState(sampleData?.sampleName || 'Abrham Yalew');
  const [paymentYear, setPaymentYear] = useState('2025');
  const [paymentMonth, setPaymentMonth] = useState('12');
  const [strict, setStrict] = useState(true);
  const [proxy, setProxy] = useState(false);

  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);

  // Sync if sampleData changes
  React.useEffect(() => {
    if (sampleData) {
      setReceipt(sampleData.sampleId);
      setProvider(sampleData.id);
      setRecipientAccount(sampleData.samplePhone);
      setRecipientName(sampleData.sampleName);
    }
  }, [sampleData]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);

    const payload = {
      receipt,
      provider: provider === 'auto' ? undefined : provider,
      expected: {
        amount: amount ? parseFloat(amount) : undefined,
        minAmount: minAmount ? parseFloat(minAmount) : undefined,
        maxAmount: maxAmount ? parseFloat(maxAmount) : undefined,
        recipientAccount: recipientAccount || undefined,
        recipientName: recipientName || undefined,
        paymentYear: paymentYear || undefined,
        paymentMonth: paymentMonth || undefined,
      },
      defaultVerification: true,
      strict,
      proxy,
    };

    try {
      const data = await verifySingleReceipt(payload);
      setResult(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 mb-12">
      {/* Form Panel */}
      <div className="lg:col-span-6 glass-panel p-6">
        <div className="flex items-center justify-between mb-6 pb-4 border-b border-white/5">
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <Sparkles className="text-indigo-400" size={20} />
            Receipt Verification Request
          </h2>
          <span className="text-xs px-2.5 py-1 rounded-full bg-indigo-500/10 border border-indigo-500/30 text-indigo-300 font-mono">
            POST /api/verify
          </span>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Receipt Input */}
          <div>
            <label className="form-label">Receipt ID or Full Receipt URL</label>
            <input
              type="text"
              value={receipt}
              onChange={(e) => setReceipt(e.target.value)}
              placeholder="e.g. CJP9OSP9WZ or FT25292FRPWD89873710"
              className="form-input font-mono"
              required
            />
          </div>

          {/* Provider Selection */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="form-label">Provider Routing</label>
              <select
                value={provider}
                onChange={(e) => setProvider(e.target.value)}
                className="form-input"
              >
                <option value="auto">✨ Auto-Detect Provider</option>
                <option value="telebirr">Telebirr Wallet</option>
                <option value="cbe">CBE Bank</option>
                <option value="boa">Bank of Abyssinia (BOA)</option>
                <option value="amharabank">Amhara Bank</option>
              </select>
            </div>

            <div>
              <label className="form-label">Expected Amount (Birr)</label>
              <input
                type="number"
                step="any"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="e.g. 100 or 230"
                className="form-input"
              />
            </div>
          </div>

          {/* Range Amounts */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="form-label">Min Amount (Optional)</label>
              <input
                type="number"
                step="any"
                value={minAmount}
                onChange={(e) => setMinAmount(e.target.value)}
                placeholder="e.g. 100"
                className="form-input text-xs"
              />
            </div>
            <div>
              <label className="form-label">Max Amount (Optional)</label>
              <input
                type="number"
                step="any"
                value={maxAmount}
                onChange={(e) => setMaxAmount(e.target.value)}
                placeholder="e.g. 500"
                className="form-input text-xs"
              />
            </div>
          </div>

          {/* Account & Name */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="form-label">Target Account / Phone</label>
              <input
                type="text"
                value={recipientAccount}
                onChange={(e) => setRecipientAccount(e.target.value)}
                placeholder="0912345678 or 1000..."
                className="form-input"
              />
            </div>

            <div>
              <label className="form-label">Target Name</label>
              <input
                type="text"
                value={recipientName}
                onChange={(e) => setRecipientName(e.target.value)}
                placeholder="Abrham Yalew"
                className="form-input"
              />
            </div>
          </div>

          {/* Date Options */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="form-label">Expected Year</label>
              <input
                type="text"
                value={paymentYear}
                onChange={(e) => setPaymentYear(e.target.value)}
                placeholder="2025"
                className="form-input"
              />
            </div>

            <div>
              <label className="form-label">Expected Month</label>
              <input
                type="text"
                value={paymentMonth}
                onChange={(e) => setPaymentMonth(e.target.value)}
                placeholder="12"
                className="form-input"
              />
            </div>
          </div>

          {/* Toggles */}
          <div className="pt-2 flex items-center justify-between bg-slate-900/40 p-3 rounded-xl border border-white/5">
            <label className="flex items-center gap-2 text-xs font-semibold text-gray-300 cursor-pointer">
              <input
                type="checkbox"
                checked={strict}
                onChange={(e) => setStrict(e.target.checked)}
                className="rounded border-gray-700 text-indigo-600 focus:ring-indigo-500 h-4 w-4"
              />
              Strict Mode (Fail HTTP 400 on Mismatch)
            </label>

            <label className="flex items-center gap-2 text-xs font-semibold text-gray-300 cursor-pointer">
              <input
                type="checkbox"
                checked={proxy}
                onChange={(e) => setProxy(e.target.checked)}
                className="rounded border-gray-700 text-indigo-600 focus:ring-indigo-500 h-4 w-4"
              />
              Use Telebirr Proxy
            </label>
          </div>

          {/* Submit */}
          <button
            type="submit"
            disabled={loading}
            className="btn-primary w-full justify-center py-3 text-base mt-2"
          >
            {loading ? (
              <>
                <RefreshCw size={20} className="animate-spin" />
                Fetching & Validating Live Receipt...
              </>
            ) : (
              <>
                <Play size={20} />
                Verify Payment Receipt Now
              </>
            )}
          </button>
        </form>
      </div>

      {/* Result Display Panel */}
      <div className="lg:col-span-6 glass-panel p-6 flex flex-col justify-between">
        <div>
          <div className="flex items-center justify-between mb-6 pb-4 border-b border-white/5">
            <h3 className="text-lg font-bold text-white flex items-center gap-2">
              Verification Results Breakdown
            </h3>
            {result && (
              <span
                className={
                  result.status === 'valid' ? 'badge-valid' : 'badge-mismatch'
                }
              >
                {result.status === 'valid' ? (
                  <>
                    <CheckCircle2 size={14} /> VALID RECEIPT
                  </>
                ) : (
                  <>
                    <XCircle size={14} /> MISMATCH DETECTED
                  </>
                )}
              </span>
            )}
          </div>

          {!result && !error && !loading && (
            <div className="h-64 flex flex-col items-center justify-center text-center p-6 border-2 border-dashed border-white/10 rounded-xl">
              <Sparkles size={40} className="text-gray-600 mb-3" />
              <h4 className="text-sm font-semibold text-gray-300 mb-1">
                Ready for Verification
              </h4>
              <p className="text-xs text-gray-500 max-w-xs">
                Enter a Telebirr, CBE, BOA, or Amhara Bank receipt ID on the left and click "Verify Payment Receipt Now".
              </p>
            </div>
          )}

          {error && (
            <div className="p-4 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-400 animate-fade-in">
              <div className="flex items-center gap-2 font-bold mb-1">
                <AlertTriangle size={18} />
                Verification Error
              </div>
              <p className="text-sm">{error}</p>
            </div>
          )}

          {result && (
            <div className="space-y-6 animate-fade-in">
              {/* Message Banner */}
              <div
                className={`p-4 rounded-xl border ${
                  result.status === 'valid'
                    ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300'
                    : 'bg-rose-500/10 border-rose-500/30 text-rose-300'
                }`}
              >
                <p className="text-sm font-semibold">{result.message}</p>
              </div>

              {/* Parsed Metadata Grid */}
              <div className="bg-slate-900/60 p-4 rounded-xl border border-white/5 space-y-3">
                <h4 className="text-xs font-bold text-gray-400 uppercase tracking-wider">
                  Live Parsed Receipt Data ({result.provider})
                </h4>

                <div className="grid grid-cols-2 gap-3 text-xs">
                  <div>
                    <span className="text-gray-500 block">Amount Paid:</span>
                    <span className="font-bold text-white text-sm">
                      {result.parsed?.amount} ETB
                    </span>
                  </div>

                  <div>
                    <span className="text-gray-500 block">Recipient Account:</span>
                    <span className="font-mono text-gray-200">
                      {result.parsed?.accountNumber || 'N/A'}
                    </span>
                  </div>

                  <div>
                    <span className="text-gray-500 block">Recipient Name:</span>
                    <span className="font-semibold text-gray-200">
                      {result.parsed?.recipientName || 'N/A'}
                    </span>
                  </div>

                  <div>
                    <span className="text-gray-500 block">Transaction Date:</span>
                    <span className="font-mono text-gray-200">
                      {result.parsed?.date || 'N/A'}
                    </span>
                  </div>
                </div>
              </div>

              {/* Check Breakdown */}
              {result.checks && (
                <div>
                  <h4 className="text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">
                    Individual Verification Checks
                  </h4>
                  <div className="grid grid-cols-2 gap-2">
                    <CheckItem
                      label="Amount Match"
                      matched={result.checks.amountMatched}
                    />
                    <CheckItem
                      label="Account Match"
                      matched={result.checks.accountNumberMatched}
                    />
                    <CheckItem
                      label="Recipient Name"
                      matched={result.checks.recipientNameMatched}
                    />
                    <CheckItem
                      label="Payment Date Window"
                      matched={result.checks.dateMatched}
                    />
                  </div>
                </div>
              )}

              {/* Mismatches List */}
              {result.mismatches && result.mismatches.length > 0 && (
                <div className="bg-rose-950/30 p-3 rounded-xl border border-rose-500/20">
                  <h4 className="text-xs font-bold text-rose-400 uppercase tracking-wider mb-1 flex items-center gap-1.5">
                    <ShieldAlert size={14} /> Mismatch Details
                  </h4>
                  <ul className="text-xs text-rose-300 space-y-1 list-disc pl-4">
                    {result.mismatches.map((m, idx) => (
                      <li key={idx}>{m}</li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Footer info */}
        <div className="pt-4 border-t border-white/5 text-[11px] text-gray-500 flex items-center justify-between">
          <span>Powered by Go Connection Pools</span>
          <span className="font-mono">JSON Endpoint: /api/verify</span>
        </div>
      </div>
    </div>
  );
}

function CheckItem({ label, matched }) {
  return (
    <div
      className={`p-2.5 rounded-lg border text-xs font-medium flex items-center justify-between ${
        matched
          ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-400'
          : 'bg-rose-500/10 border-rose-500/20 text-rose-400'
      }`}
    >
      <span>{label}</span>
      {matched ? <CheckCircle2 size={16} /> : <XCircle size={16} />}
    </div>
  );
}
