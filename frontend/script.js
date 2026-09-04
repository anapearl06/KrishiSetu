// ============================================================
// KRISHISETU - BACKEND API CONNECTED SCRIPT
// Backend Render URL: https://krishisetu-api-tiau.onrender.com
// Neon DB Endpoint Connected
// ============================================================

const API_BASE_URL = "https://krishisetu-api-tiau.onrender.com";

function getTokenRole() {
  const token = localStorage.getItem("token");
  if (!token) return null;
  try {
    const payload = token.split(".")[1];
    const decoded = JSON.parse(atob(payload.replace(/-/g, "+").replace(/_/g, "/")));
    return decoded.role || null;
  } catch (err) {
    return null;
  }
}

function extractErrorMessage(data, fallback) {
  if (!data) return fallback;
  if (typeof data.error === "object" && data.error.message) {
    return data.error.message;
  }
  if (typeof data.error === "string" && data.error) return data.error;
  if (typeof data.message === "string" && data.message) return data.message;
  if (typeof data.detail === "string" && data.detail) return data.detail;
  return fallback;
}

// ============================================================
// ANIMATED FARM SCENERY (login / register backgrounds)
// ============================================================
function buildFarmland() {
  const container = document.getElementById("farmScenery");
  if (!container) return;

  let html = '<div class="sun"></div>';
  html += '<div class="farm-hill-back"></div>';
  html += '<div class="farm-field"></div>';

  // Back row (small, higher) + front row (bigger, lower)
  const rows = [
    { top: 12, count: 10, size: 15 },
    { top: 32, count: 7, size: 23 },
  ];

  rows.forEach(function (row) {
    for (let i = 0; i < row.count; i++) {
      const frac = (i + 0.5) / row.count;
      const x = frac * 100 + (Math.random() * 6 - 3);
      const delay = (Math.random() * 2.5).toFixed(2);
      html +=
        '<div class="crop" style="left:' +
        x.toFixed(1) +
        "%;top:" +
        row.top +
        "%;font-size:" +
        row.size +
        "px;--delay:" +
        delay +
        's">';
      html +=
        '<div class="stalk"></div><div class="leaf l"></div><div class="leaf r"></div><div class="head"></div>';
      html += "</div>";
    }
  });

  // Small floating sparkles
  for (let s = 0; s < 4; s++) {
    const sx = 10 + Math.random() * 80;
    const sy = 6 + Math.random() * 26;
    const sd = (Math.random() * 3).toFixed(2);
    const size = 4 + Math.random() * 6;
    html +=
      '<div class="sparkle" style="left:' +
      sx.toFixed(1) +
      "%;top:" +
      sy.toFixed(1) +
      "%;width:" +
      size.toFixed(1) +
      "px;height:" +
      size.toFixed(1) +
      "px;--delay:" +
      sd +
      's"></div>';
  }

  container.innerHTML = html;
}

// ============================================================
// SHOW / HIDE PASSWORD TOGGLES
// ============================================================
function setupPasswordToggles() {
  document.querySelectorAll(".password-toggle").forEach(function (btn) {
    btn.addEventListener("click", function () {
      const id = this.getAttribute("data-target");
      const input = document.getElementById(id);
      if (!input) return;
      const wasHidden = input.type === "password";
      input.type = wasHidden ? "text" : "password";
      this.textContent = wasHidden ? "🙈" : "👁️";
    });
  });
}

// ============================================================
// AUTH PAGES (LOGIN + REGISTER)
// ============================================================
document.addEventListener("DOMContentLoaded", function () {
  // Build the animated farm scene (login / register pages only)
  buildFarmland();

  // Wire up password eye toggles
  setupPasswordToggles();

  // HELPER: Fetch role from URL query params
  function getRoleFromURL() {
    const params = new URLSearchParams(window.location.search);
    return params.get("role");
  }

  // ============================================================
  // LOGIN PAGE LOGIC
  // ============================================================
  const loginFarmerTab = document.getElementById("loginFarmerTab");
  const loginBuyerTab = document.getElementById("loginBuyerTab");
  const farmerLoginForm = document.getElementById("farmerLoginForm");
  const buyerLoginForm = document.getElementById("buyerLoginForm");
  const loginTitle = document.getElementById("loginTitle");
  const loginRegisterLink = document.getElementById("loginRegisterLink");

  function activateFarmerLogin() {
    if (!farmerLoginForm || !buyerLoginForm) return;
    farmerLoginForm.classList.remove("hidden");
    buyerLoginForm.classList.add("hidden");
    if (loginFarmerTab) loginFarmerTab.classList.add("active");
    if (loginBuyerTab) loginBuyerTab.classList.remove("active");
    if (loginTitle) loginTitle.textContent = "Farmer Login";
    if (loginRegisterLink)
      loginRegisterLink.href = "./register.html?role=farmer";
  }

  function activateBuyerLogin() {
    if (!farmerLoginForm || !buyerLoginForm) return;
    buyerLoginForm.classList.remove("hidden");
    farmerLoginForm.classList.add("hidden");
    if (loginBuyerTab) loginBuyerTab.classList.add("active");
    if (loginFarmerTab) loginFarmerTab.classList.remove("active");
    if (loginTitle) loginTitle.textContent = "Buyer Login";
    if (loginRegisterLink)
      loginRegisterLink.href = "./register.html?role=buyer";
  }

  if (loginFarmerTab)
    loginFarmerTab.addEventListener("click", activateFarmerLogin);
  if (loginBuyerTab)
    loginBuyerTab.addEventListener("click", activateBuyerLogin);

  if (loginFarmerTab && loginBuyerTab) {
    const role = getRoleFromURL();
    if (role === "buyer") activateBuyerLogin();
    else activateFarmerLogin();
  }

  // FARMER LOGIN BACKEND CALL
  if (farmerLoginForm) {
    farmerLoginForm.addEventListener("submit", async function (event) {
      event.preventDefault();
      const phone =
        document.getElementById("farmerLoginPhone")?.value ||
        document.getElementById("farmerPhone")?.value;
      const password =
        document.getElementById("farmerLoginPassword")?.value ||
        document.getElementById("farmerPassword")?.value;

      try {
        const response = await fetch(`${API_BASE_URL}/api/v1/farmers/login`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ phone, password }),
        });

        const data = await response.json();
        if (response.ok) {
          alert("Farmer Login Successful!");
          const token = data.data?.token || data.token;
          if (token) localStorage.setItem("token", token);
          window.location.href = "./farmer-dashboard.html";
        } else {
          alert(
            extractErrorMessage(data, "Login failed! Check credentials."),
          );
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error! Please check your connection and try again.");
      }
    });
  }

  // BUYER LOGIN BACKEND CALL
  if (buyerLoginForm) {
    buyerLoginForm.addEventListener("submit", async function (event) {
      event.preventDefault();
      const phone =
        document.getElementById("buyerLoginPhone")?.value ||
        document.getElementById("buyerPhone")?.value;
      const password =
        document.getElementById("buyerLoginPassword")?.value ||
        document.getElementById("buyerPassword")?.value;

      try {
        const response = await fetch(`${API_BASE_URL}/api/v1/buyers/login`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ phone, password }),
        });

        const data = await response.json();
        if (response.ok) {
          alert("Buyer Login Successful!");
          const token = data.data?.token || data.token;
          if (token) localStorage.setItem("token", token);
          window.location.href = "./buyer-dashboard.html";
        } else {
          alert(
            extractErrorMessage(data, "Login failed! Check credentials."),
          );
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error! Please check your connection and try again.");
      }
    });
  }

  // ============================================================
  // REGISTER PAGE LOGIC (WITH REQUIRED DISTRICT FIELD)
  // ============================================================
  const registerFarmerTab = document.getElementById("registerFarmerTab");
  const registerBuyerTab = document.getElementById("registerBuyerTab");
  const farmerForm = document.getElementById("farmerForm");
  const buyerForm = document.getElementById("buyerForm");

  function activateFarmerRegister() {
    if (!farmerForm || !buyerForm) return;
    farmerForm.classList.remove("hidden");
    farmerForm.classList.add("flex");
    buyerForm.classList.add("hidden");
    buyerForm.classList.remove("flex");
    if (registerFarmerTab) registerFarmerTab.classList.add("active");
    if (registerBuyerTab) registerBuyerTab.classList.remove("active");
  }

  function activateBuyerRegister() {
    if (!farmerForm || !buyerForm) return;
    buyerForm.classList.remove("hidden");
    buyerForm.classList.add("flex");
    farmerForm.classList.add("hidden");
    farmerForm.classList.remove("flex");
    if (registerBuyerTab) registerBuyerTab.classList.add("active");
    if (registerFarmerTab) registerFarmerTab.classList.remove("active");
  }

  if (registerFarmerTab)
    registerFarmerTab.addEventListener("click", activateFarmerRegister);
  if (registerBuyerTab)
    registerBuyerTab.addEventListener("click", activateBuyerRegister);

  if (registerFarmerTab && registerBuyerTab) {
    const role = getRoleFromURL();
    if (role === "buyer") activateBuyerRegister();
    else activateFarmerRegister();
  }

  // FARMER REGISTER BACKEND CALL
  if (farmerForm) {
    farmerForm.addEventListener("submit", async function (event) {
      event.preventDefault();
      const name = document.getElementById("farmerName")?.value || "";
      const phone = document.getElementById("farmerPhone")?.value || "";
      const district =
        document.getElementById("farmerDistrict")?.value || "Default District";
      const state =
        document.getElementById("farmerState")?.value || "Uttar Pradesh";
      const password = document.getElementById("farmerPassword")?.value || "";
      const confirmPassword =
        document.getElementById("farmerConfirmPassword")?.value || "";

      if (password !== confirmPassword) {
        alert("Passwords do not match!");
        return;
      }

      try {
        const response = await fetch(
          `${API_BASE_URL}/api/v1/farmers/register`,
          {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              name: name,
              phone: phone,
              district: district,
              state: state,
              password: password,
            }),
          },
        );

        const data = await response.json();
        if (response.ok) {
          alert("Farmer Registration Successful! Redirecting to Login...");
          window.location.href = "./login.html?role=farmer";
        } else {
          alert(
            extractErrorMessage(data, "Registration failed! Please try again."),
          );
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error! Please check your connection and try again.");
      }
    });
  }

  // BUYER REGISTER BACKEND CALL
  if (buyerForm) {
    buyerForm.addEventListener("submit", async function (event) {
      event.preventDefault();
      const name = document.getElementById("buyerName")?.value || "";
      const businessName =
        document.getElementById("businessName")?.value || name;
      const businessType =
        document.getElementById("businessType")?.value || "Retailer";
      const phone = document.getElementById("buyerPhone")?.value || "";
      const district =
        document.getElementById("buyerDistrict")?.value || "Default District";
      const state =
        document.getElementById("buyerState")?.value || "Uttar Pradesh";
      const password = document.getElementById("buyerPassword")?.value || "";
      const confirmPassword =
        document.getElementById("buyerConfirmPassword")?.value || "";

      if (password !== confirmPassword) {
        alert("Passwords do not match!");
        return;
      }

      try {
        const response = await fetch(`${API_BASE_URL}/api/v1/buyers/register`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name: name,
            business_name: businessName,
            business_type: businessType,
            phone: phone,
            district: district,
            state: state,
            password: password,
          }),
        });

        const data = await response.json();
        if (response.ok) {
          alert("Buyer Registration Successful! Redirecting to Login...");
          window.location.href = "./login.html?role=buyer";
        } else {
          alert(
            extractErrorMessage(data, "Registration failed! Please try again."),
          );
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error! Please check your connection and try again.");
      }
    });
  }
});

