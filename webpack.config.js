var webpack = require('webpack');
var CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
  entry: {
    app: './src/main/js/main.jsx'
  },
  output: {
    path: './build/exploded-app',
    filename: '[name].js'
  },
  externals: {
    react: 'React'
  },
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loader: 'babel-loader?presets[]=es2015&presets[]=react'
      },
      {
        test: /\.json$/,
        loader: 'json-loader'
      },
      {
        test: /\.less$/,
        loader: 'style!css!less?compress'
      }
    ]
  },
  plugins: [
    new webpack.optimize.UglifyJsPlugin(),
    new CopyWebpackPlugin([
      {from: './static'},
      {from: './node_modules/react/dist/react.min.js'},
      {from: './node_modules/bootswatch/flatly/bootstrap.min.css'},
      {from: './node_modules/bootswatch/fonts', to: 'fonts'},
      {from: './src/main/groovlet', to: 'WEB-INF/groovy'}
    ])
  ]
};
