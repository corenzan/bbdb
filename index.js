require('newrelic');

var express = require("express");
var logfmt  = require("logfmt");
var slug    = require("slug");
var data    = require(__dirname + '/data');
var app     = express();

app.use(logfmt.requestLogger());

app.get('/', function(request, response) {
  if('q' in request.query) {
    response.send(data.filter(function(doc) {
      return slug(doc.code + ' ' + doc.name).toLowerCase().indexOf(request.query.q) > -1;
    }));
  } else {
    response.send(data);
  }
});

app.listen(Number(process.env.PORT || 5000));