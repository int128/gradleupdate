import domain.LatestGradle

final latestGradle = new LatestGradle()

latestGradle.checkIfNewStableVersionIsAvailable { gradleVersion ->
    defaultQueue.add(
            url: '/api/found-new-gradle/index.task.groovy',
            params: [gradle_version: gradleVersion]
    )
}
