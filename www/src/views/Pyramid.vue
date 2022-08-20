<template>
  <div>
    <v-container>
      <v-dialog v-model="targetedDialog" width="500" :persistent="targeted">
        <v-card>
          <v-card-title color="danger">
            <v-list-item three-line>
              <v-list-item-content>
                <div class="text-overline mb-4">TARGET</div>
                <v-list-item-title class="text-h5 mb-1">
                  {{
                    targetFrom === sessionId ? "You" : playerNames[targetFrom]
                  }}
                  sent
                  {{ targetTo === sessionId ? "you" : playerNames[targetTo] }} a
                  card!
                </v-list-item-title>
                <v-list-item-subtitle v-if="targetCard">
                  {{
                    targetTo === sessionId ? "you" : playerNames[targetTo]
                  }}
                  opened the card!
                </v-list-item-subtitle>
                <v-list-item-subtitle v-else-if="!targeted">
                  {{
                    targetTo === sessionId ? "you" : playerNames[targetTo]
                  }}
                  rejected the card!
                </v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col>
                <Card v-bind="gameCards(round - 1)" />
              </v-col>
              <v-col>
                <Card v-bind="targetOpenCard" />
              </v-col>

              <v-dialog v-model="targetOpenCardDialog" width="500">
                <v-card>
                  <v-card-title color="danger">
                    Are you sure to open?
                  </v-card-title>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="primary" text @click="acceptTarget">
                      Yes
                    </v-btn>
                    <v-btn
                      color="primary"
                      text
                      @click="targetOpenCardDialog = false"
                    >
                      No
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-col> </v-col>
            </v-row>
          </v-card-text>

          <v-card-actions>
            <v-spacer></v-spacer>

            <v-btn
              :color="targetCard || !targeted ? 'primary' : 'red lighten-2'"
              v-if="targetTo === sessionId || !targeted"
              @click="
                targetCard || !targeted
                  ? (targetedDialog = false)
                  : (rejectTargetDialog = true)
              "
            >
              Leave
            </v-btn>
            <v-dialog v-model="rejectTargetDialog" max-width="500px">
              <v-card>
                <v-card-title>
                  <span>Don't reveal the card?</span>
                  <v-spacer></v-spacer>
                </v-card-title>
                <v-card-actions>
                  <v-btn color="primary" text @click="rejectTarget">
                    Confirm
                  </v-btn>
                  <v-btn
                    color="primary"
                    text
                    @click="rejectTargetDialog = false"
                  >
                    Cancel
                  </v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-card-actions>
        </v-card> </v-dialog
      ><v-dialog v-model="openDialog" width="500">
        <v-card>
          <v-card-title color="danger"
            >Do you want to reveal the next card?</v-card-title
          >
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="primary"
              text
              @click="
                open();
                openDialog = false;
              "
            >
              Yes
            </v-btn>
            <v-btn color="primary" text @click="openDialog = false"> No </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-dialog v-model="sendDialog">
        <v-card>
          <v-card-title color="danger">Select Target</v-card-title>
          <v-card-text>
            <v-card>
              <v-list>
                <v-list-item-group v-model="selectedPlayer">
                  <!-- eslint-disable vue/no-use-v-if-with-v-for -->

                  <v-list-item
                    v-for="(player, key) in playerNames"
                    v-if="key !== sessionId"
                    :key="key"
                    :value="key"
                  >
                    <!-- eslint-enable -->
                    <!--
          <v-list-item-icon>
            <v-icon v-text="item.icon"></v-icon>
          </v-list-item-icon>
          -->
                    <v-list-item-content>
                      <v-list-item-title v-text="player"></v-list-item-title>
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-card>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="primary" text @click="verifyAndSend()"> Yes </v-btn>
            <v-btn color="primary" text @click="sendDialog = false"> No </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
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
          <h1 class="text-h2">Pyramid</h1>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="6">
          <v-card class="mx-4 blue darken-4" dark rounded="lg" elevation="0">
            <v-card-title>
              <v-icon>mdi-home</v-icon>
              <div class="ml-3">
                <div class="caption">Room:</div>
                <div>
                  <strong class="text-uppercase text--lighten-1">{{
                    id
                  }}</strong>
                </div>
              </div>
            </v-card-title>
          </v-card>
        </v-col>
        <v-col cols="6">
          <v-card class="mx-4 blue darken-4" dark rounded="lg" elevation="0">
            <v-card-title>
              <v-icon>mdi-account</v-icon>
              <div class="ml-3">
                <div class="caption">Players</div>
                <div>
                  <strong class="text-uppercase text--lighten-1">
                    {{
                      status === STATUS.DISCONNECTED
                        ? "Not Connected"
                        : Object.keys(playerNames).length
                    }}
                  </strong>
                </div>
              </div>
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="3" sm="auto">
          <Card v-bind="gameCards(9)" />
        </v-col>
      </v-row>
      <v-row justify="center" class="mt-n12">
        <v-col cols="3" sm="auto">
          <Card v-bind="gameCards(8)" />
        </v-col>
        <v-col cols="3" sm="auto">
          <Card v-bind="gameCards(7)" />
        </v-col>
      </v-row>
      <v-row justify="center" class="mt-n12">
        <v-col cols="3" sm="auto">
          <Card v-bind="gameCards(6)" />
        </v-col>
        <v-col cols="3" sm="auto"> <Card v-bind="gameCards(5)" /> </v-col>
        <v-col cols="3" sm="auto">
          <Card v-bind="gameCards(4)" />
        </v-col>
      </v-row>
      <v-row justify="center" class="mt-n12">
        <v-col cols="3" sm="auto">
          <Card v-bind="gameCards(3)" />
        </v-col>
        <v-col cols="3" sm="auto">
          <Card v-bind="gameCards(2)" />
        </v-col>
        <v-col cols="3" sm="auto"> <Card v-bind="gameCards(1)" /> </v-col
        ><v-col cols="3" sm="auto">
          <Card v-bind="gameCards(0)" />
        </v-col>
      </v-row>
      <v-divider class="mb-4" />
      <v-divider />
      <v-row justify="center" class="mt-4">
        <v-col cols="3" sm="auto" v-for="(card, index) in cards" :key="index">
          <Card
            :suit="card.suit"
            :number="card.number"
            :onClick="
              () => {
                sendDialog = true;
                sendCard = card;
              }
            "
          />
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col cols="12" v-if="status === STATUS.DISCONNECTED">
          <v-btn color="success" block @click.stop="connect">Reconnect</v-btn>
        </v-col>
        <v-col cols="12" v-else-if="round === -1">
          <v-btn color="success" block @click.stop="start">start</v-btn>
        </v-col>
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
        <v-btn small text @click.stop="start">Restart</v-btn>
      </div>
    </v-container>
  </div>