// Restrict phone input to numbers only
const phoneInputs = document.querySelectorAll('input[type="tel"]');
phoneInputs.forEach((input) => {
  input.addEventListener("input", function () {
    input.value = input.value.replace(/\D/g, "").slice(0, 10);
  });
});

// ============================================================
// F5 — CREATE CROP LISTING (POST /api/v1/listings)
// ============================================================
function openCreateListingDrawer() {
  const drawer = document.getElementById("createListingDrawer");
  const backdrop = document.getElementById("drawerBackdrop");
  if (backdrop) {
    backdrop.classList.remove("hidden");
    setTimeout(() => backdrop.classList.remove("opacity-0"), 10);
  }
  if (drawer) {
    drawer.classList.remove("translate-x-full");
  }
}

function closeCreateListingDrawer() {
  const drawer = document.getElementById("createListingDrawer");
  const backdrop = document.getElementById("drawerBackdrop");
  if (drawer) drawer.classList.add("translate-x-full");
  if (backdrop) {
    backdrop.classList.add("opacity-0");
    setTimeout(() => backdrop.classList.add("hidden"), 300);
  }
}

document
  .getElementById("drawerBackdrop")
  ?.addEventListener("click", closeCreateListingDrawer);

document
  .getElementById("createListingForm")
  ?.addEventListener("submit", async function (e) {
    e.preventDefault();

    const token = localStorage.getItem("token");
    if (!token) {
      alert("Session expired. Please login again.");
      window.location.href = "./login.html?role=farmer";
      return;
    }

    const listingPayload = {
      crop: document.getElementById("produceName")?.value.trim(),
      quantity: parseFloat(document.getElementById("produceQuantity")?.value),
      unit: document.getElementById("produceUnit")?.value,
      price: parseFloat(document.getElementById("producePrice")?.value),
      state: document.getElementById("produceLocation")?.value.trim(),
      district: document.getElementById("produceDistrict")?.value.trim(),
      description: document.getElementById("produceDesc")?.value || "",
    };

    try {
      const submitBtn = e.target.querySelector('button[type="submit"]');
      const originalText = submitBtn.textContent;
      submitBtn.textContent = "Publishing...";
      submitBtn.disabled = true;

      const response = await fetch(`${API_BASE_URL}/api/v1/listings`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(listingPayload),
      });

      const data = await response.json();
      submitBtn.textContent = originalText;
      submitBtn.disabled = false;

      if (response.ok) {
        alert("Crop listed successfully! 🌾");
        const drawer = document.getElementById("createListingDrawer");
        const backdrop = document.getElementById("drawerBackdrop");
        if (drawer) drawer.classList.add("translate-x-full");
        if (backdrop) {
          backdrop.classList.add("opacity-0");
          setTimeout(() => backdrop.classList.add("hidden"), 300);
        }
        e.target.reset();
        window.dispatchEvent(new Event("listingCreated"));
      } else {
        alert(extractErrorMessage(data, "Failed to create crop listing."));
      }
    } catch (err) {
      console.error("Error creating listing:", err);
      alert("Server error. Please check your connection and try again.");
    }
  });

// ============================================================
// MY PRODUCE (F6, F8, F9) — GET / EDIT / DELETE
// ============================================================
let currentListings = [];

async function renderMyProduce() {
  const produceGrid = document.getElementById("produceGrid");
  if (!produceGrid) return;

  const statusFilter =
    typeof document.getElementById("statusFilter")?.value === "string"
      ? document.getElementById("statusFilter").value
      : "";
  const search =
    typeof document.getElementById("searchInput")?.value === "string"
      ? document.getElementById("searchInput").value
      : "";

  const token = localStorage.getItem("token");
  if (!token) {
    produceGrid.innerHTML = `<p class="text-red-600 font-medium">Please login to view your produce listings.</p>`;
    return;
  }

  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/listings/my`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (response.ok && Array.isArray(data)) {
      currentListings = data;
      const total = data.length;
      const active = data.filter((l) => l.status === "ACTIVE").length;
      const pending = data.filter((l) => l.status === "PENDING").length;
      const sold = data.filter(
        (l) => l.status === "SOLD" || l.status === "COMPLETED",
      ).length;

      const elTotal = document.getElementById("statTotal");
      const elActive = document.getElementById("statActive");
      const elPending = document.getElementById("statPending");
      const elSold = document.getElementById("statSold");
      if (elTotal) elTotal.textContent = total;
      if (elActive) elActive.textContent = active;
      if (elPending) elPending.textContent = pending;
      if (elSold) elSold.textContent = sold;

      const status = (statusFilter || "").toUpperCase();
      const term = (search || "").trim().toLowerCase();

      const visible = data.filter((item) => {
        const matchStatus = !status || (item.status || "ACTIVE") === status;
        const matchSearch =
          !term || (item.crop || "").toLowerCase().includes(term);
        return matchStatus && matchSearch;
      });

      if (visible.length === 0) {
        produceGrid.innerHTML = `
          <div class="col-span-full flex flex-col items-center justify-center py-16 text-center">
            <div class="text-4xl mb-4">🌾</div>
            <h3 class="text-lg font-bold text-[#181D17] mb-1">${
              total === 0 ? "No produce listed yet" : "No results found"
            }</h3>
            <p class="text-sm text-[#40493D]">${
              total === 0
                ? 'Use the "+ Add Produce" button to post your first crop listing to the marketplace.'
                : "Try adjusting your status filter or search term."
            }</p>
          </div>`;
        return;
      }

      produceGrid.innerHTML = visible
        .map(
          (item, idx) => `
        <div class="glass-card rounded-2xl border border-[#E0E4DA]/60 p-5 shadow-sm hover:-translate-y-1.5 hover:shadow-lg ksetu-fade-up" style="animation-delay:${idx * 60}ms">
          <div class="flex justify-between items-start mb-2">
            <span class="status-badge ${
              item.status === "ACTIVE"
                ? "accepted"
                : item.status === "PENDING"
                  ? "pending"
                  : "delivered"
            }">${item.status || "ACTIVE"}</span>
            <span class="text-xs text-[#40493D]">📍 ${item.district || ""}, ${item.state || ""}</span>
          </div>
          <h3 class="text-lg font-bold text-[#181D17]">${item.crop}</h3>
          <p class="text-xs text-[#40493D] mt-1">${item.description || "No description provided"}</p>

          <div class="mt-4 pt-4 border-t border-[#F1F5EB] flex justify-between items-center">
            <div>
              <p class="text-xs text-[#40493D]">Quantity: <strong class="text-[#181D17]">${item.quantity} ${item.unit}</strong></p>
              <p class="text-lg font-bold brand-name">₹${item.price} / ${item.unit}</p>
            </div>
            <div class="flex gap-2">
              <button onclick="openEditModal('${item.id}', ${item.price}, ${item.quantity})" class="px-3 py-1.5 bg-amber-50 text-amber-700 border border-amber-200 text-xs font-semibold rounded-xl hover:bg-amber-100 transition-colors">✏️ Edit</button>
              ${
                item.status === "SOLD"
                  ? ""
                  : `<button onclick="openListingMatches('${item.id}')" class="px-3 py-1.5 bg-[#0D631B] text-white text-xs font-semibold rounded-xl hover:bg-[#2E7D32] transition-colors">🔍 Find Matches</button>`
              }
              <button onclick="deleteListing('${item.id}')" class="px-3 py-1.5 bg-red-50 text-red-600 border border-red-200 text-xs font-semibold rounded-xl hover:bg-red-100 transition-colors">🗑 Cancel</button>
            </div>
          </div>
        </div>
      `,
        )
        .join("");
    } else {
      produceGrid.innerHTML = `<p class="text-red-600">Failed to load produce: ${data.message || "Error occurred"}</p>`;
    }
  } catch (err) {
    console.error("Error fetching My Produce:", err);
    produceGrid.innerHTML = `<p class="text-red-600">Server error loading produce grid.</p>`;
  }
}

// Open / Close Edit Modal (F8)
function openEditModal(id, price, quantity) {
  document.getElementById("editId").value = id;
  document.getElementById("editPrice").value = price;
  document.getElementById("editQuantity").value = quantity;
  document.getElementById("editModal").classList.remove("hidden");
}

function closeEditModal() {
  document.getElementById("editModal").classList.add("hidden");
}

// Edit Form Submit (PUT /api/v1/listings/:id) (F8)
document
  .getElementById("editForm")
  ?.addEventListener("submit", async function (e) {
    e.preventDefault();
    const id = document.getElementById("editId").value;
    const token = localStorage.getItem("token");

    const payload = {
      price: parseFloat(document.getElementById("editPrice").value),
      quantity: parseFloat(document.getElementById("editQuantity").value),
    };

    try {
      const res = await fetch(`${API_BASE_URL}/api/v1/listings/${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        alert("Listing updated successfully!");
        closeEditModal();
        renderMyProduce();
      } else {
        alert("Failed to update listing.");
      }
    } catch (err) {
      console.error(err);
    }
  });

