new Vue({
  el: '#app',

  // Using ES6 string litteral delimiters because
  // Go uses {{ . }} for its templates
  delimiters: ['${', '}'],

  components: {
    'vue-pull-refresh': VuePullRefresh
  },

  data: {
    ws: null,
    message: '',
    chat: [],
    users: [],
    showUsers: false,
  },

  created: function() {
    let protocol = (window.location.protocol === 'https:') ? 'wss:' : 'ws:'
    this.ws = new WebSocket(`${protocol}//${window.location.host}/ws`)

    this.ws.addEventListener('message', (e) => {
      let data = JSON.parse(e.data)

      // Bootstrap
      if (data.type === "bootstrap") {
        let { messages, users } = JSON.parse(e.data)

        // Retrieve messages
        messages.forEach((message) => {
          let { id, userId, userName, avatar, type, content, date } = message
          this.chat.push({ id, userId, userName, avatar, type, content, date: new Date(date) })
        })

        // Retrieve users
        this.users = [...users]
      }

      // Bulk messages
      else if (data.type === "messages") {
        let { messages } = JSON.parse(e.data)

        // Retrieve messages
        messages.forEach((message) => {
          let { id, userId, userName, avatar, type, content, date } = message
          this.chat.push({ id, userId, userName, avatar, type, content, date: new Date(date) })
        })
      }

      // Chat message
      else {
        let { id, userId, userName, avatar, type, content, date } = data

        this.chat.push({ id, userId, userName, avatar, type, content, date: new Date(date),})

        // Update users list if login or logout message
        if (type === 'login') {
          // As VueJS template directives don't iterate over Set or Map
          // make sure we are not adding duplicated.
          if (!this.users.some(user => user.id === userId)) {
            this.users.push({ id: userId, name: userName, avatar })
          }
        } else if (type === 'logout') {
          // As VueJS template directives don't iterate over Set or Map
          // make sure we are deleting a user that exists.
          if (this.users.some(user => user.id === userId)) {
            this.users.splice(this.users.findIndex(user => user.id === userId), 1)
          }
        }
      }

      // Make sure the messages are always sorted chronologically by date
      this.chat.sort((a, b) => 
        (new Date(a.date) > new Date(b.date)) ? 1 : 
        (new Date(a.date) < new Date(b.date)) ? -1 :
        0
      )
    })

    this.ws.addEventListener('close', (e) => {
      this.chat.push({
        type: "warning",
        content: i18n["chat_warning_connection_closed"],
      })
      this.ws = null
    })

    this.ws.addEventListener('error', (e) => {
      this.chat.push({
        type: "warning",
        content: i18n["chat_warning_connection_error"],
      })
      this.ws = null
    })
  },

  updated: function () {
    this.scrollToLast()
  },

  methods: {
    send: function() {
      if (!this.message) {
        return
      } else if (this.message.length > 270) {
        return
      }

      this.ws && this.ws.send(
        JSON.stringify({
          type: 'message',
          content: this.message,
        })
      )
      this.message = ''
    },
    scrollToLast: function (force = false) {
      let $chat = this.$refs.chat
      let doScroll = force ||
        $chat.scrollTop > $chat.scrollHeight - $chat.clientHeight - 50
      
      if (doScroll) {
        $chat.scrollTop = $chat.scrollHeight - $chat.clientHeight
      }
    },
    toggleUsers: function() {
      this.showUsers = !this.showUsers
    },
    onPullRefresh: function() {
      return new Promise((resolve, reject) => {
        if (this.ws) {
          this.ws.send(
            JSON.stringify({
              type: 'request',
              content: 'olderMessages',
              date: this.chat[this.chat.length - 1].date
            })
          )
          this.ws.addEventListener('message', () => resolve())
          this.ws.addEventListener('error', () => reject())
        } else {
          reject()
        }
      })
    },
  }
})
