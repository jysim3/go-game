const { defineConfig } = require("@vue/cli-service");
module.exports = {
  devServer: {
    /*proxy: "http://localhost:8081",*/
  },
  ...defineConfig({
    transpileDependencies: true,
    outputDir: "build",
    assetsDir: "static/",
  }),
};
