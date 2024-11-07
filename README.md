enlace.space
---

Link storage service and RSS feed

## Accounts

Accounts are created on first request. Just pass the desired user name and password using the [basic authorization](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization#basic_authentication_2) scheme.


```bash
echo "erik:pass" | base64
ZXJpazpwYXNzCg==
```

## API

### Create Link

```bash
curl \
    -H "Authentication: Basic ZXJpazpwYXNzCg==" \
    -X POST \
    --data '{ "url": "https://foo.com/bar", "category": "fizzle" }' \
    https://enlace.space/links
```

### Get Feed

```
curl https://enlace.space/~erik/rss.xml
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
  <channel>
    <title>erik</title>
    <link>https://enlace.space/~erik/rss.xml</link>
    <item>
      <title>og:title</title>
      <description>og:description</description>
      <link>https://foo.com/bar</link>
      <pubDate>Sat, 07 Nov 2024 18:34:56 +0000</pubDate>
      <category>fizzle</category>
    </item>
  </channel>
</rss>
```
