new Vue({
  el: '#app',
  
  data: {
    ws: null,
    message: '',
    chatContent: '',
    email: null,
    username: null,
    joined: false,
  },

  created: function() {
    this.ws = new WebSocket(`ws://${window.location.host}/ws`)
    this.ws.addEventListener('message', (e) => {
      let message = JSON.parse(e.data)
      this.chatContent += `
        <li>
          <span>${message.username}</span>
          <span>${message.content}</span>
        </li>
      `
    })
  },

  methods: {
    send: function() {
      if (!this.message) {
        return
      }

      this.ws.send(
        JSON.stringify({
          email: this.email,
          username: this.username,
          content: this.message,
        })
      )
      this.message = ''
    },

    join: function() {
      if (!this.email || !this.username) {
        return
      }
      this.joined = true
    }
  }
})
