#!/usr/bin/env python3
## -*- coding: utf-8 -*-
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

from itertools import chain
from os import stat, mkdir
from time import strftime, localtime
from sys import argv
from glob import glob
from json import dumps
from shutil import copytree, copy, rmtree
import yaml

import frontmatter
from mako.template import Template
from mako.lookup import TemplateLookup
from bs4 import BeautifulSoup


srcdir = argv[1].rstrip('/') + '/'
dstdir = None
if len(argv) >= 3:
    dstdir = argv[2].rstrip('/') + '/'
else:
    from tempfile import mkdtemp
    dstdir = mkdtemp()

assert(srcdir != dstdir)

mylookup = TemplateLookup(directories=[srcdir + '/_templates'])
config  = dict()

def here():
    return anchor(config['baseurl'], config['name'])


def anchor(url, name=None):
    if not name:
        name = url
    return '<a rel="noopener" href="{0}">{1}</a>'.format(url, name)


def ref(url, name):
    return '<a target="_blank" rel="noopener nofollow" href="{0}">{1}</a>'.format(url, name)


def link(url, name):
    return '<a target="_blank" rel="noopener noreferrer nofollow" href="{0}">{1}</a>'.format(url, name)


def wiki(page, name):
    return link('https://en.wikipedia.org/wiki/' + page, name)


hooks = {'wiki': wiki,
         'link': link,
         'ref': ref,
         'anchor': anchor,
         'here': here}


def render(fmt, env):
    tpl = Template(fmt, lookup=mylookup)
    env = dict(env)
    env.update(hooks)
    env.update({'pages': pages, 'people': people})
    return tpl.render(**env).encode('utf-8')


class HtmlExtractor(object):
    def __init__(self):
        pass
    def exerpt(self, wot):
        soup = BeautifulSoup(wot, features="html.parser")
        for junk in soup(['nav', 'script', 'style', 'h1', 'h2', 'h3', 'footer', 'aside']):
            junk.extract()
        limit, count = 256, 0
        chunks = list()
        for l in soup.get_text().splitlines():
            l = l.strip()
            if not l:
                continue
            if count + len(l) > limit:
                chunks.append('...')
                break
            chunks.append(l)
            count += len(l)
        return ' '.join(chunks)


class PageSet(object):
    def __init__(self):
        self._pages = dict()
        self._blog = list()
        self._site = list()

    def __iter__(self):
        l = list()
        l.extend(self.site())
        l.extend(self.blog())
        return iter(l)

    def get(self, key):
        meta, data = self._pages[key]
        return meta, data

    def all(self):
        yield from self.site()
        yield from self.blog()

    def site(self):
        for k in self._site:
            meta, data = self.get(k)
            yield k, meta, data

    def blog(self):
        def _date(p):
            meta, _ = self.get(p)
            return meta['date']
        for k in sorted(self._blog, key=_date, reverse=True):
            meta, data = self.get(k)
            yield k, meta, data

    def _load(self, key, src, dst):
        content = ""
        env = dict(config)
        env.update({
            'key': key,
            'src': src,
            'dst': dst,
            'url': config['baseurl'] + '/' + dst[len(dstdir):],
            'path': dst[len(dstdir):]})

        # Load the template
        post = frontmatter.Frontmatter().read_file(src)
        if not post:
            with open(src, "r") as f:
                content = f.read()
        else:
            if post['attributes']:
                for k, v in post['attributes'].items():
                    env[k.lower()] = v
            content = post['body']

        # Ensure mandatory fields
        if 'date' not in env:
            date = stat(src).st_mtime
            date = strftime('%Y-%m-%d', localtime(date))
            env.update({"date": date})
        env['date'] = str(env['date'])

        if 'banner' not in env:
            env['banner'] = env['title']

        self._pages[key] = (env, content)

    def load_blog(self, key, src, dst):
        env = self._load(key, src, dst)
        self._blog.append(key)

    def load_site(self, key, src, dst):
        self._load(key, src, dst)
        self._site.append(key)


people = dict()
pages = PageSet()
extract_HTML = HtmlExtractor()


# Create the temporary target directory
try:
    stat(dstdir)
    rmtree(dstdir)
except:
    pass
finally:
    mkdir(dstdir)
    mkdir(dstdir + '/blog')
    mkdir(dstdir + '/docs')

