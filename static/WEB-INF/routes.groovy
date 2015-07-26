get '/stable', forward: '/get-stable-version.groovy'
all '/stable', forward: '/util/cors-options.groovy'

get '/rc',     forward: '/get-rc-version.groovy'
all '/rc',     forward: '/util/cors-options.groovy'

get '/stable/feed', forward: '/feed.groovy?filter=stable', cache: (7 * 24).hours
get '/rc/feed',     forward: '/feed.groovy?filter=rc', cache: (7 * 24).hours

post '/authorize', forward: '/exchange-oauth-token.groovy'

post '/repos/@user/@repo/gradle/@version',
    forward: '/trigger-update-user-repository.groovy?fullName=@user/@repo&gradleVersion=@version'
all  '/repos/@user/@repo/gradle/@version', forward: '/util/cors-options.groovy'

post '/webhook', forward: '/receive-github-webhook.groovy'

get '/', redirect: 'https://gradleupdate.github.io/'
