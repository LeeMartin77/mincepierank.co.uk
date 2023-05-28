<script lang="ts">
  import Breadcrumb from '$components/Breadcrumb.svelte';
</script>

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