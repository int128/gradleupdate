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

  @Override
  void apply(Project project) {
    extension = project.extensions.create('appengineSpringBoot', AppEngineSpringBootExtension)
    extension.devLoggingPropertiesTask = project.tasks.create('devLoggingProperties', DevLoggingPropertiesTask)
    extension.watchAndSyncWebAppTask = project.tasks.create('watchAndSyncWebApp', WatchAndSyncWebAppTask)
    extension.stageApplicationPropertiesTask = project.tasks.create('stageApplicationProperties', StageApplicationPropertiesTask)
    extension.initialize(project)

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

    if (appengine.run.jvmFlags == null) {
      appengine.run.jvmFlags = []
    }
    appengine.run.jvmFlags.addAll(extension.computeDevServerJvmFlags())

    project.tasks.appengineRun.dependsOn(extension.devLoggingPropertiesTask)
    project.tasks.appengineRun.dependsOn(extension.watchAndSyncWebAppTask)
    project.tasks.appengineStage.finalizedBy(extension.stageApplicationPropertiesTask)
  }

  /**
   * Configure the Spring Boot plugin.
   */
  static void configureSpringBootPlugin(Project project) {
    project.tasks.bootRepackage.enabled = false
    project.tasks.findMainClass.enabled = false
  }
}