</template>

<script>
import Card from "../components/Card";
export default {
  name: "PyramidView",
  components: {
    Card,
  },
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
      name: "",
      nameDialog: false,
      sendDialog: false,
      sendCard: null,
      round: -1,
      status: "",
      sessionId: "",
      cards: [],
      openDialog: false,
      ws: null,

      selectedPlayer: null,
      playerNames: {},

      targeted: false,
      targetedDialog: false,
      targetTo: "",
      targetFrom: "",
      targetCard: null,
      targetOpenCardDialog: false,
      rejectTargetDialog: false,
    };
  },
  computed: {
    targetOpenCard() {
      var self = this;
      return (
        this.targetCard || {
          back: true,
          onClick:
            this.sessionId === this.targetTo
              ? () => {
                  self.targetOpenCardDialog = true;
                }
              : () => {
                  console.log("spectator");
                },
        }
      );
    },
  },
  mounted() {
    this.connect();
  },
  methods: {
    gameCards(index) {
      var self = this;
      index += 1;
      if (index === this.round) {
        return this.currentRoundCards;
      } else if (index > this.round) {
        if (index === this.round + 1) {
          return {
            onClick: () => {
              self.openDialog = true;
            },
            back: true,
          };
        }
        return {
          back: true,
        };
      } else {
        return {
          hidden: true,
        };
      }
    },
    connect() {
      this.status = this.STATUS.STARTUP;
      const host = process.env.VUE_APP_API_HOST || window.location.host;
      const ws = new WebSocket(`ws://${host}/pyramid/${this.id}/ws`);
      this.ws = ws;
      var self = this;
      ws.onopen = function () {};
      ws.onmessage = function (msg) {
        const j = JSON.parse(msg.data);
        console.log(j);
        switch (j.command) {
          case "setName":
            self.name = j.data;
            self.nameDialog = true;
            break;
          case "game":
            self.round = j.data.round;
            self.currentRoundCards = j.data.currentCard;

            if (j.data.target) {
              self.targetFrom = j.data.target.from;
              self.targetTo = j.data.target.to;

              self.targetedDialog = true;
              console.log(self.targetedDialog);
              self.targetCard = j.data.target.card;
              if (j.data.target.card) {
                self.targeted = false;
              } else {
                self.targeted = true;
              }
            } else {
              self.targeted = false;
            }
            break;
          case "playerUpdate":
            self.cards = j.data.cards;
            self.sessionId = j.data.id;
            break;
          case "names":
            self.playerNames = j.data;
            break;
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
    verifyAndSend() {
      if (this.selectedPlayer) {
        this.ws.send(
          JSON.stringify({
            command: "send",
            data: {
              card: this.sendCard,
              target: this.selectedPlayer,
            },
          })
        );
      }
      this.sendDialog = false;
    },
    acceptTarget() {
      this.targetedDialog = false;
      this.targetOpenCardDialog = false;
      this.ws.send(
        JSON.stringify({
          command: "accept",
        })
      );
    },
    rejectTarget() {
      this.rejectTargetDialog = false;
      this.targetedDialog = false;
      this.ws.send(
        JSON.stringify({
          command: "reject",
        })
      );
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
    start() {
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
