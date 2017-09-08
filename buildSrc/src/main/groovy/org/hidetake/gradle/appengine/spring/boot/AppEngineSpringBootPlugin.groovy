package org.hidetake.gradle.appengine.spring.boot

import org.gradle.api.Plugin
import org.gradle.api.Project

/**
 * A Gradle plugin for App Engine and Spring Boot.
 *
 * @author Hidetake Iwata
 */
class AppEngineSpringBootPlugin implements Plugin<Project> {
  static final EXTENSION_NAME = 'appengineSpringBoot'
  static final INJECT_APP_ENGINE_WEB_XML_TASK = 'injectAppEngineWebXml'
  static final WATCH_AND_SYNC_WEB_APP_TASK = 'watchAndSyncWebApp'

  @Override
  void apply(Project project) {
    project.extensions.create(EXTENSION_NAME, AppEngineSpringBootExtension).with {
      environment = project.file('.env')
    }

    project.tasks.create(WATCH_AND_SYNC_WEB_APP_TASK, WatchAndSyncWebAppTask)
    project.tasks.create(INJECT_APP_ENGINE_WEB_XML_TASK, InjectAppEngineWebXml)

    project.afterEvaluate {
      configureAppEnginePlugin(project)
      configureSpringBootPlugin(project)
    }
  }

  /**
   * Configure the App Engine plugin.
   */
  static void configureAppEnginePlugin(Project project) {
    assert project.hasProperty('appengine'), 'appengine-gradle-plugin must be applied'
    final extension = project.extensions.getByType(AppEngineSpringBootExtension)

    if (project.appengine.run.jvmFlags == null) {
      project.appengine.run.jvmFlags = []
    }
    if (extension.debugPort > 0) {
      project.appengine.run.jvmFlags.add(
        "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=${extension.debugPort}".toString())
    }

    if (project.appengine.run.environment == null) {
      project.appengine.run.environment = [:]
    }
    project.appengine.run.environment.putAll(extension.springBootDevProperties)

    project.tasks[INJECT_APP_ENGINE_WEB_XML_TASK].dependsOn('explodeWar')
    project.tasks.appengineRun.dependsOn(INJECT_APP_ENGINE_WEB_XML_TASK)
    project.tasks.appengineRun.dependsOn(WATCH_AND_SYNC_WEB_APP_TASK)
    project.tasks.appengineStage.dependsOn(INJECT_APP_ENGINE_WEB_XML_TASK)
  }

  /**
   * Configure the Spring Boot plugin.
   */
  static void configureSpringBootPlugin(Project project) {
    project.tasks.bootRepackage.enabled = false
    project.tasks.findMainClass.enabled = false
  }
}
