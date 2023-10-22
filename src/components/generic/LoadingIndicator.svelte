<script lang="ts">
  export let width = 48;
</script>

<div style="display: flex; flex-direction: column; width: 100%; align-items:center;">
  <img alt="ranking icon" src="/rankicon.svg" {width} class="spinny" />
</div>

<style>
  .spinny {
    animation: rotation 4s infinite linear;
  }
  @keyframes rotation {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
