import { api } from "./api.js";
import { auth } from "./auth.js";

let productsList = [];

document.addEventListener("DOMContentLoaded", () => {
  initMenu();
  initProductModal();
  loadProducts();

  const searchInput = document.getElementById("search-input");
  if (searchInput) {
    searchInput.addEventListener("input", (e) => {
      filterProducts(e.target.value);
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

function initProductModal() {
  const openBtn = document.getElementById("open-add-modal");
  const closeBtn = document.getElementById("close-modal");
  const modal = document.getElementById("product-modal");
  const form = document.getElementById("product-form");

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
      const name = document.getElementById("prod-name").value;
      const category = document.getElementById("prod-category").value;
      const description = document.getElementById("prod-desc").value;
      const quantity =
        parseInt(document.getElementById("prod-qty").value, 10) || 0;

      try {
        await api.post("/products", { name, category, description, quantity });
        modal.classList.remove("active");
        form.reset();
        loadProducts();
      } catch (err) {
        alert(err.message || "Failed to save product");
      }
    });
  }
}

async function loadProducts() {
  try {
    productsList = await api.get("/products");
    renderProducts(productsList);
  } catch (error) {
    console.error("Failed to load products:", error);
    const container = document.getElementById("products-list");
    if (container) {
      container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem; text-align: center;">Failed to load catalog.</div>`;
    }
  }
}

function filterProducts(query) {
  const q = query.toLowerCase();
  const filtered = productsList.filter(
    (p) =>
      p.name.toLowerCase().includes(q) || p.category.toLowerCase().includes(q),
  );
  renderProducts(filtered);
}

function renderProducts(items) {
  const container = document.getElementById("products-list");
  if (!container) return;

  if (items.length === 0) {
    container.innerHTML = `<div style="color: var(--text-secondary); font-size: 0.85rem; text-align: center; padding: 2rem 0;">No products found.</div>`;
    return;
  }

  container.innerHTML = items
    .map(
      (item) => `
    <div class="card" style="margin-bottom: 0; display: flex; justify-content: space-between; align-items: center;">
      <div>
        <div style="font-weight: 600; color: var(--text-primary); font-size: 0.95rem;">${item.name}</div>
        <div style="font-size: 0.8rem; color: var(--text-secondary); margin-top: 2px;">${item.category} • Stock: <span class="font-mono" style="font-weight: 600;">${item.quantity}</span></div>
      </div>
      <div>
        <span class="badge ${item.quantity > 10 ? "badge-in-stock" : item.quantity > 0 ? "badge-low" : "badge-out"}">
          ${item.quantity > 0 ? `${item.quantity} units` : "Out"}
        </span>
      </div>
    </div>
  `,
    )
    .join("");
}
