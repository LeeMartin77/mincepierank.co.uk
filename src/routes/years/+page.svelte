<script lang="ts">
  import type { PageData } from './$types';
  export let data: PageData;
  const years = data.years.sort().reverse();
</script>

<h1>Years</h1>

{#each years as year}
  <a class="year-link" href={`/years/${year}`}>{year}</a>
{/each}

<style>
  h1 {
    text-align: center;
  }
  .year-link {
    display: block;
    padding: 1em;
    color: #111;
    text-align: center;
    text-decoration: none;
    font-weight: bold;
    font-size: 2em;
    padding: 0.5em;
    border: 1px solid black;
    border-radius: 1em;
  }
  .year-link:visited {
    color: #111;
  }
</style>
