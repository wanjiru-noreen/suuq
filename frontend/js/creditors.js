import { api } from "./api.js";
import { auth } from "./auth.js";

let creditorsList = [];

document.addEventListener("DOMContentLoaded", () => {
  initMenu();
  initCreditorModal();
  loadCreditors();

  const searchInput = document.getElementById("creditor-search-input");
  if (searchInput) {
    searchInput.addEventListener("input", (e) => {
      filterCreditors(e.target.value);
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

function initCreditorModal() {
  const openBtn = document.getElementById("open-creditor-modal");
  const closeBtn = document.getElementById("close-creditor-modal");
  const modal = document.getElementById("creditor-modal");
  const form = document.getElementById("creditor-form");

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
      const supplierName = document.getElementById("creditor-name").value;
      const amount =
        parseFloat(document.getElementById("creditor-amount").value) || 0;
      const notes = document.getElementById("creditor-notes").value;

      try {
        await api.post("/creditors", { supplierName, amount, notes });
        modal.classList.remove("active");
        form.reset();
        loadCreditors();
      } catch (err) {
        alert(err.message || "Failed to record creditor");
      }
    });
  }
}

async function loadCreditors() {
  try {
    creditorsList = await api.get("/creditors");
    renderCreditors(creditorsList);
  } catch (error) {
    console.error("Failed to load creditors:", error);
    const container = document.getElementById("creditors-list");
    if (container) {
      container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem; text-align: center;">Failed to load creditors ledger.</div>`;
    }
  }
}

function filterCreditors(query) {
  const q = query.toLowerCase();
  const filtered = creditorsList.filter((c) =>
    c.supplierName.toLowerCase().includes(q),
  );
  renderCreditors(filtered);
}

function renderCreditors(items) {
  const container = document.getElementById("creditors-list");
  if (!container) return;

  if (items.length === 0) {
    container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem; text-align: center; padding: 2rem 0;">No creditor records found.</div>`;
    return;
  }

  container.innerHTML = items
    .map(
      (item) => `
    <div class="card" style="margin-bottom: 0; display: flex; justify-content: space-between; align-items: center;">
      <div>
        <div style="font-weight: 600; color: var(--text-primary); font-size: 0.95rem;">${item.supplierName}</div>
        <div style="font-size: 0.8rem; color: var(--text-secondary); margin-top: 2px;">Owed: <span class="font-mono" style="font-weight: 600;">KSh ${item.amount.toLocaleString()}</span></div>
      </div>
      <span class="badge ${item.amount > 0 ? "badge-low" : "badge-in-stock"}">
        ${item.amount > 0 ? "Pending" : "Settled"}
      </span>
    </div>
  `,
    )
    .join("");
}
