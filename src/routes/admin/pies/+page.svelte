<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<a href={'/admin/pies/edit'}>Add a new pie</a>
<h4>Looking at {data.year}</h4>

{#each data.years as year}
  <a style="padding: 0.5em" href={'/admin/pies?year=' + year}>{year}</a>
{/each}

<h5>Unvalidated Pies</h5>
<table>
  <thead>
    <tr>
      <th>Maker Id</th>
      <th>Display Name</th>
      <th>Edit</th>
    </tr>
  </thead>
  <tbody>
    {#each data.unvalidatedPies as pie}
      <tr>
        <td>{pie.makerid}</td>
        <td>{pie.displayname}</td>
        <td
          ><a href={`/admin/pies/edit?year=${+pie.year}&makerid=${pie.makerid}&id=${pie.id}`}
            >Edit</a
          ></td
        >
      </tr>
    {/each}
  </tbody>
</table>

<h5>Validated Pies</h5>
<table>
  <thead>
    <tr>
      <th>Maker Id</th>
      <th>Display Name</th>
      <th>Edit</th>
    </tr>
  </thead>
  <tbody>
    {#each data.validatedPies as pie}
      <tr>
        <td>{pie.makerid}</td>
        <td>{pie.displayname}</td>
        <td
          ><a href={`/admin/pies/edit?year=${+pie.year}&makerid=${pie.makerid}&id=${pie.id}`}
            >Edit</a
          ></td
        >
      </tr>
    {/each}
  </tbody>
</table>

<style>
  th,
  td {
    border: grey solid 1px;
    padding: 0.5em;
    margin: 0;
  }
</style>
