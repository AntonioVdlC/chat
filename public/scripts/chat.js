new Vue({
  el: '#app',

  // Using ES6 string litteral delimiters because
  // Go uses {{ . }} for its templates
  delimiters: ['${', '}'],

  data: {
    ws: null,
    message: '',
    chat: [],
  },

  created: function() {
    this.ws = new WebSocket("wss://" + window.location.host + "/ws")
    this.ws.addEventListener('message', (e) => {
      let { user, avatar, content } = JSON.parse(e.data)
      this.chat.push({ user, avatar, content })
    })
  },

  methods: {
    send: function() {
      if (!this.message) {
        return
      }

      this.ws.send(
        JSON.stringify({
          content: this.message,
        })
      )
      this.message = ''
    },
  }
})