# Load the main configuration and set default values
with open(srcdir + "_config.yml") as f:
    config = yaml.load(f)

if 'site_description' not in config:
    config['site_description'] = 'NOT-SET'
if 'description' not in config:
    config['description'] = config['site_description']

if 'prev' not in config:
    config['prev'] = 'NOT-SET.html'
if 'next' not in config:
    config['next'] = 'NOT-SET.html'
if 'author' not in config:
    config['author'] = 'NOT-SET'
if 'title' not in config:
    config['title'] = "NOT-SET"
if 'name' not in config:
    config['name'] = config['title']

# Load the people involved
for src in glob(srcdir + '_people/*.html'):
    meta = dict()
    post = frontmatter.Frontmatter().read_file(src)
    for k, v in post['attributes'].items():
        meta[k.strip().lower()] = v
    data = post['body']
    people[meta['nickname']] = (meta, post['body'])

# Load the Blog
for src in glob(srcdir + '_posts/*.html'):
    dst = src.replace(srcdir + '_posts/', dstdir + 'blog/')
    key = src.replace(srcdir + '_posts/', 'blog/')
    pages.load_blog(key, src, dst)

# Load the website
for src in chain(
        glob(srcdir + '/*.html'),
        glob(srcdir + 'blog/*.html'),
        glob(srcdir + 'docs/*.html')):
    dst = src.replace(srcdir, dstdir)
    key = src.replace(srcdir, '')
    pages.load_site(key, src, dst)

# Resolve the authors, the references, the next & prev
for key, meta, data in pages.all():
    print("patching>", repr(key), "\n", repr(meta))
    n, _ = pages.get(meta['next'])
    p, _ = pages.get(meta['prev'])
    if 'next_title' not in meta:
        meta['next_title'] = n['title']
    if 'next_url' not in meta:
        meta['next_url'] = '/' + n['key']
    if 'prev_title' not in meta:
        meta['prev_title'] = p['title']
    if 'prev_url' not in meta:
        meta['prev_url'] = '/' + p['key']
    if 'author_url' not in meta:
        meta['author_url'] = '/about.html#' + meta['author']

for nick, meta in people.items():
    print("\nPeople>", nick, "\n", dumps(meta, indent=' '))
for key, meta, _ in pages.all():
    print("\nPage>", key, "\n", dumps(meta, indent=' '))


# Generate the static html files of the website
# 1/ generate the blog articles, we need to generate an exerpt
#    of each article.
# 2/ generate the website, some pages might need the exerpt,
#    e.g. to generate an index of the blog posts.
for key, meta, data in pages.blog():
    print("\nprocessing>", key, dumps(meta, indent=' '))
    with open(meta['dst'], 'wb') as fout:
        rendered = render(data, meta)
        meta['exerpt'] = extract_HTML.exerpt(rendered)
        fout.write(rendered)
for key, meta, data in pages.site():
    print("\nprocessing>", key, dumps(meta, indent=' '))
    with open(meta['dst'], 'wb') as fout:
        fout.write(render(data, meta))


# Build the Sitemap
with open(dstdir + 'sitemap.xml', 'wb') as f:
    fmt = '''<url>
    <loc>${url}</loc>
    <changefreq>weekly</changefreq>
    <priority>1.0</priority>
    <lastmod>${date}</lastmod>
  </url>'''
    f.write(b'''<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">''')
    for key, meta, _ in pages.all():
        if 'draft' in meta:
            continue
        f.write(render(fmt, meta))
    f.write(b"\n</urlset>")


# Build the RSS feed
with open(dstdir + 'feed.xml', 'wb') as f:
    fmt = '''<item>
      <title>${title}</title>
      <link>${url}</link>
      <description>${description}</description>
    </item>'''
    f.write(render('''<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
  <channel>
    <title>${title}</title>
    <link>${baseurl}</link>
    <description>${description}</description>''', config))
    for key, meta, _ in pages.blog():
        if 'draft' in meta:
            continue
        f.write(render(fmt, meta))
    f.write(b"\n  </channel>\n</rss>")


# Copy the static files
special_files = (
    'robots.txt',
)
copytree(srcdir + 'static', dstdir + 'static')
for f in special_files:
    copy(srcdir + '/' + f, dstdir + '/' + f)

