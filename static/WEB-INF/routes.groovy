all '/rss', forward: '/rss.groovy', cache: (7 * 24).hours

post '/authorize', forward: '/exchange-oauth-token.groovy'

post '/webhook', forward: '/receive-github-webhook.groovy'
post '/webhook_test', forward: '/logging-github-webhook.groovy'

get '/@owner/@repo/status', forward: '/status.groovy?full_name=@owner/@repo'

get '/@owner/@repo/status.svg', forward: '/badge.groovy?full_name=@owner/@repo'
