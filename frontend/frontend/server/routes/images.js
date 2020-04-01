const express = require('express');
const router = express.Router();
const mysqlx = require('@mysql/xdevapi');

/* GET users listing. */
router.get('/', (req, res, next) => {
  mysqlx.getSession({user: "root", password: "password123", host: "localhost", port: '33060'}).then(session => {
    session.getSchema("kube").getTable("images").select().execute().then(row => {
      res.json({images: row.toString()});
    }).catch((e) => {throw e})
  }).catch((e) => {throw e})
});

module.exports = router;
