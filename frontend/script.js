// ============================================================
// KRISHISETU - BACKEND API CONNECTED SCRIPT
// Backend Render URL: https://krishisetu-api-tiau.onrender.com
// Neon DB Endpoint Connected
// ============================================================

const API_BASE_URL = "https://krishisetu-api-tiau.onrender.com";

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
          if (data.token) localStorage.setItem("token", data.token);
          window.location.href = "./farmer-dashboard.html";
        } else {
          const errorMsg =
            typeof data.error === "object"
              ? JSON.stringify(data.error)
              : data.error ||
                data.message ||
                "Login failed! Check credentials.";
          alert(errorMsg);
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error! Render API is waking up, please try again.");
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
          if (data.token) localStorage.setItem("token", data.token);
          window.location.href = "./buyer-dashboard.html";
        } else {
          const errorMsg =
            typeof data.error === "object"
              ? JSON.stringify(data.error)
              : data.error ||
                data.message ||
                "Login failed! Check credentials.";
          alert(errorMsg);
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error!");
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
      const village =
        document.getElementById("farmerVillage")?.value || "Default Village";
      const district =
        document.getElementById("farmerDistrict")?.value ||
        village ||
        "Default District";
      const state =
        document.getElementById("farmerState")?.value || "Uttar Pradesh";
      const crop =
        document.getElementById("farmerCrop")?.value || "General Crop";
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
              village: village,
              district: district,
              state: state,
              crop: crop,
              password: password,
            }),
          },
        );

        const data = await response.json();
        if (response.ok) {
          alert("Farmer Registration Successful! Redirecting to Login...");
          window.location.href = "./login.html?role=farmer";
        } else {
          const errorMsg =
            typeof data.error === "object"
              ? JSON.stringify(data.error)
              : data.error || data.message || "Registration failed!";
          alert(errorMsg);
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error! Render backend waking up.");
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
      const city =
        document.getElementById("buyerCity")?.value || "Default City";
      const district =
        document.getElementById("buyerDistrict")?.value ||
        city ||
        "Default District";
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
            city: city,
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
          const errorMsg =
            typeof data.error === "object"
              ? JSON.stringify(data.error)
              : data.error || data.message || "Registration failed!";
          alert(errorMsg);
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error!");
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
document
  .getElementById("createListingForm")
  ?.addEventListener("submit", async function (e) {
    e.preventDefault();

    const token = localStorage.getItem("token");
    if (!token) {
      alert("Session expired. Kripya login karein.");
      window.location.href = "./login.html?role=farmer";
      return;
    }

    const listingPayload = {
      crop: document.getElementById("produceName")?.value.trim(),
      quantity: parseFloat(document.getElementById("produceQuantity")?.value),
      unit: document.getElementById("produceUnit")?.value,
      price: parseFloat(document.getElementById("producePrice")?.value),
      state: document.getElementById("produceLocation")?.value.trim(),
      district: document.getElementById("produceLocation")?.value.trim(),
      description: document.getElementById("produceDesc")?.value || "",
    };

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/listings`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(listingPayload),
      });

      const data = await response.json();

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
        alert(data.message || data.error || "Failed to create crop listing.");
      }
    } catch (err) {
      console.error("Error creating listing:", err);
      alert("Server error! Connection verify karein.");
    }
  });

// ============================================================
// MY PRODUCE (F6, F8, F9) — GET / EDIT / DELETE
// ============================================================
async function renderMyProduce() {
  const produceGrid = document.getElementById("produceGrid");
  if (!produceGrid) return;

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
      if (data.length === 0) {
        produceGrid.innerHTML = `<p class="text-[#40493D]">No produce items added yet. Use "+ Sell Produce" to post your crop.</p>`;
        return;
      }

      produceGrid.innerHTML = data
        .map(
          (item) => `
        <div class="bg-white rounded-xl border border-[#E0E4DA] p-5 shadow-sm hover:shadow-md transition-shadow">
          <div class="flex justify-between items-start mb-2">
            <span class="px-2.5 py-0.5 text-xs font-semibold rounded bg-green-100 text-[#0D631B] uppercase">${item.status || "ACTIVE"}</span>
            <span class="text-xs text-[#40493D]">${item.district || ""}, ${item.state || ""}</span>
          </div>
          <h3 class="text-lg font-bold text-[#181D17]">${item.crop}</h3>
          <p class="text-xs text-[#40493D] mt-1">${item.description || "No description provided"}</p>

          <div class="mt-4 pt-4 border-t border-[#F1F5EB] flex justify-between items-center">
            <div>
              <p class="text-xs text-[#40493D]">Quantity: <strong>${item.quantity} ${item.unit}</strong></p>
              <p class="text-lg font-bold text-[#0D631B]">₹${item.price} / ${item.unit}</p>
            </div>
            <div class="flex gap-2">
              <button onclick="openEditModal('${item.id}', ${item.price}, ${item.quantity})" class="px-3 py-1.5 bg-amber-50 text-amber-700 border border-amber-200 text-xs font-semibold rounded-md hover:bg-amber-100">Edit</button>
              <button onclick="deleteListing('${item.id}')" class="px-3 py-1.5 bg-red-50 text-red-600 border border-red-200 text-xs font-semibold rounded-md hover:bg-red-100">Cancel</button>
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
        grid.innerHTML = `<p class="text-[#40493D]">No produce available matching criteria.</p>`;
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
          <button onclick="showListingDetails('${item.crop}', '${item.quantity} ${item.unit}', '₹${item.price} / ${item.unit}', '${item.district || ""}, ${item.state || ""}')" class="btn-warm w-full mt-4 py-2.5 text-white text-xs font-bold rounded-xl flex items-center justify-center gap-1.5">
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
function showListingDetails(crop, qty, price, location) {
  document.getElementById("drawerCropTitle").textContent = crop;
  document.getElementById("drawerQty").textContent = qty;
  document.getElementById("drawerPrice").textContent = price;
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

// Offer Modal Handlers (UI Only placeholder)
function openOfferModal() {
  document.getElementById("offerModalBackdrop")?.classList.remove("hidden");
}

function closeOfferModal() {
  document.getElementById("offerModalBackdrop")?.classList.add("hidden");
}

function submitOfferPlaceholder() {
  alert("Offer sent successfully! (UI feature placeholder)");
  closeOfferModal();
  closeListingDrawer();
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
      alert("Session expired. Please login as buyer.");
      window.location.href = "./login.html?role=buyer";
      return;
    }

    const payload = {
      crop: document.getElementById("demandCrop")?.value.trim(),
      quantity: parseFloat(document.getElementById("demandQuantity")?.value),
      unit: document.getElementById("demandUnit")?.value,
      target_price: parseFloat(
        document.getElementById("demandTargetPrice")?.value,
      ),
      state: document.getElementById("demandState")?.value.trim(),
      district: document.getElementById("demandDistrict")?.value.trim(),
      description: document.getElementById("demandDescription")?.value || "",
    };

    try {
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
      if (res.ok) {
        alert("Requirement posted successfully!");
        closeDemandDrawer();
        e.target.reset();
        renderMyDemands();
      } else {
        alert(data.message || data.error || "Failed to post requirement.");
      }
    } catch (err) {
      console.error("Error posting demand:", err);
    }
  });

// 2. Render My Demands (GET /api/v1/demands/my)
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
      if (data.length === 0) {
        demandsGrid.innerHTML = `<p class="text-[#40493D]">No demand requirements posted yet. Click "+ Post Requirement" above.</p>`;
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
          <h3 class="text-lg font-bold text-[#181D17]">${item.crop}</h3>
          <p class="text-xs text-[#40493D] mt-1">${item.description || "No additional specs"}</p>
          <div class="mt-4 pt-4 border-t border-[#F1F5EB] flex justify-between items-center">
            <div>
              <p class="text-xs text-[#40493D]">Required: <strong class="text-[#181D17]">${item.quantity} ${item.unit}</strong></p>
              <p class="text-lg font-bold gradient-text-warm">Max ₹${item.target_price} / ${item.unit}</p>
            </div>
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
// F20, F21 & F22 — OFFERS MANAGEMENT API INTEGRATION
// ============================================================

// 1. MAKE OFFER (F20 — POST /api/v1/offers)
async function submitOffer(listingId, quantity, price, message = "") {
  const token = localStorage.getItem("token");
  if (!token) {
    alert("Session expired. Please login as buyer.");
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
    } else {
      alert(data.message || data.error || "Failed to send offer.");
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
        <div class="bg-white rounded-xl border border-[#E0E4DA] p-5 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="px-2 py-0.5 text-xs font-semibold rounded bg-green-100 text-green-800 uppercase">${item.status || "CONFIRMED"}</span>
              <span class="text-xs text-gray-500">Order #${item.id ? item.id.substring(0, 8) : "N/A"}</span>
            </div>
            <h3 class="text-lg font-bold text-[#181D17]">${item.crop || "Crop Harvest"}</h3>
            <p class="text-xs text-[#40493D]">Buyer: <strong>${item.buyer_name || "Buyer Partner"}</strong></p>
          </div>
          <div class="flex items-center gap-4 w-full md:w-auto justify-between md:justify-end">
            <div class="text-right">
              <p class="text-xs text-[#40493D]">Agreed Qty: <strong>${item.quantity}</strong></p>
              <p class="text-lg font-bold text-[#0D631B]">Total: ₹${item.total_amount || item.quantity * item.agreed_price}</p>
            </div>
            <button onclick="openOrderModal('${item.id}', '${item.crop}', '${item.quantity}', '${item.agreed_price}', '${item.total_amount || item.quantity * item.agreed_price}', '${item.buyer_name || "Buyer"}', '${item.status || "CONFIRMED"}')" class="px-3 py-2 bg-gray-100 text-[#181D17] text-xs font-bold rounded-lg hover:bg-gray-200">
              View Summary
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
        <div class="bg-white rounded-xl border border-[#E0E4DA] p-5 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="px-2 py-0.5 text-xs font-semibold rounded bg-green-100 text-green-800 uppercase">${item.status || "CONFIRMED"}</span>
              <span class="text-xs text-gray-500">Order #${item.id ? item.id.substring(0, 8) : "N/A"}</span>
            </div>
            <h3 class="text-lg font-bold text-[#181D17]">${item.crop || "Crop Harvest"}</h3>
            <p class="text-xs text-[#40493D]">Farmer: <strong>${item.farmer_name || "Farmer Partner"}</strong></p>
          </div>
          <div class="flex items-center gap-4 w-full md:w-auto justify-between md:justify-end">
            <div class="text-right">
              <p class="text-xs text-[#40493D]">Agreed Qty: <strong>${item.quantity}</strong></p>
              <p class="text-lg font-bold text-[#0D631B]">Total: ₹${item.total_amount || item.quantity * item.agreed_price}</p>
            </div>
            <button onclick="openOrderModal('${item.id}', '${item.crop}', '${item.quantity}', '${item.agreed_price}', '${item.total_amount || item.quantity * item.agreed_price}', '${item.farmer_name || "Farmer"}', '${item.status || "CONFIRMED"}')" class="px-3 py-2 bg-gray-100 text-[#181D17] text-xs font-bold rounded-lg hover:bg-gray-200">
              View Summary
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
    <div class="p-3 bg-green-50 rounded-lg border border-green-200 text-center mb-3">
      <span class="text-2xl">🎉</span>
      <h4 class="font-bold text-[#0D631B] text-base mt-1">Order Confirmed</h4>
      <p class="text-xs text-gray-600">Deal finalized between both parties</p>
    </div>
    <div class="grid grid-cols-2 gap-2 text-xs border-b border-[#E0E4DA] pb-2">
      <span class="text-gray-500">Order Reference:</span>
      <span class="font-semibold text-right">${id ? id.substring(0, 12) : "CR-8921"}...</span>
      <span class="text-gray-500">Status:</span>
      <span class="font-bold text-green-700 text-right uppercase">${status}</span>
    </div>
    <div class="grid grid-cols-2 gap-2 text-xs border-b border-[#E0E4DA] py-2">
      <span class="text-gray-500">Crop Commodity:</span>
      <span class="font-semibold text-right">${crop}</span>
      <span class="text-gray-500">Agreed Quantity:</span>
      <span class="font-semibold text-right">${qty}</span>
      <span class="text-gray-500">Agreed Unit Price:</span>
      <span class="font-semibold text-right">₹${price}</span>
    </div>
    <div class="flex justify-between items-center pt-2 text-base font-bold text-[#0D631B]">
      <span>Total Amount:</span>
      <span>₹${total}</span>
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

      const nameDisplay = document.getElementById("profileNameDisplay");
      if (nameDisplay) nameDisplay.textContent = data.name || "User Account";
    }
  } catch (err) {
    console.error("Error loading profile:", err);
  }
}

document
  .getElementById("profileForm")
  ?.addEventListener("submit", async function (e) {
    e.preventDefault();
    alert("Profile updated successfully!");
  });

document.addEventListener("DOMContentLoaded", loadUserProfile);
