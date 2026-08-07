import React from 'react';
import { Smartphone, Landmark, Building2, CreditCard, ExternalLink } from 'lucide-react';

export default function ProviderCards({ onSelectSample }) {
  const providers = [
    {
      id: 'telebirr',
      name: 'Telebirr',
      tagline: 'Ethio Telecom Wallet',
      color: '#ff8c00',
      bgColor: 'rgba(255, 140, 0, 0.1)',
      borderColor: 'rgba(255, 140, 0, 0.3)',
      icon: Smartphone,
      accountType: 'Phone Number (09... / +251...)',
      sampleId: 'CJP9OSP9WZ',
      samplePhone: '0912345678',
      sampleName: 'Abrham Yalew',
    },
    {
      id: 'cbe',
      name: 'CBE Bank',
      tagline: 'Commercial Bank of Ethiopia',
      color: '#8b5cf6',
      bgColor: 'rgba(139, 92, 246, 0.1)',
      borderColor: 'rgba(139, 92, 246, 0.3)',
      icon: Landmark,
      accountType: '13-Digit Account (1000...)',
      sampleId: 'FT25292FRPWD89873710',
      samplePhone: '1****1234',
      sampleName: 'ABRHAM YALEW',
    },
    {
      id: 'boa',
      name: 'Abyssinia (BOA)',
      tagline: 'Bank of Abyssinia',
      color: '#06b6d4',
      bgColor: 'rgba(6, 182, 212, 0.1)',
      borderColor: 'rgba(6, 182, 212, 0.3)',
      icon: Building2,
      accountType: 'Abyssinia Account / Slip ID',
      sampleId: 'FT25284X11PS79328',
      samplePhone: '1******95',
      sampleName: 'TEWODROS HULGIZIE TEMESGEN',
    },
    {
      id: 'amharabank',
      name: 'Amhara Bank',
      tagline: 'Amhara Bank SC',
      color: '#10b981',
      bgColor: 'rgba(16, 185, 129, 0.1)',
      borderColor: 'rgba(16, 185, 129, 0.3)',
      icon: CreditCard,
      accountType: 'Amhara Account ID (ETB...)',
      sampleId: 'AB1234567890',
      samplePhone: 'ETB1251800010003',
      sampleName: 'Abrham Yalew',
    },
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      {providers.map((p) => {
        const IconComponent = p.icon;
        return (
          <div
            key={p.id}
            className="glass-panel glass-panel-interactive p-5 flex flex-col justify-between cursor-pointer"
            style={{
              borderColor: p.borderColor,
            }}
            onClick={() => onSelectSample(p)}
          >
            <div>
              <div className="flex items-center justify-between mb-3">
                <div
                  className="p-3 rounded-xl"
                  style={{ background: p.bgColor, color: p.color }}
                >
                  <IconComponent size={24} />
                </div>
                <span className="text-[10px] uppercase font-bold tracking-wider px-2 py-0.5 rounded-full border border-white/10 text-gray-400">
                  Supported
                </span>
              </div>
              <h3 className="text-lg font-bold text-white mb-1" style={{ color: p.color }}>
                {p.name}
              </h3>
              <p className="text-xs text-gray-400 mb-3">{p.tagline}</p>
              <div className="text-[11px] text-gray-400 bg-slate-900/60 p-2 rounded-lg border border-white/5 font-mono mb-4">
                <span className="text-gray-500">Account:</span> {p.accountType}
              </div>
            </div>

            <button
              type="button"
              className="w-full text-xs font-semibold py-2 px-3 rounded-lg bg-white/5 hover:bg-white/10 text-gray-300 flex items-center justify-center gap-1.5 transition-colors"
            >
              <span>Try Sample Receipt</span>
              <ExternalLink size={12} />
            </button>
          </div>
        );
      })}
    </div>
  );
}
