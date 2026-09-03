import { api } from './api.js';
import { auth } from './auth.js';

document.addEventListener('DOMContentLoaded', () => {
  initMenu();
  initAIButton();
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

function initAIButton() {
  const runBtn = document.getElementById('run-ai-btn');
  const resultsCard = document.getElementById('ai-results-card');
  const outputDiv = document.getElementById('ai-output');

  if (runBtn && resultsCard && outputDiv) {
    runBtn.addEventListener('click', async () => {
      runBtn.disabled = true;
      runBtn.textContent = 'Analyzing patterns...';
      resultsCard.style.display = 'block';
      outputDiv.textContent = 'Running diagnostic pattern analysis on operational metrics...';

      try {
        const response = await api.post('/ai/analyze', {});
        outputDiv.textContent = response.analysis || 'Analysis complete. No operational anomalies detected.';
      } catch (err) {
        outputDiv.textContent = err.message || 'Failed to complete AI diagnostic analysis.';
      } finally {
        runBtn.disabled = false;
        runBtn.textContent = 'Run Diagnostics';
      }
    });
  }
}