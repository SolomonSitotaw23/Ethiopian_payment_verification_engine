import React, { useState } from 'react';
import { Terminal, Copy, Check } from 'lucide-react';

export default function CodeGenerator() {
  const [lang, setLang] = useState('curl');
  const [copied, setCopied] = useState(false);

  const snippets = {
    curl: `curl -X POST http://localhost:5000/api/verify \\
  -H "Content-Type: application/json" \\
  -d '{
    "receipt": "CJP9OSP9WZ",
    "provider": "telebirr",
    "expected": {
      "amount": 230,
      "recipientAccount": "0912345678",
      "recipientName": "Abrham Yalew"
    },
    "defaultVerification": true
  }'`,

    javascript: `// Node.js / JavaScript Fetch
async function verifyPayment() {
  const response = await fetch('http://localhost:5000/api/verify', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      receipt: 'CJP9OSP9WZ',
      provider: 'telebirr',
      expected: {
        amount: 230,
        recipientAccount: '0912345678',
        recipientName: 'Abrham Yalew'
      },
      defaultVerification: true
    })
  });
  const data = await response.json();
  console.log(data);
}

verifyPayment();`,

    go: `package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	payload := map[string]interface{}{
		"receipt": "CJP9OSP9WZ",
		"provider": "telebirr",
		"expected": map[string]interface{}{
			"amount": 230,
			"recipientAccount": "0912345678",
			"recipientName": "Abrham Yalew",
		},
		"defaultVerification": true,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:5000/api/verify", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBytes))
}`,

    python: `import requests

url = "http://localhost:5000/api/verify"
payload = {
    "receipt": "CJP9OSP9WZ",
    "provider": "telebirr",
    "expected": {
        "amount": 230,
        "recipientAccount": "0912345678",
        "recipientName": "Abrham Yalew"
    },
    "defaultVerification": True
}

response = requests.post(url, json=payload)
print(response.json())`,

    php: `<?php
$curl = curl_init();

curl_setopt_array($curl, array(
  CURLOPT_URL => 'http://localhost:5000/api/verify',
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_CUSTOMREQUEST => 'POST',
  CURLOPT_POSTFIELDS => json_encode([
    "receipt" => "CJP9OSP9WZ",
    "provider" => "telebirr",
    "expected" => [
      "amount" => 230,
      "recipientAccount" => "0912345678",
      "recipientName" => "Abrham Yalew"
    ],
    "defaultVerification" => true
  ]),
  CURLOPT_HTTPHEADER => array(
    'Content-Type: application/json'
  ),
));

$response = curl_exec($curl);
curl_close($curl);
echo $response;
?>`,
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(snippets[lang]);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="glass-panel p-6 mb-12">
      <div className="flex items-center justify-between mb-6 pb-4 border-b border-white/5">
        <div>
          <h2 className="text-lg font-bold text-white flex items-center gap-2">
            <Terminal className="text-indigo-400" size={20} />
            API Integration Code Generator
          </h2>
          <p className="text-xs text-gray-400">
            Copy production-ready integration code snippets for your preferred language or backend stack.
          </p>
        </div>

        <button
          onClick={handleCopy}
          className="btn-secondary text-xs"
        >
          {copied ? (
            <>
              <Check size={14} className="text-emerald-400" />
              Copied to Clipboard!
            </>
          ) : (
            <>
              <Copy size={14} />
              Copy Snippet
            </>
          )}
        </button>
      </div>

      {/* Language Selector Tabs */}
      <div className="flex items-center gap-2 mb-4 overflow-x-auto pb-2">
        {['curl', 'javascript', 'go', 'python', 'php'].map((l) => (
          <button
            key={l}
            onClick={() => setLang(l)}
            className={`px-3 py-1.5 rounded-lg text-xs font-mono font-semibold uppercase transition-all ${
              lang === l
                ? 'bg-indigo-600 text-white shadow-md'
                : 'bg-white/5 text-gray-400 hover:text-white'
            }`}
          >
            {l}
          </button>
        ))}
      </div>

      {/* Code Block */}
      <div className="code-block relative">
        <pre>{snippets[lang]}</pre>
      </div>
    </div>
  );
}
