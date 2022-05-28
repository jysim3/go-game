<template>
  <div>
    <v-container>
      <v-row justify="center" class="text-center">
        <v-col cols="auto">
          <h1 class="text-h2">
            Dice - room <code>{{ id }}</code>
          </h1>
        </v-col>
      </v-row>
      <v-row v-if="dice.length > 5">
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
        <v-col cols="12">
          <v-btn color="success" block v-if="diceFlush" @click.stop="onReroll">Reroll</v-btn>
        </v-col>
        <v-col cols="12">
          <v-btn color="success" block v-if="status === STATUS.DISCONNECTED" @click.stop="connect">Reconnect</v-btn>
        </v-col>
        <v-col cols="auto">
          <v-btn :color="!rerollAllowed ? 'primary' : 'success'" @click.stop="onReroll">Restart Dice</v-btn>
        </v-col>
        <v-col cols="auto">
          <v-btn
            color="red lighten-2"
            dark
            :disabled="status == STATUS.DISCONNECTED"
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
      <v-btn fab @click.stop="backdoor" small bottom  absolute text class="ma-6" />
    </v-container>
  </div>
</template>

<script>
// import HelloWorld from "../components/HelloWorld";
    function diceResult(dice, withOnes) {
      return dice.reduce(
        (resultArray, item) => {
          if (item === 1 && withOnes) {
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
    }

export default {
  name: "DiceView",
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  data() {
    const STATUS = {
  DISCONNECTED: "DISCONNECTED",
  STARTUP: "STARTUP",
  PLAYING: "PLAYING",
  END: "END"
};
    return {
      STATUS,
      msg: "xxx",
      dice: [1, 2, 3, 4, 5],
      ws: null,
      hide: false,
      status: STATUS.STARTUP,
      action: null,
      dialog: false,
      withOnes: true,
      highlight: false,

      startup: true,
      rerollTimeout: null,
    };
  },
  computed: {
    diceFlush() {
return Object.values(diceResult(this.dice, false)).every(e => e <= 1);
    },
    rerollAllowed() {
return this.dice.length > 5 || this.status === this.STATUS.STARTUP ||  this.diceFlush;
    },
    diceResult() {
      console.log(diceResult(this.dice, this.withOnes));
      return diceResult(this.dice, this.withOnes);
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
    console.log(this.status === this.STATUS.STARTUP);
    this.connect();
  },
  methods: {
    connect() {
    this.status = this.STATUS.STARTUP;
    const host =
      process.env.VUE_APP_API_HOST || window.location.host;
      const ws = new WebSocket(`ws://${host}/dice/${this.id}/ws`);
    this.ws = ws;
    var self = this;
    ws.onmessage = function (msg) {
      const j = JSON.parse(msg.data);
      console.log(j);
      if (self.dice.lenght > 5) {
      this.status = this.STATUS.STARTUP;
      }
      if (j.command == "start") {
        self.dice = j.data;
      }
    };
    ws.onclose = function () {
      console.log("HIHI close");
      self.status = self.STATUS.DISCONNECTED;
    };
    ws.onerror = function () {
      console.log("HIHI");
      self.status = self.STATUS.DISCONNECTED;
    };

    },
    backdoor() {
      this.ws.send(
        JSON.stringify({
          command: "backdoor",
          data: [1,2,3,4,5]
        })
      );
    },
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
    onReroll() {
      if (this.ws.readyState !== this.ws.OPEN) {
        this.status = this.STATUS.DISCONNECTED;
      }
      if (this.rerollAllowed) {
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
          this.status = this.STATUS.PLAYING;
        }, 1000);
      }
    },
  },
};
</script>
