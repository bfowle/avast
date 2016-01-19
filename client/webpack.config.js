var path = require('path');
// webpack plugins
var ProvidePlugin = require('webpack/lib/ProvidePlugin');
var DefinePlugin = require('webpack/lib/DefinePlugin');
var CommonsChunkPlugin = require('webpack/lib/optimize/CommonsChunkPlugin');
var CopyWebpackPlugin = require('copy-webpack-plugin');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var ENV = process.env.ENV = process.env.NODE_ENV = 'development';

var metadata = {
    baseUrl: '/',
    host: 'localhost',
    port: 3000,
    ENV: ENV,
    title: 'Avast',
    description: 'Raspberry Pirates tinkerapp',
    author: 'Brett Fowle <brettfowle@gmail.com>'
};

module.exports = {
    // static data for index.html
    metadata: metadata,
    // for faster builds use 'eval'
    devtool: 'source-map',
    debug: true,

    entry: {
        vendor: './src/vendor.ts',
        main: './src/main.ts'
    },

    // config for our build files
    output: {
        path: root('dist'),
        filename: '[name].bundle.js',
        sourceMapFilename: '[name].map',
        chunkFilename: '[id].chunk.js'
    },

    resolve: {
        // ensure loader extensions match
        extensions: [
            '',
            '.ts',
            '.js',
            '.json',
            '.css',
            '.html'
        ]
    },

    module: {
        preLoaders: [{
            test: /\.ts$/,
            loader: 'tslint-loader',
            exclude: [/node_modules/]
        }],
        loaders: [
            // support for *.ts files.
            {
                test: /\.ts$/,
                loader: 'ts-loader',
                query: {
                    'ignoreDiagnostics': [
                        2403, // 2403 -> Subsequent variable declarations
                        2300, // 2300 -> Duplicate identifier
                        2374, // 2374 -> Duplicate number index signature
                        2375  // 2375 -> Duplicate string index signature
                    ]
                },
                exclude: [
                    /\.(spec|e2e)\.ts$/,
                    /node_modules\/(?!(ng2-.+))/
                ]
            },
            // support for *.json files.
            {
                test: /\.json$/,
                loader: 'json-loader'
            },
            // support for *.css as raw text
            {
                test: /\.css$/,
                loader: 'raw-loader'
            },
            // support for *.html as raw text
            {
                test: /\.html$/,
                loader: 'raw-loader'
            }
        ]
    },

    plugins: [
        new CommonsChunkPlugin({
            name: 'vendor',
            filename: 'vendor.bundle.js',
            minChunks: Infinity
        }),
        // static assets
        new CopyWebpackPlugin([{
            from: 'src/assets',
            to: 'assets'
        }]),
        // generating html
        new HtmlWebpackPlugin({
            template: 'src/index.html',
            inject: false
        }),
        // replace
        new DefinePlugin({
            'process.env': {
                ENV: JSON.stringify(metadata.ENV),
                NODE_ENV: JSON.stringify(metadata.ENV)
            }
        })
    ],

    // other module loader config
    tslint: {
        emitErrors: false,
        failOnHint: false
    },

    // our webpack development server config
    devServer: {
        host: metadata.host,
        port: metadata.port,
        proxy: {
            '/api': {
                target: 'http://localhost:8080',
                secure: false
            },
            '/ws': {
                target: 'ws://localhost:8080',
                ws: true
            }
        },
        historyApiFallback: true,
        watchOptions: {
            aggregateTimeout: 300,
            poll: 1000
        }
    },

    // we need this due to problems with es6-shim
    node: {
        global: 'window',
        progress: false,
        crypto: 'empty',
        module: false,
        clearImmediate: false,
        setImmediate: false
    }
};

function root(args) {
    args = Array.prototype.slice.call(arguments, 0);
    return path.join.apply(path, [__dirname].concat(args));
}

function rootNode(args) {
    args = Array.prototype.slice.call(arguments, 0);
    return root.apply(path, ['node_modules'].concat(args));
}
