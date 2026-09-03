import { api } from "./api.js";
import { auth } from "./auth.js";

let debtorsList = [];

document.addEventListener("DOMContentLoaded", () => {
  initMenu();
  initDebtModal();
  loadDebtors();

  const searchInput = document.getElementById("debtor-search-input");
  if (searchInput) {
    searchInput.addEventListener("input", (e) => {
      filterDebtors(e.target.value);
    });
  }
});

function initMenu() {
  const moreBtn = document.getElementById("more-menu-btn");
  const moreModal = document.getElementById("more-modal");
  const logoutBtn = document.getElementById("logout-btn");

  if (moreBtn && moreModal) {
    moreBtn.addEventListener("click", () => moreModal.classList.add("active"));
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
    logoutBtn.addEventListener("click", () => auth.logout());
  }
}

function initDebtModal() {
  const openBtn = document.getElementById("open-debt-modal");
  const closeBtn = document.getElementById("close-debt-modal");
  const modal = document.getElementById("debt-modal");
  const form = document.getElementById("debt-form");

  if (openBtn && modal) {
    openBtn.addEventListener("click", () => modal.classList.add("active"));
  }

  if (closeBtn && modal) {
    closeBtn.addEventListener("click", () => modal.classList.remove("active"));
  }

  if (modal) {
    modal.addEventListener("click", (e) => {
      if (e.target === modal || e.target.classList.contains("sheet-handle")) {
        modal.classList.remove("active");
      }
    });
  }

  if (form) {
    form.addEventListener("submit", async (e) => {
      e.preventDefault();
      const customerName = document.getElementById("debtor-name").value;
      const amount =
        parseFloat(document.getElementById("debt-amount").value) || 0;
      const notes = document.getElementById("debt-notes").value;

      try {
        await api.post("/debtors", { customerName, amount, notes });
        modal.classList.remove("active");
        form.reset();
        loadDebtors();
      } catch (err) {
        alert(err.message || "Failed to record debt");
      }
    });
  }
}

async function loadDebtors() {
  try {
    debtorsList = await api.get("/debtors");
    renderDebtors(debtorsList);
  } catch (error) {
    console.error("Failed to load debtors:", error);
    const container = document.getElementById("debtors-list");
    if (container) {
      container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem; text-align: center;">Failed to load debtors ledger.</div>`;
    }
  }
}

function filterDebtors(query) {
  const q = query.toLowerCase();
  const filtered = debtorsList.filter((d) =>
    d.customerName.toLowerCase().includes(q),
  );
  renderDebtors(filtered);
}

function renderDebtors(items) {
  const container = document.getElementById("debtors-list");
  if (!container) return;

  if (items.length === 0) {
    container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem; text-align: center; padding: 2rem 0;">No debtor records found.</div>`;
    return;
  }

  container.innerHTML = items
    .map(
      (item) => `
    <div class="card" style="margin-bottom: 0; display: flex; justify-content: space-between; align-items: center;">
      <div>
        <div style="font-weight: 600; color: var(--text-primary); font-size: 0.95rem;">${item.customerName}</div>
        <div style="font-size: 0.8rem; color: var(--text-secondary); margin-top: 2px;">Due: <span class="font-mono" style="font-weight: 600;">KSh ${item.amount.toLocaleString()}</span></div>
      </div>
      <span class="badge ${item.amount > 0 ? "badge-critical" : "badge-in-stock"}">
        ${item.amount > 0 ? "Unpaid" : "Settled"}
      </span>
    </div>
  `,
    )
    .join("");
}
