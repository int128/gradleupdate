post '/api/authorize', forward: '/api/auth/exchange-token.groovy'

post '/api/@owner/@repo/update', forward: '/api/github-repository/update.groovy?full_name=@owner/@repo'

post '/api/webhook', forward: '/api/github-webhook/receive-event.groovy'
post '/api/webhook/test', forward: '/api/github-webhook/receive-test.groovy'

get '/@owner/@repo/status', forward: '/github-repository/status/get.groovy?full_name=@owner/@repo'

get '/@owner/@repo/status.svg', forward: '/github-repository/badge/get.groovy?full_name=@owner/@repo'
