(function initCodeHighlight() {
  const update = () => hljs.highlightAll();

  update();
  document.body.addEventListener("htmx:load", update);
})();

/**
 * Initializes the dark mode functionality.
 * This function sets up the event listeners and initial state for dark mode based on user preference or system settings.
 */
(function initDarkMode() {
  const themeToggle = document.getElementById("theme-toggle");

  const setDarkModeStyles = (isDark) => {
    document.documentElement.classList.toggle("wa-dark", isDark);
  };

  const isDarkModeOn = () => {
    return (
      localStorage.getItem("theme") === "dark" ||
      (!localStorage.getItem("theme") &&
        window.matchMedia("(prefers-color-scheme: dark)").matches)
    );
  };

  const setDarkMode = (isDark) => {
    setDarkModeStyles(isDark);
    themeToggle.setAttribute("checked", isDark);
    localStorage.setItem("theme", isDark ? "dark" : "light");
  };

  if (isDarkModeOn()) {
    themeToggle.setAttribute("checked", true);
    setDarkModeStyles();
  }

  themeToggle.addEventListener("wa-change", () => {
    console.log("changed", themeToggle.checked);
    setDarkMode(themeToggle.checked);
  });
})();

(function initHamburgerMenu() {
  const hamburgerMenu = document.querySelector(".hamburger-menu");
  const sidebar = document.querySelector(".sidebar");
  const overlay = document.querySelector(".sidebar-overlay");
  const sidebarClose = document.querySelector(".close-sidebar");

  function toggleSidebar() {
    sidebar.classList.toggle("active");
    overlay.classList.toggle("active");
  }

  function close() {
    sidebar.classList.remove("active");
    overlay.classList.remove("active");
  }

  hamburgerMenu.addEventListener("click", toggleSidebar);
  overlay.addEventListener("click", toggleSidebar);
  sidebarClose.addEventListener("click", toggleSidebar);

  window.addEventListener("resize", () => {
    if (window.innerWidth >= 1024) {
      sidebar.classList.remove("active");
      overlay.classList.remove("active");
    }
  });

  document.body.addEventListener("htmx:load", close);
})();