// Cancel/Delete Listing (DELETE /api/v1/listings/:id) (F9)
async function deleteListing(id) {
  if (!confirm("Are you sure you want to cancel this crop listing?")) return;

  const token = localStorage.getItem("token");
  try {
    const res = await fetch(`${API_BASE_URL}/api/v1/listings/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (res.ok) {
      alert("Listing canceled successfully!");
      renderMyProduce();
    } else {
      alert("Failed to cancel listing.");
    }
  } catch (err) {
    console.error(err);
  }
}

document.addEventListener("DOMContentLoaded", renderMyProduce);
window.addEventListener("listingCreated", renderMyProduce);

document.addEventListener("DOMContentLoaded", () => {
  const statusFilter = document.getElementById("statusFilter");
  const searchInput = document.getElementById("searchInput");
  if (!statusFilter && !searchInput) return;
  statusFilter?.addEventListener("change", renderMyProduce);
  searchInput?.addEventListener("input", renderMyProduce);
});

// ============================================================
// F10, F11, F12 — BUYER MARKETPLACE CATALOG & DETAILS DRAWER
// ============================================================
async function loadBrowseCatalog(crop = "", state = "") {
  const grid = document.getElementById("marketplaceCatalogGrid");
  if (!grid) return;

  const token = localStorage.getItem("token");

  let params = new URLSearchParams();
  if (crop) params.append("crop", crop);
  if (state) params.append("state", state);
  params.append("status", "ACTIVE");

  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/listings?${params.toString()}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: token ? `Bearer ${token}` : "",
        },
      },
    );

    const data = await res.json();

    if (res.ok && Array.isArray(data)) {
      if (data.length === 0) {
        grid.innerHTML = `
          <div class="col-span-full flex flex-col items-center justify-center py-16 text-center">
            <div class="text-4xl mb-4">📦</div>
            <h3 class="text-lg font-bold text-[#181D17] mb-1">No listings found</h3>
            <p class="text-sm text-[#40493D]">No produce matches your search. Try adjusting your filter criteria or check back later.</p>
          </div>`;
        return;
      }

      grid.innerHTML = data
        .map(
          (item, idx) => `
        <div class="glass-card rounded-2xl border border-[#E0E4DA]/60 p-5 shadow-sm flex flex-col justify-between hover:-translate-y-1.5 ksetu-fade-up" style="animation-delay:${idx * 60}ms">
          <div>
            <div class="flex justify-between items-start mb-2">
              <span class="px-2.5 py-0.5 text-xs font-semibold rounded-full bg-green-100 text-[#0D631B] flex items-center gap-1"><span class="w-1.5 h-1.5 rounded-full bg-[#0D631B]"></span>VERIFIED FARMER</span>
              <span class="text-xs text-[#40493D]">${item.district || ""}, ${item.state || ""}</span>
            </div>
            <h3 class="text-lg font-bold text-[#181D17]">${item.crop}</h3>
            <p class="text-xs text-[#40493D] mt-1 line-clamp-2">${item.description || "Fresh farm harvest"}</p>
            <div class="mt-4 pt-4 border-t border-[#F1F5EB] flex items-end justify-between">
              <p class="text-xs text-[#40493D]">Quantity: <strong class="text-[#181D17]">${item.quantity} ${item.unit}</strong></p>
              <p class="text-lg font-bold gradient-text-warm">₹${item.price}/<span class="text-xs">${item.unit}</span></p>
            </div>
          </div>
          <button onclick="showListingDetails(${item.id}, '${item.crop}', '${item.quantity} ${item.unit}', '₹${item.price} / ${item.unit}', '${item.district || ""}, ${item.state || ""}')" class="btn-warm w-full mt-4 py-2.5 text-white text-xs font-bold rounded-xl flex items-center justify-center gap-1.5">
            <span>👀</span> View Details & Offer
          </button>
        </div>
      `,
        )
        .join("");
    } else {
      grid.innerHTML = `<p class="text-red-600">Failed to fetch marketplace catalog.</p>`;
    }
  } catch (err) {
    console.error("Fetch Catalog Error:", err);
  }
}

// Search Filter Form Event Listener (F11)
document
  .getElementById("marketplaceFilterForm")
  ?.addEventListener("submit", function (e) {
    e.preventDefault();
    const crop = document.getElementById("searchCrop")?.value.trim() || "";
    const state = document.getElementById("searchState")?.value.trim() || "";
    loadBrowseCatalog(crop, state);
  });

// Drawer Controls (F12)
let currentOfferListingId = null;

function showListingDetails(listingId, crop, qtyLabel, priceLabel, location) {
  const listingContainer = document.getElementById("listingDetailsDrawer");
  currentOfferListingId = listingId;
  if (listingContainer) {
    const listingInfoEl = document.getElementById("offerListingInfo");
    if (listingInfoEl) {
      listingInfoEl.textContent = `${crop} — ${qtyLabel} at ${priceLabel}`;
    }
  }

  document.getElementById("drawerCropTitle").textContent = crop;
  document.getElementById("drawerQty").textContent = qtyLabel;
  document.getElementById("drawerPrice").textContent = priceLabel;
  document.getElementById("drawerLocation").textContent = location;

  const backdrop = document.getElementById("listingDrawerBackdrop");
  const drawer = document.getElementById("listingDetailsDrawer");

  backdrop?.classList.remove("hidden");
  setTimeout(() => {
    backdrop?.classList.remove("opacity-0");
    drawer?.classList.remove("translate-x-full");
  }, 10);
}

function closeListingDrawer() {
  const backdrop = document.getElementById("listingDrawerBackdrop");
  const drawer = document.getElementById("listingDetailsDrawer");

  drawer?.classList.add("translate-x-full");
  backdrop?.classList.add("opacity-0");
  setTimeout(() => backdrop?.classList.add("hidden"), 300);
}

// Offer Modal Handlers
function openOfferModal() {
  const backdrop = document.getElementById("offerModalBackdrop");
  const qty = document.getElementById("offerQuantity");
  const price = document.getElementById("offerPrice");
  const message = document.getElementById("offerMessage");
  if (qty) qty.value = "";
  if (price) price.value = "";
  if (message) message.value = "";
  backdrop?.classList.remove("hidden");
}

function closeOfferModal() {
  document.getElementById("offerModalBackdrop")?.classList.add("hidden");
}

async function submitOfferFromModal() {
  const quantity = parseFloat(document.getElementById("offerQuantity")?.value);
  const price = parseFloat(document.getElementById("offerPrice")?.value);
  const message = document.getElementById("offerMessage")?.value || "";

  if (!currentOfferListingId) {
    alert("Missing listing information. Please try again.");
    return;
  }

  if (!quantity || quantity <= 0) {
    alert("Please enter a valid quantity.");
    return;
  }

  if (!price || price < 0) {
    alert("Please enter a valid offered price.");
    return;
  }

  await submitOffer(currentOfferListingId, quantity, price, message);
}

document.addEventListener("DOMContentLoaded", () => {
  if (document.getElementById("marketplaceCatalogGrid")) {
    loadBrowseCatalog();
  }
});
// ============================================================
// F16 & F17 — BUYER DEMANDS CREATION & LISTING INTEGRATION
// ============================================================

function openDemandDrawer() {
  const backdrop = document.getElementById("demandDrawerBackdrop");
  const drawer = document.getElementById("createDemandDrawer");
  if (backdrop && drawer) {
    backdrop.classList.remove("hidden");
    setTimeout(() => {
      backdrop.classList.remove("opacity-0");
      drawer.classList.remove("translate-x-full");
    }, 10);
  }
}

function closeDemandDrawer() {
  const backdrop = document.getElementById("demandDrawerBackdrop");
  const drawer = document.getElementById("createDemandDrawer");
  if (backdrop && drawer) {
    drawer.classList.add("translate-x-full");
    backdrop.classList.add("opacity-0");
    setTimeout(() => {
      backdrop.classList.add("hidden");
    }, 300);
  }
}

// 1. Submit Demand (POST /api/v1/demands)
document
  .getElementById("createDemandForm")
  ?.addEventListener("submit", async function (e) {
    e.preventDefault();
    const token = localStorage.getItem("token");

    if (!token) {
      alert("Session expired. Please login again.");
      window.location.href = "./login.html?role=buyer";
      return;
    }

    const payload = {
      crop_name: document.getElementById("demandCrop")?.value.trim(),
      quantity: parseFloat(document.getElementById("demandQuantity")?.value),
      unit: document.getElementById("demandUnit")?.value,
      target_price: parseFloat(
        document.getElementById("demandTargetPrice")?.value,
      ),
      state: document.getElementById("demandState")?.value.trim(),
      district: document.getElementById("demandDistrict")?.value.trim(),
      required_by: document.getElementById("demandRequiredBy")?.value || "",
    };

    try {
      const submitBtn = e.target.querySelector('button[type="submit"]');
      const originalText = submitBtn.textContent;
      submitBtn.textContent = "Publishing...";
      submitBtn.disabled = true;

      const res = await fetch(
        `${API_BASE_URL}/api/v1/demands`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(payload),
        },
      );

      const data = await res.json();
      submitBtn.textContent = originalText;
      submitBtn.disabled = false;
      if (res.ok) {
        alert("Requirement posted successfully!");
        closeDemandDrawer();
        e.target.reset();
        renderMyDemands();
      } else {
        alert(extractErrorMessage(data, "Failed to post requirement."));
      }
    } catch (err) {
      console.error("Error posting demand:", err);
    }
  });

// 2. Render My Demands (GET /api/v1/demands/my)
let currentDemands = [];

async function renderMyDemands() {
  const demandsGrid = document.getElementById("demandsGrid");
  if (!demandsGrid) return;

  const token = localStorage.getItem("token");
  if (!token) {
    demandsGrid.innerHTML = `<p class="text-red-600 font-medium">Please login to view requirements.</p>`;
    return;
  }

  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/demands/my`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      },
    );

    const data = await res.json();
    if (res.ok && Array.isArray(data)) {
      currentDemands = data;
      if (data.length === 0) {
        demandsGrid.innerHTML = `
          <div class="col-span-full flex flex-col items-center justify-center py-16 text-center">
            <div class="text-4xl mb-4">📢</div>
            <h3 class="text-lg font-bold text-[#181D17] mb-1">No demands posted</h3>
            <p class="text-sm text-[#40493D]">Click "+ Post Requirement" to broadcast your first crop demand to farmers.</p>
          </div>`;
        return;
      }

      demandsGrid.innerHTML = data
        .map(
          (item, idx) => `
        <div class="glass-card rounded-2xl border border-[#E0E4DA]/60 p-5 shadow-sm hover:-translate-y-1 ksetu-fade-up" style="animation-delay:${idx * 60}ms">
          <div class="flex justify-between items-start mb-2">
            <span class="status-badge pending">${item.status || "OPEN"}</span>
            <span class="text-xs text-[#40493D]">📍 ${item.district || ""}, ${item.state || ""}</span>
          </div>
          <h3 class="text-lg font-bold text-[#181D17]">${item.crop_name}</h3>
          <p class="text-xs text-[#40493D] mt-1">Requirement by ${item.required_by ? String(item.required_by).slice(0, 10) : "ASAP"}</p>
          <div class="mt-4 pt-4 border-t border-[#F1F5EB] flex justify-between items-center">
            <div>
              <p class="text-xs text-[#40493D]">Required: <strong class="text-[#181D17]">${item.quantity} ${item.unit}</strong></p>
              <p class="text-lg font-bold gradient-text-warm">Max ₹${item.target_price} / ${item.unit}</p>
            </div>
            ${
              item.status === "ACTIVE"
                ? `<div class="flex gap-2">
              <button onclick="openDemandEditModal('${item.id}')" class="px-3 py-1.5 bg-amber-50 text-amber-700 border border-amber-200 text-xs font-semibold rounded-xl hover:bg-amber-100 transition-colors">✏️ Edit</button>
              <button onclick="openDemandMatches('${item.id}')" class="px-3 py-1.5 bg-[#75584D] text-white text-xs font-semibold rounded-xl hover:bg-[#5c443b] transition-colors">🔍 Find Matches</button>
            </div>`
                : `<a href="./browse-produce.html" class="px-3 py-1.5 bg-[#75584D] text-white text-xs font-semibold rounded-xl hover:bg-[#5c443b] transition-colors">🔍 Find Matches</a>`
            }
          </div>
        </div>
      `,
        )
        .join("");
    } else {
      demandsGrid.innerHTML = `<p class="text-red-600">Failed to load demands.</p>`;
    }
  } catch (err) {
    console.error("Error fetching demands:", err);
  }
}

