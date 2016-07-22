import entity.PullRequestForLatestGradleWrapper
import groovy.json.JsonBuilder

def pullRequests = PullRequestForLatestGradleWrapper.findAll {
    select all from 'PullRequestForLatestGradleWrapper'
    sort desc by 'createdAt'
    limit 20
}

response.contentType = 'application/json'
println new JsonBuilder(pullRequests)
