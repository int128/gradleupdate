get '/feed/@filter', forward: '/feed.groovy?filter=@filter', cache: (7 * 24).hours

get '/internal/purge-feed-cache', forward: '/purge-feed-cache.groovy'
