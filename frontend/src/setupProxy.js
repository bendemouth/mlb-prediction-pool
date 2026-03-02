const { createProxyMiddleware } = require('http-proxy-middleware');

const backendProxy = createProxyMiddleware({
    target: 'http://backend:8080',
    changeOrigin: true,
    logLevel: 'debug',
});

module.exports = function(app) {
    app.use([
        '/health',
        '/leaderboard',
        '/predictions',
        '/users',
    ], backendProxy);
};