const express = require('express');
const router = express.Router();
const mysqlx = require('@mysql/xdevapi');
const config = require('../config');

/* GET persons listing. */
router.get('/', function(req, res, next) {
  mysqlx.getSession({user: config.USER, password: config.PASSWORD, host: config.HOST, port: config.PORT}).then(session => {
    session.getSchema("kube").getTable("person").select().execute().then(row => {
      let persons = []
      let result = row.fetchAll()
      result.forEach(e => {
        persons.push({id: e[0], name: e[1]})
      })
      res.json(persons)
    }).catch((e) => {throw e})
  }).catch((e) => {throw e})
});

module.exports = router;
