<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<a class="year-link" href="/profile/rankings/2022">2022</a>

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
</style>
