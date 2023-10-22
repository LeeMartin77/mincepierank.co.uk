<script lang="ts">
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

{#if data.editing && !data.pie}
  <h3>Pie not found</h3>
{:else if data.editing}
  <h3>Edit a pie</h3>
  <a href={`/years/${data.pie?.year}/brands/${data.pie?.makerid}/${data.pie?.id}`}>Direct link</a>
  <form method="POST" action="?/edit" enctype="multipart/form-data">
    <div>
      <label for="year">{data.pie?.year}</label>
      <input type="hidden" name="year" value={data.pie?.year} />
    </div>
    <div>
      <label for="makerid">{data.pie?.makerid}</label>
      <input type="hidden" name="makerid" value={data.pie?.makerid} />
    </div>
    <div>
      <label for="id">{data.pie?.id}</label>
      <input type="hidden" name="id" value={data.pie?.id} />
    </div>
    <div>
      <label for="displayname">displayname</label>
      <input name="displayname" value={data.pie?.displayname} />
    </div>
    <div>
      <label for="fresh">fresh</label>
      <input type="checkbox" name="fresh" checked={data.pie?.fresh} />
    </div>
    <div>
      <label for="labels">labels</label>
      <input name="labels" value={data.pie?.labels.join(',')} />
    </div>
    <div>
      <label for="web_link">web_link</label>
      <input name="web_link" value={data.pie?.web_link} />
    </div>
    <div>
      <label for="pack_count">pack_count</label>
      <input type="number" name="pack_count" value={data.pie?.pack_count} />
    </div>
    <div>
      <label for="pack_price_in_pence">pack_price_in_pence</label>
      <input type="number" name="pack_price_in_pence" value={data.pie?.pack_price_in_pence} />
    </div>
    <div>
      <img
        src={imgprssrPrefix(data.pie?.image_file + '?filter=gaussian&width=250', data.imgprssr)}
        alt={data.pie?.displayname}
      />
    </div>
    <div>
      <input type="file" name="image" />
    </div>
    <div>
      <label for="validated">validated</label>
      <input type="checkbox" name="validated" checked={data.pie?.validated} />
    </div>
    <button>Edit</button>
  </form>
  <h2>Delete this?</h2>
  <form method="POST" action="?/delete" enctype="multipart/form-data">
    <input type="hidden" name="year" value={data.pie?.year} />
    <input type="hidden" name="makerid" value={data.pie?.makerid} />
    <input type="hidden" name="id" value={data.pie?.id} />
    <button>Delete</button>
  </form>
{:else}
  <h3>Slam a pie onto cassandra</h3>

  <h4>You need to fill in all fields AND provide an image</h4>
  <form method="POST" action="?/upload" enctype="multipart/form-data">
    <div>
      <label for="year">Year</label>
      <input name="year" value={data.activeYear} type="number" />
    </div>
    <div>
      <label for="makerid">makerid</label>
      <select name="makerid">
        {#each data.makers as maker}
          <option value={maker.id}>{maker.name}</option>
        {/each}
      </select>
    </div>
    <div>
      <label for="id">id</label>
      <input name="id" />
    </div>
    <div>
      <label for="displayname">displayname</label>
      <input name="displayname" />
    </div>
    <div>
      <label for="fresh">fresh</label>
      <input type="checkbox" name="fresh" />
    </div>
    <div>
      <label for="labels">labels</label>
      <input name="labels" />
    </div>
    <div>
      <label for="web_link">web_link</label>
      <input name="web_link" />
    </div>
    <div>
      <label for="pack_count">pack_count</label>
      <input type="number" name="pack_count" />
    </div>
    <div>
      <label for="pack_price_in_pence">pack_price_in_pence</label>
      <input type="number" name="pack_price_in_pence" />
    </div>
    <div>
      <input type="file" name="image" />
    </div>
    <div>
      <label for="validated">validated</label>
      <input type="checkbox" name="validated" />
    </div>
    <button>Upload</button>
  </form>
{/if}

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 0.5em;
  }
  form div {
    display: flex;
  }
  form label {
    margin: 0;
  }
  form input,
  form select {
    border: 1px solid black;
  }
  button {
    padding: 1em;
  }
</style>
