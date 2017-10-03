package org.hidetake.gradle.appengine.spring.boot

import org.gradle.api.Plugin
import org.gradle.api.Project

/**
 * A Gradle plugin for App Engine and Spring Boot.
 *
 * @author Hidetake Iwata
 */
class AppEngineSpringBootPlugin implements Plugin<Project> {
  private AppEngineSpringBootExtension extension
  private InjectLoggingProperties injectLoggingProperties
  private InjectAppEngineWebXml injectAppEngineWebXml
  private WatchAndSyncWebAppTask watchAndSyncWebApp

  @Override
  void apply(Project project) {
    extension = project.extensions.create('appengineSpringBoot', AppEngineSpringBootExtension)
    extension.dotEnv = project.file('.env')
    watchAndSyncWebApp = project.tasks.create('watchAndSyncWebApp', WatchAndSyncWebAppTask)
    injectAppEngineWebXml = project.tasks.create('injectAppEngineWebXml', InjectAppEngineWebXml)
    injectLoggingProperties = project.tasks.create('injectLoggingProperties', InjectLoggingProperties)

    project.afterEvaluate {
      configureAppEnginePlugin(project)
      configureSpringBootPlugin(project)
    }
  }

  /**
   * Configure the App Engine plugin.
   */
  void configureAppEnginePlugin(Project project) {
    final appengine = project.findProperty('appengine')
    assert appengine, 'appengine-gradle-plugin must be applied'

    if (extension.debugPort > 0) {
      if (appengine.run.jvmFlags == null) {
        appengine.run.jvmFlags = []
      }
      appengine.run.jvmFlags.add(
        "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=${extension.debugPort}".toString())
    }

    if (injectLoggingProperties.enabled) {
      appengine.run.jvmFlags.add(
        "-Djava.util.logging.config.file=${injectLoggingProperties.loggingProperties}".toString())
    }

    if (appengine.run.environment == null) {
      appengine.run.environment = [:]
    }
    appengine.run.environment.putAll(extension.springBootDevProperties)
    appengine.run.environment.putAll(DotEnv.loadOrEmpty(extension.dotEnv))

    if (injectAppEngineWebXml.dotEnv == null) {
      injectAppEngineWebXml.dotEnv = extension.dotEnv
    }
    if (injectAppEngineWebXml.appEngineWebXml == null) {
      injectAppEngineWebXml.appEngineWebXml =
        project.file("${appengine.stage.stagingDirectory}/WEB-INF/appengine-web.xml")
    }

    project.tasks.appengineRun.dependsOn(injectLoggingProperties)
    project.tasks.appengineRun.dependsOn(watchAndSyncWebApp)
    project.tasks.appengineStage.finalizedBy(injectAppEngineWebXml)
    project.tasks.appengineDeploy.dependsOn(injectAppEngineWebXml)
  }

  /**
   * Configure the Spring Boot plugin.
   */
  static void configureSpringBootPlugin(Project project) {
    project.tasks.bootRepackage.enabled = false
    project.tasks.findMainClass.enabled = false
  }
}
