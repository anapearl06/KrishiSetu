// ============================================================
// KRISHISETU - BACKEND API CONNECTED SCRIPT
// Backend Render URL: https://krishisetu-api-tiau.onrender.com
// ============================================================

const API_BASE_URL = "https://krishisetu-api-tiau.onrender.com";

document.addEventListener("DOMContentLoaded", function () {
  // HELPER: URL se role fetch karne ke liye
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
          alert(
            data.message || "Login failed! Please check phone and password.",
          );
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error! Render backend waking up.");
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
          alert(
            data.message || "Login failed! Please check phone and password.",
          );
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error!");
      }
    });
  }

  // ============================================================
  // REGISTER PAGE LOGIC
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
      const name = document.getElementById("farmerName")?.value;
      const phone = document.getElementById("farmerPhone")?.value;
      const village = document.getElementById("farmerVillage")?.value;
      const state = document.getElementById("farmerState")?.value;
      const crop = document.getElementById("farmerCrop")?.value;
      const password = document.getElementById("farmerPassword")?.value;
      const confirmPassword = document.getElementById(
        "farmerConfirmPassword",
      )?.value;

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
              name,
              phone,
              village,
              state,
              crop,
              password,
            }),
          },
        );

        const data = await response.json();
        if (response.ok) {
          alert("Farmer Registration Successful! Redirecting to Login...");
          window.location.href = "./login.html?role=farmer";
        } else {
          alert(data.message || "Registration failed!");
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error!");
      }
    });
  }

  // BUYER REGISTER BACKEND CALL
  if (buyerForm) {
    buyerForm.addEventListener("submit", async function (event) {
      event.preventDefault();
      const name = document.getElementById("buyerName")?.value;
      const businessName = document.getElementById("businessName")?.value;
      const businessType = document.getElementById("businessType")?.value;
      const phone = document.getElementById("buyerPhone")?.value;
      const city = document.getElementById("buyerCity")?.value;
      const state = document.getElementById("buyerState")?.value;
      const password = document.getElementById("buyerPassword")?.value;
      const confirmPassword = document.getElementById(
        "buyerConfirmPassword",
      )?.value;

      if (password !== confirmPassword) {
        alert("Passwords do not match!");
        return;
      }

      try {
        const response = await fetch(`${API_BASE_URL}/api/v1/buyers/register`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name,
            businessName,
            businessType,
            phone,
            city,
            state,
            password,
          }),
        });

        const data = await response.json();
        if (response.ok) {
          alert("Buyer Registration Successful! Redirecting to Login...");
          window.location.href = "./login.html?role=buyer";
        } else {
          alert(data.message || "Registration failed!");
        }
      } catch (error) {
        console.error("Error connecting to backend:", error);
        alert("Server Error!");
      }
    });
  }

  // Phone input restricts to numbers only
  const phoneInputs = document.querySelectorAll('input[type="tel"]');
  phoneInputs.forEach((input) => {
    input.addEventListener("input", function () {
      input.value = input.value.replace(/\D/g, "").slice(0, 10);
    });
  });
});
