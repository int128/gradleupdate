package org.hidetake.gradle.appengine.spring.boot

import groovy.xml.XmlUtil
import org.gradle.api.Plugin
import org.gradle.api.Project

/**
 * A Gradle plugin for App Engine and Spring Boot.
 *
 * @author Hidetake Iwata
 */
class AppEngineSpringBootPlugin implements Plugin<Project> {
  Project project

  @Override
  void apply(Project project) {
    this.project = project

    configureExtension()

    project.afterEvaluate {
      configureAppEnginePlugin()
      configureSpringBootPlugin()
    }
  }

  /**
   * Configure an extension of the plugin.
   */
  void configureExtension() {
    final extension = project.extensions.create('appengineSpringBoot', AppEngineSpringBootExtension)
    extension.environment = project.file('.env')
  }

  /**
   * Configure the App Engine plugin.
   */
  void configureAppEnginePlugin() {
    if (!project.hasProperty('appengine')) {
      throw new IllegalStateException('appengine-gradle-plugin must be applied')
    }
    if (project.appengine.run.environment == null) {
      project.appengine.run.environment = [:]
    }
    if (project.appengine.run.jvmFlags == null) {
      project.appengine.run.jvmFlags = []
    }
    configureAppEngineEnvironmentForSpringBoot()
    configureDotEnvForAppEngineRun()
    configureDebugJvmFlag()
    configureDotEnvForAppEngineStage()
  }

  /**
   * Configure environment variables for Spring Boot dev.
   */
  void configureAppEngineEnvironmentForSpringBoot() {
    final extension = project.extensions.getByType(AppEngineSpringBootExtension)
    project.appengine.run.environment.putAll(extension.springBootDevProperties)
  }

  /**
   * Import user defined .env for run
   */
  void configureDotEnvForAppEngineRun() {
    final extension = project.extensions.getByType(AppEngineSpringBootExtension)
    final dotEnv = extension.loadEnvironmentOrNull()
    if (dotEnv) {
      project.appengine.run.environment.putAll(dotEnv)
    }
  }

  /**
   * Configure JVM flags for debug.
   */
  void configureDebugJvmFlag() {
    final extension = project.extensions.getByType(AppEngineSpringBootExtension)
    if (extension.debugPort > 0) {
      project.appengine.run.jvmFlags.add(
        "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=${extension.debugPort}".toString())
    }
  }

  /**
   * Import user defined .env for stage
   */
  void configureDotEnvForAppEngineStage() {
    project.tasks.appengineStage.doLast {
      final extension = project.extensions.getByType(AppEngineSpringBootExtension)
      final dotEnv = extension.loadEnvironmentOrNull()
      if (dotEnv) {
        final xml = project.file("${project.appengine.stage.stagingDirectory}/WEB-INF/appengine-web.xml")
        final root = new XmlParser().parse(xml)
        final envVariablesNode = root.get('env-variables').find() ?: root.appendNode('env-variables')
        dotEnv.each { k, v ->
          envVariablesNode.appendNode('env-var', [name: k, value: v])
        }
        xml.text = XmlUtil.serialize(root)
      }
    }
  }

  /**
   * Configure the Spring Boot plugin.
   */
  void configureSpringBootPlugin() {
    project.tasks.bootRepackage.enabled = false
    project.tasks.findMainClass.enabled = false
  }
}
