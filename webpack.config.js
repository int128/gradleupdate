const webpack = require('webpack');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
  entry: {
    app: './src/main/js/main.jsx'
  },
  output: {
    path: './build/exploded-app',
    filename: '[name].js'
  },
  externals: {
    'react': 'React',
    'react-dom': 'ReactDOM',
    'react-router': 'ReactRouter'
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
        exclude: /node_modules/,
        loader: 'json-loader'
      },
      {
        test: /\.less$/,
        exclude: /node_modules/,
        loader: 'style!css!less?compress'
      }
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': {
        'NODE_ENV': JSON.stringify(process.env.NODE_ENV)
      }
    }),
    new CopyWebpackPlugin([
      {from: './static'},
      {from: './node_modules/react/dist/react.min.js'},
      {from: './node_modules/react-dom/dist/react-dom.min.js'},
      {from: './node_modules/react-router/umd/ReactRouter.min.js'},
      {from: './node_modules/bootswatch/lumen/bootstrap.min.css'},
      {from: './node_modules/bootswatch/fonts', to: 'fonts'},
      {from: './src/main/groovlet', to: 'WEB-INF/groovy'}
    ])
  ],
  devServer: {
    historyApiFallback: true,
    contentBase: './static'
  }
};
