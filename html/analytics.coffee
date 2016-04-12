

TRACK_URL = 'http://localhost:8001/golang.gif'
COOKIE_NAME = '__goanalytics'

get_uuid = ->
  d = (new Date).getTime()
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace /[xy]/g, (c) ->
    r = (d + Math.random() * 16) % 16 | 0
    d = Math.floor(d / 16)
    return (if c == 'x' then r else r & 0x3 | 0x8).toString 16

get_cookie = (c_name) ->
  if document.cookie.length > 0
    c_start = document.cookie.indexOf(c_name + '=')
    if c_start != -1
      c_start = c_start + c_name.length + 1
      c_end = document.cookie.indexOf(';', c_start)
      if c_end == -1
        c_end = document.cookie.length
      return unescape(document.cookie.substring(c_start, c_end))
  return ''

set_cookie = (c_name, value) ->
  document.cookie = c_name + '=' + escape(value)

get_or_set_cookie = (c_name) ->
  uuid = get_cookie(c_name)
  if not uuid
    uuid = get_uuid()
    set_cookie(c_name, uuid)
  return uuid

# https://github.com/piwik/piwik/blob/master/js/piwik.js#L2421
# http://google-analytics.com/ga.js

send_data = ->
  params =
    id: get_uuid()
    cookieid: get_or_set_cookie(COOKIE_NAME)
    referer: encodeURIComponent(document.referrer)
    width: window.screen.width
    height: window.screen.height
    color: window.screen.colorDepth
    title: encodeURIComponent(document.title)
    lang: if navigator.browserLanguage then navigator.browserLanguage else navigator.language

  console.log params

  quary_string = ''
  for k of params
    quary_string += '&' + k + '=' + params[k]
  src = TRACK_URL + '?' + quary_string.slice(1, quary_string.length)

  image = new Image(1, 1)
  image.src = src

send_data()
