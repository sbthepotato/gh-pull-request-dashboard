<script>
  export let name;
  export let color;
  export let width = "16px";
  export let height = "16px";
</script>

<img
  {height}
  {width}
  src="../../src/assets/icons/{name}.svg"
  alt={name}
  class={color}
/>

<style>
  .green {
    filter: brightness(0) saturate(100%) invert(60%) sepia(67%) saturate(442%)
      hue-rotate(77deg) brightness(87%) contrast(93%);
  }

  .red {
    filter: brightness(0) saturate(100%) invert(40%) sepia(50%) saturate(711%)
      hue-rotate(314deg) brightness(105%) contrast(119%);
  }

  .yellow {
    filter: invert(85%) sepia(34%) saturate(4575%) hue-rotate(348deg)
      brightness(89%) contrast(83%);
  }

  .grey {
    filter: brightness(0) saturate(100%) invert(64%) sepia(10%) saturate(278%)
      hue-rotate(174deg) brightness(92%) contrast(87%);
  }

  .white {
    filter: invert(99%) sepia(1%) saturate(1998%) hue-rotate(66deg)
      brightness(117%) contrast(100%);
  }
</style>
