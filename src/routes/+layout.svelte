<script lang="ts">
  /* eslint-disable @typescript-eslint/ban-ts-comment */
  import Breadcrumb from '$components/Breadcrumb.svelte';
  import { onMount } from 'svelte';
  import type { LayoutData } from './$types';
  import Menu from '$components/Menu.svelte';
  export let data: LayoutData;

  let noticeVisible = false;
  let noticestring = `${data.readonly ? 'READONLY' : ''}:${data.notice ? data.notice : ''}`;

  let cookieNotice = false;

  const noticekey = 'noticestring';

  const cookiekey = 'cookiekey';

  onMount(() => {
    if (data.session?.user?.email) {
      // need to look into some form of event for this
      fetch(`/api/claim-anonymous`, { method: 'POST' });
    }

    noticeVisible = data.readonly || !!data.notice;

    let lsstr = localStorage.getItem(noticekey);
    if (lsstr && lsstr === noticestring) {
      noticeVisible = false;
    }

    cookieNotice = !window.location.href.endsWith('/about/cookies');
    if (localStorage.getItem(cookiekey) === 'accepted') {
      cookieNotice = false;
    }

    // snow is a modified version of https://stackoverflow.com/a/13983965
    (function () {
      var requestAnimationFrame =
        window.requestAnimationFrame ||
        function (callback: any) {
          window.setTimeout(callback, 1000 / 60);
        };
      window.requestAnimationFrame = requestAnimationFrame;
    })();

    var flakes: any[] = [],
      canvas = document.getElementById('canvas'),
      // @ts-ignore
      ctx = canvas.getContext('2d'),
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

        ctx.fillStyle = 'rgba(255,255,255,' + flake.opacity + ')';
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

    function reset(flake: any) {
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
          opacity: opacity
        });
      }

      snow();
    }

    init();

    addEventListener('resize', () => {
      canvas = document.getElementById('canvas');
      // @ts-ignore
      ctx = canvas.getContext('2d');
      // @ts-ignore
      canvas.width = window.innerWidth;
      // @ts-ignore
      canvas.height = window.innerHeight;
    });
  });
</script>

<canvas id="canvas" />
<div class="page-container">
  <Menu activeYear={data.activeYear} customPiesEnabled={data.customPiesEnabled} />
  <div class="content-container">
    {#if cookieNotice}
      <div class="cookie-notice">
        <div>
          <p>We use cookies for authentication and associating mince pie rankings to a user</p>
          <button
            on:click={() => {
              localStorage.setItem(cookiekey, 'accepted');
              cookieNotice = false;
            }}>Accept</button
          >
          <button
            on:click={() => {
              cookieNotice = false;
            }}>Dismiss</button
          >
          <a href="/about/cookies">Read Full Policy</a>
        </div>
      </div>
    {/if}
    {#if noticeVisible}
      <div class="notices">
        <button
          on:click={() => {
            localStorage.setItem(noticekey, noticestring);
            noticeVisible = false;
          }}>X</button
        >
        <div>
          {#if data.readonly}
            <p>Mince Pie Rank is read-only until the next Christmas period</p>
          {/if}
          {#if data.notice}
            <p>{data.notice}</p>
          {/if}
        </div>
      </div>
    {/if}
    <Breadcrumb />
    <div class="content">
      <slot />
    </div>
    <a style="padding-bottom: 6em" href="/about">About</a>
  </div>
</div>

<style>
  .page-container {
    display: flex;
    flex-direction: row;
    top: 0;
    left: 0;
  }
  .content-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-left: auto;
    margin-right: auto;
    padding: 1em;
    overflow-y: auto;
    max-height: calc(100dvh - 2em);
  }
  .notices {
    text-align: center;
    border-radius: 1em;
    position: fixed;
    width: calc(100vw - 4em);
    box-shadow: 0.2em 0.1em 0.5em rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    top: 1em;
    left: 1em;
    display: flex;
    padding: 1em;
    align-items: center;
    gap: 1em;
    z-index: 9999999999998;
  }

  .notices button {
    height: 2em;
    width: 2em;
    background: none;
    border: 1px solid black;
    border-radius: 0.5em;
  }

  .cookie-notice {
    text-align: center;
    border-radius: 1em;
    position: fixed;
    width: calc(100vw - 4em);
    box-shadow: 0.2em 0.1em 0.5em rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    bottom: 1em;
    left: 1em;
    display: flex;
    flex-direction: column;
    padding: 1em;
    align-items: center;
    gap: 1em;
    z-index: 9999999999999;
  }

  .cookie-notice button {
    display: block;
    width: 100%;
    cursor: pointer;
    border-radius: 1em;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.2);
    background: white;
    border: none;
    padding: 1em;
    text-align: center;
    font-weight: bold;
    color: #111;
    margin-bottom: 1em;
  }

  .cookie-notice a {
    display: block;
    border-radius: 1em;
    padding: 0.5em;
    text-align: center;
    font-weight: bold;
    color: #111;
  }
  .cookie-notice a:visited {
    color: #111;
  }

  .notices button:hover {
    cursor: pointer;
  }

  .content {
    padding-bottom: 6em;
  }

  @media only screen and (max-width: 640px) {
    .content {
      min-width: 320px;
      max-width: 460px;
    }
  }
  @media only screen and (min-width: 640px) and (max-width: 960px) {
    .content {
      min-width: 460px;
      max-width: 960px;
    }
  }

  @media only screen and (min-width: 1020px) {
    .content {
      min-width: 840px;
      max-width: 1020px;
    }
  }

  canvas {
    position: fixed;
    z-index: -9999999;
    background-image: linear-gradient(bottom, #c0e1ff 0%, #8ec8ff 84%);

    background-image: -webkit-gradient(
      linear,
      left bottom,
      left top,
      color-stop(0, #c0e1ff),
      color-stop(0.84, #8ec8ff)
    );
  }
</style>
