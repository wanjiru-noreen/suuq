import { api } from "./api.js";
import { auth } from "./auth.js";

document.addEventListener("DOMContentLoaded", () => {
  initMenu();
  loadDashboardData();
});

function initMenu() {
  const moreBtn = document.getElementById("more-menu-btn");
  const moreModal = document.getElementById("more-modal");
  const logoutBtn = document.getElementById("logout-btn");

  if (moreBtn && moreModal) {
    moreBtn.addEventListener("click", () => {
      moreModal.classList.add("active");
    });

    moreModal.addEventListener("click", (e) => {
      if (
        e.target === moreModal ||
        e.target.classList.contains("sheet-handle")
      ) {
        moreModal.classList.remove("active");
      }
    });
  }

  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      auth.logout();
    });
  }
}

async function loadDashboardData() {
  try {
    const data = await api.get("/dashboard/summary");

    if (data.userName) {
      document.getElementById("user-greeting").textContent =
        `Good morning, ${data.userName} 👋`;
    }

    document.getElementById("stat-products").textContent =
      data.totalProducts ?? 0;
    document.getElementById("stat-stock").textContent =
      data.totalStockItems ?? 0;

    document.getElementById("stock-in").textContent =
      data.stockBreakdown?.inStock ?? 0;
    document.getElementById("stock-low").textContent =
      data.stockBreakdown?.low ?? 0;
    document.getElementById("stock-critical").textContent =
      data.stockBreakdown?.critical ?? 0;
    document.getElementById("stock-out").textContent =
      data.stockBreakdown?.out ?? 0;

    document.getElementById("stat-debtors").textContent =
      `KSh ${(data.totalDebtorsOutstanding ?? 0).toLocaleString()}`;
    document.getElementById("stat-creditors").textContent =
      `KSh ${(data.totalCreditorsOutstanding ?? 0).toLocaleString()}`;

    renderLowStockAlerts(data.lowStockItems || []);
  } catch (error) {
    console.error("Failed to load dashboard metrics:", error);
    const alertContainer = document.getElementById("low-stock-list");
    if (alertContainer) {
      alertContainer.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem;">Unable to load alerts.</div>`;
    }
  }
}

function renderLowStockAlerts(items) {
  const container = document.getElementById("low-stock-list");
  if (!container) return;

  if (items.length === 0) {
    container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem;">No critical or low stock items.</div>`;
    return;
  }

  container.innerHTML = items
    .map((item) => {
      let badgeClass = "badge-low";
      let badgeText = "🟡 Low";

      if (item.quantity === 0) {
        badgeClass = "badge-out";
        badgeText = "⚫ Out";
      } else if (item.status === "critical") {
        badgeClass = "badge-critical";
        badgeText = "🔴 Critical";
      }

      return `
      <div style="display: flex; justify-content: space-between; align-items: center; padding: 0.5rem 0; border-bottom: 1px solid var(--border-color);">
        <div>
          <div style="font-weight: 600; color: var(--text-primary); font-size: 0.9rem;">${item.name}</div>
          <div style="font-size: 0.75rem; color: var(--text-secondary);">Qty: ${item.quantity}</div>
        </div>
        <span class="badge ${badgeClass}">${badgeText}</span>
      </div>
    `;
    })
    .join("");
}
