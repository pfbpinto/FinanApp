function togglePasswordVisibility(inputId, toggleBtnId, iconId) {
  const passwordInput = document.getElementById(inputId);
  const toggleButton = document.getElementById(toggleBtnId);
  const toggleIcon = document.getElementById(iconId);

  toggleButton.addEventListener("click", () => {
    const type =
      passwordInput.getAttribute("type") === "password" ? "text" : "password";
    passwordInput.setAttribute("type", type);

    // Alternar o ícone
    if (type === "password") {
      toggleIcon.classList.remove("fa-eye-slash");
      toggleIcon.classList.add("fa-eye");
    } else {
      toggleIcon.classList.remove("fa-eye");
      toggleIcon.classList.add("fa-eye-slash");
    }
  });
}