document.addEventListener("DOMContentLoaded", renderMyDemands);

// ============================================================
// DEMAND EDIT (PUT /api/v1/demands/:id)
// ============================================================
function findDemandInCache(id) {
  return currentDemands.find((d) => String(d.id) === String(id)) || null;
}

function openDemandEditModal(id) {
  const modal = document.getElementById("demandEditModal");
  const demand = findDemandInCache(id);
  if (!modal || !demand) return;

  document.getElementById("demandEditId").value = demand.id;
  document.getElementById("demandEditCrop").value = demand.crop_name || "";
  document.getElementById("demandEditQuantity").value =
    demand.quantity != null ? demand.quantity : "";
  document.getElementById("demandEditUnit").value = demand.unit || "";
  document.getElementById("demandEditTargetPrice").value =
    demand.target_price != null ? demand.target_price : "";
  document.getElementById("demandEditState").value = demand.state || "";
  document.getElementById("demandEditDistrict").value =
    demand.district || "";
  const rb = demand.required_by
    ? String(demand.required_by).slice(0, 10)
    : "";
  document.getElementById("demandEditRequiredBy").value = rb;

  modal.classList.remove("hidden");
}

function closeDemandEditModal() {
  document.getElementById("demandEditModal")?.classList.add("hidden");
}

