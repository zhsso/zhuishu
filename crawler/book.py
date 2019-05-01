# -*- coding: utf-8 -*-
from gevent import monkey; monkey.patch_all()
from bottle import route, run
from bs4 import BeautifulSoup
import requests
import redis
import os
import json

redisaddr = os.environ["REDIS_PORT_6379_TCP_ADDR"]
redisport = os.environ["REDIS_PORT_6379_TCP_PORT"]
redisclient = redis.Redis(host=redisaddr, port=redisport, db=0)

def searchQidian(keyword):
    content = requests.get("https://www.qidian.com/search?kw=" + keyword).content
    soup = BeautifulSoup(content, "html.parser")
    books = soup.find_all(name='div', attrs={'class': 'book-mid-info'})
    infos = []
    for book in books:
        title = book.find(name='h4').find(name='a')
        url = "https://m.qidian.com/book/" + title['data-bid']
        book_id = title['data-bid']
        title = title.text
        author = book.find(name='a', attrs={'class': 'name'}).text
        infos.append({
            'platform': 'qidian',
            'title': title,
            "url": url,
            "book_id": int(book_id),
            "author": author
        })
    return infos

def searchZongheng(keyword):
    content = requests.get("http://search.zongheng.com/s?keyword=" + keyword).content
    soup = BeautifulSoup(content, "html.parser")
    books = soup.find_all(name='div', attrs={'class': 'fl se-result-infos'})
    infos = []
    for book in books:
        title = book.find(name='h2', attrs={'class': 'tit'}).find('a')
        url = title['href']
        book_id = json.loads(title['data-sa-d'])["book_id"]
        title = title.text
        author = book.find(name='div', attrs={'class': 'bookinfo'}).find(name='a').text
        infos.append({
            'platform': 'zongheng',
            'title': title,
            "url": url,
            "book_id": int(book_id),
            "author": author
        })
    return infos

def fetchNewestQidian(bookId):
    url = "https://m.qidian.com/book/{0}/catalog".format(bookId)
    content = requests.get(url).content.decode("utf-8")
    volumes = ""
    for line in content.split("\n"):
        if "g_data.volumes" in line:
            volumes = line
    vvv = []
    volumes = volumes.replace("true", "True")
    volumes = volumes.replace("false", "False")
    volumes = volumes.replace("g_data.volumes", "g_volumes")
    volumes = volumes.replace(";", "")
    volumes = "# coding=utf-8\n" + volumes + '\nvvv.extend(g_volumes)'
    exec(volumes)
    chapters = []
    for volume in vvv:
        for chapter in volume['cs']:
            chapters.append({
                'title': chapter['cN'],
                'url': "https://m.qidian.com/book/{0}/{1}".format(bookId, chapter['id'])
            })
    return chapters

def fetchNewestZongheng(bookId):
    url = "http://book.zongheng.com/showchapter/{0}.html".format(bookId)
    content = requests.get(url).content
    soup = BeautifulSoup(content, "html.parser")
    volumes = soup.find_all(name='ul', attrs={'class': 'chapter-list clearfix'})
    chapters = []
    for volume in volumes:
        cs = volume.find_all('li')
        for c in cs:
            c = c.find('a')
            chapters.append({
                'title': c.text,
                'url': c['href']
            })
    return chapters

def filterChapter(platform, bookId, chapters):
    skey = "{0}/{1}".format(platform, bookId)
    ret = []
    for chapter in chapters[::-1]:
        if redisclient.sadd(skey, chapter['url']) == 1:
            ret.append(chapter)
    return ret

def fetchNewest(platform, bookId):
    chapters = []
    if platform == "qidian":
        chapters = fetchNewestQidian(bookId)
    elif platform == "zongheng":
        chapters = fetchNewestZongheng(bookId)
    return filterChapter(platform,bookId, chapters)

@route('/search/<name>')
def searchBook(name):
    infos = searchQidian(name)
    infos.extend(searchZongheng(name))
    return json.dumps(infos)


@route('/fetch/<platform>/<bookId>')
def fetchBook(platform, bookId):
    infos = fetchNewest(platform, bookId)
    return json.dumps(infos)

@route('/mark/<platform>/<bookId>')
def markBook(platform, bookId):
    fetchNewest(platform, bookId)
    return "done"

run(host='0.0.0.0', port=8080, server='gevent')
