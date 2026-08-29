// ========================================
// KRISHISETU ROLE TAB SWITCHING
// LOGIN + REGISTRATION
// ========================================

document.addEventListener("DOMContentLoaded", function () {
  // ========================================
  // LOGIN PAGE
  // ========================================

  const loginFarmerTab = document.getElementById("loginFarmerTab");
  const loginBuyerTab = document.getElementById("loginBuyerTab");

  const farmerLoginForm = document.getElementById("farmerLoginForm");
  const buyerLoginForm = document.getElementById("buyerLoginForm");

  // ----------------------------------------
  // LOGIN - FARMER
  // ----------------------------------------

  if (loginFarmerTab && loginBuyerTab && farmerLoginForm && buyerLoginForm) {
    loginFarmerTab.addEventListener("click", function () {
      // Show Farmer form
      farmerLoginForm.classList.remove("hidden");
      buyerLoginForm.classList.add("hidden");

      // Farmer active
      loginFarmerTab.classList.add(
        "text-[#0D631B]",
        "font-bold",
        "border-b-2",
        "border-[#0D631B]",
        "bg-[#0D631B]/5",
      );

      // Buyer inactive
      loginBuyerTab.classList.remove(
        "text-[#0D631B]",
        "font-bold",
        "border-b-2",
        "border-[#0D631B]",
        "bg-[#0D631B]/5",
      );

      loginBuyerTab.classList.add("text-[#40493D]", "font-medium");
    });

    // ----------------------------------------
    // LOGIN - BUYER
    // ----------------------------------------

    loginBuyerTab.addEventListener("click", function () {
      // Show Buyer form
      buyerLoginForm.classList.remove("hidden");
      farmerLoginForm.classList.add("hidden");

      // Buyer active
      loginBuyerTab.classList.add(
        "text-[#0D631B]",
        "font-bold",
        "border-b-2",
        "border-[#0D631B]",
        "bg-[#0D631B]/5",
      );

      // Farmer inactive
      loginFarmerTab.classList.remove(
        "text-[#0D631B]",
        "font-bold",
        "border-b-2",
        "border-[#0D631B]",
        "bg-[#0D631B]/5",
      );

      loginFarmerTab.classList.add("text-[#40493D]", "font-medium");
    });

    // ----------------------------------------
    // FARMER LOGIN SUBMIT
    // ----------------------------------------

    farmerLoginForm.addEventListener("submit", function (event) {
      event.preventDefault();

      alert("Farmer Login UI is ready.");
    });

    // ----------------------------------------
    // BUYER LOGIN SUBMIT
    // ----------------------------------------

    buyerLoginForm.addEventListener("submit", function (event) {
      event.preventDefault();

      alert("Buyer Login UI is ready.");
    });
  }

  // ========================================
  // REGISTRATION PAGE
  // ========================================

  const registerFarmerTab = document.getElementById("registerFarmerTab");
  const registerBuyerTab = document.getElementById("registerBuyerTab");

  const farmerForm = document.getElementById("farmerForm");
  const buyerForm = document.getElementById("buyerForm");

  // ----------------------------------------
  // REGISTRATION - FARMER
  // ----------------------------------------

  if (registerFarmerTab && registerBuyerTab && farmerForm && buyerForm) {
    registerFarmerTab.addEventListener("click", function () {
      // Show Farmer form
      farmerForm.classList.remove("hidden");
      buyerForm.classList.add("hidden");

      // Farmer active
      registerFarmerTab.classList.add(
        "bg-white",
        "text-[#0D631B]",
        "border",
        "border-[#BFCABA]",
        "shadow-sm",
        "font-semibold",
      );

      // Buyer inactive
      registerBuyerTab.classList.remove(
        "bg-white",
        "text-[#0D631B]",
        "border",
        "border-[#BFCABA]",
        "shadow-sm",
        "font-semibold",
      );

      registerBuyerTab.classList.add("text-[#40493D]", "font-medium");
    });

    // ----------------------------------------
    // REGISTRATION - BUYER
    // ----------------------------------------

    registerBuyerTab.addEventListener("click", function () {
      // Show Buyer form
      buyerForm.classList.remove("hidden");
      farmerForm.classList.add("hidden");

      // Buyer active
      registerBuyerTab.classList.add(
        "bg-white",
        "text-[#0D631B]",
        "border",
        "border-[#BFCABA]",
        "shadow-sm",
        "font-semibold",
      );

      // Farmer inactive
      registerFarmerTab.classList.remove(
        "bg-white",
        "text-[#0D631B]",
        "border",
        "border-[#BFCABA]",
        "shadow-sm",
        "font-semibold",
      );

      registerFarmerTab.classList.add("text-[#40493D]", "font-medium");
    });

    // ----------------------------------------
    // FARMER REGISTRATION SUBMIT
    // ----------------------------------------

    farmerForm.addEventListener("submit", function (event) {
      event.preventDefault();

      alert("Farmer registration UI is ready.");
    });

    // ----------------------------------------
    // BUYER REGISTRATION SUBMIT
    // ----------------------------------------

    buyerForm.addEventListener("submit", function (event) {
      event.preventDefault();

      alert("Buyer registration UI is ready.");
    });
  }
});