document
  .getElementById("demandEditForm")
  ?.addEventListener("submit", async function (e) {
    e.preventDefault();
    const id = document.getElementById("demandEditId").value;
    const token = localStorage.getItem("token");
    if (!token) {
      alert("Session expired. Please login again.");
      window.location.href = "./login.html?role=buyer";
      return;
    }

    const payload = {
      crop_name: document.getElementById("demandEditCrop").value.trim(),
      quantity: parseFloat(
        document.getElementById("demandEditQuantity").value,
      ),
      unit: document.getElementById("demandEditUnit").value,
      target_price: parseFloat(
        document.getElementById("demandEditTargetPrice").value,
      ),
      state: document.getElementById("demandEditState").value.trim(),
      district: document.getElementById("demandEditDistrict").value.trim(),
      required_by:
        document.getElementById("demandEditRequiredBy").value || "",
    };

    try {
      const res = await fetch(
        `${API_BASE_URL}/api/v1/demands/${id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(payload),
        },
      );
      const data = await res.json();

      if (res.ok) {
        alert("Requirement updated successfully!");
        closeDemandEditModal();
        renderMyDemands();
      } else {
        alert(
          extractErrorMessage(data, "Failed to update requirement."),
        );
      }
    } catch (err) {
      console.error("Error updating demand:", err);
      alert("Server error. Please check your connection and try again.");
    }
  });

// ============================================================
// MATCHING UI (farmer → matched buyer demands / buyer → matched farmer listings)
// ============================================================
function findListingInCache(id) {
  return currentListings.find((l) => String(l.id) === String(id)) || null;
}

function matchLevelClass(score) {
  if (score >= 90) return "accepted";
  if (score >= 75) return "pending";
  if (score >= 60) return "";
  return "rejected";
}

function matchScoreColor(pct) {
  if (pct >= 75) return "bg-[#0D631B]";
  if (pct >= 60) return "bg-amber-500";
  return "bg-[#9B7B6C]";
}

function matchCardHtml(match, detail, mode) {
  const pct = Math.round(match.score);
  const loc = detail
    ? `${detail.district || ""}${detail.state ? ", " + detail.state : ""}`
    : "";
  const title =
    mode === "farmer"
      ? detail
        ? detail.crop_name || "Buyer Requirement"
        : `Buyer Demand #${match.demand_id}`
      : detail
        ? detail.crop || "Farmer Listing"
        : `Farmer Listing #${match.listing_id}`;

  const reasonList =
    match.reasons && match.reasons.length
      ? match.reasons
          .map(
            (r) =>
              `<li class="flex items-start gap-2 text-xs text-[#40493D]"><span class="text-[#0D631B] font-bold mt-0.5">✓</span><span>${r}</span></li>`,
          )
          .join("")
      : `<li class="flex items-start gap-2 text-xs text-[#40493D]">No match notes available.</li>`;

  return `
    <div class="glass-card rounded-2xl border border-[#E0E4DA]/60 p-5 shadow-sm ksetu-fade-up">
      <div class="flex justify-between items-start gap-3">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 rounded-xl flex items-center justify-center text-lg font-bold text-white ${matchScoreColor(pct)}">${pct}%</div>
          <div>
            <h4 class="font-bold text-[#181D17]">${title}</h4>
            <span class="status-badge ${matchLevelClass(pct)}">${match.level || "Match"}</span>
          </div>
        </div>
        <span class="text-xs text-[#40493D]">📍 ${loc || "Location unknown"}</span>
      </div>
      ${
        detail
          ? `<div class="mt-3 pt-3 border-t border-[#F1F5EB] grid grid-cols-2 gap-2 text-xs">
              ${
                mode === "farmer"
                  ? `<span class="text-[#40493D]">Required Qty: <strong class="text-[#181D17]">${detail.quantity} ${detail.unit}</strong></span>
                     <span class="text-[#40493D]">Max Budget: <strong class="text-[#181D17]">₹${detail.target_price} / ${detail.unit}</strong></span>`
                  : `<span class="text-[#40493D]">Qty: <strong class="text-[#181D17]">${detail.quantity} ${detail.unit}</strong></span>
                     <span class="text-[#40493D]">Price: <strong class="text-[#181D17]">₹${detail.price} / ${detail.unit}</strong></span>`
              }
            </div>`
          : ""
      }
      <ul class="mt-3 space-y-1.5">${reasonList}</ul>
      <div class="mt-3 pt-3 border-t border-[#F1F5EB] grid grid-cols-2 sm:grid-cols-5 gap-2 text-center text-[10px] text-[#40493D]">
        <div><div class="text-sm font-bold text-[#181D17]">${Math.round(match.commodity_score || 0)}</div>Commodity</div>
        <div><div class="text-sm font-bold text-[#181D17]">${Math.round(match.quantity_score || 0)}</div>Quantity</div>
        <div><div class="text-sm font-bold text-[#181D17]">${Math.round(match.location_score || 0)}</div>Location</div>
        <div><div class="text-sm font-bold text-[#181D17]">${Math.round(match.price_score || 0)}</div>Price</div>
        <div><div class="text-sm font-bold text-[#181D17]">${Math.round(match.grade_score || 0)}</div>Grade</div>
      </div>
      ${
        mode === "buyer" && detail
          ? `<a href="./browse-produce.html" class="btn-warm w-full mt-4 py-2.5 text-white text-xs font-bold rounded-xl flex items-center justify-center gap-1.5"><span>🤝</span> View & Make Offer</a>`
          : ""
      }
    </div>`;
}

async function openListingMatches(listingId) {
  const modal = document.getElementById("matchesModal");
  const body = document.getElementById("matchesModalBody");
  const title = document.getElementById("matchesModalTitle");
  if (!modal || !body) return;

  const token = localStorage.getItem("token");
  if (!token) {
    alert("Session expired. Please login again.");
    window.location.href = "./login.html?role=farmer";
    return;
  }

  title.textContent = "Matching Buyers";
  body.innerHTML = `<p class="text-sm text-[#40493D] py-6 text-center">Finding suitable buyers for this listing…</p>`;
  modal.classList.remove("hidden");

  try {
    const genRes = await fetch(
      `${API_BASE_URL}/api/v1/matching/listings/${listingId}/generate`,
      {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    const genData = await genRes.json();
    const matches = genRes.ok && genData.data ? genData.data : [];

    if (!genRes.ok) {
      body.innerHTML = `<p class="text-sm text-red-600 py-6 text-center">Could not generate matches. Please try again.</p>`;
      return;
    }

    if (!matches.length) {
      body.innerHTML = `
        <div class="flex flex-col items-center py-12 text-center">
          <div class="text-4xl mb-3">🤝</div>
          <h3 class="text-lg font-bold text-[#181D17] mb-1">No suitable matches found yet</h3>
          <p class="text-sm text-[#40493D]">No active buyer requirements match this listing right now. Check back later.</p>
        </div>`;
      return;
    }

    let cards = "";
    for (const m of matches) {
      let detail = null;
      try {
        const dRes = await fetch(
          `${API_BASE_URL}/api/v1/demands/${m.demand_id}`,
          { headers: { Authorization: `Bearer ${token}` } },
        );
        if (dRes.ok) detail = await dRes.json();
      } catch (e) {
        /* ignore enrichment errors */
      }
      cards += matchCardHtml(m, detail, "farmer");
    }
    body.innerHTML = `<div class="grid grid-cols-1 gap-4">${cards}</div>`;
  } catch (err) {
    console.error("Error generating listing matches:", err);
    body.innerHTML = `<p class="text-sm text-red-600 py-6 text-center">Server error. Please check your connection and try again.</p>`;
  }
}

async function openDemandMatches(demandId) {
  const modal = document.getElementById("matchesModal");
  const body = document.getElementById("matchesModalBody");
  const title = document.getElementById("matchesModalTitle");
  if (!modal || !body) return;

  const token = localStorage.getItem("token");
  if (!token) {
    alert("Session expired. Please login again.");
    window.location.href = "./login.html?role=buyer";
    return;
  }

  title.textContent = "Matching Farmers";
  body.innerHTML = `<p class="text-sm text-[#40493D] py-6 text-center">Finding suitable farmers for this requirement…</p>`;
  modal.classList.remove("hidden");

  try {
    const genRes = await fetch(
      `${API_BASE_URL}/api/v1/matching/demands/${demandId}/generate`,
      {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    const genData = await genRes.json();
    const matches = genRes.ok && genData.data ? genData.data : [];

    if (!genRes.ok) {
      body.innerHTML = `<p class="text-sm text-red-600 py-6 text-center">Could not generate matches. Please try again.</p>`;
      return;
    }

    if (!matches.length) {
      body.innerHTML = `
        <div class="flex flex-col items-center py-12 text-center">
          <div class="text-4xl mb-3">🤝</div>
          <h3 class="text-lg font-bold text-[#181D17] mb-1">No suitable matches found yet</h3>
          <p class="text-sm text-[#40493D]">No active farmer listings match this requirement right now. Check back later.</p>
        </div>`;
      return;
    }

    let listings = [];
    try {
      const lRes = await fetch(
        `${API_BASE_URL}/api/v1/listings?${new URLSearchParams({ status: "ACTIVE" }).toString()}`,
        { headers: { Authorization: `Bearer ${token}` } },
      );
      if (lRes.ok) {
        const arr = await lRes.json();
        if (Array.isArray(arr)) listings = arr;
      }
    } catch (e) {
      /* ignore enrichment errors */
    }
    const listingMap = {};
    listings.forEach((l) => {
      listingMap[String(l.id)] = l;
    });

    let cards = "";
    matches.forEach((m) => {
      const detail = listingMap[String(m.listing_id)] || null;
      cards += matchCardHtml(m, detail, "buyer");
    });
    body.innerHTML = `<div class="grid grid-cols-1 gap-4">${cards}</div>`;
  } catch (err) {
    console.error("Error generating demand matches:", err);
    body.innerHTML = `<p class="text-sm text-red-600 py-6 text-center">Server error. Please check your connection and try again.</p>`;
  }
}

function closeMatchesModal() {
  document.getElementById("matchesModal")?.classList.add("hidden");
}

// ============================================================
// F20, F21 & F22 — OFFERS MANAGEMENT API INTEGRATION
// ============================================================

// 1. MAKE OFFER (F20 — POST /api/v1/offers)
async function submitOffer(listingId, quantity, price, message = "") {
  const token = localStorage.getItem("token");
  if (!token) {
    alert("Session expired. Please login again.");
    window.location.href = "./login.html?role=buyer";
    return;
  }

  const payload = {
    listing_id: listingId,
    quantity: parseFloat(quantity),
    offered_price: parseFloat(price),
    message: message,
  };

  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/offers`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      },
    );

    const data = await res.json();
    if (res.ok) {
      alert("Offer sent successfully to the farmer! 🎉");
      if (typeof closeOfferModal === "function") closeOfferModal();
      if (typeof closeListingDrawer === "function") closeListingDrawer();
    } else {
      alert(extractErrorMessage(data, "Failed to send offer."));
    }
  } catch (err) {
    console.error("Error sending offer:", err);
  }
}

// 2. BUYER SENT OFFERS (F21 — GET /api/v1/offers/buyer)
async function renderBuyerOffers() {
  const container = document.getElementById("buyerOffersList");
  if (!container) return;

  const token = localStorage.getItem("token");
  if (!token) {
    container.innerHTML = `<p class="text-red-600 font-medium">Please login to view sent offers.</p>`;
    return;
  }

  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/offers/buyer`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      },
    );

    const data = await res.json();
    if (res.ok && Array.isArray(data)) {
      if (data.length === 0) {
        container.innerHTML = `<p class="text-[#40493D]">No offers sent yet.</p>`;
        return;
      }

      container.innerHTML = data
        .map(
          (item) => `
        <div class="glass-card rounded-2xl border border-[#E0E4DA]/60 p-5 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-4 ksetu-fade-up">
          <div>
            <span class="status-badge ${
              item.status === "ACCEPTED"
                ? "accepted"
                : item.status === "REJECTED"
                  ? "rejected"
                  : "pending"
            }">${item.status || "PENDING"}</span>
            <h3 class="text-lg font-bold text-[#181D17] mt-2">${item.crop || "Produce Listing"}</h3>
            <p class="text-xs text-[#40493D] mt-1">💬 ${item.message || "No additional message"}</p>
          </div>
          <div class="text-right gap-2">
            <p class="text-xs text-[#40493D]">Offered Qty: <strong class="text-[#181D17]">${item.quantity}</strong></p>
            <p class="text-lg font-bold gradient-text-warm">₹${item.offered_price} / unit</p>
          </div>
        </div>
      `,
        )
        .join("");
    } else {
      container.innerHTML = `<p class="text-red-600">Failed to load offers.</p>`;
    }
  } catch (err) {
    console.error("Error fetching buyer offers:", err);
  }
}

// 3. FARMER RECEIVED OFFERS (F22 — GET /api/v1/offers/farmer & POST Accept/Reject)
async function renderFarmerOffers() {
  const container = document.getElementById("farmerOffersList");
  if (!container) return;

  const token = localStorage.getItem("token");
  if (!token) {
    container.innerHTML = `<p class="text-red-600 font-medium">Please login to view incoming offers.</p>`;
    return;
  }

  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/offers/farmer`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      },
    );

    const data = await res.json();
    if (res.ok && Array.isArray(data)) {
      if (data.length === 0) {
        container.innerHTML = `<p class="text-[#40493D]">No incoming offers received yet.</p>`;
        return;
      }

      container.innerHTML = data
        .map(
          (item) => `
        <div class="glass-card rounded-2xl border border-[#E0E4DA]/60 p-5 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-4 ksetu-fade-up">
          <div>
            <span class="status-badge ${
              item.status === "PENDING"
                ? "pending"
                : item.status === "ACCEPTED"
                  ? "accepted"
                  : "rejected"
            }">${item.status || "PENDING"}</span>
            <h3 class="text-lg font-bold text-[#181D17] mt-2">${item.crop || "Produce Listing"}</h3>
            <p class="text-xs text-[#40493D] mt-0.5">Offer from Buyer: <strong class="text-[#181D17]">${item.buyer_name || "Buyer"}</strong></p>
            <p class="text-xs text-[#40493D] mt-1">💬 ${item.message || "No note from buyer"}</p>
          </div>
          <div class="flex items-center gap-4 w-full md:w-auto justify-between md:justify-end">
            <div class="text-right">
              <p class="text-xs text-[#40493D]">Qty: <strong class="text-[#181D17]">${item.quantity}</strong></p>
              <p class="text-lg font-bold brand-name">₹${item.offered_price} / unit</p>
            </div>
            ${
              item.status === "PENDING"
                ? `
              <div class="flex gap-2">
                <button onclick="handleOfferAction('${item.id}', 'accept')" class="btn-warm px-5 py-2.5 text-white text-xs font-bold rounded-xl flex items-center gap-1"><span>✅</span> Accept</button>
                <button onclick="handleOfferAction('${item.id}', 'reject')" class="px-5 py-2.5 bg-red-50 text-red-600 border border-red-200 text-xs font-bold rounded-xl hover:bg-red-100 transition-colors">✕ Reject</button>
              </div>
            `
                : ""
            }
          </div>
        </div>
      `,
        )
        .join("");
    } else {
      container.innerHTML = `<p class="text-red-600">Failed to load received offers.</p>`;
    }
  } catch (err) {
    console.error("Error fetching farmer offers:", err);
  }
}

// 4. ACCEPT / REJECT OFFER ACTION
async function handleOfferAction(offerId, action) {
  const token = localStorage.getItem("token");
  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/offers/${offerId}/${action}`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    );

    if (res.ok) {
      alert(`Offer ${action}ed successfully!`);
      renderFarmerOffers();
    } else {
      alert(`Failed to ${action} offer.`);
    }
  } catch (err) {
    console.error(`Error during offer ${action}:`, err);
  }
}

document.addEventListener("DOMContentLoaded", () => {
  renderBuyerOffers();
  renderFarmerOffers();
});
// ============================================================
// F24, F25, F26, F27 — ORDERS MANAGEMENT API INTEGRATION
// ============================================================

// 1. FARMER ORDERS (F25 — GET /api/v1/orders/farmer)
async function renderFarmerOrders() {
  const container = document.getElementById("farmerOrdersList");
  if (!container) return;

  const token = localStorage.getItem("token");
  if (!token) {
    container.innerHTML = `<p class="text-red-600 font-medium">Please login to view orders.</p>`;
    return;
  }

  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/orders/farmer`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      },
    );

    const data = await res.json();
    if (res.ok && Array.isArray(data)) {
      if (data.length === 0) {
        container.innerHTML = `<p class="text-[#40493D]">No confirmed orders found.</p>`;
        return;
      }

      container.innerHTML = data
        .map(
          (item) => `
        <div class="glass-card rounded-2xl border border-[#E0E4DA]/60 p-5 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-4 ksetu-fade-up">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="status-badge accepted">${item.status || "CONFIRMED"}</span>
              <span class="text-xs text-gray-500">Order #${item.id ? String(item.id).slice(0, 8) : "N/A"}</span>
            </div>
            <h3 class="text-lg font-bold text-[#181D17]">${item.crop || "Crop Harvest"}</h3>
            <p class="text-xs text-[#40493D]">Buyer: <strong class="text-[#181D17]">${item.buyer_name || "Buyer Partner"}</strong></p>
          </div>
          <div class="flex items-center gap-4 w-full md:w-auto justify-between md:justify-end">
            <div class="text-right">
              <p class="text-xs text-[#40493D]">Agreed Qty: <strong class="text-[#181D17]">${item.quantity}</strong></p>
              <p class="text-lg font-bold brand-name">Total: ₹${item.total_amount || item.quantity * item.agreed_price}</p>
            </div>
            <button onclick="openOrderModal('${item.id}', '${item.crop}', '${item.quantity}', '${item.agreed_price}', '${item.total_amount || item.quantity * item.agreed_price}', '${item.buyer_name || "Buyer"}', '${item.status || "CONFIRMED"}')" class="px-4 py-2.5 bg-[#0D631B]/10 text-[#0D631B] text-xs font-bold rounded-xl hover:bg-[#0D631B] hover:text-white transition-all">
              📄 View Summary
            </button>
          </div>
        </div>
      `,
        )
        .join("");
    } else {
      container.innerHTML = `<p class="text-red-600">Failed to load farmer orders.</p>`;
    }
  } catch (err) {
    console.error("Error fetching farmer orders:", err);
  }
}

// 2. BUYER ORDERS (F26 — GET /api/v1/orders/buyer)
async function renderBuyerOrders() {
  const container = document.getElementById("buyerOrdersList");
  if (!container) return;

  const token = localStorage.getItem("token");
  if (!token) {
    container.innerHTML = `<p class="text-red-600 font-medium">Please login to view orders.</p>`;
    return;
  }

  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/orders/buyer`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      },
    );

    const data = await res.json();
    if (res.ok && Array.isArray(data)) {
      if (data.length === 0) {
        container.innerHTML = `<p class="text-[#40493D]">No purchase orders found.</p>`;
        return;
      }

      container.innerHTML = data
        .map(
          (item) => `
        <div class="glass-card rounded-2xl border border-[#E0E4DA]/60 p-5 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-4 ksetu-fade-up">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="status-badge accepted">${item.status || "CONFIRMED"}</span>
              <span class="text-xs text-gray-500">Order #${item.id ? String(item.id).slice(0, 8) : "N/A"}</span>
            </div>
            <h3 class="text-lg font-bold text-[#181D17]">${item.crop || "Crop Harvest"}</h3>
            <p class="text-xs text-[#40493D]">Farmer: <strong class="text-[#181D17]">${item.farmer_name || "Farmer Partner"}</strong></p>
          </div>
          <div class="flex items-center gap-4 w-full md:w-auto justify-between md:justify-end">
            <div class="text-right">
              <p class="text-xs text-[#40493D]">Agreed Qty: <strong class="text-[#181D17]">${item.quantity}</strong></p>
              <p class="text-lg font-bold gradient-text-warm">Total: ₹${item.total_amount || item.quantity * item.agreed_price}</p>
            </div>
            <button onclick="openOrderModal('${item.id}', '${item.crop}', '${item.quantity}', '${item.agreed_price}', '${item.total_amount || item.quantity * item.agreed_price}', '${item.farmer_name || "Farmer"}', '${item.status || "CONFIRMED"}')" class="px-4 py-2.5 bg-[#75584D]/10 text-[#75584D] text-xs font-bold rounded-xl hover:bg-[#75584D] hover:text-white transition-all">
              📄 View Summary
            </button>
          </div>
        </div>
      `,
        )
        .join("");
    } else {
      container.innerHTML = `<p class="text-red-600">Failed to load buyer orders.</p>`;
    }
  } catch (err) {
    console.error("Error fetching buyer orders:", err);
  }
}

