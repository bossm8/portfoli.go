document.addEventListener("DOMContentLoaded", function () {
  var toggle = document.getElementById("nav-toggle");
  var links = document.getElementById("nav-links");
  var icon = document.getElementById("nav-toggle-icon");
  if (!toggle || !links) return;

  toggle.addEventListener("click", function () {
    var isOpen = links.classList.toggle("is-open");
    toggle.setAttribute("aria-expanded", isOpen ? "true" : "false");
    if (icon) {
      icon.className = isOpen ? "bi-x-lg" : "bi-list";
    }
  });
});
