<script lang="ts">
  import Breadcrumb from '$components/Breadcrumb.svelte';
</script>

<Breadcrumb />

<h1>Years</h1>

<a class="year-link" href="/years/2022">2022</a>

<style>
  h1 {
    text-align: center;
  }
  .year-link {
    display: block;
    padding: 1em;
    color: initial;
    text-align: center;
    text-decoration: none;
    font-weight: bold;
    font-size: 2em;
    padding: 0.5em;
    border: 1px solid black;
    border-radius: 1em;
  }
  .year-link:visited {
    color: initial;
  }
</style>
