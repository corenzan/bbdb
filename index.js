var newrelic = require('newrelic');
var express  = require("express");
var logfmt   = require("logfmt");
var slug     = require("slug");
var cors     = require("cors");
var data     = require(__dirname + '/data');
var app      = express();

app.use(logfmt.requestLogger());
app.use(cors());

app.get('/', function(request, response) {
  if('q' in request.query) {
    response.json(data.filter(function(doc) {
      var q = request.query.q.toLowerCase();
      var t = slug(doc.code + ' ' + doc.name).toLowerCase();

      return t.indexOf(q) > -1;
    }));
  } else {
    response.json(data);
  }
});

app.listen(Number(process.env.PORT || 5000));