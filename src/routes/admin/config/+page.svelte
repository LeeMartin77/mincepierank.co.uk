<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h3>Only do one thing at a time - this is not beautiful</h3>

<table>
  <thead>
    <th> Key </th>
    <th> Value </th>
  </thead>
  <tbody>
    {#each Object.entries(data.config) as [key, value]}
      <tr>
        <td>{key}</td>
        <td
          ><form method="POST" action="?/set">
            <input name="key" type="hidden" value={key} /><input name="value" {value} /><button
              >Save</button
            >
          </form></td
        >
        <td
          ><form method="POST" action="?/delete">
            <input name="key" type="hidden" value={key} /><button>Delete</button>
          </form></td
        >
      </tr>
    {/each}
  </tbody>
</table>

<h3>Add a new thing</h3>
<form method="POST" action="?/set">
  <label for="key">Key</label><input name="key" />
  <label for="value">Value</label><input name="value" />
  <button>Add</button>
</form>
