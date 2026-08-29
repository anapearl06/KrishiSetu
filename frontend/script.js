// ============================================================
// KRISHISETU
// LOGIN + REGISTER
// FARMER / BUYER ROLE SWITCHING
// ============================================================

document.addEventListener("DOMContentLoaded", function () {
  // ============================================================
  // HELPER
  // ============================================================

  function getRoleFromURL() {
    const params = new URLSearchParams(window.location.search);
    return params.get("role");
  }

  // ============================================================
  // LOGIN PAGE
  // ============================================================

  const loginFarmerTab = document.getElementById("loginFarmerTab");
  const loginBuyerTab = document.getElementById("loginBuyerTab");

  const farmerLoginForm = document.getElementById("farmerLoginForm");
  const buyerLoginForm = document.getElementById("buyerLoginForm");

  const loginTitle = document.getElementById("loginTitle");
  const loginRegisterLink = document.getElementById("loginRegisterLink");

  // ============================================================
  // LOGIN — FARMER
  // ============================================================

  function activateFarmerLogin() {
    if (!farmerLoginForm || !buyerLoginForm) return;

    farmerLoginForm.classList.remove("hidden");
    buyerLoginForm.classList.add("hidden");

    if (loginFarmerTab) {
      loginFarmerTab.classList.add("active");
    }

    if (loginBuyerTab) {
      loginBuyerTab.classList.remove("active");
    }

    if (loginTitle) {
      loginTitle.textContent = "Farmer Login";
    }

    if (loginRegisterLink) {
      loginRegisterLink.href = "./register.html?role=farmer";
    }

    document.body.classList.remove("buyer-theme");
    document.body.classList.add("farmer-theme");
  }

  // ============================================================
  // LOGIN — BUYER
  // ============================================================

  function activateBuyerLogin() {
    if (!farmerLoginForm || !buyerLoginForm) return;

    buyerLoginForm.classList.remove("hidden");
    farmerLoginForm.classList.add("hidden");

    if (loginBuyerTab) {
      loginBuyerTab.classList.add("active");
    }

    if (loginFarmerTab) {
      loginFarmerTab.classList.remove("active");
    }

    if (loginTitle) {
      loginTitle.textContent = "Buyer Login";
    }

    if (loginRegisterLink) {
      loginRegisterLink.href = "./register.html?role=buyer";
    }

    document.body.classList.remove("farmer-theme");
    document.body.classList.add("buyer-theme");
  }

  // ============================================================
  // LOGIN TAB EVENTS
  // ============================================================

  if (loginFarmerTab) {
    loginFarmerTab.addEventListener("click", activateFarmerLogin);
  }

  if (loginBuyerTab) {
    loginBuyerTab.addEventListener("click", activateBuyerLogin);
  }

  // ============================================================
  // LOGIN URL ROLE
  // ============================================================

  if (loginFarmerTab && loginBuyerTab) {
    const role = getRoleFromURL();

    if (role === "buyer") {
      activateBuyerLogin();
    } else {
      activateFarmerLogin();
    }
  }

  // ============================================================
  // LOGIN FORM SUBMIT
  // ============================================================

  if (farmerLoginForm) {
    farmerLoginForm.addEventListener("submit", function (event) {
      event.preventDefault();

      alert("Farmer login UI is ready.");
    });
  }

  if (buyerLoginForm) {
    buyerLoginForm.addEventListener("submit", function (event) {
      event.preventDefault();

      alert("Buyer login UI is ready.");
    });
  }

  // ============================================================
  // REGISTER PAGE
  // ============================================================

  const registerFarmerTab = document.getElementById("registerFarmerTab");

  const registerBuyerTab = document.getElementById("registerBuyerTab");

  const farmerForm = document.getElementById("farmerForm");

  const buyerForm = document.getElementById("buyerForm");

  // ============================================================
  // REGISTER — FARMER
  // ============================================================

  function activateFarmerRegister() {
    if (!farmerForm || !buyerForm) return;

    farmerForm.classList.remove("hidden");
    farmerForm.classList.add("flex");

    buyerForm.classList.add("hidden");
    buyerForm.classList.remove("flex");

    if (registerFarmerTab) {
      registerFarmerTab.classList.add("active");
    }

    if (registerBuyerTab) {
      registerBuyerTab.classList.remove("active");
    }

    document.body.classList.remove("buyer-theme");
    document.body.classList.add("farmer-theme");
  }

  // ============================================================
  // REGISTER — BUYER
  // ============================================================

  function activateBuyerRegister() {
    if (!farmerForm || !buyerForm) return;

    buyerForm.classList.remove("hidden");
    buyerForm.classList.add("flex");

    farmerForm.classList.add("hidden");
    farmerForm.classList.remove("flex");

    if (registerBuyerTab) {
      registerBuyerTab.classList.add("active");
    }

    if (registerFarmerTab) {
      registerFarmerTab.classList.remove("active");
    }

    document.body.classList.remove("farmer-theme");
    document.body.classList.add("buyer-theme");
  }

  // ============================================================
  // REGISTER TAB EVENTS
  // ============================================================

  if (registerFarmerTab) {
    registerFarmerTab.addEventListener("click", activateFarmerRegister);
  }

  if (registerBuyerTab) {
    registerBuyerTab.addEventListener("click", activateBuyerRegister);
  }

  // ============================================================
  // REGISTER URL ROLE
  // ============================================================

  if (registerFarmerTab && registerBuyerTab) {
    const role = getRoleFromURL();

    if (role === "buyer") {
      activateBuyerRegister();
    } else {
      activateFarmerRegister();
    }
  }

  // ============================================================
  // REGISTER FORM SUBMIT
  // ============================================================

  if (farmerForm) {
    farmerForm.addEventListener("submit", function (event) {
      event.preventDefault();

      const password = document.getElementById("farmerPassword");

      const confirmPassword = document.getElementById("farmerConfirmPassword");

      if (
        password &&
        confirmPassword &&
        password.value !== confirmPassword.value
      ) {
        alert("Passwords do not match.");
        confirmPassword.focus();
        return;
      }

      alert("Farmer registration UI is ready.");
    });
  }

  if (buyerForm) {
    buyerForm.addEventListener("submit", function (event) {
      event.preventDefault();

      const password = document.getElementById("buyerPassword");

      const confirmPassword = document.getElementById("buyerConfirmPassword");

      if (
        password &&
        confirmPassword &&
        password.value !== confirmPassword.value
      ) {
        alert("Passwords do not match.");
        confirmPassword.focus();
        return;
      }

      alert("Buyer registration UI is ready.");
    });
  }

  // ============================================================
  // PHONE NUMBER — ONLY NUMBERS
  // ============================================================

  const phoneInputs = document.querySelectorAll('input[type="tel"]');

  phoneInputs.forEach(function (input) {
    input.addEventListener("input", function () {
      input.value = input.value.replace(/\D/g, "").slice(0, 10);
    });
  });
});
