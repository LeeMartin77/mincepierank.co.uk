<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;
</script>

{#each data.years as year}
  <a class="year-link" href="/profile/pies/{year}">{year}</a>
{/each}

<style>
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

  @media (prefers-color-scheme: dark) {
    .year-link {
      color: #fff;
      border: 1px solid #fff;
    }
    .year-link:visited {
      color: #ddd;
    }
  }
</style>