// 3. ORDER DETAILS MODAL HANDLERS (F27)
function openOrderModal(id, crop, qty, price, total, partner, status) {
  const modal = document.getElementById("orderDetailsModal");
  const body = document.getElementById("orderDetailsBody");
  if (!modal || !body) return;

  body.innerHTML = `
    <div class="p-4 bg-gradient-to-br from-[#0D631B]/10 to-[#4CAF50]/10 rounded-xl border border-[#0D631B]/20 text-center mb-3">
      <span class="text-3xl">🎉</span>
      <h4 class="font-bold gradient-heading text-base mt-2">Order Confirmed</h4>
      <p class="text-xs text-[#40493D] mt-1">Deal finalized between both parties</p>
    </div>
    <div class="grid grid-cols-2 gap-2 text-xs border-b border-[#CDBDB4]/50 pb-2">
      <span class="text-[#40493D]">Order Reference:</span>
      <span class="font-semibold text-right text-[#181D17]">${id ? id.substring(0, 12) : "CR-8921"}...</span>
      <span class="text-[#40493D]">Status:</span>
      <span class="status-badge accepted ml-auto">${status}</span>
    </div>
    <div class="grid grid-cols-2 gap-2 text-xs border-b border-[#CDBDB4]/50 py-2">
      <span class="text-[#40493D]">Crop Commodity:</span>
      <span class="font-semibold text-right text-[#181D17]">${crop}</span>
      <span class="text-[#40493D]">Agreed Quantity:</span>
      <span class="font-semibold text-right text-[#181D17]">${qty}</span>
      <span class="text-[#40493D]">Agreed Unit Price:</span>
      <span class="font-semibold text-right text-[#181D17]">₹${price}</span>
    </div>
    <div class="flex justify-between items-center pt-3 text-base font-bold brand-name">
      <span>Total Amount:</span>
      <span class="text-xl">₹${total}</span>
    </div>
  `;

  modal.classList.remove("hidden");
}

