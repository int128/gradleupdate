post '/api/exchange-oauth-token', forward: '/api/exchange-oauth-token.groovy'

post '/api/@owner/@repo/update', forward: '/api/create-pull-request-for-update-gradle.groovy?full_name=@owner/@repo'

get '/@owner/@repo/status', forward: '/render/repository-status-get.groovy?full_name=@owner/@repo'
get '/@owner/@repo/status.svg', forward: '/render/repository-badge-get.groovy?full_name=@owner/@repo'
