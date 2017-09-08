package org.hidetake.gradle.appengine.spring.boot

import groovy.xml.XmlUtil
import org.gradle.api.DefaultTask
import org.gradle.api.tasks.InputFile
import org.gradle.api.tasks.OutputFile
import org.gradle.api.tasks.TaskAction

/**
 * A task to inject env-variables into appengine-web.xml
 */
class InjectAppEngineWebXml extends DefaultTask {
  @InputFile
  File dotEnv

  @OutputFile
  File appEngineWebXml

  @TaskAction
  void mutate() {
    final dotEnvProperties = DotEnv.loadOrEmpty(dotEnv)
    if (dotEnvProperties) {
      println("Using environment variables in $dotEnv")
      final root = new XmlParser().parse(appEngineWebXml)
      final envVariablesNode = root.get('env-variables').find() ?: root.appendNode('env-variables')
      dotEnvProperties.each { k, v ->
        envVariablesNode.appendNode('env-var', [name: k, value: v])
      }
      appEngineWebXml.text = XmlUtil.serialize(root)
    }
  }
}
