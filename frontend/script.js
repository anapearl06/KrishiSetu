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

  // Restrict phone input to numbers only
  const phoneInputs = document.querySelectorAll('input[type="tel"]');
  phoneInputs.forEach((input) => {
    input.addEventListener("input", function () {
      input.value = input.value.replace(/\D/g, "").slice(0, 10);
    });
  });
});
// ============================================================
// F5 — CREATE CROP LISTING (POST /api/v1/listings)
// ============================================================
const API_BASE_URL = "https://krishisetu-api-tiau.onrender.com";

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
        if (drawer) drawer.classList.add("translate-x-full");
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
