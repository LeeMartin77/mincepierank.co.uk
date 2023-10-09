<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  export let name: string;
  export let displayName: string;
  export let value: number | undefined;

  const dispatch = createEventDispatcher();
  const change = () => dispatch('change');
</script>

<label for={name}>{displayName}</label>
<input {name} type="number" min="1" max="5" bind:value on:change={() => change()} />
