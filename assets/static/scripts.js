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

// Snow is a modified version of https://stackoverflow.com/a/13983965
//
// The canvas element only exists on public pages while the year is active,
// and hx-boost swaps the whole <body> (re-creating the canvas and re-running
// this script), so: init happens once, the loop is single and stops itself
// when no canvas is present, and the canvas is re-looked-up and re-sized on
// swap and resize.
(function () {
  if (window.matchMedia &&
      window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
    return; // honour reduced motion: gradient backdrop, no falling flakes
  }

  var canvas = null,
    ctx = null,
    flakes = [],
    flakeCount = 50,
    running = false;

  function rebind() {
    canvas = document.getElementById("canvas");
    if (!canvas) {
      ctx = null;
      return false;
    }
    ctx = canvas.getContext("2d");
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    return true;
  }

  function reset(flake) {
    flake.x = Math.floor(Math.random() * canvas.width);
    flake.y = 0;
    flake.size = Math.random() * 3 + 2;
    flake.speed = Math.random() * 1 + 0.5;
    flake.velY = flake.speed;
    flake.velX = 0;
    flake.opacity = Math.random() * 0.5 + 0.3;
  }

  function seed() {
    for (var i = 0; i < flakeCount; i++) {
      var speed = Math.random() * 1 + 0.5;
      flakes.push({
        x: Math.floor(Math.random() * canvas.width),
        y: Math.floor(Math.random() * canvas.height),
        size: Math.random() * 3 + 2,
        speed: speed,
        velY: speed,
        velX: 0,
        opacity: Math.random() * 0.5 + 0.3,
      });
    }
  }

  function snow() {
    if (!ctx) {
      running = false; // snow switched off (or page swapped away); resume on next swap
      return;
    }

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    for (var i = 0; i < flakes.length; i++) {
      var flake = flakes[i];

      ctx.fillStyle = "rgba(255,255,255," + flake.opacity + ")";
      flake.y += flake.velY;
      flake.x += flake.velX;

      if (flake.y >= canvas.height || flake.y <= 0) {
        reset(flake);
      }

      if (flake.x >= canvas.width || flake.x <= 0) {
        reset(flake);
      }

      ctx.beginPath();
      ctx.arc(flake.x, flake.y, flake.size, 0, Math.PI * 2);
      ctx.fill();
    }

    requestAnimationFrame(snow);
  }

  function start() {
    if (running || !rebind()) {
      return;
    }
    if (flakes.length === 0) {
      seed();
    }
    running = true;
    requestAnimationFrame(snow);
  }

  if (!window.mprSnowInit) {
    window.mprSnowInit = true;

    document.addEventListener("htmx:afterSwap", start);
    window.addEventListener("resize", function () {
      rebind(); // in-place resize: keep flakes, just re-fit the drawing surface
    });
  }

  start();
})();
