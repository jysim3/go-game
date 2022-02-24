<template>
  <div>
    <hello-world />
    <v-container>
      <v-row justify="center">
        <v-btn color="primary">Start a room</v-btn>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import HelloWorld from "../components/HelloWorld";

export default {
  name: "HomeView",

  components: {
    HelloWorld,
  },
};
</script>
