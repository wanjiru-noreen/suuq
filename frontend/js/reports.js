import { api } from './api.js';
import { auth } from './auth.js';

document.addEventListener('DOMContentLoaded', () => {
  initMenu();
  loadReportData();
});

function initMenu() {
  const moreBtn = document.getElementById('more-menu-btn');
  const moreModal = document.getElementById('more-modal');
  const logoutBtn = document.getElementById('logout-btn');

  if (moreBtn && moreModal) {
    moreBtn.addEventListener('click', () => moreModal.classList.add('active'));
    moreModal.addEventListener('click', (e) => {
      if (e.target === moreModal || e.target.classList.contains('sheet-handle')) {
        moreModal.classList.remove('active');
      }
    });
  }

  if (logoutBtn) {
    logoutBtn.addEventListener('click', () => auth.logout());
  }
}

async function loadReportData() {
  try {
    const data = await api.get('/dashboard/summary');
    
    document.getElementById('rep-products').textContent = data.totalProducts ?? 0;
    document.getElementById('rep-stock').textContent = data.totalStockItems ?? 0;
    document.getElementById('rep-debtors').textContent = `KSh ${(data.totalDebtorsOutstanding ?? 0).toLocaleString()}`;
    document.getElementById('rep-creditors').textContent = `KSh ${(data.totalCreditorsOutstanding ?? 0).toLocaleString()}`;

    document.getElementById('health-in').textContent = data.stockBreakdown?.inStock ?? 0;
    document.getElementById('health-low').textContent = data.stockBreakdown?.low ?? 0;
    document.getElementById('health-critical').textContent = data.stockBreakdown?.critical ?? 0;
    document.getElementById('health-out').textContent = data.stockBreakdown?.out ?? 0;
  } catch (error) {
    console.error('Failed to load report metrics:', error);
  }
}