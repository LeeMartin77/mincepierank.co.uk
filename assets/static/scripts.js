// snow is a modified version of https://stackoverflow.com/a/13983965
(function () {
  var requestAnimationFrame =
    window.requestAnimationFrame ||
    function (callback) {
      window.setTimeout(callback, 1000 / 60);
    };
  window.requestAnimationFrame = requestAnimationFrame;
})();

var flakes = [],
  canvas = document.getElementById("canvas"),
  // @ts-ignore
  ctx = canvas.getContext("2d"),
  flakeCount = 50,
  mX = -100,
  mY = -100;
// @ts-ignore

canvas.width = window.innerWidth;
// @ts-ignore
canvas.height = window.innerHeight;

function snow() {
  // @ts-ignore

  ctx.clearRect(0, 0, canvas.width, canvas.height);

  for (var i = 0; i < flakeCount; i++) {
    var flake = flakes[i];

    ctx.fillStyle = "rgba(255,255,255," + flake.opacity + ")";
    flake.y += flake.velY;
    flake.x += flake.velX;
    // @ts-ignore

    if (flake.y >= canvas.height || flake.y <= 0) {
      reset(flake);
    }

    // @ts-ignore

    if (flake.x >= canvas.width || flake.x <= 0) {
      reset(flake);
    }

    ctx.beginPath();
    ctx.arc(flake.x, flake.y, flake.size, 0, Math.PI * 2);
    ctx.fill();
  }
  requestAnimationFrame(snow);
}

function reset(flake) {
  // @ts-ignore
  flake.x = Math.floor(Math.random() * canvas.width);
  flake.y = 0;
  flake.size = Math.random() * 3 + 2;
  flake.speed = Math.random() * 1 + 0.5;
  flake.velY = flake.speed;
  flake.velX = 0;
  flake.opacity = Math.random() * 0.5 + 0.3;
}

function init() {
  for (var i = 0; i < flakeCount; i++) {
    // @ts-ignore

    var x = Math.floor(Math.random() * canvas.width),
      // @ts-ignore

      y = Math.floor(Math.random() * canvas.height),
      size = Math.random() * 3 + 2,
      speed = Math.random() * 1 + 0.5,
      opacity = Math.random() * 0.5 + 0.3;

    flakes.push({
      speed: speed,
      velY: speed,
      velX: 0,
      x: x,
      y: y,
      size: size,
      stepSize: Math.random() / 30,
      step: 0,
      angle: 180,
      opacity: opacity,
    });
  }

  snow();
}

init();

addEventListener("resize", () => {
  canvas = document.getElementById("canvas");
  // @ts-ignore
  ctx = canvas.getContext("2d");
  // @ts-ignore
  canvas.width = window.innerWidth;
  // @ts-ignore
  canvas.height = window.innerHeight;
});

menuOpen = false;
