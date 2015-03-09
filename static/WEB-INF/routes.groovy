get '/feed/@filter', forward: '/feed.groovy?filter=@filter', cache: (7 * 24).hours

get '/internal/@api', forward: '/@api.groovy'
