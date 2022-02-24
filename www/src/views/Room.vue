<template>
  <div>
    <v-container>
      <v-row justify="center">
        <v-col>
          <v-simple-table>
            <thead>
              <tr>
                <th class="text-left">Recipient</th>
                <th class="text-left">Sender</th>
                <th class="text-left">Amount</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in transactions" :key="index">
                <td>{{ item.recipient }}</td>
                <td>{{ item.sender }}</td>
                <td>{{ item.amount }}</td>
              </tr>
            </tbody>
          </v-simple-table>
        </v-col>
      </v-row>
      <v-row>
        <v-btn color="primary" @click="send">{{ id }}</v-btn>
      </v-row>
    </v-container>
  </div>
</template>

<script>
// import HelloWorld from "../components/HelloWorld";

export default {
  name: "RoomView",
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      msg: "xxx",
      transactions: [],
      ws: null,
    };
  },
  mounted() {
    const ws = new WebSocket(`ws://localhost:5000/channel/${this.id}/ws`);
    this.ws = ws;
    this.msg = "hellooo";
    var self = this;
    ws.onmessage = function (msg) {
      console.log(msg);
      self.transactions.push(JSON.parse(msg.data));
    };
  },
  methods: {
    send() {
      this.ws.send(
        JSON.stringify({
          recipient: 1,
          sender: 2,
          amount: 10,
        })
      );
    },
  },
};
</script>
