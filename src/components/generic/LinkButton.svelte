<script lang="ts">
  export let href: string;
</script>

<a class="link-button" {href}><slot /></a>

<style>
  .link-button {
    display: block;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.2);
    border-radius: 1em;
    padding: 1em;
    text-decoration: none;
    text-align: center;
    font-weight: bold;
    color: #111;
  }
  .link-button:visited {
    color: #111;
  }
</style>
