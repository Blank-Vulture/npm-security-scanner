const express = require('express');
const _ = require('lodash');

const app = express();
const port = 3000;

app.get('/', (req, res) => {
  const data = {
    message: 'Hello from demo NPM project!',
    timestamp: new Date().toISOString(),
    lodashVersion: _.VERSION
  };
  
  res.json(data);
});

app.listen(port, () => {
  console.log(`Demo app listening at http://localhost:${port}`);
});
