import React, { useState } from 'react';
import Navbar from './components/Navbar';
import ProviderCards from './components/ProviderCards';
import SingleVerifier from './components/SingleVerifier';
import BatchVerifier from './components/BatchVerifier';
import CodeGenerator from './components/CodeGenerator';
import SystemMetrics from './components/SystemMetrics';

export default function App() {
  const [activeTab, setActiveTab] = useState('single');
  const [selectedSample, setSelectedSample] = useState(null);

  const handleSelectSample = (sample) => {
    setSelectedSample(sample);
    setActiveTab('single');
  };

  return (
    <div className="min-h-screen pb-16">
      {/* Top Navbar */}
      <Navbar activeTab={activeTab} setActiveTab={setActiveTab} />

      <main className="max-w-7xl mx-auto px-6">
        {/* Interactive Provider Cards Header */}
        <ProviderCards onSelectSample={handleSelectSample} />

        {/* Tab Content */}
        {activeTab === 'single' && <SingleVerifier sampleData={selectedSample} />}
        {activeTab === 'batch' && <BatchVerifier />}
        {activeTab === 'code' && <CodeGenerator />}
        {activeTab === 'metrics' && <SystemMetrics />}
      </main>

      {/* Footer */}
      <footer className="max-w-7xl mx-auto px-6 mt-16 pt-8 border-t border-white/5 flex flex-col sm:flex-row items-center justify-between text-xs text-gray-500 gap-4">
        <div>
          <span>© 2026 Ethiopian Payment Receipt Verifier. Open Source MIT License.</span>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-gray-400">Supported: Telebirr • CBE • BOA • Amhara Bank</span>
        </div>
      </footer>
    </div>
  );
}
