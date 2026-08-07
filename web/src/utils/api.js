export async function verifySingleReceipt(payload) {
  const response = await fetch('/api/verify', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });

  const data = await response.json();
  if (!response.ok && !data.details) {
    throw new Error(data.error || 'Failed to verify receipt');
  }
  return data;
}

export async function verifyBatchReceipts(payload) {
  const response = await fetch('/api/verify/batch', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });

  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || 'Failed to process batch verification');
  }
  return data;
}

export async function getHealth() {
  try {
    const response = await fetch('/health');
    if (!response.ok) return { status: 'DOWN' };
    return await response.json();
  } catch {
    return { status: 'OFFLINE' };
  }
}

export async function getMetrics() {
  try {
    const response = await fetch('/metrics');
    if (!response.ok) return null;
    return await response.json();
  } catch {
    return null;
  }
}
