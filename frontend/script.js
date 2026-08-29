// ===============================
// KRISHISETU ROLE TAB SWITCHING
// ===============================

document.addEventListener("DOMContentLoaded", function () {
  // LOGIN PAGE
  const loginFarmerTab = document.getElementById("farmerTab");
  const loginBuyerTab = document.getElementById("buyerTab");
  const loginForm = document.getElementById("loginForm");

  // REGISTRATION PAGE
  const registerFarmerTab = document.getElementById("farmerTab");
  const registerBuyerTab = document.getElementById("buyerTab");
  const farmerForm = document.getElementById("farmerForm");
  const buyerForm = document.getElementById("buyerForm");

  // ==========================================
  // REGISTRATION PAGE
  // ==========================================

  if (farmerForm && buyerForm && registerFarmerTab && registerBuyerTab) {
    registerFarmerTab.addEventListener("click", function () {
      // Show Farmer
      farmerForm.classList.remove("hidden");
      farmerForm.classList.add("flex");

      // Hide Buyer
      buyerForm.classList.add("hidden");
      buyerForm.classList.remove("flex");

      // Farmer active
      registerFarmerTab.classList.add(
        "bg-white",
        "text-[#0D631B]",
        "border",
        "border-[#BFCABA]",
        "shadow-sm",
        "font-semibold"
      );

      // Buyer inactive
      registerBuyerTab.classList.remove(
        "bg-white",
        "text-[#0D631B]",
        "border",
        "border-[#BFCABA]",
        "shadow-sm",
        "font-semibold"
      );

      registerBuyerTab.classList.add(
        "text-[#40493D]",
        "font-medium"
      );
    });

    registerBuyerTab.addEventListener("click", function () {
      // Show Buyer
      buyerForm.classList.remove("hidden");
      buyerForm.classList.add("flex");

      // Hide Farmer
      farmerForm.classList.add("hidden");
      farmerForm.classList.remove("flex");

      // Buyer active
      registerBuyerTab.classList.add(
        "bg-white",
        "text-[#0D631B]",
        "border",
        "border-[#BFCABA]",
        "shadow-sm",
        "font-semibold"
      );

      // Farmer inactive
      registerFarmerTab.classList.remove(
        "bg-white",
        "text-[#0D631B]",
        "border",
        "border-[#BFCABA]",
        "shadow-sm",
        "font-semibold"
      );

      registerFarmerTab.classList.add(
        "text-[#40493D]",
        "font-medium"
      );
    });

    // Farmer registration
    farmerForm.addEventListener("submit", function (event) {
      event.preventDefault();
      alert("Farmer registration UI is ready.");
    });

    // Buyer registration
    buyerForm.addEventListener("submit", function (event) {
      event.preventDefault();
      alert("Buyer registration UI is ready.");
    });
  }

  // ==========================================
  // LOGIN PAGE
  // ==========================================

  if (loginFarmerTab && loginBuyerTab && loginForm) {
    loginFarmerTab.addEventListener("click", function () {
      loginFarmerTab.classList.add(
        "text-[#0D631B]",
        "font-bold",
        "border-b-2",
        "border-[#0D631B]",
        "bg-[#0D631B]/5"
      );

      loginFarmerTab.classList.remove(
        "text-[#40493D]",
        "font-medium"
      );

      loginBuyerTab.classList.remove(
        "text-[#0D631B]",
        "font-bold",
        "border-b-2",
        "border-[#0D631B]",
        "bg-[#0D631B]/5"
      );

      loginBuyerTab.classList.add(
        "text-[#40493D]",
        "font-medium"
      );
    });

    loginBuyerTab.addEventListener("click", function () {
      loginBuyerTab.classList.add(
        "text-[#0D631B]",
        "font-bold",
        "border-b-2",
        "border-[#0D631B]",
        "bg-[#0D631B]/5"
      );

      loginBuyerTab.classList.remove(
        "text-[#40493D]",
        "font-medium"
      );

      loginFarmerTab.classList.remove(
        "text-[#0D631B]",
        "font-bold",
        "border-b-2",
        "border-[#0D631B]",
        "bg-[#0D631B]/5"
      );

      loginFarmerTab.classList.add(
        "text-[#40493D]",
        "font-medium"
      );
    });

    loginForm.addEventListener("submit", function (event) {
      event.preventDefault();
      alert("Login UI is ready.");
    });
  }
});
