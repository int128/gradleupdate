import domain.LatestGradle
import groovy.json.JsonBuilder

response.contentType = 'application/json'
println new JsonBuilder({
    version new LatestGradle().version.string
})
