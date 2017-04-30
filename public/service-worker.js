const version = "v1"

self.addEventListener("install", function(event) {
  event.waitUntil(
    caches.open(version).then(function(cache) {
      return cache.addAll([
        "/",
        "/public/scripts/vendor/vue.min.js",
        "/public/scripts/chat.js",
        "/public/styles/main.css",
        "/public/styles/chat.css",
        "/public/styles/login.css",
      ])
    })
  )
})

self.addEventListener("fetch", function(event) {
  event.waitUntil(
    caches.match(event.request).then(function(cached) {
      const networked = fetch(event.request)
        .then(fetchedFromNetwork, unableToResolve)
        .catch(unableToResolve)

      return cached || networked

      function fetchedFromNetwork(response) {
        caches.open(version).then(function(cache) {
          cache.put(event.request, response.clone())
        })

        return response
      }

      function unableToResolve(error) {
        return new Reponse("<h1>Service unavailable</h1>", {
          status: 503,
          statusText: "Service unavailable",
          headers: new Headers({ "Content-Type": "text/html" })
        })
      }
    })
  )
})

self.addEventListener("activate", function(event) {
  event.waitUntil(
    caches.keys().then(function(keys) {
      return Promise.all(keys.map((key) => {
        if (key !== version) return caches.delete(key)
      }))
    })
  )
})
