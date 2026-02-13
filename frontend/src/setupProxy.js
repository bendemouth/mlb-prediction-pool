const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
    app.use(
        ['/health'],
        createProxyMiddleware({
            target: 'http://backend:8080',
            changeOrigin: true,
            logLevel: 'debug',
        })
    );
};