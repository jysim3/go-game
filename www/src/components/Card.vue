<template>
    <img :src="require(`@/assets/cards/${card}.svg`)" 
    v-if="!hidden"
    contain
    @click="onClick"
    class="elevate"
    />
</template>
<script>
export default {
  name: "CardComponent",
  props: {
    suit: [String, Number], // suit 0 --> special cards
    number: Number,
    hidden: {
      type: Boolean,
      default: false
    },
    back: {
      type: Boolean,
      default: false
    },
    onClick: {
      type: Function,
      default: () => {console.log("onClick")}
    }
  },
  computed: {
    card() {
      if (this.back) {
        return "1B";
      }
      const suits = ["S", "H", "C", "D"];
      let number = this.number
      if (number === 1) {
        number = "A"
      } else if (number > 9) {
        number = ["T", "J", "Q", "K"][number - 10];
      }
      return number + suits[this.suit];
    },
  },
};
</script>
<style scoped>
.elevate {
  box-shadow: 0px 3px 1px -2px rgb(0 0 0 / 20%), 0px 2px 2px 0px rgb(0 0 0 / 14%), 0px 1px 5px 0px rgb(0 0 0 / 12%);
  max-height: 100px;
}

</style>