function closeOrderModal() {
  const modal = document.getElementById("orderDetailsModal");
  if (modal) modal.classList.add("hidden");
}

document.addEventListener("DOMContentLoaded", () => {
  renderFarmerOrders();
  renderBuyerOrders();
});
// ============================================================
// PROFILE MANAGEMENT API INTEGRATION
// ============================================================
async function loadUserProfile() {
  const nameInput = document.getElementById("profileName");
  if (!nameInput) return;

  const token = localStorage.getItem("token");
  if (!token) {
    alert("Session expired. Please login again.");
    window.location.href = "./login.html";
    return;
  }

  try {
    const res = await fetch(
      `${API_BASE_URL}/api/v1/auth/me`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      },
    );

    const data = await res.json();
    if (res.ok && data) {
      document.getElementById("profileName").value = data.name || "";
      document.getElementById("profilePhone").value = data.phone || "";
      document.getElementById("profileState").value = data.state || "";
      document.getElementById("profileDistrict").value =
        data.district || data.city || "";
      document.getElementById("profileAddress").value =
        data.village || data.address || "";

      const businessField = document.getElementById("businessNameField");
      const businessInput = document.getElementById("profileBusiness");
      const businessTypeInput = document.getElementById("profileBusinessType");
      if (getTokenRole() === "buyer") {
        if (businessField) businessField.classList.remove("hidden");
        if (businessInput) businessInput.value = data.business_name || "";
        if (businessTypeInput)
          businessTypeInput.value = data.business_type || "";
      }

      const nameDisplay = document.getElementById("profileNameDisplay");
      if (nameDisplay) nameDisplay.textContent = data.name || "User Account";

      const profileAvatar = document.getElementById("profileAvatar");
      const profileRoleBadge = document.getElementById("profileRoleBadge");
      const displayName = data.business_name || data.name || "K";
      if (profileAvatar) profileAvatar.textContent = displayName.charAt(0).toUpperCase();
      if (profileRoleBadge) {
        profileRoleBadge.textContent =
          getTokenRole() === "buyer" ? "VERIFIED BUYER" : "VERIFIED FARMER";
      }
    }
  } catch (err) {
    console.error("Error loading profile:", err);
  }
}

