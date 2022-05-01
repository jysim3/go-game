<template>
  <div>
    <v-container>
      <v-row justify="center" class="text-center">
        <v-col cols="auto">
          <h1 class="text-h1">
            Dice - room <code>{{ id }}</code>
          </h1>
        </v-col>
      </v-row>
      <v-row v-if="dice.length > 5">
        <v-col cols="12">
          {{ this.diceResult }}
        </v-col>
        <v-col cols="12" class="d-flex justify-center">
          <v-switch v-model="withOnes" label="With Joker(1)" />
        </v-col>
        <v-col
          cols="12"
          v-for="die in [1, 2, 3, 4, 5, 6]"
          :key="die"
          class="d-flex justify-center"
        >
          <v-img
            :src="require(`@/assets/dice/${die}.png`)"
            class="rounded-lg"
            max-width="30"
            :gradient="
              highlight === die
                ? 'to top right, rgba(255,0,0,.33), rgba(56,0,0,.7)'
                : ''
            "
            @click="highlight = die"
          />
          <span class="text-h6 ml-5">{{ diceResult[die] || "0" }}</span>
        </v-col>
      </v-row>
      <div v-if="!hide">
        <v-row
          justify="center"
          v-for="(dice, index) in diceFormatted"
          :key="index"
        >
          <v-col cols="2" v-for="(die, index) in dice" :key="index">
            <v-img
              @click="highlight = die"
              :gradient="
                highlight === die || (die === 1 && withOnes)
                  ? 'to top right, rgba(255,0,0,.33), rgba(180,0,0,.7)'
                  : ''
              "
              :src="require(`@/assets/dice/${die}.png`)"
              class="rounded-lg"
              max-width="150"
            />
          </v-col>
        </v-row>
      </div>
      <v-row justify="center">
        <v-col cols="auto">
          <v-btn color="primary" @click.stop="rerollButton">Reroll</v-btn>
        </v-col>
        <v-col cols="auto">
          <v-btn
            color="red lighten-2"
            dark
            @click.stop="
              action = open;
              dialog = true;
            "
            >Open</v-btn
          >
        </v-col>
        <v-col cols="auto">
          <v-btn
            color="primary"
            @click="
              action = reset;
              dialog = true;
            "
            >Reset</v-btn
          >
        </v-col>
        <v-col cols="auto">
          <v-switch v-model="hide" label="Hide dice" />
        </v-col>
        <v-dialog v-model="dialog" width="500">
          <v-card>
            <v-card-title color="danger"> Are you sure? </v-card-title>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                color="primary"
                text
                @click="
                  action();
                  dialog = false;
                "
              >
                Yes
              </v-btn>
              <v-btn color="primary" text @click="dialog = false"> No </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
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
      hide: false,
      action: null,
      dialog: false,
      withOnes: true,
      highlight: false,

      startup: true,
      rerollTimeout: null,
    };
  },
  computed: {
    diceResult() {
      var self = this;
      return this.dice.reduce(
        (resultArray, item) => {
          console.log(item, item === 1, self.withOnes);
          if (item === 1 && self.withOnes) {
            resultArray[2] += 1;
            resultArray[3] += 1;
            resultArray[4] += 1;
            resultArray[5] += 1;
            resultArray[6] += 1;
          }
          resultArray[item] += 1;
          return resultArray;
        },
        { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0 }
      );
    },
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
    const host =
      process.env.VUE_APP_API_HOST || "localhost:8081" || window.location.host;
    const ws = new WebSocket(`ws://${host}/dice/${this.id}/ws`);
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
      console.log("open");
      this.ws.send(
        JSON.stringify({
          command: "open",
          data: "",
        })
      );
    },
    reset() {
      console.log("reset");
      this.ws.send(
        JSON.stringify({
          command: "reset",
          data: "",
        })
      );
    },
    rerollButton() {
      if (this.dice.length > 5 || this.startup) {
        this.send();
      } else {
        this.action = this.send;
        this.dialog = true;
      }
    },
    send() {
      this.ws.send(
        JSON.stringify({
          command: "start",
          data: "",
        })
      );
      if (this.rerollTimeout === null) {
        this.rerollTimeout = window.setTimeout(() => {
          this.startup = false;
        }, 1000);
      }
    },
  },
};
</script>
