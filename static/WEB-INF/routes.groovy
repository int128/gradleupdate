get '/feed/@filter', forward: '/feed.groovy?filter=@filter', cache: (7 * 24).hours

post '/authorize', forward: '/exchange-oauth-token.groovy'

get '/', redirect: 'https://gradleupdate.github.io/'