document
  .getElementById("profileForm")
  ?.addEventListener("submit", async function (e) {
    e.preventDefault();

    const token = localStorage.getItem("token");
    if (!token) {
      alert("Session expired. Please login again.");
      window.location.href = "./login.html";
      return;
    }

    const name = document.getElementById("profileName").value.trim();
    const state = document.getElementById("profileState").value.trim();
    const district = document.getElementById("profileDistrict").value.trim();

    if (!name || !state || !district) {
      alert("Name, state and district are required.");
      return;
    }

    const role = getTokenRole();
    let endpoint = `${API_BASE_URL}/api/v1/farmers/me`;
    const payload = { name, state, district };

    if (role === "buyer") {
      endpoint = `${API_BASE_URL}/api/v1/buyers/me`;
      const businessName = document
        .getElementById("profileBusiness")
        ?.value.trim();
      const businessType = document
        .getElementById("profileBusinessType")
        ?.value.trim();
      if (!businessName) {
        alert("Business name is required for buyers.");
        return;
      }
      payload.business_name = businessName;
      payload.business_type = businessType;
    }

    try {
      const submitBtn = e.target.querySelector('button[type="submit"]');
      const originalText = submitBtn.textContent;
      submitBtn.textContent = "Saving...";
      submitBtn.disabled = true;

      const res = await fetch(endpoint, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      submitBtn.textContent = originalText;
      submitBtn.disabled = false;

      if (res.ok) {
        alert("Profile updated successfully! ✅");
        loadUserProfile();
      } else {
        alert(
          extractErrorMessage(data, "Failed to update profile. Please try again."),
        );
      }
    } catch (err) {
      console.error("Error updating profile:", err);
      alert("Server error. Please check your connection and try again.");
    }
  });

document.addEventListener("DOMContentLoaded", loadUserProfile);

// ============================================================
// FARMER DASHBOARD — STATS, SIDEBAR & RECENT ACTIVITY
// ============================================================
async function loadFarmerDashboard() {
  const nameEl = document.getElementById("farmerName");
  if (!nameEl) return;

  const token = localStorage.getItem("token");
  if (!token) {
    nameEl.textContent = "Guest Farmer";
    return;
  }

  const authedFetch = (url) =>
    fetch(`${API_BASE_URL}${url}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });

  try {
    const meRes = await authedFetch("/api/v1/auth/me");
    const me = meRes.ok ? await meRes.json() : null;
    if (me && me.name) {
      nameEl.textContent = me.name;
    }
  } catch (err) {
    console.error("Error loading farmer profile:", err);
  }

  try {
    const [listingsRes, offersRes, ordersRes] = await Promise.all([
      authedFetch("/api/v1/listings/my"),
      authedFetch("/api/v1/offers/farmer"),
      authedFetch("/api/v1/orders/farmer"),
    ]);

    const listings = listingsRes.ok ? await listingsRes.json() : [];
    const offers = offersRes.ok ? await offersRes.json() : [];
    const orders = ordersRes.ok ? await ordersRes.json() : [];

    const activeListings = Array.isArray(listings)
      ? listings.filter((l) => l.status === "ACTIVE").length
      : 0;
    const pendingOffers = Array.isArray(offers)
      ? offers.filter((o) => o.status === "PENDING").length
      : 0;
    const orderCount = Array.isArray(orders) ? orders.length : 0;
    let revenue = 0;
    if (Array.isArray(orders)) {
      orders.forEach((o) => {
        revenue += parseFloat(o.total_amount || 0) || 0;
      });
    }

    const statActive = document.getElementById("statActive");
    const statOffers = document.getElementById("statOffers");
    const statOrders = document.getElementById("statOrders");
    const statRevenue = document.getElementById("statRevenue");
    if (statActive) statActive.textContent = activeListings;
    if (statOffers) statOffers.textContent = pendingOffers;
    if (statOrders) statOrders.textContent = orderCount;
    if (statRevenue)
      statRevenue.textContent = `₹${revenue.toLocaleString("en-IN")}`;

    const recent = document.getElementById("recentActivity");
    if (recent) {
      if (
        activeListings === 0 &&
        pendingOffers === 0 &&
        orderCount === 0
      ) {
        return;
      }
      const items = [];
      if (Array.isArray(listings)) {
        listings.slice(0, 3).forEach((l) => {
          items.push(`
            <div class="flex items-center gap-3 p-3 rounded-xl bg-[#F4F8F0]/50">
              <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-[#0D631B] to-[#2E7D32] flex items-center justify-center text-white text-lg">🌾</div>
              <div class="flex-1">
                <p class="font-medium text-[#181D17]">${l.crop || "Crop"} listed ${l.status === "ACTIVE" ? "on marketplace" : `(${l.status})`}</p>
                <p class="text-sm text-[#40493D]">${l.quantity} ${l.unit} at ₹${l.price}</p>
              </div>
            </div>`);
        });
      }
      if (Array.isArray(orders)) {
        orders.slice(0, 2).forEach((o) => {
          items.push(`
            <div class="flex items-center gap-3 p-3 rounded-xl bg-[#F4F8F0]/50">
              <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-[#3B82F6] to-[#06B6D4] flex items-center justify-center text-white text-lg">📦</div>
              <div class="flex-1">
                <p class="font-medium text-[#181D17]">Order for ${o.crop || "produce"} ${o.status || "CONFIRMED"}</p>
                <p class="text-sm text-[#40493D]">${o.buyer_name || "Buyer"} • ₹${o.total_amount || ""}</p>
              </div>
            </div>`);
        });
      }
      recent.innerHTML = items.join("");
    }
  } catch (err) {
    console.error("Error loading farmer dashboard stats:", err);
  }
}

// ============================================================
// BUYER DASHBOARD — STATS, SIDEBAR & FEATURED LISTINGS
// ============================================================
async function loadBuyerDashboard() {
  const featured = document.getElementById("featuredListings");
  if (!featured && !document.getElementById("statPendingOffers")) return;

  const token = localStorage.getItem("token");

  const authedFetch = (url) =>
    fetch(`${API_BASE_URL}${url}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: token ? `Bearer ${token}` : "",
      },
    });

  try {
    if (token) {
      const [offersRes, ordersRes] = await Promise.all([
        authedFetch("/api/v1/offers/buyer"),
        authedFetch("/api/v1/orders/buyer"),
      ]);

      const offers = offersRes.ok ? await offersRes.json() : [];
      const orders = ordersRes.ok ? await ordersRes.json() : [];

      const pendingOffers = Array.isArray(offers)
        ? offers.filter((o) => o.status === "PENDING").length
        : 0;
      const acceptedDeals = Array.isArray(orders)
        ? orders.filter((o) => (o.status || "ACCEPTED") !== "REJECTED").length
        : 0;
      const deliveries = Array.isArray(orders)
        ? orders.filter((o) => (o.status || "").toUpperCase() !== "COMPLETED")
            .length
        : 0;

      const el1 = document.getElementById("statPendingOffers");
      const el2 = document.getElementById("statAcceptedDeals");
      const el3 = document.getElementById("statDeliveries");
      if (el1) el1.textContent = pendingOffers;
      if (el2) el2.textContent = acceptedDeals;
      if (el3) el3.textContent = deliveries;
    }
  } catch (err) {
    console.error("Error loading buyer dashboard stats:", err);
  }

  if (featured) {
    try {
      const res = await fetch(
        `${API_BASE_URL}/api/v1/listings?${new URLSearchParams({ status: "ACTIVE" }).toString()}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: token ? `Bearer ${token}` : "",
          },
        },
      );
      const data = res.ok ? await res.json() : [];
      if (Array.isArray(data) && data.length > 0) {
        featured.innerHTML = data
          .slice(0, 6)
          .map(
            (item) => `
          <div class="min-w-[260px] max-w-[280px] glass-card rounded-xl border border-[#E0E4DA]/60 p-4 flex flex-col justify-between hover:-translate-y-1 hover:shadow-lg transition-all duration-300 animate-fade-in-up">
            <div>
              <span class="inline-block px-2 py-0.5 text-xs font-semibold rounded bg-green-100 text-[#0D631B] mb-2">Verified Farmer</span>
              <h4 class="font-bold text-[#181D17] text-base">${item.crop || "Fresh Crop"}</h4>
              <p class="text-xs text-[#40493D] mb-3">${item.district || ""}${item.state ? ", " + item.state : ""}</p>
              <div class="space-y-1 text-sm">
                <p class="text-[#40493D]">Qty: <strong class="text-[#181D17]">${item.quantity} ${item.unit}</strong></p>
                <p class="text-[#75584D] font-bold text-base">₹${item.price} / ${item.unit}</p>
              </div>
            </div>
            <a href="./browse-produce.html" class="text-center w-full mt-4 py-2 bg-[#75584D] text-white text-xs font-bold rounded hover:bg-[#5c443b] transition-colors">View Details & Offer</a>
          </div>`,
          )
          .join("");
      } else {
        featured.innerHTML = `<p class="text-sm text-[#40493D]">No produce available right now. Check back soon! 🌱</p>`;
      }
    } catch (err) {
      console.error("Error loading featured listings:", err);
    }
  }
}

// ============================================================
// MOBILE MENU + LOGOUT (shared across app pages)
// ============================================================
function setupAppChrome() {
  const menuBtn = document.getElementById("mobileMenuBtn");
  const overlay = document.getElementById("mobileOverlay");
  const mobileMenu = document.getElementById("mobileMenu");

  function toggleMenu(open) {
    if (!mobileMenu || !overlay) return;
    mobileMenu.classList.toggle("hidden", !open);
    overlay.classList.toggle("hidden", !open);
    document.body.style.overflow = open ? "hidden" : "";
  }

  menuBtn?.addEventListener("click", () =>
    toggleMenu(mobileMenu.classList.contains("hidden")),
  );
  overlay?.addEventListener("click", () => toggleMenu(false));
  mobileMenu?.addEventListener("click", (e) => {
    if (e.target.tagName === "A") toggleMenu(false);
  });

  const wireLogout = (sel) => {
    const el = document.querySelector(sel);
    el?.addEventListener("click", (e) => {
      e.preventDefault();
      localStorage.clear();
      window.location.href = "./login.html";
    });
  };
  ["#sidebarLogout", "#mobileLogout", "#buyerLogout"].forEach(wireLogout);
}

// ============================================================
// SHARED SIDEBAR HYDRATION (all app pages)
// ============================================================
async function hydrateSidebar() {
  const nameEl = document.getElementById("sidebarName");
  const userNameEl = document.getElementById("sidebarUserName");
  const roleEl = document.getElementById("sidebarUserRole");
  const avatarEl = document.getElementById("sidebarAvatar");
  const emailEl = document.getElementById("sidebarEmail");
  if (!nameEl && !userNameEl && !avatarEl) return;

  const token = localStorage.getItem("token");
  if (!token) {
    document.querySelectorAll(".sidebar-logout").forEach((a) => {
      a.addEventListener("click", (e) => {
        e.preventDefault();
        localStorage.clear();
        window.location.href = "./login.html";
      });
    });
    return;
  }

  try {
    const res = await fetch(`${API_BASE_URL}/api/v1/auth/me`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });
    const me = res.ok ? await res.json() : null;
    if (me) {
      const role = getTokenRole();
      const displayName =
        me.business_name || me.name || (role === "buyer" ? "Buyer" : "Farmer");
      if (nameEl) nameEl.textContent = displayName;
      if (userNameEl) userNameEl.textContent = displayName;
      if (avatarEl)
        avatarEl.textContent = (displayName || "K").charAt(0).toUpperCase();
      if (roleEl)
        roleEl.textContent =
          role === "buyer" ? "Buyer Account" : "Farmer Account";
      if (emailEl) emailEl.textContent = me.phone ? `+91 ${me.phone}` : "";
    }
  } catch (err) {
    console.error("Error hydrating sidebar:", err);
  }

  document.querySelectorAll(".sidebar-logout").forEach((a) => {
    a.addEventListener("click", (e) => {
      e.preventDefault();
      localStorage.clear();
      window.location.href = "./login.html";
    });
  });
}

document.addEventListener("DOMContentLoaded", () => {
  hydrateSidebar();
  setupAppChrome();
  loadFarmerDashboard();
  loadBuyerDashboard();
});

// ============================================================
// MARKET INTELLIGENCE — LATEST REPORTED MANDI PRICES (index.html)
// ============================================================
async function loadMarketPrices() {
  const grid = document.getElementById("marketPriceGrid");
  if (!grid) return;

  const commodities = ["Wheat", "Potato", "Tomato", "Onion", "Paddy(Common)"];

  let html = "";
  for (const commodity of commodities) {
    let info = null;
    try {
      const res = await fetch(
        `${API_BASE_URL}/api/v1/market/prices/intelligence?${new URLSearchParams({ commodity })}`,
      );
      const data = await res.json();
      if (res.ok && data && data.data) info = data.data;
    } catch (err) {
      console.error("Error loading market price for", commodity, err);
    }

    if (info) {
      const freshnessLabel =
        info.freshness === "Today" || info.freshness === "1 day old"
          ? "Recent"
          : info.freshness || "—";
      html += `
        <div class="premium-card rounded-xl p-6" style="--card-accent: linear-gradient(90deg,#4e99d9,#9cc7ee);--card-glow:rgba(78,153,217,0.14);--card-shadow:rgba(78,153,217,0.16);">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-[#1E1E1E]">${info.commodity || commodity}</h3>
            <span class="inline-block px-2.5 py-0.5 text-[10px] font-semibold rounded-full bg-blue-100 text-[#2E6BA6]">${freshnessLabel}</span>
          </div>
          <p class="text-sm text-[#40493D] mt-1">Reported ₹<strong class="text-[#181D17]">${info.current_price ? info.current_price.toLocaleString("en-IN") : "—"}</strong> / quintal</p>
          <p class="text-xs text-[#40493D] mt-2">Range: ₹${info.min_price ? info.min_price.toLocaleString("en-IN") : "—"} – ₹${info.max_price ? info.max_price.toLocaleString("en-IN") : "—"}</p>
          <div class="mt-4 pt-3 border-t border-[#E0E4DA]/60 text-[11px] text-[#40493D] space-y-1">
            <p>📅 Reported on: ${info.reported_date || "—"}</p>
            <p>🏪 Source: ${info.source || "Mandi"}</p>
          </div>
        </div>`;
    } else {
      html += `
        <div class="premium-card rounded-xl p-6 text-center" style="--card-accent: linear-gradient(90deg,#4e99d9,#9cc7ee);--card-glow:rgba(78,153,217,0.14);--card-shadow:rgba(78,153,217,0.16);">
          <h3 class="text-lg font-semibold text-[#1E1E1E]">${commodity}</h3>
          <p class="text-sm text-[#40493D] mt-2">Latest report not available yet.</p>
        </div>`;
    }
  }

  grid.innerHTML = html;
}

document.addEventListener("DOMContentLoaded", loadMarketPrices);
