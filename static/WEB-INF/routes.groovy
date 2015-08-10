get '/rss', forward: '/rss.groovy', cache: (7 * 24).hours

post '/authorize', forward: '/exchange-oauth-token.groovy'

post '/webhook', forward: '/receive-github-webhook.groovy'

get '/', redirect: 'https://gradleupdate.github.io/'
