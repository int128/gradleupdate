package org.hidetake.gradle.appengine.spring.boot

import groovy.xml.XmlUtil
import org.gradle.api.DefaultTask
import org.gradle.api.tasks.TaskAction

/**
 * A task to inject env-variables into appengine-web.xml
 */
class InjectAppEngineWebXml extends DefaultTask {
  @TaskAction
  void mutate() {
    final xml = project.file("${project.appengine.stage.sourceDirectory}/WEB-INF/appengine-web.xml")
    final extension = project.extensions.getByType(AppEngineSpringBootExtension)
    final dotEnv = extension.loadEnvironmentOrNull()
    if (dotEnv) {
      println("Using environment variables in ${extension.environment}")
      final root = new XmlParser().parse(xml)
      final envVariablesNode = root.get('env-variables').find() ?: root.appendNode('env-variables')
      dotEnv.each { k, v ->
        envVariablesNode.appendNode('env-var', [name: k, value: v])
      }
      xml.text = XmlUtil.serialize(root)
    }
  }
}
