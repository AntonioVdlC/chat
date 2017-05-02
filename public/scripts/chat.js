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
    let protocol = (window.location.protocol === 'https:') ? 'wss:' : 'ws:'
    this.ws = new WebSocket(`${protocol}//${window.location.host}/ws`)

    this.ws.addEventListener('message', (e) => {
      let {
        id,
        userId,
        userName,
        avatar,
        type,
        content,
        date,
      } = JSON.parse(e.data)
      
      this.chat.push({ id, userId, userName, avatar, type, content, date })
      
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
  }
})
