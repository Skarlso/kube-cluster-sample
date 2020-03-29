var express = require('express');
var router = express.Router();
var mysqlx = require('@mysql/xdevapi');

/* GET persons listing. */
router.get('/', function(req, res, next) {
  var connection = mysqlx.getSession({user: "root", password: "password123", host: "localhost", port: '33060'})
  connection.then(() => {
    res.json({persons: 'shit'});
  }).catch((e) => {alert(e)})
});

module.exports = router;
