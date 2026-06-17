const path = require('path');
module.exports = {entry: './src/plugin.tsx', output: {path: path.resolve(__dirname, 'dist'), filename: 'main.js', library: 'ReadReceiptPlugin', libraryTarget: 'window'}, resolve: {extensions: ['.ts', '.tsx', '.js']}, module: {rules: [{test: /\.tsx?$/, use: 'ts-loader', exclude: /node_modules/}]}, externals: {react: 'React', 'react-dom': 'ReactDOM'}};
