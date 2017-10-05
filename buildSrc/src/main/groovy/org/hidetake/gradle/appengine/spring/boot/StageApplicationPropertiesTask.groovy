package org.hidetake.gradle.appengine.spring.boot

import org.gradle.api.DefaultTask
import org.gradle.api.tasks.OutputFile
import org.gradle.api.tasks.TaskAction

class StageApplicationPropertiesTask extends DefaultTask {
  @OutputFile
  def applicationPropertiesFile = {
    new File(project.appengine.stage.stagingDirectory as File, 'WEB-INF/classes/config/application.properties')
  }

  @TaskAction
  void stage() {
    final extension = project.extensions.getByType(AppEngineSpringBootExtension)
    project.file(applicationPropertiesFile).withWriter { writer ->
      final properties = new Properties()
      properties.putAll(extension.computeProductionProperties())
      properties.store(writer, null)
    }
  }
}
