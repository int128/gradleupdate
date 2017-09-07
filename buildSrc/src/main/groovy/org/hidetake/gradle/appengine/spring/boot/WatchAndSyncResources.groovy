package org.hidetake.gradle.appengine.spring.boot

import org.gradle.api.DefaultTask
import org.gradle.api.tasks.TaskAction

class WatchAndSyncResources extends DefaultTask {
  @TaskAction
  void background() {
    assert project.appengine.stage.sourceDirectory instanceof File
    final webInfClasses = new File(project.appengine.stage.sourceDirectory as File, 'WEB-INF/classes')
    project.sourceSets.main.resources.srcDirs.each { File srcDir ->
      final watchAndSync = new WatchAndSync(srcDir.toPath(), webInfClasses.toPath())
      new Thread(watchAndSync).start()
    }
  }
}
