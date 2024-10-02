<script>
  export let name;
  export let color;
  export let width = "16px";
  export let height = "16px";
</script>

<img
  {height}
  {width}
  src="/icons/{name}.svg"
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

  .rainbow {
    filter: invert(100%) sepia(0%);
    animation: 32s rainbowz linear infinite;
  }

  @keyframes rainbowz {
    0% {
      filter: invert(18%) sepia(82%) saturate(6300%) hue-rotate(356deg)
        brightness(106%) contrast(124%);
    }
    10% {
      filter: invert(58%) sepia(26%) saturate(2593%) hue-rotate(5deg)
        brightness(110%) contrast(103%);
    }
    20% {
      filter: invert(84%) sepia(81%) saturate(3918%) hue-rotate(354deg)
        brightness(105%) contrast(106%);
    }
    30% {
      filter: invert(72%) sepia(20%) saturate(3443%) hue-rotate(43deg)
        brightness(114%) contrast(116%);
    }
    40% {
      filter: invert(28%) sepia(84%) saturate(1461%) hue-rotate(92deg)
        brightness(95%) contrast(107%);
    }
    50% {
      filter: invert(100%) sepia(92%) saturate(3697%) hue-rotate(102deg)
        brightness(104%) contrast(103%);
    }
    60% {
      filter: invert(9%) sepia(100%) saturate(6501%) hue-rotate(248deg)
        brightness(91%) contrast(146%);
    }
    70% {
      filter: invert(15%) sepia(43%) saturate(4774%) hue-rotate(266deg)
        brightness(86%) contrast(129%);
    }
    80% {
      filter: invert(75%) sepia(35%) saturate(2086%) hue-rotate(244deg)
        brightness(98%) contrast(90%);
    }
    90% {
      filter: invert(19%) sepia(44%) saturate(7411%) hue-rotate(294deg)
        brightness(115%) contrast(126%);
    }
    100% {
      filter: invert(25%) sepia(66%) saturate(7356%) hue-rotate(293deg)
        brightness(112%) contrast(126%);
    }
  }
</style>
