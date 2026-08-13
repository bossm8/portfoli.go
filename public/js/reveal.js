(function () {
  // The .js-reveal marker is added synchronously in <head> (before first
  // paint) so sections never render visible and then flash to hidden here.
  // If it's missing, reduced motion is on, or the head script didn't run
  // (JS partially unavailable) - leave content visible, do nothing.
  if (!document.documentElement.classList.contains("js-reveal")) {
    return;
  }

  document.addEventListener("DOMContentLoaded", function () {
    var els = document.querySelectorAll(".reveal");
    if (!els.length) return;

    if (!("IntersectionObserver" in window)) {
      els.forEach(function (el) {
        el.classList.add("is-visible");
      });
      return;
    }

    // A small negative bottom rootMargin means a section only reveals once
    // it's genuinely on screen (not hundreds of pixels early, which made the
    // fade finish before it was ever seen - and not so late that it needs a
    // large fraction already scrolled past), so the transition plays out as
    // part of the scroll motion instead of being invisible either way.
    var observer = new IntersectionObserver(
      function (entries) {
        entries.forEach(function (entry) {
          if (entry.isIntersecting) {
            entry.target.classList.add("is-visible");
            observer.unobserve(entry.target);
          }
        });
      },
      { threshold: 0, rootMargin: "0px 0px -80px 0px" }
    );

    els.forEach(function (el) {
      observer.observe(el);
    });
  });
})();
