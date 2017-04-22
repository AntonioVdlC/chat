new Vue({
  el: '#app',

  data: {
    ws: null,
    message: '',
    chatContent: '',
  },

  created: function() {
    this.ws = new WebSocket("ws://" + window.location.host + "/ws")
    this.ws.addEventListener('message', (e) => {
      let { user, avatar, content } = JSON.parse(e.data)
      // This may need to be re-written sometime in the near future ... :P
      this.chatContent += "<li><img class='avatar' src='" + avatar + "' alt='" + user + "' title='" + user + "'><span class='message'>" + content + "</span></li>"
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
