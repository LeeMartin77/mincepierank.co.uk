<script lang="ts">
  import type { MakerPie } from '$lib/storage';
  import { ppCategory } from './utilities/formatCategory';
  import type { PieListRanking } from './utilities/mapPiesAndRankings';
  export let pie: MakerPie;
  export let pieListRanking: PieListRanking | undefined = undefined;
</script>

<div class="pie-card">
  <div class="pie-card-top">
    <div style="width: 30%">
      <img
        width="100%"
        src={`${pie.image_file}?width=250&height=250&filter=gaussian`}
        alt={pie.displayname}
      />
    </div>
    <div>
      <dl>
        <dt>Average</dt>
        <dd>{(pieListRanking?.average || 0).toFixed(2)}</dd>
        <dl>
          <dt>Votes</dt>
          <dd>{pieListRanking?.count || 0}</dd>
        </dl>
      </dl>
    </div>
  </div>
  <div class="category-links">
    {#each pie.labels as label}
      <a class="category-link" href={`/years/${pie.year}/categories/${label}`}
        >{ppCategory(label)}</a
      >
    {/each}
  </div>
  <a class="pie-link" href={`/years/${pie.year}/brands/${pie.makerid}/${pie.id}`}
    >{pie.displayname}</a
  >
</div>

<style>
  .pie-card {
    padding: 1em;
    border: 1px solid #111;
    border-radius: 1em;
    display: flex;
    flex-direction: column;
  }
  .pie-card-top {
    display: flex;
    gap: 2em;
    justify-content: center;
  }
  .pie-link {
    display: block;
    border: 1px solid #111;
    border-radius: 1em;
    padding: 1em;
    text-decoration: none;
    text-align: center;
    font-weight: bold;
    color: #111;
  }
  .pie-link:visited {
    color: #111;
  }

  .category-links {
    display: flex;
    justify-content: center;
    gap: 0.5em;
    margin-top: 0.5em;
    margin-bottom: 0.5em;
  }

  .category-link {
    display: block;
    border: 1px solid #111;
    font-size: 0.8em;
    border-radius: 2em;
    padding: 0.5em;
    text-decoration: none;
    text-align: center;
    color: initial;
  }
  .category-link:visited {
    color: initial;
  }

  @media (prefers-color-scheme: dark) {
    .pie-card {
      border: 1px solid #fff;
    }
    .pie-link {
      border: 1px solid #fff;
      color: #fff;
    }
    .pie-link:visited {
      border: 1px solid #fff;
      color: #fff;
    }
    .category-link {
      border: 1px solid #fff;
      color: #fff;
    }
    .category-link:visited {
      border: 1px solid #fff;
      color: #fff;
    }
  }
</style>
