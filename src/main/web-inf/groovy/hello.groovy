import groovy.json.JsonBuilder

response.contentType = 'application/json'
println new JsonBuilder({
    id 'hello'
    description 'Example REST API'
})
