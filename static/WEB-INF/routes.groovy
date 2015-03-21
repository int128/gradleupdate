get '/stable', forward: '/get-stable-version.groovy'
all '/stable', forward: '/util/cors-options.groovy'

get '/rc',     forward: '/get-rc-version.groovy'
all '/rc',     forward: '/util/cors-options.groovy'

get '/stable/feed', forward: '/feed.groovy?filter=stable', cache: (7 * 24).hours
get '/rc/feed',     forward: '/feed.groovy?filter=rc', cache: (7 * 24).hours

post '/authorize', forward: '/exchange-oauth-token.groovy'

get  '/repos/@user/@repo', forward: '/get-user-repository.groovy?fullName=@user/@repo'
post '/repos/@user/@repo', forward: '/save-user-repository.groovy?fullName=@user/@repo'
all  '/repos/@user/@repo', forward: '/util/cors-options.groovy'

get '/', redirect: 'https://gradleupdate.github.io/'
