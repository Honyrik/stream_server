module.exports = {
    transpileDependencies: ["vuetify"],
    devServer: {
        proxy: {
            "^/api": {
                target: "http://localhost:6090",
                ws: true,
                changeOrigin: true,
            },
        },
    },
};
