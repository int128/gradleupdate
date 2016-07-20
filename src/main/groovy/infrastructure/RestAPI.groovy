package infrastructure

import groovy.transform.Canonical
import groovy.util.logging.Log
import wslite.rest.ContentType
import wslite.rest.RESTClient
import wslite.rest.Response

@Log
@Canonical
class RestAPI<E> {

    final Class<E> entityClass
    final String resourcePath
    final def parentEntity
    final RESTClient client

    static <E> RestAPI<E> of(Class<E> entityClass, String resourcePath, parentEntity, RESTClient client) {
        new RestAPI<E>(entityClass, resourcePath, parentEntity, client)
    }

    E get(Map queryMap = null, String key) {
        assert key
        log.info("Fetching $entityClass.simpleName($key) of $parentEntity")
        def response = client.get(path: "$resourcePath/$key", query: queryMap)
        switch (response.statusCode) {
            case 200: return withLog(response, entityClass.newInstance(parentEntity, response.json))
            default:  throw new RestAPIException(response)
        }
    }

    E find(Map queryMap = null, String key) {
        assert key
        log.info("Finding $entityClass.simpleName($key) of $parentEntity")
        def response = client.get(path: "$resourcePath/$key", query: queryMap)
        switch (response.statusCode) {
            case 200: return withLog(response, entityClass.newInstance(parentEntity, response.json))
            case 404: return withLog(response, null)
            default:  throw new RestAPIException(response)
        }
    }

    List<E> findAll(Map queryMap = null) {
        log.info("Finding $entityClass.simpleName of $parentEntity by $queryMap")
        def response = client.get(path: resourcePath, query: queryMap)
        switch (response.statusCode) {
            case 200: return withLog(response, response.json.collect { item ->
                entityClass.newInstance(parentEntity, item)
            })
            case 404: return withLog(response, null)
            default:  throw new RestAPIException(response)
        }
    }

    E create(Map contentJson) {
        log.info("Creating $entityClass.simpleName on $parentEntity")
        def response = client.post(path: resourcePath) {
            type ContentType.JSON
            json contentJson
        }
        switch (response.statusCode) {
            case 201: return withLog(response, entityClass.newInstance(parentEntity, response.json))
            default:  throw new RestAPIException(response)
        }
    }

    E update(Map contentJson, String key) {
        log.info("Updating $entityClass.simpleName($key) on $parentEntity")
        def response = client.patch(path: "$resourcePath/$key") {
            type ContentType.JSON
            json contentJson
        }
        switch (response.statusCode) {
            case 200: return withLog(response, entityClass.newInstance(parentEntity, response.json))
            default:  throw new RestAPIException(response)
        }
    }

    boolean delete(String key) {
        assert key
        log.info("Deleting $entityClass.simpleName($key) from $parentEntity")
        def response = client.get(path: "$resourcePath/$key")
        switch (response.statusCode) {
            case 200: return withLog(response, true)
            case 422: return withLog(response, false)
            default:  throw new RestAPIException(response)
        }
    }

    E invoke(Closure content = null) {
        log.info("Invoking $entityClass.simpleName on $parentEntity")
        def response = client.post(path: resourcePath, content)
        switch (response.statusCode) {
            case 202: return withLog(response, entityClass.newInstance(parentEntity, response.json))
            default:  throw new RestAPIException(response)
        }
    }

    private static <T> T withLog(Response response, T t) {
        log.info("$response.statusCode $response.statusMessage $t")
        t
    }

}
