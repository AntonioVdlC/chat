module.exports = {
  map: {
    inline: true
  },
  plugins: [
    require("postcss-import")(),
    require("postcss-cssnext")({
     browsers: [ "> 5%" ] 
    }),
    require("cssnano")(),
  ]
}
