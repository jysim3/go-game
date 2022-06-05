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
            KingJoker - room <code>{{ id }}</code>
          </h2>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="auto" v-if="status === ''">
          <v-card
            :color="opponentColor"
            width="200"
            height="200"
            class="d-flex align-center justify-center pa-4 mx-auto"
          >
            <v-progress-circular
              indeterminate
              color="primary"
            ></v-progress-circular>
          </v-card>
        </v-col>

        <v-col cols="4" v-else>
          <v-card v-if="(4 - opponentCount) > 0">
            <v-img
              class="otherCards-cards"
              eager
              :src="require(`@/assets/cards/1B.svg`)"
              v-for="(move, index) in (4-opponentCount)"
              :key="index"
              :style="`bottom: ${(3-opponentCount-(index)) * 5}px`"
            />
          </v-card> </v-col
        ><v-col cols="4">
          <v-card>
            <v-img :src="require(`@/assets/cards/1B.svg`)" />
          </v-card>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col
          cols="1"
          class="pa-0"
          v-for="(move, index) in historyFormatted"
          :key="index"
        >
          <v-card>
            <v-img :src="require(`@/assets/cards/${move[0]}.svg`)" />
          </v-card>
          <v-card>
            <v-img :src="require(`@/assets/cards/${move[1]}.svg`)" />
          </v-card>
        </v-col>
        <v-col cols="12" md="7">
          <v-alert
            v-if="status === 'won' || status === 'lost'"
            :type="status === 'won' ? 'success' : 'error'"
            prominent
          >
            <v-row align="center">
              <v-col class="grow"> You {{ status }} </v-col>
              <v-col class="shrink">
                <v-btn @click="start">Start again</v-btn>
              </v-col>
            </v-row>
          </v-alert>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="4">
          <v-hover v-slot="{ hover }">
            <v-card :disabled="status !== 'ready'">
              <div style="position: relative">
                <v-img
                  class="otherCards-cards"
                  @click="send('N')"
                  :src="require(`@/assets/cards/${move}.svg`)"
                  v-for="(move, index) in moves.others"
                  :key="move"
                  :style="`top: -${(moves.others.length - index - 1) * 5}px`"
                />
              </div>
              <v-fade-transition>
                <v-overlay v-if="hover" absolute z-index="5">
                  <v-btn @click="send('N')">Send others</v-btn>
                </v-overlay>
              </v-fade-transition>
            </v-card>
          </v-hover> </v-col
        ><v-col cols="4">
          <v-hover v-slot="{ hover }">
            <v-card :disabled="status !== 'ready'">
              <v-img
                @click="send('x')"
                :src="require(`@/assets/cards/${moves.key}.svg`)"
              />
              <v-fade-transition>
                <v-overlay v-if="hover" absolute>
                  <v-btn @click="send('x')">Send {{ moves.label }}</v-btn>
                </v-overlay>
              </v-fade-transition>
            </v-card>
          </v-hover>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
// import HelloWorld from "../components/HelloWorld";

function moves(isJoker, n = 4) {
  const special = isJoker ? "J" : "K";
  const suit = isJoker ? "S" : "H";
  const ret = {
    key: special + suit,
    label: isJoker ? "Joker" : "King",
    others:
      n > 0
        ? Array(n)
            .fill(0)
            .map((_, index) => 5 - index + suit)
        : [],
  };
  return ret;
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
    return {
      error: false,
      player: "",
      playerCount: 0,
      opponentCount: 0,
      history: [],
      opponent: "",
      status: "",
      ws: null,
      dialog: false,
    };
  },
  computed: {
    historyFormatted() {
      return Array(5)
        .fill(0)
        .map((_, index) => {
          const ret = [];
          const turn = this.history[index];
          if (!turn) {
            return ["1B", "1B"];
          }
          if (turn[0] === "N") {
            ret.push(index + 2 + "H");
          } else {
            ret.push("K" + "H");
          }
          if (!turn[1]) {
            ret.push("1B");
          }
          if (turn[1] === "n") {
            ret.push(index + 2 + "S");
          } else {
            ret.push("J" + "S");
          }
          if (this.player !== "joker") {
            ret[1] = [ret[0], (ret[0] = ret[1])][0];
          }
          return ret;
        });
    },
    moves() {
      return moves(this.player === "joker", 4 - (this.playerCount || 0));
    },
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
      // const host = process.env.VUE_APP_API_HOST || window.location.host;
      // const host = 'localhost:8081';
      const host =
        process.env.VUE_APP_API_HOST ||
        window.location.host;
      const ws = new WebSocket(`ws://${host}/joker/${this.id}/ws`);
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
        console.log(j);
        self.status = j.status;
        self.player = j.data.player;
        self.playerCount = j.data.player_past_moves;
        self.opponentCount = j.data.opponent_past_moves;
        self.history = j.data.history;
        console.log(self)
      };
    },
    reconnect() {
      this.error = false;
      this.status = "";
      this.setUpWebSocket();
    },
    reset() {
      this.error = false;
      const host = process.env.VUE_APP_API_HOST || window.location.host;
      fetch(`http://${host}/joker/${this.id}/reset`, {
        method: "POST",
      })
        .then(this.setUpWebSocket)
        .then(() => {
          this.status = "";
        });
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
    send(move) {
      if (move) {
        this.player = move;
        this.status = "sent";
        this.ws.send(
          JSON.stringify({
            command: "send",
            data: move,
          })
        );
      }
    },
  },
};
</script>
<style scoped>
.otherCards-cards:not(:first-child) {
  position: absolute;
  width: 100%;
  z-index: 0;
}
</style>
