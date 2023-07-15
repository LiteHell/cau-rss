let express = require('express'),
    config = require('../../data/config.json'),
    path = require('path'),
    app = express(),
    CauBoardFeed = require('../feed/cau/board'),
    redis = require('redis'),
    redisClient = redis.createClient(config.redisConfig),
    rssMime = 'application/rss+xml; charset=utf-8',
    atomMime = 'application/atom+xml; charset=utf-8';

const titles = {
    'cse': '소프트웨어학부 공지사항',
    'abeek': '공학인증혁신센터 공지사항',
    'sw': '다빈치sw교욱원 공지사항',
    'dormitory': '중앙대학교 서울캠퍼스 기숙사 공지사항'
};

app.set('views', path.join(__dirname, '../../views'));
app.set('view engine', 'pug');
app.use(express.static(path.join(__dirname, '../../static')));

app.get('/cau/:parserName/:feedtype', (req, res, next) => {
    let {parserName, feedtype} = req.params
    req.cacheKey = `cau/${parserName}/${feedtype}`;
    redisClient.get(req.cacheKey, (err, reply) => {
        if (err || reply == null)
            next();
        else
            switch(feedtype) {
                case 'rss':
                    res.type(rssMime).end(reply)
                    break;
                case 'atom':
                    res.type(atomMime).end(reply)
                    break;
                default:
                    next(new Error("Unsupported feed type"));
            }
    })
}, (req, res, next) => {
    let {parserName, feedtype} = req.params;
    // parserName, title, description, link
    if(!titles[parserName])
        throw new Error("Unsupported board");
    if(feedtype !== 'rss' && feedtype !== 'atom')
        throw new Error("Unsupported feed type");

    let boardFeed = new CauBoardFeed(parserName, titles[parserName], titles[parserName], `https://${feedtype}.cau.ac.kr`);
    boardFeed[feedtype]().then(result => {
        redisClient.set(req.cacheKey, result, 'EX', 300);
        res.type(feedtype === 'rss' ? rssMime : atomMime).end(result)
    }).catch(next)
    
});
app.get('/', (req, res) => {
    res.render('index', {cauTitles: titles});
});

app.use(function (err, req, res, next) {
    console.error(err.stack)
    if (!res.headersSent)
        res.status(500).render('error');
});
module.exports = app;