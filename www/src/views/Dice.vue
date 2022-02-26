<template>
  <div>
    <v-container>
      <v-row justify="center" class="text-center">
        <v-col cols="auto">
          <h1 class="text-h1">Dice - room <code>{{ id }}</code></h1>
        </v-col>
      </v-row>
      <v-row
        justify="center"
        v-for="(dice, index) in diceFormatted"
        :key="index"
      >
        <v-col cols="2" v-for="(die, index) in dice" :key="index">
          <v-img :src="require(`@/assets/dice/${die}.png`)" class="rounded-lg" max-width="150" />
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="auto">
          <v-btn color="primary" @click="send">Start</v-btn>
        </v-col>
        <v-col cols="auto">
          <v-dialog v-model="dialog" width="500">
            <template v-slot:activator="{ on, attrs }">
              <v-btn color="red lighten-2" dark v-bind="attrs" v-on="on"
                >Open</v-btn
              >
            </template>

            <v-card>
              <v-card-title color="danger"> Are you sure? </v-card-title>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="primary"
                  text
                  @click="
                    open();
                    dialog = false;
                  "
                >
                  Yes
                </v-btn>
                <v-btn color="primary" text @click="dialog = false"> No </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-col>
        <v-col cols="auto">
          <v-btn color="primary" @click="reset">Reset</v-btn>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
// import HelloWorld from "../components/HelloWorld";

export default {
  name: "DiceView",
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      msg: "xxx",
      dice: [1, 2, 3, 4, 5],
      ws: null,
      dialog: false,
    };
  },
  computed: {
    diceFormatted() {
      return this.dice.reduce((resultArray, item, index) => {
        const chunkIndex = Math.floor(index / 5);
        if (!resultArray[chunkIndex]) {
          resultArray[chunkIndex] = []; // start a new chunk
        }
        resultArray[chunkIndex].push(item);

        return resultArray;
      }, []);
    },
  },
  mounted() {
    const ws = new WebSocket(`ws://${window.location.host}/dice/${this.id}/ws`);
    this.ws = ws;
    this.msg = "hellooo";
    var self = this;
    ws.onmessage = function (msg) {
      const j = JSON.parse(msg.data);
      if (j.command == "start") {
        self.dice = j.data;
      }
    };
  },
  methods: {
    open() {
      this.ws.send(
        JSON.stringify({
          command: "open",
          data: "",
        })
      );
    },
    reset() {
      this.ws.send(
        JSON.stringify({
          command: "reset",
          data: "",
        })
      );
    },
    send() {
      this.ws.send(
        JSON.stringify({
          command: "start",
          data: "",
        })
      );
    },
  },
};
</script>
