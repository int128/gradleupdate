package domain

import groovy.transform.Immutable

import static infrastructure.HTTPClientExceptionUtil.nullIfResourceIsNotFound

@Immutable
class GradleVersion {

    final String string

    static GradleVersion get(GHBranch branch) {
        def content = GHContent.get(branch, 'gradle/wrapper/gradle-wrapper.properties')
        new GradleVersion(parseGradleWrapperProperties(content.contentAsString))
    }

    static GradleVersion getOrNull(GHBranch branch) {
        nullIfResourceIsNotFound {
            get(branch)
        }
    }

    static String parseGradleWrapperProperties(String gradleWrapperProperties) {
        assert gradleWrapperProperties
        try {
            def m = gradleWrapperProperties =~ /distributionUrl=.+?\/gradle-(.+?)-.+?\.zip/
            assert m
            def m0 = m[0]
            assert m0 instanceof List
            assert m0.size() == 2
            m0[1] as String
        } catch (AssertionError ignore) {
            null
        }
    }

    @Override
    String toString() {
        string
    }

}
