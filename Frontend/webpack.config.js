const path = require("path");

module.exports = {
    mode: "development",
    entry: "./src/index.js",
    output: {
        path: path.resolve(__dirname, "public"),
        filename: "main.js",
        publicPath: "/"
    },
    target: "web",
    devServer: {
        port: "9500",
        static: {
            directory: path.join(__dirname, "public"),
        },
        host: '0.0.0.0',
        allowedHosts: 'all',
        open: false,
        hot: true,
        liveReload: true,
        historyApiFallback: true
    },
    resolve: {
        extensions: ['.js', '.jsx', '.json']
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: 'babel-loader'
            },
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader', 'postcss-loader']
            }
        ]
    }
};
