<script lang="ts">
  //import type { PageData } from './$types';

  //export let data: PageData;
</script>

<h3>Slam a pie onto cassandra</h3>

<h4>You need to fill in all fields AND provide an image</h4>
<form method="POST" action="?/upload" enctype="multipart/form-data">
  <div>
    <label for="year">Year</label>
    <input name="year" value="2023"/>
  </div>
  <div>
    <label for="makerid">makerid</label>
    <input name="makerid" />
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
  <input type="file" name="image" />
  <button>Upload</button>
</form>
