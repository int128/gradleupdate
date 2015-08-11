get '/rss', forward: '/rss.groovy', cache: (7 * 24).hours

post '/authorize', forward: '/exchange-oauth-token.groovy'

post '/webhook', forward: '/receive-github-webhook.groovy'
post '/webhook_test', forward: '/logging-github-webhook.groovy'
