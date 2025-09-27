module.exports = {
  api: {
    input: './swagger/api.json',
    output: {
      client: 'axios',
      target: './src/api/api.ts',
      schemas: './src/api/model',
      baseUrl: 'http://localhost:8080',
    },
  }
};
