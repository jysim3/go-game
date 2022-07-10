<template>
  <div>
    <v-container>
      <v-dialog v-model="nameDialog">
        <v-card>
          <v-card-title>
            <span class="text-h5">User Profile</span>
          </v-card-title>
          <v-card-text>
            <v-text-field label="Name" v-model="name"></v-text-field>
          </v-card-text>

          <v-card-actions>
            <v-btn color="blue darken-1" text @click="setName"> Save </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-row justify="center" class="text-center">
        <v-col cols="auto">
          <h1 class="text-h2">Dice Game</h1>
        </v-col>
        <v-col cols="2" class="d-flex align-center">
          <v-img
            :src="require(`@/assets/dice/6.png`)"
            class="rounded-lg"
            max-width="50"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-card class="mx-4" rounded="lg">
            <v-card-title>
              <v-btn color="yellow accent-4" fab elevation="0">
                <v-img
                  :src="require('@/assets/room-icon.png')"
                  class="rounded-lg"
                  max-width="30"
                />
              </v-btn>
              <div class="ml-3">
                Room:
                <strong class="red--text text--lighten-1">{{ id }}</strong>
              </div>
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="auto">
          Number of players:
          {{ n_players === null ? "Not Connected" : n_players }}
        </v-col>
      </v-row>
      <v-row v-if="roundEnded">
        <v-col>
          <v-card rounded="xl">
            <v-card-title>
              <h3>Results</h3>
            </v-card-title>
            <v-card-text>
              <v-row>
                <v-col
                  v-for="die in [1, 2, 3, 4, 5, 6]"
                  :key="die"
                  class="d-flex justify-center flex-column text-center align-center"
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
                  <span class="text-h6">{{ diceSum[die] || "0" }}</span>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <div v-if="!hide">
        <v-row v-if="roundEnded">
          <v-col><h3 class="text-center">Player Dices</h3></v-col>
        </v-row>
        <v-row
          justify="center"
          v-for="(dice, name) in diceFormatted"
          :key="name"
        >
          <v-col
            cols="2"
            v-if="roundEnded"
            class="overflow-hidden"
            style="box-shadow: inset rgb(0 0 0 / 20%) -9px 0px 9px 0px"
          >
            {{ name }}</v-col
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
          <v-btn
            color="success"
            block
            @click.stop="connect"
            v-if="status === STATUS.DISCONNECTED"
            >Reconnect</v-btn
          >
          <v-btn
            color="success"
            block
            @click.stop="onReroll"
            v-else-if="rerollAllowed"
            >Reroll</v-btn
          >
          <v-btn
            v-else
            color="red lighten-2"
            dark
            block
            @click.stop="
              action = open;
              dialog = true;
            "
            >Open</v-btn
          >
        </v-col>
        <v-col cols="auto">
          <v-switch v-model="hide" label="Hide dice" />
        </v-col>
        <v-col cols="auto">
          <v-switch v-model="withOnes" label="With 1s" />
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
      <v-btn fab @click.stop="backdoor" small absolute text right />
      <div style="position: absolute; bottom: 0">
        <v-btn
          small
          text
          @click="
            action = reset;
            dialog = true;
          "
          >Reset</v-btn
        >
        <v-btn small text @click="nameDialog = true">Rename</v-btn>
      </div>
    </v-container>
  </div>
</template>

<script>
// import HelloWorld from "../components/HelloWorld";
function diceSum(dice, withOnes) {
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
      CONNECTED: "CONNECTED",
      PLAYING: "PLAYING",
      END: "END",
    };
    return {
      STATUS,
      msg: "xxx",
      dice: [1, 2, 3, 4, 5],
      diceResult: null,
      name: "",
      nameDialog: false,
      ws: null,
      wsReconnect: null,
      hide: false,
      status: STATUS.STARTUP,
      action: null,
      dialog: false,
      withOnes: true,
      highlight: false,
      n_players: null,

      startup: true,
      rerollTimeout: null,
    };
  },
  computed: {
    diceFlush() {
      return Object.values(diceSum(this.dice, false)).every((e) => e <= 1);
    },
    rerollAllowed() {
      return this.diceResult != null || this.diceFlush;
    },
    diceSum() {
      return diceSum(
        Object.values(this.diceResult).reduce((a, c) => a.concat(c), []),
        this.withOnes
      );
    },
    diceFormatted() {
      if (this.diceResult) {
        return this.diceResult;
      }
      return {
        "": this.dice,
      };
    },
    roundEnded() {
      return this.diceResult !== null;
    },
  },
  mounted() {
    this.connect();
  },
  methods: {
    connect() {
      this.status = this.STATUS.STARTUP;
      const host = process.env.VUE_APP_API_HOST || window.location.host;
      const ws = new WebSocket(`ws://${host}/dice/${this.id}/ws`);
      this.ws = ws;
      var self = this;
      ws.onopen = function () {
        if (self.wsReconnect) {
          clearTimeout(self.wsReconnect);
          self.wsReconnect = null;
        }
        self.status = self.STATUS.CONNECTED;
      };
      ws.onmessage = function (msg) {
        const j = JSON.parse(msg.data);
        console.log(j);
        if (self.dice.lenght > 5) {
          self.status = self.STATUS.STARTUP;
        }
        if (j.command == "setName") {
          self.nameDialog = true;
        }
        if (j.command == "start") {
          self.dice = j.data;
          self.diceResult = null;
        }
        if (j.command == "open") {
          self.diceResult = j.data;
        }
        if (j.command == "players") {
          self.n_players = j.data;
        }
      };
      ws.onclose = function () {
        self.status = self.STATUS.DISCONNECTED;
        if (!self.wsReconnect) {
          // self.wsReconnect = setTimeout(self.connect, 1000);
        }
      };
      ws.onerror = function () {
        self.status = self.STATUS.DISCONNECTED;
        if (!self.wsReconnect) {
          // self.wsReconnect = setTimeout(self.connect, 1000);
        }
      };
    },

    setName() {
      this.nameDialog = false;
      this.ws.send(
        JSON.stringify({
          command: "setName",
          data: this.name,
        })
      );
    },
    backdoor() {
      this.ws.send(
        JSON.stringify({
          command: "backdoor",
          data: [1, 2, 3, 4, 5],
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
