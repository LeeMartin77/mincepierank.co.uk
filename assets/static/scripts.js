// Mince Pie Rank public-page navigation.
// Below the desktop breakpoint the drawer is off-canvas and opens from the
// bottom pill; at >= 992px (Bootstrap lg) CSS pins it as a permanent left rail
// and hides the pill, so open/closed state only has an effect on mobile.
//
// hx-boost re-renders the whole <body> (including this fragment and possibly
// this script) on every navigation, so: state lives on window, document-level
// listeners are attached once, and the menu element is looked up fresh by id
// each time.
(function () {
  function applyMenuState() {
    var menu = document.getElementById("mpr-menu");
    if (!menu) {
      return; // admin pages (and non-nav pages) have no drawer
    }
    var open = !!window.mprMenuOpen;
    menu.classList.toggle("is-open", open);
    document.body.classList.toggle("mpr-menu-open", open);
    var pill = document.getElementById("mpr-menu-pill");
    if (pill) {
      pill.setAttribute("aria-expanded", open ? "true" : "false");
    }
  }

  function setMenu(open, focusPill) {
    window.mprMenuOpen = open;
    applyMenuState();
    if (focusPill) {
      var pill = document.getElementById("mpr-menu-pill");
      if (pill) {
        pill.focus();
      }
    }
  }

  window.mprApplyMenuState = applyMenuState;

  if (!window.mprMenuInit) {
    window.mprMenuInit = true;

    document.addEventListener("click", function (event) {
      if (event.target.closest("#mpr-menu-pill")) {
        setMenu(!window.mprMenuOpen, true);
      } else if (event.target.closest("#mpr-menu .mpr-menu-close")) {
        setMenu(false, true);
      } else if (event.target.closest("#mpr-menu a")) {
        // navigating away; the drawer must not stay open behind the new page
        setMenu(false);
      }
    });

    document.addEventListener("keydown", function (event) {
      if (event.key === "Escape" && window.mprMenuOpen) {
        setMenu(false, true);
      }
    });

    document.addEventListener("htmx:afterSwap", function () {
      applyMenuState();
    });
  }

  applyMenuState();
})();
