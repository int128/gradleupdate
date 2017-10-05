package org.hidetake.gradle.appengine.spring.boot

import org.gradle.api.Project

class AppEngineSpringBootExtension {
  /**
   * {@link DevLoggingPropertiesTask} task.
   */
  protected DevLoggingPropertiesTask devLoggingPropertiesTask

  /**
   * {@link WatchAndSyncWebAppTask} task.
   */
  protected WatchAndSyncWebAppTask watchAndSyncWebAppTask

  /**
   * {@link StageApplicationPropertiesTask} task.
   */
  protected StageApplicationPropertiesTask stageApplicationPropertiesTask

  /**
   * JVM debug port
   */
  int debugPort = 5005

  /**
   * @see org.springframework.boot.devtools.env.DevToolsPropertyDefaultsPostProcessor
   */
  Map<String, String> devProperties = [
    'spring.output.ansi.enabled': 'always',
    'spring.thymeleaf.cache': 'false',
    'spring.freemarker.cache': 'false',
    'spring.groovy.template.cache': 'false',
    'spring.mustache.cache': 'false',
    'server.session.persistent': 'true',
    'spring.h2.console.enabled': 'true',
    'spring.resources.cache-period': '0',
    'spring.resources.chain.cache': 'false',
    'spring.template.provider.cache': 'false',
    'spring.mvc.log-resolved-exception': 'true',
    'server.jsp-servlet.init-parameters.development': 'true',
  ]

  /**
   * Application properties for dev server.
   */
  File devPropertiesFile

  /**
   * Application properties for production.
   */
  Map<String, String> productionProperties = [:]

  /**
   * Application properties for production.
   */
  File productionPropertiesFile

  /**
   * Initialize properties.
   * @param project
   */
  protected void initialize(Project project) {
    devPropertiesFile = project.file('.properties')
    productionPropertiesFile = project.file('.properties')
  }

  /**
   * Compute JVM flags for dev server.
   */
  protected List<String> computeDevServerJvmFlags() {
    final jvmFlags = ["-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=$debugPort"]

    if (devLoggingPropertiesTask.enabled) {
      jvmFlags << "-Djava.util.logging.config.file=${devLoggingPropertiesTask.loggingProperties}"
    }
    if (devPropertiesFile?.exists()) {
      devPropertiesFile.withReader { reader ->
        final properties = new Properties()
        properties.load(reader)
        properties.each { key, value -> jvmFlags << "-D$key=$value" }
      }
    }
    devProperties.each { key, value -> jvmFlags << "-D$key=$value" }

    jvmFlags*.toString()
  }

  /**
   * Compute application.properties for production.
   */
  protected Map<String, String> computeProductionProperties() {
    final applicationProperties = [:]
    if (productionPropertiesFile?.exists()) {
      productionPropertiesFile.withReader { reader ->
        final properties = new Properties()
        properties.load(reader)
        applicationProperties.putAll(properties)
      }
    }
    applicationProperties << productionProperties
    applicationProperties
  }
}
