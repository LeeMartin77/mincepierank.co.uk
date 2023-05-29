<script lang="ts">
  import Breadcrumb from '$components/Breadcrumb.svelte';
</script>

<div class="notices">
  <p>Mince Pie Rank is read-only until the next Christmas period</p>
  <p>Things are a bit sparse right now - getting ready for the new season!</p>
</div>

<Breadcrumb />

<div class="content-container">
  <div class="content">
    <slot />
  </div>
</div>

<style>
  .content-container {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .notices {
    text-align: center;
    border: 1px dashed grey;
    border-radius: 1em;
  }

  @media only screen and (max-width: 640px) {
    .content {
      max-width: 460px;
    }
  }
  @media only screen and (min-width: 640px) and (max-width: 960px) {
    .content {
      min-width: 640px;
      max-width: 960px;
    }
  }

  @media only screen and (min-width: 960px) {
    .content {
      min-width: 960px;
      max-width: 1020px;
    }
  }
</style>
