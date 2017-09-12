const fastify = require('fastify');
const slug = require('slug');
const database = require('./database');

const app = fastify();

app.use(require('response-time')());
app.use(require('cors')());

app.get('/', (req, reply) => {
  if('q' in req.query) {
    reply.send(database.filter(function(doc) {
      const q = slug(req.query.q).toLowerCase();
      const t = slug(doc.code + ' ' + doc.name).toLowerCase();
      return t.indexOf(q) > -1;
    }));
  } else {
    reply.send(database);
  }
});

app.listen(Number(process.env.PORT || 5000), (err) => {
  if (err) {
    throw err;
  }
  console.log('Listening...');
});
