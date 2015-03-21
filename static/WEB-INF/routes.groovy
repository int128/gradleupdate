get '/feed/@filter', forward: '/feed.groovy?filter=@filter', cache: (7 * 24).hours

post '/authorize', forward: '/exchange-oauth-token.groovy'

get  '/repos/@user/@repo', forward: '/get-user-repository.groovy?fullName=@user/@repo'
post '/repos/@user/@repo', forward: '/save-user-repository.groovy?fullName=@user/@repo'
all  '/repos/@user/@repo', forward: '/util/cors-options.groovy'

get '/', redirect: 'https://gradleupdate.github.io/'
