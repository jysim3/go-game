<template>
  <div>
    <v-container>
      <v-dialog v-model="error" persistent width="500">
        <v-card>
          <v-card-title class="text-h5 grey lighten-2">
            Room is full or server error. Try resetting room below
          </v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="primary" text @click="reset"> Reset </v-btn>
            <v-btn color="primary" text @click="reconnect">Reconnect</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <v-row justify="center" class="text-center">
        <v-col cols="auto">
          <h2 class="text-h2">
            RockPaperScissor - room <code>{{ id }}</code>
          </h2>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="auto">
          <v-card
            :color="opponentColor"
            width="200"
            height="200"
            class="d-flex align-center justify-center pa-4 mx-auto"
          >
            <v-progress-circular
              indeterminate
              v-if="status === ''"
              color="primary"
            ></v-progress-circular>
            <h1 class="text-h1" v-else-if="opponent === ''">?</h1>
            <v-img v-else :src="require(`@/assets/rock/${opponent}.png`)" />
          </v-card>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="auto">
          <h2 class="text-h2">VS</h2>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="4" v-for="move in moves" :key="move">
          <v-hover v-slot="{ hover }">
            <v-card
              :disabled="status !== 'start' && move !== player"
              :color="move === player ? playerColor : 'white'"
            >
              <v-img
                @click="send"
                :src="require(`@/assets/rock/${move}.png`)"
              />
              <v-fade-transition>
                <v-overlay v-if="hover" absolute>
                  <v-btn @click="send(move)">Send {{ move }}</v-btn>
                </v-overlay>
              </v-fade-transition>
            </v-card>
          </v-hover>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="auto">
          <v-btn
            v-if="status === 'win' || status === 'loss' || status === 'draw'"
            @click="start"
          >
            Restart
          </v-btn>
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
      error: false,
      moves: ["paper", "rock", "scissors"],
      player: "",
      opponent: "",
      status: "",
      ws: null,
      dialog: false,
    };
  },
  computed: {
    opponentColor() {
      return (
        {
          draw: "grey",
          loss: "success",
          win: "error",
        }[this.status] || "white"
      );
    },
    playerColor() {
      return {
        sent: "blue",
        draw: "grey",
        win: "success",
        loss: "error",
      }[this.status];
    },
  },
  unmounted() {
    this.ws.close();
  },
  mounted() {
    this.setUpWebSocket();
  },
  methods: {
    setUpWebSocket() {
      //const host = window.location.host;
      const host = 'localhost:8081';
      const ws = new WebSocket(
        `ws://${host}/rock/${this.id}/ws`
      );
      this.ws = ws;
      var self = this;
      ws.onerror = function (event) {
        self.status = "error";
        self.error = true;
        console.log(event);
      };
      ws.onclose = function (event) {
        console.log(event);
        self.error = true;
      };
      ws.onmessage = function (msg) {
        const j = JSON.parse(msg.data);
        if (j.command == "open") {
          self.opponent = j.data;
          if (self.opponent === self.player) {
            self.status = "draw";
          } else if (self.opponent === "rock") {
            self.status = self.player === "paper" ? "win" : "loss";
          } else if (self.opponent === "paper") {
            self.status = self.player === "scissors" ? "win" : "loss";
          } else if (self.opponent === "scissors") {
            self.status = self.player === "rock" ? "win" : "loss";
          }
        } else if (j.command === "start") {
          self.player = "";
          self.opponent = "";
          self.status = "start";
        }
      };
    },
    start() {
      this.player = "";
      this.opponent = "";
      this.status = "start";
      this.ws.send(
        JSON.stringify({
          command: "start",
        })
      );
    },
    reconnect() {
      this.error = false;
      this.status = "";
      this.setUpWebSocket();
    },
    reset() {
      this.error = false;
      fetch(`http://${window.location.host}/rock/${this.id}/reset`, {
        method: "POST",
      })
        .then(this.setUpWebSocket)
        .then(() => {
          this.status = "";
        });
    },
    send(move) {
      this.player = move;
      this.status = "sent";
      this.ws.send(
        JSON.stringify({
          command: "send",
          data: move,
        })
      );
    },
  },
};
</script>
