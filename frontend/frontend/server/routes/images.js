const express = require('express');
const router = express.Router();
const mysqlx = require('@mysql/xdevapi');
const config = require('../config')

/* GET users listing. */
router.get('/', (req, res, next) => {
  mysqlx.getSession({user: config.USER, password: config.PASSWORD, host: config.HOST, port: config.PORT}).then(session => {
    session.getSchema("kube").getTable("images").select().execute().then(row => {
      res.json({images: row.fetchAll()});
    }).catch((e) => {throw e})
  }).catch((e) => {throw e})
});

module.exports = router;
