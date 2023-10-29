<script lang="ts">
  export let href: string;
  export let style: string | undefined = undefined;
</script>

<a class="link-button" {style} {href}><slot /></a>

<style>
  .link-button {
    display: block;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.2);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 1em;
    padding: 1em;
    text-align: center;
    font-weight: bold;
    color: #111;
  }
  .link-button:visited {
    color: #111;
  }
</style>
