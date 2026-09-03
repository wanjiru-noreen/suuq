import { api } from './api.js';
import { auth } from './auth.js';

let stockList = [];

document.addEventListener('DOMContentLoaded', () => {
  initMenu();
  initRestockModal();
  loadStock();

  const searchInput = document.getElementById('stock-search-input');
  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      filterStock(e.target.value);
    });
  }
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

function initRestockModal() {
  const openBtn = document.getElementById('open-restock-modal');
  const closeBtn = document.getElementById('close-restock-modal');
  const modal = document.getElementById('restock-modal');
  const form = document.getElementById('restock-form');

  if (openBtn && modal) {
    openBtn.addEventListener('click', async () => {
      await loadProductDropdown();
      modal.classList.add('active');
    });
  }

  if (closeBtn && modal) {
    closeBtn.addEventListener('click', () => modal.classList.remove('active'));
  }

  if (modal) {
    modal.addEventListener('click', (e) => {
      if (e.target === modal || e.target.classList.contains('sheet-handle')) {
        modal.classList.remove('active');
      }
    });
  }

  if (form) {
    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      const productId = document.getElementById('restock-product').value;
      const quantity = parseInt(document.getElementById('restock-qty').value, 10);

      try {
        await api.put(`/stock/${productId}`, { quantity });
        modal.classList.remove('active');
        form.reset();
        loadStock();
      } catch (err) {
        alert(err.message || 'Failed to update stock');
      }
    });
  }
}

async function loadProductDropdown() {
  const select = document.getElementById('restock-product');
  if (!select) return;

  try {
    const products = await api.get('/products');
    select.innerHTML = '<option value="">-- Choose product --</option>' + 
      products.map(p => `<option value="${p.id}">${p.name} (Current: ${p.quantity})</option>`).join('');
  } catch (error) {
    console.error('Failed to load products for dropdown:', error);
  }
}

async function loadStock() {
  try {
    stockList = await api.get('/stock');
    renderStock(stockList);
  } catch (error) {
    console.error('Failed to load stock:', error);
    const container = document.getElementById('stock-list');
    if (container) {
      container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem; text-align: center;">Failed to load stock inventory.</div>`;
    }
  }
}

function filterStock(query) {
  const q = query.toLowerCase();
  const filtered = stockList.filter(s => s.name.toLowerCase().includes(q));
  renderStock(filtered);
}

function renderStock(items) {
  const container = document.getElementById('stock-list');
  if (!container) return;

  if (items.length === 0) {
    container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem; text-align: center; padding: 2rem 0;">No stock items found.</div>`;
    return;
  }

  container.innerHTML = items.map(item => {
    let badgeClass = 'badge-in-stock';
    let badgeText = '🟢 In Stock';

    if (item.quantity === 0) {
      badgeClass = 'badge-out';
      badgeText = '⚫ Out of Stock';
    } else if (item.quantity <= 5) {
      badgeClass = 'badge-critical';
      badgeText = '🔴 Critical';
    } else if (item.quantity <= 15) {
      badgeClass = 'badge-low';
      badgeText = '🟡 Low';
    }

    return `
      <div class="card" style="margin-bottom: 0; display: flex; justify-content: space-between; align-items: center;">
        <div>
          <div style="font-weight: 600; color: var(--text-primary); font-size: 0.95rem;">${item.name}</div>
          <div style="font-size: 0.8rem; color: var(--text-secondary); margin-top: 2px;">Quantity: <span class="font-mono" style="font-weight: 600;">${item.quantity} units</span></div>
        </div>
        <span class="badge ${badgeClass}">${badgeText}</span>
      </div>
    `;
  }).join('');
}