<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h3>Upload aritary images - be somewhat careful maybe?</h3>

<h4>Do not start with slash - will be unceremoniously dumped into imgprssr directory</h4>
<form method="POST" action="?/upload" enctype="multipart/form-data">
  <label for="path">Path</label>
  <input name="path" />
  <input type="file" name="image" />
  <button>Upload</button>
</form>
