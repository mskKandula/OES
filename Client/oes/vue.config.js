
module.exports = {
    devServer: {
        port: 8085, // CHANGE YOUR PORT HERE!
        https: true,
        proxy: {
            '/api': {
                target: 'http://localhost:9000', // provide proxy  for your project
                ws: true,
                changeOrigin: true,
                pathRewrite: {
                    '^/api': ''
                }
            },
            '/cdn': {
                target: 'http://localhost:8887/', // provide proxy  for your project
                ws: true,
                changeOrigin: true,
                pathRewrite: {
                    '^/cdn': ''
                }
            },
            '^/': {
                target: 'http://localhost:9000/', // provide proxy  for your project
                ws: true,
                changeOrigin: true
            }
        }
    }
}