(function () {
  var STORAGE_KEY = "theme";
  var root = document.documentElement;

  function effectiveTheme() {
    return (
      root.getAttribute("data-theme") ||
      (window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light")
    );
  }

  function updateToggle(theme) {
    var btn = document.getElementById("theme-toggle");
    var icon = document.getElementById("theme-toggle-icon");
    if (!btn || !icon) return;
    icon.className = theme === "dark" ? "bi-sun" : "bi-moon-stars";
    btn.setAttribute("aria-pressed", theme === "dark" ? "true" : "false");
    btn.setAttribute(
      "aria-label",
      theme === "dark" ? "Switch to light theme" : "Switch to dark theme"
    );
  }

  document.addEventListener("DOMContentLoaded", function () {
    updateToggle(effectiveTheme());

    var btn = document.getElementById("theme-toggle");
    if (!btn) return;

    btn.addEventListener("click", function () {
      var next = effectiveTheme() === "dark" ? "light" : "dark";
      root.setAttribute("data-theme", next);
      try {
        localStorage.setItem(STORAGE_KEY, next);
      } catch (e) {
        /* localStorage unavailable (e.g. private browsing) - theme just won't persist */
      }
      updateToggle(next);
    });
  });
})();
