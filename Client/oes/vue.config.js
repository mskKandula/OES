
module.exports = {
 devServer: {
    port: 8085, // CHANGE YOUR PORT HERE!
    proxy: {
        '/api': {
                target: 'http://localhost:9000', // provide proxy  for your project
                ws: true,
                changeOrigin: true,
                pathRewrite: {
                    '^/api': ''
            }
        }
    }
}
}