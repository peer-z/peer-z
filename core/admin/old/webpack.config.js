/*
 * Copyright 2020 PeerGum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
// const Manifest = require('webpack-manifest-plugin');
// const ManifestOptions = {};
const {CleanWebpackPlugin} = require('clean-webpack-plugin');
const {VueLoaderPlugin} = require('vue-loader')
const CopyWebpackPlugin = require('copy-webpack-plugin')

module.exports = {
    mode: 'development',
    // target: 'node',
    // resolve: {
    //     alias: {
    //         // If using the runtime only build
    //         vue$: 'vue/dist/vue.runtime.esm.js' // 'vue/dist/vue.runtime.common.js' for webpack 1
    //         // Or if using full build of Vue (runtime + compiler)
    //         // vue$: 'vue/dist/vue.esm.js'      // 'vue/dist/vue.common.js' for webpack 1
    //     }
    // },
    entry: './src/index.js',
    plugins: [
        // new HtmlWebpackPlugin({
        //     template: './src/index.html',
        //     scriptLoading: "defer",
        //     inject: "head",
        // }),
        // new Manifest(ManifestOptions),
        // new CleanWebpackPlugin({ cleanStaleWebpackAssets: false }),
        new VueLoaderPlugin(),
        new CopyWebpackPlugin({
            patterns: [
                {
                    from: '*/*.html',
                    context: 'src'
                },
            ],
        }),
    ],
    output: {
        // filename: '[name].bundle.js',
        filename: 'js/bundle.js',
        path: path.resolve(__dirname, 'dist'),
        // assetModuleFilename: 'images/[hash][ext][query]',
    },
    devServer: {
        contentBase: './dist',
    },
    module: {
        rules: [
            {
                test: /\.vue$/,
                loader: 'vue-loader'
            },
            {
                test: /\.s[ac]ss$/i,
                use: [
                    // Creates `style` nodes from JS strings
                    // {
                    "style-loader",
                    // options: {
                    //     injectType: 'linkTag',
                    // },
                    // },
                    // Translates CSS into CommonJS
                    "css-loader",
                    // Compiles Sass to CSS
                    "sass-loader",
                ],
            },
            // {
            //     test: /\.(html)$/i,
            //     // loader: 'file-loader',
            // },
            {
                test: /\.(png|svg|jpe?g|gif)$/i,
                loader: 'file-loader',
                options: {
                    name: 'images/[name].[ext]',
                },
            },
            // {
            //     test: /\.(png|svg|jpg|jpeg|gif)$/i,
            //     // generator: {
            //     //     filename: 'images/[hash][ext][query]'
            //     // }
            // },
            // {
            //     test: /\.html/,
            //     type: 'asset/resource',
            //     generator: {
            //         filename: 'static/[hash][ext][query]'
            //     }
            // },
            // {
            //     test: /\.json5$/i,
            //     type: 'json',
            //     parser: {
            //         parse: json5.parse,
            //     },
            // },
        ],
    },
};