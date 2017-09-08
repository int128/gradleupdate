package org.hidetake.gradle.appengine.spring.boot

class AppEngineSpringBootExtension {
  /**
   * User defined .env file
   */
  File dotEnv

  /**
   * JVM debug port
   */
  int debugPort = 5005

  /**
   * @see org.springframework.boot.devtools.env.DevToolsPropertyDefaultsPostProcessor
   */
  Map springBootDevProperties = [
    "spring.thymeleaf.cache": "false",
    "spring.freemarker.cache": "false",
    "spring.groovy.template.cache": "false",
    "spring.mustache.cache": "false",
    "server.session.persistent": "true",
    "spring.h2.console.enabled": "true",
    "spring.resources.cache-period": "0",
    "spring.resources.chain.cache": "false",
    "spring.template.provider.cache": "false",
    "spring.mvc.log-resolved-exception": "true",
    "server.jsp-servlet.init-parameters.development": "true",
  ]
}
